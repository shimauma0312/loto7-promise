package consecutivePattern

import (
	"testing"
)

// TestFindConsecutiveGroups 連番検出ロジックのテスト
func TestFindConsecutiveGroups(t *testing.T) {
	testCases := []struct {
		name     string
		numbers  []int
		expected [][]int
	}{
		{
			name:     "2連番_単一",
			numbers:  []int{1, 2, 5, 10, 15, 20, 25},
			expected: [][]int{{1, 2}},
		},
		{
			name:     "3連番_単一",
			numbers:  []int{1, 5, 10, 11, 12, 20, 25},
			expected: [][]int{{10, 11, 12}},
		},
		{
			name:     "複数の連番",
			numbers:  []int{1, 2, 3, 10, 11, 20, 25},
			expected: [][]int{{1, 2, 3}, {10, 11}},
		},
		{
			name:     "連番なし",
			numbers:  []int{1, 5, 10, 15, 20, 25, 30},
			expected: nil,
		},
		{
			name:     "全て連番",
			numbers:  []int{1, 2, 3, 4, 5, 6, 7},
			expected: [][]int{{1, 2, 3, 4, 5, 6, 7}},
		},
		{
			name:     "空の配列",
			numbers:  []int{},
			expected: nil,
		},
		{
			name:     "要素1個",
			numbers:  []int{5},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := findConsecutiveGroups(tc.numbers)
			
			if len(result) != len(tc.expected) {
				t.Errorf("結果の配列長が異なります。期待値: %d, 実際: %d", 
					len(tc.expected), len(result))
				return
			}
			
			for i, expectedGroup := range tc.expected {
				if len(result[i]) != len(expectedGroup) {
					t.Errorf("グループ%dの長さが異なります。期待値: %d, 実際: %d", 
						i, len(expectedGroup), len(result[i]))
					continue
				}
				
				for j, expectedNum := range expectedGroup {
					if result[i][j] != expectedNum {
						t.Errorf("グループ%dの要素%dが異なります。期待値: %d, 実際: %d", 
							i, j, expectedNum, result[i][j])
					}
				}
			}
		})
	}
}

// TestFindSameDigitGroups ゾロ目検出ロジックのテスト
func TestFindSameDigitGroups(t *testing.T) {
	testCases := []struct {
		name     string
		numbers  []int
		expected [][]int
	}{
		{
			name:     "一の位1のゾロ目",
			numbers:  []int{1, 11, 21, 5, 10, 15, 25},
			expected: [][]int{{1, 11, 21}},
		},
		{
			name:     "複数の一の位でゾロ目",
			numbers:  []int{1, 11, 2, 12, 5, 15, 25},
			expected: [][]int{{1, 11}, {2, 12}, {5, 15, 25}},
		},
		{
			name:     "ゾロ目なし",
			numbers:  []int{1, 5, 10, 16, 22, 28, 34},
			expected: nil,
		},
		{
			name:     "全て同じ一の位",
			numbers:  []int{1, 11, 21, 31},
			expected: [][]int{{1, 11, 21, 31}},
		},
		{
			name:     "空の配列",
			numbers:  []int{},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := findSameDigitGroups(tc.numbers)
			
			// 結果とexpectedの長さをチェック
			if len(result) != len(tc.expected) {
				t.Errorf("結果の配列長が異なります。期待値: %d, 実際: %d", 
					len(tc.expected), len(result))
				return
			}
			
			// 順序が異なる可能性があるため、各期待値グループが結果に含まれているかチェック
			for _, expectedGroup := range tc.expected {
				found := false
				for _, resultGroup := range result {
					if slicesEqual(expectedGroup, resultGroup) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("期待されるグループ %v が結果に含まれていません", expectedGroup)
				}
			}
		})
	}
}

// TestGeneratePatternKey パターンキー生成のテスト
func TestGeneratePatternKey(t *testing.T) {
	testCases := []struct {
		name        string
		pattern     []int
		patternType string
		expected    string
	}{
		{
			name:        "連番パターン",
			pattern:     []int{1, 2, 3},
			patternType: "consecutive",
			expected:    "consecutive_1_2_3",
		},
		{
			name:        "ゾロ目パターン",
			pattern:     []int{1, 11, 21},
			patternType: "same_digit",
			expected:    "same_digit_1_11_21",
		},
		{
			name:        "単一要素",
			pattern:     []int{5},
			patternType: "test",
			expected:    "test_5",
		},
		{
			name:        "空パターン",
			pattern:     []int{},
			patternType: "empty",
			expected:    "empty_",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := generatePatternKey(tc.pattern, tc.patternType)
			if result != tc.expected {
				t.Errorf("パターンキーが異なります。期待値: %s, 実際: %s", 
					tc.expected, result)
			}
		})
	}
}

// TestAnalyzeConsecutivePatterns_InvalidInput 不正な入力のテスト
func TestAnalyzeConsecutivePatterns_InvalidInput(t *testing.T) {
	testCases := []struct {
		name        string
		searchRange int
		expectError bool
	}{
		{
			name:        "検索範囲が0",
			searchRange: 0,
			expectError: true,
		},
		{
			name:        "検索範囲が負数",
			searchRange: -1,
			expectError: true,
		},
		{
			name:        "正常な入力",
			searchRange: 1,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := AnalyzeConsecutivePatterns(tc.searchRange)
			
			if tc.expectError && err == nil {
				t.Errorf("エラーが期待されたのに発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーが発生しました: %v", err)
			}
		})
	}
}

// BenchmarkAnalyzeConsecutivePatterns パフォーマンステスト
func BenchmarkAnalyzeConsecutivePatterns(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := AnalyzeConsecutivePatterns(10)
		if err != nil {
			b.Fatalf("ベンチマークでエラー: %v", err)
		}
	}
}

// ヘルパー関数: スライスが等しいかチェック
func slicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	
	// ソートされた状態で比較
	aCopy := make([]int, len(a))
	bCopy := make([]int, len(b))
	copy(aCopy, a)
	copy(bCopy, b)
	
	// 簡易ソート（バブルソート）
	for i := 0; i < len(aCopy); i++ {
		for j := 0; j < len(aCopy)-i-1; j++ {
			if aCopy[j] > aCopy[j+1] {
				aCopy[j], aCopy[j+1] = aCopy[j+1], aCopy[j]
			}
		}
	}
	
	for i := 0; i < len(bCopy); i++ {
		for j := 0; j < len(bCopy)-i-1; j++ {
			if bCopy[j] > bCopy[j+1] {
				bCopy[j], bCopy[j+1] = bCopy[j+1], bCopy[j]
			}
		}
	}
	
	for i := range aCopy {
		if aCopy[i] != bCopy[i] {
			return false
		}
	}
	
	return true
}
