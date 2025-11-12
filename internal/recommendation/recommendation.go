package recommendation

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shimauma0312/loto7-promise/internal/result"
)

// RecommendationConfig 推薦設定
type RecommendationConfig struct {
	MaxRecommendations int     // 最大推薦数
	RecentAvoidCount   int     // 直近の回避回数
	ConsecutiveBoost   int     // 連続出現ブースト判定回数
	HistoryLookback    int     // 過去の履歴確認回数
	FrequencyWeight    float64 // 出現頻度の重み
	RecentWeight       float64 // 最近の出現の重み
	ConsecutiveWeight  float64 // 連続出現の重み
	PositionWeight     float64 // 位置による重み
}

// DefaultConfig デフォルト設定
func DefaultConfig() RecommendationConfig {
	return RecommendationConfig{
		MaxRecommendations: 5,
		RecentAvoidCount:   3,
		ConsecutiveBoost:   3,
		HistoryLookback:    100,
		FrequencyWeight:    0.3,
		RecentWeight:       0.4,
		ConsecutiveWeight:  0.2,
		PositionWeight:     0.1,
	}
}

// NumberScore 数字のスコア
type NumberScore struct {
	Number           int
	FrequencyScore   float64
	RecentScore      float64
	ConsecutiveScore float64
	PositionScore    float64
	TotalScore       float64
	IsRecentlyDrawn  bool
}

// StatisticalAnalysis 統計分析結果
type StatisticalAnalysis struct {
	FrequentNumbers   []int       // 出現回数が多い数字
	RareNumbers       []int       // 出現回数が少ない数字
	HotNumbers        []int       // 直近10回以内で連続して出ている数字
	RevivalCandidates []int       // 20回以上出ていない復活候補
	FrequencyMap      map[int]int // 各数字の出現回数
	LastAppearance    map[int]int // 各数字の最後の出現回数
}

// CombinationHistory 組み合わせ履歴
type CombinationHistory struct {
	Pairs map[string]int // "num1,num2"形式のペア出現回数
}

// RecommendationEngine 推薦エンジンのすとらくと
type RecommendationEngine struct {
	config     RecommendationConfig
	results    [][]string
	history    CombinationHistory
	rng        *rand.Rand
	statistics StatisticalAnalysis
}

// NewRecommendationEngine 推薦エンジンを作成
func NewRecommendationEngine(config RecommendationConfig) *RecommendationEngine {
	// 現在時刻をシードにしてランダムジェネレータ初期化
	source := rand.NewSource(time.Now().UnixNano())
	return &RecommendationEngine{
		config: config,
		history: CombinationHistory{
			Pairs: make(map[string]int),
		},
		rng: rand.New(source),
	}
}

// LoadData データを読み込み
func (re *RecommendationEngine) LoadData(count int) error {
	re.results = result.GetResult(count)
	if len(re.results) == 0 {
		return fmt.Errorf("データを取得できませんでした")
	}

	// 組み合わせ履歴を構築
	re.buildCombinationHistory()

	// 統計分析を実行
	re.buildStatisticalAnalysis()

	return nil
}

// buildStatisticalAnalysis 統計分析を実行
func (re *RecommendationEngine) buildStatisticalAnalysis() {
	re.statistics.FrequencyMap = make(map[int]int)
	re.statistics.LastAppearance = make(map[int]int)

	// 過去100回分の出現回数をカウント
	lookback := 100
	if lookback > len(re.results) {
		lookback = len(re.results)
	}

	// 各数字の出現回数と最終出現位置を記録
	for i := 0; i < lookback; i++ {
		for _, numStr := range re.results[i] {
			num, _ := strconv.Atoi(numStr)
			re.statistics.FrequencyMap[num]++

			// 最終出現位置を更新（最新の結果が0）
			if _, exists := re.statistics.LastAppearance[num]; !exists {
				re.statistics.LastAppearance[num] = i
			}
		}
	}

	// 出現回数でソート
	type numFreq struct {
		num  int
		freq int
	}
	var frequencies []numFreq
	for num := 1; num <= 37; num++ {
		frequencies = append(frequencies, numFreq{num: num, freq: re.statistics.FrequencyMap[num]})
	}
	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].freq > frequencies[j].freq
	})

	// 出現回数が多い数字（上位10個）
	for i := 0; i < 10 && i < len(frequencies); i++ {
		re.statistics.FrequentNumbers = append(re.statistics.FrequentNumbers, frequencies[i].num)
	}

	// 出現回数が少ない数字（下位10個）
	for i := len(frequencies) - 1; i >= len(frequencies)-10 && i >= 0; i-- {
		re.statistics.RareNumbers = append(re.statistics.RareNumbers, frequencies[i].num)
	}

	// 直近10回以内で連続して出ている数字（波が来ている）
	hotCheckRange := 10
	if hotCheckRange > len(re.results) {
		hotCheckRange = len(re.results)
	}
	hotCount := make(map[int]int)
	for i := 0; i < hotCheckRange; i++ {
		for _, numStr := range re.results[i] {
			num, _ := strconv.Atoi(numStr)
			hotCount[num]++
		}
	}
	// 直近10回で3回以上出現している数字を「波が来ている」と判断
	for num, count := range hotCount {
		if count >= 3 {
			re.statistics.HotNumbers = append(re.statistics.HotNumbers, num)
		}
	}

	// 20回以上出ていない復活候補
	for num := 1; num <= 37; num++ {
		lastAppear, exists := re.statistics.LastAppearance[num]
		if !exists || lastAppear >= 20 {
			// ただし、出現回数が少ない数字のリストの上位5個に含まれている場合は無視
			isRare := false
			for i := 0; i < 5 && i < len(re.statistics.RareNumbers); i++ {
				if re.statistics.RareNumbers[i] == num {
					isRare = true
					break
				}
			}
			if !isRare {
				re.statistics.RevivalCandidates = append(re.statistics.RevivalCandidates, num)
			}
		}
	}
}

