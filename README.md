# Loto7 Promise

ロト7あてる

## Docker環境構築
```bash

// Docker環境の起動
docker compose up -d --build

// コンテナアクセス
docker compose exec dev sh

// 環境の停止
docker compose down

```

## APIサーバー構築呼び出し

### 1. APIサーバー起動
```bash
go run cmd/api/main.go
```

### 2. APIエンドポイント

ベースURL: `http://localhost:8080`

#### ヘルスチェック
```bash
curl http://localhost:8080/health
```

#### 抽選結果取得
```bash
# 過去10回分の結果取得
curl http://localhost:8080/api/results?count=10

# 過去100回分の結果取得
curl http://localhost:8080/api/results?count=100
```

#### ヒートマップ分析
```bash
# 基本ヒートマップ（過去50回）
curl http://localhost:8080/api/heatmap

# 詳細ヒートマップ（過去100回）
curl http://localhost:8080/api/heatmap?range=100
```

#### 推薦番号生成
```bash
# 推薦番号の生成
curl http://localhost:8080/api/recommendation
```

## APIデプロイ

### 前提条件
- Go 1.23 以上がインストールされてること
- 

### 1. ビルド
```bash
# APIサーバーをビルド
go build -o loto7-api ./cmd/api

# Linux
GOOS=linux GOARCH=arm64 go build -o loto7-api-linux ./cmd/api

# Windows
GOOS=windows GOARCH=arm64 go build -o loto7-api.exe ./cmd/api
```

### 3. サーバー配置
```bash
# ディレクトリを作成
mkdir -p /opt/loto7-api
cd /opt/loto7-api

# バイナリとキャッシュディレクトリをコピー
cp /path/to/loto7-api-linux ./loto7-api
cp -r /path/to/cache ./cache

# 実行権限
chmod +x loto7-api
```

### 4. APIサーバー起動

#### フォアグラウンド
```bash
# 本番モードで起動
GIN_MODE=release ./loto7-api
```

#### バックグラウンド
```bash
# （nohup
nohup GIN_MODE=release ./loto7-api > loto7-api.log 2>&1 &

# プロセスID
ps aux | grep loto7-api
```

#### systemdサービスとして起動
`/etc/systemd/system/loto7-api.service` を作成：

```ini
[Unit]
Description=Loto7 API Server
After=network.target

[Service]
Type=simple
User=loto7
Group=loto7
WorkingDirectory=/opt/loto7-api
ExecStart=/opt/loto7-api/loto7-api
Environment=GIN_MODE=release
Environment=PORT=8080
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

```bash
# 有効化
sudo systemctl enable loto7-api

# 起動
sudo systemctl start loto7-api

# サービス状態確認
sudo systemctl status loto7-api

# ログ確認
sudo journalctl -u loto7-api -f
```

### 5. サーバー管理コマンド

#### サーバーの停止
```bash
# プロセスIDを確認して終了
pkill loto7-api

# systemdの場合
sudo systemctl stop loto7-api
```

#### サーバー再起動
```bash
# systemdの場合
sudo systemctl restart loto7-api
```

#### サーバーのアップデート
```bash
# 最新コードをプル
git pull origin main

# 新しいバイナリをビルド
go build -o loto7-api-new ./cmd/api

# サービスを停止
sudo systemctl stop loto7-api

# バイナリを置き換え
mv loto7-api-new loto7-api

# サービスを再起動
sudo systemctl start loto7-api
```

#### ログ監視
```bash
# nohupの場合
tail -f loto7-api.log

# systemdの場合
sudo journalctl -u loto7-api -f
```

### 6. エンドポイント確認

#### API エンドポイントテスト
```bash
# ヘルスチェック
curl http://localhost:8080/health

# 抽選結果取得
curl http://localhost:8080/api/results?count=10

# ヒートマップ
curl http://localhost:8080/api/heatmap

# 推薦番号生成
curl http://localhost:8080/api/recommendation
```
