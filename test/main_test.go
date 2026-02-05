package test

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// TestMain 全テストのエントリーポイント
func TestMain(m *testing.M) {
	fmt.Println("=== ロト7 Promise 全テスト開始 ===")
	fmt.Println("テスト実行時刻:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()

	// テストの前処理
	setupTests()

	// テスト実行
	exitCode := m.Run()

	// テストの後処理
	teardownTests()

	fmt.Println()
	if exitCode == 0 {
		fmt.Println("✅ 全テストが成功しました！")
	} else {
		fmt.Println("❌ テストに失敗しました")
	}
	fmt.Println("=== ロト7 Promise 全テスト終了 ===")

	os.Exit(exitCode)
}

// setupTests テスト前の準備
func setupTests() {
	fmt.Println("🔧 テスト環境を準備中...")
	// 必要に応じてテスト前の準備処理をここに追加
	fmt.Println("✅ テスト環境の準備完了")
}

// teardownTests テスト後のクリーンアップ
func teardownTests() {
	fmt.Println("🧹 テスト環境をクリーンアップ中...")
	// 必要に応じてテスト後のクリーンアップ処理をここに追加
	fmt.Println("✅ クリーンアップ完了")
}

// TestIntegration 全テスト
func TestIntegration(t *testing.T) {
	t.Run("全機能全テスト", func(t *testing.T) {
		t.Log("結果取得、ヒートマップ生成、推薦機能の全テストを実行")

		// 各機能が最低限動作することを確認
		t.Run("結果取得機能", func(t *testing.T) {
			TestResultGetData(t)
		})

		t.Run("ヒートマップ機能", func(t *testing.T) {
			TestHeatmapGenerate(t)
		})

		t.Run("推薦機能", func(t *testing.T) {
			TestRecommendationDefaultConfig_StandardValues(t)
		})

		t.Log("全機能の基本動作を確認しました")
	})
}

// TestPerformance パフォーマンステスト
func TestPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("パフォーマンステストをスキップします（-short フラグが指定されています）")
	}

	t.Run("推薦生成パフォーマンス", func(t *testing.T) {
		BenchmarkRecommendationGenerate(&testing.B{})
		t.Log("推薦生成のパフォーマンステストを実行しました")
	})
}

// TestErrorHandling エラーハンドリングテスト
func TestErrorHandling(t *testing.T) {
	t.Run("不正入力のエラーハンドリング", func(t *testing.T) {
		t.Run("ヒートマップ不正入力", func(t *testing.T) {
			TestHeatmapInvalidInput(t)
		})

		t.Log("エラーハンドリングテストを実行しました")
	})
}

// TestDataIntegrity データ整合性テスト
func TestDataIntegrity(t *testing.T) {
	t.Run("データフォーマット整合性", func(t *testing.T) {
		t.Run("結果データフォーマット", func(t *testing.T) {
			TestResultDataFormat(t)
		})

		t.Run("ヒートマップデータ整合性", func(t *testing.T) {
			TestHeatmapStatisticalCalculations(t)
		})

		t.Log("データ整合性テストを実行しました")
	})
}
