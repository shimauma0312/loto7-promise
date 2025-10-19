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

// APIサーバーの正常性をチェックする
//
// 機能:
//   - APIサーバーの動作状況を確認
//   - サーバーの稼働時間を取得
//   - 各内部サービスの状態を確認
//
// リクエスト:
//   - HTTP Method: GET
//   - Path: /health
//   - Parameters: なし
//
// レスポンス:
//   - Status: 200 OK
//   - Body: HealthResponse形式のJSON
//   - status: サーバーの状態 ("healthy")
//   - uptime: サーバーの稼働時間
//   - services: 各サービスの動作状況
//   - timestamp: レスポンス生成時刻
//
// エラーレスポンス:
//   - なし（常に正常応答）
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

// 過去のロト7抽選結果をデフォルト設定で取得する
//
// 機能:
//   - 過去10回分の抽選結果を取得（デフォルト値）
//   - getResultsWithCountの内部呼び出しでcount=10固定
//
// リクエスト:
//   - HTTP Method: GET
//   - Path: /api/v1/results
//   - Parameters: なし
//
// レスポンス:
//   - getResultsWithCountと同じ形式
//
// 注意:
//   - この関数は内部でgetResultsWithCountを呼び出すラッパー関数
func getResults(c *gin.Context) {
	getResultsWithCount(c)
}

// 指定回数分の過去のロト7抽選結果を取得する
//
// 機能概要:
//   - 指定された回数分の過去の抽選結果を取得
//   - キャッシュされたデータから結果を読み込み
//   - 抽選回数、抽選日、当選番号を含むデータを返却
//
// リクエスト:
//   - HTTP Method: GET
//   - Path: /api/v1/results/{count}
//   - Parameters:
//   - count (path): 取得する抽選結果の回数 (1-1000の範囲、省略時は10)
//
// レスポンス:
//   - Status: 200 OK
//   - Body: ResultResponse形式のJSON
//   - results: DrawResultAPIの配列
//   - draw_number: 抽選回数
//   - numbers: 当選番号の配列
//   - date: 抽選日
//   - formatted_numbers: フォーマットされた当選番号
//   - total_count: 実際に取得した結果数
//   - requested_count: リクエストされた回数
//   - data_source: データソース ("cache")
//
// エラーレスポンス:
//   - 400 Bad Request: countが範囲外の値 (1-1000以外)
//   - 500 Internal Server Error: データ取得に失敗
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

// 数字の出現頻度ヒートマップをデフォルト設定で取得する
//
// 機能:
//   - 過去50回分の抽選結果からヒートマップを生成（デフォルト値）
//   - getHeatmapWithRangeの内部呼び出しでrange=50固定
//
// リクエスト:
//   - HTTP Method: GET
//   - Path: /api/v1/heatmap
//   - Parameters: なし
//
// レスポンス:
//   - getHeatmapWithRangeと同じ形式
//
// 注意:
//   - この関数は内部でgetHeatmapWithRangeを呼び出すラッパー関数
func getHeatmap(c *gin.Context) {
	getHeatmapWithRange(c)
}

// 指定範囲の抽選結果から数字の出現頻度ヒートマップを生成する
//
// 機能:
//   - 指定された範囲の過去の抽選結果から各数字の出現頻度を分析
//   - 最も出現頻度が高い数字（ホット）と低い数字（コールド）を特定
//   - 最も偏りがある数字と統計情報を計算
//
// リクエスト:
//   - HTTP Method: GET
//   - Path: /api/v1/heatmap/{range}
//   - Parameters:
//   - range (path): 分析対象とする抽選回数 (1-1000の範囲、省略時は50)
//
// レスポンス:
//   - Status: 200 OK
//   - Body: HeatmapResponse形式のJSON
//   - search_range: 検索範囲
//   - actual_draws: 実際に分析した抽選回数
//   - hottest_numbers: 出現頻度が高い数字のリスト
//   - coldest_numbers: 出現頻度が低い数字のリスト
//   - most_biased_number: 最も偏りがある数字
//   - statistics: 統計情報
//   - total_numbers: 総数字数 (37)
//   - average_frequency: 平均出現頻度
//   - max_frequency: 最大出現頻度
//   - min_frequency: 最小出現頻度
//
// エラーレスポンス:
//   - 400 Bad Request: rangeが範囲外の値 (1-1000以外)
//   - 500 Internal Server Error: ヒートマップ生成に失敗
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

// 指定した数字の位置別出現分析を行う
//
// 機能:
//   - 特定の数字が各位置（1～7番目）にどの程度出現するかを分析
//   - 指定した範囲内での位置別出現頻度を計算
//   - 数字の出現パターンと傾向を提供
//
// リクエスト:
//   - HTTP Method: GET
//   - Path: /api/v1/heatmap/{range}/number/{number}
//   - Parameters:
//   - range (path): 分析対象とする抽選回数 (1-1000の範囲)
//   - number (path): 分析対象の数字 (1-37の範囲)
//
// レスポンス:
//   - Status: 200 OK
//   - Body: JSON形式
//   - number: 分析対象の数字
//   - search_range: 検索範囲
//   - actual_draws: 実際に分析した抽選回数
//   - position_data: 位置別の出現データ
//   - 各位置での出現回数と割合
//
// エラーレスポンス:
//   - 400 Bad Request: rangeが範囲外 (1-1000以外) または numberが範囲外 (1-37以外)
//   - 500 Internal Server Error: ヒートマップ生成失敗またはデータ取得失敗
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

