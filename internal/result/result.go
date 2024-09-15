package result

import (
	"fmt"
)

func GetResult(repeatNum int) {
	newNum := NewNumber()
	for i := 0; i < repeatNum; i++ {
		cnt := newNum - i
		records, err := GetCsv(cnt)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("第%d回：%s\n", cnt, records)
	}
}
