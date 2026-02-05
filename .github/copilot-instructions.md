# GitHub Copilot Instructions - Loto7 Recommendation System

## Architecture Guidelines

### モジュール構成原則

**疎結合・疎依存**を重視したモジュラー設計を採用しています。

#### ディレクトリ構造
```
cmd/               # エントリーポイント（各コマンドは独立）
├── api/          # APIサーバー
├── heatmap/      # ヒートマップ生成
├── recommendation/ # 推薦エンジン
└── result/       # 結果取得

internal/         # 内部パッケージ（外部公開しない）
├── api/          # API応答処理
├── cache/        # キャッシュ管理
├── heatmap/      # ヒートマップロジック
├── recommendation/ # 推薦ロジック（モジュール分割推奨）
└── result/       # 結果処理ロジック
```

#### モジュール設計原則

1. **単一責任原則（SRP）**
   - 1つのパッケージ/ファイルは1つの責務のみを持つ
   - 例: `analyzer.go`は統計分析のみ、`validator.go`はバリデーションのみ

2. **依存関係の方向性**
   - 依存は常に一方向（循環依存を禁止）
   - `internal`パッケージは`cmd`に依存しない
   - 共通型定義は`types.go`に集約

3. **Interface駆動設計**
   - 具象型ではなくinterfaceに依存
   - モック化・テスタビリティを重視
   ```go
   // Good: interface経由の依存
   type Analyzer interface {
       Analyze(results [][]string) StatisticalAnalysis
   }
   
   // Bad: 具象型への直接依存
   type Engine struct {
       analyzer *DefaultAnalyzer // NG
   }
   ```

4. **モジュール構成（recommendation パッケージ）**
   ```
   internal/recommendation/
   ├── types.go          # 共通型定義（依存なし）
   ├── analyzer.go       # 統計分析（types.goに依存）
   ├── scorer.go         # スコアリング（types.goに依存）
   ├── selector.go       # 候補選択（scorer interfaceに依存）
   ├── validator.go      # バリデーション（独立）
   ├── formatter.go      # フォーマット（types.goに依存）
   └── recommendation.go # オーケストレーション（全interfaceに依存）
   ```

---

## Documentation Comments (DOCコメント)

### 必須ルール

1. **すべてのエクスポート関数・型にDOCコメントを記載**
   ```go
   // 数字の優先度を計算する
   func CalculatePriority(num int) float64 {
       // ...
   }
   ```

2. **コメントの形式**
   - 日本語で簡潔に（1行が理想、最大3行）
   - 「何をするか」を説明（「どのように」は不要）

3. **端的且つ正確に**
   ```go
   // Good: 端的で正確
   // 過去の抽選結果から統計分析を実行する
   func Analyze(results [][]string) StatisticalAnalysis
   
   // Bad: 冗長で不正確
   // Analyze はresultsを受け取って、いろいろな分析をして、
   // その結果をStatisticalAnalysisという構造体に詰めて返します
   func Analyze(results [][]string) StatisticalAnalysis
   ```

4. **不要なコメントは生成しない**
   ```go
   // Good: 自明なコードにはコメント不要
   sum := 0
   for _, num := range numbers {
       sum += num
   }
   
   // Bad: 自明なことを繰り返さない
   // sumに0を代入
   sum := 0
   // numbersをループ
   for _, num := range numbers {
       // sumにnumを足す
       sum += num
   }
   ```

5. **パッケージコメント**
   - 各パッケージの最初のファイル（通常は`doc.go`またはメインファイル）にパッケージの説明を記載
   ```go
   // ロト7の推薦番号生成機能を提供する
   //
   // 統計分析、スコアリング、バリデーションを組み合わせて
   // 最適な数字の組み合わせを生成する
   package recommendation
   ```

### 構造体・インターフェースのコメント

