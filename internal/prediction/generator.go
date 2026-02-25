// 候補組み合わせの生成とハードフィルタリングを担当するモジュール
package prediction

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

const (
	maxNumber = 37
	comboSize = 7
	// ハードフィルター後に残す最低候補数
	minFilteredPool = 500
)

// generateCandidatePool はランダムな候補組み合わせのプールを生成する
func generateCandidatePool(rng *rand.Rand, size int) [][]int {
	pool := make([][]int, 0, size)
	seen := make(map[string]bool, size)

	maxAttempts := size * 15
	for i := 0; i < maxAttempts && len(pool) < size; i++ {
		combo := randomCombo(rng)
		key := comboKey(combo)
		if !seen[key] {
			seen[key] = true
			pool = append(pool, combo)
		}
	}
	return pool
}

// randomCombo はランダムな7数字の組み合わせを生成する
func randomCombo(rng *rand.Rand) []int {
	nums := make([]int, maxNumber)
	for i := range nums {
		nums[i] = i + 1
	}
	rng.Shuffle(len(nums), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})
	combo := make([]int, comboSize)
	copy(combo, nums[:comboSize])
	sort.Ints(combo)
	return combo
}

// comboKey は組み合わせをユニークな文字列キーに変換する
func comboKey(combo []int) string {
	var b strings.Builder
	for i, n := range combo {
		if i > 0 {
			b.WriteByte('-')
		}
		b.WriteString(strconv.Itoa(n))
	}
	return b.String()
}

// applyHardFilters はハードフィルターを順次適用し、通過した組み合わせを返す
// フィルター通過率が低すぎる場合は段階的に条件を緩和する
func applyHardFilters(pool [][]int, data AnalysisData) [][]int {
	out := make([][]int, 0, len(pool)/4)
	for _, combo := range pool {
		if !passesParityFilter(combo) {
			continue
		}
		if !passesSizeFilter(combo) {
			continue
		}
		if !passesSumFilter(combo, data.LastDrawSum) {
			continue
		}
		if !passesTensGroupFilter(combo) {
			continue
		}
		if !passesLastDigitFilter(combo) {
			continue
		}
		out = append(out, combo)
	}

	// フィルター後の候補が少なすぎる場合は合計・一の位フィルターを緩和
	if len(out) < minFilteredPool {
		out = applyRelaxedFilters(pool, data)
	}

	return out
}

// applyRelaxedFilters は合計・一の位フィルターを外した緩和条件を適用する
func applyRelaxedFilters(pool [][]int, data AnalysisData) [][]int {
	out := make([][]int, 0, len(pool)/2)
	for _, combo := range pool {
		if !passesParityFilter(combo) {
			continue
		}
		if !passesSizeFilter(combo) {
			continue
		}
		if !passesTensGroupFilter(combo) {
			continue
		}
		out = append(out, combo)
	}
	return out
}
