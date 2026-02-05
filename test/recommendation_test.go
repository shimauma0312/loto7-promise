package test

import (
	"fmt"
	"testing"

	"github.com/shimauma0312/loto7-promise/internal/recommendation"
)

// TestRecommendationDefaultConfig_StandardValues デフォルト設定の標準値をテスト
func TestRecommendationDefaultConfig_StandardValues(t *testing.T) {
	config := recommendation.DefaultConfig()

	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"MaxRecommendations", config.MaxRecommendations, 5},
		{"RecentAvoidCount", config.RecentAvoidCount, 3},
		{"ConsecutiveBoost", config.ConsecutiveBoost, 3},
		{"HistoryLookback", config.HistoryLookback, 100},
		{"FrequencyWeight", config.FrequencyWeight, 0.3},
		{"RecentWeight", config.RecentWeight, 0.4},
		{"ConsecutiveWeight", config.ConsecutiveWeight, 0.2},
		{"PositionWeight", config.PositionWeight, 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

// TestRecommendationNewEngine_NotNil エンジン作成が成功することをテスト
func TestRecommendationNewEngine_NotNil(t *testing.T) {
	config := recommendation.DefaultConfig()
	engine := recommendation.NewRecommendationEngine(config)

	if engine == nil {
		t.Fatal("NewRecommendationEngine returned nil")
	}
}

// TestRecommendationNewEngine_ConfigIsSet エンジンに設定が反映されることをテスト
func TestRecommendationNewEngine_ConfigIsSet(t *testing.T) {
	config := recommendation.DefaultConfig()
	config.MaxRecommendations = 10

	engine := recommendation.NewRecommendationEngine(config)
	engineConfig := engine.GetConfig()

	if engineConfig.MaxRecommendations != 10 {
		t.Errorf("Engine config MaxRecommendations = %d, want 10", engineConfig.MaxRecommendations)
	}
}

// TestLoadData_Success データ読み込みが成功することをテスト
func TestLoadData_Success(t *testing.T) {
	config := recommendation.DefaultConfig()
	engine := recommendation.NewRecommendationEngine(config)

	err := engine.LoadData(50)
	if err != nil {
		t.Fatalf("LoadData failed: %v", err)
	}

	count := engine.GetAnalyzedDrawsCount()
	if count == 0 {
		t.Error("LoadData did not load any data")
	}
}

// TestLoadData_WithDifferentCounts 異なるデータ件数で読み込みをテスト
func TestLoadData_WithDifferentCounts(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{"Small dataset", 10},
		{"Medium dataset", 50},
		{"Large dataset", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := recommendation.DefaultConfig()
			engine := recommendation.NewRecommendationEngine(config)

			err := engine.LoadData(tt.count)
			if err != nil {
				t.Fatalf("LoadData(%d) failed: %v", tt.count, err)
			}

			analyzedCount := engine.GetAnalyzedDrawsCount()
			if analyzedCount == 0 {
				t.Errorf("LoadData(%d) loaded 0 results", tt.count)
			}
		})
	}
}

// TestGenerateRecommendations_Success 推薦生成が成功することをテスト
func TestGenerateRecommendations_Success(t *testing.T) {
	config := recommendation.DefaultConfig()
	config.MaxRecommendations = 3
	engine := recommendation.NewRecommendationEngine(config)

	err := engine.LoadData(50)
	if err != nil {
		t.Fatalf("LoadData failed: %v", err)
	}

	recommendations, err := engine.GenerateRecommendations()
	if err != nil {
		t.Fatalf("GenerateRecommendations failed: %v", err)
	}

	if len(recommendations) == 0 {
		t.Error("GenerateRecommendations returned 0 recommendations")
	}

	if len(recommendations) > config.MaxRecommendations {
		t.Errorf("GenerateRecommendations returned %d recommendations, max is %d",
			len(recommendations), config.MaxRecommendations)
	}
}

// TestGenerateRecommendations_WithoutLoadData データ未読み込みでエラーとなることをテスト
func TestGenerateRecommendations_WithoutLoadData(t *testing.T) {
	config := recommendation.DefaultConfig()
	engine := recommendation.NewRecommendationEngine(config)

	_, err := engine.GenerateRecommendations()
	if err == nil {
		t.Error("GenerateRecommendations should return error when data is not loaded")
	}
}

