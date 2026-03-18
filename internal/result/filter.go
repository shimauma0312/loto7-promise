package result

import "strconv"

const (
	// HighNumberThreshold は「高い数字」とみなす境界値
	HighNumberThreshold = 29
	// MinHighNumberCount は正常な抽選結果に必要な高い数字の最小個数
	MinHighNumberCount = 2
)

// FilterAbnormalResults は変異体（大阪抽選など）とみなされる結果を除去する。
//
// 東京抽選では29以上の数字が2個以上選ばれる傾向が強い。
// 29以上の数字が2個未満の結果は異常とみなして取り除く。
func FilterAbnormalResults(results [][]string) [][]string {
	filtered := make([][]string, 0, len(results))
	for _, draw := range results {
		if countHighNumbers(draw) >= MinHighNumberCount {
			filtered = append(filtered, draw)
		}
	}
	return filtered
}

// countHighNumbers は1回の抽選結果に含まれる HighNumberThreshold 以上の数字の個数を返す。
func countHighNumbers(draw []string) int {
	count := 0
	for _, numStr := range draw {
		n, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}
		if n >= HighNumberThreshold {
			count++
		}
	}
	return count
}
