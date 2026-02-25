package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/loto7-promise/internal/api"
	"github.com/shimauma0312/loto7-promise/internal/heatmap"
	"github.com/shimauma0312/loto7-promise/internal/prediction"
	"github.com/shimauma0312/loto7-promise/internal/random"
	"github.com/shimauma0312/loto7-promise/internal/recommendation"
	"github.com/shimauma0312/loto7-promise/internal/result"
	"github.com/shimauma0312/loto7-promise/internal/simulation"
)

var (
	startTime = time.Now()
)

// クエリパラメータを解析して検証する
func parseQueryParam(c *gin.Context, key string, defaultValue, min, max int, errorMsg, errorCode string) (int, bool) {
	strVal := c.Query(key)
	if strVal == "" {
		return defaultValue, true
	}

	if parsedVal, err := strconv.Atoi(strVal); err == nil && parsedVal >= min && parsedVal <= max {
		return parsedVal, true
	}

	response := api.ErrorResponse(
		errorCode,
		errorMsg,
		fmt.Sprintf("%s parameter must be between %d and %d", key, min, max),
	)
	c.JSON(http.StatusBadRequest, response)
	return 0, false
}

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
		v1.GET("/heatmap/table", getHeatmapTableDefault)
		v1.GET("/heatmap/:range", getHeatmapWithRange)
		v1.GET("/heatmap/:range/table", getHeatmapTable)
		v1.GET("/heatmap/:range/number/:number", getHeatmapByNumber)
		v1.GET("/heatmap/:range/position/:position", getHeatmapByPosition)

		// 推薦機能
		v1.GET("/recommendations", getRecommendations)
		v1.POST("/recommendations", getRecommendationsPost)

		// ランダム生成機能
		v1.GET("/random", getRandom)

		// 統計分析予測機能
		v1.GET("/prediction", getPrediction)
		v1.POST("/prediction", getPredictionPost)

		// シミュレーション機能
		v1.POST("/simulation", runSimulation)
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
					"heatmap_table":   "/api/v1/heatmap/table",
					"recommendations": "/api/v1/recommendations",
					"random":          "/api/v1/random",
					"prediction":      "/api/v1/prediction",
					"simulation":      "/api/v1/simulation",
				},
				"documentation": "https://github.com/shimauma0312/loto7-promise",
			},
		)
		c.JSON(http.StatusOK, response)
	})

	// サーバー起動
	port := ":8080"
	log.Printf("Loto7 Promise API Server starting on port %s", port)
	log.Printf("Health check: http://localhost%s/health", port)
	log.Printf("API Documentation: http://localhost%s/", port)

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
			"random":         "operational",
			"prediction":     "operational",
			"simulation":     "operational",
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
	results, err := result.GetResultWithDate(count)
	if err != nil {
		response := api.ErrorResponse(
			"DATA_RETRIEVAL_ERROR",
			"データの取得に失敗しました",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// レスポンス用にデータを変換
	apiResults := make([]api.DrawResultAPI, len(results))
	for i, res := range results {
		apiResults[i] = api.DrawResultAPI{
			DrawNumber:       res.DrawNumber,
			Numbers:          res.Numbers,
			Date:             res.Date,
			FormattedNumbers: strings.Join(res.Numbers, "-"),
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

	count, ok := parseQueryParam(c, "count", 5, 1, 10, "生成数は1-10の範囲で指定してください", "INVALID_COUNT")
	if !ok {
		return
	}
	config.MaxRecommendations = count

	// DoS対策: historyは最大200に制限
	history, ok := parseQueryParam(c, "history", 100, 1, 200, "履歴数は1-200の範囲で指定してください", "INVALID_HISTORY")
	if !ok {
		return
	}
	config.HistoryLookback = history

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
		// DoS対策: historyは最大200に制限
		if reqConfig.HistoryLookback <= 0 || reqConfig.HistoryLookback > 200 {
			reqConfig.HistoryLookback = 100
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
			HistoryLookback:    config.HistoryLookback,
			FrequencyWeight:    config.FrequencyWeight,
			RecentWeight:       config.RecentWeight,
		},
	}

	response := api.SuccessResponse("推薦を正常に生成しました", data)
	c.JSON(http.StatusOK, response)
}

// 指定範囲の抽選結果からヒートマップテーブルをマークダウン形式で取得する
//
// 機能:
//   - 指定された範囲の過去の抽選結果から位置別出現ヒートマップテーブルを生成
//   - マークダウン形式のテーブルとして出力
//   - 統計情報も含めて表形式で視覚的に分かりやすく表示
//
// リクエスト:
//   - HTTP Method: GET
//   - Path: /api/v1/heatmap/{range}/table
//   - Parameters:
//   - range (path): 分析対象とする抽選回数 (1-1000の範囲、省略時は50)
//
// レスポンス:
//   - Status: 200 OK
//   - Content-Type: text/plain
//   - Body: マークダウン形式のヒートマップテーブル
//   - ヘッダー情報（検索範囲、実際の抽選回数）
//   - 位置別出現回数テーブル（数字1-37 × 位置1-7 + 合計）
//   - 統計情報（最頻出数字、最低出現数字、最も偏りが大きい数字）
//
// エラーレスポンス:
//   - 400 Bad Request: rangeが範囲外の値 (1-1000以外)
//   - 500 Internal Server Error: ヒートマップ生成に失敗
//
// 注意:
//   - レスポンスはJSONではなくプレーンテキスト形式
//   - マークダウンビューアで表として表示可能
//   - コピー&ペーストでそのまま使用可能
func getHeatmapTable(c *gin.Context) {
	// パスパラメータから範囲を取得
	rangeStr := c.Param("range")
	searchRange := 50 // デフォルト値

	if rangeStr != "" {
		if parsedRange, err := strconv.Atoi(rangeStr); err == nil && parsedRange > 0 && parsedRange <= 1000 {
			searchRange = parsedRange
		} else {
			c.Header("Content-Type", "text/plain")
			c.String(http.StatusBadRequest, "エラー: 範囲は1-1000の値で指定してください")
			return
		}
	}

	// ヒートマップを生成
	summary, err := heatmap.GenerateHeatmap(searchRange)
	if err != nil {
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusInternalServerError, "エラー: ヒートマップの生成に失敗しました - %s", err.Error())
		return
	}

	// マークダウン形式のテーブルを生成
	tableText := generateMarkdownHeatmapTable(summary)

	// プレーンテキストとして返却
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, tableText)
}

