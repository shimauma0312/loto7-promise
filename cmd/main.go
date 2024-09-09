package main

// pkg/Parser.goを読み取り
import (
	"fmt"
	"loto7-results/pkg"
)

func main() {
	// TODOParseに呼び出し回数指定
	records, err := pkg.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(records)
}
