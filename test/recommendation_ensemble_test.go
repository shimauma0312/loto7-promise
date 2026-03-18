package test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/shimauma0312/loto7-promise/internal/recommendation"
)

func buildEnsembleScorer(results [][]string, lookback int) *recommendation.EnsembleScorer {
	config := recommendation.DefaultConfig()
	base := recommendation.NewWeightedScorer(config)

	bayes := recommendation.NewBayesianEstimator(1.0)
	bayes.Build(results)

	rng := rand.New(rand.NewSource(42))
	rf := recommendation.NewRandomForest(5, 3, rng)
	rf.Train(results, lookback)

	return recommendation.NewEnsembleScorer(base, bayes, rf, recommendation.DefaultEnsembleWeights(), lookback)
}

func TestEnsembleScorer_CalculatePriority_FiniteValue(t *testing.T) {
	results := buildRFTestResults(30, 77)
	scorer := buildEnsembleScorer(results, 20)
	stats := recommendation.StatisticalAnalysis{
		FrequencyMap: map[int]int{},
	}
	for num := 1; num <= 37; num++ {
		stats.FrequencyMap[num] = 10
	}

	for num := 1; num <= 37; num++ {
		score := scorer.CalculatePriority(num, 0, stats, results, nil)
		if math.IsNaN(score) || math.IsInf(score, 0) {
			t.Errorf("CalculatePriority(%d) = %v, want finite", num, score)
		}
	}
}

func TestEnsembleScorer_CalculatePriority_WithSelected_DiffersFromUnselected(t *testing.T) {
	results := buildRFTestResults(50, 55)
	scorer := buildEnsembleScorer(results, 30)

	stats := recommendation.StatisticalAnalysis{
		FrequencyMap: map[int]int{},
	}
	for num := 1; num <= 37; num++ {
		stats.FrequencyMap[num] = 15
	}

	// selected がある場合とない場合でスコアに差が出る数字が存在するはず
	differentCount := 0
	for num := 1; num <= 37; num++ {
		scoreNoSelected := scorer.CalculatePriority(num, 0, stats, results, nil)
		scoreWithSelected := scorer.CalculatePriority(num, 0, stats, results, []int{5, 10})
		if math.Abs(scoreNoSelected-scoreWithSelected) > 1e-9 {
			differentCount++
		}
	}
	if differentCount == 0 {
		t.Error("selected あり/なしでスコアに差がない: Bayesian 条件付き確率が機能していない可能性")
	}
}

func TestEnsembleScorer_CalculatePriority_PreservesHardPenalty(t *testing.T) {
	results := buildRFTestResults(30, 77)
	scorer := buildEnsembleScorer(results, 20)

	// 直近1回目に出現した数字は強制的に ScoreRecent1Penalty を返すはず
	stats := recommendation.StatisticalAnalysis{
		FrequencyMap: map[int]int{3: 10},
		PositionLastAppearance: []map[int]int{
			{3: 0}, // position 0 で直近1回目に出現
		},
	}

	score := scorer.CalculatePriority(3, 0, stats, results, nil)
	if score != recommendation.ScoreRecent1Penalty {
		t.Errorf("directly-penalized score = %v, want %v", score, recommendation.ScoreRecent1Penalty)
	}
}
