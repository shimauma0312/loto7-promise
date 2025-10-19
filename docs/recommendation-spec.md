# ロト7推薦機能仕様書

ロト7の過去の抽選結果を分析し、推薦番号を生成する

## 機能説明

### 1. 基本コンセプト

### 2. 分析要素

#### 2.1 出現頻度分析 (FrequencyScore)
- **基準**: 過去100回以内の出現回数
- **重み**: 30%
- **ロジック**: 
  - 過去100回以内に出現していない数字は候補から除外
  - 出現回数が多い数字ほど高スコア

#### 2.2 最近の出現分析 (RecentScore)  
- **基準**: 直近の出現パターン
- **重み**: 40%
- **ロジック**:
  - 直前数回（デフォルト3回）で出た数字は大幅減点 (-1.0)
  - 過去50回以内に出現した数字は加点 (+0.5)
  - それ以外は中立 (0.0)

#### 2.3 連続出現分析 (ConsecutiveScore)
- **基準**: 直近での連続出現回数
- **重み**: 20%
- **ロジック**:
  - 3回以上連続で出現している数字はブーストされる
  - 連続回数 × 0.2 のスコア加算

#### 2.4 位置による分析 (PositionScore)
- **基準**: 最初の数字による全体傾向
- **重み**: 10%
- **ロジック**:
  - 前回の最初の数字が10以下の場合、1-20の数字に加点、21-37に減点
  - それ以外は中立

### 3. 組み合わせ制約

#### 3.1 ペア回避
- 過去に2回以上出現した数字ペアは避ける
- 例: 3と9が過去に組み合わされた場合、再度の組み合わせを回避

#### 3.2 重複度制限
- 複数の推薦組み合わせ間で4個以上の数字が重複することを避ける
- 異なる推薦組み合わせを保証

### 4. 設定パラメータ

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

## API仕様

### RecommendationEngine

#### NewRecommendationEngine(config RecommendationConfig) *RecommendationEngine
推薦エンジンを作成します。

#### LoadData(count int) error
過去の抽選データを読み込みます。
- `count`: 分析する過去の抽選回数

#### GenerateRecommendations() ([][]int, error)
推薦組み合わせを生成します。
- 戻り値: 推薦番号の配列（各組み合わせは7個の数字）

#### FormatRecommendations(recommendations [][]int) string
推薦結果を見やすい形式でフォーマットします。

## 使用例

### Go コード例

```go
package main

import (
    "fmt"
    "github.com/shimauma0312/loto7-promise/internal/recommendation"
)

func main() {
    // 設定を作成
    config := recommendation.DefaultConfig()
    config.MaxRecommendations = 3
    
    // エンジンを作成
    engine := recommendation.NewRecommendationEngine(config)
    
    // データを読み込み（過去100回分）
    err := engine.LoadData(100)
    if err != nil {
        panic(err)
    }
    
    // 推薦を生成
    recommendations, err := engine.GenerateRecommendations()
    if err != nil {
        panic(err)
    }
    
    // 結果を表示
    result := engine.FormatRecommendations(recommendations)
    fmt.Print(result)
}
```

### CLIツール使用例

```bash
# 基本的な使用
./recommendation

# 過去150回分を分析し、3つの推薦を生成
./recommendation -count 150 -max 3

# 直近5回を避け、連続4回でブースト
./recommendation -avoid 5 -boost 4

# スコア詳細表示
./recommendation -scores
```

## 制限事項

1. **データ依存性**: キャッシュされた過去の抽選結果に依存
2. **確率的手法ではない**: 統計的確率ではなく、経験則に基づく
3. **保証なし**: 当選を保証するものではない

## 将来の拡張

1. **機械学習モデル**: より高度な予測モデルの導入
2. **パターン分析**: より複雑な出現パターンの分析
3. **外部要因**: 季節性や曜日効果の考慮
4. **ユーザー設定**: より細かな重み調整機能

## テスト

推薦機能には包括的なユニットテストが含まれています：

```bash
# テスト実行
go test ./internal/recommendation/

# ベンチマーク実行  
go test -bench=. ./internal/recommendation/
```

## パフォーマンス

- **メモリ使用量**: 過去の抽選データサイズに依存（通常数MB以下）
- **処理時間**: 100回分の分析で通常1秒以下
- **スケーラビリティ**: 1000回分までの分析で良好なパフォーマンス