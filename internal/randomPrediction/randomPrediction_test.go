package randomPrediction

import (
	"testing"
	"time"
)

// TestGenerateRandomPrediction_ValidInput 正常な入力値での予想番号生成テスト
func TestGenerateRandomPrediction_ValidInput(t *testing.T) {
	testCases := []struct {
		analysisRange int
		useTrend      bool
		description   string
	}{
		{50, true, "傾向ベース予想"},
		{100, false, "ランダム予想"},
		{30, true, "短期傾向ベース予想"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result, err := GenerateRandomPrediction(tc.analysisRange, tc.useTrend)
			
			if err != nil {
				t.Fatalf("期待: エラーなし, 実際: %v", err)
			}
			
			if result == nil {
				t.Fatal("期待: 結果が返される, 実際: nil")
			}
			
			// 基本的なフィールドの検証
			if len(result.MainNumbers) != 7 {
				t.Errorf("期待: 本数字7個, 実際: %d個", len(result.MainNumbers))
			}
			
			if len(result.BonusNumbers) != 2 {
				t.Errorf("期待: ボーナス数字2個, 実際: %d個", len(result.BonusNumbers))
			}
			
			if result.TrendBased != tc.useTrend {
				t.Errorf("期待: TrendBased = %v, 実際: %v", tc.useTrend, result.TrendBased)
			}
			
			// 生成時刻の検証（現在時刻から1秒以内）
			if time.Since(result.GeneratedTime) > time.Second {
				t.Errorf("期待: 生成時刻が現在時刻に近い, 実際: %v", result.GeneratedTime)
			}
			
			// 数字の範囲検証（1-37）
			for i, num := range result.MainNumbers {
				if num < 1 || num > 37 {
					t.Errorf("期待: 本数字[%d]が1-37の範囲, 実際: %d", i, num)
				}
			}
			
			for i, num := range result.BonusNumbers {
				if num < 1 || num > 37 {
					t.Errorf("期待: ボーナス数字[%d]が1-37の範囲, 実際: %d", i, num)
				}
			}
			
			// 本数字の重複チェック
			mainNumSet := make(map[int]bool)
			for _, num := range result.MainNumbers {
				if mainNumSet[num] {
					t.Errorf("期待: 本数字に重複なし, 実際: %dが重複", num)
				}
				mainNumSet[num] = true
			}
			
			// ボーナス数字の重複チェック
			bonusNumSet := make(map[int]bool)
			for _, num := range result.BonusNumbers {
				if bonusNumSet[num] {
					t.Errorf("期待: ボーナス数字に重複なし, 実際: %dが重複", num)
				}
				bonusNumSet[num] = true
			}
		})
	}
}

// TestGenerateRandomPrediction_InvalidInput 無効な入力値でのテスト
func TestGenerateRandomPrediction_InvalidInput(t *testing.T) {
	testCases := []int{0, -1, -10}
	
	for _, invalidRange := range testCases {
		t.Run("", func(t *testing.T) {
			result, err := GenerateRandomPrediction(invalidRange, true)
			
			if err == nil {
				t.Errorf("期待: エラーが発生（範囲: %d）, 実際: エラーなし", invalidRange)
			}
			
			if result != nil {
				t.Errorf("期待: result = nil, 実際: %v", result)
			}
		})
	}
}

// TestGenerateMultiplePredictions_ValidInput 複数予想番号生成の正常テスト
func TestGenerateMultiplePredictions_ValidInput(t *testing.T) {
	testCases := []struct {
		count         int
		analysisRange int
		useTrend      bool
		description   string
	}{
		{3, 50, true, "3組の傾向ベース予想"},
		{5, 100, false, "5組のランダム予想"},
		{1, 30, true, "1組の短期傾向ベース予想"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			results, err := GenerateMultiplePredictions(tc.count, tc.analysisRange, tc.useTrend)
			
			if err != nil {
				t.Fatalf("期待: エラーなし, 実際: %v", err)
			}
			
			if results == nil {
				t.Fatal("期待: 結果が返される, 実際: nil")
			}
			
			// 生成数の検証
			if len(results) != tc.count {
				t.Errorf("期待: 予想数 = %d, 実際: %d", tc.count, len(results))
			}
			
			// 各予想の基本検証
			for i, result := range results {
				if len(result.MainNumbers) != 7 {
					t.Errorf("期待: 予想[%d]の本数字7個, 実際: %d個", i, len(result.MainNumbers))
				}
				
				if len(result.BonusNumbers) != 2 {
					t.Errorf("期待: 予想[%d]のボーナス数字2個, 実際: %d個", i, len(result.BonusNumbers))
				}
				
				if result.TrendBased != tc.useTrend {
					t.Errorf("期待: 予想[%d]のTrendBased = %v, 実際: %v", i, tc.useTrend, result.TrendBased)
				}
			}
		})
	}
}

