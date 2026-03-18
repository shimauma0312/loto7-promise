package recommendation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/shimauma0312/loto7-promise/internal/result"
)

// 推薦エンジンのオーケストレーター
type RecommendationEngine struct {
	config     RecommendationConfig
	analyzer   Analyzer
	scorer     Scorer
	selector   Selector
	validator  Validator
	formatter  Formatter
	results    [][]string
	history    CombinationHistory
	statistics StatisticalAnalysis
}

// 推薦エンジンを作成する
func NewRecommendationEngine(config RecommendationConfig) *RecommendationEngine {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	analyzer := NewDefaultAnalyzer()
	scorer := NewWeightedScorer(config)
	selector := NewPoolBasedSelector(rng)
	formatter := &TextFormatter{}

	return &RecommendationEngine{
		config:    config,
		analyzer:  analyzer,
		scorer:    scorer,
		selector:  selector,
		validator: nil,
		formatter: formatter,
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
