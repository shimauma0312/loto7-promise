// ハードフィルター群：条件を満たさない組み合わせを除外する
package prediction

// passesParityFilter は奇数/偶数バランスをチェックする
// 許容パターン: 奇4:偶3 / 奇3:偶4 / 奇5:偶2 / 奇2:偶5
func passesParityFilter(combo []int) bool {
	odd := 0
	for _, n := range combo {
		if n%2 == 1 {
			odd++
		}
	}
	even := len(combo) - odd
	switch {
	case odd == 4 && even == 3,
		odd == 3 && even == 4,
		odd == 5 && even == 2,
		odd == 2 && even == 5:
		return true
	default:
		return false
	}
}

// passesSizeFilter は大小バランスをチェックする
// 小:1〜19 / 大:20〜37
// 許容パターン: 小3:大4 / 小4:大3 / 小2:大5 / 小5:大2
func passesSizeFilter(combo []int) bool {
	small := 0
	for _, n := range combo {
		if n <= 19 {
			small++
		}
	}
	large := len(combo) - small
	switch {
	case small == 3 && large == 4,
		small == 4 && large == 3,
		small == 2 && large == 5,
		small == 5 && large == 2:
		return true
	default:
		return false
	}
}

// passesSumFilter は合計値の適正範囲をチェックする
// 範囲: 95〜170、かつ前回合計との差が60以内
func passesSumFilter(combo []int, lastSum int) bool {
	s := sumInts(combo)
	if s < 95 || s > 170 {
		return false
	}
	if lastSum > 0 {
		diff := s - lastSum
		if diff < 0 {
			diff = -diff
		}
		if diff > 60 {
			return false
		}
	}
	return true
}

// passesTensGroupFilter は十の位グループの「欠落」ルールをチェックする
// グループ: 0番台(1-9), 10番台(10-19), 20番台(20-29), 30番台(30-37)
// 少なくとも1グループが欠落していること
func passesTensGroupFilter(combo []int) bool {
	groups := make(map[int]bool, 4)
	for _, n := range combo {
		groups[tensGroup(n)] = true
	}
	return len(groups) < 4
}

// passesLastDigitFilter は一の位の種類数をチェックする
// 出現しやすい「5種類」または「6種類」のみ許容
func passesLastDigitFilter(combo []int) bool {
	lastDigits := make(map[int]bool, 7)
	for _, n := range combo {
		lastDigits[n%10] = true
	}
	kinds := len(lastDigits)
	return kinds == 5 || kinds == 6
}

// tensGroup は数字の十の位グループを返す（0/1/2/3）
func tensGroup(n int) int {
	switch {
	case n <= 9:
		return 0
	case n <= 19:
		return 1
	case n <= 29:
		return 2
	default:
		return 3
	}
}

// intAbs は整数の絶対値を返す
func intAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// toSet はスライスをセット(map)に変換する
func toSet(nums []int) map[int]bool {
	set := make(map[int]bool, len(nums))
	for _, n := range nums {
		set[n] = true
	}
	return set
}