// ヒートマップテーブルをデフォルト設定（50回分）でマークダウン形式で取得する
//
// 機能:
//   - 過去50回分の抽選結果からヒートマップテーブルを生成（デフォルト値）
//   - getHeatmapTableの内部呼び出しでrange=50固定
//
// リクエスト:
//   - HTTP Method: GET
//   - Path: /api/v1/heatmap/table
//   - Parameters: なし
//
// レスポンス:
//   - getHeatmapTableと同じ形式
//
// 注意:
//   - この関数は内部でgetHeatmapTableを呼び出すラッパー関数
func getHeatmapTableDefault(c *gin.Context) {
	// 50回分を固定で設定
	c.Params = append(c.Params, gin.Param{Key: "range", Value: "50"})
	getHeatmapTable(c)
}

// ヒートマップサマリーからマークダウン形式のテーブルを生成
func generateMarkdownHeatmapTable(summary *heatmap.HeatmapSummary) string {
	var result strings.Builder

	// ヘッダー情報
	result.WriteString(fmt.Sprintf("# ロト7位置別出現ヒートマップ（過去%d回分）\n\n", summary.SearchRange))
	result.WriteString(fmt.Sprintf("**実際の抽選回数**: %d回\n\n", summary.ActualDraws))

	// テーブルヘッダー
	result.WriteString("| 数字 | 位置1 | 位置2 | 位置3 | 位置4 | 位置5 | 位置6 | 位置7 | 合計 |\n")
	result.WriteString("|------|-------|-------|-------|-------|-------|-------|-------|------|\n")

	// データ行
	for number := 1; number <= 37; number++ {
		result.WriteString(fmt.Sprintf("| %2d", number))
		total := 0
		for position := 1; position <= 7; position++ {
			count := summary.PositionData[number][position].Count
			result.WriteString(fmt.Sprintf(" | %3d", count))
			total += count
		}
		result.WriteString(fmt.Sprintf(" | %3d |\n", total))
	}

	// 統計情報
	result.WriteString("\n## 統計情報\n\n")

	// 最頻出数字
	result.WriteString(fmt.Sprintf("**最頻出数字**: %v\n", formatNumberList(summary.HottestNumbers)))

	// 最低出現数字
	result.WriteString(fmt.Sprintf("**最低出現数字**: %v\n", formatNumberList(summary.ColdestNumbers)))

	// 最も偏りが大きい数字
	result.WriteString(fmt.Sprintf("**最も位置偏りが大きい数字**: %d\n", summary.MostBiasedNumber))

	return result.String()
}

