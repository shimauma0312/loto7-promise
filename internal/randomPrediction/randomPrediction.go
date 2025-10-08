package randomPrediction

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/shimauma0312/loto7-promise/internal/numberCount"
)

// PredictionResult 予想番号の結果を格納する構造体
type PredictionResult struct {
	MainNumbers   []int `json:"main_numbers"`   // 本数字7個
	BonusNumbers  []int `json:"bonus_numbers"`  // ボーナス数字2個
	TrendBased    bool  `json:"trend_based"`    // 傾向ベースの予想かどうか
	GeneratedTime time.Time `json:"generated_time"` // 生成時刻
	AnalysisInfo  *AnalysisInfo `json:"analysis_info,omitempty"` // 分析情報（デバッグ用）
}

// WeightedNumber 重み付き数字を表す構造体
type WeightedNumber struct {
	Number int     // 数字
	Weight float64 // 重み（出現頻度等から算出）
	Reason string  // 重み付けの理由（デバッグ用）
}

// AnalysisInfo 統計分析の詳細情報を格納する構造体
type AnalysisInfo struct {
	FrequencyWeights    map[int]float64 `json:"frequency_weights"`    // 出現頻度重み
	ConsecutiveWeights  map[int]float64 `json:"consecutive_weights"`  // 連続性重み
	IntervalWeights     map[int]float64 `json:"interval_weights"`     // 間隔重み
	OddEvenBalance      float64         `json:"odd_even_balance"`     // 奇数偶数バランス
	RangeBalance        float64         `json:"range_balance"`        // 前半後半バランス
	PeriodicityWeights  map[int]float64 `json:"periodicity_weights"`  // 周期性重み
	CorrelationWeights  map[int]float64 `json:"correlation_weights"`  // 相関重み
	TotalAnalysisTime   time.Duration   `json:"total_analysis_time"`  // 分析時間
}

// StatisticalPattern 統計パターンの設定
type StatisticalPattern struct {
	UseFrequency     bool    `json:"use_frequency"`     // 出現頻度
	UseConsecutive   bool    `json:"use_consecutive"`   // 連続性
	UseInterval      bool    `json:"use_interval"`      // 間隔分析
	UseOddEven       bool    `json:"use_odd_even"`      // 奇数偶数バランス
	UseRange         bool    `json:"use_range"`         // 範囲バランス
	UsePeriodicity   bool    `json:"use_periodicity"`   // 周期性
	UseCorrelation   bool    `json:"use_correlation"`   // 相関
	FrequencyWeight  float64 `json:"frequency_weight"`  // 出現頻度の重み係数
	ConsecutiveWeight float64 `json:"consecutive_weight"` // 連続性の重み係数
	IntervalWeight   float64 `json:"interval_weight"`   // 間隔の重み係数
	OddEvenWeight    float64 `json:"odd_even_weight"`   // 奇数偶数の重み係数
	RangeWeight      float64 `json:"range_weight"`      // 範囲の重み係数
	PeriodicityWeight float64 `json:"periodicity_weight"` // 周期性の重み係数
	CorrelationWeight float64 `json:"correlation_weight"` // 相関の重み係数
}

// ランダムジェネレータ（初期化時に一度だけシード設定）
var globalRand *rand.Rand

