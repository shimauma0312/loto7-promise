# ================================
# Loto7 Promise - 超シンプル版
# ================================

.PHONY: help up down shell clean api-up api-down api-logs api-test

# デフォルトのヘルプ表示
help:
	@echo "Loto7 Promise - 超シンプル使い方"
	@echo ""
	@echo "📋 開発作業:"
	@echo "  1. make up    - 開発用コンテナ起動"
	@echo "  2. make shell - コンテナに入る"
	@echo "  3. コンテナ内で: go run ./cmd/result"
	@echo "  4. make down  - コンテナ停止"
	@echo ""
	@echo "🌐 API サーバー:"
	@echo "  make api-up   - API + nginx 起動"
	@echo "  make api-down - API サーバー停止"
	@echo "  make api-logs - API ログ確認"
	@echo "  make api-test - API 動作テスト"
	@echo ""
	@echo "🧹 その他:"
	@echo "  make clean    - 全部お掃除"
	@echo ""

# コンテナ起動（バックグラウンドで）
up:
	@echo "開発用コンテナ起動中..."
	docker-compose up -d --build dev
	@echo "起動完了！次は 'make shell' でコンテナに入ってください"

# API サーバー起動
api-up:
	@echo "🚀 API サーバー + nginx を起動中..."
	docker-compose up -d --build api nginx
	@echo "✅ API サーバー起動完了！"
	@echo "📍 アクセス先:"
	@echo "   - nginx経由: http://localhost/api/v1/info"
	@echo "   - 直接アクセス: http://localhost:8080/api/v1/info"
	@echo "   - ヘルスチェック: http://localhost/health"

# API サーバー停止
api-down:
	@echo "🛑 API サーバーを停止します..."
	docker-compose down api nginx
	@echo "停止完了"

# API ログ確認
api-logs:
	@echo "📊 API サーバーのログを表示します..."
	docker-compose logs -f api

# API 動作テスト
api-test:
	@echo "🧪 API 動作テスト中..."
	@echo ""
	@echo "⚕️  ヘルスチェック:"
	@curl -s http://localhost/health | head -c 200 && echo "..." || echo "❌ API が起動していません"
	@echo ""
	@echo "📋 API情報取得:"
	@curl -s http://localhost/api/v1/info | head -c 200 && echo "..." || echo "❌ API情報取得失敗"
	@echo ""
	@echo "🎲 単一予想生成:"
	@curl -s "http://localhost/api/v1/prediction?use_trend=true" | head -c 200 && echo "..." || echo "❌ 予想生成失敗"
	@echo ""
	@echo "✅ テスト完了！詳細は 'make api-logs' でログを確認してください"

# コンテナに入る
shell:
	@echo "🔧 コンテナに入ります..."
	@echo "💡 コンテナ内でのコマンド例:"
	@echo "   go run ./cmd/result"
	@echo "   go run ./cmd/frequentNumbers"
	@echo "   go run ./cmd/numberCount -number=7 -range=100"
	@echo "   go run ./cmd/consecutivePattern -range=100"
	@echo "   go run ./cmd/heatmap -range=50"
	@echo "   go run ./cmd/heatmap -number=7 -range=100"
	@echo "   go run ./cmd/heatmap -position=1 -json"
	@echo "   go run ./cmd/randomPrediction -count=5 -range=200"
	@echo "   go run ./cmd/randomPrediction -random -count=3"
	@echo "   go run ./cmd/bonusNumber -trend -top=15 -range=100"
	@echo "   go run ./cmd/bonusNumber -number=7 -range=50"
	@echo "   go build ./cmd/result"
	@echo "   go test ./..."
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
