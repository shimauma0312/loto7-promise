# ロト7推薦機能仕様書

ロト7の過去の抽選結果を分析し、推薦番号を生成する

## 機能説明

### 1. 基本コンセプト

過去の抽選結果を統計分析し、各数字の出現パターンを評価してスコアリングを行う。  
スコアが高い数字を優先的に選択し、複数の制約条件を満たす推薦組み合わせを生成する。

### 2. 分析要素

#### 2.1 直近1回の抽選結果の除外 (Recent1Penalty)
- **基準**: 直近1回の抽選結果に出現した数字
- **重み**: -100.0（限りなく低い）
- **ロジック**: 
  - 直近1回で出現した数字は実質的に選択されない

#### 2.2 直近2回までの抽選結果の減点 (Recent2Penalty)
- **基準**: 直近2回までの抽選結果に出現した数字
- **重み**: -10.0（それなりに低い）
- **ロジック**:
  - 直近2回目に出現した数字は大幅に減点される

#### 2.3 50回以内の出現ブースト (Within50Boost)
- **基準**: 直近3回を除き、50回以内に出現した数字
- **重み**: 出現回数 × 0.5
- **ロジック**:
  - 直近3回は除外し、4回目から50回目までの範囲で出現した数字をブースト
  - 出現回数が多いほど高スコア

#### 2.4 低頻度数字の減点 (LowFrequencyPenalty)
- **基準**: 直近100回で出現回数が5回未満の数字
- **重み**: -1.0
- **ロジック**:
  - 出現回数が少ない数字は減点される

#### 2.5 未出現数字の除外 (NoAppearancePenalty)
- **基準**: 直近100回で出現していない数字
- **重み**: -1000.0（出現率0）
- **ロジック**:
  - 出現していない数字は実質的に選択されない

### 3. 組み合わせ制約

#### 3.1 ゾーン分布
- ゾーン1 (1-13): 3-4個
- ゾーン2 (14-26): 残り
- ゾーン3 (27-37): 2-3個

#### 3.2 奇数偶数比率
- 奇4:偶3 または 奇3:偶4

#### 3.3 合計値
- 100-160の範囲内

#### 3.4 平均値
- 16.0-23.0の範囲内

#### 3.5 滑らかさ
- 隣接数字間の差分の分散が50.0以下

#### 3.6 近距離ペア
- 差1-3のペアが少なくとも1組存在

#### 3.7 過去の組み合わせチェック
- 推薦番号が出たとき、同じ番号が出た過去の抽選結果を直近100回分まで洗い出す
- 数字の組み合わせが当時と完全に同じ場合は却下
- 同じ組み合わせの再発を防ぐ

#### 3.8 重複度制限
- 複数の推薦組み合わせ間で5個以上の数字が重複することを避ける

### 4. 設定パラメータ

```go
type RecommendationConfig struct {
    MaxRecommendations int     // 最大推薦数（デフォルト: 5）
    HistoryLookback    int     // 過去履歴確認回数（デフォルト: 100）
    FrequencyWeight    float64 // 出現頻度重み（デフォルト: 0.5）
    RecentWeight       float64 // 最近出現重み（デフォルト: 0.5）
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

# スコア詳細表示
./recommendation -scores
```
## テスト
```bash
# テスト実行
go test ./internal/recommendation/

# ベンチマーク実行  
go test -bench=. ./internal/recommendation/
```
