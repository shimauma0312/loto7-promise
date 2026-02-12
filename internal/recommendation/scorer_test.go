package recommendation

import "testing"

// TestWeightedScorer_CalculatePriority 数字の優先度計算をテスト
func TestWeightedScorer_CalculatePriority(t *testing.T) {
	config := DefaultConfig()
	scorer := NewWeightedScorer(config)

	stats := StatisticalAnalysis{
		FrequencyMap: map[int]int{
			1: 8, 5: 18, 10: 17, 15: 16, 20: 15, // 5回以上
			3: 3, 7: 2, 33: 4, 34: 2, 37: 1, // 5回未満
			25: 10, 30: 12,
		},
	}

	recentResults := [][]string{
		{"1", "2", "3", "4", "6", "7", "8"},        // 直近1回目
		{"9", "11", "12", "13", "14", "16", "17"},  // 直近2回目
		{"18", "19", "21", "22", "23", "24", "26"}, // 直近3回目
		{"5", "10", "15", "20", "25", "30", "35"},  // 4回目（50回以内カウント対象）
		{"5", "10", "27", "28", "29", "31", "36"},  // 5回目（50回以内カウント対象）
	}

	tests := []struct {
		name         string
		num          int
		wantPositive bool
		wantNegative bool
	}{
		{
			name:         "直近1回で出現した数字（限りなく低い）",
			num:          1,
			wantPositive: false,
			wantNegative: true,
		},
		{
			name:         "直近2回で出現した数字（それなりに低い）",
			num:          9,
			wantPositive: false,
			wantNegative: true,
		},
		{
			name:         "50回以内に複数回出現した数字（ブースト）",
			num:          5,
			wantPositive: true,
			wantNegative: false,
		},
		{
			name:         "5回未満の出現回数（低頻度ペナルティ）",
			num:          33,
			wantPositive: false,
			wantNegative: true,
		},
		{
			name:         "出現していない数字（出現率0）",
			num:          36,
			wantPositive: false,
			wantNegative: true,
		},
		{
			name:         "普通の数字（5回以上、直近には出ていない）",
			num:          25,
			wantPositive: true, // 4回目に出現しているのでブースト
			wantNegative: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.CalculatePriority(tt.num, stats, recentResults)

			if tt.wantPositive && score <= 0 {
				t.Errorf("CalculatePriority(%d) = %v, want positive", tt.num, score)
			}
			if tt.wantNegative && score >= 0 {
				t.Errorf("CalculatePriority(%d) = %v, want negative", tt.num, score)
			}
		})
	}
}

// TestWeightedScorer_Recent1Penalty 直近1回のペナルティをテスト
func TestWeightedScorer_Recent1Penalty(t *testing.T) {
	config := DefaultConfig()
	scorer := NewWeightedScorer(config)

	stats := StatisticalAnalysis{
		FrequencyMap: map[int]int{1: 10, 2: 8},
	}

	recentResults := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"}, // 直近1回目
	}

	score1 := scorer.CalculatePriority(1, stats, recentResults)
	if score1 != ScoreRecent1Penalty {
		t.Errorf("直近1回の数字のスコア = %v, want %v", score1, ScoreRecent1Penalty)
	}
}

// TestWeightedScorer_Recent2Penalty 直近2回のペナルティをテスト
func TestWeightedScorer_Recent2Penalty(t *testing.T) {
	config := DefaultConfig()
	scorer := NewWeightedScorer(config)

	stats := StatisticalAnalysis{
		FrequencyMap: map[int]int{8: 10},
	}

	recentResults := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},      // 直近1回目
		{"8", "9", "10", "11", "12", "13", "14"}, // 直近2回目
	}

	score := scorer.CalculatePriority(8, stats, recentResults)
	if score >= 0 {
		t.Errorf("直近2回の数字のスコア = %v, want negative", score)
	}
}

// TestWeightedScorer_Within50Boost 50回以内ブーストをテスト
func TestWeightedScorer_Within50Boost(t *testing.T) {
	config := DefaultConfig()
	scorer := NewWeightedScorer(config)

	stats := StatisticalAnalysis{
		FrequencyMap: map[int]int{5: 10},
	}

	recentResults := [][]string{
		{"1", "2", "3", "4", "6", "7", "8"},        // 直近1回目
		{"9", "10", "11", "12", "13", "14", "15"},  // 直近2回目
		{"16", "17", "18", "19", "20", "21", "22"}, // 直近3回目
		{"5", "23", "24", "25", "26", "27", "28"},  // 4回目（ブースト対象）
		{"5", "29", "30", "31", "32", "33", "34"},  // 5回目（ブースト対象）
	}

	score := scorer.CalculatePriority(5, stats, recentResults)
	expectedBoost := 2 * ScoreWithin50Boost // 2回出現
	if score < expectedBoost {
		t.Errorf("50回以内に2回出現した数字のスコア = %v, want >= %v", score, expectedBoost)
	}
}

// TestWeightedScorer_NoAppearancePenalty 未出現ペナルティをテスト
func TestWeightedScorer_NoAppearancePenalty(t *testing.T) {
	config := DefaultConfig()
	scorer := NewWeightedScorer(config)

	stats := StatisticalAnalysis{
		FrequencyMap: map[int]int{}, // 出現なし
	}

	recentResults := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},
	}

	score := scorer.CalculatePriority(99, stats, recentResults)
	if score != ScoreNoAppearancePenalty {
		t.Errorf("未出現の数字のスコア = %v, want %v", score, ScoreNoAppearancePenalty)
	}
}

// TestWeightedScorer_LowFrequencyPenalty 低頻度ペナルティをテスト
func TestWeightedScorer_LowFrequencyPenalty(t *testing.T) {
	config := DefaultConfig()
	scorer := NewWeightedScorer(config)

	stats := StatisticalAnalysis{
		FrequencyMap: map[int]int{10: 3}, // 5回未満
	}

	recentResults := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},
	}

	score := scorer.CalculatePriority(10, stats, recentResults)
	if score >= 0 {
		t.Errorf("低頻度の数字のスコア = %v, want negative", score)
	}
}
