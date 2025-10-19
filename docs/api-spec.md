# ロト7 Promise API仕様書

ロト7の抽選結果分析、ヒートマップ生成、推薦番号生成機能APIです。

## ベースURL

```
http://localhost:8080
```

## 共通レスポンス形式

すべてのAPIエンドポイントは統一されたレスポンス形式を返します：

```json
{
  "success": true,
  "message": "処理成功メッセージ",
  "data": { /* 実際のデータ */ },
  "error": null,
  "timestamp": "2025-10-19T21:30:00Z",
  "version": "1.0.0"
}
```

### エラーレスポンス例

```json
{
  "success": false,
  "message": "エラーが発生しました",
  "data": null,
  "error": {
    "code": "INVALID_PARAMETER",
    "message": "パラメータが不正です",
    "details": "count parameter must be between 1 and 1000"
  },
  "timestamp": "2025-10-19T21:30:00Z",
  "version": "1.0.0"
}
```

## エンドポイント一覧

### 1. ヘルスチェック

**GET** `/health`

サーバーの稼働状況を確認します。

#### レスポンス例

```json
{
  "success": true,
  "message": "サービスは正常に動作しています",
  "data": {
    "status": "healthy",
    "uptime": "2h30m15s",
    "timestamp": "2025-10-19T21:30:00Z",
    "services": {
      "result_cache": "operational",
      "heatmap_engine": "operational",
      "recommendation": "operational"
    }
  }
}
```

### 2. 抽選結果取得

#### デフォルト取得（10回分）

**GET** `/api/v1/results`

#### 指定回数取得

**GET** `/api/v1/results/{count}`

**パラメータ:**
- `count` (1-1000): 取得する抽選回数

#### レスポンス例

```json
{
  "success": true,
  "message": "結果を正常に取得しました",
  "data": {
    "results": [
      {
        "draw_number": 642,
        "numbers": ["01", "07", "22", "23", "33", "34", "35"],
        "date": "2025-10-19",
        "formatted_numbers": "01-07-22-23-33-34-35"
      }
    ],
    "total_count": 1,
    "requested_count": 10,
    "data_source": "cache"
  }
}
```

### 3. ヒートマップ分析

#### デフォルトヒートマップ（50回分）

**GET** `/api/v1/heatmap`

#### 指定範囲ヒートマップ

**GET** `/api/v1/heatmap/{range}`

**パラメータ:**
- `range` (1-1000): 分析対象の抽選回数

#### 指定数字の分析

**GET** `/api/v1/heatmap/{range}/number/{number}`

**パラメータ:**
- `range` (1-1000): 分析対象の抽選回数
- `number` (1-37): 分析対象の数字

#### 指定位置の分析

**GET** `/api/v1/heatmap/{range}/position/{position}`

**パラメータ:**
- `range` (1-1000): 分析対象の抽選回数  
- `position` (1-7): 分析対象の位置

#### レスポンス例

```json
{
  "success": true,
  "message": "ヒートマップを正常に生成しました",
  "data": {
    "search_range": 50,
    "actual_draws": 45,
    "hottest_numbers": [7, 12, 23],
    "coldest_numbers": [3, 15, 29],
    "most_biased_number": 7,
    "statistics": {
      "total_numbers": 37,
      "average_frequency": 8.5,
      "max_frequency": 15,
      "min_frequency": 3
    }
  }
}
```

### 4. 推薦番号生成

#### デフォルト推薦（GET）

**GET** `/api/v1/recommendations`

**クエリパラメータ:**
- `count` (1-10): 生成する推薦数 (デフォルト: 5)
- `history` (1-1000): 分析対象回数 (デフォルト: 100)
- `avoid` (0-10): 直近回避回数 (デフォルト: 3)

#### カスタム設定推薦（POST）

**POST** `/api/v1/recommendations`

**リクエストボディ例:**

```json
{
  "max_recommendations": 3,
  "recent_avoid_count": 5,
  "consecutive_boost": 3,
  "history_lookback": 150,
  "frequency_weight": 0.3,
  "recent_weight": 0.4,
  "consecutive_weight": 0.2,
  "position_weight": 0.1
}
```

#### レスポンス例

```json
{
  "success": true,
  "message": "推薦を正常に生成しました",
  "data": {
    "recommendations": [
      {
        "id": 1,
        "numbers": [4, 8, 17, 18, 19, 20, 21]
      },
      {
        "id": 2,
        "numbers": [4, 5, 6, 18, 19, 20, 24]
      }
    ],
    "analysis_info": {
      "data_source": "cache",
      "analyzed_draws": 100,
      "generated_at": "2025-10-19T21:30:00Z",
      "algorithm": "human-intuition-based"
    },
    "config": {
      "max_recommendations": 5,
      "recent_avoid_count": 3,
      "consecutive_boost": 3,
      "history_lookback": 100,
      "frequency_weight": 0.3,
      "recent_weight": 0.4,
      "consecutive_weight": 0.2,
      "position_weight": 0.1
    }
  }
}
```

## エラーコード一覧

| エラーコード | 説明 |
|-------------|------|
| `INVALID_COUNT` | count パラメータが不正 |
| `INVALID_RANGE` | range パラメータが不正 |
| `INVALID_NUMBER` | number パラメータが不正 |
| `INVALID_POSITION` | position パラメータが不正 |
| `DATA_FETCH_ERROR` | データ取得エラー |
| `HEATMAP_GENERATION_ERROR` | ヒートマップ生成エラー |
| `DATA_RETRIEVAL_ERROR` | データ取得エラー |
| `DATA_LOAD_ERROR` | データ読み込みエラー |
| `RECOMMENDATION_ERROR` | 推薦生成エラー |

## 使用例

### curl使用例

```bash
# ヘルスチェック
curl http://localhost:8080/health

# 過去10回分の結果取得
curl http://localhost:8080/api/v1/results

# 過去5回分の結果取得
curl http://localhost:8080/api/v1/results/5

# ヒートマップ取得（過去30回分）
curl http://localhost:8080/api/v1/heatmap/30

# 数字7の詳細分析（過去50回分）
curl http://localhost:8080/api/v1/heatmap/50/number/7

# デフォルト推薦取得
curl http://localhost:8080/api/v1/recommendations

# カスタム推薦取得（3つ、過去80回分）
curl "http://localhost:8080/api/v1/recommendations?count=3&history=80"

# POST推薦（カスタム設定）
curl -X POST http://localhost:8080/api/v1/recommendations \
  -H "Content-Type: application/json" \
  -d '{"max_recommendations": 3, "history_lookback": 150}'
```

### JavaScript使用例

```javascript
// 推薦取得
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

// ヒートマップ取得
async function getHeatmap(range = 50) {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/heatmap/${range}`);
    const data = await response.json();
    
    if (data.success) {
      console.log('ヒートマップデータ:', data.data);
    }
  } catch (error) {
    console.error('エラー:', error);
  }
}
```

## レート制限

現在のバージョンではレート制限は実装されていませんが、将来的には以下の制限を予定しています：

- 1分間に60リクエスト
- 推薦生成API: 1分間に10リクエスト

## CORS対応

すべてのオリジンからのアクセスを許可しています。本番環境では適切に制限してください。

## バージョニング

現在のAPIバージョンは `v1` です。将来的な変更に備えてバージョン管理を行っています。