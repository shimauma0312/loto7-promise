// 過去の抽選データを統計分析するモジュール
package prediction

import (
	"sort"
	"strconv"
)

// Analyzer は過去の抽選結果から統計分析を行う
type Analyzer struct{}

// Analyze は results（新しい順）を分析し AnalysisData を返す
// recentWindow は直近のホット/コールド判定に使うドロー数
func (a *Analyzer) Analyze(results [][]string, recentWindow int) AnalysisData {
	data := AnalysisData{
		FrequencyMap: make(map[int]int),
		IntervalMap:  make(map[int]int),
		AnalyzedDraws: len(results),
	}

	if len(results) == 0 {
		return data
	}

	// 全ての数字の間隔を最大値で初期化
	for num := 1; num <= 37; num++ {
		data.IntervalMap[num] = len(results)
	}

	// 全期間の出現頻度と最初の出現位置（間隔計算）を収集
	for drawIdx, draw := range results {
		numbers := parseNumberStrings(draw)
		for _, num := range numbers {
			data.FrequencyMap[num]++
			// 新しい順でイテレートしているため、最初に見つかった位置が最短間隔
			if data.IntervalMap[num] == len(results) {
				data.IntervalMap[num] = drawIdx
			}
		}
	}

	// 直近ウィンドウ内の出現頻度を計算
	window := recentWindow
	if window > len(results) {
		window = len(results)
	}
	recentFreq := make(map[int]int)
	for i := 0; i < window; i++ {
		for _, num := range parseNumberStrings(results[i]) {
			recentFreq[num]++
		}
	}

	// ホット・コールド分類
	for num := 1; num <= 37; num++ {
		if recentFreq[num] >= 2 {
			data.HotNumbers = append(data.HotNumbers, num)
		} else if recentFreq[num] == 0 {
			data.ColdNumbers = append(data.ColdNumbers, num)
		}
	}
	sort.Ints(data.HotNumbers)
	sort.Ints(data.ColdNumbers)

	// 前回のドロー情報を収集
	lastDraw := results[0]
	data.LastDrawNumbers = parseNumberStrings(lastDraw)
	data.LastDrawSum = sumInts(data.LastDrawNumbers)

	// 前々回の合計
	if len(results) >= 2 {
		data.PrevDrawSum = sumInts(parseNumberStrings(results[1]))
	}

	// 合計の連続トレンドを判定
	if len(results) >= 3 {
		olderSum := sumInts(parseNumberStrings(results[2]))
		if data.LastDrawSum > data.PrevDrawSum && data.PrevDrawSum > olderSum {
			// 2回連続上昇 → 下降ほうを優先
			data.SumTrend = -1
		} else if data.LastDrawSum < data.PrevDrawSum && data.PrevDrawSum < olderSum {
			// 2回連続下降 → 上昇ほうを優先
			data.SumTrend = 1
		}
	}

	return data
}

// parseNumberStrings は文字列スライスを 1〜37 の整数スライスに変換する
func parseNumberStrings(strs []string) []int {
	nums := make([]int, 0, len(strs))
	for _, s := range strs {
		if n, err := strconv.Atoi(s); err == nil && n >= 1 && n <= 37 {
			nums = append(nums, n)
		}
	}
	sort.Ints(nums)
	return nums
}

// sumInts は整数スライスの合計を返す
func sumInts(nums []int) int {
	s := 0
	for _, n := range nums {
		s += n
	}
	return s
}
