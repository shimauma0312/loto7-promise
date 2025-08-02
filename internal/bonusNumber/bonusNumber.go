package bonusNumber

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/shimauma0312/loto7-promise/internal/result"
)

// BonusNumberResult ボーナス数字の出現統計を格納する構造体
type BonusNumberResult struct {
	Number     int     `json:"number"`      // ボーナス数字
	Count      int     `json:"count"`       // 出現回数
	TotalDraws int     `json:"total_draws"` // 調査対象の抽選回数
	Percentage float64 `json:"percentage"`  // 出現率（%）
}

// GetBonusNumberTrend 指定した回数分のボーナス数字の出現傾向を分析する
// 引数: searchRange - 過去何回分を調査するか
// 戻り値: ボーナス数字の出現統計スライス、エラー
func GetBonusNumberTrend(searchRange int) ([]*BonusNumberResult, error) {
	// 入力値の妥当性チェック
	if searchRange <= 0 {
		return nil, fmt.Errorf("検索範囲は1以上で指定してください。入力値: %d", searchRange)
	}

	records := result.GetResult(searchRange)
	if records == nil {
		return nil, fmt.Errorf("ロト7の結果データを取得できませんでした")
	}

	// ボーナス数字の出現回数をカウント（1-37の範囲）
	bonusCount := make(map[int]int)
	totalDraws := len(records)

	for _, record := range records {
		// ボーナス数字は8番目と9番目の位置にある（インデックス7,8）
		for i := 7; i < 9 && i < len(record); i++ {
			number, err := strconv.Atoi(record[i])
			if err != nil {
				continue
			}
			
			// ロト7の数字範囲チェック（1-37）
			if number >= 1 && number <= 37 {
				bonusCount[number]++
			}
		}
	}

	// 結果を構造体に格納
	var results []*BonusNumberResult
	for number := 1; number <= 37; number++ {
		count := bonusCount[number]
		percentage := 0.0
		if totalDraws > 0 {
			percentage = (float64(count) / float64(totalDraws)) * 100
		}

		result := &BonusNumberResult{
			Number:     number,
			Count:      count,
			TotalDraws: totalDraws,
			Percentage: percentage,
		}
		results = append(results, result)
	}

	// 出現回数の降順でソート
	sort.Slice(results, func(i, j int) bool {
		return results[i].Count > results[j].Count
	})

	return results, nil
}

// GetSpecificBonusNumberCount 指定したボーナス数字の出現回数を取得
// 引数: targetNumber - 調査したい数字（1-37）, searchRange - 過去何回分を調査するか
// 戻り値: BonusNumberResult構造体のポインタ、エラー
func GetSpecificBonusNumberCount(targetNumber int, searchRange int) (*BonusNumberResult, error) {
	// 入力値の妥当性チェック
	if targetNumber < 1 || targetNumber > 37 {
		return nil, fmt.Errorf("数字は1から37の範囲で指定してください。入力値: %d", targetNumber)
	}
	
	if searchRange <= 0 {
		return nil, fmt.Errorf("検索範囲は1以上で指定してください。入力値: %d", searchRange)
	}

	records := result.GetResult(searchRange)
	if records == nil {
		return nil, fmt.Errorf("ロト7の結果データを取得できませんでした")
	}

	count := 0
	totalDraws := len(records)

	for _, record := range records {
		// ボーナス数字は8番目と9番目の位置にある（インデックス7,8）
		for i := 7; i < 9 && i < len(record); i++ {
			number, err := strconv.Atoi(record[i])
			if err != nil {
				continue
			}
			
			if number == targetNumber {
				count++
			}
		}
	}

	percentage := 0.0
	if totalDraws > 0 {
		percentage = (float64(count) / float64(totalDraws)) * 100
	}

	result := &BonusNumberResult{
		Number:     targetNumber,
		Count:      count,
		TotalDraws: totalDraws,
		Percentage: percentage,
	}

	return result, nil
}

// PrintBonusNumberTrend ボーナス数字の出現傾向を表示
// 引数: results - 表示したいボーナス数字の統計データ, displayCount - 表示する上位件数
func PrintBonusNumberTrend(results []*BonusNumberResult, displayCount int) {
	fmt.Println("ボーナス数字の出現傾向分析")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("過去%d回分の抽選データから分析\n", results[0].TotalDraws)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("%-6s %-8s %-8s %-10s\n", "数字", "出現回数", "調査回数", "出現率(%)")
	fmt.Println(strings.Repeat("-", 60))
	
	// 指定された件数まで表示
	for i := 0; i < displayCount && i < len(results); i++ {
		result := results[i]
		fmt.Printf("%-6d %-8d %-8d %-10.2f\n", 
			result.Number, result.Count, result.TotalDraws, result.Percentage)
	}
	fmt.Println(strings.Repeat("=", 60))
}

// PrintSpecificBonusNumberResult 特定のボーナス数字の結果を表示
// 引数: result - 表示したいボーナス数字の統計データ
func PrintSpecificBonusNumberResult(result *BonusNumberResult) {
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("ボーナス数字 %d の出現統計\n", result.Number)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("出現回数: %d回\n", result.Count)
	fmt.Printf("調査対象: 過去%d回分\n", result.TotalDraws)
	fmt.Printf("出現率: %.2f%%\n", result.Percentage)
	
	// 出現頻度の評価コメント
	switch {
	case result.Percentage >= 8.0:
		fmt.Println("ボーナス数字として高頻度で出現")
	case result.Percentage >= 6.0:
		fmt.Println("平均的な出現頻度")
	case result.Percentage >= 4.0:
		fmt.Println("やや低い出現頻度")
	case result.Percentage >= 2.0:
		fmt.Println("低い出現頻度")
	default:
		fmt.Println("ほとんど出現していない")
	}
	fmt.Println(strings.Repeat("=", 50))
}
