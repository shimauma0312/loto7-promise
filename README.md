# Loto7 Promise



ロト7の過去の結果を分析して、統計的根拠に基づく推薦番号を生成するAPI



## 📋 主要機能ロト7の過去の結果を分析して、統計的根拠に基づく推薦番号を生成するWebAPIシステムです。```


#### 1. **結果取得システム** (`internal/result/`)## 📋 主要機能```

- ロト7の過去の抽選結果データの取得・管理

- CSVファイルからのデータ読み込み

- 次回抽選日の自動計算

### 機能```

#### 2. **ヒートマップ分析** (`internal/heatmap/`)

- 各数字の出現頻度を視覚的に分析go build -o <app name> ./cmd/result

- 位置別数字の出現パターン分析

- 統計的データの可視化#### 1. **結果取得システム** (`internal/result/`)```



#### 3. **推薦エンジン** (`internal/recommendation/`)- ロト7の過去の抽選結果データの取得・管理

- 過去のデータに基づく推薦番号生成

- 複数の推薦組み合わせ提供- CSVファイルからのデータ読み込み```

- カスタマイズ可能な推薦設定

- 次回抽選日の自動計算go build -o <app name> ./cmd/numberCount

#### 4. **REST APIサーバー** (`cmd/api/`)

- Ginフレームワークベースの高性能WebAPI```

- 統一されたレスポンス形式

- CORS対応でフロントエンドと連携可能#### 2. **ヒートマップ分析** (`internal/heatmap/`)

- 包括的なエラーハンドリング

- 各数字の出現頻度分析```

## 🛠️ 技術構成

- 位置別数字の出現パターン分析go build -o <app name> ./cmd/consecutivePattern

- **言語**: Go 1.23

- **Webフレームワーク**: Gin v1.11.0- 統計の可視化```

- **アーキテクチャ**: Clean Architecture

- **テスト**: 統合テストスイート

- **API**: RESTful設計

#### 3. **推薦エンジン** (`internal/recommendation/`)```

## 🚀 APIエンドポイント

- 過去のデータに基づく推薦番号生成go build -o <app name> ./cmd/heatmap

### ヘルスチェック

```- いくつかの推薦組み合わせ```

GET /health

```- カスタム可能な推薦設定



### 抽選結果取得## Docker

```

GET /api/v1/results           # デフォルト10回分#### 4. **REST APIサーバー** (`cmd/api/`)

GET /api/v1/results/{count}   # 指定回数分

```- WebAPIThis project includes Docker support for easy deployment and execution.



### ヒートマップデータ- 統一されたレスポンス形式

```

GET /api/v1/heatmap                        # デフォルト50回分- CORS対応でフロントエンドと連携可能### Building the Docker image

GET /api/v1/heatmap/{range}                # 指定範囲

GET /api/v1/heatmap/{range}/number/{number}    # 特定数字の分析- 包括的なエラーハンドリング

GET /api/v1/heatmap/{range}/position/{position} # 特定位置の分析

``````bash



### 推薦番号生成## 🛠️ 技術構成docker build -t loto7-promise .

```

GET  /api/v1/recommendations  # デフォルト設定```

POST /api/v1/recommendations  # カスタム設定

```- **言語**: Go 1.23



## 📊 API使用例とレスポンス- **Webフレームワーク**: Gin v1.11.0### Running with Docker



### 1. ヘルスチェック- **アーキテクチャ**: Clean Architecture



**リクエスト:**- **テスト**: 統合テストスイートRun the result application:

```bash

curl -X GET http://localhost:8080/health- **API**: RESTful設計```bash

```

# Show usage

**レスポンス:**

```json## 🚀 APIエンドポイントdocker run --rm loto7-promise

