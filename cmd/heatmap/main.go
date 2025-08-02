package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/shimauma0312/loto7-promise/internal/heatmap"
)

func main() {
	var (
		// コマンドライン引数の定義
		searchRange = flag.Int("range", 50, "過去何回分の抽選結果を調査するか（デフォルト: 50）")
		outputJSON  = flag.Bool("json", false, "JSON形式で出力するか（デフォルト: false）")
		showNumber  = flag.Int("number", 0, "特定の数字の位置別詳細を表示（1-37、0で全体表示）")
		showPos     = flag.Int("position", 0, "特定の位置の全数字詳細を表示（1-7、0で全体表示）")
		help        = flag.Bool("help", false, "ヘルプを表示")
	)
	
	flag.Parse()

	// ヘルプ表示
	if *help {
		showHelp()
		return
	}

	// 入力値の妥当性チェック
	if *searchRange <= 0 {
		fmt.Fprintf(os.Stderr, "エラー: 検索範囲は1以上で指定してください\n")
		os.Exit(1)
	}

	if *showNumber < 0 || *showNumber > 37 {
		fmt.Fprintf(os.Stderr, "エラー: 数字は0-37の範囲で指定してください（0は全体表示）\n")
		os.Exit(1)
	}

	if *showPos < 0 || *showPos > 7 {
		fmt.Fprintf(os.Stderr, "エラー: 位置は0-7の範囲で指定してください（0は全体表示）\n")
		os.Exit(1)
	}

	// ヒートマップデータを生成
	fmt.Printf("ロト7抽選結果を取得中... (過去%d回分)\n", *searchRange)
	summary, err := heatmap.GenerateHeatmap(*searchRange)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 出力処理
	if *outputJSON {
		// JSON形式で出力
		outputJSONData(summary, *showNumber, *showPos)
	} else {
		// テキスト形式で出力
		outputText(summary, *showNumber, *showPos)
	}
}

// showHelp ヘルプメッセージを表示
func showHelp() {
	fmt.Println("ロト7位置別出現ヒートマップ生成ツール")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  heatmap [オプション]")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println("  -range int       過去何回分の抽選結果を調査するか (デフォルト: 50)")
	fmt.Println("  -json            JSON形式で出力する")
	fmt.Println("  -number int      特定の数字(1-37)の位置別詳細を表示 (0で全体表示)")
	fmt.Println("  -position int    特定の位置(1-7)の全数字詳細を表示 (0で全体表示)")
	fmt.Println("  -help            このヘルプを表示")
	fmt.Println()
	fmt.Println("使用例:")
	fmt.Println("  heatmap                    # デフォルト設定で全体ヒートマップを表示")
	fmt.Println("  heatmap -range 100         # 過去100回分でヒートマップを生成")
	fmt.Println("  heatmap -number 7          # 数字7の位置別詳細を表示")
	fmt.Println("  heatmap -position 1        # 1番目の位置の全数字詳細を表示")
	fmt.Println("  heatmap -json -range 30    # 過去30回分をJSON形式で出力")
}

// outputText テキスト形式で出力
func outputText(summary *heatmap.HeatmapSummary, showNumber, showPos int) {
	if showNumber > 0 {
		// 特定の数字の詳細表示
		showNumberDetails(summary, showNumber)
	} else if showPos > 0 {
		// 特定の位置の詳細表示
		showPositionDetails(summary, showPos)
	} else {
		// 全体のヒートマップ表示
		summary.PrintHeatmapTable()
	}
}

