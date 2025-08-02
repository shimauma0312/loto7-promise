package heatmap

import (
	"testing"
)

// TestGenerateHeatmap ヒートマップ生成の基本テスト
func TestGenerateHeatmap(t *testing.T) {
	// 正常ケース：過去10回分のデータで検証
	summary, err := GenerateHeatmap(10)
	if err != nil {
		t.Fatalf("ヒートマップ生成でエラーが発生しました: %v", err)
	}

	// 基本的な構造チェック
	if summary == nil {
		t.Fatal("ヒートマップサマリーがnilです")
	}

	if summary.SearchRange != 10 {
		t.Errorf("検索範囲が期待値と異なります。期待値: 10, 実際: %d", summary.SearchRange)
	}

	if summary.ActualDraws <= 0 {
		t.Errorf("実際の抽選回数が不正です: %d", summary.ActualDraws)
	}

	// PositionDataの構造チェック
	if len(summary.PositionData) != 38 { // インデックス0を含むため38
		t.Errorf("PositionDataのサイズが不正です。期待値: 38, 実際: %d", len(summary.PositionData))
	}

	// 数字1-37の各位置データをチェック
	for number := 1; number <= 37; number++ {
		if len(summary.PositionData[number]) != 8 { // インデックス0を含むため8
			t.Errorf("数字%dの位置データサイズが不正です。期待値: 8, 実際: %d", number, len(summary.PositionData[number]))
		}

		for position := 1; position <= 7; position++ {
			data := summary.PositionData[number][position]
			
			if data.Number != number {
				t.Errorf("数字データが不正です。期待値: %d, 実際: %d", number, data.Number)
			}
			
			if data.Position != position {
				t.Errorf("位置データが不正です。期待値: %d, 実際: %d", position, data.Position)
			}
			
			if data.Count < 0 {
				t.Errorf("出現回数が負の値です: %d", data.Count)
			}
			
			if data.TotalDraws != summary.ActualDraws {
				t.Errorf("総抽選回数が不一致です。期待値: %d, 実際: %d", summary.ActualDraws, data.TotalDraws)
			}
		}
	}
}

// TestGenerateHeatmapInvalidInput 不正な入力値のテスト
func TestGenerateHeatmapInvalidInput(t *testing.T) {
	testCases := []struct {
		name        string
		searchRange int
		expectError bool
	}{
		{"負の値", -1, true},
		{"ゼロ", 0, true},
		{"正の値", 1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GenerateHeatmap(tc.searchRange)
			
			if tc.expectError && err == nil {
				t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
			}
			
			if !tc.expectError && err != nil {
				t.Errorf("エラーが発生しました: %v", err)
			}
		})
	}
}

// TestGetPositionData 指定数字の位置別データ取得テスト
func TestGetPositionData(t *testing.T) {
	summary, err := GenerateHeatmap(5)
	if err != nil {
		t.Fatalf("ヒートマップ生成でエラーが発生しました: %v", err)
	}

	// 正常ケース
	data, err := summary.GetPositionData(1)
	if err != nil {
		t.Fatalf("位置データ取得でエラーが発生しました: %v", err)
	}

	if len(data) != 7 {
		t.Errorf("位置データの長さが不正です。期待値: 7, 実際: %d", len(data))
	}

	// 範囲外の数字
	invalidNumbers := []int{0, 38, -1, 100}
	for _, number := range invalidNumbers {
		_, err := summary.GetPositionData(number)
		if err == nil {
			t.Errorf("数字%dで範囲外エラーが期待されましたが、エラーが発生しませんでした", number)
		}
	}
}

// TestGetNumbersByPosition 指定位置での全数字データ取得テスト
func TestGetNumbersByPosition(t *testing.T) {
	summary, err := GenerateHeatmap(5)
	if err != nil {
		t.Fatalf("ヒートマップ生成でエラーが発生しました: %v", err)
	}

	// 正常ケース
	data, err := summary.GetNumbersByPosition(1)
	if err != nil {
		t.Fatalf("位置別数字データ取得でエラーが発生しました: %v", err)
	}

	if len(data) != 37 {
		t.Errorf("数字データの長さが不正です。期待値: 37, 実際: %d", len(data))
	}

	// 各データの妥当性チェック
	for i, posData := range data {
		expectedNumber := i + 1
		if posData.Number != expectedNumber {
			t.Errorf("数字が不正です。期待値: %d, 実際: %d", expectedNumber, posData.Number)
		}
		
		if posData.Position != 1 {
			t.Errorf("位置が不正です。期待値: 1, 実際: %d", posData.Position)
		}
	}

	// 範囲外の位置
	invalidPositions := []int{0, 8, -1, 10}
	for _, position := range invalidPositions {
		_, err := summary.GetNumbersByPosition(position)
		if err == nil {
			t.Errorf("位置%dで範囲外エラーが期待されましたが、エラーが発生しませんでした", position)
		}
	}
}

// TestStatisticalCalculations 統計計算のテスト
func TestStatisticalCalculations(t *testing.T) {
	summary, err := GenerateHeatmap(10)
	if err != nil {
		t.Fatalf("ヒートマップ生成でエラーが発生しました: %v", err)
	}

	// HottestNumbersとColdestNumbersが空でないことを確認
	if len(summary.HottestNumbers) == 0 {
		t.Error("最頻出数字が空です")
	}

	if len(summary.ColdestNumbers) == 0 {
		t.Error("最低出現数字が空です")
	}

	// MostBiasedNumberが妥当な範囲内であることを確認
	if summary.MostBiasedNumber < 1 || summary.MostBiasedNumber > 37 {
		t.Errorf("最偏位数字が範囲外です: %d", summary.MostBiasedNumber)
	}

	// HottestNumbersの各数字が妥当な範囲内であることを確認
	for _, number := range summary.HottestNumbers {
		if number < 1 || number > 37 {
			t.Errorf("最頻出数字が範囲外です: %d", number)
		}
	}

	// ColdestNumbersの各数字が妥当な範囲内であることを確認
	for _, number := range summary.ColdestNumbers {
		if number < 1 || number > 37 {
			t.Errorf("最低出現数字が範囲外です: %d", number)
		}
	}
}

// TestPercentageCalculations パーセンテージ計算のテスト
func TestPercentageCalculations(t *testing.T) {
	summary, err := GenerateHeatmap(5)
	if err != nil {
		t.Fatalf("ヒートマップ生成でエラーが発生しました: %v", err)
	}

	// 各数字の各位置でのパーセンテージが0-100の範囲内であることを確認
	for number := 1; number <= 37; number++ {
		for position := 1; position <= 7; position++ {
			data := summary.PositionData[number][position]
			
			if data.Percentage < 0 || data.Percentage > 100 {
				t.Errorf("出現率が範囲外です。数字%d、位置%d: %f%%", number, position, data.Percentage)
			}
			
			if data.PositionRatio < 0 || data.PositionRatio > 100 {
				t.Errorf("位置別出現率が範囲外です。数字%d、位置%d: %f%%", number, position, data.PositionRatio)
			}
		}
	}
}
