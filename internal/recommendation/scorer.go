package recommendation

import "strconv"

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
	CalculatePriority(num int, stats StatisticalAnalysis, recentResults [][]string) float64
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
func (s *WeightedScorer) CalculatePriority(num int, stats StatisticalAnalysis, recentResults [][]string) float64 {
	score := 0.0

	// 1. 直近1回の抽選結果での数字は、出現率を限りなく低いレベルで下げる
	if len(recentResults) > 0 {
		for _, numStr := range recentResults[0] {
			n, err := strconv.Atoi(numStr)
			if err != nil {
				continue
			}
			if n == num {
				return ScoreRecent1Penalty // 即座に大幅減点で返す
			}
		}
	}

	// 2. 直近2回目の抽選結果での数字は、出現率をそれなりに下げる
	if len(recentResults) > 1 {
		// 直近2回目（インデックス1）をチェック
		for _, numStr := range recentResults[1] {
			n, err := strconv.Atoi(numStr)
			if err != nil {
				continue
			}
			if n == num {
				score += ScoreRecent2Penalty
				break
			}
		}
	}

	// 3. 直近100回で出現していない数字は、出現率が0になる
	freq, exists := stats.FrequencyMap[num]
	if !exists || freq == 0 {
		return ScoreNoAppearancePenalty
	}

	// 4. 直近100回で出現回数が5回未満の数字は、少しばかし出現率をマイナスする
	if freq < LowFrequencyThreshold {
		score += ScoreLowFrequencyPenalty
	}

	// 5. 直近3回は除き、50回以内に出現した数字は、出現回数分ブーストされる
	// 直近3回を除いた範囲（3回目から50回目）をチェック
	checkStart := RecentAvoidRange
	checkEnd := Recent50CheckRange
	if checkEnd > len(recentResults) {
		checkEnd = len(recentResults)
	}

	within50Count := 0
	for i := checkStart; i < checkEnd; i++ {
		for _, numStr := range recentResults[i] {
			n, err := strconv.Atoi(numStr)
			if err != nil {
				continue
			}
			if n == num {
				within50Count++
			}
		}
	}
	score += float64(within50Count) * ScoreWithin50Boost

	return score
}
