package main

import (
	"fmt"

	"github.com/shimauma0312/loto7-promise/internal/frequentNumbers"
)

func main() {
	top10 := frequentNumbers.GetFrequentNumbers()

	// 結果を表示
	fmt.Println("出現回数多い順:")
	for i, num := range top10 {
		fmt.Printf("%d列目 : %d\n", i+1, num)
	}
}
