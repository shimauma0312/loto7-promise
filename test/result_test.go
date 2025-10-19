package main

import (
	"testing"

	"github.com/shimauma0312/loto7-promise/internal/result"
)

// TestResultGetData 結果取得のテスト
func TestResultGetData(t *testing.T) {
	// 小規模なテスト（5回分）
	results := result.GetResult(5)

	if len(results) == 0 {
		t.Error("結果データが取得できませんでした")
		return
	}

	// 各結果が7個の数字を含むかチェック
	for i, resultData := range results {
		if len(resultData) != 7 {
			t.Errorf("結果%d: 数字の個数が不正です。期待値: 7, 実際: %d", i+1, len(resultData))
		}
	}
}

// TestResultDataFormat 結果データフォーマットのテスト
func TestResultDataFormat(t *testing.T) {
	// 1回分の結果を取得
	results := result.GetResult(1)

	if len(results) == 0 {
		t.Skip("結果データが取得できないため、テストをスキップします")
		return
	}

	numbers := results[0]

	// 各数字が2桁の文字列形式かチェック
	for i, numStr := range numbers {
		if len(numStr) != 2 {
			t.Errorf("数字%d: フォーマットが不正です。期待値: 2桁文字列, 実際: '%s'", i+1, numStr)
		}

		// 数字が01-37の範囲内かチェック（文字列として）
		if numStr < "01" || numStr > "37" {
			t.Errorf("数字%d: 範囲外の値です。値: '%s'", i+1, numStr)
		}
	}
}
