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

// 最適化可能なスコアリング重みパラメータ
type ScorerWeights struct {
	Recent1Penalty         float64    // 同ポジションで直近1回出現したペナルティ
	Recent2Penalty         float64    // 同ポジションで直近2回目出現したペナルティ
	Within50Boost          float64    // 直近3-50回・同ポジション出現1回あたりのブースト
	LowFrequencyPenalty    float64    // 直近100回で 5 回未満のペナルティ
	NoAppearancePenalty    float64    // 直近100回で未出現のペナルティ
	PositionPenaltyScale   [7]float64 // ポジション別ペナルティスケール係数
	PositionFrequencyWeight float64   // ポジション別履歴分布に基づく適合度ボーナスの重み
}

// 現在の定数値をそのまま使ったデフォルト重みを返す
func DefaultScorerWeights() ScorerWeights {
	return ScorerWeights{
		Recent1Penalty:          ScoreRecent1Penalty,
		Recent2Penalty:          ScoreRecent2Penalty,
		Within50Boost:           ScoreWithin50Boost,
		LowFrequencyPenalty:     ScoreLowFrequencyPenalty,
		NoAppearancePenalty:     ScoreNoAppearancePenalty,
		PositionPenaltyScale:    [7]float64{1, 1, 1, 1, 1, 1, 1},
		PositionFrequencyWeight: 2.0, // 平均出現率の数字に +2.0 程度のボーナスがかかる設定
	}
}

// 数字の優先度計算を行うインターフェース
type Scorer interface {
	// position は 0-indexed（0=昇順1番目 … 6=昇順7番目）
	// selected は選択済みの数字（EnsembleScorer でのベイズ条件付き確率計算に使用）
	CalculatePriority(num int, position int, stats StatisticalAnalysis, recentResults [][]string, selected []int) float64
}

// 重み付きでスコアを計算する実装
type WeightedScorer struct {
	config  RecommendationConfig
	weights ScorerWeights
}

// デフォルト重みで重み付きスコアラーを作成する
func NewWeightedScorer(config RecommendationConfig) *WeightedScorer {
	return &WeightedScorer{config: config, weights: DefaultScorerWeights()}
}

// 指定した重みで重み付きスコアラーを作成する
func NewWeightedScorerWithWeights(config RecommendationConfig, weights ScorerWeights) *WeightedScorer {
	return &WeightedScorer{config: config, weights: weights}
}

// 数字の優先度を計算する
//
// position は 0-indexed。直近出現ペナルティはポジション単位で適用し、
// 同ポジションに同じ数字が出た回の近さのみをペナルティ対象とする。
// selected は EnsembleScorer のためのパラメータで、WeightedScorer では使用しない。
func (s *WeightedScorer) CalculatePriority(num int, position int, stats StatisticalAnalysis, recentResults [][]string, selected []int) float64 {
	score := 0.0

	// 1. 同ポジションでの直近出現チェック
	if position >= 0 && position < len(stats.PositionLastAppearance) {
		posLast := stats.PositionLastAppearance[position]
		if lastIdx, exists := posLast[num]; exists {
			posScale := s.weights.PositionPenaltyScale[position]
			if lastIdx == 0 {
				return s.weights.Recent1Penalty * posScale
			}
			if lastIdx == 1 {
				score += s.weights.Recent2Penalty * posScale
			}
		}
	} else {
		// PositionLastAppearance が未設定のフォールバック（全体チェック）
		fbScale := 1.0
		if position >= 0 && position < 7 {
			fbScale = s.weights.PositionPenaltyScale[position]
		}
		if len(recentResults) > 0 {
			for _, numStr := range recentResults[0] {
				n, err := strconv.Atoi(numStr)
				if err == nil && n == num {
					return s.weights.Recent1Penalty * fbScale
				}
			}
		}
		if len(recentResults) > 1 {
			for _, numStr := range recentResults[1] {
				n, err := strconv.Atoi(numStr)
				if err == nil && n == num {
					score += s.weights.Recent2Penalty * fbScale
					break
				}
			}
		}
	}

	// 2. 直近100回で出現していない数字は、出現率が0になる
	freq, exists := stats.FrequencyMap[num]
	if !exists || freq == 0 {
		return s.weights.NoAppearancePenalty
	}

	// 3. 直近100回で出現回数が5回未満の数字は、少しばかし出現率をマイナスする
	if freq < LowFrequencyThreshold {
		score += s.weights.LowFrequencyPenalty
	}

	// 4. 直近3回は除き、同ポジションで50回以内に出現した数字は、出現回数分ブーストされる
	checkStart := RecentAvoidRange
	checkEnd := Recent50CheckRange
	if checkEnd > len(recentResults) {
		checkEnd = len(recentResults)
	}

	within50Count := 0
	if len(stats.RecentSortedDraws) > 0 {
		// 事前計算済みキャッシュを使う（高速パス: int比較のみ）
		for _, sortedDraw := range stats.RecentSortedDraws {
			if position < len(sortedDraw) && sortedDraw[position] == num {
				within50Count++
			}
		}
	} else if position >= 0 && position < len(stats.PositionFrequencyMap) {
		// キャッシュなし・PositionFrequencyMap あり（従来のパース方式）
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
	score += float64(within50Count) * s.weights.Within50Boost

	// 5. ポジション別履歴分布に基づく適合度ボーナス
	// PositionFrequencyMap[pos][num] = 過去履歴で num が pos に出現した回数
	// 平均出現率（= 全ドロー数 / 37）と比べて多く出た数字ほど高ボーナス
	if s.weights.PositionFrequencyWeight != 0 &&
		position >= 0 && position < len(stats.PositionFrequencyMap) {
		posMap := stats.PositionFrequencyMap[position]
		totalAtPos := 0
		for _, f := range posMap {
			totalAtPos += f
		}
		if totalAtPos > 0 {
			posFreq := float64(posMap[num])
			// normalizedFreq: 1.0 = 均一分布と同じ、>1.0 = 平均より多い
			normalizedFreq := posFreq / float64(totalAtPos) * float64(lotoTotalNumbers)
			score += s.weights.PositionFrequencyWeight * normalizedFreq
		}
	}

	return score
}