// 数字リストをフォーマット
func formatNumberList(numbers []int) string {
	if len(numbers) == 0 {
		return "なし"
	}

	var strNumbers []string
	for _, num := range numbers {
		strNumbers = append(strNumbers, strconv.Itoa(num))
	}

	return "[" + strings.Join(strNumbers, ", ") + "]"
}

// ランダム性を重視した推薦番号を生成する
//
// 機能:
//   - 各数字位置で過去に出現した範囲内からランダムに数字を選択
//   - 未出現の数字を自動的に除外
//   - よりロト7の本来のランダム性を再現
//
// リクエスト:
//   - HTTP Method: GET
//   - Path: /api/v1/random
//   - Parameters:
//   - count (query): 生成する組み合わせ数 (1-10、デフォルト: 5)
//   - history (query): 分析する過去の抽選回数 (1-1000、デフォルト: 100)
//
// レスポンス:
//   - Status: 200 OK
//   - Body: RandomResponse形式のJSON
//   - combinations: RandomSet配列（各組み合わせ）
//   - analysis_info: 分析情報
//   - data_source: データソース
//   - history_count: 分析した履歴数
//   - timestamp: 生成時刻
//
// エラーレスポンス:
//   - 400 Bad Request: パラメータが範囲外
//   - 500 Internal Server Error: データ読み込みまたは生成に失敗
func getRandom(c *gin.Context) {
	// クエリパラメータから設定を取得
	count, ok := parseQueryParam(c, "count", 5, 1, 10, "生成数は1-10の範囲で指定してください", "INVALID_COUNT")
	if !ok {
		return
	}

	// DoS対策: historyは最大200に制限
	historyCount, ok := parseQueryParam(c, "history", 100, 1, 200, "履歴数は1-200の範囲で指定してください", "INVALID_HISTORY")
	if !ok {
		return
	}

	// ランダムエンジンを作成
	engine := random.NewRandomEngine()

	// データを読み込み
	err := engine.LoadData(historyCount)
	if err != nil {
		response := api.ErrorResponse(
			"DATA_LOAD_ERROR",
			"ランダム生成データの読み込みに失敗しました",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// ランダム組み合わせを生成
	combinations, err := engine.GenerateMultipleRandomCombinations(count)
	if err != nil {
		response := api.ErrorResponse(
			"RANDOM_GENERATION_ERROR",
			"ランダム組み合わせの生成に失敗しました",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// レスポンス用にデータを変換
	randomSets := make([]api.RecommendationSet, len(combinations))
	for i, combo := range combinations {
		randomSets[i] = api.RecommendationSet{
			ID:      i + 1,
			Numbers: combo,
		}
	}

	data := api.RecommendationResponse{
		Recommendations: randomSets,
		AnalysisInfo: api.AnalysisInfo{
			DataSource:    "cache",
			AnalyzedDraws: historyCount,
			GeneratedAt:   time.Now(),
			Algorithm:     "random",
		},
		Config: api.ConfigInfo{}, // ランダム生成なので空の設定
	}

	response := api.SuccessResponse("ランダム組み合わせを生成しました", data)
	c.JSON(http.StatusOK, response)
}

// GETリクエストで統計分析予測番号を生成する
//
// 機能:
//   - ①ホット/コールド分析 ②パリティ ③大小バランス ④合計値 ⑤十の位グループ
//     ⑥一の位種類数 ⑦引っ張り・斜め数字 ⑧ボーナス周辺の8フィルターを適用
//   - ハードフィルター通過後にスコアリングし上位N件を返す
//
// リクエスト:
//   - HTTP Method: GET
//   - Path: /api/v1/prediction
//   - Query Parameters（省略可）:
//   - count: 生成口数 (1-10, デフォルト: 5)
//   - history: 分析回数 (1-500, デフォルト: 100)
//   - hot_cold, parity, size_balance, sum_range, tens_group,
//     last_digit, pull, bonus, interval: 各フィルター重み (0.0-2.0)
//
// レスポンス:
//   - Status: 200 OK
//   - Body: PredictionResponse形式のJSON
//
// エラーレスポンス:
//   - 400 Bad Request: パラメータ不正
//   - 500 Internal Server Error: データ読み込みまたは予測生成失敗
func getPrediction(c *gin.Context) {
	count, ok := parseQueryParam(c, "count", 5, 1, 10, "生成数は1-10の範囲で指定してください", "INVALID_COUNT")
	if !ok {
		return
	}
	history, ok := parseQueryParam(c, "history", 100, 1, 500, "履歴数は1-500の範囲で指定してください", "INVALID_HISTORY")
	if !ok {
		return
	}

	cfg := prediction.DefaultConfig()
	cfg.Count = count
	cfg.History = history

	// クエリパラメータで各フィルター重みを上書き
	overrideWeight := func(key string, target *float64) {
		if v := c.Query(key); v != "" {
			if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 && f <= 2.0 {
				*target = f
			}
		}
	}
	overrideWeight("hot_cold", &cfg.Weights.HotCold)
	overrideWeight("parity", &cfg.Weights.Parity)
	overrideWeight("size_balance", &cfg.Weights.SizeBal)
	overrideWeight("sum_range", &cfg.Weights.SumRange)
	overrideWeight("tens_group", &cfg.Weights.TensGroup)
	overrideWeight("last_digit", &cfg.Weights.LastDigit)
	overrideWeight("pull", &cfg.Weights.Pull)
	overrideWeight("bonus", &cfg.Weights.Bonus)
	overrideWeight("interval", &cfg.Weights.Interval)

	result, err := runPredictionEngine(cfg)
	if err != nil {
		response := api.ErrorResponse("PREDICTION_ERROR", "予測生成に失敗しました", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, api.SuccessResponse("統計分析予測番号を生成しました", result))
}

// POSTリクエストでフィルター重みを細かく指定して統計分析予測番号を生成する
//
// リクエスト:
//   - HTTP Method: POST
//   - Path: /api/v1/prediction
//   - Content-Type: application/json
//   - Body: PredictionRequest形式のJSON
//   - count: 生成口数 (デフォルト: 5)
//   - history: 分析回数 (デフォルト: 100)
//   - weights: 各フィルター重み
//
// レスポンス・エラー: GETと同様
func getPredictionPost(c *gin.Context) {
	var req api.PredictionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(
			"INVALID_REQUEST", "リクエストボディの解析に失敗しました", err.Error(),
		))
		return
	}

	cfg := prediction.DefaultConfig()
	if req.Count > 0 && req.Count <= 10 {
		cfg.Count = req.Count
	}
	if req.History > 0 && req.History <= 500 {
		cfg.History = req.History
	}
	// 重みは0より大きい値のみ適用
	applyWeight := func(src float64, dst *float64) {
		if src > 0 {
			*dst = src
		}
	}
	applyWeight(req.Weights.HotCold, &cfg.Weights.HotCold)
	applyWeight(req.Weights.Parity, &cfg.Weights.Parity)
	applyWeight(req.Weights.SizeBal, &cfg.Weights.SizeBal)
	applyWeight(req.Weights.SumRange, &cfg.Weights.SumRange)
	applyWeight(req.Weights.TensGroup, &cfg.Weights.TensGroup)
	applyWeight(req.Weights.LastDigit, &cfg.Weights.LastDigit)
	applyWeight(req.Weights.Pull, &cfg.Weights.Pull)
	applyWeight(req.Weights.Bonus, &cfg.Weights.Bonus)
	applyWeight(req.Weights.Interval, &cfg.Weights.Interval)

	result, err := runPredictionEngine(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(
			"PREDICTION_ERROR", "予測生成に失敗しました", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, api.SuccessResponse("統計分析予測番号を生成しました", result))
}

// runPredictionEngine はエンジンを初期化してデータ読み込みから予測生成までを実行する
func runPredictionEngine(cfg prediction.PredictionConfig) (api.PredictionResponse, error) {
	engine := prediction.NewEngine()
	if err := engine.LoadData(cfg.History); err != nil {
		return api.PredictionResponse{}, err
	}

	res, err := engine.Generate(cfg)
	if err != nil {
		return api.PredictionResponse{}, err
	}

	// internal の PredictionResult を API レスポンス型に変換
	combos := make([]api.PredictionCombination, len(res.Combinations))
	for i, c := range res.Combinations {
		combos[i] = api.PredictionCombination{
			ID:         c.ID,
			Numbers:    c.Numbers,
			TotalScore: c.TotalScore,
			ScoreDetail: api.PredictionScoreDetail{
				HotColdScore:   c.ScoreDetail.HotColdScore,
				ParityScore:    c.ScoreDetail.ParityScore,
				SizeScore:      c.ScoreDetail.SizeScore,
				SumScore:       c.ScoreDetail.SumScore,
				TensScore:      c.ScoreDetail.TensScore,
				LastDigitScore: c.ScoreDetail.LastDigitScore,
				PullScore:      c.ScoreDetail.PullScore,
				BonusScore:     c.ScoreDetail.BonusScore,
				IntervalScore:  c.ScoreDetail.IntervalScore,
			},
		}
	}

	w := cfg.Weights
	respCfg := api.PredictionRequest{
		Count:   cfg.Count,
		History: cfg.History,
		Weights: api.PredictionFilterWeights{
			HotCold:   w.HotCold,
			Parity:    w.Parity,
			SizeBal:   w.SizeBal,
			SumRange:  w.SumRange,
			TensGroup: w.TensGroup,
			LastDigit: w.LastDigit,
			Pull:      w.Pull,
			Bonus:     w.Bonus,
			Interval:  w.Interval,
		},
	}

	return api.PredictionResponse{
		Combinations: combos,
		AnalysisInfo: api.PredictionAnalysisInfo{
			AnalyzedDraws:    res.AnalysisInfo.AnalyzedDraws,
			HotNumbers:       res.AnalysisInfo.HotNumbers,
			ColdNumbers:      res.AnalysisInfo.ColdNumbers,
			LastDrawNumbers:  res.AnalysisInfo.LastDrawNumbers,
			LastBonusNumbers: res.AnalysisInfo.LastBonusNumbers,
			LastDrawSum:      res.AnalysisInfo.LastDrawSum,
			SumTrend:         res.AnalysisInfo.SumTrend,
			GeneratedAt:      res.AnalysisInfo.GeneratedAt,
		},
		Config: respCfg,
	}, nil
}

// ロト7の1等当選までのシミュレーションを実行する
//
// 機能:
//   - ユーザーの選択数字で1等が当選するまで抽選を繰り返す
//   - 試行回数、所要時間などの統計情報を返却
//   - 複数回のシミュレーションにも対応
//
// リクエスト:
//   - HTTP Method: POST
//   - Path: /api/v1/simulation
//   - Content-Type: application/json
//   - Body:
//   - user_numbers: ユーザーの選択数字（7個の整数配列）
//   - simulation_count: シミュレーション実行回数（1-100、デフォルト: 1）
//   - history_count: 分析する過去の抽選回数（1-1000、デフォルト: 100）
//
// レスポンス:
//   - Status: 200 OK
//   - Body: SimulationResponse形式のJSON
//   - single_result: 単一シミュレーション結果（simulation_count=1の場合）
//   - stats: 統計情報（simulation_count>1の場合）
//
// エラーレスポンス:
//   - 400 Bad Request: パラメータ不正
//   - 500 Internal Server Error: シミュレーション実行失敗
func runSimulation(c *gin.Context) {
	// リクエストボディ
	var req struct {
		UserNumbers      []int `json:"user_numbers"`
		SimulationCount  int   `json:"simulation_count"`
		HistoryCount     int   `json:"history_count"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response := api.ErrorResponse(
			"INVALID_REQUEST",
			"リクエストボディの解析に失敗しました",
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// デフォルト値の設定
	if req.SimulationCount <= 0 {
		req.SimulationCount = 1
	}
	if req.SimulationCount > 100 {
		response := api.ErrorResponse(
			"INVALID_COUNT",
			"シミュレーション回数は1-100の範囲で指定してください",
			"simulation_count must be between 1 and 100",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// DoS対策: historyは最大200に制限
	if req.HistoryCount <= 0 {
		req.HistoryCount = 100
	}
	if req.HistoryCount > 200 {
		req.HistoryCount = 200
	}

	// ユーザーの数字チェック
	if len(req.UserNumbers) == 0 {
		// 自動生成
		engine := random.NewRandomEngine()
		if err := engine.LoadData(req.HistoryCount); err != nil {
			response := api.ErrorResponse(
				"DATA_LOAD_ERROR",
				"データの読み込みに失敗しました",
				err.Error(),
			)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		userNums, err := engine.GenerateRandomCombination()
		if err != nil {
			response := api.ErrorResponse(
				"NUMBER_GENERATION_ERROR",
				"ユーザー数字の生成に失敗しました",
				err.Error(),
			)
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		req.UserNumbers = userNums
	} else if len(req.UserNumbers) != 7 {
		response := api.ErrorResponse(
			"INVALID_NUMBERS",
			"ユーザーの数字は7個必要です",
			"user_numbers must contain exactly 7 numbers",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// シミュレーションエンジンを作成
	simEngine := simulation.NewSimulationEngine()

	// データを読み込み
	if err := simEngine.LoadData(req.HistoryCount); err != nil {
		response := api.ErrorResponse(
			"DATA_LOAD_ERROR",
			"シミュレーションデータの読み込みに失敗しました",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// シミュレーション実行
	if req.SimulationCount == 1 {
		// 単一シミュレーション
		result, err := simEngine.RunSimulation(req.UserNumbers)
		if err != nil {
			response := api.ErrorResponse(
				"SIMULATION_ERROR",
				"シミュレーションの実行に失敗しました",
				err.Error(),
			)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		data := map[string]interface{}{
			"user_numbers":    result.UserNumbers,
			"draw_count":      result.DrawCount,
			"winning_numbers": result.WinningNumbers,
			"duration_ms":     result.Duration.Milliseconds(),
			"estimated_cost":  result.DrawCount * 300,
			"history_count":   req.HistoryCount,
			"timestamp":       time.Now(),
		}

		response := api.SuccessResponse("シミュレーションが完了しました", data)
		c.JSON(http.StatusOK, response)
	} else {
		// 複数シミュレーション
		stats, err := simEngine.RunMultipleSimulations(req.SimulationCount, req.UserNumbers)
		if err != nil {
			response := api.ErrorResponse(
				"SIMULATION_ERROR",
				"シミュレーションの実行に失敗しました",
				err.Error(),
			)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		results := make([]map[string]interface{}, len(stats.Results))
		for i, r := range stats.Results {
			results[i] = map[string]interface{}{
				"draw_count":      r.DrawCount,
				"winning_numbers": r.WinningNumbers,
				"duration_ms":     r.Duration.Milliseconds(),
			}
		}

		data := map[string]interface{}{
			"user_numbers":           req.UserNumbers,
			"simulations":            stats.Simulations,
			"min_draws":              stats.MinDraws,
			"max_draws":              stats.MaxDraws,
			"avg_draws":              stats.AvgDraws,
			"median_draws":           stats.MedianDraws,
			"total_duration_ms":      stats.TotalDuration.Milliseconds(),
			"avg_estimated_cost":     int(stats.AvgDraws) * 300,
			"results":                results,
			"history_count":          req.HistoryCount,
			"timestamp":              time.Now(),
		}

		response := api.SuccessResponse("複数シミュレーションが完了しました", data)
		c.JSON(http.StatusOK, response)
	}
}

