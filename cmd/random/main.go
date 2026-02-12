package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/shimauma0312/loto7-promise/internal/random"
)

func main() {
	var (
		history = flag.Int("history", 50, "分析する過去の抽選回数")
		count   = flag.Int("count", 5, "生成するランダム組み合わせの数")
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
	fmt.Printf("過去%d回分のデータを分析中...\n", *history)
	if err := engine.LoadData(*history); err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	// 詳細表示（オプション）
	if *verbose {
		showPositionRanges(engine)
	}

	// ランダム組み合わせを生成
	fmt.Printf("\n%d組のランダム組み合わせを生成中...\n", *count)
	combinations, err := engine.GenerateMultipleRandomCombinations(*count)
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

}

func showHelp() {
	fmt.Println("ロト7ランダム推薦ツール")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  random [オプション]")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println("  -history int    分析する過去の抽選回数 (デフォルト: 100)")
	fmt.Println("  -count int      生成するランダム組み合わせの数 (デフォルト: 5)")
	fmt.Println("  -verbose        各位置の出現範囲を表示")
	fmt.Println("  -help           このヘルプを表示")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  random -history 150 -count 3")
	fmt.Println("  random -verbose")
}

func showPositionRanges(engine *random.RandomEngine) {
	fmt.Println("\n=== 各位置の出現範囲 ===")
	ranges := engine.GetPositionRanges()
	for _, r := range ranges {
		fmt.Printf("位置 %d: %2d ~ %2d (出現数字: %d個)\n",
			r.Position, r.Min, r.Max, len(r.Numbers))
		
		// 上位5個の頻出数字を表示
		if len(r.Frequency) > 0 {
			type numFreq struct {
				num  int
				freq int
			}
			freqs := make([]numFreq, 0, len(r.Frequency))
			for num, freq := range r.Frequency {
				freqs = append(freqs, numFreq{num, freq})
			}
			// 頻度の高い順にソート
			sort.Slice(freqs, func(i, j int) bool {
				if freqs[i].freq == freqs[j].freq {
					return freqs[i].num < freqs[j].num
				}
				return freqs[i].freq > freqs[j].freq
			})
			
			// 上位5個または全て表示
			displayCount := 5
			if len(freqs) < displayCount {
				displayCount = len(freqs)
			}
			fmt.Print("  上位頻出: ")
			for i := 0; i < displayCount; i++ {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%d(%d回)", freqs[i].num, freqs[i].freq)
			}
			fmt.Println()
		}
	}
}
