package internal

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func GetResult(num int) ([]string, error) {
	url := "https://www.mizuhobank.co.jp/retail/takarakuji/loto/loto7/csv/A1030" + strconv.Itoa(num) + ".CSV"
	return Parse(url)
}

func Parse(url string) ([]string, error) {
	body, err := fetchCSV(url)
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

	return trimRecordKeys(records[3]), nil
}

func fetchCSV(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error downloading the CSV: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the response body: %v", err)
	}

	return body, nil
}

func trimRecordKeys(records []string) []string {
	filteredRecords := records[1:8]

	return filteredRecords
}
