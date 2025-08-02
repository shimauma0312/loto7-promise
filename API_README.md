# ロト7予想API


## クイックスタート

### Docker Composeを使用した起動

```bash
# API + nginx構成で起動
docker-compose up -d api nginx

# 開発環境で起動
docker-compose up -d dev
docker-compose exec dev sh
```

### 直接ビルドして起動

```bash
# 依存関係をインストール
go mod tidy

# APIサーバーを起動
go run ./cmd/api
```

## API仕様

ベースURL: `http://localhost/api/v1` (nginx経由) または `http://localhost:8080/api/v1` (直接)

### エンドポイント一覧

#### 1. 単一予想生成

**GET /api/v1/prediction**

クエリパラメータで予想を生成します。

```bash
# 基本的な使用例
curl "http://localhost/api/v1/prediction"

# パラメータ指定例
curl "http://localhost/api/v1/prediction?analysis_range=50&use_trend=true"
```

**パラメータ:**
- `analysis_range` (int): 分析対象の回数 (デフォルト: 30)
- `use_trend` (bool): 傾向を使用するか (デフォルト: true)

**POST /api/v1/prediction**

JSONボディで予想を生成します。

```bash
curl -X POST http://localhost/api/v1/prediction \
  -H "Content-Type: application/json" \
  -d '{
    "analysis_range": 50,
    "use_trend": true
  }'
```

#### 2. 複数予想生成

**GET /api/v1/predictions**

```bash
# 5組の予想を生成
curl "http://localhost/api/v1/predictions?count=5"

# 詳細パラメータ指定
curl "http://localhost/api/v1/predictions?count=10&analysis_range=100&use_trend=false"
```

**パラメータ:**
- `count` (int): 生成する予想数 (最大20、デフォルト: 5)
- `analysis_range` (int): 分析対象の回数 (デフォルト: 30)
- `use_trend` (bool): 傾向を使用するか (デフォルト: true)

**POST /api/v1/predictions**

```bash
curl -X POST http://localhost/api/v1/predictions \
  -H "Content-Type: application/json" \
  -d '{
    "count": 10,
    "analysis_range": 100,
    "use_trend": false
  }'
```

#### 3. その他

**GET /health** - ヘルスチェック
**GET /api/v1/info** - API仕様情報

## レスポンス形式

すべてのAPIは以下の形式でレスポンスを返します：

```json
{
  "success": true,
  "message": "予想番号を正常に生成しました",
  "data": {
    "main_numbers": [3, 12, 18, 25, 31, 34, 37],
    "bonus_numbers": [7, 15],
    "trend_based": true,
    "generated_time": "2025-08-02T10:30:45Z"
  }
}
```

**エラー時:**
```json
{
  "success": false,
  "message": "予想番号の生成に失敗しました",
  "error": "詳細なエラーメッセージ"
}
```

## 設定

### 環境変数

- `PORT`: APIサーバーのポート番号 (デフォルト: 8080)
- `GIN_MODE`: Ginの動作モード (release/debug)

### Nginx設定

`nginx.conf` を参考にして、本番環境に合わせてドメイン名やSSL設定を調整してください。

## 予想アルゴリズム

### 傾向ベース予想 (use_trend: true)

1. 指定された回数分の過去データから各数字の出現頻度を計算
2. 直近データの重みを30%、全体データの重みを70%で組み合わせ
3. 重み付き確率に基づいて数字を選択

### ランダム予想 (use_trend: false)

- 1〜37の数字から完全ランダムに選択
- 重複なしで本数字7個、ボーナス数字2個を生成

## 開発

### ローカル開発

```bash
# 開発用コンテナで作業
docker-compose up -d dev
docker-compose exec dev sh

# コンテナ内で
go run ./cmd/api
```

### テスト実行

```bash
# 全テスト実行
go test ./...

# 特定パッケージのテスト
go test ./internal/randomPrediction/
```

### ビルド

```bash
# 本番用バイナリビルド
go build -o api ./cmd/api

# Docker イメージビルド
docker build -t loto7-api .
```

## ライセンス

このプロジェクトは [MIT License](LICENSE) の下で公開されています。

## コントリビューション

1. このリポジトリをフォーク
2. フィーチャーブランチを作成 (`git checkout -b feature/AmazingFeature`)
3. 変更をコミット (`git commit -m 'Add some AmazingFeature'`)
4. ブランチにプッシュ (`git push origin feature/AmazingFeature`)
5. プルリクエストを作成

## API使用例

### JavaScript (fetch)

```javascript
// 単一予想取得
const response = await fetch('/api/v1/prediction?use_trend=true&analysis_range=50');
const data = await response.json();
console.log(data.data.main_numbers);

// 複数予想取得
const multiResponse = await fetch('/api/v1/predictions', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    count: 10,
    use_trend: true,
    analysis_range: 30
  })
});
const multiData = await multiResponse.json();
```

### Python (requests)

```python
import requests

# 単一予想
response = requests.get('http://localhost/api/v1/prediction', 
                       params={'use_trend': True, 'analysis_range': 50})
data = response.json()
print(data['data']['main_numbers'])

# 複数予想
response = requests.post('http://localhost/api/v1/predictions',
                        json={'count': 5, 'use_trend': True})
data = response.json()
for prediction in data['data']:
    print(prediction['main_numbers'])
```

### PHP (cURL)

```php
<?php
// 単一予想
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, 'http://localhost/api/v1/prediction?use_trend=true');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
$response = curl_exec($ch);
$data = json_decode($response, true);
echo implode(', ', $data['data']['main_numbers']);
curl_close($ch);
?>
