package result

import (
	"fmt"
	"time"
)

func GetResult(repeatNum int) [][]string {
	newNum := NewNumber()

	records := make([][]string, 0)
	for i := 0; i < repeatNum; i++ {
		cnt := newNum - i
		record, err := GetCsv(cnt)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		records = append(records, record)

		if i < repeatNum-1 {
			time.Sleep(1 * time.Second)
		}
	}
	return records
}
