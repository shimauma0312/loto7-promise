// 過去の出現範囲を分析し、その範囲内でランダムに数字を選択する
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

// 各位置での数字出現範囲
type PositionRange struct {
	Position  int        // 位置（1-7）
	Min       int        // 最小値
	Max       int        // 最大値
	Numbers   []int      // 出現した数字のリスト
	Frequency map[int]int // 各数字の出現回数（重み付き乱択用）
}

// ランダム推薦エンジン
type RandomEngine struct {
	rng     *rand.Rand
	results [][]string
	ranges  []PositionRange
}

// ランダムエンジンを作成する
func NewRandomEngine() *RandomEngine {
	source := rand.NewSource(time.Now().UnixNano())
	return &RandomEngine{
		rng:    rand.New(source),
		ranges: make([]PositionRange, 7),
	}
}

// データを読み込み、各位置の出現範囲を分析する
func (re *RandomEngine) LoadData(count int) error {
	re.results = result.GetResult(count)
	if len(re.results) == 0 {
		return fmt.Errorf("データを取得できませんでした")
	}

	// 各位置の出現範囲を分析
	re.analyzePositionRanges()

	return nil
}

// 各位置での数字の出現範囲を分析する
func (re *RandomEngine) analyzePositionRanges() {
	// 各位置ごとに出現した数字と出現回数を収集
	positionNumbers := make([]map[int]bool, 7)
	positionFrequency := make([]map[int]int, 7)
	for i := 0; i < 7; i++ {
		positionNumbers[i] = make(map[int]bool)
		positionFrequency[i] = make(map[int]int)
	}

	// 過去データから各位置の数字を収集
	for _, drawResult := range re.results {
		if len(drawResult) < 7 {
			continue
		}

		// 数字を整数に変換してソート
		numbers := make([]int, 0, 7)
		for _, numStr := range drawResult[:7] { // ボーナス番号を除く
			num, err := strconv.Atoi(numStr)
			if err == nil && num >= 1 && num <= 37 {
				numbers = append(numbers, num)
			}
		}

		// ソートして位置を確定
		sort.Ints(numbers)

		// 各位置に出現した数字を記録し、出現回数をカウント
		for pos, num := range numbers {
			if pos < 7 {
				positionNumbers[pos][num] = true
				positionFrequency[pos][num]++
			}
		}
	}

	// 各位置の範囲を計算
	for pos := 0; pos < 7; pos++ {
		nums := make([]int, 0, len(positionNumbers[pos]))
		for num := range positionNumbers[pos] {
			nums = append(nums, num)
		}

		if len(nums) == 0 {
			// データがない場合のデフォルト範囲
			re.ranges[pos] = PositionRange{
				Position:  pos + 1,
				Min:       1,
				Max:       37,
				Numbers:   make([]int, 0),
				Frequency: make(map[int]int),
			}
			continue
		}

		sort.Ints(nums)

		re.ranges[pos] = PositionRange{
			Position:  pos + 1,
			Min:       nums[0],
			Max:       nums[len(nums)-1],
			Numbers:   nums,
			Frequency: positionFrequency[pos],
		}
	}
}

