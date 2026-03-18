package recommendation

// RecommendationConfig は推薦設定を保持します
type RecommendationConfig struct {
	MaxRecommendations int     // 最大推薦数
	HistoryLookback    int     // 過去の履歴確認回数
	FrequencyWeight    float64 // 出現頻度の重み
	RecentWeight       float64 // 最近の出現の重み
}

// DefaultConfig はデフォルト設定を返します
func DefaultConfig() RecommendationConfig {
	return RecommendationConfig{
		MaxRecommendations: 5,
		HistoryLookback:    100,
		FrequencyWeight:    0.5,
		RecentWeight:       0.5,
	}
}

// NumberScore は数字のスコアを保持します
type NumberScore struct {
	Number          int
	FrequencyScore  float64
	RecentScore     float64
	TotalScore      float64
	IsRecentlyDrawn bool
}

// StatisticalAnalysis は統計分析結果を保持します
type StatisticalAnalysis struct {
	FrequentNumbers   []int       // 出現回数が多い数字
	RareNumbers       []int       // 出現回数が少ない数字
	HotNumbers        []int       // 直近10回以内で連続して出ている数字
	RevivalCandidates []int       // 20回以上出ていない復活候補
	FrequencyMap      map[int]int // 各数字の出現回数（全ポジション合算）
	LastAppearance    map[int]int // 各数字の最後の出現回数（全ポジション合算）

	// ポジション別統計（index 0 = 昇順1番目、index 6 = 昇順7番目）
	PositionFrequencyMap   []map[int]int // PositionFrequencyMap[pos][num] = 出現回数
	PositionLastAppearance []map[int]int // PositionLastAppearance[pos][num] = 最終出現インデックス（0=直近）
	// scorer 最適化用キャッシュ：直近 RecentAvoidRange..Recent50CheckRange 回分のソート済みドロー
	// 事前計算することで CalculatePriority 内のパース・ソートを割愛
	RecentSortedDraws [][]int
}

// CombinationHistory は組み合わせ履歴を保持します
type CombinationHistory struct {
	Pairs map[string]int // "num1,num2"形式のペア出現回数
}
