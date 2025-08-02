package consecutivePattern

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/shimauma0312/loto7-promise/internal/result"
)

// ConsecutivePatternResult 連番・ゾロ目パターンの分析結果
type ConsecutivePatternResult struct {
	PatternType       string    `json:"pattern_type"`       // パターンタイプ（連番、ゾロ目）
	Pattern           []int     `json:"pattern"`            // 見つかったパターン
	Count             int       `json:"count"`              // 出現回数
	TotalDraws        int       `json:"total_draws"`        // 調査対象の抽選回数
	Percentage        float64   `json:"percentage"`         // 出現率（%）
	LastAppearance    *int      `json:"last_appearance"`    // 最後に出現した回（nilの場合は未出現）
	AppearanceHistory []int     `json:"appearance_history"` // 出現履歴
}

// PatternSummary パターン全体のサマリー
type PatternSummary struct {
	TotalDraws          int                         `json:"total_draws"`            // 調査対象の抽選回数
	ConsecutivePatterns []*ConsecutivePatternResult `json:"consecutive_patterns"`   // 連番パターンの結果
	SameDigitPatterns   []*ConsecutivePatternResult `json:"same_digit_patterns"`    // ゾロ目パターンの結果
	OverallStats        *OverallStats               `json:"overall_stats"`          // 全体統計
}

// OverallStats 全体統計情報
type OverallStats struct {
	TotalConsecutiveCount int     `json:"total_consecutive_count"` // 連番が含まれる抽選回数
	TotalSameDigitCount   int     `json:"total_same_digit_count"`   // ゾロ目が含まれる抽選回数
	ConsecutivePercentage float64 `json:"consecutive_percentage"`   // 連番が含まれる割合
	SameDigitPercentage   float64 `json:"same_digit_percentage"`    // ゾロ目が含まれる割合
}

// AnalyzeConsecutivePatterns 連番・ゾロ目パターンを分析する
// searchRange: 過去何回分を調査するか
func AnalyzeConsecutivePatterns(searchRange int) (*PatternSummary, error) {
	// 入力値妥当性チェック
	if searchRange <= 0 {
		return nil, fmt.Errorf("検索範囲は1以上で指定してください。入力値: %d", searchRange)
	}

	// ロト7の結果データを取得
	records := result.GetResult(searchRange)
	if records == nil {
		return nil, fmt.Errorf("ロト7の結果データを取得できませんでした")
	}

	totalDraws := len(records)

	// 各抽選結果を数値配列に変換
	drawResults := make([][]int, 0, totalDraws)
	for _, record := range records {
		numbers := make([]int, 0, 7)
		for i := 0; i < 7 && i < len(record); i++ {
			if num, err := strconv.Atoi(record[i]); err == nil {
				numbers = append(numbers, num)
			}
		}
		if len(numbers) == 7 {
			// ソートして分析しやすくする
			sort.Ints(numbers)
			drawResults = append(drawResults, numbers)
		}
	}

	// 連番パターンを分析
	consecutivePatterns := analyzeConsecutiveNumbers(drawResults)
	
	// ゾロ目パターンを分析
	sameDigitPatterns := analyzeSameDigitNumbers(drawResults)

	// 全体統計を計算
	overallStats := calculateOverallStats(drawResults, consecutivePatterns, sameDigitPatterns)

	return &PatternSummary{
		TotalDraws:          totalDraws,
		ConsecutivePatterns: consecutivePatterns,
		SameDigitPatterns:   sameDigitPatterns,
		OverallStats:        overallStats,
	}, nil
}

// analyzeConsecutiveNumbers 連番パターンを分析
func analyzeConsecutiveNumbers(drawResults [][]int) []*ConsecutivePatternResult {
	// 2連番、3連番、4連番のパターンを検出
	patternMap := make(map[string]*ConsecutivePatternResult)
	
	for drawIndex, numbers := range drawResults {
		// 連番検出ロジック
		consecutiveGroups := findConsecutiveGroups(numbers)
		
		for _, group := range consecutiveGroups {
			if len(group) < 2 {
				continue // 2個未満は連番とみなさない
			}
			
			patternKey := generatePatternKey(group, "consecutive")
			
			if pattern, exists := patternMap[patternKey]; exists {
				pattern.Count++
				pattern.AppearanceHistory = append(pattern.AppearanceHistory, drawIndex+1)
				if pattern.LastAppearance == nil || *pattern.LastAppearance < drawIndex+1 {
					latest := drawIndex + 1
					pattern.LastAppearance = &latest
				}
			} else {
				latest := drawIndex + 1
				patternMap[patternKey] = &ConsecutivePatternResult{
					PatternType:       "連番",
					Pattern:           group,
					Count:             1,
					TotalDraws:        len(drawResults),
					LastAppearance:    &latest,
					AppearanceHistory: []int{drawIndex + 1},
				}
			}
		}
	}
	
	// パーセンテージを計算
	for _, pattern := range patternMap {
		pattern.Percentage = (float64(pattern.Count) / float64(len(drawResults))) * 100
	}
	
	// 結果を配列に変換してソート
	patterns := make([]*ConsecutivePatternResult, 0, len(patternMap))
	for _, pattern := range patternMap {
		patterns = append(patterns, pattern)
	}
	
	// 出現回数の多い順にソート
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Count > patterns[j].Count
	})
	
	return patterns
}

