package result

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// CacheMetadata キャッシュのメタデータ
type CacheMetadata struct {
	LastUpdated   time.Time `json:"last_updated"`
	LatestDrawNum int       `json:"latest_draw_num"`
	TotalRecords  int       `json:"total_records"`
}

// DrawResult 抽選結果
type DrawResult struct {
	DrawNumber int      `json:"draw_number"`
	Numbers    []string `json:"numbers"`
	Date       string   `json:"date"`
}

// キャッシュディレクトリのパス
const (
	CacheDir     = "cache"
	DataFile     = "cache/loto7_results.json"
	MetadataFile = "cache/metadata.json"
)

// GetResult 過去の抽選結果を取得（キャッシュ機能付き）
func GetResult(repeatNum int) [][]string {
	// キャッシュから結果を取得
	results, err := getCachedResults(repeatNum)
	if err != nil {
		fmt.Printf("キャッシュからの取得に失敗しました。従来の方法で取得します: %v\n", err)
		return GetResultDirect(repeatNum)
	}

	return results
} // GetResultDirect 従来の直接ダウンロード方式（フォールバック用）
func GetResultDirect(repeatNum int) [][]string {
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

// initCache キャッシュディレクトリを初期化
func initCache() error {
	return os.MkdirAll(CacheDir, 0755)
}

// loadMetadata メタデータを読み込み
func loadMetadata() (*CacheMetadata, error) {
	if _, err := os.Stat(MetadataFile); os.IsNotExist(err) {
		return &CacheMetadata{}, nil
	}

	data, err := ioutil.ReadFile(MetadataFile)
	if err != nil {
		return nil, err
	}

	var metadata CacheMetadata
	err = json.Unmarshal(data, &metadata)
	return &metadata, err
}

// saveMetadata メタデータを保存
func saveMetadata(metadata *CacheMetadata) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(MetadataFile, data, 0644)
}

// loadCachedData キャッシュされたデータを読み込み
func loadCachedData() ([]DrawResult, error) {
	if _, err := os.Stat(DataFile); os.IsNotExist(err) {
		return []DrawResult{}, nil
	}

	data, err := ioutil.ReadFile(DataFile)
	if err != nil {
		return nil, err
	}

	var results []DrawResult
	err = json.Unmarshal(data, &results)
	return results, err
}

// saveCachedData キャッシュデータを保存
func saveCachedData(results []DrawResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(DataFile, data, 0644)
}

// needUpdate 更新が必要かチェック
func needUpdate() (bool, error) {
	metadata, err := loadMetadata()
	if err != nil {
		return false, err
	}

	if metadata.LastUpdated.IsZero() {
		return true, nil
	}

	currentDrawNum := NewNumber()
	return currentDrawNum > metadata.LatestDrawNum, nil
}

// updateCache キャッシュを更新
func updateCache() error {
	if err := initCache(); err != nil {
		return err
	}

	needUpd, err := needUpdate()
	if err != nil {
		return err
	}

	if !needUpd {
		fmt.Println("キャッシュは最新です。")
		return nil
	}

	cachedData, err := loadCachedData()
	if err != nil {
		return err
	}

	currentDrawNum := NewNumber()
	var latestData []DrawResult
	var startDrawNum int

	if len(cachedData) == 0 {
		startDrawNum = currentDrawNum - 9 // 初回は10回分（段階的に増やす）
		fmt.Printf("初回実行: 第%d回～第%d回を取得中...\n", startDrawNum, currentDrawNum)
	} else {
		lastCachedDrawNum := cachedData[0].DrawNumber
		startDrawNum = lastCachedDrawNum + 1
		fmt.Printf("更新: 第%d回～第%d回を取得中...\n", startDrawNum, currentDrawNum)
	}

	for drawNum := startDrawNum; drawNum <= currentDrawNum; drawNum++ {
		record, err := GetCsv(drawNum)
		if err != nil {
			fmt.Printf("第%d回の取得に失敗: %v\n", drawNum, err)
			continue
		}

		drawResult := DrawResult{
			DrawNumber: drawNum,
			Numbers:    record,
			Date:       time.Now().Format("2006-01-02"),
		}
		latestData = append(latestData, drawResult)

		if drawNum < currentDrawNum {
			time.Sleep(1 * time.Second)
		}
	}

	updatedData := append(latestData, cachedData...)

	if err := saveCachedData(updatedData); err != nil {
		return err
	}

	metadata := &CacheMetadata{
		LastUpdated:   time.Now(),
		LatestDrawNum: currentDrawNum,
		TotalRecords:  len(updatedData),
	}

	if err := saveMetadata(metadata); err != nil {
		return err
	}

	if len(latestData) > 0 {
		fmt.Printf("キャッシュを更新しました。(新規%d件、全%d件)\n", len(latestData), len(updatedData))
	}
	return nil
}

// getCachedResults キャッシュから指定回数分の結果を取得
func getCachedResults(repeatNum int) ([][]string, error) {
	if err := updateCache(); err != nil {
		return nil, err
	}

	cachedData, err := loadCachedData()
	if err != nil {
		return nil, err
	}

	if len(cachedData) == 0 {
		return nil, fmt.Errorf("キャッシュデータが存在しません")
	}

	// 要求された回数分のデータが不足している場合、追加取得
	if len(cachedData) < repeatNum {
		fmt.Printf("キャッシュ不足: %d回分必要、%d回分のみ存在。追加取得中...\n", repeatNum, len(cachedData))

		currentDrawNum := NewNumber()
		oldestCachedDrawNum := cachedData[len(cachedData)-1].DrawNumber
		startDrawNum := oldestCachedDrawNum - (repeatNum - len(cachedData))

		if startDrawNum < 1 {
			startDrawNum = 1 // 最低回数は1
		}

		var additionalData []DrawResult
		for drawNum := startDrawNum; drawNum < oldestCachedDrawNum; drawNum++ {
			record, err := GetCsv(drawNum)
			if err != nil {
				fmt.Printf("第%d回の取得に失敗: %v\n", drawNum, err)
				continue
			}

			drawResult := DrawResult{
				DrawNumber: drawNum,
				Numbers:    record,
				Date:       time.Now().Format("2006-01-02"),
			}
			additionalData = append(additionalData, drawResult)

			time.Sleep(1 * time.Second)
		}

		// 古いデータを既存キャッシュの末尾に追加
		updatedData := append(cachedData, additionalData...)

		if err := saveCachedData(updatedData); err != nil {
			return nil, err
		}

		// メタデータを更新
		metadata := &CacheMetadata{
			LastUpdated:   time.Now(),
			LatestDrawNum: currentDrawNum,
			TotalRecords:  len(updatedData),
		}

		if err := saveMetadata(metadata); err != nil {
			return nil, err
		}

		cachedData = updatedData
		fmt.Printf("追加取得完了: 全%d件\n", len(cachedData))
	}

	var results [][]string
	for i := 0; i < repeatNum && i < len(cachedData); i++ {
		results = append(results, cachedData[i].Numbers)
	}

	return results, nil
}