// TestGenerateRecommendations_CombinationFormat 生成された組み合わせの形式をテスト
func TestGenerateRecommendations_CombinationFormat(t *testing.T) {
	config := recommendation.DefaultConfig()
	config.MaxRecommendations = 5
	engine := recommendation.NewRecommendationEngine(config)

	err := engine.LoadData(100)
	if err != nil {
		t.Fatalf("LoadData failed: %v", err)
	}

	recommendations, err := engine.GenerateRecommendations()
	if err != nil {
		t.Fatalf("GenerateRecommendations failed: %v", err)
	}

	for i, combination := range recommendations {
		t.Run(formatCombinationTestName(i, combination), func(t *testing.T) {
			// 7個の数字であることを確認
			if len(combination) != 7 {
				t.Errorf("Combination %d has %d numbers, want 7", i+1, len(combination))
			}

			// 数字が1-37の範囲内であることを確認
			for _, num := range combination {
				if num < 1 || num > 37 {
					t.Errorf("Combination %d contains invalid number %d", i+1, num)
				}
			}

			// 重複がないことを確認
			seen := make(map[int]bool)
			for _, num := range combination {
				if seen[num] {
					t.Errorf("Combination %d contains duplicate number %d", i+1, num)
				}
				seen[num] = true
			}

			// ソート済みであることを確認
			for j := 1; j < len(combination); j++ {
				if combination[j] <= combination[j-1] {
					t.Errorf("Combination %d is not sorted: %v", i+1, combination)
					break
				}
			}
		})
	}
}

// TestGenerateRecommendations_MeetsConstraints 生成された組み合わせが制約を満たすことをテスト
func TestGenerateRecommendations_MeetsConstraints(t *testing.T) {
	config := recommendation.DefaultConfig()
	config.MaxRecommendations = 5
	engine := recommendation.NewRecommendationEngine(config)

	err := engine.LoadData(100)
	if err != nil {
		t.Fatalf("LoadData failed: %v", err)
	}

	recommendations, err := engine.GenerateRecommendations()
	if err != nil {
		t.Fatalf("GenerateRecommendations failed: %v", err)
	}

	for i, combination := range recommendations {
		t.Run(formatCombinationTestName(i, combination), func(t *testing.T) {
			// ゾーン分布をチェック
			zone1, zone2, zone3 := countZones(combination)

			if zone1 < 3 || zone1 > 4 {
				t.Errorf("Combination %d: zone1=%d, want 3-4", i+1, zone1)
			}
			if zone2 < 0 || zone2 > 4 {
				t.Errorf("Combination %d: zone2=%d, should be 0-4", i+1, zone2)
			}
			if zone3 < 2 || zone3 > 3 {
				t.Errorf("Combination %d: zone3=%d, want 2-3", i+1, zone3)
			}

			// 奇数偶数比率をチェック
			oddCount, evenCount := countOddEven(combination)
			validRatio := (oddCount == 4 && evenCount == 3) || (oddCount == 3 && evenCount == 4)
			if !validRatio {
				t.Errorf("Combination %d: odd=%d, even=%d, want 4:3 or 3:4", i+1, oddCount, evenCount)
			}

			// 合計値をチェック
			sum := sumNumbers(combination)
			if sum < 100 || sum > 160 {
				t.Errorf("Combination %d: sum=%d, want 100-160", i+1, sum)
			}

			// 平均値をチェック
			avg := float64(sum) / float64(len(combination))
			if avg < 16.0 || avg > 23.0 {
				t.Errorf("Combination %d: avg=%.2f, want 16.0-23.0", i+1, avg)
			}
		})
	}
}

// TestFormatRecommendations_OutputStructure フォーマット出力の構造をテスト
func TestFormatRecommendations_OutputStructure(t *testing.T) {
	config := recommendation.DefaultConfig()
	engine := recommendation.NewRecommendationEngine(config)

	recommendations := [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{8, 9, 10, 11, 12, 13, 14},
	}

	result := engine.FormatRecommendations(recommendations)

	if result == "" {
		t.Fatal("FormatRecommendations returned empty string")
	}

	requiredStrings := []string{
		"推薦 1:",
		"推薦 2:",
		"分析情報",
		"統計データ",
	}

	for _, required := range requiredStrings {
		if !containsSubstring(result, required) {
			t.Errorf("result should contain '%s'", required)
		}
	}
}

