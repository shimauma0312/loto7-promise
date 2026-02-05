package recommendation

import "testing"

// TestWeightedScorer_CalculatePriority 数字の優先度計算をテスト
func TestWeightedScorer_CalculatePriority(t *testing.T) {
	config := DefaultConfig()
	scorer := NewWeightedScorer(config)

	stats := StatisticalAnalysis{
		FrequentNumbers:   []int{1, 5, 10, 15, 20},
		RareNumbers:       []int{3, 7, 33, 34, 37},
		HotNumbers:        []int{5, 10, 20},
		RevivalCandidates: []int{25, 30},
		FrequencyMap:      map[int]int{1: 20, 5: 18, 10: 17, 15: 16, 20: 15},
		LastAppearance:    map[int]int{1: 0, 5: 1, 10: 2},
	}

	recentResults := [][]string{
		{"1", "2", "3", "4", "6", "7", "8"},
		{"9", "11", "12", "13", "14", "16", "17"},
		{"18", "19", "21", "22", "23", "24", "26"},
	}

	tests := []struct {
		name         string
		num          int
		wantPositive bool // スコアがプラスであることを期待
		wantNegative bool // スコアがマイナスであることを期待
	}{
		{
			name:         "波が来ている数字（+2.0）",
			num:          5,
			wantPositive: true,  // HotNumbersに含まれるので+2.0だが、直近に出現していれば-3.0
			wantNegative: false, // このテストでは直近に出現していない前提に修正
		},
		{
			name:         "出現頻度が高い数字（+1.0）",
			num:          1,
			wantPositive: false, // 直近に出現で-3.0
			wantNegative: true,
		},
		{
			name:         "復活候補（+1.0）",
			num:          25,
			wantPositive: true,
			wantNegative: false,
		},
		{
			name:         "出現回数が少ない数字（-0.5）",
			num:          33,
			wantPositive: false,
			wantNegative: true,
		},
		{
			name:         "直近3回で出現（-3.0）",
			num:          2,
			wantPositive: false,
			wantNegative: true,
		},
		{
			name:         "普通の数字",
			num:          22,
			wantPositive: false, // 0.0付近
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

// TestWeightedScorer_HotNumberBonus 波が来ている数字のボーナスをテスト
func TestWeightedScorer_HotNumberBonus(t *testing.T) {
	config := DefaultConfig()
	scorer := NewWeightedScorer(config)

	stats := StatisticalAnalysis{
		HotNumbers: []int{10, 20, 30},
	}

	recentResults := [][]string{} // 直近出現なし

	tests := []struct {
		name      string
		num       int
		wantBonus bool
	}{
		{
			name:      "波が来ている数字",
			num:       10,
			wantBonus: true,
		},
		{
			name:      "波が来ていない数字",
			num:       15,
			wantBonus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.CalculatePriority(tt.num, stats, recentResults)

			if tt.wantBonus && score < 2.0 {
				t.Errorf("CalculatePriority(%d) = %v, want >= 2.0", tt.num, score)
			}
			if !tt.wantBonus && score >= 2.0 {
				t.Errorf("CalculatePriority(%d) = %v, want < 2.0", tt.num, score)
			}
		})
	}
}

// TestWeightedScorer_RecentAvoidPenalty 直近出現のペナルティをテスト
func TestWeightedScorer_RecentAvoidPenalty(t *testing.T) {
	config := DefaultConfig()
	config.RecentAvoidCount = 3
	scorer := NewWeightedScorer(config)

	stats := StatisticalAnalysis{}

	recentResults := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},
		{"8", "9", "10", "11", "12", "13", "14"},
		{"15", "16", "17", "18", "19", "20", "21"},
	}

	tests := []struct {
		name        string
		num         int
		wantPenalty bool
	}{
		{
			name:        "直近1回目に出現",
			num:         1,
			wantPenalty: true,
		},
		{
			name:        "直近2回目に出現",
			num:         10,
			wantPenalty: true,
		},
		{
			name:        "直近3回目に出現",
			num:         20,
			wantPenalty: true,
		},
		{
			name:        "直近3回に未出現",
			num:         25,
			wantPenalty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.CalculatePriority(tt.num, stats, recentResults)

			if tt.wantPenalty && score > -3.0 {
				t.Errorf("CalculatePriority(%d) = %v, want <= -3.0", tt.num, score)
			}
			if !tt.wantPenalty && score < 0 {
				t.Errorf("CalculatePriority(%d) = %v, want >= 0", tt.num, score)
			}
		})
	}
}

// TestWeightedScorer_RevivalCandidateBonus 復活候補のボーナスをテスト
func TestWeightedScorer_RevivalCandidateBonus(t *testing.T) {
	config := DefaultConfig()
	scorer := NewWeightedScorer(config)

	stats := StatisticalAnalysis{
		RevivalCandidates: []int{25, 30, 35},
	}

	recentResults := [][]string{} // 直近出現なし

	tests := []struct {
		name      string
		num       int
		wantBonus bool
	}{
		{
			name:      "復活候補",
			num:       25,
			wantBonus: true,
		},
		{
			name:      "復活候補でない",
			num:       10,
			wantBonus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.CalculatePriority(tt.num, stats, recentResults)

			if tt.wantBonus && score < 1.0 {
				t.Errorf("CalculatePriority(%d) = %v, want >= 1.0", tt.num, score)
			}
			if !tt.wantBonus && score >= 1.0 {
				t.Errorf("CalculatePriority(%d) = %v, want < 1.0", tt.num, score)
			}
		})
	}
}
