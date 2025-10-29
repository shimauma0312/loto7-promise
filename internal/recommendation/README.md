過去のロト7抽選結果を統計分析して良い感じの組み合わせを出力する。

- **重みでランダム選択**: スコアが高い数字ほど選ばれやすいが、完全に固定はしない
- **範囲別選択**: 低い数字・中間・高い数字の各範囲からバランスを選択する。
- **ランダム性**: 時刻ベースシードを使ってある程度のランダム性を保つ

### コンポーネント

#### 1. RecommendationEngine
```go
type RecommendationEngine struct {
    config  RecommendationConfig  // 推薦設定
    results [][]string           // 過去の抽選結果
    history CombinationHistory   // 組み合わせ履歴
    rng     *rand.Rand          // ランダムジェネレータ
}
```

#### 2. NumberScore
```go
type NumberScore struct {
    Number           int     // 数字（1-37）
    FrequencyScore   float64 // 出現頻度スコア
    RecentScore      float64 // 最近の出現スコア
    ConsecutiveScore float64 // 連続出現スコア
    PositionScore    float64 // 位置スコア
    TotalScore       float64 // 総合スコア
}
```

## スコアリングについて

### 1. 出現頻度スコア (FrequencyScore)
過去に出現した回数で選択確率を評価する
- **計算方法**: `出現回数 / 全体の最大出現回数 * 100`
- **重み**: `FrequencyWeight` (デフォルト: 0.4)

### 2. 最近の出現スコア (RecentScore)
最近出現した数字を避ける傾向を作り出す
- **計算方法**: 直近の出現からの経過回数に基づく
- **重み**: `RecentWeight` (デフォルト: 0.3)

### 3. 連続出現スコア (ConsecutiveScore)
連続で出現している数字にボーナス/ペナルティを付与する
- **計算方法**: 連続出現回数で調整
- **重み**: `ConsecutiveWeight` (デフォルト: 0.2)

### 4. 位置スコア (PositionScore)
抽選時の位置による傾向を反映する
- **計算方法**: 各位置での出現傾向を分析
- **重み**: `PositionWeight` (デフォルト: 0.1)

### 総合スコア
```
TotalScore = FrequencyScore * 0.4 + RecentScore * 0.3 + 
             ConsecutiveScore * 0.2 + PositionScore * 0.1
```

## ランダム選択アルゴ

### 1. 範囲別分類
数字を3つの範囲に分類し、各範囲からバランス選択する：

- **低範囲**: 1-15（最低2個選択）
- **中範囲**: 16-25（最低2個選択）  
- **高範囲**: 26-37（最低2個選択）

### 2. 重み付きランダム選択
```go
func (re *RecommendationEngine) weightedRandomSelect(candidates []NumberScore, usedNumbers map[int]bool) int {
    // スコアが高いほど選ばれやすいが、完全に決定論的ではない
    totalWeight := Σ(candidate.TotalScore + 1.0)
    randomValue := re.rng.Float64() * totalWeight
    
    // 累積重みで選択
    currentWeight := 0.0
    for each candidate {
        currentWeight += (candidate.TotalScore + 1.0)
        if randomValue <= currentWeight {
            return candidate.Number
        }
    }
}
```

### 3. 制約チェック

#### 連続数字制限
- 連続した数字のペアが2つ以上にならないよう制限
- 例: `[1,2,3,4,5,6,7]` のようなアホみたいな組み合わせは回避する

#### 重複度チェック
- 既存の推薦との重複が4個以上にならないよう制限する
- 各推薦がユニークな特徴を持たせる

## 🔧 設定パラメータ

```go
type RecommendationConfig struct {
    MaxRecommendations int     // 最大推薦数（デフォルト: 5）
    RecentAvoidCount   int     // 直近回避回数（デフォルト: 3）
    ConsecutiveBoost   int     // 連続出現ブースト判定回数（デフォルト: 3）
    HistoryLookback    int     // 過去履歴確認回数（デフォルト: 100）
    FrequencyWeight    float64 // 出現頻度重み（デフォルト: 0.4）
    RecentWeight       float64 // 最近出現重み（デフォルト: 0.3）
    ConsecutiveWeight  float64 // 連続出現重み（デフォルト: 0.2）
    PositionWeight     float64 // 位置重み（デフォルト: 0.1）
}
```

## 使用例

### 基本的な使用方法
```go
// 推薦エンジン作成
engine := recommendation.NewRecommendationEngine(recommendation.DefaultConfig())

// データ読み込み
err := engine.LoadData(100) // 過去100回分
if err != nil {
    log.Fatal(err)
}

// 推薦生成
recommendations, err := engine.GenerateRecommendations()
if err != nil {
    log.Fatal(err)
}

// 結果表示
for i, rec := range recommendations {
    fmt.Printf("推薦 %d: %v\n", i+1, rec)
}
```

### 実行結果例
```
推薦 1: 04 - 11 - 16 - 18 - 31 - 32 - 33
推薦 2: 01 - 02 - 16 - 17 - 18 - 31 - 35
推薦 3: 05 - 12 - 21 - 24 - 25 - 26 - 32
推薦 4: 01 - 09 - 10 - 21 - 25 - 27 - 32
推薦 5: 02 - 04 - 10 - 21 - 25 - 36 - 37
```
## メモ

### 重要な設計判断
1. **時刻ベースシード**: `time.Now().UnixNano()`でランダム性を確保する
2. **範囲別選択**: ロト7の統計的傾向を反映した分散選択する
3. **重み付き選択**: 完全ランダムと決定論のバランスを維持する

### 既知の制限事項
1. Go 1.23以上でのコンパイルが必要
2. 過去データが少ない場合の精度低下
3. 極端な設定値での動作未検証
