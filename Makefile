# ================================
# Loto7 Promise - 超シンプル版
# ================================

.PHONY: help up down shell clean build test test-unit test-integration test-performance test-bench recommend api api-dev api-test

# デフォルトのヘルプ表示
help:
	@echo "Loto7 Promise - 超シンプル使い方"
	@echo ""
	@echo "基本の流れ:"
	@echo "  1. make up    - コンテナ起動"
	@echo "  2. make shell - コンテナに入る"
	@echo "  3. コンテナ内で: go run ./cmd/result"
	@echo "  4. make down  - コンテナ停止"
	@echo ""
	@echo "新機能:"
	@echo "  make recommend - ロト7推薦番号を生成"
	@echo "  make api       - APIサーバーを起動 (ポート8080)"
	@echo "  make api-dev   - 開発モードでAPIサーバーを起動"
	@echo "  make api-test  - APIサーバーのテストを実行"
	@echo ""
	@echo "テスト関連:"
	@echo "  make test               - 全テストを実行"
	@echo "  make test-unit          - ユニットテストのみ実行"
	@echo "  make test-integration   - 統合テストのみ実行"
	@echo "  make test-performance   - パフォーマンステスト実行"
	@echo "  make test-bench         - ベンチマークテスト実行"
	@echo ""
	@echo "その他:"
	@echo "  make build    - 全実行ファイルをビルド"
	@echo "  make clean    - 全部お掃除"
	@echo ""

# コンテナ起動（バックグラウンドで）
up:
	@echo "コンテナ起動中..."
	docker-compose up -d --build
	@echo "起動完了！次は 'make shell' でコンテナに入ってください"

# コンテナに入る
shell:
	@echo "コンテナに入ります..."
	@echo "利用可能なコマンド:"
	@echo "   go run ./cmd/result [回数]              - 過去N回分の抽選結果取得"
	@echo "   go run ./cmd/heatmap -range=50          - 過去50回分のヒートマップ分析"
	@echo "   go run ./cmd/heatmap -number=7 -range=100 - 数字7の詳細分析"
	@echo "   go run ./cmd/heatmap -position=1 -json  - 1番目位置をJSON出力"
	@echo "   go run ./cmd/recommendation             - ロト7推薦番号生成"
	@echo "   go run ./cmd/recommendation -help       - 推薦機能のヘルプ"
	@echo "   go build ./cmd/result                   - 実行ファイル生成"
	@echo "   go test ./...                          - テスト実行"
	@echo ""
	docker-compose exec dev sh

# コンテナ停止
down:
	@echo "コンテナを停止します..."
	docker-compose down
	@echo "停止完了"

# 全部お掃除
clean:
	@echo "全部お掃除中..."
	docker-compose down --volumes --remove-orphans
	docker system prune -f
	@echo "きれいになりました！"

# 実行ファイルをビルド
build:
	@echo "全実行ファイルをビルド中..."
	go build -o ./bin/result ./cmd/result
	go build -o ./bin/heatmap ./cmd/heatmap
	go build -o ./bin/recommendation ./cmd/recommendation
	@echo "ビルド完了！ ./bin/ に実行ファイルが作成されました"

# 全テストを実行
test:
	@echo "統合テストを実行中..."
	go test ./test/ -v
	@echo "全テスト完了"

# 個別テストを実行
test-unit:
	@echo "ユニットテスト実行中..."
	go test ./test/ -v -run "Test[^I]"
	@echo "ユニットテスト完了"

# 統合テストのみ実行
test-integration:
	@echo "統合テスト実行中..."
	go test ./test/ -v -run "TestIntegration"
	@echo "統合テスト完了"

# パフォーマンステストを実行
test-performance:
	@echo "パフォーマンステスト実行中..."
	go test ./test/ -v -run "TestPerformance"
	@echo "パフォーマンステスト完了"

# ベンチマークテストを実行
test-bench:
	@echo "ベンチマークテスト実行中..."
	go test ./test/ -bench=. -benchmem
	@echo "ベンチマークテスト完了"

# ロト7推薦番号を生成
recommend:
	@echo "ロト7推薦番号を生成中..."
	go run ./cmd/recommendation
	@echo "推薦完了"

# APIサーバーを起動
api:
	@echo "APIサーバーを起動中..."
	@echo "ヘルスチェック: http://localhost:8080/health"
	@echo "API仕様書: http://localhost:8080/"
	@echo "ドキュメント: docs/api-spec.md"
	go run ./cmd/api

# 開発モードでAPIサーバーを起動
api-dev:
	@echo "開発モードでAPIサーバーを起動中..."
	@echo "ヘルスチェック: http://localhost:8080/health"
	@echo "ファイル変更で自動再起動"
	GIN_MODE=debug go run ./cmd/api

# APIサーバーのテストを実行
api-test:
	@echo "APIテスト実行中..."
	go test ./test/ -v -run "TestAPI"
	@echo "APIテスト完了"