// TestFormatRecommendations_EmptyRecommendations 空の推薦リストでもエラーにならないことをテスト
func TestFormatRecommendations_EmptyRecommendations(t *testing.T) {
	config := recommendation.DefaultConfig()
	engine := recommendation.NewRecommendationEngine(config)

	recommendations := [][]int{}
	result := engine.FormatRecommendations(recommendations)

	if result == "" {
		t.Error("FormatRecommendations should return non-empty string even for empty recommendations")
	}
}

// TestGetConfig_ReturnsCorrectConfig GetConfigが正しい設定を返すことをテスト
func TestGetConfig_ReturnsCorrectConfig(t *testing.T) {
	config := recommendation.DefaultConfig()
	config.MaxRecommendations = 7
	config.RecentAvoidCount = 5

	engine := recommendation.NewRecommendationEngine(config)
	retrievedConfig := engine.GetConfig()

	if retrievedConfig.MaxRecommendations != 7 {
		t.Errorf("GetConfig().MaxRecommendations = %d, want 7", retrievedConfig.MaxRecommendations)
	}
	if retrievedConfig.RecentAvoidCount != 5 {
		t.Errorf("GetConfig().RecentAvoidCount = %d, want 5", retrievedConfig.RecentAvoidCount)
	}
}

// TestGetAnalyzedDrawsCount_AfterLoadData LoadData後に正しい件数を返すことをテスト
func TestGetAnalyzedDrawsCount_AfterLoadData(t *testing.T) {
	config := recommendation.DefaultConfig()
	engine := recommendation.NewRecommendationEngine(config)

	err := engine.LoadData(50)
	if err != nil {
		t.Fatalf("LoadData failed: %v", err)
	}

	count := engine.GetAnalyzedDrawsCount()
	if count == 0 {
		t.Error("GetAnalyzedDrawsCount returned 0 after LoadData")
	}
}

// TestGetAnalyzedDrawsCount_BeforeLoadData LoadData前は0を返すことをテスト
func TestGetAnalyzedDrawsCount_BeforeLoadData(t *testing.T) {
	config := recommendation.DefaultConfig()
	engine := recommendation.NewRecommendationEngine(config)

	count := engine.GetAnalyzedDrawsCount()
	if count != 0 {
		t.Errorf("GetAnalyzedDrawsCount = %d before LoadData, want 0", count)
	}
}

// ヘルパー関数: 文字列に部分文字列が含まれるかチェック
func containsSubstring(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ヘルパー関数: ゾーンごとの数字をカウント
func countZones(combination []int) (zone1, zone2, zone3 int) {
	for _, num := range combination {
		if num >= 1 && num <= 13 {
			zone1++
		} else if num >= 14 && num <= 26 {
			zone2++
		} else if num >= 27 && num <= 37 {
			zone3++
		}
	}
	return
}

// ヘルパー関数: 奇数偶数をカウント
func countOddEven(combination []int) (odd, even int) {
	for _, num := range combination {
		if num%2 == 0 {
			even++
		} else {
			odd++
		}
	}
	return
}

// ヘルパー関数: 数字の合計を計算
func sumNumbers(combination []int) int {
	sum := 0
	for _, num := range combination {
		sum += num
	}
	return sum
}

// ヘルパー関数: 組み合わせテスト名をフォーマット
func formatCombinationTestName(index int, combination []int) string {
	return fmt.Sprintf("組み合わせ%d: %v", index+1, combination)
}

// ベンチマークテスト: 推薦生成のパフォーマンス
func BenchmarkRecommendationGenerate(b *testing.B) {
	config := recommendation.DefaultConfig()
	config.MaxRecommendations = 3
	engine := recommendation.NewRecommendationEngine(config)

	err := engine.LoadData(20)
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

// ベンチマークテスト: データ読み込みのパフォーマンス
func BenchmarkLoadData(b *testing.B) {
	config := recommendation.DefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine := recommendation.NewRecommendationEngine(config)
		err := engine.LoadData(100)
		if err != nil {
			b.Fatal(err)
		}
	}
}