```go
// 過去データの統計分析を行うインターフェース
type Analyzer interface {
    // 指定された回数分の結果を分析する
    Analyze(results [][]string, lookback int) StatisticalAnalysis
}

// 統計分析の結果を保持する
type StatisticalAnalysis struct {
    FrequentNumbers   []int       // 出現回数が多い数字
    HotNumbers        []int       // 直近で波が来ている数字
    RevivalCandidates []int       // 復活候補（20回以上未出現）
    FrequencyMap      map[int]int // 各数字の出現回数
}
```

---

## Code Quality Standards

### 1. エラーハンドリング

```go
// Good: エラーを適切に処理
func LoadData(count int) error {
    results := result.GetResult(count)
    if len(results) == 0 {
        return fmt.Errorf("データを取得できませんでした: count=%d", count)
    }
    return nil
}

// Bad: エラーを無視
func LoadData(count int) {
    results := result.GetResult(count)
    // len(results)のチェックなし
}
```

### 2. 命名規則

- **関数名**: 動詞で始める（`Calculate`, `Validate`, `Generate`）
- **型名**: 名詞（`Analyzer`, `Validator`, `Engine`）
- **変数名**: 
  - 短いスコープ: 短い名前（`i`, `num`, `cfg`）
  - 長いスコープ: 説明的な名前（`recommendationCount`, `maxAttempts`）
- **定数**: 大文字スネークケース（`MAX_RETRY_COUNT`）または PascalCase（`MaxRetryCount`）

```go
// Good
func CalculatePriority(num int, stats StatisticalAnalysis) float64 {
    for i, hot := range stats.HotNumbers {
        // ...
    }
}

// Bad: 曖昧な命名
func calc(n int, s StatisticalAnalysis) float64 {
    for index, hotNumberValue := range s.HotNumbers {
        // ...
    }
}
```

### 3. 関数の長さ

- 1つの関数は**50行以内**を目標
- 50行を超える場合は、サブ関数に分割
- 1つの関数は1つの抽象レベルで記述

```go
// Good: 適切な抽象化
func GenerateRecommendations() ([][]int, error) {
    var recommendations [][]int
    
    for rec := 0; rec < maxCount; rec++ {
        combination := generateCombination()
        if isValid(combination) {
            recommendations = append(recommendations, combination)
        }
    }
    
    return recommendations, nil
}

// Bad: 詳細が混在（抽象レベルが混在）
func GenerateRecommendations() ([][]int, error) {
    var recommendations [][]int
    
    for rec := 0; rec < maxCount; rec++ {
        // ゾーン1から選択（詳細すぎる）
        zone1Count := 3 + rand.Intn(2)
        for i := 1; i <= 13; i++ {
            // ...
        }
        // ...
    }
}
```

### 4. テスタビリティ

- グローバル変数を避ける
- 依存は引数で注入（Dependency Injection）
- 外部システムとの接続はinterfaceで抽象化

```go
// Good: テスト可能
type Engine struct {
    analyzer Analyzer  // interface
    rng      *rand.Rand
}

func NewEngine(analyzer Analyzer, seed int64) *Engine {
    return &Engine{
        analyzer: analyzer,
        rng:      rand.New(rand.NewSource(seed)),
    }
}

// Bad: グローバル変数に依存
var globalAnalyzer *DefaultAnalyzer

func NewEngine() *Engine {
    return &Engine{
        analyzer: globalAnalyzer, // テスト時に差し替え困難
    }
}
```

### 5. パフォーマンス考慮

- 不要なメモリアロケーションを避ける
- `make`で容量を事前確保
- ループ内での重複計算を避ける

```go
// Good: 事前確保
candidates := make([]int, 0, 37)
for num := 1; num <= 37; num++ {
    candidates = append(candidates, num)
}

// Bad: 再アロケーション発生
var candidates []int
for num := 1; num <= 37; num++ {
    candidates = append(candidates, num) // 容量不足で再確保
}
```

---

## Go-Specific Best Practices

### 1. ゼロ値の活用

```go
// Good: ゼロ値を活用
type Config struct {
    MaxCount int     // 0がデフォルト
    Enabled  bool    // falseがデフォルト
}

// Bad: 不要な初期化
config := Config{
    MaxCount: 0,
    Enabled:  false,
}
```

