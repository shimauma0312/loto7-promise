package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/shimauma0312/loto7-promise/internal/numberCount"
)

func main() {
	var (
		number = flag.String("number", "", "調査したい数字（1-37）複数の場合はカンマ区切り（例: 1,7,23）")
		range_ = flag.Int("range", 50, "過去何回分を調査するか（デフォルト: 50回）")
		help   = flag.Bool("help", false, "ヘルプを表示")
	)
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "使用方法:\n")
		fmt.Fprintf(os.Stderr, "  %s -number=数字 [-range=回数]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "例:\n")
		fmt.Fprintf(os.Stderr, "  %s -number=7 -range=100    # 数字7が過去100回で何回出たか\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -number=1,7,23         # 複数数字を一括調査\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nオプション:\n")
		flag.PrintDefaults()
	}
	
	flag.Parse()

	// ヘルプ
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *number == "" {
		fmt.Fprintf(os.Stderr, "エラー: 数字を指定してください\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if *range_ <= 0 {
		fmt.Fprintf(os.Stderr, "エラー: 検索範囲は1以上で指定してください\n")
		os.Exit(1)
	}

	numberStrings := strings.Split(*number, ",")
	numbers := make([]int, 0, len(numberStrings))
	
	for _, numStr := range numberStrings {
		numStr = strings.TrimSpace(numStr)
		if numStr == "" {
			continue
		}
		
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: '%s' は有効な数字ではありません\n", numStr)
			os.Exit(1)
		}
		
		// 数字範囲チェック
		if num < 1 || num > 37 {
			fmt.Fprintf(os.Stderr, "エラー: 数字は1から37の範囲で指定してください（指定値: %d）\n", num)
			os.Exit(1)
		}
		
		numbers = append(numbers, num)
	}

	fmt.Printf("ロト7データを分析中...\n")
	fmt.Printf("調査対象: 過去%d回分\n", *range_)
	fmt.Printf("対象数字: %v\n\n", numbers)

	if len(numbers) == 1 {
		// 単一数字の場合 - 詳細表示
		result, err := numberCount.CountSpecificNumber(numbers[0], *range_)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		
		numberCount.PrintResult(result)
	} else {
		// 複数数字の場合 - 一覧表示
		results, err := numberCount.CountMultipleNumbers(numbers, *range_)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		
		numberCount.PrintMultipleResults(results)
	}

	fmt.Println("\n分析完了！")
}
