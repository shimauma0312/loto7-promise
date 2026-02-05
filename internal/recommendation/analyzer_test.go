package recommendation

import (
	"reflect"
	"testing"
)

// TestDefaultAnalyzer_Analyze 統計分析の基本機能をテスト
func TestDefaultAnalyzer_Analyze(t *testing.T) {
	analyzer := NewDefaultAnalyzer()

	// テスト用の過去データ
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

	lookback := 10
	stats := analyzer.Analyze(results, lookback)

	// FrequencyMapの存在確認
	if stats.FrequencyMap == nil {
		t.Fatal("FrequencyMap is nil")
	}

	// LastAppearanceの存在確認
	if stats.LastAppearance == nil {
		t.Fatal("LastAppearance is nil")
	}

	// FrequentNumbersが10個であることを確認
	if len(stats.FrequentNumbers) != 10 {
		t.Errorf("len(FrequentNumbers) = %d, want 10", len(stats.FrequentNumbers))
	}

	// RareNumbersが10個であることを確認
	if len(stats.RareNumbers) != 10 {
		t.Errorf("len(RareNumbers) = %d, want 10", len(stats.RareNumbers))
	}
}

// TestDefaultAnalyzer_FrequencyMap 出現回数のカウントをテスト
func TestDefaultAnalyzer_FrequencyMap(t *testing.T) {
	analyzer := NewDefaultAnalyzer()

	results := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},
		{"1", "2", "3", "4", "5", "6", "8"},
		{"1", "2", "3", "4", "5", "9", "10"},
	}

	stats := analyzer.Analyze(results, 3)

	// 数字1は3回出現
	if stats.FrequencyMap[1] != 3 {
		t.Errorf("FrequencyMap[1] = %d, want 3", stats.FrequencyMap[1])
	}

	// 数字5は3回出現
	if stats.FrequencyMap[5] != 3 {
		t.Errorf("FrequencyMap[5] = %d, want 3", stats.FrequencyMap[5])
	}

	// 数字7は1回出現
	if stats.FrequencyMap[7] != 1 {
		t.Errorf("FrequencyMap[7] = %d, want 1", stats.FrequencyMap[7])
	}

	// 数字11は0回出現
	if stats.FrequencyMap[11] != 0 {
		t.Errorf("FrequencyMap[11] = %d, want 0", stats.FrequencyMap[11])
	}
}

// TestDefaultAnalyzer_HotNumbers 波が来ている数字の検出をテスト
func TestDefaultAnalyzer_HotNumbers(t *testing.T) {
	analyzer := NewDefaultAnalyzer()

	// 直近10回で数字5が4回出現
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

	// 数字5がHotNumbersに含まれることを確認（4回出現 >= 3回）
	found := false
	for _, num := range stats.HotNumbers {
		if num == 5 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("HotNumbers does not contain 5, got %v", stats.HotNumbers)
	}
}

// TestDefaultAnalyzer_RevivalCandidates 復活候補の検出をテスト
func TestDefaultAnalyzer_RevivalCandidates(t *testing.T) {
	analyzer := NewDefaultAnalyzer()

	// 数字30が20回以上未出現になるようなテストデータを作成
	// 最初の20回は数字30を含まない結果
	results := make([][]string, 21)
	for i := 0; i < 20; i++ {
		results[i] = []string{"1", "2", "3", "4", "5", "6", "7"}
	}
	// 21回目に数字30を含める（LastAppearance[30] = 20）
	results[20] = []string{"30", "31", "32", "33", "34", "35", "36"}

	stats := analyzer.Analyze(results, 21)

	// 数字37は一度も出現していないため、復活候補に含まれるはず
	// ただし、RareNumbersの上位5個に含まれている場合は除外される
	found37 := false
	for _, num := range stats.RevivalCandidates {
		if num == 37 {
			found37 = true
			break
		}
	}

	// 37がRareNumbersの上位5個に含まれているかチェック
	isRare37 := false
	for i := 0; i < 5 && i < len(stats.RareNumbers); i++ {
		if stats.RareNumbers[i] == 37 {
			isRare37 = true
			break
		}
	}

	// 37がRareNumbersに含まれていなければ、RevivalCandidatesに含まれるべき
	if !isRare37 && !found37 {
		t.Errorf("数字37は復活候補に含まれるべきですが、含まれていません: %v", stats.RevivalCandidates)
	}

	t.Logf("RevivalCandidates: %v", stats.RevivalCandidates)
	t.Logf("RareNumbers (top 5): %v", stats.RareNumbers[:5])
}