// ランダムな組み合わせを生成する
func (re *RandomEngine) GenerateRandomCombination() ([]int, error) {
	if len(re.ranges) != 7 {
		return nil, fmt.Errorf("データが読み込まれていません")
	}

	combination := make([]int, 7)
	usedNumbers := make(map[int]bool)

	// 各位置ごとにランダムに数字を選択
	for pos := 0; pos < 7; pos++ {
		posRange := re.ranges[pos]

		// 出現した数字のリストが空の場合、最小-最大の範囲全体を使用
		var candidates []int
		if len(posRange.Numbers) > 0 {
			candidates = make([]int, len(posRange.Numbers))
			copy(candidates, posRange.Numbers)
		} else {
			// 範囲内の全数字を候補にする
			for num := posRange.Min; num <= posRange.Max; num++ {
				candidates = append(candidates, num)
			}
		}

		// 既に使用済みの数字を除外
		availableCandidates := make([]int, 0)
		for _, num := range candidates {
			if !usedNumbers[num] {
				availableCandidates = append(availableCandidates, num)
			}
		}

		// 候補がない場合は範囲を広げる
		if len(availableCandidates) == 0 {
			// 前の位置の最大値+1から次の位置の最小値-1まで
			minNum := 1
			maxNum := 37

			if pos > 0 {
				minNum = combination[pos-1] + 1
			}
			if pos < 6 {
				maxNum = re.ranges[pos+1].Min - 1
				if maxNum < minNum {
					maxNum = 37
				}
			}

			for num := minNum; num <= maxNum; num++ {
				if !usedNumbers[num] {
					availableCandidates = append(availableCandidates, num)
				}
			}
		}

		if len(availableCandidates) == 0 {
			return nil, fmt.Errorf("位置 %d で選択可能な数字がありません", pos+1)
		}

		// 重み付き乱択で選択
		selectedNum := re.weightedRandomSelect(availableCandidates, posRange.Frequency)
		combination[pos] = selectedNum
		usedNumbers[selectedNum] = true
	}

	// ソートして返す
	sort.Ints(combination)

	return combination, nil
}

// 重み付き乱択で数字を選択する
func (re *RandomEngine) weightedRandomSelect(candidates []int, frequency map[int]int) int {
	if len(candidates) == 0 {
		return 0
	}

	// 候補が1つの場合はそれを返す
	if len(candidates) == 1 {
		return candidates[0]
	}

	// 直近1回、2回の結果を取得して数字のセットを作成
	recent1Numbers := make(map[int]bool)
	recent2Numbers := make(map[int]bool)

	if len(re.results) > 0 {
		// 直近1回目の数字
		for _, numStr := range re.results[0] {
			if num, err := strconv.Atoi(numStr); err == nil {
				recent1Numbers[num] = true
			}
		}
	}

	if len(re.results) > 1 {
		// 直近2回目の数字
		for _, numStr := range re.results[1] {
			if num, err := strconv.Atoi(numStr); err == nil {
				recent2Numbers[num] = true
			}
		}
	}

	// 各候補の重みを計算（直近ペナルティ適用）
	totalWeight := 0.0
	weights := make([]float64, len(candidates))
	for i, num := range candidates {
		// 頻度の平方根を取って重みの差を緩和
		freq := float64(frequency[num])
		if freq == 0 {
			freq = 1.0
		}
		weight := math.Sqrt(freq)

		// 直近1回目に出現した数字：重みを減らす（×0.2）
		if recent1Numbers[num] {
			weight *= 0.2
		} else if recent2Numbers[num] {
			// 直近2回目に出現した数字：重みをやや減らす（×0.6）
			weight *= 0.6
		}

		weights[i] = weight
		totalWeight += weight
	}

	// 重みに基づいてランダム選択
	randomValue := re.rng.Float64() * totalWeight
	cumulativeWeight := 0.0
	for i, weight := range weights {
		cumulativeWeight += weight
		if randomValue < cumulativeWeight {
			return candidates[i]
		}
	}

	// フォールバック（通常ここには到達しない）
	return candidates[len(candidates)-1]
}

// 複数のランダム組み合わせを生成する
func (re *RandomEngine) GenerateMultipleRandomCombinations(count int) ([][]int, error) {
	if count <= 0 {
		return nil, fmt.Errorf("生成数は1以上である必要があります")
	}

	combinations := make([][]int, 0, count)
	usedCombinations := make(map[string]bool)

	maxAttempts := count * 100 // 最大試行回数

	for attempt := 0; attempt < maxAttempts && len(combinations) < count; attempt++ {
		combination, err := re.GenerateRandomCombination()
		if err != nil {
			continue
		}

		// 組み合わせの重複チェック
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

// 各位置の出現範囲を取得する
func (re *RandomEngine) GetPositionRanges() []PositionRange {
	return re.ranges
}

// 組み合わせをキー文字列に変換する
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
