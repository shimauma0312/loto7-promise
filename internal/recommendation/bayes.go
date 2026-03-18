// ベイズ推定によるペア条件付き確率の計算を提供する
//
// ペア相関系の特徴量を担当し、RandomForest の集計統計系と役割を明確に分担する。
// 「1番目が3だった場合、候補Xは何番目に偏るか」という
// P(candidate | selected) を過去の共起データから推定する。
package recommendation

import "math"

// BayesianEstimator は過去の共起データから条件付き確率 P(candidate | selected) を推定する。
//
// Beta-Binomial モデル（ラプラス平滑化）によりデータ不足時のゼロ確率問題を補正する。
type BayesianEstimator struct {
	pairCounts  map[int]map[int]int // pairCounts[a][b] = a と b が同一抽選に出現した回数
	appearances map[int]int         // appearances[a] = a が出現した抽選数
	totalDraws  int
	alpha       float64 // ラプラス平滑化係数（1.0 が均一事前分布）
}

// NewBayesianEstimator は BayesianEstimator を作成する。
//
// alpha はラプラス平滑化の強度で、1.0 が均一事前分布（Laplace smoothing）。
func NewBayesianEstimator(alpha float64) *BayesianEstimator {
	return &BayesianEstimator{
		pairCounts:  make(map[int]map[int]int),
		appearances: make(map[int]int),
		alpha:       alpha,
	}
}

// Build は過去抽選データからペア共起行列と出現回数を構築する
func (b *BayesianEstimator) Build(results [][]string) {
	b.pairCounts = make(map[int]map[int]int)
	b.appearances = make(map[int]int)
	b.totalDraws = len(results)

	for _, draw := range results {
		nums := parseDraw(draw)
		for _, n := range nums {
			b.appearances[n]++
		}
		for i, a := range nums {
			if _, ok := b.pairCounts[a]; !ok {
				b.pairCounts[a] = make(map[int]int)
			}
			for j, c := range nums {
				if i != j {
					b.pairCounts[a][c]++
				}
			}
		}
	}
}

// ConditionalProb は selected が選ばれた条件で candidate が同じ組み合わせに含まれる事後確率を返す。
//
// P(candidate | selected) をナイーブベイズの幾何平均で近似する。
// selected が空の場合は周辺確率 P(candidate) を返す。
func (b *BayesianEstimator) ConditionalProb(candidate int, selected []int) float64 {
	if b.totalDraws == 0 {
		return float64(lotoDrawCount) / float64(lotoTotalNumbers)
	}

	if len(selected) == 0 {
		freq := b.appearances[candidate]
		return (float64(freq) + b.alpha) / (float64(b.totalDraws) + 2*b.alpha)
	}

	logSum := 0.0
	validCount := 0
	for _, s := range selected {
		if s == candidate {
			continue
		}
		totalS := b.appearances[s]
		pairCount := 0
		if counts, ok := b.pairCounts[s]; ok {
			pairCount = counts[candidate]
		}
		// Beta-Binomial: P(candidate | s) = (共起数 + alpha) / (s の出現数 + 2*alpha)
		prob := (float64(pairCount) + b.alpha) / (float64(totalS) + 2*b.alpha)
		if prob > 0 {
			logSum += math.Log(prob)
			validCount++
		}
	}

	if validCount == 0 {
		freq := b.appearances[candidate]
		return (float64(freq) + b.alpha) / (float64(b.totalDraws) + 2*b.alpha)
	}
	// 幾何平均でナイーブベイズ積を安定化
	return math.Exp(logSum / float64(validCount))
}
