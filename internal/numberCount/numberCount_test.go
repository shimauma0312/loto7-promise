package numberCount

import (
	"testing"
)

// TestNumberCountResult_Validation 結果構造体の基本的な検証
func TestNumberCountResult_Validation(t *testing.T) {
	result := &NumberCountResult{
		Number:     7,
		Count:      15,
		TotalDraws: 100,
		Percentage: 15.0,
	}

	// 各フィールドテスト
	if result.Number != 7 {
		t.Errorf("Number = %d, want 7", result.Number)
	}
	if result.Count != 15 {
		t.Errorf("Count = %d, want 15", result.Count)
	}
	if result.TotalDraws != 100 {
		t.Errorf("TotalDraws = %d, want 100", result.TotalDraws)
	}
	if result.Percentage != 15.0 {
		t.Errorf("Percentage = %f, want 15.0", result.Percentage)
	}
}

// TestCountSpecificNumber_InvalidInput 不正な入力テスト
func TestCountSpecificNumber_InvalidInput(t *testing.T) {
	testCases := []struct {
		name         string
		targetNumber int
		searchRange  int
		expectError  bool
	}{
		{
			name:         "数字が範囲外_小さすぎる",
			targetNumber: 0,
			searchRange:  10,
			expectError:  true,
		},
		{
			name:         "数字が範囲外_大きすぎる",
			targetNumber: 38,
			searchRange:  10,
			expectError:  true,
		},
		{
			name:         "検索範囲が0",
			targetNumber: 7,
			searchRange:  0,
			expectError:  true,
		},
		{
			name:         "検索範囲が負数",
			targetNumber: 7,
			searchRange:  -1,
			expectError:  true,
		},
		{
			name:         "正常な入力",
			targetNumber: 7,
			searchRange:  1,
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := CountSpecificNumber(tc.targetNumber, tc.searchRange)
			
			if tc.expectError && err == nil {
				t.Errorf("エラーが期待されたのに発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーが発生しました: %v", err)
			}
		})
	}
}

// TestCountMultipleNumbers_EmptyInput 空の入力でのテスト
func TestCountMultipleNumbers_EmptyInput(t *testing.T) {
	numbers := []int{}
	_, err := CountMultipleNumbers(numbers, 10)
	
	if err == nil {
		t.Error("空の配列でエラーが発生するべきです")
	}
}

// TestCountMultipleNumbers_MixedValid 一部正常・一部異常な入力のテスト
func TestCountMultipleNumbers_MixedValid(t *testing.T) {
	// 正常な数字と異常な数字を混在させる
	numbers := []int{7, 0, 23, 38, 15} // 0と38は範囲外
	results, err := CountMultipleNumbers(numbers, 1)
	
	// 有効な結果が取得できるかチェック
	if err != nil && len(results) == 0 {
		t.Error("有効な数字があるのに結果が取得できませんでした")
	}
	
	// 有効な数字（7, 23, 15）の結果が含まれているかチェック
	expectedValidNumbers := map[int]bool{7: false, 23: false, 15: false}
	for _, result := range results {
		if _, exists := expectedValidNumbers[result.Number]; exists {
			expectedValidNumbers[result.Number] = true
		}
	}
	
	for number, found := range expectedValidNumbers {
		if !found {
			t.Errorf("有効な数字 %d の結果が見つかりませんでした", number)
		}
	}
}

// BenchmarkCountSpecificNumber パフォーマンステスト
func BenchmarkCountSpecificNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := CountSpecificNumber(7, 10)
		if err != nil {
			b.Fatalf("ベンチマークでエラー: %v", err)
		}
	}
}
