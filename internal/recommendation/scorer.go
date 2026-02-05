package recommendation

import "strconv"

const (
	// スコア計算用の重み
	ScoreHotNumberBonus      = 2.0  // 波が来ている数字のボーナス
	ScoreFrequentNumberBonus = 1.0  // 出現頻度が高い数字のボーナス
	ScoreRevivalBonus        = 1.0  // 復活候補のボーナス
	ScoreRarePenalty         = -0.5 // 出現頻度が少ない数字のペナルティ
	ScoreRecentPenalty       = -3.0 // 直近で出た数字のペナルティ
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

	// 1. 波が来ている数字（直近10回で3回以上）は高評価
	for _, hot := range stats.HotNumbers {
		if hot == num {
			score += ScoreHotNumberBonus
			break
		}
	}

	// 2. 出現回数が多い数字は中程度評価
	for _, freq := range stats.FrequentNumbers {
		if freq == num {
			score += ScoreFrequentNumberBonus
			break
		}
	}

	// 3. 復活候補（20回以上未出現）は中程度評価
	for _, revival := range stats.RevivalCandidates {
		if revival == num {
			score += ScoreRevivalBonus
			break
		}
	}

	// 4. 出現回数が極端に少ない数字は減点
	isRare := false
	for _, rare := range stats.RareNumbers {
		if rare == num {
			isRare = true
			break
		}
	}
	if isRare {
		score += ScoreRarePenalty
	}

	// 5. 直近3回で出た数字は大幅減点
	recentAvoid := s.config.RecentAvoidCount
	if recentAvoid > len(recentResults) {
		recentAvoid = len(recentResults)
	}
	for i := 0; i < recentAvoid; i++ {
		for _, numStr := range recentResults[i] {
			n, err := strconv.Atoi(numStr)
			if err != nil {
				// 不正なフォーマットの場合はスキップ
				continue
			}
			if n == num {
				score += ScoreRecentPenalty
				break
			}
		}
	}

	return score
}