func init() {
	// 開始時に一度だけランダムシードを設定
	globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// GetDefaultStatisticalPattern デフォルトの統計パターン設定を返す
func GetDefaultStatisticalPattern() *StatisticalPattern {
	return &StatisticalPattern{
		UseFrequency:      true,
		UseConsecutive:    true,
		UseInterval:       true,
		UseOddEven:        true,
		UseRange:          true,
		UsePeriodicity:    true,
		UseCorrelation:    true,
		FrequencyWeight:   0.25, // 25% - 出現頻度
		ConsecutiveWeight: 0.15, // 15% - 連続性
		IntervalWeight:    0.15, // 15% - 間隔
		OddEvenWeight:     0.15, // 15% - 奇数偶数バランス
		RangeWeight:       0.15, // 15% - 範囲バランス
		PeriodicityWeight: 0.10, // 10% - 周期性
		CorrelationWeight: 0.05, // 5%  - 相関（最新データ依存のため軽め）
	}
}

// GenerateRandomPrediction 過去データの傾向を加味したランダム予想番号を生成
// 引数: analysisRange - 傾向分析に使用する過去データの範囲, useTrend - 傾向を考慮するかどうか
// 戻り値: PredictionResult構造体のポインタ、エラー
func GenerateRandomPrediction(analysisRange int, useTrend bool) (*PredictionResult, error) {
	return GenerateRandomPredictionWithPattern(analysisRange, useTrend, GetDefaultStatisticalPattern())
}

// GenerateRandomPredictionWithPattern 統計パターンを指定して予想番号を生成
// 引数: analysisRange - 分析範囲, useTrend - 傾向使用フラグ, pattern - 統計パターン設定
// 戻り値: PredictionResult構造体のポインタ、エラー
func GenerateRandomPredictionWithPattern(analysisRange int, useTrend bool, pattern *StatisticalPattern) (*PredictionResult, error) {
	startTime := time.Now()
	
	// 入力値の妥当性チェック
	if analysisRange <= 0 {
		return nil, fmt.Errorf("分析範囲は1以上で指定してください。入力値: %d", analysisRange)
	}

	if pattern == nil {
		pattern = GetDefaultStatisticalPattern()
	}

	var mainNumbers []int
	var bonusNumbers []int
	var analysisInfo *AnalysisInfo

	if useTrend {
		// 高度な統計分析による予想番号生成
		weightedNumbers, analysis, err := calculateAdvancedNumberWeights(analysisRange, pattern)
		if err != nil {
			return nil, fmt.Errorf("統計分析に失敗しました: %v", err)
		}

		analysisInfo = analysis
		
		// 本数字の選択
		mainNumbers = selectDiverseWeightedNumbers(weightedNumbers, 7, true)
		
		// ボーナス数字の選択（本数字と重複しないよう注意）
		bonusNumbers = selectBonusNumbers(weightedNumbers, mainNumbers, 2)
	} else {
		// ランダム生成
		mainNumbers = generateImprovedRandomNumbers(7)
		bonusNumbers = generateImprovedRandomNumbers(2)
	}

	// 数字をソート
	sort.Ints(mainNumbers)
	sort.Ints(bonusNumbers)

	if analysisInfo != nil {
		analysisInfo.TotalAnalysisTime = time.Since(startTime)
	}

	result := &PredictionResult{
		MainNumbers:   mainNumbers,
		BonusNumbers:  bonusNumbers,
		TrendBased:    useTrend,
		GeneratedTime: time.Now(),
		AnalysisInfo:  analysisInfo,
	}

	return result, nil
}

// calculateAdvancedNumberWeights 
// 引数: analysisRange - 分析範囲, pattern - 統計パターン設定
// 戻り値: 重み付き数字配列、分析情報、エラー
func calculateAdvancedNumberWeights(analysisRange int, pattern *StatisticalPattern) ([]WeightedNumber, *AnalysisInfo, error) {
	analysis := &AnalysisInfo{
		FrequencyWeights:   make(map[int]float64),
		ConsecutiveWeights: make(map[int]float64),
		IntervalWeights:    make(map[int]float64),
		PeriodicityWeights: make(map[int]float64),
		CorrelationWeights: make(map[int]float64),
	}
	
	var weightedNumbers []WeightedNumber

	// 各数字（1-37）に対して統計分析を実行
	for i := 1; i <= 37; i++ {
		var totalWeight float64 = 0.0
		var reasons []string

		// 1. 出現頻度分析
		if pattern.UseFrequency {
			frequencyWeight, err := calculateFrequencyWeight(i, analysisRange)
			if err != nil {
				fmt.Printf("警告: 数字%dの出現頻度取得に失敗: %v\n", i, err)
				frequencyWeight = 0.027 // 1/37の期待値
			}
			analysis.FrequencyWeights[i] = frequencyWeight
			totalWeight += frequencyWeight * pattern.FrequencyWeight
			reasons = append(reasons, fmt.Sprintf("頻度:%.3f", frequencyWeight))
		}

		// 2. 連続数字分析
		if pattern.UseConsecutive {
			consecutiveWeight := calculateConsecutiveWeight(i)
			analysis.ConsecutiveWeights[i] = consecutiveWeight
			totalWeight += consecutiveWeight * pattern.ConsecutiveWeight
			reasons = append(reasons, fmt.Sprintf("連続:%.3f", consecutiveWeight))
		}

		// 3. 数字間隔分析
		if pattern.UseInterval {
			intervalWeight := calculateIntervalWeight(i, analysisRange)
			analysis.IntervalWeights[i] = intervalWeight
			totalWeight += intervalWeight * pattern.IntervalWeight
			reasons = append(reasons, fmt.Sprintf("間隔:%.3f", intervalWeight))
		}

		// 4. 奇数偶数バランス重み
		if pattern.UseOddEven {
			oddEvenWeight := calculateOddEvenWeight(i)
			totalWeight += oddEvenWeight * pattern.OddEvenWeight
			reasons = append(reasons, fmt.Sprintf("奇偶:%.3f", oddEvenWeight))
		}

		// 5. 範囲バランス重み（前半1-18, 後半19-37）
		if pattern.UseRange {
			rangeWeight := calculateRangeWeight(i)
			totalWeight += rangeWeight * pattern.RangeWeight
			reasons = append(reasons, fmt.Sprintf("範囲:%.3f", rangeWeight))
		}

		// 6. 周期性分析
		if pattern.UsePeriodicity {
			periodicityWeight := calculatePeriodicityWeight(i, analysisRange)
			analysis.PeriodicityWeights[i] = periodicityWeight
			totalWeight += periodicityWeight * pattern.PeriodicityWeight
			reasons = append(reasons, fmt.Sprintf("周期:%.3f", periodicityWeight))
		}

		// 7. 前回との相関分析
		if pattern.UseCorrelation {
			correlationWeight := calculateCorrelationWeight(i, analysisRange)
			analysis.CorrelationWeights[i] = correlationWeight
			totalWeight += correlationWeight * pattern.CorrelationWeight
			reasons = append(reasons, fmt.Sprintf("相関:%.3f", correlationWeight))
		}

		// 最低重みを保証（完全に0にならないよう調整）
		if totalWeight < 0.001 {
			totalWeight = 0.001
		}

		weightedNumbers = append(weightedNumbers, WeightedNumber{
			Number: i,
			Weight: totalWeight,
			Reason: strings.Join(reasons, ", "),
		})
	}

	// 奇数偶数バランスの分析結果を保存
	analysis.OddEvenBalance = calculateOverallOddEvenBalance(analysisRange)
	analysis.RangeBalance = calculateOverallRangeBalance(analysisRange)

	return weightedNumbers, analysis, nil
}

// calculateFrequencyWeight 出現頻度重みを計算
func calculateFrequencyWeight(number int, analysisRange int) (float64, error) {
	countResult, err := numberCount.CountSpecificNumber(number, analysisRange)
	if err != nil {
		return 0.0, err
	}
	
	// 基本重みは出現率をベースに設定
	weight := countResult.Percentage / 100.0
	
	// 最近の傾向を重視（直近データの重みを増加）
	recentWeight := calculateRecentTrend(number, analysisRange/3)
	weight = (weight * 0.7) + (recentWeight * 0.3)
	
	return weight, nil
}

// calculateConsecutiveWeight 連続数字の重みを計算
func calculateConsecutiveWeight(number int) float64 {
	// ロト7では連続数字が出現しやすい傾向があるとされるため
	// 隣接する数字との組み合わせで重みを調整
	baseWeight := 1.0
	
	// 両端の数字（1, 37）は連続が作りにくいため重み減少
	if number == 1 || number == 37 {
		baseWeight = 0.8
	} else if number == 2 || number == 36 {
		baseWeight = 0.9
	}
	
	// 中央付近の数字は連続を作りやすいため重み増加
	if number >= 15 && number <= 23 {
		baseWeight = 1.2
	}
	
	return baseWeight
}

// calculateIntervalWeight 数字間隔の重みを計算
func calculateIntervalWeight(number int, analysisRange int) float64 {
	// 数字間隔のパターン分析
	// 一般的に等間隔で並ぶ数字や、特定の間隔パターンを分析
	baseWeight := 1.0
	
	// 5の倍数周辺は選ばれやすい傾向があるとされる
	if number%5 == 0 || (number-1)%5 == 0 || (number+1)%5 == 0 {
		baseWeight = 1.1
	}
	
	// 7の倍数周辺も同様
	if number%7 == 0 || (number-1)%7 == 0 || (number+1)%7 == 0 {
		baseWeight *= 1.05
	}
	
	return baseWeight
}

// calculateOddEvenWeight 奇数偶数バランスの重みを計算
func calculateOddEvenWeight(number int) float64 {
	// ロト7では一般的に奇数偶数が3:4または4:3の比率で出現することが多い
	// 完全に偏らせず、バランスを取る重み付け
	if number%2 == 1 { // 奇数
		return 1.02 // 奇数をわずかに優遇（統計的に若干多い傾向）
	} else { // 偶数
		return 0.98
	}
}

// calculateRangeWeight 範囲バランスの重みを計算
func calculateRangeWeight(number int) float64 {
	// 前半（1-18）と後半（19-37）のバランス分析
	// 一般的に前半後半がバランス良く選ばれる傾向
	if number <= 18 { // 前半
		return 1.01
	} else { // 後半
		return 0.99
	}
}

// calculatePeriodicityWeight 周期性の重みを計算
func calculatePeriodicityWeight(number int, analysisRange int) float64 {
	// 特定の周期での出現パターンを分析
	// 簡易実装として、数字の性質による重み付け
	baseWeight := 1.0
	
	// 素数は統計的に選ばれやすいとされる
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}
	for _, prime := range primes {
		if number == prime {
			baseWeight = 1.05
			break
		}
	}
	
	// 完全数や特殊な数字の重み調整
	specialNumbers := []int{6, 28} // 6と28は完全数
	for _, special := range specialNumbers {
		if number == special {
			baseWeight *= 1.03
			break
		}
	}
	
	return baseWeight
}

