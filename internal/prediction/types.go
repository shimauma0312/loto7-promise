// 統計分析予測機能の共通型定義
package prediction

import "time"

// FilterWeights は各フィルターの重みを保持する（0.0〜1.0）
type FilterWeights struct {
	HotCold   float64 `json:"hot_cold"`
	Parity    float64 `json:"parity"`
	SizeBal   float64 `json:"size_balance"`
	SumRange  float64 `json:"sum_range"`
	TensGroup float64 `json:"tens_group"`
	LastDigit float64 `json:"last_digit"`
	Pull      float64 `json:"pull"`
	Bonus     float64 `json:"bonus"`
	Interval  float64 `json:"interval"`
}

// DefaultWeights はデフォルトのフィルター重みを返す
func DefaultWeights() FilterWeights {
	return FilterWeights{
		HotCold:   1.0,
		Parity:    1.0,
		SizeBal:   1.0,
		SumRange:  1.0,
		TensGroup: 1.0,
		LastDigit: 1.0,
		Pull:      0.8,
		Bonus:     0.8,
		Interval:  0.6,
	}
}

// PredictionConfig は予測生成の設定を保持する
type PredictionConfig struct {
	Count   int           `json:"count"`
	History int           `json:"history"`
	Weights FilterWeights `json:"weights"`
}

// DefaultConfig はデフォルト設定を返す
func DefaultConfig() PredictionConfig {
	return PredictionConfig{
		Count:   5,
		History: 100,
		Weights: DefaultWeights(),
	}
}

// AnalysisData は過去データの統計分析結果を保持する
type AnalysisData struct {
	HotNumbers       []int       // 直近ウィンドウ内で2回以上出現した数字
	ColdNumbers      []int       // 直近ウィンドウ内で未出現の数字
	LastDrawNumbers  []int       // 前回の本数字（7個）
	LastBonusNumbers []int       // 前回のボーナス数字（存在する場合）
	LastDrawSum      int         // 前回の合計値
	PrevDrawSum      int         // 前々回の合計値
	SumTrend         int         // -1=2回連続上昇中, +1=2回連続下降中, 0=変動なし
	FrequencyMap     map[int]int // 全期間の出現頻度
	IntervalMap      map[int]int // 各数字が現在何回連続で出ていないか（最新から）
	AnalyzedDraws    int
}

// ScoreDetail はフィルターごとのスコア内訳を保持する
type ScoreDetail struct {
	HotColdScore   float64 `json:"hot_cold_score"`
	ParityScore    float64 `json:"parity_score"`
	SizeScore      float64 `json:"size_score"`
	SumScore       float64 `json:"sum_score"`
	TensScore      float64 `json:"tens_score"`
	LastDigitScore float64 `json:"last_digit_score"`
	PullScore      float64 `json:"pull_score"`
	BonusScore     float64 `json:"bonus_score"`
	IntervalScore  float64 `json:"interval_score"`
}

// ScoredCombination はスコア付きの組み合わせを表す
type ScoredCombination struct {
	ID          int        `json:"id"`
	Numbers     []int      `json:"numbers"`
	TotalScore  float64    `json:"total_score"`
	ScoreDetail ScoreDetail `json:"score_detail"`
}

// PredictionAnalysisInfo は予測生成時の分析メタデータを保持する
type PredictionAnalysisInfo struct {
	AnalyzedDraws    int       `json:"analyzed_draws"`
	HotNumbers       []int     `json:"hot_numbers"`
	ColdNumbers      []int     `json:"cold_numbers"`
	LastDrawNumbers  []int     `json:"last_draw_numbers"`
	LastBonusNumbers []int     `json:"last_bonus_numbers"`
	LastDrawSum      int       `json:"last_draw_sum"`
	SumTrend         string    `json:"sum_trend"`
	GeneratedAt      time.Time `json:"generated_at"`
}

// PredictionResult は予測生成の最終結果を保持する
type PredictionResult struct {
	Combinations []ScoredCombination   `json:"combinations"`
	AnalysisInfo PredictionAnalysisInfo `json:"analysis_info"`
	Config       PredictionConfig       `json:"config"`
}