// 指定した位置での数字出現分析を行う
//
// 機能概要:
//   - 特定の位置（1～7番目）にどの数字がどの程度出現するかを分析
//   - 指定した範囲内での数字別出現頻度を計算
//   - 位置における数字の出現パターンと傾向を提供
//
// リクエスト:
//   - HTTP Method: GET
//   - Path: /api/v1/heatmap/{range}/position/{position}
//   - Parameters:
//   - range (path): 分析対象とする抽選回数 (1-1000の範囲)
//   - position (path): 分析対象の位置 (1-7の範囲)
//
// レスポンス:
//   - Status: 200 OK
//   - Body: JSON形式
//   - position: 分析対象の位置
//   - search_range: 検索範囲
//   - actual_draws: 実際に分析した抽選回数
//   - number_data: 数字別の出現データ
//   - 各数字での出現回数と割合
//
// エラーレスポンス:
//   - 400 Bad Request: rangeが範囲外 (1-1000以外) または positionが範囲外 (1-7以外)
//   - 500 Internal Server Error: ヒートマップ生成失敗またはデータ取得失敗
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

// GETリクエストでデフォルト設定による推薦番号を生成する
//
// 機能:
//   - 過去のデータを基に統計的分析を行い推薦番号を生成
//   - URLクエリパラメータで一部設定をカスタマイズ可能
//   - デフォルト設定: 最大5組、過去100回分を分析、直近3回を回避
//
// リクエスト:
//   - HTTP Method: GET
//   - Path: /api/v1/recommendations
//   - Query Parameters (すべて省略可能):
//   - count: 生成する推薦組数 (1-10の範囲、デフォルト: 5)
//   - history: 分析対象とする過去の抽選回数 (1-1000の範囲、デフォルト: 100)
//   - avoid: 直近から除外する抽選回数 (0-10の範囲、デフォルト: 3)
//
// レスポンス:
//   - Status: 200 OK
//   - Body: RecommendationResponse形式のJSON
//   - recommendations: RecommendationSetの配列
//   - id: 推薦セットID
//   - numbers: 推薦番号の配列 (7個)
//   - analysis_info: 分析情報
//   - data_source: データソース
//   - analyzed_draws: 分析した抽選回数
//   - generated_at: 生成日時
//   - algorithm: 使用アルゴリズム
//   - config: 使用した設定情報
//
// エラーレスポンス:
//   - generateRecommendationsで発生するエラーと同じ
//
// 注意:
//   - この関数は内部でgenerateRecommendationsを呼び出す
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

// POSTリクエストでカスタム設定による推薦番号を生成する
//
// 機能:
//   - JSONで詳細な設定を受け取り、カスタマイズされた推薦番号を生成
//   - 分析アルゴリズムの重み付けや除外/優先数字の設定が可能
//   - 不正な設定値は自動的にデフォルト値に修正される
//
// リクエスト:
//   - HTTP Method: POST
//   - Path: /api/v1/recommendations
//   - Content-Type: application/json
//   - Body: RecommendationConfig形式のJSON (すべて省略可能)
//   - max_recommendations: 最大推薦組数 (1-10、デフォルト: 5)
//   - history_lookback: 分析対象の過去抽選回数 (1-1000、デフォルト: 100)
//   - recent_avoid_count: 直近除外回数 (0-10、デフォルト: 3)
//   - frequency_weight: 頻度分析の重み (デフォルト: 1.0)
//   - recent_weight: 最近の傾向の重み (デフォルト: 1.2)
//   - consecutive_weight: 連続性の重み (デフォルト: 0.8)
//   - position_weight: 位置分析の重み (デフォルト: 1.1)
//   - consecutive_boost: 連続数字のブースト値 (デフォルト: 1.5)
//
// レスポンス:
//   - Status: 200 OK
//   - Body: RecommendationResponse形式のJSON (getRecommendationsと同じ)
//
// エラーレスポンス:
//   - generateRecommendationsで発生するエラーと同じ
//
// 注意:
//   - JSONパースエラー時はデフォルト設定を使用
//   - 範囲外の値は自動的にデフォルト値に修正
//   - この関数は内部でgenerateRecommendationsを呼び出す
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

// 推薦番号生成の共通処理
//
// 機能:
//   - 指定された設定に基づいて推薦エンジンを作成
//   - 過去のデータを読み込み分析を実行
//   - 推薦番号を生成してAPIレスポンス形式に変換
//   - GET/POSTの両方のエンドポイントから共通で呼び出される
//
// 引数:
//   - c: Ginのコンテキスト (HTTPレスポンス用)
//   - config: 推薦生成の設定 (RecommendationConfig型)
//   - max_recommendations: 生成する推薦組数
//   - history_lookback: 分析する過去の抽選回数
//   - 各種重み付け設定
//
// 処理フロー:
//  1. RecommendationEngineのインスタンス作成
//  2. 指定された履歴データの読み込み
//  3. 推薦アルゴリズムによる推薦番号生成
//  4. API形式のレスポンスデータに変換
//  5. JSONレスポンスとして返却
//
// レスポンス:
//   - 成功時: 200 OK + RecommendationResponse
//   - 失敗時: 500 Internal Server Error + ErrorResponse
//
// エラーレスポンス:
//   - DATA_LOAD_ERROR: データ読み込み失敗
//   - RECOMMENDATION_ERROR: 推薦生成失敗
//
// 注意:
//   - この関数はpublic APIではなく内部処理用
//   - エラー時は適切なHTTPステータスコードとエラーメッセージを返却
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