// calculateCorrelationWeight 前回との相関重みを計算
func calculateCorrelationWeight(number int, analysisRange int) float64 {
	// 前回の当選番号との相関を分析
	// 実装上は前回と同じ数字は避ける傾向、近い数字は選ばれやすい傾向を模擬
	baseWeight := 1.0
	
	// 簡易実装: ランダム要素を含む相関重み
	// 実際の実装では過去の当選データを取得して分析する必要がある
	variation := globalRand.Float64()*0.2 - 0.1 // -0.1 to +0.1の範囲
	baseWeight += variation
	
	// 最低重みを保証
	if baseWeight < 0.5 {
		baseWeight = 0.5
	}
	
	return baseWeight
}

// calculateOverallOddEvenBalance 全体の奇数偶数バランスを分析
func calculateOverallOddEvenBalance(analysisRange int) float64 {
	// 過去データの奇数偶数比率を分析
	// 簡易実装として理想的なバランス値を返す
	return 0.53 // 奇数がわずかに多い傾向を表現
}

// calculateOverallRangeBalance 全体の範囲バランスを分析
func calculateOverallRangeBalance(analysisRange int) float64 {
	// 過去データの前半後半比率を分析
	// 簡易実装として理想的なバランス値を返す
	return 0.51 // 前半がわずかに多い傾向を表現
}