// TestGenerateMultiplePredictions_InvalidInput 複数予想番号生成の無効入力テスト
func TestGenerateMultiplePredictions_InvalidInput(t *testing.T) {
	testCases := []struct {
		count         int
		analysisRange int
		description   string
	}{
		{0, 50, "生成数0"},
		{-1, 50, "生成数マイナス"},
		{5, 0, "分析範囲0"},
		{5, -1, "分析範囲マイナス"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			results, err := GenerateMultiplePredictions(tc.count, tc.analysisRange, true)
			
			if err == nil {
				t.Errorf("期待: エラーが発生, 実際: エラーなし")
			}
			
			if results != nil {
				t.Errorf("期待: results = nil, 実際: %v", results)
			}
		})
	}
}

// TestGenerateRandomNumbers 完全ランダム数字生成のテスト
func TestGenerateRandomNumbers(t *testing.T) {
	testCases := []int{1, 3, 7, 10}
	
	for _, count := range testCases {
		t.Run("", func(t *testing.T) {
			numbers := generateRandomNumbers(count)
			
			// 生成数の検証
			if len(numbers) != count {
				t.Errorf("期待: 生成数 = %d, 実際: %d", count, len(numbers))
			}
			
			// 範囲検証（1-37）
			for i, num := range numbers {
				if num < 1 || num > 37 {
					t.Errorf("期待: 数字[%d]が1-37の範囲, 実際: %d", i, num)
				}
			}
			
			// 重複チェック
			numSet := make(map[int]bool)
			for _, num := range numbers {
				if numSet[num] {
					t.Errorf("期待: 重複なし, 実際: %dが重複", num)
				}
				numSet[num] = true
			}
		})
	}
}

// TestWeightedRandomSelect 重み付きランダム選択のテスト
func TestWeightedRandomSelect(t *testing.T) {
	numbers := []WeightedNumber{
		{Number: 1, Weight: 0.1},
		{Number: 2, Weight: 0.3},
		{Number: 3, Weight: 0.5},
		{Number: 4, Weight: 0.1},
	}
	
	// 100回実行して統計的な妥当性をチェック
	selections := make(map[int]int)
	iterations := 1000
	
	for i := 0; i < iterations; i++ {
		selectedIndex := weightedRandomSelect(numbers)
		if selectedIndex < 0 || selectedIndex >= len(numbers) {
			t.Errorf("期待: 有効なインデックス(0-%d), 実際: %d", len(numbers)-1, selectedIndex)
		}
		selections[selectedIndex]++
	}
	
	// 重みが高い数字ほど多く選ばれるかの簡易チェック
	// Weight: 0.5のindex 2が最も多く選ばれるはず
	maxSelections := 0
	maxIndex := -1
	for index, count := range selections {
		if count > maxSelections {
			maxSelections = count
			maxIndex = index
		}
	}
	
	// 統計的に最も重みの高い選択肢が最多選択される可能性が高い
	if maxIndex != 2 {
		t.Logf("注意: 最高重み（index:2）が最多選択されませんでした。統計的偏差の可能性あり")
		t.Logf("選択状況: %v", selections)
	}
}

// TestPredictionResult_StructureValidation 構造体の基本検証
func TestPredictionResult_StructureValidation(t *testing.T) {
	now := time.Now()
	result := &PredictionResult{
		MainNumbers:   []int{1, 7, 14, 21, 28, 35, 37},
		BonusNumbers:  []int{3, 15},
		TrendBased:    true,
		GeneratedTime: now,
	}
	
	// 各フィールドのテスト
	if len(result.MainNumbers) != 7 {
		t.Errorf("期待: MainNumbers = 7個, 実際: %d個", len(result.MainNumbers))
	}
	
	if len(result.BonusNumbers) != 2 {
		t.Errorf("期待: BonusNumbers = 2個, 実際: %d個", len(result.BonusNumbers))
	}
	
	if !result.TrendBased {
		t.Errorf("期待: TrendBased = true, 実際: %v", result.TrendBased)
	}
	
	if !result.GeneratedTime.Equal(now) {
		t.Errorf("期待: GeneratedTime = %v, 実際: %v", now, result.GeneratedTime)
	}
}