{

  "success": true,

  "message": "Loto7 Promise API is healthy",

  "data": {### ヘルスチェック# Get past 10 results  

    "status": "ok",

    "uptime": "5m30s",```docker run --rm loto7-promise ./result 10

    "version": "1.0.0"

  },GET /health```

  "timestamp": "2025-10-19T22:00:00Z",

  "version": "1.0.0",```

  "error": null

}Run the frequent numbers application:

```

### 抽選結果取得```bash

### 2. 抽選結果取得

```docker run --rm loto7-promise ./frequentNumbers

**リクエスト (デフォルト10回分):**

```bashGET /api/v1/results           # デフォルト10回分```

curl -X GET http://localhost:8080/api/v1/results

```GET /api/v1/results/{count}   # 指定回数分



**リクエスト (指定回数):**```Run the number count application:

```bash

curl -X GET http://localhost:8080/api/v1/results/5```bash

```

### ヒートマップデータ# Count specific number appearances

**レスポンス:**

```json```docker run --rm loto7-promise ./numberCount -number=7 -range=100

{

  "success": true,GET /api/v1/heatmap                        # デフォルト50回分

  "message": "抽選結果を5回分取得しました",

  "data": {GET /api/v1/heatmap/{range}                # 指定範囲# Count multiple numbers

    "results": [

      {GET /api/v1/heatmap/{range}/number/{number}    # 特定数字の分析docker run --rm loto7-promise ./numberCount -number=1,7,23 -range=50

        "round": 550,

        "date": "2025-10-18",GET /api/v1/heatmap/{range}/position/{position} # 特定位置の分析```

        "numbers": [3, 7, 14, 21, 28, 35, 37],

        "bonus": [2, 15]```

      },

      {Run the consecutive pattern analysis application:

        "round": 549,

        "date": "2025-10-11",### 推薦番号生成```bash

        "numbers": [1, 8, 15, 22, 29, 33, 39],

        "bonus": [5, 18]```# Analyze consecutive and same-digit patterns (default: 100 draws)

      }

    ],GET  /api/v1/recommendations  # デフォルト設定docker run --rm loto7-promise ./consecutivePattern

    "count": 5,

    "next_draw_date": "2025-10-25"POST /api/v1/recommendations  # カスタム設定

  },

  "timestamp": "2025-10-19T22:00:00Z",```# Analyze patterns for specific range

  "version": "1.0.0",

  "error": nulldocker run --rm loto7-promise ./consecutivePattern -range=200

}

```## 📦 プロジェクト構造```



### 3. ヒートマップデータ取得



**リクエスト (デフォルト50回分):**```Run the heatmap analysis application:

```bash

curl -X GET http://localhost:8080/api/v1/heatmaploto7/```bash

```

├── cmd/# Generate position-based heatmap (default: 50 draws)

**リクエスト (指定範囲):**

```bash│   ├── api/            # APIサーバーdocker run --rm loto7-promise ./heatmap

curl -X GET http://localhost:8080/api/v1/heatmap/100

```│   ├── heatmap/        # ヒートマップ生成CLI



**レスポンス:**│   └── result/         # 結果取得CLI# Generate heatmap for specific range

```json

{├── internal/docker run --rm loto7-promise ./heatmap -range=100

  "success": true,

  "message": "ヒートマップデータを100回分生成しました",│   ├── api/            # API共通機能

  "data": {

    "range": 100,│   ├── cache/          # キャッシュ機能# Show specific number's position details

    "number_frequency": {

      "1": 15, "2": 12, "3": 18, "7": 22, "14": 16,│   ├── heatmap/        # ヒートマップ分析docker run --rm loto7-promise ./heatmap -number=7 -range=100

      "21": 20, "28": 14, "35": 17, "37": 13

    },│   ├── recommendation/ # 推薦エンジン

    "position_data": {

      "position_1": {"most_frequent": 3, "frequency": 8},│   └── result/         # 結果取得システム# Show specific position's number details

      "position_2": {"most_frequent": 7, "frequency": 9},

      "position_3": {"most_frequent": 14, "frequency": 7}├── test/               # 統合テストスイートdocker run --rm loto7-promise ./heatmap -position=1

    },

    "statistics": {├── cache/              # データキャッシュ

      "total_draws": 100,

      "avg_frequency": 14.3,└── docs/               # ドキュメント# Output in JSON format

      "most_frequent_number": 7,

      "least_frequent_number": 2```docker run --rm loto7-promise ./heatmap -json -range=30

    }

  },```

  "timestamp": "2025-10-19T22:00:00Z",

  "version": "1.0.0",## 🧪 テスト実行

  "error": null

}### Using Docker Compose

```

```bash

### 4. 特定数字の分析

# 全テスト実行For easier management, use Docker Compose:

**リクエスト:**

```bashgo test ./test -v

curl -X GET http://localhost:8080/api/v1/heatmap/50/number/7

``````bash



**レスポンス:**# 特定機能テスト# Build all services

```json

{go test ./test -run="TestHeatmap" -vdocker compose build

  "success": true,

  "message": "数字7の分析結果を50回分取得しました",go test ./test -run="TestRecommendation" -v

  "data": {

    "number": 7,go test ./test -run="TestResult" -v# Run result service with argument

    "range": 50,

    "total_appearances": 11,docker compose run --rm result ./result 5

    "frequency_percentage": 22.0,

    "position_distribution": {# APIテスト（基本機能）

      "1": 2, "2": 3, "3": 2, "4": 1, "5": 2, "6": 1, "7": 0

    },go test ./test -run="TestAPIEndpoints|TestAPIResponseFormat" -v# Run frequent numbers service

    "recent_appearances": [548, 545, 542, 539, 537],

    "prediction_score": 0.78```docker compose run --rm frequent

  },

  "timestamp": "2025-10-19T22:00:00Z",

  "version": "1.0.0",

  "error": null## 🏃‍♂️ 起動方法# Run interactive shell

}

```docker compose run --rm interactive



### 5. 特定位置の分析### APIサーバー起動```



**リクエスト:**```bash

```bash

curl -X GET http://localhost:8080/api/v1/heatmap/30/position/1go run ./cmd/api### Image Details

```

```

**レスポンス:**

```json- **Base images**: golang:1.23-alpine (build) + alpine:latest (runtime)

{

  "success": true,### CLI実行- **Image size**: ~39MB

  "message": "位置1の分析結果を30回分取得しました",

  "data": {```bash- **Applications**: Both `result` and `frequentNumbers` binaries included

    "position": 1,

    "range": 30,# ヒートマップ生成- **Security**: Runs as non-root user `appuser`

    "number_distribution": {

      "1": 4, "2": 2, "3": 5, "4": 3, "5": 2, "6": 1, "7": 0, "8": 3go run ./cmd/heatmap- **Dependencies**: All Go dependencies vendored for offline builds

    },

    "most_frequent_number": 3,

    "frequency": 5,# 結果取得

    "percentage": 16.67,go run ./cmd/result

    "statistics": {

      "total_numbers": 30,# 推薦番号生成

      "unique_numbers": 7,go run ./cmd/recommendation

      "avg_value": 3.2```

    }

  },## 📊 レスポンス形式

  "timestamp": "2025-10-19T22:00:00Z",

  "version": "1.0.0",```json

  "error": null{

}  "success": true,

```  "message": "成功メッセージ",

  "data": {

### 6. 推薦番号生成 (GET)    // 実際のデータ

  },

**リクエスト (デフォルト設定):**  "timestamp": "2025-10-19T22:00:00Z",

```bash  "version": "1.0.0",

curl -X GET http://localhost:8080/api/v1/recommendations  "error": null

```}

```

**リクエスト (パラメータ付き):**

```bash## 🔧 開発・テスト状況

curl -X GET "http://localhost:8080/api/v1/recommendations?count=3&history=100"

```### ✅ 完了項目

- [x] 基本アーキテクチャ設計

**レスポンス:**- [x] 結果取得システム実装

```json- [x] ヒートマップ分析機能

{- [x] 推薦エンジン開発

  "success": true,- [x] REST APIサーバー構築

  "message": "推薦番号を3組生成しました",- [x] 統一レスポンス形式

  "data": {- [x] CORS対応

    "recommendations": [- [x] エラーハンドリング

      {- [x] 統合テストスイート

        "set": 1,- [x] プロジェクト文書化

        "numbers": [3, 7, 14, 21, 28, 35, 37],

        "confidence": 0.82,### 🔄 テスト状況

        "reason": "高頻度出現パターンに基づく推薦"- **機能テスト**: ✅ 全て合格

      },- **統合テスト**: ✅ 全て合格  

      {- **APIテスト**: ✅ 基本機能合格

        "set": 2,- **パフォーマンステスト**: ✅ 合格

        "numbers": [1, 8, 15, 22, 29, 33, 39],

        "confidence": 0.76,### 📈 品質指標

        "reason": "バランス型分布パターンに基づく推薦"- **コード品質**: 高品質（Clean Architecture採用）

      },- **テストカバレッジ**: 包括的な統合テスト

      {- **API設計**: RESTful準拠

        "set": 3,- **エラーハンドリング**: 包括的対応

        "numbers": [5, 12, 19, 26, 31, 36, 42],

        "confidence": 0.71,## 🎯 使用例

        "reason": "統計的予測モデルに基づく推薦"

      }### 推薦番号取得

    ],```bash

    "config": {curl http://localhost:8080/api/v1/recommendations

      "count": 3,```

      "history_lookback": 100,

      "algorithm": "hybrid"### ヒートマップデータ取得

    },```bash

    "generated_at": "2025-10-19T22:00:00Z"curl http://localhost:8080/api/v1/heatmap/30

  },```

  "timestamp": "2025-10-19T22:00:00Z",

  "version": "1.0.0",### 抽選結果取得

  "error": null```bash

}curl http://localhost:8080/api/v1/results/5

``````



### 7. 推薦番号生成 (POST - カスタム設定)## 🐳 Docker対応



**リクエスト:**### イメージビルド

```bash```bash

curl -X POST http://localhost:8080/api/v1/recommendations \docker build -t loto7-promise .

  -H "Content-Type: application/json" \```

  -d '{

    "max_recommendations": 5,### Docker実行

    "history_lookback": 200,```bash

    "use_position_analysis": true,# APIサーバー起動

    "use_frequency_analysis": true,docker run -p 8080:8080 loto7-promise

    "exclude_numbers": [13, 24],

    "prefer_numbers": [7, 14, 21]# CLI実行例

  }'docker run --rm loto7-promise ./result 10

```docker run --rm loto7-promise ./heatmap -range=100

```

**レスポンス:**

```json### Docker Compose

{```bash

  "success": true,# 全サービスビルド

  "message": "カスタム設定で推薦番号を5組生成しました",docker compose build

  "data": {

    "recommendations": [# サービス実行

      {docker compose run --rm api

        "set": 1,```

        "numbers": [7, 14, 21, 28, 35, 38, 42],

        "confidence": 0.89,## 📝 まとめ

        "reason": "優先数字と高頻度パターンの組み合わせ"

      },このプロジェクトは、ロト7の過去データを基にした統計的推薦システムとして、以下の価値を提供します：

      {

        "set": 2,1. **データ駆動**: 過去の結果データに基づく客観的分析

        "numbers": [3, 7, 16, 21, 29, 34, 41],2. **拡張性**: Clean Architectureによる保守性の高い設計

        "confidence": 0.84,3. **使いやすさ**: シンプルなREST APIインターフェース

        "reason": "位置分析と頻度分析の最適化"4. **信頼性**: 包括的なテストスイートによる品質保証

      }

    ],本システムは完全に機能する状態であり、フロントエンドアプリケーションとの統合や、さらなる機能拡張の基盤として活用できます。
    "config": {
      "max_recommendations": 5,
      "history_lookback": 200,
      "use_position_analysis": true,
      "use_frequency_analysis": true,
      "exclude_numbers": [13, 24],
      "prefer_numbers": [7, 14, 21]
    },
    "generated_at": "2025-10-19T22:00:00Z"
  },
  "timestamp": "2025-10-19T22:00:00Z",
  "version": "1.0.0",
  "error": null
}
```

### 8. エラーレスポンス例

**リクエスト (無効なパラメータ):**
```bash
curl -X GET http://localhost:8080/api/v1/results/invalid
```

**レスポンス:**
```json
{
  "success": false,
  "message": "リクエストパラメータが無効です",
  "data": null,
  "timestamp": "2025-10-19T22:00:00Z",
  "version": "1.0.0",
  "error": {
    "code": "INVALID_PARAMETER",
    "message": "count parameter must be a positive integer",
    "details": "received: 'invalid', expected: number between 1 and 1000"
  }
}
```

## 📦 プロジェクト構造

```
loto7/
├── cmd/
│   ├── api/            # APIサーバー
│   ├── heatmap/        # ヒートマップ生成CLI
│   └── result/         # 結果取得CLI
├── internal/
│   ├── api/            # API共通機能
│   ├── cache/          # キャッシュ機能
│   ├── heatmap/        # ヒートマップ分析
│   ├── recommendation/ # 推薦エンジン
│   └── result/         # 結果取得システム
├── test/               # 統合テストスイート
├── cache/              # データキャッシュ
└── docs/               # ドキュメント
```

## 🧪 テスト実行

```bash
# 全テスト実行
go test ./test -v

# 特定機能テスト
go test ./test -run="TestHeatmap" -v
go test ./test -run="TestRecommendation" -v
go test ./test -run="TestResult" -v

# APIテスト（基本機能）
go test ./test -run="TestAPIEndpoints|TestAPIResponseFormat" -v
```

## 🏃‍♂️ 起動方法

### APIサーバー起動
```bash
go run ./cmd/api
```
サーバー起動後、以下のURLでアクセス可能：
- API Base URL: `http://localhost:8080`
- Health Check: `http://localhost:8080/health`
- API Documentation: `http://localhost:8080/` (今後実装予定)

### CLI実行
```bash
# ヒートマップ生成
go run ./cmd/heatmap

# 結果取得
go run ./cmd/result

# 推薦番号生成
go run ./cmd/recommendation
```

## 🐳 Docker対応

### イメージビルド
```bash
docker build -t loto7-promise .
```

### Docker実行
```bash
# APIサーバー起動
docker run -p 8080:8080 loto7-promise

# CLI実行例
docker run --rm loto7-promise ./result 10
docker run --rm loto7-promise ./heatmap -range=100
```

### Docker Compose
```bash
# 全サービスビルド
docker compose build

# サービス実行
docker compose run --rm api
```

## 🔧 開発・テスト状況

### ✅ 完了項目
- [x] 基本アーキテクチャ設計
- [x] 結果取得システム実装
- [x] ヒートマップ分析機能
- [x] 推薦エンジン開発
- [x] REST APIサーバー構築
- [x] 統一レスポンス形式
- [x] CORS対応
- [x] エラーハンドリング
- [x] 統合テストスイート
- [x] プロジェクト文書化

### 🔄 テスト状況
- **機能テスト**: ✅ 全て合格
- **統合テスト**: ✅ 全て合格  
- **APIテスト**: ✅ 基本機能合格
- **パフォーマンステスト**: ✅ 合格

### 📈 品質指標
- **コード品質**: 高品質（Clean Architecture採用）
- **テストカバレッジ**: 包括的な統合テスト
- **API設計**: RESTful準拠
- **エラーハンドリング**: 包括的対応

## 💡 フロントエンド統合例

### JavaScript (Fetch API)
```javascript
// 推薦番号取得
async function getRecommendations() {
  try {
    const response = await fetch('http://localhost:8080/api/v1/recommendations');
    const data = await response.json();
    
    if (data.success) {
      console.log('推薦番号:', data.data.recommendations);
    } else {
      console.error('エラー:', data.error);
    }
  } catch (error) {
    console.error('通信エラー:', error);
  }
}

// ヒートマップデータ取得
async function getHeatmapData(range = 50) {
  const response = await fetch(`http://localhost:8080/api/v1/heatmap/${range}`);
  const data = await response.json();
  return data.data;
}
```

### Python (requests)
```python
import requests
import json

# APIベースURL
BASE_URL = "http://localhost:8080"

def get_recommendations(count=5, history=100):
    """推薦番号を取得"""
    url = f"{BASE_URL}/api/v1/recommendations"
    params = {"count": count, "history": history}
    
    response = requests.get(url, params=params)
    data = response.json()
    
    if data["success"]:
        return data["data"]["recommendations"]
    else:
        raise Exception(f"API Error: {data['error']}")

def get_heatmap_data(range_size=50):
    """ヒートマップデータを取得"""
    url = f"{BASE_URL}/api/v1/heatmap/{range_size}"
    response = requests.get(url)
    return response.json()["data"]

# 使用例
recommendations = get_recommendations(count=3, history=200)
for i, rec in enumerate(recommendations, 1):
    print(f"推薦{i}: {rec['numbers']} (信頼度: {rec['confidence']})")
```

## 📝 まとめ

このプロジェクトは、ロト7の過去データを基にした統計的推薦システムとして、以下の価値を提供します：

1. **データ駆動**: 過去の結果データに基づく客観的分析
2. **拡張性**: Clean Architectureによる保守性の高い設計
3. **使いやすさ**: シンプルなREST APIインターフェース
4. **信頼性**: 包括的なテストスイートによる品質保証

本システムは完全に機能する状態であり、フロントエンドアプリケーションとの統合や、さらなる機能拡張の基盤として活用できます。