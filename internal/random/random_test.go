package random

import (
	"sort"
	"testing"
)

// engineWithFrequency はテスト用に頻度マップを直接セットしたエンジンを返す
func engineWithFrequency(freq map[int]int) *RandomEngine {
	e := NewRandomEngine()
	e.frequency = freq
	return e
}

// defaultFreq は 1-37 全数字に均等な頻度を返す
func defaultFreq() map[int]int {
	freq := make(map[int]int, TotalNumbers)
	for i := 1; i <= TotalNumbers; i++ {
		freq[i] = 10
	}
	return freq
}

// TestGenerateRandomCombination_Returns7Numbers 正常に7個生成されることを確認
func TestGenerateRandomCombination_Returns7Numbers(t *testing.T) {
	engine := engineWithFrequency(defaultFreq())
	comb, err := engine.GenerateRandomCombination()
	if err != nil {
		t.Fatalf("GenerateRandomCombination() error: %v", err)
	}
	if len(comb) != DrawCount {
		t.Errorf("len(combination) = %d, want %d", len(comb), DrawCount)
	}
}

// TestGenerateRandomCombination_NoDuplicates 重複数字が含まれないことを確認
func TestGenerateRandomCombination_NoDuplicates(t *testing.T) {
	engine := engineWithFrequency(defaultFreq())
	for trial := 0; trial < 100; trial++ {
		comb, err := engine.GenerateRandomCombination()
		if err != nil {
			t.Fatalf("trial %d: error: %v", trial, err)
		}
		seen := make(map[int]bool)
		for _, n := range comb {
			if seen[n] {
				t.Errorf("trial %d: duplicate number %d in %v", trial, n, comb)
			}
			seen[n] = true
		}
	}
}

// TestGenerateRandomCombination_NumbersInRange 全数字が 1-37 に収まることを確認
func TestGenerateRandomCombination_NumbersInRange(t *testing.T) {
	engine := engineWithFrequency(defaultFreq())
	for trial := 0; trial < 100; trial++ {
		comb, err := engine.GenerateRandomCombination()
		if err != nil {
			t.Fatalf("trial %d: error: %v", trial, err)
		}
		for _, n := range comb {
			if n < 1 || n > TotalNumbers {
				t.Errorf("trial %d: number %d out of range [1, %d]", trial, n, TotalNumbers)
			}
		}
	}
}

// TestGenerateRandomCombination_IsSorted 結果が昇順ソート済みであることを確認
func TestGenerateRandomCombination_IsSorted(t *testing.T) {
	engine := engineWithFrequency(defaultFreq())
	for trial := 0; trial < 50; trial++ {
		comb, err := engine.GenerateRandomCombination()
		if err != nil {
			t.Fatalf("trial %d: error: %v", trial, err)
		}
		if !sort.IntsAreSorted(comb) {
			t.Errorf("trial %d: combination not sorted: %v", trial, comb)
		}
	}
}

// TestGenerateRandomCombination_NoDataReturnsError データ未ロード時にエラーになることを確認
func TestGenerateRandomCombination_NoDataReturnsError(t *testing.T) {
	engine := NewRandomEngine()
	_, err := engine.GenerateRandomCombination()
	if err == nil {
		t.Error("expected error for empty frequency map, got nil")
	}
}

// TestAnalyzeFrequency_CountsCorrectly 頻度集計が正確であることを確認
func TestAnalyzeFrequency_CountsCorrectly(t *testing.T) {
	engine := NewRandomEngine()
	engine.results = [][]string{
		{"01", "05", "10", "15", "29", "33", "37"},
		{"01", "05", "10", "20", "29", "33", "37"},
	}
	engine.analyzeFrequency()

	tests := []struct {
		num  int
		want int
	}{
		{1, 2},
		{5, 2},
		{10, 2},
		{15, 1},
		{20, 1},
		{29, 2},
		{33, 2},
		{37, 2},
	}
	for _, tt := range tests {
		if got := engine.frequency[tt.num]; got != tt.want {
			t.Errorf("frequency[%d] = %d, want %d", tt.num, got, tt.want)
		}
	}
}

// TestGetNumberStats_SortedByFrequency 出現回数降順でソートされることを確認
func TestGetNumberStats_SortedByFrequency(t *testing.T) {
	freq := make(map[int]int)
	freq[1] = 5
	freq[2] = 10
	freq[3] = 3
	engine := engineWithFrequency(freq)
	stats := engine.GetNumberStats()

	if stats[0].Number != 2 || stats[0].Frequency != 10 {
		t.Errorf("stats[0] = {%d, %d}, want {2, 10}", stats[0].Number, stats[0].Frequency)
	}
}

// TestGenerateMultipleRandomCombinations_CountMatches 指定数通りに生成できることを確認
func TestGenerateMultipleRandomCombinations_CountMatches(t *testing.T) {
	engine := engineWithFrequency(defaultFreq())
	combs, err := engine.GenerateMultipleRandomCombinations(5)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(combs) != 5 {
		t.Errorf("len = %d, want 5", len(combs))
	}
}

// TestGenerateMultipleRandomCombinations_NoDuplicateSets 重複する組み合わせセットが含まれないことを確認
func TestGenerateMultipleRandomCombinations_NoDuplicateSets(t *testing.T) {
	engine := engineWithFrequency(defaultFreq())
	combs, err := engine.GenerateMultipleRandomCombinations(10)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	keys := make(map[string]bool)
	for _, c := range combs {
		k := combinationKey(c)
		if keys[k] {
			t.Errorf("duplicate combination: %v", c)
		}
		keys[k] = true
	}
}
