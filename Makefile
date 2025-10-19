# ================================
# Loto7 Promise - 超シンプル版
# ================================

.PHONY: help up down shell clean build test recommend

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
	@echo "  make recommend - ロト7推薦番号を生成"
	@echo ""
	@echo "その他:"
	@echo "  make build    - 全実行ファイルをビルド"
	@echo "  make test     - 全テストを実行"
	@echo "  make clean    - 全部お掃除"
	@echo ""

# コンテナ起動（バックグラウンドで）
up:
	@echo "コンテナ起動中..."
	docker-compose up -d --build
	@echo "起動完了！次は 'make shell' でコンテナに入ってください"

# コンテナに入る
shell:
	@echo "🔧 コンテナに入ります..."
	@echo "💡 利用可能なコマンド:"
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
	@echo "� コンテナを停止します..."
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
	@echo "📦 全実行ファイルをビルド中..."
	go build -o ./bin/result ./cmd/result
	go build -o ./bin/heatmap ./cmd/heatmap
	go build -o ./bin/recommendation ./cmd/recommendation
	@echo "✅ ビルド完了！ ./bin/ に実行ファイルが作成されました"

# 全テストを実行
test:
	@echo "🧪 テスト実行中..."
	go test ./...
	@echo "✅ テスト完了"

# ロト7推薦番号を生成
recommend:
	@echo "🎲 ロト7推薦番号を生成中..."
	go run ./cmd/recommendation
	@echo "✅ 推薦完了"
