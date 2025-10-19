package main

import (
	"testing"

	"github.com/shimauma0312/loto7-promise/internal/recommendation"
)

// TestRecommendationDefaultConfig デフォルト設定のテスト
func TestRecommendationDefaultConfig(t *testing.T) {
	config := recommendation.DefaultConfig()

	if config.MaxRecommendations != 5 {
		t.Errorf("MaxRecommendations = %d, want 5", config.MaxRecommendations)
	}

	if config.RecentAvoidCount != 3 {
		t.Errorf("RecentAvoidCount = %d, want 3", config.RecentAvoidCount)
	}

	if config.ConsecutiveBoost != 3 {
		t.Errorf("ConsecutiveBoost = %d, want 3", config.ConsecutiveBoost)
	}

	if config.HistoryLookback != 100 {
		t.Errorf("HistoryLookback = %d, want 100", config.HistoryLookback)
	}
}

// TestRecommendationNewEngine エンジン作成のテスト
func TestRecommendationNewEngine(t *testing.T) {
	config := recommendation.DefaultConfig()
	engine := recommendation.NewRecommendationEngine(config)

	if engine == nil {
		t.Fatal("NewRecommendationEngine returned nil")
	}
}

// TestRecommendationFormatOutput フォーマット機能のテスト
func TestRecommendationFormatOutput(t *testing.T) {
	config := recommendation.DefaultConfig()
	engine := recommendation.NewRecommendationEngine(config)

	// ダミーデータでテスト
	recommendations := [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{8, 9, 10, 11, 12, 13, 14},
	}

	result := engine.FormatRecommendations(recommendations)

	if result == "" {
		t.Error("FormatRecommendations returned empty string")
	}

	// 推薦結果が含まれているかチェック
	if !containsSubstring(result, "推薦 1:") {
		t.Error("result should contain '推薦 1:'")
	}

	if !containsSubstring(result, "推薦 2:") {
		t.Error("result should contain '推薦 2:'")
	}

	if !containsSubstring(result, "分析情報") {
		t.Error("result should contain '分析情報'")
	}
}

// ヘルパー関数
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ベンチマークテスト
func BenchmarkRecommendationGenerate(b *testing.B) {
	config := recommendation.DefaultConfig()
	config.MaxRecommendations = 3 // ベンチマーク用に数を少なく
	engine := recommendation.NewRecommendationEngine(config)

	// 小規模なテストデータでデータ読み込み
	err := engine.LoadData(20) // 20回分のデータを使用
	if err != nil {
		b.Skip("データの読み込みに失敗したため、ベンチマークをスキップします:", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.GenerateRecommendations()
		if err != nil {
			b.Fatal(err)
		}
	}
}
