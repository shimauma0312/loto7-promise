package randomPrediction

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/shimauma0312/loto7-promise/internal/numberCount"
)

// PredictionResult 予想番号の結果を格納する構造体
type PredictionResult struct {
	MainNumbers   []int `json:"main_numbers"`   // 本数字7個
	BonusNumbers  []int `json:"bonus_numbers"`  // ボーナス数字2個
	TrendBased    bool  `json:"trend_based"`    // 傾向ベースの予想かどうか
	GeneratedTime time.Time `json:"generated_time"` // 生成時刻
}

// WeightedNumber 重み付き数字を表す構造体
type WeightedNumber struct {
	Number int     // 数字
	Weight float64 // 重み（出現頻度等から算出）
}

// GenerateRandomPrediction 過去データの傾向を加味したランダム予想番号を生成
// 引数: analysisRange - 傾向分析に使用する過去データの範囲, useTrend - 傾向を考慮するかどうか
// 戻り値: PredictionResult構造体のポインタ、エラー
func GenerateRandomPrediction(analysisRange int, useTrend bool) (*PredictionResult, error) {
	// 入力値の妥当性チェック
	if analysisRange <= 0 {
		return nil, fmt.Errorf("分析範囲は1以上で指定してください。入力値: %d", analysisRange)
	}

	rand.Seed(time.Now().UnixNano())

	var mainNumbers []int
	var bonusNumbers []int

	if useTrend {
		// 傾向を考慮した予想番号生成
		weightedNumbers, err := calculateNumberWeights(analysisRange)
		if err != nil {
			return nil, fmt.Errorf("傾向分析に失敗しました: %v", err)
		}

		mainNumbers = selectWeightedNumbers(weightedNumbers, 7)
		bonusNumbers = selectWeightedNumbers(weightedNumbers, 2)
	} else {
		// 完全ランダム生成
		mainNumbers = generateRandomNumbers(7)
		bonusNumbers = generateRandomNumbers(2)
	}

	// 数字をソート
	sort.Ints(mainNumbers)
	sort.Ints(bonusNumbers)

	result := &PredictionResult{
		MainNumbers:   mainNumbers,
		BonusNumbers:  bonusNumbers,
		TrendBased:    useTrend,
		GeneratedTime: time.Now(),
	}

	return result, nil
}

// calculateNumberWeights 過去データから各数字の重みを算出
// 引数: analysisRange - 分析対象の回数
// 戻り値: 重み付き数字のスライス、エラー
func calculateNumberWeights(analysisRange int) ([]WeightedNumber, error) {
	var weightedNumbers []WeightedNumber

	// 各数字（1-37）の出現頻度を取得
	for i := 1; i <= 37; i++ {
		countResult, err := numberCount.CountSpecificNumber(i, analysisRange)
		if err != nil {
			return nil, fmt.Errorf("数字%dの出現回数取得に失敗: %v", i, err)
		}

		// 基本重みは出現率をベースに設定
		weight := countResult.Percentage / 100.0
		
		// 最低重みを設定（完全に0にならないよう調整）
		if weight < 0.01 {
			weight = 0.01
		}
		
		// 最近の傾向を重視（直近データの重みを増加）
		recentWeight := calculateRecentTrend(i, analysisRange/3)
		weight = (weight * 0.7) + (recentWeight * 0.3)

		weightedNumbers = append(weightedNumbers, WeightedNumber{
			Number: i,
			Weight: weight,
		})
	}

	return weightedNumbers, nil
}

// calculateRecentTrend 最近の出現傾向を計算
// 引数: number - 対象数字, recentRange - 最近の範囲
// 戻り値: 最近の傾向重み
func calculateRecentTrend(number int, recentRange int) float64 {
	if recentRange <= 0 {
		recentRange = 10 // デフォルト値
	}

	countResult, err := numberCount.CountSpecificNumber(number, recentRange)
	if err != nil {
		return 0.01 // エラー時は最低重み
	}

	// 最近の出現率を重みとして返す
	recentWeight := countResult.Percentage / 100.0
	if recentWeight < 0.01 {
		recentWeight = 0.01
	}

	return recentWeight
}

