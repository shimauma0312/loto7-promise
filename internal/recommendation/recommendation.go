package recommendation

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

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

// CombinationHistory 組み合わせ履歴
type CombinationHistory struct {
	Pairs map[string]int // "num1,num2"形式のペア出現回数
}

// RecommendationEngine 推薦エンジン
type RecommendationEngine struct {
	config  RecommendationConfig
	results [][]string
	history CombinationHistory
}

// NewRecommendationEngine 推薦エンジンを作成
func NewRecommendationEngine(config RecommendationConfig) *RecommendationEngine {
	return &RecommendationEngine{
		config: config,
		history: CombinationHistory{
			Pairs: make(map[string]int),
		},
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
	return nil
}

// buildCombinationHistory 組み合わせ履歴を構築
func (re *RecommendationEngine) buildCombinationHistory() {
	for _, numbers := range re.results {
		// 7個の数字から2個ずつの組み合わせを生成
		for i := 0; i < len(numbers); i++ {
			for j := i + 1; j < len(numbers); j++ {
				num1, _ := strconv.Atoi(numbers[i])
				num2, _ := strconv.Atoi(numbers[j])

				// ペアキーを作成（小さい数字,大きい数字の順）
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

// calculateNumberScores 各数字のスコアを計算
func (re *RecommendationEngine) calculateNumberScores() []NumberScore {
	scores := make([]NumberScore, 37) // 1-37の数字

	// 各数字を初期化
	for i := 1; i <= 37; i++ {
		scores[i-1].Number = i
	}

	// 出現頻度の計算
	re.calculateFrequencyScores(scores)

	// 最近の出現スコア計算
	re.calculateRecentScores(scores)

	// 連続出現スコア計算
	re.calculateConsecutiveScores(scores)

	// 位置スコア計算
	re.calculatePositionScores(scores)

	// 総合スコア計算
	for i := range scores {
		scores[i].TotalScore = scores[i].FrequencyScore*re.config.FrequencyWeight +
			scores[i].RecentScore*re.config.RecentWeight +
			scores[i].ConsecutiveScore*re.config.ConsecutiveWeight +
			scores[i].PositionScore*re.config.PositionWeight
	}

	return scores
}

// calculateFrequencyScores 出現頻度スコアを計算
func (re *RecommendationEngine) calculateFrequencyScores(scores []NumberScore) {
	frequency := make(map[int]int)

	// 過去100回以内の出現回数をカウント
	lookback := re.config.HistoryLookback
	if lookback > len(re.results) {
		lookback = len(re.results)
	}

	for i := 0; i < lookback; i++ {
		for _, numStr := range re.results[i] {
			num, _ := strconv.Atoi(numStr)
			frequency[num]++
		}
	}

	// スコア化（過去100回以内に出現していない数字はスコア0）
	for i := range scores {
		num := scores[i].Number
		if count, exists := frequency[num]; exists {
			scores[i].FrequencyScore = float64(count) / float64(lookback)
		} else {
			scores[i].FrequencyScore = 0 // 過去100回以内に出現していない
		}
	}
}

// calculateRecentScores 最近の出現スコアを計算
func (re *RecommendationEngine) calculateRecentScores(scores []NumberScore) {
	// 過去50回以内の出現は加点
	recent50 := make(map[int]bool)
	lookback50 := 50
	if lookback50 > len(re.results) {
		lookback50 = len(re.results)
	}

	for i := 0; i < lookback50; i++ {
		for _, numStr := range re.results[i] {
			num, _ := strconv.Atoi(numStr)
			recent50[num] = true
		}
	}

	// 直近の回避（直前数回で出た数字は減点）
	recentDrawn := make(map[int]bool)
	avoidCount := re.config.RecentAvoidCount
	if avoidCount > len(re.results) {
		avoidCount = len(re.results)
	}

	for i := 0; i < avoidCount; i++ {
		for _, numStr := range re.results[i] {
			num, _ := strconv.Atoi(numStr)
			recentDrawn[num] = true
		}
	}

	for i := range scores {
		num := scores[i].Number

		// 直近で出た数字は大幅減点
		if recentDrawn[num] {
			scores[i].RecentScore = -1.0
			scores[i].IsRecentlyDrawn = true
		} else if recent50[num] {
			// 過去50回以内で出た数字は加点
			scores[i].RecentScore = 0.5
		} else {
			scores[i].RecentScore = 0.0
		}
	}
}

// calculateConsecutiveScores 連続出現スコアを計算
func (re *RecommendationEngine) calculateConsecutiveScores(scores []NumberScore) {
	for i := range scores {
		num := scores[i].Number

		// 直近での連続出現をチェック
		consecutiveCount := 0
		for j := 0; j < len(re.results) && j < 20; j++ { // 直近20回をチェック
			found := false
			for _, numStr := range re.results[j] {
				if n, _ := strconv.Atoi(numStr); n == num {
					found = true
					break
				}
			}
			if found {
				consecutiveCount++
			} else {
				break
			}
		}

		// 3回以上連続で出ている場合はブースト
		if consecutiveCount >= re.config.ConsecutiveBoost {
			scores[i].ConsecutiveScore = float64(consecutiveCount) * 0.2
		} else {
			scores[i].ConsecutiveScore = 0.0
		}
	}
}

// calculatePositionScores 位置によるスコアを計算
func (re *RecommendationEngine) calculatePositionScores(scores []NumberScore) {
	// 最初の数字が1に近い場合の考慮
	if len(re.results) > 0 && len(re.results[0]) > 0 {
		firstNum, _ := strconv.Atoi(re.results[0][0])

		// 最初の数字が小さい場合、小さい数字群にボーナス
		if firstNum <= 10 {
			for i := range scores {
				if scores[i].Number <= 20 {
					scores[i].PositionScore = 0.3
				} else {
					scores[i].PositionScore = -0.1
				}
			}
		} else {
			// 均等なスコア
			for i := range scores {
				scores[i].PositionScore = 0.0
			}
		}
	}
}

// isPairFrequentlyUsed ペアが頻繁に使用されているかチェック
func (re *RecommendationEngine) isPairFrequentlyUsed(num1, num2 int) bool {
	var pairKey string
	if num1 < num2 {
		pairKey = fmt.Sprintf("%d,%d", num1, num2)
	} else {
		pairKey = fmt.Sprintf("%d,%d", num2, num1)
	}

	// 過去に2回以上出現している組み合わせは避ける
	return re.history.Pairs[pairKey] >= 2
}

// GenerateRecommendations 推薦組み合わせを生成
func (re *RecommendationEngine) GenerateRecommendations() ([][]int, error) {
	if len(re.results) == 0 {
		return nil, fmt.Errorf("データが読み込まれていません")
	}

	scores := re.calculateNumberScores()

	// スコア順にソート
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].TotalScore > scores[j].TotalScore
	})

	var recommendations [][]int

	for rec := 0; rec < re.config.MaxRecommendations; rec++ {
		// 各推薦で異なる組み合わせを生成するため、試行回数を制限
		maxAttempts := 50
		for attempt := 0; attempt < maxAttempts; attempt++ {
			combination := re.generateSingleRecommendation(scores, recommendations, rec+attempt)
			if len(combination) == 7 && !re.isDuplicateCombination(combination, recommendations) {
				recommendations = append(recommendations, combination)
				break
			}
		}
	}

	return recommendations, nil
}

// generateSingleRecommendation 単一の推薦組み合わせを生成
func (re *RecommendationEngine) generateSingleRecommendation(scores []NumberScore, existing [][]int, variation int) []int {
	var combination []int
	usedNumbers := make(map[int]bool)

	// バリエーションを作るため、開始インデックスを変更
	startIndex := variation % len(scores)

	// スコア順に候補を選択（ただし制約を考慮）
	for i := 0; i < len(scores); i++ {
		if len(combination) >= 7 {
			break
		}

		scoreIndex := (startIndex + i) % len(scores)
		score := scores[scoreIndex]
		num := score.Number

		// 既に使用済みの数字はスキップ
		if usedNumbers[num] {
			continue
		}

		// 過去100回以内に出現していない数字はスキップ
		if score.FrequencyScore == 0 {
			continue
		}

		// 組み合わせが既存と重複しないかチェック
		if re.isValidAddition(num, combination, existing) {
			combination = append(combination, num)
			usedNumbers[num] = true
		}
	}

	// 7個に満たない場合は、制約を緩めて補完
	if len(combination) < 7 {
		combination = re.fillCombination(combination, scores, existing)
	}

	sort.Ints(combination)
	return combination
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

// isValidAddition 数字を追加しても良いかチェック
func (re *RecommendationEngine) isValidAddition(num int, currentComb []int, existing [][]int) bool {
	// 現在の組み合わせ内での重複ペアチェック
	for _, existingNum := range currentComb {
		if re.isPairFrequentlyUsed(num, existingNum) {
			return false
		}
	}

	// 既存の推薦との重複度チェック
	for _, existingComb := range existing {
		overlap := 0
		testComb := append(currentComb, num)
		for _, testNum := range testComb {
			for _, existingNum := range existingComb {
				if testNum == existingNum {
					overlap++
				}
			}
		}
		// 4個以上重複する場合は避ける
		if overlap >= 4 {
			return false
		}
	}

	return true
}

// fillCombination 組み合わせを7個まで補完
func (re *RecommendationEngine) fillCombination(combination []int, scores []NumberScore, existing [][]int) []int {
	usedNumbers := make(map[int]bool)
	for _, num := range combination {
		usedNumbers[num] = true
	}

	// 制約を緩めて残りを埋める
	for _, score := range scores {
		if len(combination) >= 7 {
			break
		}

		num := score.Number
		if usedNumbers[num] || score.FrequencyScore == 0 {
			continue
		}

		combination = append(combination, num)
		usedNumbers[num] = true
	}

	return combination
}

// FormatRecommendations 推薦結果をフォーマット
func (re *RecommendationEngine) FormatRecommendations(recommendations [][]int) string {
	var result strings.Builder

	result.WriteString("=== ロト7 推薦番号 ===\n\n")

	for i, combination := range recommendations {
		result.WriteString(fmt.Sprintf("推薦 %d: ", i+1))
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

	return result.String()
}
