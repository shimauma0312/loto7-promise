package heatmap

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shimauma0312/loto7-promise/internal/result"
)

// PositionHeatmapData ヒートマップ用の位置別出現データ
type PositionHeatmapData struct {
	Number        int     `json:"number"`        // 数字（1-37）
	Position      int     `json:"position"`      // 位置（1-7）
	Count         int     `json:"count"`         // その位置での出現回数
	TotalDraws    int     `json:"total_draws"`   // 調査対象の抽選回数
	Percentage    float64 `json:"percentage"`    // その位置での出現率（%）
	OverallCount  int     `json:"overall_count"` // 全位置での出現回数
	PositionRatio float64 `json:"position_ratio"` // 全出現回数に対するその位置での出現率（%）
}

// HeatmapSummary ヒートマップの全体サマリー
type HeatmapSummary struct {
	SearchRange      int                    `json:"search_range"`       // 検索範囲
	ActualDraws      int                    `json:"actual_draws"`       // 実際の抽選回数
	PositionData     [][]PositionHeatmapData `json:"position_data"`     // [number][position]の2次元データ
	HottestNumbers   []int                  `json:"hottest_numbers"`    // 最も出現回数の多い数字
	ColdestNumbers   []int                  `json:"coldest_numbers"`    // 最も出現回数の少ない数字
	MostBiasedNumber int                    `json:"most_biased_number"` // 最も位置による偏りが大きい数字
}

// GenerateHeatmap 指定した検索範囲でヒートマップデータを生成する
// searchRange: 過去何回分の抽選結果を調査するか
func GenerateHeatmap(searchRange int) (*HeatmapSummary, error) {
	if searchRange <= 0 {
		return nil, fmt.Errorf("検索範囲は1以上で指定してください。入力値: %d", searchRange)
	}

	records := result.GetResult(searchRange)
	if records == nil {
		return nil, fmt.Errorf("ロト7の結果データを取得できませんでした")
	}

	actualDraws := len(records)
	
	// 位置別出現回数を集計するマップ [number][position] = count
	positionCounts := make(map[int]map[int]int)
	overallCounts := make(map[int]int) // 各数字の全体出現回数
	
	// 初期化 - 全数字・全位置を0で初期化
	for number := 1; number <= 37; number++ {
		positionCounts[number] = make(map[int]int)
		overallCounts[number] = 0
		for position := 1; position <= 7; position++ {
			positionCounts[number][position] = 0
		}
	}

	// 実際のデータを集計
	for _, record := range records {
		for position := 0; position < 7 && position < len(record); position++ {
			numberStr := strings.TrimSpace(record[position])
			if numberStr == "" {
				continue
			}
			
			number, err := strconv.Atoi(numberStr)
			if err != nil {
				continue // パースエラーはスキップ
			}
			
			if number >= 1 && number <= 37 {
				positionCounts[number][position+1]++ // positionは1から始まる
				overallCounts[number]++
			}
		}
	}

	// PositionHeatmapDataの2次元スライスを生成
	positionData := make([][]PositionHeatmapData, 38) // インデックス0は使わない
	for number := 1; number <= 37; number++ {
		positionData[number] = make([]PositionHeatmapData, 8) // インデックス0は使わない
		
		for position := 1; position <= 7; position++ {
			count := positionCounts[number][position]
			percentage := 0.0
			positionRatio := 0.0
			
			if actualDraws > 0 {
				percentage = (float64(count) / float64(actualDraws)) * 100
			}
			
			if overallCounts[number] > 0 {
				positionRatio = (float64(count) / float64(overallCounts[number])) * 100
			}
			
			positionData[number][position] = PositionHeatmapData{
				Number:        number,
				Position:      position,
				Count:         count,
				TotalDraws:    actualDraws,
				Percentage:    percentage,
				OverallCount:  overallCounts[number],
				PositionRatio: positionRatio,
			}
		}
	}

	// 統計情報を計算
	hottestNumbers := findHottestNumbers(overallCounts)
	coldestNumbers := findColdestNumbers(overallCounts)
	mostBiasedNumber := findMostBiasedNumber(positionCounts, overallCounts)

	return &HeatmapSummary{
		SearchRange:      searchRange,
		ActualDraws:      actualDraws,
		PositionData:     positionData,
		HottestNumbers:   hottestNumbers,
		ColdestNumbers:   coldestNumbers,
		MostBiasedNumber: mostBiasedNumber,
	}, nil
}

