package bonusNumber

import (
	"testing"
)

// TestGetSpecificBonusNumberCount_ValidInput 正常な入力値でのテスト
func TestGetSpecificBonusNumberCount_ValidInput(t *testing.T) {
	// 数字7の過去50回分の出現回数を取得
	result, err := GetSpecificBonusNumberCount(7, 50)
	
	if err != nil {
		t.Fatalf("期待: エラーなし, 実際: %v", err)
	}
	
	if result == nil {
		t.Fatal("期待: 結果が返される, 実際: nil")
	}
	
	// 基本的なフィールドの検証
	if result.Number != 7 {
		t.Errorf("期待: Number = 7, 実際: %d", result.Number)
	}
	
	if result.TotalDraws != 50 {
		t.Errorf("期待: TotalDraws = 50, 実際: %d", result.TotalDraws)
	}
	
	if result.Count < 0 {
		t.Errorf("期待: Count >= 0, 実際: %d", result.Count)
	}
	
	if result.Percentage < 0 || result.Percentage > 100 {
		t.Errorf("期待: 0 <= Percentage <= 100, 実際: %.2f", result.Percentage)
	}
}

// TestGetSpecificBonusNumberCount_InvalidNumber 無効な数字でのテスト
func TestGetSpecificBonusNumberCount_InvalidNumber(t *testing.T) {
	testCases := []struct {
		number      int
		description string
	}{
		{0, "0（下限値未満）"},
		{38, "38（上限値超過）"},
		{-1, "-1（負の値）"},
		{100, "100（範囲外）"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result, err := GetSpecificBonusNumberCount(tc.number, 10)
			
			if err == nil {
				t.Errorf("期待: エラーが発生, 実際: エラーなし")
			}
			
			if result != nil {
				t.Errorf("期待: result = nil, 実際: %v", result)
			}
		})
	}
}

// TestGetSpecificBonusNumberCount_InvalidRange 無効な範囲でのテスト
func TestGetSpecificBonusNumberCount_InvalidRange(t *testing.T) {
	testCases := []int{0, -1, -10}
	
	for _, invalidRange := range testCases {
		t.Run("", func(t *testing.T) {
			result, err := GetSpecificBonusNumberCount(7, invalidRange)
			
			if err == nil {
				t.Errorf("期待: エラーが発生（範囲: %d）, 実際: エラーなし", invalidRange)
			}
			
			if result != nil {
				t.Errorf("期待: result = nil, 実際: %v", result)
			}
		})
	}
}

// TestGetBonusNumberTrend_ValidInput 正常な入力値での傾向分析テスト
func TestGetBonusNumberTrend_ValidInput(t *testing.T) {
	results, err := GetBonusNumberTrend(30)
	
	if err != nil {
		t.Fatalf("期待: エラーなし, 実際: %v", err)
	}
	
	if results == nil {
		t.Fatal("期待: 結果が返される, 実際: nil")
	}
	
	// 1-37の全ての数字が含まれているかチェック
	if len(results) != 37 {
		t.Errorf("期待: 結果数 = 37, 実際: %d", len(results))
	}
	
	// ソートが正しく行われているかチェック（降順）
	for i := 0; i < len(results)-1; i++ {
		if results[i].Count < results[i+1].Count {
			t.Errorf("期待: ソート順が正しい, 実際: [%d]の出現回数%d > [%d]の出現回数%d", 
				i, results[i].Count, i+1, results[i+1].Count)
		}
	}
	
	// 各結果の基本検証
	for i, result := range results {
		if result.Number < 1 || result.Number > 37 {
			t.Errorf("期待: 1 <= Number <= 37（インデックス%d）, 実際: %d", i, result.Number)
		}
		
		if result.Count < 0 {
			t.Errorf("期待: Count >= 0（数字%d）, 実際: %d", result.Number, result.Count)
		}
		
		if result.TotalDraws != 30 {
			t.Errorf("期待: TotalDraws = 30（数字%d）, 実際: %d", result.Number, result.TotalDraws)
		}
		
		if result.Percentage < 0 || result.Percentage > 100 {
			t.Errorf("期待: 0 <= Percentage <= 100（数字%d）, 実際: %.2f", result.Number, result.Percentage)
		}
	}
}

// TestGetBonusNumberTrend_InvalidRange 無効な範囲での傾向分析テスト
func TestGetBonusNumberTrend_InvalidRange(t *testing.T) {
	testCases := []int{0, -1, -5}
	
	for _, invalidRange := range testCases {
		t.Run("", func(t *testing.T) {
			results, err := GetBonusNumberTrend(invalidRange)
			
			if err == nil {
				t.Errorf("期待: エラーが発生（範囲: %d）, 実際: エラーなし", invalidRange)
			}
			
			if results != nil {
				t.Errorf("期待: results = nil, 実際: %v", results)
			}
		})
	}
}

// TestBonusNumberResult_StructureValidation 構造体の基本検証
func TestBonusNumberResult_StructureValidation(t *testing.T) {
	result := &BonusNumberResult{
		Number:     15,
		Count:      8,
		TotalDraws: 100,
		Percentage: 8.0,
	}
	
	// 各フィールドのテスト
	if result.Number != 15 {
		t.Errorf("期待: Number = 15, 実際: %d", result.Number)
	}
	
	if result.Count != 8 {
		t.Errorf("期待: Count = 8, 実際: %d", result.Count)
	}
	
	if result.TotalDraws != 100 {
		t.Errorf("期待: TotalDraws = 100, 実際: %d", result.TotalDraws)
	}
	
	if result.Percentage != 8.0 {
		t.Errorf("期待: Percentage = 8.0, 実際: %.2f", result.Percentage)
	}
}
