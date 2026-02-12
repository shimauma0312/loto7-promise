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
	DataSource    string    `json:"data_source"`
	AnalyzedDraws int       `json:"analyzed_draws"`
	GeneratedAt   time.Time `json:"generated_at"`
	Algorithm     string    `json:"algorithm"`
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
