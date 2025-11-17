# Loto7 Promise

ロト7あてる

## Docker環境
```bash
docker compose up -d --build

docker compose exec dev sh

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

#### ヒートマップ
```bash
curl http://localhost:8080/api/heatmap

curl http://localhost:8080/api/heatmap?range=100
```

#### 予測番号生成
```bash
# 番号の生成
curl http://localhost:8080/api/recommendation
```

## APIデプロイ

### 前提条件
- Go 1.23 以上がインストールされてること
- 

### 1. ビルド
```bash
go build -o loto7-api ./cmd/api
```

### 3. サーバー配置
```bash
# ディレクトリ作る
mkdir -p /opt/loto7-api
cd /opt/loto7-api

# バイナリとキャッシュをコピー
cp /path/to/loto7-api-linux ./loto7-api
cp -r /path/to/cache ./cache

# 実行権限
chmod +x loto7-api
```

### 4. systemd APIサーバー化

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
ExecStart=/opt/loto7-api/loto7-api # Pathにあたる
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

### 5. 鯖管

#### サーバー停止
```bash
sudo systemctl stop loto7-api
```

#### サーバー再起動
```bash
sudo systemctl restart loto7-api
```

#### サーバーのアップデート
```bash
git pull origin master

go build -o loto7-api-new ./cmd/api

sudo systemctl stop loto7-api

mv loto7-api-new /path/loto7-api

sudo systemctl start loto7-api
```

#### ログ見る
```bash
sudo journalctl -u loto7-api -f
```