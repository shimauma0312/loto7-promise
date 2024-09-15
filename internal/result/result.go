package result

import (
	"fmt"

	"github.com/shimauma0312/loto7-results/internal"
)

func GetResult(repeatNum int) {
	newNum := internal.NewNumber()
	for i := 0; i < repeatNum; i++ {
		cnt := newNum - i
		records, err := internal.GetCsv(cnt)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("第%d回：%s\n", cnt, records)
	}
}