// outputJSONData JSON形式で出力
func outputJSONData(summary *heatmap.HeatmapSummary, showNumber, showPos int) {
	var data interface{}
	
	if showNumber > 0 {
		// 特定の数字のデータを取得
		positionData, err := summary.GetPositionData(showNumber)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		data = map[string]interface{}{
			"number":        showNumber,
			"position_data": positionData,
			"search_range":  summary.SearchRange,
			"actual_draws":  summary.ActualDraws,
		}
	} else if showPos > 0 {
		// 特定の位置のデータを取得
		numberData, err := summary.GetNumbersByPosition(showPos)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		data = map[string]interface{}{
			"position":     showPos,
			"number_data":  numberData,
			"search_range": summary.SearchRange,
			"actual_draws": summary.ActualDraws,
		}
	} else {
		// 全体データ
		data = summary
	}
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON出力エラー: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println(string(jsonData))
}

// showNumberDetails 特定の数字の詳細を表示
func showNumberDetails(summary *heatmap.HeatmapSummary, number int) {
	positionData, err := summary.GetPositionData(number)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n=== 数字 %d の位置別出現詳細（過去%d回分） ===\n", number, summary.SearchRange)
	fmt.Printf("実際の抽選回数: %d回\n\n", summary.ActualDraws)

	totalCount := 0
	fmt.Printf("位置\t出現回数\t出現率\t\t位置内率\n")
	fmt.Printf("----\t--------\t------\t\t--------\n")
	
	for _, data := range positionData {
		fmt.Printf("%d\t%d\t\t%.2f%%\t\t%.2f%%\n", 
			data.Position, data.Count, data.Percentage, data.PositionRatio)
		totalCount += data.Count
	}
	
	fmt.Printf("\n合計出現回数: %d回\n", totalCount)
	if summary.ActualDraws > 0 {
		overallRate := (float64(totalCount) / float64(summary.ActualDraws)) * 100
		fmt.Printf("全体出現率: %.2f%%\n", overallRate)
	}

	// 最も出現頻度の高い位置を表示
	maxCount := 0
	var bestPositions []int
	for _, data := range positionData {
		if data.Count > maxCount {
			maxCount = data.Count
			bestPositions = []int{data.Position}
		} else if data.Count == maxCount {
			bestPositions = append(bestPositions, data.Position)
		}
	}
	
	if maxCount > 0 {
		fmt.Printf("最頻出位置: %v (出現回数: %d回)\n", bestPositions, maxCount)
	}
}

// showPositionDetails 特定の位置の詳細を表示
func showPositionDetails(summary *heatmap.HeatmapSummary, position int) {
	numberData, err := summary.GetNumbersByPosition(position)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n=== 位置 %d の全数字出現詳細（過去%d回分） ===\n", position, summary.SearchRange)
	fmt.Printf("実際の抽選回数: %d回\n\n", summary.ActualDraws)

	// 出現回数でソート用のスライスを作成
	type NumberCount struct {
		Number int
		Count  int
		Percentage float64
	}
	
	var numbers []NumberCount
	for _, data := range numberData {
		numbers = append(numbers, NumberCount{
			Number: data.Number,
			Count:  data.Count,
			Percentage: data.Percentage,
		})
	}

	// 出現回数で降順ソート
	for i := 0; i < len(numbers)-1; i++ {
		for j := i + 1; j < len(numbers); j++ {
			if numbers[i].Count < numbers[j].Count {
				numbers[i], numbers[j] = numbers[j], numbers[i]
			}
		}
	}

	fmt.Printf("数字\t出現回数\t出現率\n")
	fmt.Printf("----\t--------\t------\n")
	
	for _, num := range numbers {
		if num.Count > 0 { // 出現回数が0より大きいもののみ表示
			fmt.Printf("%2d\t%d\t\t%.2f%%\n", num.Number, num.Count, num.Percentage)
		}
	}

	// 上位5位を表示
	fmt.Printf("\n=== TOP 5 ===\n")
	for i := 0; i < 5 && i < len(numbers) && numbers[i].Count > 0; i++ {
		fmt.Printf("%d位: 数字%d (%d回, %.2f%%)\n", 
			i+1, numbers[i].Number, numbers[i].Count, numbers[i].Percentage)
	}
}
