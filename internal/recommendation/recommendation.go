package recommendation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/shimauma0312/loto7-promise/internal/result"
)

// 推薦エンジンのオーケストレーター
type RecommendationEngine struct {
	config         RecommendationConfig
	analyzer       Analyzer
	baseScorer     *WeightedScorer // ルールベーススコアラー（ウォークフォワード最適化対象）
	scorer         Scorer          // アクティブスコアラー（LoadData 後は EnsembleScorer）
	selector       Selector
	validator      Validator
	formatter      Formatter
	results        [][]string
	history        CombinationHistory
	statistics     StatisticalAnalysis
	learnedWeights bool // 学習済み重みを使用中かどうか
	trainedOnDraw  int  // 学習時の最新抽箌回号（0=不明）
}

// 推薦エンジンを作成する
//
// cache/weights.json が存在する場合は自動的に学習済み重みを読み込む。
// ファイルが存在しない場合や新回抽選後は、LoadData 呼び出し時に自動学習する。
func NewRecommendationEngine(config RecommendationConfig) *RecommendationEngine {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	weights := DefaultScorerWeights()
	usingLearned := false
	trainedOnDraw := 0
	if wf, err := LoadWeights(WeightsFilePath); err == nil {
		weights = wf.Weights
		usingLearned = true
		trainedOnDraw = wf.TrainedOnDraw
	}

	analyzer := NewDefaultAnalyzer()
	baseScorer := NewWeightedScorerWithWeights(config, weights)
	selector := NewPoolBasedSelector(rng)
	formatter := &TextFormatter{}

	return &RecommendationEngine{
		config:         config,
		analyzer:       analyzer,
		baseScorer:     baseScorer,
		scorer:         baseScorer, // LoadData 前はベーススコアラーを屐用する
		selector:       selector,
		validator:      nil,
		formatter:      formatter,
		learnedWeights: usingLearned,
		trainedOnDraw:  trainedOnDraw,
		history: CombinationHistory{
			Pairs: make(map[string]int),
		},
	}
}

// データを読み込む
func (re *RecommendationEngine) LoadData(count int) error {
	re.results = FilterAbnormalResults(result.GetResult(count))
	if len(re.results) == 0 {
		return fmt.Errorf("データを取得できませんでした")
	}

	// 組み合わせ履歴を構築
	re.history = re.analyzer.BuildCombinationHistory(re.results)

	// 統計分析を実行
	re.statistics = re.analyzer.Analyze(re.results, re.config.HistoryLookback)

	// 新回の抽選後に重みが古い場合は自動再学習
	if re.shouldRetrain() {
		re.autoRetrain()
	}
	// ベイズ推定とランダムフォレストを訓練し EnsembleScorerを構築する
	bayes := NewBayesianEstimator(1.0)
	bayes.Build(re.results)
	rf := NewRandomForest(rfTreeCount, rfMaxDepth, rand.New(rand.NewSource(42)))
	rf.Train(re.results, re.config.HistoryLookback)
	re.scorer = NewEnsembleScorer(re.baseScorer, bayes, rf, DefaultEnsembleWeights(), re.config.HistoryLookback)

	// バリデータを設定（過去の組み合わせチェックを含む）
	re.validator = NewCompositeValidator(
		&ZoneDistributionValidator{},
		&OddEvenRatioValidator{},
		&TotalSumValidator{},
		&AverageValueValidator{},
		&SmoothnessValidator{},
		&CloseNumberPairValidator{},
		NewPastCombinationValidator(re.results, 100), // 直近100回分をチェック
	)

	return nil
}

// 推薦組み合わせを生成する
func (re *RecommendationEngine) GenerateRecommendations() ([][]int, error) {
	if len(re.results) == 0 {
		return nil, fmt.Errorf("データが読み込まれていません")
	}

	var recommendations [][]int
	maxAttempts := 2000

	for rec := 0; rec < re.config.MaxRecommendations; rec++ {
		for attempt := 0; attempt < maxAttempts; attempt++ {
			combination := re.selector.GenerateCombination(re.scorer, re.statistics, re.results)

			// 全制約をパスした組み合わせを採用
			if len(combination) == 7 &&
				re.validator.Validate(combination) &&
				!re.isDuplicateCombination(combination, recommendations) {
				recommendations = append(recommendations, combination)
				break
			}
		}
	}

	return recommendations, nil
}

func (re *RecommendationEngine) isDuplicateCombination(newComb []int, existing [][]int) bool {
	for _, existingComb := range existing {
		if len(newComb) != len(existingComb) {
			continue
		}

		matches := 0
		for _, num1 := range newComb {
			for _, num2 := range existingComb {
				if num1 == num2 {
					matches++
					break
				}
			}
		}

		// 5個以上一致する組み合わせは重複とみなす
		if matches >= 5 {
			return true
		}
	}
	return false
}

// 推薦結果をフォーマットする
func (re *RecommendationEngine) FormatRecommendations(recommendations [][]int) string {
	return re.formatter.Format(recommendations, re.statistics, len(re.results))
}

// 設定情報を取得する
func (re *RecommendationEngine) GetConfig() RecommendationConfig {
	return re.config
}

// 分析した抽選回数を取得する
func (re *RecommendationEngine) GetAnalyzedDrawsCount() int {
	return len(re.results)
}

// 学習済み重みを使用中かどうかを返す
func (re *RecommendationEngine) IsUsingLearnedWeights() bool {
	return re.learnedWeights
}

// NeedsWeightUpdate は新回抽選後に重みが未更新の場合 true を返す
//
// LoadData 前でも呼び出せる。
func (re *RecommendationEngine) NeedsWeightUpdate() bool {
	return re.shouldRetrain()
}

// shouldRetrain は重みの再学習が必要かどうかを返す
//
// weights.json が未作成（初回起動）の場合も true を返す。
func (re *RecommendationEngine) shouldRetrain() bool {
	if !re.learnedWeights {
		return true // weights.json 未作成 → 初回学習
	}
	return re.trainedOnDraw < result.NewNumber()
}

// autoRetrain は読み込み済みデータで軽量再最適化を実行し weights.json を更新する
func (re *RecommendationEngine) autoRetrain() {
	// walk-forward テストケース数 = len(results) - lookback
	// config.HistoryLookback と同じデータ数の場合テストケースが 0 になるため、
	// 必ず autoRetrainMinTestCases 件以上のテストケースを確保する lookback を計算する
	effectiveLookback := re.config.HistoryLookback
	if len(re.results)-effectiveLookback < autoRetrainMinTestCases {
		effectiveLookback = len(re.results) - autoRetrainMinTestCases
	}
	if effectiveLookback < autoRetrainMinTestCases {
		// データ不足でバックテスト不可能—次回に再試行
		return
	}

	backtester := NewBacktester(re.results, effectiveLookback)
	backtester.SetTrialCount(autoRetrainTrials)
	optimizer := NewHillClimbOptimizer(backtester, time.Now().UnixNano())

	current := DefaultScorerWeights()
	if wf, err := LoadWeights(WeightsFilePath); err == nil {
		current = wf.Weights
	}

	optimized := optimizer.Optimize(current, autoRetrainIterations, nil)

	newDraw := result.NewNumber()
	if err := SaveWeights(optimized, newDraw, WeightsFilePath); err == nil {
		re.baseScorer = NewWeightedScorerWithWeights(re.config, optimized)
		// re.scorer は LoadData 内で EnsembleScorer として再構築されるため更新不要
		re.trainedOnDraw = newDraw
		re.learnedWeights = true
	}
}
