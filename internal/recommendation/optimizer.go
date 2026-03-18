// walk-forward バックテストによるスコアリング重みの最適化機能を提供する
package recommendation

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// WeightsFilePath は学習済み重みの保存・読み込みパス（cache ディレクトリに統一）
const WeightsFilePath = "cache/weights.json"

// autoRetrainIterations は新回検出時の自動再最適化の反復回数
const autoRetrainIterations = 100

// autoRetrainTrials は自動再最適化の評価試行数
const autoRetrainTrials = 5

// autoRetrainMinTestCases は walk-forward で最低限必要なテストケース数
// この数が確保できない場合は自動再学習をスキップする
const autoRetrainMinTestCases = 15

// WeightsFile は重みと学習メタデータをまとめたファイル構造
type WeightsFile struct {
	Weights       ScorerWeights `json:"weights"`
	TrainedOnDraw int           `json:"trained_on_draw"` // 学習時の最新抽選回号
	TrainedAt     string        `json:"trained_at"`      // 学習日時（RFC3339）
}

// Backtester は walk-forward 方式で ScorerWeights の性能を評価する
//
// 評価指標: 各テスト回でtrialCount個の推薦を生成し、
// 実際の当選番号と最もよく一致した組み合わせのヒット数（0-7）を計上した平均値。
type Backtester struct {
	results    [][]string
	lookback   int
	trialCount int
}

// バックテスターを作成する
func NewBacktester(results [][]string, lookback int) *Backtester {
	return &Backtester{
		results:    results,
		lookback:   lookback,
		trialCount: 10,
	}
}

// 1 回の評価で生成する推薦数を設定する
func (b *Backtester) SetTrialCount(n int) {
	if n > 0 {
		b.trialCount = n
	}
}

// 指定した重みの平均ヒット数（0.0〜7.0）を返す
func (b *Backtester) Evaluate(weights ScorerWeights, rng *rand.Rand) float64 {
	if len(b.results) <= b.lookback {
		return 0
	}

	analyzer := NewDefaultAnalyzer()
	totalHits := 0
	testCount := 0

	for i := 0; i+b.lookback < len(b.results); i++ {
		target := b.results[i]
		train := b.results[i+1 : i+1+b.lookback]

		stats := analyzer.Analyze(train, b.lookback)
		scorer := NewWeightedScorerWithWeights(DefaultConfig(), weights)
		selector := NewPoolBasedSelector(rng)

		bestHits := 0
		for t := 0; t < b.trialCount; t++ {
			comb := selector.GenerateCombination(scorer, stats, train)
			if h := countHits(comb, target); h > bestHits {
				bestHits = h
			}
		}
		totalHits += bestHits
		testCount++
	}

	if testCount == 0 {
		return 0
	}
	return float64(totalHits) / float64(testCount)
}

// countHits は生成した組み合わせと当選番号の一致数を返す
func countHits(combination []int, draw []string) int {
	drawSet := make(map[int]bool, len(draw))
	for _, s := range draw {
		if n, err := strconv.Atoi(s); err == nil {
			drawSet[n] = true
		}
	}
	hits := 0
	for _, n := range combination {
		if drawSet[n] {
			hits++
		}
	}
	return hits
}

// HillClimbOptimizer はランダムヒルクライミングで ScorerWeights を最適化する
//
// 各反復でランダムに重みを摂動し、評価スコアが改善する場合のみ採用する。
// 評価は固定シードの RNG を使うことでノイズを抑制し比較を安定させる。
type HillClimbOptimizer struct {
	backtester *Backtester
	rng        *rand.Rand
}

// ヒルクライミングオプティマイザを作成する
func NewHillClimbOptimizer(backtester *Backtester, seed int64) *HillClimbOptimizer {
	return &HillClimbOptimizer{
		backtester: backtester,
		rng:        rand.New(rand.NewSource(seed)),
	}
}

// iterations 回のランダムヒルクライミングを実行し、最良の重みを返す
//
// progress が nil でなければ各反復後に呼ばれる（進捗表示に使用）。
func (o *HillClimbOptimizer) Optimize(initial ScorerWeights, iterations int, progress func(i int, bestScore float64)) ScorerWeights {
	newEvalRNG := func() *rand.Rand { return rand.New(rand.NewSource(42)) }

	best := initial
	bestScore := o.backtester.Evaluate(best, newEvalRNG())

	for i := 0; i < iterations; i++ {
		candidate := perturbWeights(best, o.rng)
		score := o.backtester.Evaluate(candidate, newEvalRNG())
		if score > bestScore {
			best = candidate
			bestScore = score
		}
		if progress != nil {
			progress(i, bestScore)
		}
	}
	return best
}

// SaveWeights は重みと学習回号を JSON ファイルに保存する
func SaveWeights(weights ScorerWeights, trainedOnDraw int, path string) error {
	wf := WeightsFile{
		Weights:       weights,
		TrainedOnDraw: trainedOnDraw,
		TrainedAt:     time.Now().Format(time.RFC3339),
	}
	data, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadWeights は JSON ファイルから WeightsFile を読み込む
//
// 旧フォーマット（ScorerWeights 直接）の場合は TrainedOnDraw=0 で返す（→次回起動時に自動再学習）。
func LoadWeights(path string) (WeightsFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return WeightsFile{}, err
	}
	var wf WeightsFile
	if err := json.Unmarshal(data, &wf); err != nil {
		return WeightsFile{}, err
	}
	// 旧フォーマット: "weights" キーが存在せず Weights が全ゼロの場合は直接パース
	if wf.Weights.Recent1Penalty == 0 && wf.Weights.Recent2Penalty == 0 {
		var legacy ScorerWeights
		if err := json.Unmarshal(data, &legacy); err == nil && legacy.Recent1Penalty != 0 {
			wf.Weights = legacy
			wf.TrainedOnDraw = 0 // 不明なので 0 → 次回起動時に再学習
		}
	}
	return wf, nil
}


// perturbWeights は各パラメータに ±20% のランダムノイズを加えた新しい ScorerWeights を返す
func perturbWeights(w ScorerWeights, rng *rand.Rand) ScorerWeights {
	noise := func(v float64) float64 {
		return v * (1.0 + (rng.Float64()*2-1)*0.2)
	}
	clampNeg := func(v float64) float64 {
		if v > -0.01 {
			return -0.01
		}
		return v
	}
	clampPos := func(v float64) float64 {
		if v < 0 {
			return 0
		}
		return v
	}
	clampScale := func(v float64) float64 {
		switch {
		case v < 0.1:
			return 0.1
		case v > 10.0:
			return 10.0
		default:
			return v
		}
	}

	newScales := [7]float64{}
	for i := 0; i < 7; i++ {
		newScales[i] = clampScale(noise(w.PositionPenaltyScale[i]))
	}

	return ScorerWeights{
		Recent1Penalty:          clampNeg(noise(w.Recent1Penalty)),
		Recent2Penalty:          clampNeg(noise(w.Recent2Penalty)),
		Within50Boost:           clampPos(noise(w.Within50Boost)),
		LowFrequencyPenalty:     clampNeg(noise(w.LowFrequencyPenalty)),
		NoAppearancePenalty:     w.NoAppearancePenalty, // 未出現は固定
		PositionPenaltyScale:    newScales,
		PositionFrequencyWeight: clampPos(noise(w.PositionFrequencyWeight)),
	}
}
