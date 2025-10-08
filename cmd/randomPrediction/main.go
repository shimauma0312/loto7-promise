package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shimauma0312/loto7-promise/internal/randomPrediction"
)

func main() {
	// コマンドラインフラグの定義
	var (
		count     = flag.Int("count", 1, "生成する予想番号の組数")
		useTrend  = flag.Bool("trend", true, "過去データの傾向を考慮するかどうか")
		range_    = flag.Int("range", 100, "傾向分析に使用する過去回数")
		random    = flag.Bool("random", false, "完全ランダム生成（傾向を無視）")
		details   = flag.Bool("details", false, "分析情報表示")
		verbose   = flag.Bool("verbose", false, "詳細ログ表示")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "ロト7予想番号生成ツール\n\n")
		fmt.Fprintf(os.Stderr, "使用方法:\n")
		fmt.Fprintf(os.Stderr, "  %s [オプション]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "オプション:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n例:\n")
		fmt.Fprintf(os.Stderr, "  %s                                   # 統計分析で1組生成\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -count 5 -range 200              # 過去200回の傾向で5組生成\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -random -count 3                 # 完全ランダムで3組生成\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -details -verbose                # 分析情報付きで生成\n", os.Args[0])
	}

	flag.Parse()

	// 引数の妥当性チェック
	if *count <= 0 {
		fmt.Println("エラー: 生成数は1以上で指定してください")
		os.Exit(1)
	}

	if *range_ <= 0 {
		fmt.Println("エラー: 分析範囲は1以上で指定してください")
		os.Exit(1)
	}

	// randomフラグが指定された場合は傾向を無視
	if *random {
		*useTrend = false
	}

	// 詳細ログ表示
	if *verbose {
		fmt.Printf("🔧 設定: 生成数=%d, 分析範囲=%d, 傾向分析=%t, 詳細表示=%t\n", 
			*count, *range_, *useTrend, *details)
		if *useTrend {
			fmt.Println("統計分析モードで実行します...")
		} else {
			fmt.Println("ランダムモードで実行します...")
		}
		fmt.Println()
	}

	// 単一の予想番号生成
	if *count == 1 {
		result, err := randomPrediction.GenerateRandomPrediction(*range_, *useTrend)
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
			os.Exit(1)
		}

		randomPrediction.PrintPredictionResultWithDetails(result, *details)
		return
	}

	// 複数の予想番号生成
	results, err := randomPrediction.GenerateMultiplePredictions(*count, *range_, *useTrend)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	randomPrediction.PrintMultiplePredictionsWithDetails(results, *details)
}
