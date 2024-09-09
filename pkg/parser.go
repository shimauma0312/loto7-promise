package pkg

import (
	"encoding/csv"
	"loto7-results/api"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// TODO呼び出し引数
func Parse() ([]string, error) {

	url := "https://www.mizuhobank.co.jp/retail/takarakuji/loto/loto7/csv/A1030" + strconv.Itoa(calculateFridayNumber()) + ".CSV"
	body, err := api.FetchCSV(url)
	if err != nil {
		return nil, err
	}

	reader := transform.NewReader(strings.NewReader(string(body)), japanese.ShiftJIS.NewDecoder())

	r := csv.NewReader(reader)
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return records[3], nil
}
