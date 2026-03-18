// 過去の抽選頻度に基づいた重み付き非復元抽出でロト7の組み合わせを生成する。
package random

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shimauma0312/loto7-promise/internal/result"
)

const (
	// TotalNumbers はロト7の総数字数
	TotalNumbers = 37
	// DrawCount は1回の抽選で引く数字数
	DrawCount = 7
)

// 数字ごとの出現統計
type NumberStat struct {
	Number    int
	Frequency int // 過去データでの出現回数
}

// 過去の出現頻度を重みとした重み付き非復元抽出で組み合わせを生成するエンジン
type RandomEngine struct {
	rng       *rand.Rand
	results   [][]string
	frequency map[int]int // 各数字の出現回数（1-37）
}

// 現在時刻をシードに初期化したエンジンを返す
func NewRandomEngine() *RandomEngine {
	source := rand.NewSource(time.Now().UnixNano())
	return &RandomEngine{
		rng:       rand.New(source),
		frequency: make(map[int]int),
	}
}

// 過去count回の抽選結果を読み込み、変異体を除外したうえで頻度を集計する
func (re *RandomEngine) LoadData(count int) error {
	re.results = result.FilterAbnormalResults(result.GetResult(count))
	if len(re.results) == 0 {
		return fmt.Errorf("データを取得できませんでした")
	}
	re.analyzeFrequency()
	return nil
}

func (re *RandomEngine) analyzeFrequency() {
	re.frequency = make(map[int]int)
	for _, draw := range re.results {
		for _, numStr := range draw {
			n, err := strconv.Atoi(numStr)
			if err == nil && n >= 1 && n <= TotalNumbers {
				re.frequency[n]++
			}
		}
	}
}

// 頻度を重み（√頻度）とした非復元ルーレット選択で7個の数字を生成する
func (re *RandomEngine) GenerateRandomCombination() ([]int, error) {
	if len(re.frequency) == 0 {
		return nil, fmt.Errorf("データが読み込まれていません")
	}

	pool := make([]int, TotalNumbers)
	weights := make([]float64, TotalNumbers)
	for i := 0; i < TotalNumbers; i++ {
		num := i + 1
		pool[i] = num
		freq := float64(re.frequency[num])
		if freq == 0 {
			freq = 1.0
		}
		weights[i] = math.Sqrt(freq)
	}

	combination := make([]int, 0, DrawCount)
	size := TotalNumbers

	for len(combination) < DrawCount {
		total := 0.0
		for i := 0; i < size; i++ {
			total += weights[i]
		}

		r := re.rng.Float64() * total
		cumulative := 0.0
		chosen := size - 1
		for i := 0; i < size; i++ {
			cumulative += weights[i]
			if r < cumulative {
				chosen = i
				break
			}
		}

		combination = append(combination, pool[chosen])

		// 末尾と交換して除去
		size--
		pool[chosen] = pool[size]
		weights[chosen] = weights[size]
	}

	sort.Ints(combination)
	return combination, nil
}

// 出現回数の降順で全37数字の統計を返す
func (re *RandomEngine) GetNumberStats() []NumberStat {
	stats := make([]NumberStat, TotalNumbers)
	for i := 0; i < TotalNumbers; i++ {
		stats[i] = NumberStat{
			Number:    i + 1,
			Frequency: re.frequency[i+1],
		}
	}
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Frequency == stats[j].Frequency {
			return stats[i].Number < stats[j].Number
		}
		return stats[i].Frequency > stats[j].Frequency
	})
	return stats
}

// 重複のない組み合わせをcount個生成する
func (re *RandomEngine) GenerateMultipleRandomCombinations(count int) ([][]int, error) {
	if count <= 0 {
		return nil, fmt.Errorf("生成数は1以上である必要があります")
	}

	combinations := make([][]int, 0, count)
	usedCombinations := make(map[string]bool)

	maxAttempts := count * 100

	for attempt := 0; attempt < maxAttempts && len(combinations) < count; attempt++ {
		combination, err := re.GenerateRandomCombination()
		if err != nil {
			continue
		}

		key := combinationKey(combination)
		if !usedCombinations[key] {
			combinations = append(combinations, combination)
			usedCombinations[key] = true
		}
	}

	if len(combinations) == 0 {
		return nil, fmt.Errorf("組み合わせを生成できませんでした")
	}

	return combinations, nil
}

func combinationKey(combination []int) string {
	var builder strings.Builder
	for i, num := range combination {
		if i > 0 {
			builder.WriteString("-")
		}
		builder.WriteString(strconv.Itoa(num))
	}
	return builder.String()
}
