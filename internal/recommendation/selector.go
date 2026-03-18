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
// 各選択ラウンドで残候補に対して動的ポジション（選択済み数字のうち num より小さいものの数）を
// 用いて再スコアリングすることでポジション推定の精度を高める。
func (s *PoolBasedSelector) GenerateCombination(scorer Scorer, stats StatisticalAnalysis, recentResults [][]string) []int {
	pool := make([]int, lotoTotalNumbers)
	weights := make([]float64, lotoTotalNumbers)
	for i := 0; i < lotoTotalNumbers; i++ {
		pool[i] = i + 1
	}

	combination := make([]int, 0, lotoDrawCount)
	size := lotoTotalNumbers

	for len(combination) < lotoDrawCount {
		// 残残候補を現在の選択済みセットを元に動的ポジションで再スコアリング
		for i := 0; i < size; i++ {
			num := pool[i]
			pos := dynamicPosition(num, combination)
			priority := scorer.CalculatePriority(num, pos, stats, recentResults, combination)
			weight := priority + BaseWeight
			if weight < MinWeight {
				weight = MinWeight
			}
			weights[i] = weight
		}

		chosen := weightedRouletteSelect(s.rng, pool[:size], weights[:size])

		// pool から除去（末尾と交換）
		for i := 0; i < size; i++ {
			if pool[i] == chosen {
				size--
				pool[i] = pool[size]
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

// num が最終的に占めるポジション（0-indexed）を確率的に推定する。
//
// selected が空のとき pos=0 を返す単純計算では全候補がポジション0で評価されてしまい、
// PositionFrequencyMap の pos=0 スコアが小さい数字（特に1）に有利に働く偏りが生じる。
// そこで「残り選択数 × 残プールで num より小さくなりうる数字の割合」を期待値として加算し、
// 選択が進むにつれて推定ポジションが実際の最終位置に収束するようにする。
func dynamicPosition(num int, selected []int) int {
	alreadyBelow := 0
	for _, s := range selected {
		if s < num {
			alreadyBelow++
		}
	}
	remainingPicks := lotoDrawCount - len(selected) - 1
	remainingPool := lotoTotalNumbers - len(selected) - 1
	if remainingPool <= 0 || remainingPicks <= 0 {
		pos := alreadyBelow
		if pos >= lotoDrawCount {
			pos = lotoDrawCount - 1
		}
		return pos
	}
	potentiallyBelow := (num - 1) - alreadyBelow
	if potentiallyBelow < 0 {
		potentiallyBelow = 0
	}
	expectedPos := float64(alreadyBelow) + float64(remainingPicks)*float64(potentiallyBelow)/float64(remainingPool)
	pos := int(expectedPos + 0.5) // 四捨五入
	if pos >= lotoDrawCount {
		pos = lotoDrawCount - 1
	}
	return pos
}
