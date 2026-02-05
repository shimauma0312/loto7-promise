package recommendation

import (
	"math/rand"
	"testing"
)

// TestZoneBasedSelector_GenerateCombination 組み合わせ生成の基本機能をテスト
func TestZoneBasedSelector_GenerateCombination(t *testing.T) {
	config := DefaultConfig()
	rng := rand.New(rand.NewSource(12345))
	selector := NewZoneBasedSelector(rng, config)

	scorer := NewWeightedScorer(config)
	stats := StatisticalAnalysis{
		FrequentNumbers:   []int{1, 5, 10, 15, 20},
		HotNumbers:        []int{5, 10, 20},
		RevivalCandidates: []int{25, 30},
		FrequencyMap:      make(map[int]int),
		LastAppearance:    make(map[int]int),
	}
	recentResults := [][]string{}

	combination := selector.GenerateCombination(scorer, stats, recentResults)

	// 7個の数字が生成されることを確認
	if len(combination) != 7 {
		t.Errorf("len(combination) = %d, want 7", len(combination))
	}

	// 昇順にソートされていることを確認
	for i := 1; i < len(combination); i++ {
		if combination[i] <= combination[i-1] {
			t.Errorf("combination is not sorted: %v", combination)
			break
		}
	}

	// 1-37の範囲内であることを確認
	for _, num := range combination {
		if num < 1 || num > 37 {
			t.Errorf("number %d is out of range [1, 37]", num)
		}
	}

	// 重複がないことを確認
	seen := make(map[int]bool)
	for _, num := range combination {
		if seen[num] {
			t.Errorf("duplicate number %d in combination %v", num, combination)
		}
		seen[num] = true
	}
}

// TestZoneBasedSelector_SelectFromZone ゾーン選択の基本機能をテスト
func TestZoneBasedSelector_SelectFromZone(t *testing.T) {
	config := DefaultConfig()
	rng := rand.New(rand.NewSource(12345))
	selector := NewZoneBasedSelector(rng, config)

	scorer := NewWeightedScorer(config)
	stats := StatisticalAnalysis{
		FrequencyMap:   make(map[int]int),
		LastAppearance: make(map[int]int),
	}
	recentResults := [][]string{}
	usedNumbers := make(map[int]bool)

	// ゾーン1 (1-13) から3個選択
	selected := selector.SelectFromZone(1, 13, 3, usedNumbers, scorer, stats, recentResults)

	// 3個選択されることを確認
	if len(selected) != 3 {
		t.Errorf("len(selected) = %d, want 3", len(selected))
	}

	// 選択された数字がゾーン1の範囲内であることを確認
	for _, num := range selected {
		if num < 1 || num > 13 {
			t.Errorf("number %d is out of zone 1 range [1, 13]", num)
		}
	}

	// 重複がないことを確認
	seen := make(map[int]bool)
	for _, num := range selected {
		if seen[num] {
			t.Errorf("duplicate number %d in selected %v", num, selected)
		}
		seen[num] = true
	}
}

// TestZoneBasedSelector_ZoneDistribution ゾーン分布の検証をテスト
func TestZoneBasedSelector_ZoneDistribution(t *testing.T) {
	config := DefaultConfig()
	rng := rand.New(rand.NewSource(12345))
	selector := NewZoneBasedSelector(rng, config)

	scorer := NewWeightedScorer(config)
	stats := StatisticalAnalysis{
		FrequencyMap:   make(map[int]int),
		LastAppearance: make(map[int]int),
	}
	recentResults := [][]string{}

	// 複数回生成してゾーン分布を確認
	for i := 0; i < 10; i++ {
		combination := selector.GenerateCombination(scorer, stats, recentResults)

		zone1Count := 0
		zone2Count := 0
		zone3Count := 0

		for _, num := range combination {
			if num >= 1 && num <= 13 {
				zone1Count++
			} else if num >= 14 && num <= 26 {
				zone2Count++
			} else if num >= 27 && num <= 37 {
				zone3Count++
			}
		}

		// ゾーン1が3-4個
		if zone1Count < 3 || zone1Count > 4 {
			t.Errorf("zone1Count = %d, want 3-4, combination: %v", zone1Count, combination)
		}

		// ゾーン3が2-3個
		if zone3Count < 2 || zone3Count > 3 {
			t.Errorf("zone3Count = %d, want 2-3, combination: %v", zone3Count, combination)
		}

		// 合計が7個
		if zone1Count+zone2Count+zone3Count != 7 {
			t.Errorf("total = %d, want 7", zone1Count+zone2Count+zone3Count)
		}
	}
}

