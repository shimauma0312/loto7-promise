package test

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/shimauma0312/loto7-promise/internal/recommendation"
)

func TestPoolBasedSelector_Returns7Numbers(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))
	selector := recommendation.NewPoolBasedSelector(rng)
	scorer := recommendation.NewWeightedScorer(recommendation.DefaultConfig())
	stats := recommendation.StatisticalAnalysis{
		FrequencyMap:   make(map[int]int),
		LastAppearance: make(map[int]int),
	}

	combination := selector.GenerateCombination(scorer, stats, [][]string{})

	if len(combination) != 7 {
		t.Errorf("len(combination) = %d, want 7", len(combination))
	}
}

func TestPoolBasedSelector_NoDuplicates(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))
	selector := recommendation.NewPoolBasedSelector(rng)
	scorer := recommendation.NewWeightedScorer(recommendation.DefaultConfig())
	stats := recommendation.StatisticalAnalysis{
		FrequencyMap:   make(map[int]int),
		LastAppearance: make(map[int]int),
	}

	for trial := 0; trial < 50; trial++ {
		comb := selector.GenerateCombination(scorer, stats, [][]string{})
		seen := make(map[int]bool)
		for _, n := range comb {
			if seen[n] {
				t.Errorf("trial %d: duplicate %d in %v", trial, n, comb)
			}
			seen[n] = true
		}
	}
}

func TestPoolBasedSelector_NumbersInRange(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))
	selector := recommendation.NewPoolBasedSelector(rng)
	scorer := recommendation.NewWeightedScorer(recommendation.DefaultConfig())
	stats := recommendation.StatisticalAnalysis{
		FrequencyMap:   make(map[int]int),
		LastAppearance: make(map[int]int),
	}

	for trial := 0; trial < 50; trial++ {
		comb := selector.GenerateCombination(scorer, stats, [][]string{})
		for _, n := range comb {
			if n < 1 || n > 37 {
				t.Errorf("trial %d: number %d out of range [1, 37]", trial, n)
			}
		}
	}
}

func TestPoolBasedSelector_IsSorted(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))
	selector := recommendation.NewPoolBasedSelector(rng)
	scorer := recommendation.NewWeightedScorer(recommendation.DefaultConfig())
	stats := recommendation.StatisticalAnalysis{
		FrequencyMap:   make(map[int]int),
		LastAppearance: make(map[int]int),
	}

	for trial := 0; trial < 50; trial++ {
		comb := selector.GenerateCombination(scorer, stats, [][]string{})
		if !sort.IntsAreSorted(comb) {
			t.Errorf("trial %d: combination not sorted: %v", trial, comb)
		}
	}
}

func TestPoolBasedSelector_ScoredWeightsAffectSelection(t *testing.T) {
	// 数字1だけ高頻度にして、1が選ばれやすいことを確認
	rng := rand.New(rand.NewSource(42))
	selector := recommendation.NewPoolBasedSelector(rng)
	scorer := recommendation.NewWeightedScorer(recommendation.DefaultConfig())

	freq := make(map[int]int)
	for i := 1; i <= 37; i++ {
		freq[i] = 0 // no appearance → ScoreNoAppearancePenalty
	}
	freq[1] = 50
	freq[2] = 50
	freq[3] = 50
	freq[4] = 50
	freq[5] = 50
	freq[6] = 50
	freq[7] = 50

	stats := recommendation.StatisticalAnalysis{
		FrequencyMap:   freq,
		LastAppearance: make(map[int]int),
	}

	hit := 0
	trials := 200
	for i := 0; i < trials; i++ {
		comb := selector.GenerateCombination(scorer, stats, [][]string{})
		for _, n := range comb {
			if n >= 1 && n <= 7 {
				hit++
			}
		}
	}

	// 1-7 が 7 枠中ほぼ独占するはず（期待: hit >= trials*5 ≈ 1000）
	if hit < trials*4 {
		t.Errorf("高スコア数字の選択率が低すぎます: hit=%d / %d trials", hit, trials)
	}
}
