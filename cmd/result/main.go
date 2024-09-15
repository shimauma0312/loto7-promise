package main

import (
	"flag"
	"fmt"
	"loto7-results/pkg"
	"os"
	"strconv"
)

func main() {
	resultFlag := flag.Bool("result", false, "repeat_numを指定してresultメソッドを実行")

	flag.Parse()

	if len(os.Args) == 1 {
		printUsage()
		return
	}

	if *resultFlag {
		if flag.NArg() < 1 {
			fmt.Println("Error: resultには過去回数を指定してください")
			printUsage()
			return
		}
		repeatNum := flag.Arg(0)
		num, err := strconv.Atoi(repeatNum)

		if err != nil {
			fmt.Println("Error: 過去回数は数値で指定してください")
			return
		}
		result(num)
	} else {
		fmt.Println("Error: 不明な引数")
		printUsage()
	}
}

func result(repeatNum int) {
	newNum := pkg.NewNumber()
	for i := 0; i < repeatNum; i++ {
		cnt := newNum - i
		records, err := pkg.GetResult(cnt)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("第%d回：%s\n", cnt, records)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  -result <過去回数>  過去回数分の当選番号を出力します")
}
