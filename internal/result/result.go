package result

import (
	"fmt"
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
	}
	return records
}
