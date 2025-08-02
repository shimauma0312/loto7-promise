package numberCount

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shimauma0312/loto7-promise/internal/result"
)

type NumberCountResult struct {
	Number     int   `json:"number"`      // 指定された数字
	Count      int   `json:"count"`       // 出現回数
	TotalDraws int   `json:"total_draws"` // 調査対象の抽選回数
	Percentage float64 `json:"percentage"` // 出現率（%）
}

// CountSpecificNumber 指定した数字が過去何回出現したかを集計する
// 引数: targetNumber - 調査したい数字（1-37）, searchRange - 過去何回分を調査するか
// 戻り値: NumberCountResult構造体のポインタ、エラー
func CountSpecificNumber(targetNumber int, searchRange int) (*NumberCountResult, error) {
	// 入力値の妥当性チェック - 当たり前だけど範囲外はダメよ
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

	count := 0 // 出現回数カウンター
	totalDraws := len(records) // 実際に取得できた抽選回数

	for _, record := range records {
		for i := 0; i < 7 && i < len(record); i++ {
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

	result := &NumberCountResult{
		Number:     targetNumber,
		Count:      count,
		TotalDraws: totalDraws,
		Percentage: percentage,
	}

	return result, nil
}

// CountMultipleNumbers 複数の数字の出現回数を一括で調査する
// 引数: targetNumbers - 調査したい数字のスライス, searchRange - 過去何回分を調査するか
// 戻り値: NumberCountResult構造体のスライス、エラー
func CountMultipleNumbers(targetNumbers []int, searchRange int) ([]*NumberCountResult, error) {
	results := make([]*NumberCountResult, 0, len(targetNumbers))
	
	// 各数字について個別に集計 - 力技だけど確実よ
	for _, number := range targetNumbers {
		result, err := CountSpecificNumber(number, searchRange)
		if err != nil {
			// エラーが発生した数字はスキップして続行 - 完璧を求めない
			fmt.Printf("数字 %d の集計でエラー: %v\n", number, err)
			continue
		}
		results = append(results, result)
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("有効な結果が取得できませんでした")
	}
	
	return results, nil
}

// 引数: result - 表示したい結果データ
func PrintResult(result *NumberCountResult) {
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("数字 %d の出現統計\n", result.Number)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("出現回数: %d回\n", result.Count)
	fmt.Printf("調査対象: 過去%d回分\n", result.TotalDraws)
	fmt.Printf("出現率: %.2f%%\n", result.Percentage)
	
	switch {
	case result.Percentage >= 20.0:
		fmt.Println("よく出てる")
	case result.Percentage >= 15.0:
		fmt.Println("まあまあ出てる")
	case result.Percentage >= 10.0:
		fmt.Println("平均的な出現率")
	case result.Percentage >= 5.0:
		fmt.Println("あまり出てない")
	default:
		fmt.Println("ゴミ")
	}
	fmt.Println(strings.Repeat("=", 50))
}

// PrintMultipleResults 複数の結果を一覧表示する
// 引数: results - 表示したい結果データのスライス
func PrintMultipleResults(results []*NumberCountResult) {
	fmt.Println("📋 複数数字の出現統計一覧")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("%-6s %-8s %-8s %-10s\n", "数字", "出現回数", "調査回数", "出現率(%)")
	fmt.Println(strings.Repeat("-", 60))
	
	for _, result := range results {
		fmt.Printf("%-6d %-8d %-8d %-10.2f\n", 
			result.Number, result.Count, result.TotalDraws, result.Percentage)
	}
	fmt.Println(strings.Repeat("=", 60))
}