// analyzeSameDigitNumbers ゾロ目パターンを分析
func analyzeSameDigitNumbers(drawResults [][]int) []*ConsecutivePatternResult {
	patternMap := make(map[string]*ConsecutivePatternResult)
	
	for drawIndex, numbers := range drawResults {
		// ゾロ目検出ロジック
		sameDigitGroups := findSameDigitGroups(numbers)
		
		for _, group := range sameDigitGroups {
			if len(group) < 2 {
				continue // 2個未満はゾロ目とみなさない
			}
			
			patternKey := generatePatternKey(group, "same_digit")
			
			if pattern, exists := patternMap[patternKey]; exists {
				pattern.Count++
				pattern.AppearanceHistory = append(pattern.AppearanceHistory, drawIndex+1)
				if pattern.LastAppearance == nil || *pattern.LastAppearance < drawIndex+1 {
					latest := drawIndex + 1
					pattern.LastAppearance = &latest
				}
			} else {
				latest := drawIndex + 1
				patternMap[patternKey] = &ConsecutivePatternResult{
					PatternType:       "ゾロ目",
					Pattern:           group,
					Count:             1,
					TotalDraws:        len(drawResults),
					LastAppearance:    &latest,
					AppearanceHistory: []int{drawIndex + 1},
				}
			}
		}
	}
	
	// パーセンテージを計算
	for _, pattern := range patternMap {
		pattern.Percentage = (float64(pattern.Count) / float64(len(drawResults))) * 100
	}
	
	// 結果を配列に変換してソート
	patterns := make([]*ConsecutivePatternResult, 0, len(patternMap))
	for _, pattern := range patternMap {
		patterns = append(patterns, pattern)
	}
	
	// 出現回数の多い順にソート
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Count > patterns[j].Count
	})
	
	return patterns
}

// findConsecutiveGroups 数値配列から連番グループを見つける
func findConsecutiveGroups(numbers []int) [][]int {
	if len(numbers) < 2 {
		return nil
	}
	
	var groups [][]int
	var currentGroup []int
	
	for i, num := range numbers {
		if len(currentGroup) == 0 {
			currentGroup = append(currentGroup, num)
		} else {
			lastNum := currentGroup[len(currentGroup)-1]
			if num == lastNum+1 {
				// 連続している
				currentGroup = append(currentGroup, num)
			} else {
				// 連続が途切れた
				if len(currentGroup) >= 2 {
					groups = append(groups, make([]int, len(currentGroup)))
					copy(groups[len(groups)-1], currentGroup)
				}
				currentGroup = []int{num}
			}
		}
		
		// 最後の要素の場合
		if i == len(numbers)-1 && len(currentGroup) >= 2 {
			groups = append(groups, make([]int, len(currentGroup)))
			copy(groups[len(groups)-1], currentGroup)
		}
	}
	
	return groups
}

// findSameDigitGroups 同じ一の位を持つ数字のグループを見つける
func findSameDigitGroups(numbers []int) [][]int {
	digitMap := make(map[int][]int)
	
	// 一の位でグループ化
	for _, num := range numbers {
		digit := num % 10
		digitMap[digit] = append(digitMap[digit], num)
	}
	
	var groups [][]int
	for _, group := range digitMap {
		if len(group) >= 2 {
			// ソートして見やすくする
			sort.Ints(group)
			groups = append(groups, group)
		}
	}
	
	return groups
}

// generatePatternKey パターンの一意キーを生成
func generatePatternKey(pattern []int, patternType string) string {
	var parts []string
	for _, num := range pattern {
		parts = append(parts, strconv.Itoa(num))
	}
	return fmt.Sprintf("%s_%s", patternType, strings.Join(parts, "_"))
}