### 2. エラー型のカスタマイズ

```go
// Good: 独自エラー型
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// 使用例
if !isValid(zone1Count) {
    return nil, &ValidationError{
        Field:   "zone1Count",
        Message: "must be between 3 and 4",
    }
}
```

### 3. context.Contextの使用

長時間処理やAPIハンドラーでは`context.Context`を使用

```go
// Good: contextを受け取る
func Analyze(ctx context.Context, results [][]string) (StatisticalAnalysis, error) {
    select {
    case <-ctx.Done():
        return StatisticalAnalysis{}, ctx.Err()
    default:
        // 分析処理
    }
}
```

---

## Prohibited Patterns（禁止パターン）

### ❌ 避けるべきパターン

1. **panicの乱用**
   ```go
   // Bad
   if err != nil {
       panic(err) // 回復不可能なエラー以外でpanicを使わない
   }
   
   // Good
   if err != nil {
       return nil, fmt.Errorf("処理に失敗: %w", err)
   }
   ```

2. **マジックナンバー**
   ```go
   // Bad
   if count > 100 {  // 100は何？
   
   // Good
   const MaxHistoryLookback = 100
   if count > MaxHistoryLookback {
   ```

3. **肥大化した構造体**
   ```go
   // Bad: 責務が多すぎる
   type Engine struct {
       results [][]string
       stats   StatisticalAnalysis
       scores  []NumberScore
       validators []Validator
       formatters []Formatter
       cache  map[string]interface{}
       config Config
       // ...20個以上のフィールド
   }
   
   // Good: 分割
   type Engine struct {
       analyzer  Analyzer
       selector  Selector
       validator Validator
   }
   ```

4. **深いネスト**
   ```go
   // Bad: 4段以上のネスト
   if condition1 {
       if condition2 {
           if condition3 {
               if condition4 {
                   // ...
               }
           }
       }
   }
   
   // Good: 早期return
   if !condition1 {
       return
   }
   if !condition2 {
       return
   }
   // フラットな構造
   ```

---

## Testing Guidelines

### テスト命名規則

```go
// TestFunction_Scenario_ExpectedResult パターン
func TestCalculatePriority_HotNumber_ReturnsHighScore(t *testing.T) {
    // ...
}

func TestValidateZoneDistribution_EmptyZone_ReturnsFalse(t *testing.T) {
    // ...
}
```

### テーブル駆動テスト

```go
func TestValidateTotalSum(t *testing.T) {
    tests := []struct {
        name        string
        combination []int
        want        bool
    }{
        {
            name:        "合計が範囲内",
            combination: []int{1, 2, 15, 20, 25, 30, 37},
            want:        true,
        },
        {
            name:        "合計が小さすぎる",
            combination: []int{1, 2, 3, 4, 5, 6, 7},
            want:        false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := validateTotalSum(tt.combination)
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

## Version Control

### コミットメッセージ

```
<type>: <subject>

<body>

例:
feat: recommendation パッケージにスコアリングモジュールを追加

- WeightedScorer 構造体を実装
- 波、頻度、復活候補の優先度計算ロジックを分離
- ユニットテストを追加

fix: ゾーン分布のバリデーションロジックを修正

zone3の上限を3から4に変更し、より柔軟な組み合わせを許容
```

**Type一覧:**
- `feat`: 新機能
- `fix`: バグ修正
- `refactor`: リファクタリング
- `docs`: ドキュメント
- `test`: テスト追加・修正
- `chore`: ビルド、ツール設定等

---

## Summary

✅ **必ず守ること:**
1. モジュールは単一責任で設計
2. すべてのエクスポート関数・型にDOCコメント
3. Interface駆動設計でテスタビリティを確保
4. エラーは適切にハンドリング
5. 関数は50行以内を目標

⚠️ **避けること:**
1. 循環依存
2. マジックナンバー
3. グローバル変数
4. 不要な（自明な）コメント
5. 深いネスト（4段以上）

このガイドラインに従うことで、保守性が高く、テストしやすい、品質の高いコードを維持できます。
