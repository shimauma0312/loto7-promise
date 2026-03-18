// ベイズ推定・ランダムフォレスト・重み付きスコアラーのアンサンブルを提供する
//
// BayesianEstimator にはペア相関系、RandomForest には集計統計系の特徴量を分担させ、
// 互いの相関を抑えることでアンサンブル効果を最大化する。
package recommendation

import "math"

// EnsembleWeights はアンサンブルスコアリングの各モデル寄与度を保持する。
//
// ベイズ・RF のオフセットがベーススコアに加算される量を制御する。
// BayesScale=5.0 の場合、事前確率の2倍の条件付き確率を持つ数字に最大+5点を付与する。
type EnsembleWeights struct {
	BayesScale float64 // ベイズスコアのオフセットスケール
	RFScale    float64 // ランダムフォレストスコアのオフセットスケール
}

// DefaultEnsembleWeights はデフォルトのアンサンブル重みを返す
func DefaultEnsembleWeights() EnsembleWeights {
	return EnsembleWeights{
		BayesScale: 5.0,
		RFScale:    5.0,
	}
}

// EnsembleScorer は WeightedScorer に Bayes・RF オフセットを加算する Scorer 実装
type EnsembleScorer struct {
	base     *WeightedScorer
	bayes    *BayesianEstimator
	rf       *RandomForest
	weights  EnsembleWeights
	lookback int
}

// NewEnsembleScorer は EnsembleScorer を作成する
func NewEnsembleScorer(
	base *WeightedScorer,
	bayes *BayesianEstimator,
	rf *RandomForest,
	weights EnsembleWeights,
	lookback int,
) *EnsembleScorer {
	return &EnsembleScorer{
		base:     base,
		bayes:    bayes,
		rf:       rf,
		weights:  weights,
		lookback: lookback,
	}
}

// CalculatePriority はルールベーススコアに Bayes・RF の相対スコアオフセットを加算する。
//
// selected は PoolBasedSelector から渡される選択済み数字であり、
// P(candidate | selected) の条件として使用される（例: selected={3} → P(X | 1番目=3) 近似）。
// 直近1回の直接ペナルティ等の強制スコアは保持する。
func (e *EnsembleScorer) CalculatePriority(num int, position int, stats StatisticalAnalysis, recentResults [][]string, selected []int) float64 {
	baseScore := e.base.CalculatePriority(num, position, stats, recentResults, nil)

	// 強いペナルティはそのまま返す（直近1回出現等）
	if baseScore <= ScoreRecent1Penalty {
		return baseScore
	}

	priorProb := float64(lotoDrawCount) / float64(lotoTotalNumbers)

	// ベイズ条件付き確率 P(candidate | selected) と事前確率の相対差をオフセットに変換
	bayesProb := e.bayes.ConditionalProb(num, selected)
	bayesRelative := bayesProb/priorProb - 1.0
	bayesOffset := math.Max(-e.weights.BayesScale, math.Min(e.weights.BayesScale, e.weights.BayesScale*bayesRelative))

	// RF 予測確率と事前確率の相対差をオフセットに変換
	rfProb := e.rf.Predict(num, recentResults, e.lookback)
	rfRelative := rfProb/priorProb - 1.0
	rfOffset := math.Max(-e.weights.RFScale, math.Min(e.weights.RFScale, e.weights.RFScale*rfRelative))

	return baseScore + bayesOffset + rfOffset
}
