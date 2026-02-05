package recommendation

import (
	"fmt"
	"sort"
	"strconv"
)

// 統計分析を行うインターフェース
type Analyzer interface {
	Analyze(results [][]string, lookback int) StatisticalAnalysis
	BuildCombinationHistory(results [][]string) CombinationHistory
}

// デフォルトの統計分析実装
type DefaultAnalyzer struct{}

// デフォルトアナライザーを作成する
func NewDefaultAnalyzer() *DefaultAnalyzer {
	return &DefaultAnalyzer{}
}

// 過去データの統計分析を実行する
func (a *DefaultAnalyzer) Analyze(results [][]string, lookback int) StatisticalAnalysis {
	stats := StatisticalAnalysis{
		FrequencyMap:   make(map[int]int),
		LastAppearance: make(map[int]int),
	}

	// 過去lookback回分の出現回数をカウント
	if lookback > len(results) {
		lookback = len(results)
	}

	// 各数字の出現回数と最終出現位置を記録
	for i := 0; i < lookback; i++ {
		for _, numStr := range results[i] {
			num, _ := strconv.Atoi(numStr)
			stats.FrequencyMap[num]++

			// 最終出現位置を更新（最新の結果が0）
			if _, exists := stats.LastAppearance[num]; !exists {
				stats.LastAppearance[num] = i
			}
		}
	}

	// 出現回数でソート
	type numFreq struct {
		num  int
		freq int
	}
	var frequencies []numFreq
	for num := 1; num <= 37; num++ {
		frequencies = append(frequencies, numFreq{num: num, freq: stats.FrequencyMap[num]})
	}
	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].freq > frequencies[j].freq
	})

	// 出現回数が多い数字（上位10個）
	for i := 0; i < 10 && i < len(frequencies); i++ {
		stats.FrequentNumbers = append(stats.FrequentNumbers, frequencies[i].num)
	}

	// 出現回数が少ない数字（下位10個）
	for i := len(frequencies) - 1; i >= len(frequencies)-10 && i >= 0; i-- {
		stats.RareNumbers = append(stats.RareNumbers, frequencies[i].num)
	}

	// 直近10回以内で連続して出ている数字（波が来ている）
	hotCheckRange := 10
	if hotCheckRange > len(results) {
		hotCheckRange = len(results)
	}
	hotCount := make(map[int]int)
	for i := 0; i < hotCheckRange; i++ {
		for _, numStr := range results[i] {
			num, _ := strconv.Atoi(numStr)
			hotCount[num]++
		}
	}
	// 直近10回で3回以上出現している数字を「波が来ている」と判断
	for num, count := range hotCount {
		if count >= 3 {
			stats.HotNumbers = append(stats.HotNumbers, num)
		}
	}

	// 20回以上出ていない復活候補
	for num := 1; num <= 37; num++ {
		lastAppear, exists := stats.LastAppearance[num]
		if !exists || lastAppear >= 20 {
			// ただし、出現回数が少ない数字のリストの上位5個に含まれている場合は無視
			isRare := false
			for i := 0; i < 5 && i < len(stats.RareNumbers); i++ {
				if stats.RareNumbers[i] == num {
					isRare = true
					break
				}
			}
			if !isRare {
				stats.RevivalCandidates = append(stats.RevivalCandidates, num)
			}
		}
	}

	return stats
}

// 組み合わせ履歴を構築する
func (a *DefaultAnalyzer) BuildCombinationHistory(results [][]string) CombinationHistory {
	history := CombinationHistory{
		Pairs: make(map[string]int),
	}

	for _, numbers := range results {
		// 7個の数字から2個ずつの組み合わせを生成
		for i := 0; i < len(numbers); i++ {
			for j := i + 1; j < len(numbers); j++ {
				num1, _ := strconv.Atoi(numbers[i])
				num2, _ := strconv.Atoi(numbers[j])

				// ペアキー作成（小さい数字,大きい数字の順）
				var pairKey string
				if num1 < num2 {
					pairKey = fmt.Sprintf("%d,%d", num1, num2)
				} else {
					pairKey = fmt.Sprintf("%d,%d", num2, num1)
				}

				history.Pairs[pairKey]++
			}
		}
	}

	return history
}
