package recommendation

import (
	"math/rand"
	"sort"
)

const (
	BaseWeight = 1.0 // スコアのベース重み
	MinWeight  = 0.1 // 最小重み
)

const (
	lotoTotalNumbers = 37
	lotoDrawCount    = 7
)

// 組み合わせ生成を行うインターフェース
type Selector interface {
	GenerateCombination(scorer Scorer, stats StatisticalAnalysis, recentResults [][]string) []int
}

// 統計スコアを重みとしたプール非復元抽出で候補を選択する実装
type PoolBasedSelector struct {
	rng *rand.Rand
}

// 現在時刻シード済みの RNG を受け取り PoolBasedSelector を返す
func NewPoolBasedSelector(rng *rand.Rand) *PoolBasedSelector {
	return &PoolBasedSelector{rng: rng}
}

// 1-37 全数字をプールとし、統計スコアを重みとした非復元ルーレット選択で7個を生成する。
//
// 各数字のポジションは値ベースで推定する（数字の大きさ → ソート後の位置を線形マッピング）ため、
// 順次充填方式の選択バイアスを回避しつつポジション単位のスコアリングを実現する。
func (s *PoolBasedSelector) GenerateCombination(scorer Scorer, stats StatisticalAnalysis, recentResults [][]string) []int {
	pool := make([]int, lotoTotalNumbers)
	weights := make([]float64, lotoTotalNumbers)
	for i := 0; i < lotoTotalNumbers; i++ {
		num := i + 1
		pool[i] = num
		pos := estimatedPosition(num)
		priority := scorer.CalculatePriority(num, pos, stats, recentResults)
		weight := priority + BaseWeight
		if weight < MinWeight {
			weight = MinWeight
		}
		weights[i] = weight
	}

	combination := make([]int, 0, lotoDrawCount)
	size := lotoTotalNumbers

	for len(combination) < lotoDrawCount {
		chosen := weightedRouletteSelect(s.rng, pool[:size], weights[:size])

		// pool から除去（末尾と交換）
		for i := 0; i < size; i++ {
			if pool[i] == chosen {
				size--
				pool[i] = pool[size]
				weights[i] = weights[size]
				break
			}
		}
		combination = append(combination, chosen)
	}

	sort.Ints(combination)
	return combination
}

// 重み付きルーレット選択で candidates の中から1つを選ぶ
func weightedRouletteSelect(rng *rand.Rand, candidates []int, weights []float64) int {
	total := 0.0
	for _, w := range weights {
		total += w
	}
	r := rng.Float64() * total
	cumulative := 0.0
	for i, w := range weights {
		cumulative += w
		if r < cumulative {
			return candidates[i]
		}
	}
	return candidates[len(candidates)-1]
}

// 数字の値から「ソート後に何番目になるか」を線形推定する（0-indexed）
//
// 全37数字を7ポジションに均等マッピング。例: num=1-6→pos0, 7-11→pos1, ...
func estimatedPosition(num int) int {
	pos := (num - 1) * lotoDrawCount / lotoTotalNumbers
	if pos >= lotoDrawCount {
		pos = lotoDrawCount - 1
	}
	return pos
}
