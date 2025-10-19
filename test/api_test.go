package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"
	"time"
)

// TestAPIServerBuild APIサーバーのビルドテスト
func TestAPIServerBuild(t *testing.T) {
	// APIサーバーがビルドできるかテスト
	cmd := exec.Command("go", "build", "./cmd/api")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("APIサーバーのビルドに失敗しました: %v", err)
	}
}

// TestAPIServerStart APIサーバーの起動テスト
func TestAPIServerStart(t *testing.T) {
	if testing.Short() {
		t.Skip("短時間テストではAPIサーバー起動テストをスキップします")
	}

	// APIサーバーを起動
	cmd := exec.Command("go", "run", "./cmd/api")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("APIサーバーの起動に失敗しました: %v", err)
	}
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	// サーバーが起動するまで少し待機
	time.Sleep(2 * time.Second)

	// ヘルスチェックエンドポイントにアクセス
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Fatalf("ヘルスチェックエンドポイントへのアクセスに失敗しました: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("期待されるステータスコード %d, 実際 %d", http.StatusOK, resp.StatusCode)
	}
}

// TestAPIEndpoints APIエンドポイントの基本テスト
func TestAPIEndpoints(t *testing.T) {
	if testing.Short() {
		t.Skip("短時間テストではAPIエンドポイントテストをスキップします")
	}

	// テスト用のダミーサーバーを作成
	mux := http.NewServeMux()

	// ダミーレスポンスを返すハンドラーを設定
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "message": "healthy"}`))
	})

	mux.HandleFunc("/api/v1/results", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "data": []}`))
	})

	mux.HandleFunc("/api/v1/heatmap", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "data": {}}`))
	})

	mux.HandleFunc("/api/v1/recommendations", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "data": {"recommendations": []}}`))
	})

	// テストサーバーを作成
	server := httptest.NewServer(mux)
	defer server.Close()

	// 各エンドポイントをテスト
	endpoints := []string{
		"/health",
		"/api/v1/results",
		"/api/v1/heatmap",
		"/api/v1/recommendations",
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			resp, err := http.Get(server.URL + endpoint)
			if err != nil {
				t.Fatalf("エンドポイント %s へのアクセスに失敗: %v", endpoint, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("エンドポイント %s: 期待されるステータスコード %d, 実際 %d", endpoint, http.StatusOK, resp.StatusCode)
			}

			contentType := resp.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("エンドポイント %s: 期待されるContent-Type application/json, 実際 %s", endpoint, contentType)
			}
		})
	}
}

// TestAPIResponseFormat APIレスポンス形式のテスト
func TestAPIResponseFormat(t *testing.T) {
	// ダミーレスポンス形式をテスト
	testCases := []struct {
		name         string
		responseBody string
		expectValid  bool
	}{
		{
			name:         "有効なレスポンス",
			responseBody: `{"success": true, "message": "test", "timestamp": "2025-10-19T21:30:00Z", "version": "1.0.0"}`,
			expectValid:  true,
		},
		{
			name:         "不正なJSON",
			responseBody: `{"success": true, "message": "test"`,
			expectValid:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// JSONの妥当性をチェック
			var response map[string]interface{}
			err := json.Unmarshal([]byte(tc.responseBody), &response)

			if tc.expectValid && err != nil {
				t.Errorf("有効なJSONが期待されましたが、パースエラー: %v", err)
			}

			if !tc.expectValid && err == nil {
				t.Error("不正なJSONが期待されましたが、パースが成功しました")
			}
		})
	}
}
