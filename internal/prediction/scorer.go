// 各フィルター基準に基づいて組み合わせをスコアリングするモジュール
package prediction

import "math"

// Score は組み合わせの総合スコアと内訳を計算する
func Score(combo []int, data AnalysisData, weights FilterWeights) (float64, ScoreDetail) {
	detail := ScoreDetail{}

	detail.HotColdScore = scoreHotCold(combo, data)
	detail.ParityScore = scorePassFail(passesParityFilter(combo))
	detail.SizeScore = scorePassFail(passesSizeFilter(combo))
	detail.SumScore = scoreSumRange(combo, data)
	detail.TensScore = scoreTensGroup(combo)
	detail.LastDigitScore = scoreLastDigit(combo)
	detail.PullScore = scorePull(combo, data)
	detail.BonusScore = scoreBonus(combo, data)
	detail.IntervalScore = scoreInterval(combo, data)

	total := detail.HotColdScore*weights.HotCold +
		detail.ParityScore*weights.Parity +
		detail.SizeScore*weights.SizeBal +
		detail.SumScore*weights.SumRange +
		detail.TensScore*weights.TensGroup +
		detail.LastDigitScore*weights.LastDigit +
		detail.PullScore*weights.Pull +
		detail.BonusScore*weights.Bonus +
		detail.IntervalScore*weights.Interval

	return total, detail
}

// ① ホット/コールドスコア
// 理想: ホット5〜6個 + コールド1〜2個
func scoreHotCold(combo []int, data AnalysisData) float64 {
	hotSet := toSet(data.HotNumbers)
	coldSet := toSet(data.ColdNumbers)

	hotCount, coldCount := 0, 0
	for _, n := range combo {
		if hotSet[n] {
			hotCount++
		}
		if coldSet[n] {
			coldCount++
		}
	}

	// ホット数が5に近いほど高スコア
	hotScore := math.Max(0, 1.0-math.Abs(float64(hotCount)-5.0)/5.0)

	// コールドが1〜2個なら満点
	coldScore := 0.0
	if coldCount >= 1 && coldCount <= 2 {
		coldScore = 1.0
	} else if coldCount == 0 {
		coldScore = 0.2
	}

	return (hotScore + coldScore) / 2.0
}

// passFailスコア: 通過=1.0, 不合格=0.0
func scorePassFail(pass bool) float64 {
	if pass {
		return 1.0
	}
	return 0.0
}

// ④ 合計スコア: 範囲中央(132.5)に近いほど高スコア、トレンド補正あり
func scoreSumRange(combo []int, data AnalysisData) float64 {
	s := float64(sumInts(combo))

	// 範囲(95-170)内での距離
	const center = 132.5
	const halfRange = 37.5
	dist := math.Abs(s-center) / halfRange
	score := math.Max(0, 1.0-dist)

	// 合計トレンド補正
	switch data.SumTrend {
	case -1: // 2回連続上昇 → 下降傾向ほうを優先
		if int(s) < data.LastDrawSum {
			score = math.Min(1.0, score*1.2)
		} else {
			score *= 0.8
		}
	case 1: // 2回連続下降 → 上昇傾向を優先
		if int(s) > data.LastDrawSum {
			score = math.Min(1.0, score*1.2)
		} else {
			score *= 0.8
		}
	}

	return score
}

// ⑤ 十の位グループスコア: 欠落グループ数が多いほど高スコア
// 0番台(1-9)を含む場合はボーナス
func scoreTensGroup(combo []int) float64 {
	groups := make(map[int]bool, 4)
	for _, n := range combo {
		groups[tensGroup(n)] = true
	}

	missing := 4 - len(groups)
	score := 0.0
	switch missing {
	case 2:
		score = 1.0
	case 1:
		score = 0.8
	default:
		score = 0.0
	}

	// 0番台(1-9)が含まれていればボーナス
	if groups[0] {
		score = math.Min(1.0, score+0.1)
	}

	return score
}

// ⑥ 一の位の種類数スコア
// 5種類=最高点, 6種類=次点, それ以外=低点
func scoreLastDigit(combo []int) float64 {
	lastDigits := make(map[int]bool, 7)
	for _, n := range combo {
		lastDigits[n%10] = true
	}
	switch len(lastDigits) {
	case 5:
		return 1.0
	case 6:
		return 0.9
	case 4:
		return 0.4
	default:
		return 0.0
	}
}

// ⑦ 引っ張り・斜め数字スコア
// 前回から1〜2個継続 + 前回±1の数字も加点
func scorePull(combo []int, data AnalysisData) float64 {
	if len(data.LastDrawNumbers) == 0 {
		return 0.5
	}

	lastSet := toSet(data.LastDrawNumbers)
	pullCount, diagCount := 0, 0

	for _, n := range combo {
		if lastSet[n] {
			pullCount++
		}
		if lastSet[n-1] || lastSet[n+1] {
			diagCount++
		}
	}

	// 理想: 1〜2個の引っ張り
	score := 0.0
	switch {
	case pullCount >= 1 && pullCount <= 2:
		score = 1.0
	case pullCount == 3:
		score = 0.5
	case pullCount == 0:
		score = 0.3
	default:
		score = 0.2
	}

	// 斜め数字ボーナス（最大+0.2）
	if diagCount > 0 {
		score = math.Min(1.0, score+0.1*math.Min(float64(diagCount), 2))
	}

	return score
}

// ⑧ ボーナス数字周辺スコア
// 前回ボーナス数字±3の範囲にある数字が含まれているほど高スコア
func scoreBonus(combo []int, data AnalysisData) float64 {
	if len(data.LastBonusNumbers) == 0 {
		return 0.5
	}

	zoneCount := 0
	for _, n := range combo {
		for _, bonus := range data.LastBonusNumbers {
			if intAbs(n-bonus) <= 3 {
				zoneCount++
				break
			}
		}
	}

	if zoneCount >= 2 {
		return 1.0
	} else if zoneCount == 1 {
		return 0.7
	}
	return 0.2
}

// 当選間隔スコア
// 理想: 間隔10以内が5〜6個、間隔20以上が1〜2個
func scoreInterval(combo []int, data AnalysisData) float64 {
	shortCount, longCount := 0, 0
	for _, n := range combo {
		interval := data.IntervalMap[n]
		if interval <= 10 {
			shortCount++
		} else if interval >= 20 {
			longCount++
		}
	}

	switch {
	case shortCount >= 5 && longCount >= 1:
		return 1.0
	case shortCount >= 4 && longCount >= 1:
		return 0.8
	case shortCount >= 5:
		return 0.6
	case shortCount >= 4:
		return 0.5
	default:
		return 0.2
	}
}