// calculateOverallStats 全体統計を計算
func calculateOverallStats(drawResults [][]int, consecutivePatterns, sameDigitPatterns []*ConsecutivePatternResult) *OverallStats {
	totalDraws := len(drawResults)
	consecutiveCount := 0
	sameDigitCount := 0
	
	// 各抽選回について連番・ゾロ目の有無をチェック
	for _, numbers := range drawResults {
		hasConsecutive := len(findConsecutiveGroups(numbers)) > 0
		hasSameDigit := len(findSameDigitGroups(numbers)) > 0
		
		if hasConsecutive {
			consecutiveCount++
		}
		if hasSameDigit {
			sameDigitCount++
		}
	}
	
	consecutivePercentage := 0.0
	sameDigitPercentage := 0.0
	
	if totalDraws > 0 {
		consecutivePercentage = (float64(consecutiveCount) / float64(totalDraws)) * 100
		sameDigitPercentage = (float64(sameDigitCount) / float64(totalDraws)) * 100
	}
	
	return &OverallStats{
		TotalConsecutiveCount: consecutiveCount,
		TotalSameDigitCount:   sameDigitCount,
		ConsecutivePercentage: consecutivePercentage,
		SameDigitPercentage:   sameDigitPercentage,
	}
}

// PrintPatternSummary 分析結果を整理して表示
func PrintPatternSummary(summary *PatternSummary) {
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("ロト7 連番・ゾロ目パターン分析結果")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("調査対象: 過去%d回分\n\n", summary.TotalDraws)
	
	// 全体統計
	stats := summary.OverallStats
	fmt.Println("全体統計")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("連番が含まれる抽選: %d回 (%.2f%%)\n", 
		stats.TotalConsecutiveCount, stats.ConsecutivePercentage)
	fmt.Printf("ゾロ目が含まれる抽選: %d回 (%.2f%%)\n\n", 
		stats.TotalSameDigitCount, stats.SameDigitPercentage)
	
	// 連番パターン詳細
	if len(summary.ConsecutivePatterns) > 0 {
		fmt.Println("連番パターン詳細（出現回数順）")
		fmt.Println(strings.Repeat("-", 60))
		fmt.Printf("%-20s %-8s %-10s %-12s\n", "パターン", "出現回数", "出現率(%)", "最終出現回")
		fmt.Println(strings.Repeat("-", 60))
		
		for _, pattern := range summary.ConsecutivePatterns {
			patternStr := formatPattern(pattern.Pattern)
			lastAppearance := "未出現"
			if pattern.LastAppearance != nil {
				lastAppearance = fmt.Sprintf("%d回前", *pattern.LastAppearance)
			}
			fmt.Printf("%-20s %-8d %-10.2f %-12s\n", 
				patternStr, pattern.Count, pattern.Percentage, lastAppearance)
		}
		fmt.Println()
	} else {
		fmt.Println("連番パターン: 検出されませんでした\n")
	}
	
	// ゾロ目パターン詳細
	if len(summary.SameDigitPatterns) > 0 {
		fmt.Println("ゾロ目パターン詳細（出現回数順）")
		fmt.Println(strings.Repeat("-", 60))
		fmt.Printf("%-20s %-8s %-10s %-12s\n", "パターン", "出現回数", "出現率(%)", "最終出現回")
		fmt.Println(strings.Repeat("-", 60))
		
		for _, pattern := range summary.SameDigitPatterns {
			patternStr := formatPattern(pattern.Pattern)
			lastAppearance := "未出現"
			if pattern.LastAppearance != nil {
				lastAppearance = fmt.Sprintf("%d回前", *pattern.LastAppearance)
			}
			fmt.Printf("%-20s %-8d %-10.2f %-12s\n", 
				patternStr, pattern.Count, pattern.Percentage, lastAppearance)
		}
		fmt.Println()
	} else {
		fmt.Println("ゾロ目パターン: 検出されませんでした\n")
	}
	
	// 分析コメント
	fmt.Println("分析コメント")
	fmt.Println(strings.Repeat("-", 40))
	
	if stats.ConsecutivePercentage > 30.0 {
		fmt.Println("連番：予想以上に頻繁に出現しています")
	} else if stats.ConsecutivePercentage > 15.0 {
		fmt.Println("連番：一定の頻度で出現しています")
	} else {
		fmt.Println("連番：出現頻度は低めです")
	}
	
	if stats.SameDigitPercentage > 20.0 {
		fmt.Println("ゾロ目：かなり頻繁に出現しています")
	} else if stats.SameDigitPercentage > 10.0 {
		fmt.Println("ゾロ目：定期的に出現しています")
	} else {
		fmt.Println("ゾロ目：出現頻度は低めです")
	}
	
	fmt.Println(strings.Repeat("=", 70))
}

// formatPattern パターンを見やすい文字列に変換
func formatPattern(pattern []int) string {
	if len(pattern) == 0 {
		return "-"
	}
	
	var parts []string
	for _, num := range pattern {
		parts = append(parts, fmt.Sprintf("%02d", num))
	}
	return strings.Join(parts, ",")
}