// calculateRecentTrend 最近の出現傾向を計算
// 引数: number - 対象数字, recentRange - 最近の範囲
// 戻り値: 最近の傾向重み
func calculateRecentTrend(number int, recentRange int) float64 {
	if recentRange <= 0 {
		recentRange = 10 // デフォルト値
	}

	countResult, err := numberCount.CountSpecificNumber(number, recentRange)
	if err != nil {
		return 0.027 // エラー時は期待値（1/37）
	}

	// 最近の出現率を重みとして返す
	recentWeight := countResult.Percentage / 100.0
	if recentWeight < 0.01 {
		recentWeight = 0.01
	}

	return recentWeight
}

// selectDiverseWeightedNumbers 多様性重み付き数字選択
// 引数: weightedNumbers - 重み付き数字配列, count - 選択個数, avoidDuplicatePatterns - 重複パターン回避
// 戻り値: 選択された数字配列
func selectDiverseWeightedNumbers(weightedNumbers []WeightedNumber, count int, avoidDuplicatePatterns bool) []int {
	var selected []int
	available := make([]WeightedNumber, len(weightedNumbers))
	copy(available, weightedNumbers)

	for len(selected) < count && len(available) > 0 {
		// 重みに基づいて数字を選択
		selectedIndex := improvedWeightedRandomSelect(available)
		selectedNumber := available[selectedIndex].Number
		selected = append(selected, selectedNumber)
		
		// 選択された数字を利用可能リストから除去
		available = append(available[:selectedIndex], available[selectedIndex+1:]...)
		
		// 多様性を確保するため、連続数字の重みを一時的に調整
		if avoidDuplicatePatterns && len(selected) < count {
			available = adjustWeightsForDiversity(available, selectedNumber)
		}
	}

	return selected
}