// selectWeightedNumbers 重み付きから指定個数の数字を選択
// 引数: weightedNumbers - 重み付き数字のスライス, count - 選択する個数
// 戻り値: 選択された数字のスライス
func selectWeightedNumbers(weightedNumbers []WeightedNumber, count int) []int {
	var selected []int
	available := make([]WeightedNumber, len(weightedNumbers))
	copy(available, weightedNumbers)

	for len(selected) < count && len(available) > 0 {
		// 重みに基づいて数字を選択
		selectedIndex := weightedRandomSelect(available)
		selected = append(selected, available[selectedIndex].Number)
		
		// 選択された数字を利用可能リストから除去
		available = append(available[:selectedIndex], available[selectedIndex+1:]...)
	}

	return selected
}

// weightedRandomSelect 重み付きランダム選択
// 引数: numbers - 重み付き数字のスライス
// 戻り値: 選択されたインデックス
func weightedRandomSelect(numbers []WeightedNumber) int {
	// 重みの合計を計算
	totalWeight := 0.0
	for _, num := range numbers {
		totalWeight += num.Weight
	}

	// ランダム値を生成
	randomValue := rand.Float64() * totalWeight

	// 重みに基づいてインデックスを選択
	currentWeight := 0.0
	for i, num := range numbers {
		currentWeight += num.Weight
		if randomValue <= currentWeight {
			return i
		}
	}

	// フォールバック（通常は到達しない）
	return rand.Intn(len(numbers))
}

// generateRandomNumbers 完全ランダムな数字を生成
// 引数: count - 生成する数字の個数
// 戻り値: 生成された数字のスライス
func generateRandomNumbers(count int) []int {
	var numbers []int
	used := make(map[int]bool)

	for len(numbers) < count {
		number := rand.Intn(37) + 1 // 1-37の範囲
		if !used[number] {
			numbers = append(numbers, number)
			used[number] = true
		}
	}

	return numbers
}

// GenerateMultiplePredictions 複数の予想番号を生成
// 引数: count - 生成する予想の数, analysisRange - 分析範囲, useTrend - 傾向を使用するか
// 戻り値: 予想結果のスライス、エラー
func GenerateMultiplePredictions(count int, analysisRange int, useTrend bool) ([]*PredictionResult, error) {
	if count <= 0 {
		return nil, fmt.Errorf("生成数は1以上で指定してください。入力値: %d", count)
	}

	var predictions []*PredictionResult

	for i := 0; i < count; i++ {
		prediction, err := GenerateRandomPrediction(analysisRange, useTrend)
		if err != nil {
			return nil, fmt.Errorf("予想番号生成に失敗しました（%d回目）: %v", i+1, err)
		}
		predictions = append(predictions, prediction)
	}

	return predictions, nil
}

// PrintPredictionResult 予想結果を表示
// 引数: result - 表示したい予想結果
func PrintPredictionResult(result *PredictionResult) {
	fmt.Println("ロト7予想番号")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("生成時刻: %s\n", result.GeneratedTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("予想タイプ: ")
	if result.TrendBased {
		fmt.Println("傾向分析ベース")
	} else {
		fmt.Println("完全ランダム")
	}
	fmt.Println(strings.Repeat("-", 50))
	
	fmt.Printf("本数字: ")
	for i, num := range result.MainNumbers {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%02d", num)
	}
	fmt.Println()
	
	fmt.Printf("ボーナス数字: ")
	for i, num := range result.BonusNumbers {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%02d", num)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("=", 50))
}

// PrintMultiplePredictions 複数の予想結果を表示
// 引数: results - 表示したい予想結果のスライス
func PrintMultiplePredictions(results []*PredictionResult) {
	fmt.Println("ロト7複数予想番号")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("総予想数: %d組\n", len(results))
	if len(results) > 0 {
		fmt.Printf("予想タイプ: ")
		if results[0].TrendBased {
			fmt.Println("傾向分析ベース")
		} else {
			fmt.Println("完全ランダム")
		}
	}
	fmt.Println(strings.Repeat("-", 60))
	
	for i, result := range results {
		fmt.Printf("【予想%d】", i+1)
		fmt.Printf(" 本数字: ")
		for j, num := range result.MainNumbers {
			if j > 0 {
				fmt.Printf(",")
			}
			fmt.Printf("%02d", num)
		}
		fmt.Printf(" ボーナス: ")
		for j, num := range result.BonusNumbers {
			if j > 0 {
				fmt.Printf(",")
			}
			fmt.Printf("%02d", num)
		}
		fmt.Println()
	}
	fmt.Println(strings.Repeat("=", 60))
}
