# ================================
# Loto7 Promise - 超シンプル版
# ================================

.PHONY: help up down shell clean

# デフォルトのヘルプ表示
help:
	@echo "🎯 Loto7 Promise - 超シンプル使い方"
	@echo ""
	@echo "基本の流れ:"
	@echo "  1. make up    - コンテナ起動"
	@echo "  2. make shell - コンテナに入る"
	@echo "  3. コンテナ内で: go run ./cmd/result"
	@echo "  4. make down  - コンテナ停止"
	@echo ""
	@echo "その他:"
	@echo "  make clean    - 全部お掃除"
	@echo ""

# コンテナ起動（バックグラウンドで）
up:
	@echo "🚀 コンテナ起動中..."
	docker-compose up -d --build
	@echo "✅ 起動完了！次は 'make shell' でコンテナに入ってください"

# コンテナに入る
shell:
	@echo "🔧 コンテナに入ります..."
	@echo "💡 コンテナ内でのコマンド例:"
	@echo "   go run ./cmd/result"
	@echo "   go run ./cmd/frequentNumbers"
	@echo "   go run ./cmd/numberCount -number=7 -range=100"
	@echo "   go build ./cmd/result"
	@echo "   go test ./..."
	@echo ""
	docker-compose exec dev sh

# コンテナ停止
down:
	@echo "� コンテナを停止します..."
	docker-compose down
	@echo "✅ 停止完了"

# 全部お掃除
clean:
	@echo "🧹 全部お掃除中..."
	docker-compose down --volumes --remove-orphans
	docker system prune -f
	@echo "✨ きれいになりました！"
