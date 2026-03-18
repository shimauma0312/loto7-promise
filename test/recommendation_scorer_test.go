package test

import (
	"testing"

	"github.com/shimauma0312/loto7-promise/internal/recommendation"
)

func TestWeightedScorer_CalculatePriority(t *testing.T) {
	config := recommendation.DefaultConfig()
	scorer := recommendation.NewWeightedScorer(config)

	// PositionLastAppearance / PositionFrequencyMap は未設定 → グローバルフォールバックで動作
	stats := recommendation.StatisticalAnalysis{
		FrequencyMap: map[int]int{
			1: 8, 5: 18, 10: 17, 15: 16, 20: 15,
			3: 3, 7: 2, 33: 4, 34: 2, 37: 1,
			25: 10, 30: 12,
		},
	}

	recentResults := [][]string{
		{"1", "2", "3", "4", "6", "7", "8"},
		{"9", "11", "12", "13", "14", "16", "17"},
		{"18", "19", "21", "22", "23", "24", "26"},
		{"5", "10", "15", "20", "25", "30", "35"},
		{"5", "10", "27", "28", "29", "31", "36"},
	}

	tests := []struct {
		name         string
		num          int
		wantPositive bool
		wantNegative bool
	}{
		{"直近1回で出現した数字", 1, false, true},
		{"直近2回で出現した数字", 9, false, true},
		{"50回以内に複数回出現した数字（ブースト）", 5, true, false},
		{"5回未満の出現回数（低頻度ペナルティ）", 33, false, true},
		{"出現していない数字", 36, false, true},
		{"普通の数字（5回以上・直近未出現）", 25, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.CalculatePriority(tt.num, 0, stats, recentResults)
			if tt.wantPositive && score <= 0 {
				t.Errorf("CalculatePriority(%d) = %v, want positive", tt.num, score)
			}
			if tt.wantNegative && score >= 0 {
				t.Errorf("CalculatePriority(%d) = %v, want negative", tt.num, score)
			}
		})
	}
}

func TestWeightedScorer_Recent1ReturnsExactPenalty(t *testing.T) {
	scorer := recommendation.NewWeightedScorer(recommendation.DefaultConfig())
	// 数字1がポジション0（最小値枠）で直近1回目に出現
	stats := recommendation.StatisticalAnalysis{
		FrequencyMap: map[int]int{1: 10},
		PositionLastAppearance: []map[int]int{
			{1: 0}, // position 0: 直近1回目に出現
		},
	}
	recentResults := [][]string{{"1", "2", "3", "4", "5", "6", "7"}}

	score := scorer.CalculatePriority(1, 0, stats, recentResults)
	if score != recommendation.ScoreRecent1Penalty {
		t.Errorf("score = %v, want %v", score, recommendation.ScoreRecent1Penalty)
	}
}

func TestWeightedScorer_Recent2IsNegative(t *testing.T) {
	scorer := recommendation.NewWeightedScorer(recommendation.DefaultConfig())
	// 数字8がポジション0で直近2回目に出現
	stats := recommendation.StatisticalAnalysis{
		FrequencyMap: map[int]int{8: 10},
		PositionLastAppearance: []map[int]int{
			{8: 1}, // position 0: 直近2回目（インデックス1）に出現
		},
	}
	recentResults := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},
		{"8", "9", "10", "11", "12", "13", "14"},
	}

	score := scorer.CalculatePriority(8, 0, stats, recentResults)
	if score >= 0 {
		t.Errorf("直近2回の数字のスコア = %v, want negative", score)
	}
}

func TestWeightedScorer_Within50Boost(t *testing.T) {
	scorer := recommendation.NewWeightedScorer(recommendation.DefaultConfig())
	stats := recommendation.StatisticalAnalysis{
		FrequencyMap: map[int]int{5: 10},
	}
	recentResults := [][]string{
		{"1", "2", "3", "4", "6", "7", "8"},
		{"9", "10", "11", "12", "13", "14", "15"},
		{"16", "17", "18", "19", "20", "21", "22"},
		{"5", "23", "24", "25", "26", "27", "28"},
		{"5", "29", "30", "31", "32", "33", "34"},
	}

	score := scorer.CalculatePriority(5, 0, stats, recentResults)
	expected := 2 * recommendation.ScoreWithin50Boost
	if score < expected {
		t.Errorf("score = %v, want >= %v (2回分ブースト)", score, expected)
	}
}

func TestWeightedScorer_NoAppearanceReturnsExactPenalty(t *testing.T) {
	scorer := recommendation.NewWeightedScorer(recommendation.DefaultConfig())
	stats := recommendation.StatisticalAnalysis{FrequencyMap: map[int]int{}}
	recentResults := [][]string{{"1", "2", "3", "4", "5", "6", "7"}}

	score := scorer.CalculatePriority(37, 6, stats, recentResults)
	if score != recommendation.ScoreNoAppearancePenalty {
		t.Errorf("score = %v, want %v", score, recommendation.ScoreNoAppearancePenalty)
	}
}

func TestWeightedScorer_LowFrequencyIsNegative(t *testing.T) {
	scorer := recommendation.NewWeightedScorer(recommendation.DefaultConfig())
	stats := recommendation.StatisticalAnalysis{
		FrequencyMap: map[int]int{10: 3},
	}
	recentResults := [][]string{{"1", "2", "3", "4", "5", "6", "7"}}

	score := scorer.CalculatePriority(10, 2, stats, recentResults)
	if score >= 0 {
		t.Errorf("低頻度の数字のスコア = %v, want negative", score)
	}
}