// buildCombinationHistory 組み合わせ履歴を構築
func (re *RecommendationEngine) buildCombinationHistory() {
	for _, numbers := range re.results {
		// 7個の数字から2個ずつの組み合わせを生成
		for i := 0; i < len(numbers); i++ {
			for j := i + 1; j < len(numbers); j++ {
				num1, _ := strconv.Atoi(numbers[i])
				num2, _ := strconv.Atoi(numbers[j])

				// ペアキー作成（小さい数字,大きい数字の順）
				var pairKey string
				if num1 < num2 {
					pairKey = fmt.Sprintf("%d,%d", num1, num2)
				} else {
					pairKey = fmt.Sprintf("%d,%d", num2, num1)
				}

				re.history.Pairs[pairKey]++
			}
		}
	}
}

// GenerateRecommendations 推薦組み合わせを生成
func (re *RecommendationEngine) GenerateRecommendations() ([][]int, error) {
	if len(re.results) == 0 {
		return nil, fmt.Errorf("データが読み込まれていません")
	}

	var recommendations [][]int
	maxAttempts := 2000 // 試行回数を増やす（制約が厳しくなったため）

	// 各推薦を生成
	for rec := 0; rec < re.config.MaxRecommendations; rec++ {
		for attempt := 0; attempt < maxAttempts; attempt++ {
			combination := re.generateAdvancedRecommendation()

			// すべての制約をチェック
			if len(combination) == 7 &&
				re.validateZoneDistribution(combination) &&
				re.validateOddEvenRatio(combination) &&
				re.validateTotalSum(combination) &&
				re.validateAverageValue(combination) &&
				re.validateSmoothness(combination) &&
				re.hasCloseNumberPair(combination) &&
				!re.isDuplicateCombination(combination, recommendations) {
				recommendations = append(recommendations, combination)
				break
			}
		}
	}

	return recommendations, nil
}

// generateAdvancedRecommendation 高度な分析に基づいた推薦を生成
func (re *RecommendationEngine) generateAdvancedRecommendation() []int {
	var combination []int
	usedNumbers := make(map[int]bool)

	// ゾーン1 (1-13): 3-4個選択
	zone1Count := 3 + re.rng.Intn(2) // 3 or 4
	combination = re.selectFromZone(combination, usedNumbers, 1, 13, zone1Count)

	// ゾーン3 (27-37): 2-3個選択（まれに3個）
	zone3Count := 2
	if re.rng.Float64() < 0.2 { // 20%の確率で3個
		zone3Count = 3
	}
	combination = re.selectFromZone(combination, usedNumbers, 27, 37, zone3Count)

	// ゾーン2 (14-26): 残りを埋める
	zone2Count := 7 - len(combination)
	combination = re.selectFromZone(combination, usedNumbers, 14, 26, zone2Count)

	// 7個に満たない場合は補完
	if len(combination) < 7 {
		combination = re.fillRemainingNumbers(combination, usedNumbers)
	}

	sort.Ints(combination)
	return combination
}

// selectFromZone 指定ゾーンから数字を選択
func (re *RecommendationEngine) selectFromZone(combination []int, usedNumbers map[int]bool, min, max, count int) []int {
	candidates := re.getZoneCandidates(min, max, usedNumbers)

	// 候補が不足している場合は調整
	if len(candidates) < count {
		count = len(candidates)
	}

	selected := 0
	attempts := 0
	maxAttempts := count * 10

	for selected < count && attempts < maxAttempts {
		attempts++

		if len(candidates) == 0 {
			break
		}

		// 重み付きランダム選択
		num := re.weightedSelectFromCandidates(candidates, usedNumbers)
		if num == 0 {
			continue
		}

		if !usedNumbers[num] {
			combination = append(combination, num)
			usedNumbers[num] = true
			selected++

			// 選択済みの候補を削除
			newCandidates := []int{}
			for _, c := range candidates {
				if c != num {
					newCandidates = append(newCandidates, c)
				}
			}
			candidates = newCandidates
		}
	}

	return combination
}

