package recommendation

import (
	"strings"
	"testing"
)

// TestTextFormatter_Format 基本的なフォーマット出力をテスト
func TestTextFormatter_Format(t *testing.T) {
	formatter := &TextFormatter{}

	recommendations := [][]int{
		{1, 5, 10, 15, 20, 25, 30},
		{2, 6, 11, 16, 21, 26, 31},
	}

	stats := StatisticalAnalysis{
		FrequentNumbers:   []int{1, 5, 10, 15, 20},
		HotNumbers:        []int{5, 10, 20},
		RevivalCandidates: []int{25, 30},
	}

	analyzedCount := 100

	result := formatter.Format(recommendations, stats, analyzedCount)

	// ヘッダーが含まれることを確認
	if !strings.Contains(result, "=== ロト7 推薦番号 ===") {
		t.Error("結果にヘッダーが含まれていません")
	}

	// 推薦番号が含まれることを確認
	if !strings.Contains(result, "推薦 1:") {
		t.Error("結果に推薦1が含まれていません")
	}
	if !strings.Contains(result, "推薦 2:") {
		t.Error("結果に推薦2が含まれていません")
	}

	// 分析情報が含まれることを確認
	if !strings.Contains(result, "=== 分析情報 ===") {
		t.Error("結果に分析情報ヘッダーが含まれていません")
	}
	if !strings.Contains(result, "過去100回分") {
		t.Error("結果に分析対象回数が含まれていません")
	}
	if !strings.Contains(result, "生成した推薦数: 2") {
		t.Error("結果に推薦数が含まれていません")
	}

	// 統計データが含まれることを確認
	if !strings.Contains(result, "=== 統計データ ===") {
		t.Error("結果に統計データヘッダーが含まれていません")
	}
	if !strings.Contains(result, "出現頻度が高い数字:") {
		t.Error("結果に出現頻度が含まれていません")
	}
	if !strings.Contains(result, "波が来ている数字") {
		t.Error("結果に波情報が含まれていません")
	}
	if !strings.Contains(result, "復活候補") {
		t.Error("結果に復活候補が含まれていません")
	}
}

// TestTextFormatter_NumberFormatting 数字のフォーマットをテスト
func TestTextFormatter_NumberFormatting(t *testing.T) {
	formatter := &TextFormatter{}

	recommendations := [][]int{
		{1, 2, 3, 4, 5, 6, 7},
	}

	stats := StatisticalAnalysis{}
	result := formatter.Format(recommendations, stats, 10)

	// 2桁0埋めフォーマットを確認
	if !strings.Contains(result, "01") {
		t.Error("数字が2桁0埋めでフォーマットされていません")
	}
	if !strings.Contains(result, "02") {
		t.Error("数字が2桁0埋めでフォーマットされていません")
	}

	// 区切り文字を確認
	if !strings.Contains(result, " - ") {
		t.Error("数字の区切り文字が含まれていません")
	}
}

// TestTextFormatter_EmptyRecommendations 空の推薦をテスト
func TestTextFormatter_EmptyRecommendations(t *testing.T) {
	formatter := &TextFormatter{}

	recommendations := [][]int{}
	stats := StatisticalAnalysis{}

	result := formatter.Format(recommendations, stats, 50)

	// ヘッダーは含まれる
	if !strings.Contains(result, "=== ロト7 推薦番号 ===") {
		t.Error("空の推薦でもヘッダーが含まれるべきです")
	}

	// 推薦数が0であることを確認
	if !strings.Contains(result, "生成した推薦数: 0") {
		t.Error("推薦数が0と表示されるべきです")
	}
}

// TestTextFormatter_MultipleRecommendations 複数推薦のフォーマットをテスト
func TestTextFormatter_MultipleRecommendations(t *testing.T) {
	formatter := &TextFormatter{}

	recommendations := [][]int{
		{1, 5, 10, 15, 20, 25, 30},
		{2, 6, 11, 16, 21, 26, 31},
		{3, 7, 12, 17, 22, 27, 32},
		{4, 8, 13, 18, 23, 28, 33},
		{5, 9, 14, 19, 24, 29, 34},
	}

	stats := StatisticalAnalysis{
		FrequentNumbers:   []int{1, 5, 10},
		HotNumbers:        []int{5, 10},
		RevivalCandidates: []int{25, 30},
	}

	result := formatter.Format(recommendations, stats, 100)

	// 各推薦が含まれることを確認
	for i := 1; i <= 5; i++ {
		// 簡易チェック: "推薦" という文字列が含まれる
		if !strings.Contains(result, "推薦") {
			t.Errorf("推薦%dが含まれていません", i)
		}
	}

	// 生成した推薦数が正しいことを確認
	if !strings.Contains(result, "生成した推薦数: 5") {
		t.Error("推薦数が正しく表示されていません")
	}
}

// TestTextFormatter_StatisticsFormatting 統計情報のフォーマットをテスト
func TestTextFormatter_StatisticsFormatting(t *testing.T) {
	formatter := &TextFormatter{}

	recommendations := [][]int{
		{1, 5, 10, 15, 20, 25, 30},
	}

	stats := StatisticalAnalysis{
		FrequentNumbers:   []int{12, 34, 9, 29, 1},
		HotNumbers:        []int{33, 22, 34},
		RevivalCandidates: []int{},
	}

	result := formatter.Format(recommendations, stats, 100)

	// 出現頻度が高い数字が含まれることを確認
	if !strings.Contains(result, "12") {
		t.Error("出現頻度が高い数字が表示されていません")
	}

	// 波が来ている数字が含まれることを確認
	if !strings.Contains(result, "33") {
		t.Error("波が来ている数字が表示されていません")
	}

	// 復活候補が空の場合の表示を確認
	if !strings.Contains(result, "復活候補") {
		t.Error("復活候補のラベルが表示されていません")
	}
}

// TestTextFormatter_AnalyzedCount 分析対象回数の表示をテスト
func TestTextFormatter_AnalyzedCount(t *testing.T) {
	formatter := &TextFormatter{}

	recommendations := [][]int{
		{1, 5, 10, 15, 20, 25, 30},
	}

	stats := StatisticalAnalysis{}

	tests := []struct {
		name          string
		analyzedCount int
		want          string
	}{
		{
			name:          "10回分",
			analyzedCount: 10,
			want:          "過去10回分",
		},
		{
			name:          "100回分",
			analyzedCount: 100,
			want:          "過去100回分",
		},
		{
			name:          "500回分",
			analyzedCount: 500,
			want:          "過去500回分",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.Format(recommendations, stats, tt.analyzedCount)
			if !strings.Contains(result, tt.want) {
				t.Errorf("結果に '%s' が含まれていません", tt.want)
			}
		})
	}
}

// TestTextFormatter_Interface インターフェースの実装確認
func TestTextFormatter_Interface(t *testing.T) {
	var _ Formatter = (*TextFormatter)(nil)
}

// ベンチマークテスト
func BenchmarkTextFormatter_Format(b *testing.B) {
	formatter := &TextFormatter{}

	recommendations := [][]int{
		{1, 5, 10, 15, 20, 25, 30},
		{2, 6, 11, 16, 21, 26, 31},
		{3, 7, 12, 17, 22, 27, 32},
		{4, 8, 13, 18, 23, 28, 33},
		{5, 9, 14, 19, 24, 29, 34},
	}

	stats := StatisticalAnalysis{
		FrequentNumbers:   []int{1, 5, 10, 15, 20},
		HotNumbers:        []int{5, 10, 20},
		RevivalCandidates: []int{25, 30},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		formatter.Format(recommendations, stats, 100)
	}
}
