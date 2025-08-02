package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/loto7-promise/internal/randomPrediction"
)

// APIレスポンス用の共通構造体
// エラーハンドリングとメッセージ表示を統一するために使用
type APIResponse struct {
	Success bool        `json:"success"`         // 処理成功フラグ
	Message string      `json:"message"`         // レスポンスメッセージ
	Data    interface{} `json:"data,omitempty"`  // レスポンスデータ（エラー時は空）
	Error   string      `json:"error,omitempty"` // エラーメッセージ（成功時は空）
}

// 単一予想のリクエストパラメータ
type PredictionRequest struct {
	AnalysisRange int  `json:"analysis_range" form:"analysis_range"` // 分析対象の回数（デフォルト：30）
	UseTrend      bool `json:"use_trend" form:"use_trend"`           // 傾向を使用するか（デフォルト：true）
}

// 複数予想のリクエストパラメータ
type MultiplePredictionRequest struct {
	Count         int  `json:"count" form:"count"`                   // 生成する予想数（デフォルト：5）
	AnalysisRange int  `json:"analysis_range" form:"analysis_range"` // 分析対象の回数（デフォルト：30）
	UseTrend      bool `json:"use_trend" form:"use_trend"`           // 傾向を使用するか（デフォルト：true）
}

func main() {
	// Ginのデバッグモードを本番環境では無効化
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// ルーター初期化
	router := gin.Default()

	// CORS設定（必要に応じて）
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// ヘルスチェック用エンドポイント
	// nginxのアップストリームヘルスチェックで使用
	router.GET("/health", healthCheck)

	// APIバージョン1のグループ
	v1 := router.Group("/api/v1")
	{
		// 単一予想生成API
		v1.GET("/prediction", getSinglePrediction)
		v1.POST("/prediction", postSinglePrediction)

		// 複数予想生成API
		v1.GET("/predictions", getMultiplePredictions)
		v1.POST("/predictions", postMultiplePredictions)

		// API仕様確認用
		v1.GET("/info", getAPIInfo)
	}

	// ポート番号の設定（環境変数またはデフォルト）
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("ロト7予想APIサーバーを開始します... ポート: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// ヘルスチェックハンドラー
// nginx等のロードバランサーからの生存確認に使用
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "ロト7予想APIサーバーは正常に稼働中です",
	})
}

// GET方式での単一予想取得
// クエリパラメータでリクエストを受け取る
func getSinglePrediction(c *gin.Context) {
	// デフォルト値設定
	analysisRange := 30
	useTrend := true

	// クエリパラメータの解析
	if rangeStr := c.Query("analysis_range"); rangeStr != "" {
		if parsed, err := strconv.Atoi(rangeStr); err == nil && parsed > 0 {
			analysisRange = parsed
		}
	}
	
	if trendStr := c.Query("use_trend"); trendStr != "" {
		if parsed, err := strconv.ParseBool(trendStr); err == nil {
			useTrend = parsed
		}
	}

	// 予想生成実行
	prediction, err := randomPrediction.GenerateRandomPrediction(analysisRange, useTrend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "予想番号の生成に失敗しました",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "予想番号を正常に生成しました",
		Data:    prediction,
	})
}

// POST方式での単一予想取得
// JSONボディでリクエストを受け取る
func postSinglePrediction(c *gin.Context) {
	var req PredictionRequest
	
	// デフォルト値設定
	req.AnalysisRange = 30
	req.UseTrend = true

	// JSONバインド（エラーは無視してデフォルト値使用）
	c.ShouldBindJSON(&req)

	// 入力値検証
	if req.AnalysisRange <= 0 {
		req.AnalysisRange = 30
	}

	// 予想生成実行
	prediction, err := randomPrediction.GenerateRandomPrediction(req.AnalysisRange, req.UseTrend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "予想番号の生成に失敗しました",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "予想番号を正常に生成しました",
		Data:    prediction,
	})
}

// GET方式での複数予想取得
func getMultiplePredictions(c *gin.Context) {
	// デフォルト値設定
	count := 5
	analysisRange := 30
	useTrend := true

	// クエリパラメータの解析
	if countStr := c.Query("count"); countStr != "" {
		if parsed, err := strconv.Atoi(countStr); err == nil && parsed > 0 && parsed <= 20 {
			count = parsed
		}
	}
	
	if rangeStr := c.Query("analysis_range"); rangeStr != "" {
		if parsed, err := strconv.Atoi(rangeStr); err == nil && parsed > 0 {
			analysisRange = parsed
		}
	}
	
	if trendStr := c.Query("use_trend"); trendStr != "" {
		if parsed, err := strconv.ParseBool(trendStr); err == nil {
			useTrend = parsed
		}
	}

	// 複数予想生成実行
	predictions, err := randomPrediction.GenerateMultiplePredictions(count, analysisRange, useTrend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "複数予想番号の生成に失敗しました",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "複数予想番号を正常に生成しました",
		Data:    predictions,
	})
}

// POST方式での複数予想取得
func postMultiplePredictions(c *gin.Context) {
	var req MultiplePredictionRequest
	
	// デフォルト値設定
	req.Count = 5
	req.AnalysisRange = 30
	req.UseTrend = true

	// JSONバインド（エラーは無視してデフォルト値使用）
	c.ShouldBindJSON(&req)

	// 入力値検証
	if req.Count <= 0 || req.Count > 20 {
		req.Count = 5 // 上限を20に制限
	}
	if req.AnalysisRange <= 0 {
		req.AnalysisRange = 30
	}

	// 複数予想生成実行
	predictions, err := randomPrediction.GenerateMultiplePredictions(req.Count, req.AnalysisRange, req.UseTrend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "複数予想番号の生成に失敗しました",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "複数予想番号を正常に生成しました",
		Data:    predictions,
	})
}

// API仕様情報の取得
func getAPIInfo(c *gin.Context) {
	info := map[string]interface{}{
		"name":        "ロト7予想API",
		"version":     "1.0.0",
		"description": "ロト7の予想番号を生成するAPI",
		"endpoints": map[string]interface{}{
			"GET /api/v1/prediction": map[string]interface{}{
				"description": "単一の予想番号を生成",
				"parameters": map[string]string{
					"analysis_range": "分析対象の回数（デフォルト：30）",
					"use_trend":      "傾向を使用するか（デフォルト：true）",
				},
			},
			"POST /api/v1/prediction": map[string]interface{}{
				"description": "単一の予想番号を生成（JSONボディ）",
				"body": map[string]string{
					"analysis_range": "分析対象の回数",
					"use_trend":      "傾向を使用するか",
				},
			},
			"GET /api/v1/predictions": map[string]interface{}{
				"description": "複数の予想番号を生成",
				"parameters": map[string]string{
					"count":          "生成する予想数（最大20、デフォルト：5）",
					"analysis_range": "分析対象の回数（デフォルト：30）",
					"use_trend":      "傾向を使用するか（デフォルト：true）",
				},
			},
			"POST /api/v1/predictions": map[string]interface{}{
				"description": "複数の予想番号を生成（JSONボディ）",
				"body": map[string]string{
					"count":          "生成する予想数",
					"analysis_range": "分析対象の回数",
					"use_trend":      "傾向を使用するか",
				},
			},
		},
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "API仕様情報を取得しました",
		Data:    info,
	})
}