// getZoneCandidates ゾーンから候補数字を取得
func (re *RecommendationEngine) getZoneCandidates(min, max int, usedNumbers map[int]bool) []int {
	var candidates []int

	for num := min; num <= max; num++ {
		if usedNumbers[num] {
			continue
		}

		// スコア計算
		score := re.calculateNumberPriority(num)

		// スコアが一定以上の数字を候補に
		if score > -0.5 { // 閾値
			candidates = append(candidates, num)
		}
	}

	// 候補が少なすぎる場合は範囲内すべてを候補に
	if len(candidates) < 3 {
		candidates = []int{}
		for num := min; num <= max; num++ {
			if !usedNumbers[num] {
				candidates = append(candidates, num)
			}
		}
	}

	return candidates
}

// calculateNumberPriority 数字の優先度を計算
func (re *RecommendationEngine) calculateNumberPriority(num int) float64 {
	score := 0.0

	// 1. 波が来ている数字（直近10回で3回以上）は高評価
	for _, hot := range re.statistics.HotNumbers {
		if hot == num {
			score += 2.0
			break
		}
	}

	// 2. 出現回数が多い数字は中程度評価
	for _, freq := range re.statistics.FrequentNumbers {
		if freq == num {
			score += 1.0
			break
		}
	}

	// 3. 復活候補（20回以上未出現）は中程度評価
	for _, revival := range re.statistics.RevivalCandidates {
		if revival == num {
			score += 1.0
			break
		}
	}

	// 4. 出現回数が極端に少ない数字は減点
	isRare := false
	for _, rare := range re.statistics.RareNumbers {
		if rare == num {
			isRare = true
			break
		}
	}
	if isRare {
		score -= 0.5
	}

	// 5. 直近3回で出た数字は大幅減点
	recentAvoid := 3
	if recentAvoid > len(re.results) {
		recentAvoid = len(re.results)
	}
	for i := 0; i < recentAvoid; i++ {
		for _, numStr := range re.results[i] {
			if n, _ := strconv.Atoi(numStr); n == num {
				score -= 3.0
				break
			}
		}
	}

	return score
}

// weightedSelectFromCandidates 候補から重み付き選択
func (re *RecommendationEngine) weightedSelectFromCandidates(candidates []int, usedNumbers map[int]bool) int {
	if len(candidates) == 0 {
		return 0
	}

	type weightedNum struct {
		num    int
		weight float64
	}

	var weighted []weightedNum
	totalWeight := 0.0

	for _, num := range candidates {
		if usedNumbers[num] {
			continue
		}

		priority := re.calculateNumberPriority(num)
		weight := priority + 1.0 // 最低限の重み
		if weight < 0.1 {
			weight = 0.1
		}

		weighted = append(weighted, weightedNum{num: num, weight: weight})
		totalWeight += weight
	}

	if len(weighted) == 0 || totalWeight == 0 {
		return 0
	}

	// ルーレット選択
	randomValue := re.rng.Float64() * totalWeight
	currentWeight := 0.0

	for _, wn := range weighted {
		currentWeight += wn.weight
		if randomValue <= currentWeight {
			return wn.num
		}
	}

	// フォールバック
	if len(weighted) > 0 {
		return weighted[re.rng.Intn(len(weighted))].num
	}

	return 0
}

// fillRemainingNumbers 残りの数字を補完
func (re *RecommendationEngine) fillRemainingNumbers(combination []int, usedNumbers map[int]bool) []int {
	for num := 1; num <= 37 && len(combination) < 7; num++ {
		if !usedNumbers[num] {
			combination = append(combination, num)
			usedNumbers[num] = true
		}
	}
	return combination
}

// validateZoneDistribution ゾーン分布をバリデーション
func (re *RecommendationEngine) validateZoneDistribution(combination []int) bool {
	zone1 := 0 // 1-13
	zone2 := 0 // 14-26
	zone3 := 0 // 27-37

	for _, num := range combination {
		if num >= 1 && num <= 13 {
			zone1++
		} else if num >= 14 && num <= 26 {
			zone2++
		} else if num >= 27 && num <= 37 {
			zone3++
		}
	}

	// 各ゾーンが空でないこと
	if zone1 == 0 || zone2 == 0 || zone3 == 0 {
		return false
	}

	// ゾーン1は3-4個
	if zone1 < 3 || zone1 > 4 {
		return false
	}

	// ゾーン3は2-3個
	if zone3 < 2 || zone3 > 3 {
		return false
	}

	return true
}

