# Loto7 Promise

ロト7あてる

## Docker環境
```bash
docker compose up -d --build

docker compose exec dev sh

docker compose down

```

## コマンドラインツール

### ランダム番号生成
各数字位置で過去に出現した範囲内から**重み付き乱択**で数字を選択します。
出現頻度の高い数字ほど選ばれやすくなり、実際のロト7の傾向を再現します。

**特徴:**
- 各位置での過去の出現頻度を分析
- 頻出数字に高い重みを付与（例: 1番目の位置では1-5が選ばれやすい）
- **直近ペナルティ**: 直近1回、2回で出現した数字は選ばれにくくなる
  - 直近1回目: 重み×0.05（大幅に選ばれにくい）
  - 直近2回目: 重み×0.3（選ばれにくい）
- 未出現の数字は自動的に除外
- より本来のロト7のランダム性を再現

```bash
# 基本的な使い方（過去100回分を分析、5組生成）
go run cmd/random/main.go

# 過去150回分を分析、3組生成
go run cmd/random/main.go -history 150 -count 3

# 各位置の出現範囲と上位頻出数字を表示
go run cmd/random/main.go -verbose

# ヘルプを表示
go run cmd/random/main.go -help
```

### 統計ベース推薦番号生成
統計分析に基づいて番号を推薦します。

```bash
# 基本的な使い方
go run cmd/recommendation/main.go

# 過去150回分を分析、3組生成
go run cmd/recommendation/main.go -history 150 -count 3
```

### 1等当選シミュレーター
ユーザーの選択数字で1等が当選するまでのシミュレーションを行います。

```bash
# 基本的な使い方（自動生成された数字で1回シミュレーション）
go run cmd/simulation/main.go

# 指定した数字で1回シミュレーション
go run cmd/simulation/main.go -numbers 1,5,10,15,20,25,30

# 10回シミュレーションして統計を取る
go run cmd/simulation/main.go -simulations 10

# 指定した数字で10回シミュレーション
go run cmd/simulation/main.go -numbers 1,5,10,15,20,25,30 -simulations 10

# ヘルプを表示
go run cmd/simulation/main.go -help
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

#### ランダム番号生成
```bash
# 重み付き乱択による番号生成
curl http://localhost:8080/api/v1/random

# 過去の抽選データ範囲と生成数を指定
curl "http://localhost:8080/api/v1/random?history=150&count=3"
```
**注:** 各位置での出現頻度に応じた重み付きランダム選択を使用

#### 1等当選シミュレーション
```bash
# 自動生成された数字で1回シミュレーション
curl -X POST http://localhost:8080/api/v1/simulation \
  -H "Content-Type: application/json" \
  -d '{}'

# 指定した数字で1回シミュレーション
curl -X POST http://localhost:8080/api/v1/simulation \
  -H "Content-Type: application/json" \
  -d '{
    "user_numbers": [1, 5, 10, 15, 20, 25, 30]
  }'

# 10回シミュレーションして統計を取る
curl -X POST http://localhost:8080/api/v1/simulation \
  -H "Content-Type: application/json" \
  -d '{
    "user_numbers": [1, 5, 10, 15, 20, 25, 30],
    "simulation_count": 10
  }'

# 詳細設定
curl -X POST http://localhost:8080/api/v1/simulation \
  -H "Content-Type: application/json" \
  -d '{
    "user_numbers": [1, 5, 10, 15, 20, 25, 30],
    "simulation_count": 5,
    "history_count": 150
  }'
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