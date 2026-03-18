package test

import (
	"strings"
	"testing"

	"github.com/shimauma0312/loto7-promise/internal/recommendation"
)

func TestTextFormatter_Format_ContainsHeader(t *testing.T) {
	formatter := &recommendation.TextFormatter{}

	recs := [][]int{
		{1, 5, 10, 15, 20, 25, 30},
		{2, 6, 11, 16, 21, 26, 31},
	}
	stats := recommendation.StatisticalAnalysis{
		FrequentNumbers:   []int{1, 5, 10, 15, 20},
		HotNumbers:        []int{5, 10, 20},
		RevivalCandidates: []int{25, 30},
	}

	result := formatter.Format(recs, stats, 100)

	checks := []struct {
		substr string
		desc   string
	}{
		{"=== ロト7 推薦番号 ===", "ヘッダー"},
		{"推薦 1:", "推薦1"},
		{"推薦 2:", "推薦2"},
		{"=== 分析情報 ===", "分析情報ヘッダー"},
		{"過去100回分", "分析対象回数"},
		{"生成した推薦数: 2", "推薦数"},
		{"=== 統計データ ===", "統計データヘッダー"},
		{"出現頻度が高い数字:", "頻度情報"},
		{"波が来ている数字", "波情報"},
		{"復活候補", "復活候補"},
	}

	for _, c := range checks {
		if !strings.Contains(result, c.substr) {
			t.Errorf("結果に%sが含まれていません（期待: %q）", c.desc, c.substr)
		}
	}
}

func TestTextFormatter_NumberFormatting_ZeroPadded(t *testing.T) {
	formatter := &recommendation.TextFormatter{}

	result := formatter.Format([][]int{{1, 2, 3, 4, 5, 6, 7}}, recommendation.StatisticalAnalysis{}, 10)

	if !strings.Contains(result, "01") {
		t.Error("数字が2桁0埋めでフォーマットされていません")
	}
}
