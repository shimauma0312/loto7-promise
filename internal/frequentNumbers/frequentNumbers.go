package frequentNumbers

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/shimauma0312/loto7-promise/internal/result"
)

var results = result.GetResult(100)

func GetFrequentNumbers() [][]int {

	convertToIntSlice := func(numbers []string) []int {
		intNumbers := make([]int, len(numbers))
		for i, num := range numbers {
			intNum, err := strconv.Atoi(num)
			if err != nil {
				fmt.Println("error converting string to int:", err)
				return nil
			}
			intNumbers[i] = intNum
		}
		return intNumbers
	}

	sorts := make([][]int, 0)
	for i := 0; i < 7; i++ {
		numbers := getNumbers(i)
		sorts = append(sorts, countTop10Numbers(convertToIntSlice(numbers)))
	}

	return sorts
}

func getNumbers(key int) []string {
	numbers := make([]string, 0)
	for _, record := range results {
		numbers = append(numbers, record[key])
	}
	return numbers
}

func countTop10Numbers(slice []int) []int {
	// 出現回数をマップで記録
	countMap := make(map[int]int)
	for _, num := range slice {
		countMap[num]++
	}

	// マップをスライスに変換 (数値とその出現回数をペアで格納)
	type numberCount struct {
		Number int
		Count  int
	}
	var counts []numberCount
	for num, count := range countMap {
		counts = append(counts, numberCount{Number: num, Count: count})
	}

	// 出現回数降順にソート
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	//トップ10の数値をスライスに格納
	top10 := []int{}
	for i := 0; i < 10 && i < len(counts); i++ {
		top10 = append(top10, counts[i].Number)
	}

	return top10
}
