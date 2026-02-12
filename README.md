# Loto7 Promise

ロト7あてる

## Docker環境
```bash
docker compose up -d --build

docker compose exec dev sh

docker compose down

```

## コマンドラ

### ランダム番号生成
各数字位置で過去に出現した範囲内から**重み付き乱択**で数字を選択します。
出現頻度の高い数字ほど選ばれやすくなり、実際のロト7の傾向を再現します。

```bash
# 基本的な使い方（過去100回分を分析、5組生成）
go run cmd/random/main.go

# 過去150回分を分析、3組生成
go run cmd/random/main.go -history 150 -count 3

# 各位置の出現範囲と上位頻出数字を表示
go run cmd/random/main.go -verbose

# ヘルプ
go run cmd/random/main.go -help
```

### 統計ベース推薦番号生成
統計に基づいて番号を推薦します。

```bash
# 基本的な使い方
go run cmd/recommendation/main.go

# 過去150回分を分析、3組生成
go run cmd/recommendation/main.go -history 150 -count 3
```

### 1等当選シミュレーター
ユーザーの選択数字で1等が当選するまでのシミュレーションを行う

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

## デプロイ

本番へのデプロイ手順は [docs/deploy.md](docs/deploy.md) を参照