package recommendation

import "github.com/shimauma0312/loto7-promise/internal/result"

// FilterAbnormalResults は変異体（大阪抽選など）とみなされる結果を除去する。
// 実装は result.FilterAbnormalResults に委譲する。
func FilterAbnormalResults(results [][]string) [][]string {
	return result.FilterAbnormalResults(results)
}
