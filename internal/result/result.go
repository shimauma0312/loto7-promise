package result

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
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

// DrawResultWithDate 抽選結果
type DrawResultWithDate struct {
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

// GetResult 過去の抽選結果を取得
func GetResult(repeatNum int) [][]string {
	// キャッシュから結果を取得
	results, err := getCachedResults(repeatNum)
	if err != nil {
		// キャッシュからの取得に失敗した場合は従来の方法で取得
		return GetResultDirect(repeatNum)
	}

	return results
}

// GetResultWithDate 過去の抽選結果を日付情報付きで取得
func GetResultWithDate(repeatNum int) ([]DrawResultWithDate, error) {
	// キャッシュから結果を取得
	cachedData, err := getCachedResultsWithDate(repeatNum)
	if err != nil {
		// キャッシュからの取得に失敗した場合は従来の方法で取得
		return getResultDirectWithDate(repeatNum), nil
	}

	return cachedData, nil
}

// GetResultDirect 従来の直接ダウンロード
func GetResultDirect(repeatNum int) [][]string {
	newNum := NewNumber()

	records := make([][]string, 0)
	for i := 0; i < repeatNum; i++ {
		cnt := newNum - i
		record, err := GetCsv(cnt)
		if err != nil {
			// エラーが発生した場合は処理を中断
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

	data, err := os.ReadFile(MetadataFile)
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
	return os.WriteFile(MetadataFile, data, 0644)
}

// loadCachedData キャッシュされたデータを読み込み
func loadCachedData() ([]DrawResult, error) {
	if _, err := os.Stat(DataFile); os.IsNotExist(err) {
		return []DrawResult{}, nil
	}

	data, err := os.ReadFile(DataFile)
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
	return os.WriteFile(DataFile, data, 0644)
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
		// キャッシュは最新
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
	} else {
		lastCachedDrawNum := cachedData[0].DrawNumber
		startDrawNum = lastCachedDrawNum + 1
	}

	for drawNum := startDrawNum; drawNum <= currentDrawNum; drawNum++ {
		record, err := GetCsv(drawNum)
		if err != nil {
			// 取得に失敗した場合はスキップして次へ
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

	// latestDataを降順にソート（新しい→古い）
	sort.Slice(latestData, func(i, j int) bool {
		return latestData[i].DrawNumber > latestData[j].DrawNumber
	})
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
		// キャッシュ更新完了
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
		// キャッシュ不足のため追加取得を実行

		currentDrawNum := NewNumber()
		oldestCachedDrawNum := cachedData[len(cachedData)-1].DrawNumber
		startDrawNum := oldestCachedDrawNum - (repeatNum - len(cachedData))

		if startDrawNum < 1 {
			startDrawNum = 1 // 最低回数は1
		}

		var additionalData []DrawResult
		for drawNum := oldestCachedDrawNum - 1; drawNum >= startDrawNum; drawNum-- {
			record, err := GetCsv(drawNum)
			if err != nil {
				// 取得に失敗した場合はスキップして次へ
				continue
			}

			drawResult := DrawResult{
				DrawNumber: drawNum,
				Numbers:    record,
				Date:       GetDrawDate(drawNum), // 実際の抽選日を計算
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
		// 追加取得完了
	}

	var results [][]string
	for i := 0; i < repeatNum && i < len(cachedData); i++ {
		results = append(results, cachedData[i].Numbers)
	}

	return results, nil
}

// getCachedResultsWithDate キャッシュから指定回数分の結果を日付情報付きで取得
func getCachedResultsWithDate(repeatNum int) ([]DrawResultWithDate, error) {
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
		currentDrawNum := NewNumber()
		oldestCachedDrawNum := cachedData[len(cachedData)-1].DrawNumber
		startDrawNum := oldestCachedDrawNum - (repeatNum - len(cachedData))

		if startDrawNum < 1 {
			startDrawNum = 1
		}

		var additionalData []DrawResult
		for drawNum := oldestCachedDrawNum - 1; drawNum >= startDrawNum; drawNum-- {
			record, err := GetCsv(drawNum)
			if err != nil {
				continue
			}

			drawResult := DrawResult{
				DrawNumber: drawNum,
				Numbers:    record,
				Date:       GetDrawDate(drawNum), // 実際の抽選日を計算
			}
			additionalData = append(additionalData, drawResult)

			time.Sleep(1 * time.Second)
		}

		// データを更新
		updatedData := append(cachedData, additionalData...)
		if err := saveCachedData(updatedData); err != nil {
			return nil, err
		}

		metadata := &CacheMetadata{
			LastUpdated:   time.Now(),
			LatestDrawNum: currentDrawNum,
			TotalRecords:  len(updatedData),
		}

		if err := saveMetadata(metadata); err != nil {
			return nil, err
		}

		cachedData = updatedData
	}

	var results []DrawResultWithDate
	for i := 0; i < repeatNum && i < len(cachedData); i++ {
		results = append(results, DrawResultWithDate{
			DrawNumber: cachedData[i].DrawNumber,
			Numbers:    cachedData[i].Numbers,
			Date:       cachedData[i].Date,
		})
	}

	return results, nil
}

// getResultDirectWithDate 従来の直接ダウンロード（日付情報付き）
func getResultDirectWithDate(repeatNum int) []DrawResultWithDate {
	newNum := NewNumber()

	var results []DrawResultWithDate
	for i := 0; i < repeatNum; i++ {
		drawNum := newNum - i
		if drawNum < 1 {
			break
		}

		record, err := GetCsv(drawNum)
		if err != nil {
			continue
		}

		result := DrawResultWithDate{
			DrawNumber: drawNum,
			Numbers:    record,
			Date:       GetDrawDate(drawNum), // 実際の抽選日を計算
		}
		results = append(results, result)
	}

	return results
}
