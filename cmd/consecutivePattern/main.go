package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shimauma0312/loto7-promise/internal/consecutivePattern"
)

func main() {
	var (
		searchRange = flag.Int("range", 100, "過去何回分を調査するか（デフォルト: 100回）")
		help        = flag.Bool("help", false, "ヘルプを表示")
	)
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "使用方法:\n")
		fmt.Fprintf(os.Stderr, "  %s [-range=回数]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "説明:\n")
		fmt.Fprintf(os.Stderr, "  ロト7の過去データから連番・ゾロ目パターンの出現傾向を分析します。\n")
		fmt.Fprintf(os.Stderr, "  連番: 11,12,13 のような連続した数字\n")
		fmt.Fprintf(os.Stderr, "  ゾロ目: 01,11,21 のような一の位が同じ数字\n\n")
		fmt.Fprintf(os.Stderr, "例:\n")
		fmt.Fprintf(os.Stderr, "  %s                    # 過去100回分を分析\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -range=200         # 過去200回分を分析\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -range=50          # 過去50回分を分析\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nオプション:\n")
		flag.PrintDefaults()
	}
	
	flag.Parse()

	// ヘルプ表示
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// 入力値チェック
	if *searchRange <= 0 {
		fmt.Fprintf(os.Stderr, "エラー: 検索範囲は1以上で指定してください\n")
		os.Exit(1)
	}

	fmt.Printf("ロト7 連番・ゾロ目パターン分析を開始します...\n")
	fmt.Printf("調査対象: 過去%d回分\n\n", *searchRange)

	// パターン分析を実行
	summary, err := consecutivePattern.AnalyzeConsecutivePatterns(*searchRange)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を表示
	consecutivePattern.PrintPatternSummary(summary)

	fmt.Println("分析完了！")
}