// selectBonusNumbers ボーナス数字を選択（本数字と重複回避）
// 引数: weightedNumbers - 重み付き数字配列, mainNumbers - 本数字配列, count - 選択個数
// 戻り値: 選択されたボーナス数字配列
func selectBonusNumbers(weightedNumbers []WeightedNumber, mainNumbers []int, count int) []int {
	// 本数字を除外したリストを作成
	var available []WeightedNumber
	mainSet := make(map[int]bool)
	for _, num := range mainNumbers {
		mainSet[num] = true
	}
	
	for _, wn := range weightedNumbers {
		if !mainSet[wn.Number] {
			available = append(available, wn)
		}
	}
	
	// ボーナス数字は本数字とは違う統計特性を持つため、重みを調整
	for i := range available {
		// ボーナス数字は本数字より出現頻度が低いため重み調整
		available[i].Weight *= 0.8
	}
	
	return selectDiverseWeightedNumbers(available, count, false)
}

// adjustWeightsForDiversity 多様性確保のため重みを調整
func adjustWeightsForDiversity(numbers []WeightedNumber, lastSelected int) []WeightedNumber {
	adjusted := make([]WeightedNumber, len(numbers))
	copy(adjusted, numbers)
	
	for i := range adjusted {
		num := adjusted[i].Number
		
		// 連続数字の重みを下げる
		if math.Abs(float64(num-lastSelected)) <= 2 {
			adjusted[i].Weight *= 0.7
		}
		
		// 同じ末尾数字の重みを下げる
		if num%10 == lastSelected%10 {
			adjusted[i].Weight *= 0.8
		}
	}
	
	return adjusted
}

// improvedWeightedRandomSelect 改良版重み付きランダム選択
// 引数: numbers - 重み付き数字配列
// 戻り値: 選択されたインデックス
func improvedWeightedRandomSelect(numbers []WeightedNumber) int {
	if len(numbers) == 0 {
		return 0
	}
	
	if len(numbers) == 1 {
		return 0
	}
	
	// 重みの合計を計算
	totalWeight := 0.0
	for _, num := range numbers {
		totalWeight += num.Weight
	}
	
	// 重みがすべて0の場合は均等選択
	if totalWeight <= 0 {
		return globalRand.Intn(len(numbers))
	}

	// ランダム値を生成
	randomValue := globalRand.Float64() * totalWeight

	// 重みに基づいてインデックスを選択
	currentWeight := 0.0
	for i, num := range numbers {
		currentWeight += num.Weight
		if randomValue <= currentWeight {
			return i
		}
	}

	// フォールバック（通常は到達しない）
	return globalRand.Intn(len(numbers))
}

