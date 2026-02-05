package recommendation

import "strconv"

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
			score += 2.0
			break
		}
	}

	// 2. 出現回数が多い数字は中程度評価
	for _, freq := range stats.FrequentNumbers {
		if freq == num {
			score += 1.0
			break
		}
	}

	// 3. 復活候補（20回以上未出現）は中程度評価
	for _, revival := range stats.RevivalCandidates {
		if revival == num {
			score += 1.0
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
		score -= 0.5
	}

	// 5. 直近3回で出た数字は大幅減点
	recentAvoid := s.config.RecentAvoidCount
	if recentAvoid > len(recentResults) {
		recentAvoid = len(recentResults)
	}
	for i := 0; i < recentAvoid; i++ {
		for _, numStr := range recentResults[i] {
			if n, _ := strconv.Atoi(numStr); n == num {
				score -= 3.0
				break
			}
		}
	}

	return score
}
