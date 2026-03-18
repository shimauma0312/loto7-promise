package api

import (
	"time"
)

// APIResponse 統一APIレスポンス構造
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Version   string      `json:"version"`
}

// APIError エラー情報
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse 成功レスポンスを作成
func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}
}

// ErrorResponse エラーレスポンスを作成
func ErrorResponse(code, message, details string) APIResponse {
	return APIResponse{
		Success: false,
		Message: "エラーが発生しました",
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}
}

// HeatmapResponse ヒートマップ機能のレスポンス
type HeatmapResponse struct {
	SearchRange      int                    `json:"search_range"`
	ActualDraws      int                    `json:"actual_draws"`
	HottestNumbers   []int                  `json:"hottest_numbers"`
	ColdestNumbers   []int                  `json:"coldest_numbers"`
	MostBiasedNumber int                    `json:"most_biased_number"`
	PositionData     map[string]interface{} `json:"position_data,omitempty"`
	Statistics       HeatmapStatistics      `json:"statistics"`
}

// HeatmapStatistics ヒートマップの統計情報
type HeatmapStatistics struct {
	TotalNumbers     int     `json:"total_numbers"`
	AverageFrequency float64 `json:"average_frequency"`
	MaxFrequency     int     `json:"max_frequency"`
	MinFrequency     int     `json:"min_frequency"`
}

// RecommendationResponse 推薦機能のレスポンス
type RecommendationResponse struct {
	Recommendations []RecommendationSet `json:"recommendations"`
	AnalysisInfo    AnalysisInfo        `json:"analysis_info"`
	Config          ConfigInfo          `json:"config"`
}

// RecommendationSet 個別の推薦セット
type RecommendationSet struct {
	ID      int     `json:"id"`
	Numbers []int   `json:"numbers"`
	Score   float64 `json:"score,omitempty"`
}

// AnalysisInfo 分析情報
type AnalysisInfo struct {
	DataSource     string    `json:"data_source"`
	AnalyzedDraws  int       `json:"analyzed_draws"`
	GeneratedAt    time.Time `json:"generated_at"`
	Algorithm      string    `json:"algorithm"`
	LearnedWeights bool      `json:"learned_weights"` // 学習済み重みを使用中かどうか
}

// ConfigInfo 設定情報
type ConfigInfo struct {
	MaxRecommendations int     `json:"max_recommendations"`
	HistoryLookback    int     `json:"history_lookback"`
	FrequencyWeight    float64 `json:"frequency_weight"`
	RecentWeight       float64 `json:"recent_weight"`
}

// ResultResponse 結果取得機能のレスポンス
type ResultResponse struct {
	Results        []DrawResultAPI `json:"results"`
	TotalCount     int             `json:"total_count"`
	RequestedCount int             `json:"requested_count"`
	DataSource     string          `json:"data_source"`
}

// DrawResultAPI 抽選結果API用
type DrawResultAPI struct {
	DrawNumber       int      `json:"draw_number"`
	Numbers          []string `json:"numbers"`
	Date             string   `json:"date"`
	FormattedNumbers string   `json:"formatted_numbers"`
}

// HealthResponse ヘルスチェックレスポンス
type HealthResponse struct {
	Status    string            `json:"status"`
	Uptime    string            `json:"uptime"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// PredictionFilterWeights は各フィルターの重み設定
type PredictionFilterWeights struct {
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

// PredictionRequest はPOST /api/v1/prediction のリクエストボディ
type PredictionRequest struct {
	Count   int                     `json:"count"`
	History int                     `json:"history"`
	Weights PredictionFilterWeights `json:"weights"`
}

// PredictionScoreDetail はフィルターごとのスコア内訳
type PredictionScoreDetail struct {
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

// PredictionCombination はスコア付き推薦組み合わせ
type PredictionCombination struct {
	ID          int                   `json:"id"`
	Numbers     []int                 `json:"numbers"`
	TotalScore  float64               `json:"total_score"`
	ScoreDetail PredictionScoreDetail `json:"score_detail"`
}

// PredictionAnalysisInfo は予測生成時の分析メタデータ
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

// PredictionResponse は /api/v1/prediction のレスポンス
type PredictionResponse struct {
	Combinations []PredictionCombination `json:"combinations"`
	AnalysisInfo PredictionAnalysisInfo  `json:"analysis_info"`
	Config       PredictionRequest       `json:"config"`
}
