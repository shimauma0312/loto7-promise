package test

import (
	"testing"

	"github.com/shimauma0312/loto7-promise/internal/recommendation"
)

func TestDefaultAnalyzer_Analyze_BasicStats(t *testing.T) {
	analyzer := recommendation.NewDefaultAnalyzer()

	results := [][]string{
		{"1", "5", "10", "15", "20", "25", "30"},
		{"1", "5", "10", "16", "21", "26", "31"},
		{"1", "5", "11", "17", "22", "27", "32"},
		{"2", "6", "12", "18", "23", "28", "33"},
		{"2", "6", "12", "18", "23", "28", "33"},
		{"3", "7", "13", "19", "24", "29", "34"},
		{"3", "7", "13", "19", "24", "29", "34"},
		{"4", "8", "14", "20", "25", "30", "35"},
		{"4", "8", "14", "20", "25", "30", "35"},
		{"5", "9", "15", "21", "26", "31", "36"},
	}

	stats := analyzer.Analyze(results, 10)

	if stats.FrequencyMap == nil {
		t.Fatal("FrequencyMap is nil")
	}
	if stats.LastAppearance == nil {
		t.Fatal("LastAppearance is nil")
	}
	if len(stats.FrequentNumbers) != 10 {
		t.Errorf("len(FrequentNumbers) = %d, want 10", len(stats.FrequentNumbers))
	}
	if len(stats.RareNumbers) != 10 {
		t.Errorf("len(RareNumbers) = %d, want 10", len(stats.RareNumbers))
	}
}

func TestDefaultAnalyzer_FrequencyMap_CountsCorrectly(t *testing.T) {
	analyzer := recommendation.NewDefaultAnalyzer()

	results := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},
		{"1", "2", "3", "4", "5", "6", "8"},
		{"1", "2", "3", "4", "5", "9", "10"},
	}

	stats := analyzer.Analyze(results, 3)

	tests := []struct{ num, want int }{
		{1, 3}, {5, 3}, {7, 1}, {11, 0},
	}
	for _, tt := range tests {
		if got := stats.FrequencyMap[tt.num]; got != tt.want {
			t.Errorf("FrequencyMap[%d] = %d, want %d", tt.num, got, tt.want)
		}
	}
}

func TestDefaultAnalyzer_HotNumbers_Detected(t *testing.T) {
	analyzer := recommendation.NewDefaultAnalyzer()

	results := [][]string{
		{"5", "10", "15", "20", "25", "30", "35"},
		{"5", "11", "16", "21", "26", "31", "36"},
		{"1", "12", "17", "22", "27", "32", "37"},
		{"5", "13", "18", "23", "28", "33", "34"},
		{"2", "14", "19", "24", "29", "30", "31"},
		{"5", "15", "20", "25", "26", "27", "28"},
		{"3", "16", "21", "22", "23", "24", "25"},
		{"4", "17", "18", "19", "20", "21", "22"},
		{"6", "10", "11", "12", "13", "14", "15"},
		{"7", "8", "9", "16", "17", "18", "19"},
	}

	stats := analyzer.Analyze(results, 10)

	found := false
	for _, num := range stats.HotNumbers {
		if num == 5 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("HotNumbers does not contain 5 (appeared 4 times), got %v", stats.HotNumbers)
	}
}

func TestDefaultAnalyzer_LookbackLimit_Clamped(t *testing.T) {
	analyzer := recommendation.NewDefaultAnalyzer()

	results := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},
		{"8", "9", "10", "11", "12", "13", "14"},
	}

	stats := analyzer.Analyze(results, 100)

	if len(stats.FrequencyMap) == 0 {
		t.Error("FrequencyMap is empty")
	}
	if stats.FrequencyMap[1] != 1 {
		t.Errorf("FrequencyMap[1] = %d, want 1", stats.FrequencyMap[1])
	}
}

func TestDefaultAnalyzer_EmptyResults(t *testing.T) {
	analyzer := recommendation.NewDefaultAnalyzer()
	stats := analyzer.Analyze([][]string{}, 10)

	if len(stats.FrequencyMap) != 0 {
		t.Errorf("FrequencyMap should be empty, got %v", stats.FrequencyMap)
	}
	if len(stats.HotNumbers) != 0 {
		t.Errorf("HotNumbers should be empty, got %v", stats.HotNumbers)
	}
}

func TestDefaultAnalyzer_FrequentNumbers_TopIsCorrect(t *testing.T) {
	analyzer := recommendation.NewDefaultAnalyzer()

	results := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},
		{"1", "8", "9", "10", "11", "12", "13"},
		{"1", "14", "15", "16", "17", "18", "19"},
		{"1", "20", "21", "22", "23", "24", "25"},
		{"1", "26", "27", "28", "29", "30", "31"},
	}

	stats := analyzer.Analyze(results, 5)

	if len(stats.FrequentNumbers) == 0 {
		t.Fatal("FrequentNumbers is empty")
	}
	if stats.FrequentNumbers[0] != 1 {
		t.Errorf("FrequentNumbers[0] = %d, want 1", stats.FrequentNumbers[0])
	}
}

func TestDefaultAnalyzer_BuildCombinationHistory(t *testing.T) {
	analyzer := recommendation.NewDefaultAnalyzer()

	results := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},
		{"1", "2", "8", "9", "10", "11", "12"},
	}

	history := analyzer.BuildCombinationHistory(results)

	if history.Pairs == nil {
		t.Fatal("Pairs is nil")
	}
	if history.Pairs["1,2"] != 2 {
		t.Errorf("Pairs[\"1,2\"] = %d, want 2", history.Pairs["1,2"])
	}
	if history.Pairs["1,3"] != 1 {
		t.Errorf("Pairs[\"1,3\"] = %d, want 1", history.Pairs["1,3"])
	}
}