// generateImprovedRandomNumbers 改良版完全ランダム数字生成
// 引数: count - 生成する数字の個数
// 戻り値: 生成された数字配列
func generateImprovedRandomNumbers(count int) []int {
	var numbers []int
	used := make(map[int]bool)

	for len(numbers) < count {
		number := globalRand.Intn(37) + 1 // 1-37の範囲
		if !used[number] {
			numbers = append(numbers, number)
			used[number] = true
		}
	}

	return numbers
}

// GenerateMultiplePredictions 複数の予想番号を生成（改良版）
// 引数: count - 生成する予想の数, analysisRange - 分析範囲, useTrend - 傾向を使用するか
// 戻り値: 予想結果配列、エラー
func GenerateMultiplePredictions(count int, analysisRange int, useTrend bool) ([]*PredictionResult, error) {
	return GenerateMultiplePredictionsWithPattern(count, analysisRange, useTrend, GetDefaultStatisticalPattern())
}

// GenerateMultiplePredictionsWithPattern パターン指定での複数予想生成
// 引数: count - 生成数, analysisRange - 分析範囲, useTrend - 傾向使用, pattern - 統計パターン
// 戻り値: 予想結果配列、エラー
func GenerateMultiplePredictionsWithPattern(count int, analysisRange int, useTrend bool, pattern *StatisticalPattern) ([]*PredictionResult, error) {
	if count <= 0 {
		return nil, fmt.Errorf("生成数は1以上で指定してください。入力値: %d", count)
	}

	var predictions []*PredictionResult
	
	// 多様性確保のため、若干異なる設定で複数パターンを生成
	for i := 0; i < count; i++ {
		// 各予想で微細な設定変更を加えて多様性を確保
		adjustedPattern := createVariedPattern(pattern, i, count)
		
		prediction, err := GenerateRandomPredictionWithPattern(analysisRange, useTrend, adjustedPattern)
		if err != nil {
			// 個別の生成に失敗した場合はログ出力して継続
			fmt.Printf("警告: 予想番号生成に失敗しました（%d回目）: %v\n", i+1, err)
			
			// フォールバック: 完全ランダム生成を試行
			fallbackPrediction, fallbackErr := GenerateRandomPredictionWithPattern(analysisRange, false, nil)
			if fallbackErr != nil {
				return nil, fmt.Errorf("予想番号生成に失敗しました（%d回目、フォールバックも失敗）: %v", i+1, fallbackErr)
			}
			predictions = append(predictions, fallbackPrediction)
		} else {
			predictions = append(predictions, prediction)
		}
		
		// 短時間での重複を避けるため、わずかに待機
		time.Sleep(time.Millisecond * 10)
	}

	// 生成された予想の重複チェックと修正
	predictions = ensureDiversePredictions(predictions)

	return predictions, nil
}

// createVariedPattern 多様性確保のため設定を微調整したパターンを作成
func createVariedPattern(basePattern *StatisticalPattern, index int, total int) *StatisticalPattern {
	if basePattern == nil {
		basePattern = GetDefaultStatisticalPattern()
	}
	
	// ベースパターンをコピー
	varied := *basePattern
	
	// インデックスに基づいて重み配分を微調整
	variation := float64(index) / float64(total) * 0.1 // 最大10%の変動
	
	// 各重みを微調整（合計が1.0になるよう正規化）
	varied.FrequencyWeight += variation * (globalRand.Float64()*2 - 1) // -variation to +variation
	varied.ConsecutiveWeight += variation * (globalRand.Float64()*2 - 1)
	varied.IntervalWeight += variation * (globalRand.Float64()*2 - 1)
	
	// 重みが負にならないよう調整
	if varied.FrequencyWeight < 0.05 {
		varied.FrequencyWeight = 0.05
	}
	if varied.ConsecutiveWeight < 0.05 {
		varied.ConsecutiveWeight = 0.05
	}
	if varied.IntervalWeight < 0.05 {
		varied.IntervalWeight = 0.05
	}
	
	// 重みの合計を正規化
	totalWeight := varied.FrequencyWeight + varied.ConsecutiveWeight + varied.IntervalWeight + 
					varied.OddEvenWeight + varied.RangeWeight + varied.PeriodicityWeight + varied.CorrelationWeight
	
	if totalWeight > 0 {
		varied.FrequencyWeight /= totalWeight
		varied.ConsecutiveWeight /= totalWeight
		varied.IntervalWeight /= totalWeight
		varied.OddEvenWeight /= totalWeight
		varied.RangeWeight /= totalWeight
		varied.PeriodicityWeight /= totalWeight
		varied.CorrelationWeight /= totalWeight
	}
	
	return &varied
}

