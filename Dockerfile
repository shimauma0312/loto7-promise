# ===== 開発用ステージ =====
# 開発時はこのステージまででOK
FROM golang:1.23-alpine AS development

# 開発に必要なツールをインストール
RUN apk add --no-cache git

# 作業ディレクトリ設定
WORKDIR /app

# Go モジュールの依存関係をコピー（キャッシュ効率化のため）
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# デフォルトコマンド（開発時は通常オーバーライドされる）
CMD ["go", "run", "./cmd/result"]

# ===== ビルド用ステージ =====
FROM development AS builder

# 本番用バイナリをビルド 
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /result ./cmd/result && \
    CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /frequentNumbers ./cmd/frequentNumbers

# ===== 本番用ステージ =====
FROM alpine:latest AS production

# セキュリティ：最新パッケージにアップデート & 必要最小限のパッケージ
RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

# 作業ディレクトリ設定
WORKDIR /app

# ビルドしたバイナリをコピー
COPY --from=builder /result .
COPY --from=builder /frequentNumbers .

# セキュリティ：非rootユーザーを作成して使用
RUN adduser -D -s /bin/sh appuser && \
    chown -R appuser:appuser /app
USER appuser

# デフォルトコマンド
CMD ["./result"]