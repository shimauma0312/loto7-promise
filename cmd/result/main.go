package main

import (
	"flag"
	"fmt"
	"loto7-results/internal"
	"os"
	"strconv"
)

func main() {

	flag.Parse()

	if len(os.Args) == 1 {
		printUsage()
		return
	}

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

}

func result(repeatNum int) {
	newNum := internal.NewNumber()
	for i := 0; i < repeatNum; i++ {
		cnt := newNum - i
		records, err := internal.GetResult(cnt)
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
