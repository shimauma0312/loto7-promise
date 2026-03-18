package recommendation

import (
	"sort"
	"strconv"
)

const (
	// スコア計算用の定数
	ScoreRecent1Penalty      = -100.0  // 直近1回の抽選結果のペナルティ（限りなく低い）
	ScoreRecent2Penalty      = -10.0   // 直近2回までのペナルティ（それなりに低い）
	ScoreWithin50Boost       = 0.5     // 直近3回を除き、50回以内の出現ブースト（出現回数分）
	ScoreLowFrequencyPenalty = -1.0    // 直近100回で5回未満のペナルティ
	ScoreNoAppearancePenalty = -1000.0 // 直近100回で未出現のペナルティ（出現率0）
	Recent50CheckRange       = 50      // 50回以内チェック範囲
	RecentAvoidRange         = 3       // 直近3回は除く
	LowFrequencyThreshold    = 5       // 低頻度判定閾値
)

// 数字の優先度計算を行うインターフェース
type Scorer interface {
	// position は 0-indexed（0=昇順1番目 … 6=昇順7番目）
	CalculatePriority(num int, position int, stats StatisticalAnalysis, recentResults [][]string) float64
}

// 重み付きでスコアを計算する実装
type WeightedScorer struct {
	config RecommendationConfig
}

// 重み付きスコアラーを作成する
func NewWeightedScorer(config RecommendationConfig) *WeightedScorer {
	return &WeightedScorer{
		config: config,
	}
}

// 数字の優先度を計算する
//
// position は 0-indexed。直近出現ペナルティはポジション単位で適用し、
// 同ポジションに同じ数字が出た回の近さのみをペナルティ対象とする。
func (s *WeightedScorer) CalculatePriority(num int, position int, stats StatisticalAnalysis, recentResults [][]string) float64 {
	score := 0.0

	// 1. 同ポジションでの直近出現チェック
	if position >= 0 && position < len(stats.PositionLastAppearance) {
		posLast := stats.PositionLastAppearance[position]
		if lastIdx, exists := posLast[num]; exists {
			if lastIdx == 0 {
				return ScoreRecent1Penalty
			}
			if lastIdx == 1 {
				score += ScoreRecent2Penalty
			}
		}
	} else {
		// PositionLastAppearance が未設定のフォールバック（全体チェック）
		if len(recentResults) > 0 {
			for _, numStr := range recentResults[0] {
				n, err := strconv.Atoi(numStr)
				if err == nil && n == num {
					return ScoreRecent1Penalty
				}
			}
		}
		if len(recentResults) > 1 {
			for _, numStr := range recentResults[1] {
				n, err := strconv.Atoi(numStr)
				if err == nil && n == num {
					score += ScoreRecent2Penalty
					break
				}
			}
		}
	}

	// 2. 直近100回で出現していない数字は、出現率が0になる
	freq, exists := stats.FrequencyMap[num]
	if !exists || freq == 0 {
		return ScoreNoAppearancePenalty
	}

	// 3. 直近100回で出現回数が5回未満の数字は、少しばかし出現率をマイナスする
	if freq < LowFrequencyThreshold {
		score += ScoreLowFrequencyPenalty
	}

	// 4. 直近3回は除き、同ポジションで50回以内に出現した数字は、出現回数分ブーストされる
	checkStart := RecentAvoidRange
	checkEnd := Recent50CheckRange
	if checkEnd > len(recentResults) {
		checkEnd = len(recentResults)
	}

	within50Count := 0
	if position >= 0 && position < len(stats.PositionFrequencyMap) {
		// ポジション単位：同ポジションに出現した回数のみカウント
		for i := checkStart; i < checkEnd; i++ {
			draw := parseDraw(recentResults[i])
			sortedDraw := make([]int, len(draw))
			copy(sortedDraw, draw)
			sort.Ints(sortedDraw)
			if position < len(sortedDraw) && sortedDraw[position] == num {
				within50Count++
			}
		}
	} else {
		// グローバルフォールバック（PositionFrequencyMap 未設定時）
		for i := checkStart; i < checkEnd; i++ {
			for _, numStr := range recentResults[i] {
				n, err := strconv.Atoi(numStr)
				if err == nil && n == num {
					within50Count++
				}
			}
		}
	}
	score += float64(within50Count) * ScoreWithin50Boost

	return score
}