// validateOddEvenRatio 奇数偶数比率をバリデーション
func (re *RecommendationEngine) validateOddEvenRatio(combination []int) bool {
	oddCount := 0
	evenCount := 0

	for _, num := range combination {
		if num%2 == 0 {
			evenCount++
		} else {
			oddCount++
		}
	}

	// 奇4:偶3 または 奇3:偶4
	return (oddCount == 4 && evenCount == 3) || (oddCount == 3 && evenCount == 4)
}

// validateTotalSum 合計値をバリデーション
func (re *RecommendationEngine) validateTotalSum(combination []int) bool {
	sum := 0
	for _, num := range combination {
		sum += num
	}
	return sum >= 100 && sum <= 160
}

// validateAverageValue 平均値をバリデーション（16~23付近）
func (re *RecommendationEngine) validateAverageValue(combination []int) bool {
	if len(combination) == 0 {
		return false
	}
	sum := 0
	for _, num := range combination {
		sum += num
	}
	average := float64(sum) / float64(len(combination))
	// 平均値が16~23の範囲内
	return average >= 16.0 && average <= 23.0
}

// validateSmoothness 数字の並び滑らかさをバリデーション
func (re *RecommendationEngine) validateSmoothness(combination []int) bool {
	if len(combination) < 2 {
		return true
	}

	// 隣接する数字間の差分を計算
	diffs := []int{}
	for i := 1; i < len(combination); i++ {
		diffs = append(diffs, combination[i]-combination[i-1])
	}

	// 差分の標準偏差を計算
	// 標準偏差が小さいほど滑らかな並び
	avgDiff := 0.0
	for _, diff := range diffs {
		avgDiff += float64(diff)
	}
	avgDiff /= float64(len(diffs))

	variance := 0.0
	for _, diff := range diffs {
		variance += (float64(diff) - avgDiff) * (float64(diff) - avgDiff)
	}
	variance /= float64(len(diffs))

	stdDev := variance // 平方根は計算せず分散で判定

	// 分散が大きすぎる（ジグザグが激しい）場合は却下
	// 経験的に分散が50以下だと滑らかな並び
	return stdDev <= 50.0
}

// hasCloseNumberPair 近距離ペア（差1-3）が含まれているかチェック
func (re *RecommendationEngine) hasCloseNumberPair(combination []int) bool {
	for i := 0; i < len(combination); i++ {
		for j := i + 1; j < len(combination); j++ {
			diff := abs(combination[i] - combination[j])
			if diff >= 1 && diff <= 3 {
				return true
			}
		}
	}
	return false
}

// isDuplicateCombination 重複する組み合わせかチェック
func (re *RecommendationEngine) isDuplicateCombination(newComb []int, existing [][]int) bool {
	for _, existingComb := range existing {
		if len(newComb) != len(existingComb) {
			continue
		}

		matches := 0
		for _, num1 := range newComb {
			for _, num2 := range existingComb {
				if num1 == num2 {
					matches++
					break
				}
			}
		}

		// 5個以上一致する場合は重複とみなす
		if matches >= 5 {
			return true
		}
	}
	return false
}

// abs 絶対値計算
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// FormatRecommendations 推薦結果をフォーマット
func (re *RecommendationEngine) FormatRecommendations(recommendations [][]int) string {
	var result strings.Builder

	result.WriteString("=== ロト7 推薦番号 ===\n\n")

	for i, combination := range recommendations {
		result.WriteString(fmt.Sprintf("推薦 %d: ", i+1))

		// 数字を表示
		for j, num := range combination {
			if j > 0 {
				result.WriteString(" - ")
			}
			result.WriteString(fmt.Sprintf("%02d", num))
		}

		result.WriteString("\n")
	}

	result.WriteString("\n=== 分析情報 ===\n")
	result.WriteString(fmt.Sprintf("分析対象: 過去%d回分\n", len(re.results)))
	result.WriteString(fmt.Sprintf("生成した推薦数: %d\n", len(recommendations)))

	// 統計情報を追加
	result.WriteString("\n=== 統計データ ===\n")
	result.WriteString(fmt.Sprintf("出現頻度が高い数字: %v\n", re.statistics.FrequentNumbers))
	result.WriteString(fmt.Sprintf("波が来ている数字 (直近10回): %v\n", re.statistics.HotNumbers))
	result.WriteString(fmt.Sprintf("復活候補 (20回以上未出現): %v\n", re.statistics.RevivalCandidates))

	return result.String()
}

// GetConfig 設定情報を取得
func (re *RecommendationEngine) GetConfig() RecommendationConfig {
	return re.config
}

// GetAnalyzedDrawsCount 分析した抽選回数を取得
func (re *RecommendationEngine) GetAnalyzedDrawsCount() int {
	return len(re.results)
}