// ensureDiversePredictions 予想結果の多様性を確保
func ensureDiversePredictions(predictions []*PredictionResult) []*PredictionResult {
	if len(predictions) <= 1 {
		return predictions
	}
	
	// 重複度をチェックし、類似度が高い予想を再生成
	for i := 0; i < len(predictions); i++ {
		for j := i + 1; j < len(predictions); j++ {
			similarity := calculatePredictionSimilarity(predictions[i], predictions[j])
			
			// 類似度が80%以上の場合は再生成
			if similarity > 0.8 {
				fmt.Printf("類似予想を検出（類似度: %.2f%%）、予想%dを再生成します\n", similarity*100, j+1)
				
				// 再生成（より多様な設定で）
				newPattern := GetDefaultStatisticalPattern()
				newPattern.FrequencyWeight *= 0.5 // 頻度重みを下げて多様性重視
				newPattern.ConsecutiveWeight *= 1.5 // 連続性重みを上げる
				
				newPrediction, err := GenerateRandomPredictionWithPattern(100, true, newPattern)
				if err == nil {
					predictions[j] = newPrediction
				}
			}
		}
	}
	
	return predictions
}

// calculatePredictionSimilarity 予想間の類似度を計算
func calculatePredictionSimilarity(p1, p2 *PredictionResult) float64 {
	if p1 == nil || p2 == nil {
		return 0.0
	}
	
	// 本数字の重複数をカウント
	mainOverlap := 0
	for _, num1 := range p1.MainNumbers {
		for _, num2 := range p2.MainNumbers {
			if num1 == num2 {
				mainOverlap++
				break
			}
		}
	}
	
	// ボーナス数字の重複数をカウント
	bonusOverlap := 0
	for _, num1 := range p1.BonusNumbers {
		for _, num2 := range p2.BonusNumbers {
			if num1 == num2 {
				bonusOverlap++
				break
			}
		}
	}
	
	// 類似度を計算（本数字重視）
	mainSimilarity := float64(mainOverlap) / 7.0
	bonusSimilarity := float64(bonusOverlap) / 2.0
	
	// 重み付け平均（本数字80%, ボーナス20%）
	totalSimilarity := (mainSimilarity * 0.8) + (bonusSimilarity * 0.2)
	
	return totalSimilarity
}

// PrintPredictionResult 予想結果を表示（改良版）
// 引数: result - 表示したい予想結果
func PrintPredictionResult(result *PredictionResult) {
	PrintPredictionResultWithDetails(result, false)
}

