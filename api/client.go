// api/client.go
package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func FetchCSV(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error downloading the CSV: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading the response body: %v", err)
	}

	return body, nil
}
