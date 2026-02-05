package recommendation

import (
	"math/rand"
	"sort"
)

const (
	// ゾーン選択用の定数
	Zone3SelectProbability = 0.2  // ゾーン3から3個選択する確率
	MaxAttemptsFactor      = 10   // 最大試行回数の倍数
	ScoreThreshold         = -0.5 // スコアの閾値
	MinCandidates          = 3    // 最小候補数
	BaseWeight             = 1.0  // ベースウェイト
	MinWeight              = 0.1  // 最小ウェイト
)

// 候補選択を行うインターフェース
type Selector interface {
	SelectFromZone(min, max, count int, usedNumbers map[int]bool, scorer Scorer, stats StatisticalAnalysis, recentResults [][]string) []int
	GenerateCombination(scorer Scorer, stats StatisticalAnalysis, recentResults [][]string) []int
}

// ゾーンベースで候補を選択する実装
type ZoneBasedSelector struct {
	rng    *rand.Rand
	config RecommendationConfig
}

// ゾーンベースセレクターを作成する
func NewZoneBasedSelector(rng *rand.Rand, config RecommendationConfig) *ZoneBasedSelector {
	return &ZoneBasedSelector{
		rng:    rng,
		config: config,
	}
}

// 高度な分析に基づいた推薦を生成する
func (s *ZoneBasedSelector) GenerateCombination(scorer Scorer, stats StatisticalAnalysis, recentResults [][]string) []int {
	var combination []int
	usedNumbers := make(map[int]bool)

	// ゾーン1 (1-13): 3-4個選択
	zone1Count := 3 + s.rng.Intn(2) // 3 or 4
	zone1 := s.SelectFromZone(1, 13, zone1Count, usedNumbers, scorer, stats, recentResults)
	combination = append(combination, zone1...)
	for _, num := range zone1 {
		usedNumbers[num] = true
	}

	// ゾーン3 (27-37): 2-3個選択（まれに3個）
	zone3Count := 2
	if s.rng.Float64() < Zone3SelectProbability { // 20%の確率て3個
		zone3Count = 3
	}
	zone3 := s.SelectFromZone(27, 37, zone3Count, usedNumbers, scorer, stats, recentResults)
	combination = append(combination, zone3...)
	for _, num := range zone3 {
		usedNumbers[num] = true
	}

	// ゾーン2 (14-26): 残りを埋める
	zone2Count := 7 - len(combination)
	zone2 := s.SelectFromZone(14, 26, zone2Count, usedNumbers, scorer, stats, recentResults)
	combination = append(combination, zone2...)
	for _, num := range zone2 {
		usedNumbers[num] = true
	}

	// 7個に満たない場合は補完
	if len(combination) < 7 {
		combination = s.fillRemainingNumbers(combination, usedNumbers)
	}

	sort.Ints(combination)
	return combination
}

// 指定ゾーンから数字を選択する
func (s *ZoneBasedSelector) SelectFromZone(min, max, count int, usedNumbers map[int]bool, scorer Scorer, stats StatisticalAnalysis, recentResults [][]string) []int {
	var selected []int
	candidates := s.getZoneCandidates(min, max, usedNumbers, scorer, stats, recentResults)

	// 候補が不足している場合は調整
	if len(candidates) < count {
		count = len(candidates)
	}

	attempts := 0
	maxAttempts := count * MaxAttemptsFactor

	for len(selected) < count && attempts < maxAttempts {
		attempts++

		if len(candidates) == 0 {
			break
		}

		// 重み付きランダム選択
		num := s.weightedSelectFromCandidates(candidates, usedNumbers, scorer, stats, recentResults)
		if num == 0 {
			continue
		}

		if !usedNumbers[num] {
			selected = append(selected, num)
			usedNumbers[num] = true

			// 選択済みの候補を削除
			newCandidates := []int{}
			for _, c := range candidates {
				if c != num {
					newCandidates = append(newCandidates, c)
				}
			}
			candidates = newCandidates
		}
	}

	return selected
}

// ゾーンから候補数字を取得する
func (s *ZoneBasedSelector) getZoneCandidates(min, max int, usedNumbers map[int]bool, scorer Scorer, stats StatisticalAnalysis, recentResults [][]string) []int {
	var candidates []int

	for num := min; num <= max; num++ {
		if usedNumbers[num] {
			continue
		}

		// スコア計算
		score := scorer.CalculatePriority(num, stats, recentResults)

		// スコアが一定以上の数字を候補に
		if score > ScoreThreshold { // 闾値
			candidates = append(candidates, num)
		}
	}

	// 候補が少なすぎる場合は範囲内すべてを候補に
	if len(candidates) < MinCandidates {
		candidates = []int{}
		for num := min; num <= max; num++ {
			if !usedNumbers[num] {
				candidates = append(candidates, num)
			}
		}
	}

	return candidates
}

// 候補から重み付き選択する
func (s *ZoneBasedSelector) weightedSelectFromCandidates(candidates []int, usedNumbers map[int]bool, scorer Scorer, stats StatisticalAnalysis, recentResults [][]string) int {
	if len(candidates) == 0 {
		return 0
	}

	type weightedNum struct {
		num    int
		weight float64
	}

	var weighted []weightedNum
	totalWeight := 0.0

	for _, num := range candidates {
		if usedNumbers[num] {
			continue
		}

		priority := scorer.CalculatePriority(num, stats, recentResults)
		weight := priority + BaseWeight // 最低限の重み
		if weight < MinWeight {
			weight = MinWeight
		}

		weighted = append(weighted, weightedNum{num: num, weight: weight})
		totalWeight += weight
	}

	if len(weighted) == 0 || totalWeight == 0 {
		return 0
	}

	// ルーレット選択
	randomValue := s.rng.Float64() * totalWeight
	currentWeight := 0.0

	for _, wn := range weighted {
		currentWeight += wn.weight
		if randomValue <= currentWeight {
			return wn.num
		}
	}

	// フォールバック
	if len(weighted) > 0 {
		return weighted[s.rng.Intn(len(weighted))].num
	}

	return 0
}

// 残りの数字を補完する
func (s *ZoneBasedSelector) fillRemainingNumbers(combination []int, usedNumbers map[int]bool) []int {
	for num := 1; num <= 37 && len(combination) < 7; num++ {
		if !usedNumbers[num] {
			combination = append(combination, num)
			usedNumbers[num] = true
		}
	}
	return combination
}