// TestDefaultAnalyzer_BuildCombinationHistory ペア履歴の構築をテスト
func TestDefaultAnalyzer_BuildCombinationHistory(t *testing.T) {
	analyzer := NewDefaultAnalyzer()

	results := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},
		{"1", "2", "8", "9", "10", "11", "12"},
	}

	history := analyzer.BuildCombinationHistory(results)

	// Pairsが存在することを確認
	if history.Pairs == nil {
		t.Fatal("Pairs is nil")
	}

	// ペア"1,2"が2回出現
	if history.Pairs["1,2"] != 2 {
		t.Errorf("Pairs[\"1,2\"] = %d, want 2", history.Pairs["1,2"])
	}

	// ペア"1,3"が1回出現
	if history.Pairs["1,3"] != 1 {
		t.Errorf("Pairs[\"1,3\"] = %d, want 1", history.Pairs["1,3"])
	}

	// ペア"1,8"が1回出現
	if history.Pairs["1,8"] != 1 {
		t.Errorf("Pairs[\"1,8\"] = %d, want 1", history.Pairs["1,8"])
	}
}

// TestDefaultAnalyzer_LookbackLimit lookbackの上限チェックをテスト
func TestDefaultAnalyzer_LookbackLimit(t *testing.T) {
	analyzer := NewDefaultAnalyzer()

	results := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},
		{"8", "9", "10", "11", "12", "13", "14"},
		{"15", "16", "17", "18", "19", "20", "21"},
	}

	// lookbackが結果数より大きい場合、結果数に調整される
	stats := analyzer.Analyze(results, 100)

	// FrequencyMapに値が設定されていることを確認
	if len(stats.FrequencyMap) == 0 {
		t.Error("FrequencyMap is empty")
	}

	// 数字1が1回出現
	if stats.FrequencyMap[1] != 1 {
		t.Errorf("FrequencyMap[1] = %d, want 1", stats.FrequencyMap[1])
	}
}

// TestDefaultAnalyzer_EmptyResults 空の結果をテスト
func TestDefaultAnalyzer_EmptyResults(t *testing.T) {
	analyzer := NewDefaultAnalyzer()

	results := [][]string{}

	stats := analyzer.Analyze(results, 10)

	// FrequencyMapが空であることを確認
	if len(stats.FrequencyMap) != 0 {
		t.Errorf("len(FrequencyMap) = %d, want 0", len(stats.FrequencyMap))
	}

	// HotNumbersが空であることを確認
	if len(stats.HotNumbers) != 0 {
		t.Errorf("len(HotNumbers) = %d, want 0", len(stats.HotNumbers))
	}
}

// TestDefaultAnalyzer_FrequentNumbers 出現頻度が高い数字の抽出をテスト
func TestDefaultAnalyzer_FrequentNumbers(t *testing.T) {
	analyzer := NewDefaultAnalyzer()

	// 数字1が最も多く出現（5回）
	results := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},
		{"1", "8", "9", "10", "11", "12", "13"},
		{"1", "14", "15", "16", "17", "18", "19"},
		{"1", "20", "21", "22", "23", "24", "25"},
		{"1", "26", "27", "28", "29", "30", "31"},
	}

	stats := analyzer.Analyze(results, 5)

	// FrequentNumbersの最初の要素が1であることを確認
	if len(stats.FrequentNumbers) == 0 {
		t.Fatal("FrequentNumbers is empty")
	}
	if stats.FrequentNumbers[0] != 1 {
		t.Errorf("FrequentNumbers[0] = %d, want 1", stats.FrequentNumbers[0])
	}
}

// TestDefaultAnalyzer_RareNumbers 出現頻度が低い数字の抽出をテスト
func TestDefaultAnalyzer_RareNumbers(t *testing.T) {
	analyzer := NewDefaultAnalyzer()

	// 数字37は一度も出現しない
	results := [][]string{
		{"1", "2", "3", "4", "5", "6", "7"},
		{"1", "2", "3", "4", "5", "6", "8"},
		{"1", "2", "3", "4", "5", "6", "9"},
	}

	stats := analyzer.Analyze(results, 3)

	// RareNumbersに37が含まれることを確認
	found := false
	for _, num := range stats.RareNumbers {
		if num == 37 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("RareNumbers does not contain 37, got %v", stats.RareNumbers)
	}
}

// TestDefaultAnalyzer_Interface インターフェースの実装確認
func TestDefaultAnalyzer_Interface(t *testing.T) {
	var _ Analyzer = (*DefaultAnalyzer)(nil)
}

// ベンチマークテスト
func BenchmarkDefaultAnalyzer_Analyze(b *testing.B) {
	analyzer := NewDefaultAnalyzer()

	results := make([][]string, 100)
	for i := 0; i < 100; i++ {
		results[i] = []string{"1", "2", "3", "4", "5", "6", "7"}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.Analyze(results, 100)
	}
}

func BenchmarkDefaultAnalyzer_BuildCombinationHistory(b *testing.B) {
	analyzer := NewDefaultAnalyzer()

	results := make([][]string, 100)
	for i := 0; i < 100; i++ {
		results[i] = []string{"1", "2", "3", "4", "5", "6", "7"}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.BuildCombinationHistory(results)
	}
}

// テストヘルパー関数
func contains(slice []int, target int) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

func assertEqual(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
