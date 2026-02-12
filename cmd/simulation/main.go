package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/shimauma0312/loto7-promise/internal/simulation"
)

func main() {
	var (
		history      = flag.Int("history", 100, "分析に使用する過去の抽選回数")
		simulations  = flag.Int("simulations", 1, "シミュレーション実行回数（統計取得用）")
		userNumbers  = flag.String("numbers", "", "ユーザーの選択数字（カンマ区切り、例: 1,5,10,15,20,25,30）")
		autoGenerate = flag.Bool("auto", false, "ユーザーの数字を自動生成")
		help         = flag.Bool("help", false, "ヘルプを表示")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// シミュレーションエンジンを作成
	engine := simulation.NewSimulationEngine()

	// データを読み込み
	fmt.Printf("過去%d回分のデータを読み込み中...\n", *history)
	if err := engine.LoadData(*history); err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	// ユーザーの数字を準備
	var userNums []int
	var err error

	if *autoGenerate {
		// 自動生成
		fmt.Println("ユーザーの数字を自動生成中...")
		userNums, err = engine.GenerateUserNumbers()
		if err != nil {
			fmt.Printf("数字生成エラー: %v\n", err)
			os.Exit(1)
		}
	} else if *userNumbers != "" {
		// ユーザー指定
		userNums, err = parseUserNumbers(*userNumbers)
		if err != nil {
			fmt.Printf("数字の解析エラー: %v\n", err)
			os.Exit(1)
		}
	} else {
		// デフォルト: 自動生成
		fmt.Println("ユーザーの数字を自動生成中...")
		userNums, err = engine.GenerateUserNumbers()
		if err != nil {
			fmt.Printf("数字生成エラー: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("\nユーザーの選択数字: %s\n", formatNumbers(userNums))
	fmt.Println(strings.Repeat("=", 60))

	// シミュレーション実行
	if *simulations == 1 {
		// 1回のシミュレーション
		runSingleSimulation(engine, userNums)
	} else {
		// 複数回のシミュレーション
		runMultipleSimulations(engine, userNums, *simulations)
	}
}

func runSingleSimulation(engine *simulation.SimulationEngine, userNumbers []int) {
	fmt.Println("\n1等が当選するまでシミュレーション実行中...")
	fmt.Println("（これには時間がかかる場合があります）")

	result, err := engine.RunSimulation(userNumbers)
	if err != nil {
		fmt.Printf("シミュレーションエラー: %v\n", err)
		os.Exit(1)
	}

	// 結果表示
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("=== シミュレーション結果 ===")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("抽選回数: %d回\n", result.DrawCount)
	fmt.Printf("当選番号: %s\n", formatNumbers(result.WinningNumbers))
	fmt.Printf("所要時間: %v\n", result.Duration)
	fmt.Printf("\n推定金額: ¥%s (1口300円として)\n", formatMoney(result.DrawCount*300))
	fmt.Println(strings.Repeat("=", 60))
}

func runMultipleSimulations(engine *simulation.SimulationEngine, userNumbers []int, count int) {
	fmt.Printf("\n%d回のシミュレーションを実行中...\n", count)
	fmt.Println("（これには時間がかかる場合があります）")

	stats, err := engine.RunMultipleSimulations(count, userNumbers)
	if err != nil {
		fmt.Printf("シミュレーションエラー: %v\n", err)
		os.Exit(1)
	}

	// 統計結果表示
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("=== シミュレーション統計 ===")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("実行回数: %d回\n", stats.Simulations)
	fmt.Printf("最小抽選回数: %d回\n", stats.MinDraws)
	fmt.Printf("最大抽選回数: %d回\n", stats.MaxDraws)
	fmt.Printf("平均抽選回数: %.0f回\n", stats.AvgDraws)
	fmt.Printf("中央値: %d回\n", stats.MedianDraws)
	fmt.Printf("総所要時間: %v\n", stats.TotalDuration)
	fmt.Printf("\n平均推定金額: ¥%s (1口300円として)\n", formatMoney(int(stats.AvgDraws)*300))
	fmt.Println(strings.Repeat("=", 60))

	// 詳細結果（最初の5件のみ）
	if len(stats.Results) > 0 {
		fmt.Println("\n=== 詳細結果（最初の5件）===")
		displayCount := 5
		if len(stats.Results) < displayCount {
			displayCount = len(stats.Results)
		}
		for i := 0; i < displayCount; i++ {
			r := stats.Results[i]
			fmt.Printf("[%d] 抽選回数: %d回, 所要時間: %v\n", i+1, r.DrawCount, r.Duration)
		}
	}
}

func parseUserNumbers(input string) ([]int, error) {
	parts := strings.Split(input, ",")
	if len(parts) != 7 {
		return nil, fmt.Errorf("数字は7個指定してください")
	}

	numbers := make([]int, 7)
	for i, part := range parts {
		num, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("数字の解析に失敗: %s", part)
		}
		if num < 1 || num > 37 {
			return nil, fmt.Errorf("数字は1-37の範囲である必要があります: %d", num)
		}
		numbers[i] = num
	}

	return numbers, nil
}

func formatNumbers(numbers []int) string {
	parts := make([]string, len(numbers))
	for i, num := range numbers {
		parts[i] = fmt.Sprintf("%02d", num)
	}
	return strings.Join(parts, " - ")
}

func formatMoney(amount int) string {
	str := strconv.Itoa(amount)
	result := ""
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(c)
	}
	return result
}

func showHelp() {
	fmt.Println("ロト7 1等当選シミュレーター")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  simulation [オプション]")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println("  -history int        分析に使用する過去の抽選回数 (デフォルト: 100)")
	fmt.Println("  -simulations int    シミュレーション実行回数 (デフォルト: 1)")
	fmt.Println("  -numbers string     ユーザーの選択数字（カンマ区切り）")
	fmt.Println("                      例: -numbers 1,5,10,15,20,25,30")
	fmt.Println("  -auto               ユーザーの数字を自動生成（デフォルト）")
	fmt.Println("  -help               このヘルプを表示")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  # 自動生成された数字で1回シミュレーション")
	fmt.Println("  simulation")
	fmt.Println()
	fmt.Println("  # 指定した数字で1回シミュレーション")
	fmt.Println("  simulation -numbers 1,5,10,15,20,25,30")
	fmt.Println()
	fmt.Println("  # 10回シミュレーションして統計を取る")
	fmt.Println("  simulation -simulations 10")
	fmt.Println()
	fmt.Println("  # 指定した数字で10回シミュレーション")
	fmt.Println("  simulation -numbers 1,5,10,15,20,25,30 -simulations 10")
}
