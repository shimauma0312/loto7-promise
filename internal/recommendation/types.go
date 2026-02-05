package recommendation

// RecommendationConfig は推薦設定を保持します
type RecommendationConfig struct {
	MaxRecommendations int     // 最大推薦数
	RecentAvoidCount   int     // 直近の回避回数
	ConsecutiveBoost   int     // 連続出現ブースト判定回数
	HistoryLookback    int     // 過去の履歴確認回数
	FrequencyWeight    float64 // 出現頻度の重み
	RecentWeight       float64 // 最近の出現の重み
	ConsecutiveWeight  float64 // 連続出現の重み
	PositionWeight     float64 // 位置による重み
}

// DefaultConfig はデフォルト設定を返します
func DefaultConfig() RecommendationConfig {
	return RecommendationConfig{
		MaxRecommendations: 5,
		RecentAvoidCount:   3,
		ConsecutiveBoost:   3,
		HistoryLookback:    100,
		FrequencyWeight:    0.3,
		RecentWeight:       0.4,
		ConsecutiveWeight:  0.2,
		PositionWeight:     0.1,
	}
}

// NumberScore は数字のスコアを保持します
type NumberScore struct {
	Number           int
	FrequencyScore   float64
	RecentScore      float64
	ConsecutiveScore float64
	PositionScore    float64
	TotalScore       float64
	IsRecentlyDrawn  bool
}

// StatisticalAnalysis は統計分析結果を保持します
type StatisticalAnalysis struct {
	FrequentNumbers   []int       // 出現回数が多い数字
	RareNumbers       []int       // 出現回数が少ない数字
	HotNumbers        []int       // 直近10回以内で連続して出ている数字
	RevivalCandidates []int       // 20回以上出ていない復活候補
	FrequencyMap      map[int]int // 各数字の出現回数
	LastAppearance    map[int]int // 各数字の最後の出現回数
}

// CombinationHistory は組み合わせ履歴を保持します
type CombinationHistory struct {
	Pairs map[string]int // "num1,num2"形式のペア出現回数
}
