package result

import (
	"reflect"
	"testing"
)

// TestFilterAbnormalResults_RemovesAbnormalDraws 変異体除去の基本動作をテスト
func TestFilterAbnormalResults_RemovesAbnormalDraws(t *testing.T) {
	tests := []struct {
		name  string
		input [][]string
		want  [][]string
	}{
		{
			name: "29以上が2個 — 残す",
			input: [][]string{
				{"03", "08", "14", "19", "22", "29", "35"},
			},
			want: [][]string{
				{"03", "08", "14", "19", "22", "29", "35"},
			},
		},
		{
			name: "29以上が3個 — 残す",
			input: [][]string{
				{"01", "05", "10", "15", "29", "33", "37"},
			},
			want: [][]string{
				{"01", "05", "10", "15", "29", "33", "37"},
			},
		},
		{
			name: "29以上が1個 — 除去",
			input: [][]string{
				{"01", "05", "10", "15", "20", "25", "35"},
			},
			want: [][]string{},
		},
		{
			name: "29以上が0個 — 除去",
			input: [][]string{
				{"01", "02", "03", "04", "05", "06", "07"},
			},
			want: [][]string{},
		},
		{
			name: "境界値: 28は対象外、29は対象内",
			input: [][]string{
				{"01", "05", "10", "15", "20", "28", "37"}, // 29以上は37のみ → 除去
				{"01", "05", "10", "15", "20", "29", "37"}, // 29以上は29,37 → 残す
			},
			want: [][]string{
				{"01", "05", "10", "15", "20", "29", "37"},
			},
		},
		{
			name: "混在: 正常2件・異常1件",
			input: [][]string{
				{"03", "08", "14", "19", "22", "29", "35"}, // 正常
				{"01", "05", "10", "15", "20", "25", "30"}, // 29以上が30のみ → 除去
				{"02", "07", "11", "20", "29", "32", "36"}, // 正常
			},
			want: [][]string{
				{"03", "08", "14", "19", "22", "29", "35"},
				{"02", "07", "11", "20", "29", "32", "36"},
			},
		},
		{
			name:  "空のスライス — そのまま返す",
			input: [][]string{},
			want:  [][]string{},
		},
		{
			name: "すべて異常 — 空を返す",
			input: [][]string{
				{"01", "02", "03", "04", "05", "06", "07"},
				{"01", "05", "10", "15", "20", "25", "28"},
			},
			want: [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterAbnormalResults(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterAbnormalResults() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCountHighNumbers_CountsCorrectly countHighNumbers の単体テスト
func TestCountHighNumbers_CountsCorrectly(t *testing.T) {
	tests := []struct {
		name string
		draw []string
		want int
	}{
		{"すべて29未満", []string{"01", "10", "20", "28"}, 0},
		{"29が1個", []string{"01", "10", "20", "29"}, 1},
		{"29と37の2個", []string{"01", "10", "29", "37"}, 2},
		{"すべて29以上", []string{"29", "30", "31", "37"}, 4},
		{"不正フォーマットを含む", []string{"xx", "29", "37"}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countHighNumbers(tt.draw)
			if got != tt.want {
				t.Errorf("countHighNumbers() = %d, want %d", got, tt.want)
			}
		})
	}
}