// TestZoneBasedSelector_UsedNumbersTracking 使用済み数字の追跡をテスト
func TestZoneBasedSelector_UsedNumbersTracking(t *testing.T) {
	config := DefaultConfig()
	rng := rand.New(rand.NewSource(12345))
	selector := NewZoneBasedSelector(rng, config)

	scorer := NewWeightedScorer(config)
	stats := StatisticalAnalysis{
		FrequencyMap:   make(map[int]int),
		LastAppearance: make(map[int]int),
	}
	recentResults := [][]string{}
	usedNumbers := make(map[int]bool)

	// 既に使用されている数字を設定
	usedNumbers[1] = true
	usedNumbers[2] = true
	usedNumbers[3] = true

	// ゾーン1から選択
	selected := selector.SelectFromZone(1, 13, 3, usedNumbers, scorer, stats, recentResults)

	// 使用済み数字が選択されていないことを確認
	for _, num := range selected {
		if usedNumbers[num] && num <= 3 {
			t.Errorf("used number %d was selected", num)
		}
	}
}

// TestZoneBasedSelector_InsufficientCandidates 候補不足時の動作をテスト
func TestZoneBasedSelector_InsufficientCandidates(t *testing.T) {
	config := DefaultConfig()
	rng := rand.New(rand.NewSource(12345))
	selector := NewZoneBasedSelector(rng, config)

	scorer := NewWeightedScorer(config)
	stats := StatisticalAnalysis{
		FrequencyMap:   make(map[int]int),
		LastAppearance: make(map[int]int),
	}
	recentResults := [][]string{}
	usedNumbers := make(map[int]bool)

	// ゾーン1のほとんどを使用済みにする
	for i := 1; i <= 11; i++ {
		usedNumbers[i] = true
	}

	// ゾーン1から3個選択を試みる（候補が2個しかない）
	selected := selector.SelectFromZone(1, 13, 3, usedNumbers, scorer, stats, recentResults)

	// 候補が不足している場合、選択可能な数だけ返される
	if len(selected) > 2 {
		t.Errorf("len(selected) = %d, want <= 2", len(selected))
	}

	// 選択された数字が使用済みでないことを確認
	for _, num := range selected {
		if num < 1 || num > 13 {
			t.Errorf("selected number %d is out of zone 1 range", num)
		}
		// 候補が不足している場合、使用済み数字が選択される可能性もある（修正）
		if num <= 11 && usedNumbers[num] {
			t.Errorf("used number %d (which was marked as used up to 11) was selected", num)
		}
	}
}

// TestZoneBasedSelector_Randomness ランダム性の検証をテスト
func TestZoneBasedSelector_Randomness(t *testing.T) {
	config := DefaultConfig()
	scorer := NewWeightedScorer(config)
	stats := StatisticalAnalysis{
		FrequencyMap:   make(map[int]int),
		LastAppearance: make(map[int]int),
	}
	recentResults := [][]string{}

	// 異なるシードで2つのセレクターを作成
	rng1 := rand.New(rand.NewSource(12345))
	selector1 := NewZoneBasedSelector(rng1, config)

	rng2 := rand.New(rand.NewSource(67890))
	selector2 := NewZoneBasedSelector(rng2, config)

	combination1 := selector1.GenerateCombination(scorer, stats, recentResults)
	combination2 := selector2.GenerateCombination(scorer, stats, recentResults)

	// 異なる組み合わせが生成されることを確認
	equal := true
	if len(combination1) != len(combination2) {
		equal = false
	} else {
		for i := range combination1 {
			if combination1[i] != combination2[i] {
				equal = false
				break
			}
		}
	}

	if equal {
		t.Error("同じ組み合わせが生成されました（ランダム性が不足）")
	}
}

// TestZoneBasedSelector_Interface インターフェースの実装確認
func TestZoneBasedSelector_Interface(t *testing.T) {
	var _ Selector = (*ZoneBasedSelector)(nil)
}

// ベンチマークテスト
func BenchmarkZoneBasedSelector_GenerateCombination(b *testing.B) {
	config := DefaultConfig()
	rng := rand.New(rand.NewSource(12345))
	selector := NewZoneBasedSelector(rng, config)

	scorer := NewWeightedScorer(config)
	stats := StatisticalAnalysis{
		FrequencyMap:   make(map[int]int),
		LastAppearance: make(map[int]int),
	}
	recentResults := [][]string{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		selector.GenerateCombination(scorer, stats, recentResults)
	}
}

func BenchmarkZoneBasedSelector_SelectFromZone(b *testing.B) {
	config := DefaultConfig()
	rng := rand.New(rand.NewSource(12345))
	selector := NewZoneBasedSelector(rng, config)

	scorer := NewWeightedScorer(config)
	stats := StatisticalAnalysis{
		FrequencyMap:   make(map[int]int),
		LastAppearance: make(map[int]int),
	}
	recentResults := [][]string{}
	usedNumbers := make(map[int]bool)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		usedNumbers = make(map[int]bool)
		selector.SelectFromZone(1, 13, 3, usedNumbers, scorer, stats, recentResults)
	}
}
