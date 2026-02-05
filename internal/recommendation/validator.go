package recommendation

import (
	"sort"
	"strconv"
)

const (
	// バリデーション用の定数
	TotalSumMin        = 100  // 合計値の最小値
	TotalSumMax        = 160  // 合計値の最大値
	AverageValueMin    = 16.0 // 平均値の最小値
	AverageValueMax    = 23.0 // 平均値の最大値
	SmoothnessMaxVar   = 50.0 // 滑らかさの最大分散
	CloseNumberMinDiff = 1    // 近距離ペアの最小差
	CloseNumberMaxDiff = 3    // 近距離ペアの最大差
)

// 組み合わせの妥当性検証を行うインターフェース
type Validator interface {
	Validate(combination []int) bool
}

// 複数のバリデータを組み合わせて検証する実装
type CompositeValidator struct {
	validators []Validator
}

// 複合バリデータを作成する
func NewCompositeValidator(validators ...Validator) *CompositeValidator {
	return &CompositeValidator{
		validators: validators,
	}
}

// すべてのバリデータで検証を実行する
func (v *CompositeValidator) Validate(combination []int) bool {
	for _, validator := range v.validators {
		if !validator.Validate(combination) {
			return false
		}
	}
	return true
}

// ゾーン分布をバリデーションする実装
type ZoneDistributionValidator struct{}

// ゾーン分布が妥当かチェックする
func (v *ZoneDistributionValidator) Validate(combination []int) bool {
	zone1 := 0 // 1-13
	zone2 := 0 // 14-26
	zone3 := 0 // 27-37

	for _, num := range combination {
		if num >= 1 && num <= 13 {
			zone1++
		} else if num >= 14 && num <= 26 {
			zone2++
		} else if num >= 27 && num <= 37 {
			zone3++
		}
	}

	// 各ゾーンが空でないこと
	if zone1 == 0 || zone2 == 0 || zone3 == 0 {
		return false
	}

	// ゾーン1は3-4個
	if zone1 < 3 || zone1 > 4 {
		return false
	}

	// ゾーン3は2-3個
	if zone3 < 2 || zone3 > 3 {
		return false
	}

	return true
}

// 奇数偶数比率をバリデーションする実装
type OddEvenRatioValidator struct{}

// 奇数偶数比率が妥当かチェックする
func (v *OddEvenRatioValidator) Validate(combination []int) bool {
	oddCount := 0
	evenCount := 0

	for _, num := range combination {
		if num%2 == 0 {
			evenCount++
		} else {
			oddCount++
		}
	}

	// 奇4:偶3 または 奇3:偶4
	return (oddCount == 4 && evenCount == 3) || (oddCount == 3 && evenCount == 4)
}

// 合計値をバリデーションする実装
type TotalSumValidator struct{}

// 合計値が妥当かチェックする
func (v *TotalSumValidator) Validate(combination []int) bool {
	sum := 0
	for _, num := range combination {
		sum += num
	}
	return sum >= TotalSumMin && sum <= TotalSumMax
}

// 平均値をバリデーションする実装
type AverageValueValidator struct{}

// 平均値が妥当かチェックする
func (v *AverageValueValidator) Validate(combination []int) bool {
	if len(combination) == 0 {
		return false
	}
	sum := 0
	for _, num := range combination {
		sum += num
	}
	average := float64(sum) / float64(len(combination))
	// 平均値が16~23の範囲内
	return average >= AverageValueMin && average <= AverageValueMax
}

// 数字の並び滑らかさをバリデーションする実装
type SmoothnessValidator struct{}

// 数字の並びが滑らかかチェックする
func (v *SmoothnessValidator) Validate(combination []int) bool {
	if len(combination) < 2 {
		return true
	}

	// 隣接する数字間の差分を計算
	diffs := []int{}
	for i := 1; i < len(combination); i++ {
		diffs = append(diffs, combination[i]-combination[i-1])
	}

	// 差分の標準偏差を計算
	avgDiff := 0.0
	for _, diff := range diffs {
		avgDiff += float64(diff)
	}
	avgDiff /= float64(len(diffs))

	variance := 0.0
	for _, diff := range diffs {
		variance += (float64(diff) - avgDiff) * (float64(diff) - avgDiff)
	}
	variance /= float64(len(diffs))

	stdDev := variance // 平方根は計算せず分散で判定

	// 分散が大きすぎる（ジグザグが激しい）場合は却下
	return stdDev <= SmoothnessMaxVar
}

// 近距離ペアの存在をバリデーションする実装
type CloseNumberPairValidator struct{}

// 近距離ペア(差1-3)が含まれているかチェックする
func (v *CloseNumberPairValidator) Validate(combination []int) bool {
	for i := 0; i < len(combination); i++ {
		for j := i + 1; j < len(combination); j++ {
			diff := abs(combination[i] - combination[j])
			if diff >= CloseNumberMinDiff && diff <= CloseNumberMaxDiff {
				return true
			}
		}
	}
	return false
}

// 絶対値を計算する
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// 過去の組み合わせと重複していないかバリデーションする実装
type PastCombinationValidator struct {
	recentResults [][]string
	checkCount    int // 直近何回分をチェックするか
}

// 過去組み合わせバリデータを作成する
func NewPastCombinationValidator(recentResults [][]string, checkCount int) *PastCombinationValidator {
	return &PastCombinationValidator{
		recentResults: recentResults,
		checkCount:    checkCount,
	}
}

// 過去の組み合わせと重複していないかチェックする
func (v *PastCombinationValidator) Validate(combination []int) bool {
	// 組み合わせをソート
	sortedCombination := make([]int, len(combination))
	copy(sortedCombination, combination)
	sort.Ints(sortedCombination)

	// 過去の抽選結果から、推薦番号の中の任意の1つが出現した回を探す
	matchingDraws := [][]int{}
	checkRange := v.checkCount
	if checkRange > len(v.recentResults) {
		checkRange = len(v.recentResults)
	}

	for i := 0; i < checkRange; i++ {
		// その回の数字を取得
		drawNumbers := make([]int, 0, len(v.recentResults[i]))
		for _, numStr := range v.recentResults[i] {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				continue
			}
			drawNumbers = append(drawNumbers, num)
		}

		// 推薦番号のいずれかの数字が含まれているかチェック
		hasMatch := false
		for _, recNum := range sortedCombination {
			for _, drawNum := range drawNumbers {
				if recNum == drawNum {
					hasMatch = true
					break
				}
			}
			if hasMatch {
				break
			}
		}

		// マッチした場合、その回の数字を記録
		if hasMatch {
			matchingDraws = append(matchingDraws, drawNumbers)
		}

		// 直近2回分を確認したら終了
		if len(matchingDraws) >= 2 {
			break
		}
	}

	// マッチした過去の抽選結果と、推薦組み合わせが同じでないかチェック
	for _, pastDraw := range matchingDraws {
		sortedPastDraw := make([]int, len(pastDraw))
		copy(sortedPastDraw, pastDraw)
		sort.Ints(sortedPastDraw)

		// 完全一致チェック
		if len(sortedCombination) == len(sortedPastDraw) {
			isIdentical := true
			for j := 0; j < len(sortedCombination); j++ {
				if sortedCombination[j] != sortedPastDraw[j] {
					isIdentical = false
					break
				}
			}
			// 同じ組み合わせが過去に存在した場合は却下
			if isIdentical {
				return false
			}
		}
	}

	return true
}
