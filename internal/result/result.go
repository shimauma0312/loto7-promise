package result

import (
	"fmt"
	"loto7-results/internal"
)

func Result(repeatNum int) {
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
