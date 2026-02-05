package recommendation

import (
	"fmt"
	"strings"
)

// 推薦結果の整形を行うインターフェース
type Formatter interface {
	Format(recommendations [][]int, stats StatisticalAnalysis, analyzedCount int) string
}

// テキスト形式で推薦結果を整形する実装
type TextFormatter struct{}

// 推薦結果をテキスト形式で整形する
func (f *TextFormatter) Format(recommendations [][]int, stats StatisticalAnalysis, analyzedCount int) string {
	var result strings.Builder

	result.WriteString("=== ロト7 推薦番号 ===\n\n")

	for i, combination := range recommendations {
		result.WriteString(fmt.Sprintf("推薦 %d: ", i+1))

		// 数字を表示
		for j, num := range combination {
			if j > 0 {
				result.WriteString(" - ")
			}
			result.WriteString(fmt.Sprintf("%02d", num))
		}

		result.WriteString("\n")
	}

	result.WriteString("\n=== 分析情報 ===\n")
	result.WriteString(fmt.Sprintf("分析対象: 過去%d回分\n", analyzedCount))
	result.WriteString(fmt.Sprintf("生成した推薦数: %d\n", len(recommendations)))

	// 統計情報を追加
	result.WriteString("\n=== 統計データ ===\n")
	result.WriteString(fmt.Sprintf("出現頻度が高い数字: %v\n", stats.FrequentNumbers))
	result.WriteString(fmt.Sprintf("波が来ている数字 (直近10回): %v\n", stats.HotNumbers))
	result.WriteString(fmt.Sprintf("復活候補 (20回以上未出現): %v\n", stats.RevivalCandidates))

	return result.String()
}
