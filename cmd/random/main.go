package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shimauma0312/loto7-promise/internal/random"
)

func main() {
	var (
		count   = flag.Int("count", 100, "分析する過去の抽選回数")
		maxRand = flag.Int("max", 5, "生成するランダム組み合わせの数")
		verbose = flag.Bool("verbose", false, "各位置の出現範囲を表示")
		help    = flag.Bool("help", false, "ヘルプを表示")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// ランダムエンジンを作成
	engine := random.NewRandomEngine()

	// データを読み込み
	fmt.Printf("過去%d回分のデータを分析中...\n", *count)
	if err := engine.LoadData(*count); err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	// 詳細表示（オプション）
	if *verbose {
		showPositionRanges(engine)
	}

	// ランダム組み合わせを生成
	fmt.Printf("\n%d組のランダム組み合わせを生成中...\n", *maxRand)
	combinations, err := engine.GenerateMultipleRandomCombinations(*maxRand)
	if err != nil {
		fmt.Printf("生成エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を表示
	fmt.Println("\n=== ランダム推薦番号 ===")
	for i, combination := range combinations {
		fmt.Printf("[%d] ", i+1)
		for j, num := range combination {
			if j > 0 {
				fmt.Print(" - ")
			}
			fmt.Printf("%02d", num)
		}
		fmt.Println()
	}

	// ロジックの説明
	showRandomLogic()
}

func showHelp() {
	fmt.Println("ロト7ランダム推薦ツール")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  random [オプション]")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println("  -count int      分析する過去の抽選回数 (デフォルト: 100)")
	fmt.Println("  -max int        生成するランダム組み合わせの数 (デフォルト: 5)")
	fmt.Println("  -verbose        各位置の出現範囲を表示")
	fmt.Println("  -help           このヘルプを表示")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  random -count 150 -max 3")
	fmt.Println("  random -verbose")
}

func showPositionRanges(engine *random.RandomEngine) {
	fmt.Println("\n=== 各位置の出現範囲 ===")
	ranges := engine.GetPositionRanges()
	for _, r := range ranges {
		fmt.Printf("位置 %d: %2d ~ %2d (出現数字: %d個)\n",
			r.Position, r.Min, r.Max, len(r.Numbers))
	}
}

func showRandomLogic() {
	fmt.Println()
	fmt.Println("=== ランダム生成ロジック ===")
	fmt.Println("各数字位置で過去に出現した範囲内から")
	fmt.Println("ランダムに数字を選択します。")
	fmt.Println("未出現の数字は自動的に除外されます。")
}