// GetPositionData 指定した数字の位置別データを取得
func (hs *HeatmapSummary) GetPositionData(number int) ([]PositionHeatmapData, error) {
	if number < 1 || number > 37 {
		return nil, fmt.Errorf("数字は1から37の範囲で指定してください。入力値: %d", number)
	}
	
	return hs.PositionData[number][1:], nil // インデックス0をスキップして返す
}

// GetNumbersByPosition 指定した位置での全数字のデータを取得
func (hs *HeatmapSummary) GetNumbersByPosition(position int) ([]PositionHeatmapData, error) {
	if position < 1 || position > 7 {
		return nil, fmt.Errorf("位置は1から7の範囲で指定してください。入力値: %d", position)
	}
	
	result := make([]PositionHeatmapData, 37)
	for number := 1; number <= 37; number++ {
		result[number-1] = hs.PositionData[number][position]
	}
	
	return result, nil
}

// PrintHeatmapTable ヒートマップをテーブル形式で出力
func (hs *HeatmapSummary) PrintHeatmapTable() {
	fmt.Printf("\n=== ロト7位置別出現ヒートマップ（過去%d回分） ===\n", hs.SearchRange)
	fmt.Printf("実際の抽選回数: %d回\n\n", hs.ActualDraws)
	
	// ヘッダー行
	fmt.Printf("数字")
	for position := 1; position <= 7; position++ {
		fmt.Printf("\t位置%d", position)
	}
	fmt.Printf("\t合計\n")
	
	// データ行
	for number := 1; number <= 37; number++ {
		fmt.Printf("%2d", number)
		total := 0
		for position := 1; position <= 7; position++ {
			count := hs.PositionData[number][position].Count
			fmt.Printf("\t%3d", count)
			total += count
		}
		fmt.Printf("\t%3d\n", total)
	}
	
	// 統計情報
	fmt.Printf("\n=== 統計情報 ===\n")
	fmt.Printf("最頻出数字: %v\n", hs.HottestNumbers)
	fmt.Printf("最低出現数字: %v\n", hs.ColdestNumbers)
	fmt.Printf("最も位置偏りが大きい数字: %d\n", hs.MostBiasedNumber)
}

// 最も出現回数の多い数字を見つける
func findHottestNumbers(overallCounts map[int]int) []int {
	maxCount := 0
	for _, count := range overallCounts {
		if count > maxCount {
			maxCount = count
		}
	}
	
	var hottest []int
	for number := 1; number <= 37; number++ {
		if overallCounts[number] == maxCount {
			hottest = append(hottest, number)
		}
	}
	
	return hottest
}

// 最も出現回数の少ない数字を見つける
func findColdestNumbers(overallCounts map[int]int) []int {
	minCount := 999999
	for _, count := range overallCounts {
		if count < minCount {
			minCount = count
		}
	}
	
	var coldest []int
	for number := 1; number <= 37; number++ {
		if overallCounts[number] == minCount {
			coldest = append(coldest, number)
		}
	}
	
	return coldest
}

// 最も位置による偏りが大きい数字を見つける
func findMostBiasedNumber(positionCounts map[int]map[int]int, overallCounts map[int]int) int {
	maxBias := 0.0
	mostBiasedNumber := 1
	
	for number := 1; number <= 37; number++ {
		if overallCounts[number] == 0 {
			continue
		}
		
		// 各位置での出現率の分散を計算
		mean := float64(overallCounts[number]) / 7.0
		variance := 0.0
		
		for position := 1; position <= 7; position++ {
			diff := float64(positionCounts[number][position]) - mean
			variance += diff * diff
		}
		
		if variance > maxBias {
			maxBias = variance
			mostBiasedNumber = number
		}
	}
	
	return mostBiasedNumber
}
