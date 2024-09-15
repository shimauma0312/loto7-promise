package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/shimauma0312/loto7-results/internal/result"
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
	result.GetResult(num)

}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  <過去回数> 過去回数分の当選番号を出力します")
}
