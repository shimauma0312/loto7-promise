package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shimauma0312/loto7-promise/internal/bonusNumber"
)

func main() {
	// コマンドラインフラグの定義
	var (
		showTrend = flag.Bool("trend", false, "ボーナス数字の出現傾向を表示")
		topCount  = flag.Int("top", 10, "傾向表示時の上位表示件数")
		number    = flag.Int("number", 0, "特定のボーナス数字の出現回数を調査")
		range_    = flag.Int("range", 100, "調査対象の過去回数")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "ロト7ボーナス数字分析ツール\n\n")
		fmt.Fprintf(os.Stderr, "使用方法:\n")
		fmt.Fprintf(os.Stderr, "  %s [オプション]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "オプション:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n例:\n")
		fmt.Fprintf(os.Stderr, "  %s -trend -top 15 -range 200        # 過去200回の上位15位までの傾向表示\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -number 7 -range 50              # ボーナス数字7の過去50回分の出現回数\n", os.Args[0])
	}

	flag.Parse()

	// 引数の妥当性チェック
	if *range_ <= 0 {
		fmt.Println("エラー: 調査範囲は1以上で指定してください")
		os.Exit(1)
	}

	if *number != 0 {
		// 特定のボーナス数字の出現回数を調査
		if *number < 1 || *number > 37 {
			fmt.Println("エラー: 数字は1から37の範囲で指定してください")
			os.Exit(1)
		}

		result, err := bonusNumber.GetSpecificBonusNumberCount(*number, *range_)
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
			os.Exit(1)
		}

		bonusNumber.PrintSpecificBonusNumberResult(result)
		return
	}

	if *showTrend {
		// ボーナス数字の出現傾向を表示
		if *topCount <= 0 {
			fmt.Println("エラー: 表示件数は1以上で指定してください")
			os.Exit(1)
		}

		results, err := bonusNumber.GetBonusNumberTrend(*range_)
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
			os.Exit(1)
		}

		bonusNumber.PrintBonusNumberTrend(results, *topCount)
		return
	}

	// 引数が指定されていない場合はデフォルトで傾向表示
	fmt.Println("ボーナス数字の出現傾向を表示します...")
	results, err := bonusNumber.GetBonusNumberTrend(*range_)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	bonusNumber.PrintBonusNumberTrend(results, 10)
}
