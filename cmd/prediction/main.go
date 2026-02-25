package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shimauma0312/loto7-promise/internal/prediction"
)

func main() {
	var (
		history   = flag.Int("history", 100, "分析する過去の抽選回数")
		count     = flag.Int("count", 5, "生成する予測組み合わせの数")
		verbose   = flag.Bool("verbose", false, "統計分析情報を表示")
		scores    = flag.Bool("scores", false, "各組み合わせのスコア内訳を表示")
		help      = flag.Bool("help", false, "ヘルプを表示")

		// フィルター重みオプション
		wHotCold   = flag.Float64("w-hot-cold", 1.0, "ホット/コールド分析の重み")
		wParity    = flag.Float64("w-parity", 1.0, "パリティバランスの重み")
		wSize      = flag.Float64("w-size", 1.0, "大小バランスの重み")
		wSum       = flag.Float64("w-sum", 1.0, "合計値適正範囲の重み")
		wTens      = flag.Float64("w-tens", 1.0, "十の位グループの重み")
		wLastDigit = flag.Float64("w-last-digit", 1.0, "一の位種類数の重み")
		wPull      = flag.Float64("w-pull", 0.8, "引っ張り・斜め数字の重み")
		wBonus     = flag.Float64("w-bonus", 0.8, "ボーナス周辺の重み")
		wInterval  = flag.Float64("w-interval", 0.6, "当選間隔の重み")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// 設定を組み立て
	cfg := prediction.DefaultConfig()
	cfg.Count = *count
	cfg.History = *history
	cfg.Weights = prediction.FilterWeights{
		HotCold:   *wHotCold,
		Parity:    *wParity,
		SizeBal:   *wSize,
		SumRange:  *wSum,
		TensGroup: *wTens,
		LastDigit: *wLastDigit,
		Pull:      *wPull,
		Bonus:     *wBonus,
		Interval:  *wInterval,
	}

	// エンジンを作成してデータ読み込み
	engine := prediction.NewEngine()
	fmt.Printf("過去%d回分のデータを分析中...\n", *history)
	if err := engine.LoadData(*history); err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	// 予測生成
	fmt.Printf("%d口の予測組み合わせを生成中...\n", *count)
	result, err := engine.Generate(cfg)
	if err != nil {
		fmt.Printf("予測生成エラー: %v\n", err)
		os.Exit(1)
	}

	// 統計分析情報を表示（-verbose）
	if *verbose {
		info := result.AnalysisInfo
		fmt.Println()
		fmt.Println("=== 統計分析情報 ===")
		fmt.Printf("分析回数    : %d回\n", info.AnalyzedDraws)
		fmt.Printf("前回数字    : %v\n", info.LastDrawNumbers)
		fmt.Printf("前回合計値  : %d\n", info.LastDrawSum)
		fmt.Printf("合計トレンド: %s\n", info.SumTrend)
		fmt.Printf("ホット数字  : %v\n", info.HotNumbers)
		fmt.Printf("コールド数字: %v\n", info.ColdNumbers)
	}

	// 結果を表示
	fmt.Println()
	fmt.Println("=== 統計分析予測番号 ===")
	for _, combo := range result.Combinations {
		fmt.Printf("[%d] ", combo.ID)
		for j, num := range combo.Numbers {
			if j > 0 {
				fmt.Print(" - ")
			}
			fmt.Printf("%02d", num)
		}
		if *scores {
			fmt.Printf("  (スコア: %.3f)", combo.TotalScore)
		}
		fmt.Println()
	}

	// スコア内訳を表示（-scores）
	if *scores {
		fmt.Println()
		fmt.Println("=== スコア内訳 ===")
		for _, combo := range result.Combinations {
			d := combo.ScoreDetail
			fmt.Printf("[%d] 合計:%.3f  ホット/コールド:%.2f  パリティ:%.2f  大小:%.2f  合計値:%.2f  十の位:%.2f  一の位:%.2f  引っ張り:%.2f  ボーナス:%.2f  間隔:%.2f\n",
				combo.ID, combo.TotalScore,
				d.HotColdScore, d.ParityScore, d.SizeScore, d.SumScore,
				d.TensScore, d.LastDigitScore, d.PullScore, d.BonusScore, d.IntervalScore,
			)
		}
	}
}

func showHelp() {
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  prediction [オプション]")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println("  -history int      分析する過去の抽選回数 (デフォルト: 100)")
	fmt.Println("  -count int        生成する予測組み合わせの数 (デフォルト: 5)")
	fmt.Println("  -verbose          統計分析情報を表示")
	fmt.Println("  -scores           各組み合わせのスコア内訳を表示")
	fmt.Println("  -help             このヘルプを表示")
	fmt.Println()
	fmt.Println("フィルター重みオプション (0.0〜2.0):")
	fmt.Println("  -w-hot-cold       ホット/コールド分析 (デフォルト: 1.0)")
	fmt.Println("  -w-parity         パリティバランス   (デフォルト: 1.0)")
	fmt.Println("  -w-size           大小バランス       (デフォルト: 1.0)")
	fmt.Println("  -w-sum            合計値適正範囲     (デフォルト: 1.0)")
	fmt.Println("  -w-tens           十の位グループ     (デフォルト: 1.0)")
	fmt.Println("  -w-last-digit     一の位種類数       (デフォルト: 1.0)")
	fmt.Println("  -w-pull           引っ張り・斜め数字 (デフォルト: 0.8)")
	fmt.Println("  -w-bonus          ボーナス周辺       (デフォルト: 0.8)")
	fmt.Println("  -w-interval       当選間隔           (デフォルト: 0.6)")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  prediction -history 200 -count 3 -verbose -scores")
	fmt.Println("  prediction -w-hot-cold 1.5 -w-pull 1.2 -count 3")
}