// PrintPredictionResultWithDetails 詳細情報付きで予想結果を表示
// 引数: result - 予想結果, showDetails - 詳細分析情報を表示するか
func PrintPredictionResultWithDetails(result *PredictionResult, showDetails bool) {
	fmt.Println("🎯 ロト7予想番号")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("生成時刻: %s\n", result.GeneratedTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("予想タイプ: ")
	if result.TrendBased {
		fmt.Println("📊 高度統計分析")
	} else {
		fmt.Println("🎲 完全ランダム")
	}
	fmt.Println(strings.Repeat("-", 60))
	
	fmt.Printf("🎯 本数字: ")
	for i, num := range result.MainNumbers {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%02d", num)
	}
	fmt.Println()
	
	fmt.Printf("⭐ ボーナス数字: ")
	for i, num := range result.BonusNumbers {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%02d", num)
	}
	fmt.Println()
	
	// 分析情報表示
	if result.TrendBased {
		oddCount := 0
		evenCount := 0
		frontCount := 0 // 1-18
		backCount := 0  // 19-37
		
		for _, num := range result.MainNumbers {
			if num%2 == 1 {
				oddCount++
			} else {
				evenCount++
			}
			if num <= 18 {
				frontCount++
			} else {
				backCount++
			}
		}
		
		fmt.Println(strings.Repeat("-", 60))
		fmt.Printf("📈 分析情報: 奇数%d/偶数%d, 前半%d/後半%d\n", 
			oddCount, evenCount, frontCount, backCount)
		
		// 詳細分析情報表示
		if showDetails && result.AnalysisInfo != nil {
			fmt.Printf("⏱️  分析時間: %v\n", result.AnalysisInfo.TotalAnalysisTime)
			
			// 上位重み数字表示
			fmt.Printf("🏆 高重み数字（上位5位）:\n")
			var weightPairs []struct {
				number int
				weight float64
			}
			
			for num, weight := range result.AnalysisInfo.FrequencyWeights {
				weightPairs = append(weightPairs, struct {
					number int
					weight float64
				}{num, weight})
			}
			
			// 重みでソート
			sort.Slice(weightPairs, func(i, j int) bool {
				return weightPairs[i].weight > weightPairs[j].weight
			})
			
			for i := 0; i < 5 && i < len(weightPairs); i++ {
				fmt.Printf("   %d位: %02d (重み: %.3f)\n", 
					i+1, weightPairs[i].number, weightPairs[i].weight)
			}
		}
	}
	
	fmt.Println(strings.Repeat("=", 60))
}

// PrintMultiplePredictions 複数の予想結果を表示（改良版）
// 引数: results - 表示したい予想結果配列
func PrintMultiplePredictions(results []*PredictionResult) {
	PrintMultiplePredictionsWithDetails(results, false)
}

// PrintMultiplePredictionsWithDetails 詳細情報付きで複数予想結果を表示
// 引数: results - 予想結果配列, showDetails - 詳細情報を表示するか
func PrintMultiplePredictionsWithDetails(results []*PredictionResult, showDetails bool) {
	fmt.Println("🎯 ロト7複数予想番号")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("総予想数: %d組\n", len(results))
	if len(results) > 0 {
		fmt.Printf("予想タイプ: ")
		if results[0].TrendBased {
			fmt.Println("📊 高度統計分析ベース")
		} else {
			fmt.Println("🎲 完全ランダム")
		}
	}
	fmt.Println(strings.Repeat("-", 70))
	
	for i, result := range results {
		fmt.Printf("【予想%d】", i+1)
		fmt.Printf(" 本数字: ")
		for j, num := range result.MainNumbers {
			if j > 0 {
				fmt.Printf(",")
			}
			fmt.Printf("%02d", num)
		}
		fmt.Printf(" ボーナス: ")
		for j, num := range result.BonusNumbers {
			if j > 0 {
				fmt.Printf(",")
			}
			fmt.Printf("%02d", num)
		}
		
		// 簡易分析情報
		if result.TrendBased && showDetails {
			oddCount := 0
			for _, num := range result.MainNumbers {
				if num%2 == 1 {
					oddCount++
				}
			}
			fmt.Printf(" (奇:%d偶:%d)", oddCount, 7-oddCount)
		}
		
		fmt.Println()
	}
	
	// 全体的な統計情報
	if len(results) > 1 && showDetails {
		fmt.Println(strings.Repeat("-", 70))
		fmt.Println("📊 全体統計:")
		
		// 最頻出数字
		numberCounts := make(map[int]int)
		for _, result := range results {
			for _, num := range result.MainNumbers {
				numberCounts[num]++
			}
		}
		
		fmt.Printf("🔥 最頻出数字: ")
		type numberFreq struct {
			number int
			count  int
		}
		var freqs []numberFreq
		for num, count := range numberCounts {
			freqs = append(freqs, numberFreq{num, count})
		}
		
		sort.Slice(freqs, func(i, j int) bool {
			return freqs[i].count > freqs[j].count
		})
		
		for i := 0; i < 5 && i < len(freqs); i++ {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%02d(%d回)", freqs[i].number, freqs[i].count)
		}
		fmt.Println()
	}
	
	fmt.Println(strings.Repeat("=", 70))
}
