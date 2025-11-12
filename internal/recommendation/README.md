
### 1. 分析フェーズ

#### 出現頻度分析
- 過去100回の各数字の出現回数集計
- 出現回数が多い数字10リスト
- 出現回数が少ない数字10リスト

#### 波の検出
- 直近10回以内で3回以上出現した数字を波判定
- これらの数字は選択確率が高くなる

#### 復活候補の識別
- 20回以上出現していない数字を復活候補として抽出するが、ただし、全体で出現回数が極端に少ない数字は除外する

### 2. 推薦生成フェーズ

#### ゾーン別選択
```
ゾーン1 (1-13):  3-4個選択 (安定構成)
ゾーン2 (14-26): 残りを調整
ゾーン3 (27-37): 2-3個選択 (まれに3個)
```

#### 数字の優先度
- **+2.0点**: 波が来ている
- **+1.0点**: 出現頻度が高い
- **+1.0点**: 復活候補
- **-0.5点**: 出現回数が極端に少ない
- **-3.0点**: 直近3回以内に出現してる

#### 制約バリデーション
1. **ゾーン分布**: 各ゾーンが空でなく、推奨範囲内
2. **奇数偶数比率**: 奇4:偶3 または 奇3:偶4
3. **合計値**: 100 ≤ 合計 ≤ 160
4. **近距離ペア**: 差が1~3の数字ペアが最低1組存在

### コンポーネント

#### 1. RecommendationEngine
```go
type RecommendationEngine struct {
    config     RecommendationConfig  // 推薦設定
    results    [][]string           // 過去の抽選結果
    history    CombinationHistory   // 組み合わせ履歴
    rng        *rand.Rand          // ランダム
    statistics StatisticalAnalysis  // 統計分析結果
}
```

#### 2. StatisticalAnalysis
```go
type StatisticalAnalysis struct {
    FrequentNumbers    []int        // 出現回数が多い数字
    RareNumbers        []int        // 出現回数が少ない数字
    HotNumbers         []int        // 波が来ている数字
    RevivalCandidates  []int        // 復活候補
    FrequencyMap       map[int]int  // 各数字の出現回数
    LastAppearance     map[int]int  // 各数字の最終出現位置
}
```

#### 3. NumberScore
```go
type NumberScore struct {
    Number           int     // 数字（1-37）
    FrequencyScore   float64 // 出現頻度スコア
    RecentScore      float64 // 最近の出現スコア
    ConsecutiveScore float64 // 連続出現スコア
    PositionScore    float64 // 位置スコア
    TotalScore       float64 // 総合スコア
    IsRecentlyDrawn  bool    // 直近で出現したか
}
```
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
    FrequencyWeight    float64 // 出現頻度重み（デフォルト: 0.3）
    RecentWeight       float64 // 最近出現重み（デフォルト: 0.4）
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
formattedResult := engine.FormatRecommendations(recommendations)
fmt.Println(formattedResult)
```

### コマンドラインから実行
```bash
# 推薦番号を生成
make recommend

# または直接実行
go run ./cmd/recommendation
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
