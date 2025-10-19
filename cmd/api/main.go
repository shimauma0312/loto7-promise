package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/loto7-promise/internal/api"
	"github.com/shimauma0312/loto7-promise/internal/heatmap"
	"github.com/shimauma0312/loto7-promise/internal/recommendation"
	"github.com/shimauma0312/loto7-promise/internal/result"
)

var (
	startTime = time.Now()
)

func main() {
	// モードを設定（本番環境では GIN_MODE=release を設定）
	gin.SetMode(gin.DebugMode)

	// ルーターを作成
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// ヘルスチェックエンドポイント
	r.GET("/health", healthCheck)

	// APIバージョニング
	v1 := r.Group("/api/v1")
	{
		// 結果取得
		v1.GET("/results", getResults)
		v1.GET("/results/:count", getResultsWithCount)

		// ヒートマップ
		v1.GET("/heatmap", getHeatmap)
		v1.GET("/heatmap/:range", getHeatmapWithRange)
		v1.GET("/heatmap/:range/number/:number", getHeatmapByNumber)
		v1.GET("/heatmap/:range/position/:position", getHeatmapByPosition)

		// 推薦機能
		v1.GET("/recommendations", getRecommendations)
		v1.POST("/recommendations", getRecommendationsPost)
	}

	// ルート
	r.GET("/", func(c *gin.Context) {
		response := api.SuccessResponse(
			"ロト7 Promise API Server",
			map[string]interface{}{
				"name":    "Loto7 Promise API",
				"version": "1.0.0",
				"endpoints": map[string]interface{}{
					"health":          "/health",
					"results":         "/api/v1/results",
					"heatmap":         "/api/v1/heatmap",
					"recommendations": "/api/v1/recommendations",
				},
				"documentation": "https://github.com/shimauma0312/loto7-promise",
			},
		)
		c.JSON(http.StatusOK, response)
	})

	// サーバー起動
	port := ":8080"
	log.Printf("🚀 Loto7 Promise API Server starting on port %s", port)
	log.Printf("📋 Health check: http://localhost%s/health", port)
	log.Printf("📊 API Documentation: http://localhost%s/", port)

	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// healthCheck ヘルスチェック
func healthCheck(c *gin.Context) {
	uptime := time.Since(startTime)

	data := api.HealthResponse{
		Status:    "healthy",
		Uptime:    uptime.String(),
		Timestamp: time.Now(),
		Services: map[string]string{
			"result_cache":   "operational",
			"heatmap_engine": "operational",
			"recommendation": "operational",
		},
	}

	response := api.SuccessResponse("サービスは正常に動作しています", data)
	c.JSON(http.StatusOK, response)
}

// getResults 結果取得（デフォルト10回）
func getResults(c *gin.Context) {
	getResultsWithCount(c)
}

// getResultsWithCount 指定回数の結果取得
func getResultsWithCount(c *gin.Context) {
	// パスパラメータから回数を取得
	countStr := c.Param("count")
	count := 10 // デフォルト値

	if countStr != "" {
		if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 && parsedCount <= 1000 {
			count = parsedCount
		} else {
			response := api.ErrorResponse(
				"INVALID_COUNT",
				"回数は1-1000の範囲で指定してください",
				"count parameter must be between 1 and 1000",
			)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	// 結果を取得
	results := result.GetResult(count)
	if results == nil {
		response := api.ErrorResponse(
			"DATA_FETCH_ERROR",
			"データの取得に失敗しました",
			"Failed to fetch lottery results",
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// レスポンス用にデータを変換
	apiResults := make([]api.DrawResultAPI, len(results))
	for i, res := range results {
		drawNumber := 0
		if len(results) > 0 {
			// 最新の回数から逆算
			latestDrawNum := result.NewNumber()
			drawNumber = latestDrawNum - i
		}

		apiResults[i] = api.DrawResultAPI{
			DrawNumber:       drawNumber,
			Numbers:          res,
			Date:             time.Now().Format("2006-01-02"),
			FormattedNumbers: strings.Join(res, "-"),
		}
	}

	data := api.ResultResponse{
		Results:        apiResults,
		TotalCount:     len(results),
		RequestedCount: count,
		DataSource:     "cache",
	}

	response := api.SuccessResponse("結果を正常に取得しました", data)
	c.JSON(http.StatusOK, response)
}

// getHeatmap ヒートマップ取得（デフォルト50回）
func getHeatmap(c *gin.Context) {
	getHeatmapWithRange(c)
}

// getHeatmapWithRange 指定範囲のヒートマップ取得
func getHeatmapWithRange(c *gin.Context) {
	// パスパラメータから範囲を取得
	rangeStr := c.Param("range")
	searchRange := 50 // デフォルト値

	if rangeStr != "" {
		if parsedRange, err := strconv.Atoi(rangeStr); err == nil && parsedRange > 0 && parsedRange <= 1000 {
			searchRange = parsedRange
		} else {
			response := api.ErrorResponse(
				"INVALID_RANGE",
				"範囲は1-1000の値で指定してください",
				"range parameter must be between 1 and 1000",
			)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	// ヒートマップを生成
	summary, err := heatmap.GenerateHeatmap(searchRange)
	if err != nil {
		response := api.ErrorResponse(
			"HEATMAP_GENERATION_ERROR",
			"ヒートマップの生成に失敗しました",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// 統計情報を計算
	stats := api.HeatmapStatistics{
		TotalNumbers:     37,
		AverageFrequency: float64(summary.ActualDraws*7) / 37.0,
		MaxFrequency:     0,
		MinFrequency:     summary.ActualDraws,
	}

	// 最大・最小頻度を計算
	for number := 1; number <= 37; number++ {
		totalCount := 0
		for position := 1; position <= 7; position++ {
			totalCount += summary.PositionData[number][position].Count
		}
		if totalCount > stats.MaxFrequency {
			stats.MaxFrequency = totalCount
		}
		if totalCount < stats.MinFrequency {
			stats.MinFrequency = totalCount
		}
	}

	data := api.HeatmapResponse{
		SearchRange:      summary.SearchRange,
		ActualDraws:      summary.ActualDraws,
		HottestNumbers:   summary.HottestNumbers,
		ColdestNumbers:   summary.ColdestNumbers,
		MostBiasedNumber: summary.MostBiasedNumber,
		Statistics:       stats,
	}

	response := api.SuccessResponse("ヒートマップを正常に生成しました", data)
	c.JSON(http.StatusOK, response)
}

// getHeatmapByNumber 指定数字のヒートマップ取得
func getHeatmapByNumber(c *gin.Context) {
	rangeStr := c.Param("range")
	numberStr := c.Param("number")

	searchRange, err := strconv.Atoi(rangeStr)
	if err != nil || searchRange <= 0 || searchRange > 1000 {
		response := api.ErrorResponse(
			"INVALID_RANGE",
			"範囲は1-1000の値で指定してください",
			"range parameter must be between 1 and 1000",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	number, err := strconv.Atoi(numberStr)
	if err != nil || number < 1 || number > 37 {
		response := api.ErrorResponse(
			"INVALID_NUMBER",
			"数字は1-37の値で指定してください",
			"number parameter must be between 1 and 37",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	summary, err := heatmap.GenerateHeatmap(searchRange)
	if err != nil {
		response := api.ErrorResponse(
			"HEATMAP_GENERATION_ERROR",
			"ヒートマップの生成に失敗しました",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	positionData, err := summary.GetPositionData(number)
	if err != nil {
		response := api.ErrorResponse(
			"DATA_RETRIEVAL_ERROR",
			"位置データの取得に失敗しました",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	data := map[string]interface{}{
		"number":        number,
		"search_range":  searchRange,
		"actual_draws":  summary.ActualDraws,
		"position_data": positionData,
	}

	response := api.SuccessResponse("指定数字のヒートマップデータを取得しました", data)
	c.JSON(http.StatusOK, response)
}

// getHeatmapByPosition 指定位置のヒートマップ取得
func getHeatmapByPosition(c *gin.Context) {
	rangeStr := c.Param("range")
	positionStr := c.Param("position")

	searchRange, err := strconv.Atoi(rangeStr)
	if err != nil || searchRange <= 0 || searchRange > 1000 {
		response := api.ErrorResponse(
			"INVALID_RANGE",
			"範囲は1-1000の値で指定してください",
			"range parameter must be between 1 and 1000",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	position, err := strconv.Atoi(positionStr)
	if err != nil || position < 1 || position > 7 {
		response := api.ErrorResponse(
			"INVALID_POSITION",
			"位置は1-7の値で指定してください",
			"position parameter must be between 1 and 7",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	summary, err := heatmap.GenerateHeatmap(searchRange)
	if err != nil {
		response := api.ErrorResponse(
			"HEATMAP_GENERATION_ERROR",
			"ヒートマップの生成に失敗しました",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	numberData, err := summary.GetNumbersByPosition(position)
	if err != nil {
		response := api.ErrorResponse(
			"DATA_RETRIEVAL_ERROR",
			"数字データの取得に失敗しました",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	data := map[string]interface{}{
		"position":     position,
		"search_range": searchRange,
		"actual_draws": summary.ActualDraws,
		"number_data":  numberData,
	}

	response := api.SuccessResponse("指定位置のヒートマップデータを取得しました", data)
	c.JSON(http.StatusOK, response)
}

// getRecommendations 推薦取得（GETリクエスト、デフォルト設定）
func getRecommendations(c *gin.Context) {
	// クエリパラメータから設定を取得
	config := recommendation.DefaultConfig()

	if countStr := c.Query("count"); countStr != "" {
		if count, err := strconv.Atoi(countStr); err == nil && count > 0 && count <= 10 {
			config.MaxRecommendations = count
		}
	}

	if historyStr := c.Query("history"); historyStr != "" {
		if history, err := strconv.Atoi(historyStr); err == nil && history > 0 && history <= 1000 {
			config.HistoryLookback = history
		}
	}

	if avoidStr := c.Query("avoid"); avoidStr != "" {
		if avoid, err := strconv.Atoi(avoidStr); err == nil && avoid >= 0 && avoid <= 10 {
			config.RecentAvoidCount = avoid
		}
	}

	generateRecommendations(c, config)
}

// getRecommendationsPost 推薦取得（POSTリクエスト、カスタム設定）
func getRecommendationsPost(c *gin.Context) {
	var reqConfig recommendation.RecommendationConfig

	// JSONから設定を読み込み
	if err := c.ShouldBindJSON(&reqConfig); err != nil {
		// デフォルト設定を使用
		reqConfig = recommendation.DefaultConfig()
	} else {
		// 値の妥当性チェック
		if reqConfig.MaxRecommendations <= 0 || reqConfig.MaxRecommendations > 10 {
			reqConfig.MaxRecommendations = 5
		}
		if reqConfig.HistoryLookback <= 0 || reqConfig.HistoryLookback > 1000 {
			reqConfig.HistoryLookback = 100
		}
		if reqConfig.RecentAvoidCount < 0 || reqConfig.RecentAvoidCount > 10 {
			reqConfig.RecentAvoidCount = 3
		}
	}

	generateRecommendations(c, reqConfig)
}

// generateRecommendations 推薦生成の共通処理
func generateRecommendations(c *gin.Context, config recommendation.RecommendationConfig) {
	// 推薦エンジンを作成
	engine := recommendation.NewRecommendationEngine(config)

	// データを読み込み
	err := engine.LoadData(config.HistoryLookback)
	if err != nil {
		response := api.ErrorResponse(
			"DATA_LOAD_ERROR",
			"推薦データの読み込みに失敗しました",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// 推薦を生成
	recommendations, err := engine.GenerateRecommendations()
	if err != nil {
		response := api.ErrorResponse(
			"RECOMMENDATION_ERROR",
			"推薦の生成に失敗しました",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// レスポンス用にデータを変換
	recSets := make([]api.RecommendationSet, len(recommendations))
	for i, rec := range recommendations {
		recSets[i] = api.RecommendationSet{
			ID:      i + 1,
			Numbers: rec,
		}
	}

	data := api.RecommendationResponse{
		Recommendations: recSets,
		AnalysisInfo: api.AnalysisInfo{
			DataSource:    "cache",
			AnalyzedDraws: engine.GetAnalyzedDrawsCount(),
			GeneratedAt:   time.Now(),
			Algorithm:     "human-intuition-based",
		},
		Config: api.ConfigInfo{
			MaxRecommendations: config.MaxRecommendations,
			RecentAvoidCount:   config.RecentAvoidCount,
			ConsecutiveBoost:   config.ConsecutiveBoost,
			HistoryLookback:    config.HistoryLookback,
			FrequencyWeight:    config.FrequencyWeight,
			RecentWeight:       config.RecentWeight,
			ConsecutiveWeight:  config.ConsecutiveWeight,
			PositionWeight:     config.PositionWeight,
		},
	}

	response := api.SuccessResponse("推薦を正常に生成しました", data)
	c.JSON(http.StatusOK, response)
}
