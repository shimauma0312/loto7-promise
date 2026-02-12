// ロト7の1等当選までのシミュレーション機能を提供する
//
// ユーザーの数字選択と抽選を再現し、1等が当選するまでの
// 試行回数や統計情報を計算する
package simulation

import (
	"fmt"
	"sort"
	"time"

	"github.com/shimauma0312/loto7-promise/internal/random"
)

// シミュレーション結果
type SimulationResult struct {
	UserNumbers    []int         // ユーザーが選択した数字
	DrawCount      int           // 抽選回数（1等が出るまで）
	WinningNumbers []int         // 当選した抽選の番号
	Duration       time.Duration // シミュレーション所要時間
	Attempts       int64         // 総試行回数
}

// 複数回シミュレーションの統計結果
type SimulationStats struct {
	Simulations   int           // シミュレーション実行回数
	MinDraws      int           // 最小抽選回数
	MaxDraws      int           // 最大抽選回数
	AvgDraws      float64       // 平均抽選回数
	MedianDraws   int           // 中央値
	TotalDuration time.Duration // 総所要時間
	Results       []SimulationResult // 各シミュレーションの結果
}

// シミュレーションエンジン
type SimulationEngine struct {
	randomEngine *random.RandomEngine
}

// シミュレーションエンジンを作成する
func NewSimulationEngine() *SimulationEngine {
	return &SimulationEngine{
		randomEngine: random.NewRandomEngine(),
	}
}

// データを読み込む
func (se *SimulationEngine) LoadData(count int) error {
	return se.randomEngine.LoadData(count)
}

// ユーザーの数字と抽選結果を比較する
func (se *SimulationEngine) checkMatch(userNumbers, drawNumbers []int) bool {
	if len(userNumbers) != 7 || len(drawNumbers) != 7 {
		return false
	}

	// ソート済みの配列を比較
	for i := 0; i < 7; i++ {
		if userNumbers[i] != drawNumbers[i] {
			return false
		}
	}
	return true
}

// 一致した数字の個数をカウントする
func (se *SimulationEngine) countMatches(userNumbers, drawNumbers []int) int {
	matches := 0
	userMap := make(map[int]bool)
	for _, num := range userNumbers {
		userMap[num] = true
	}

	for _, num := range drawNumbers {
		if userMap[num] {
			matches++
		}
	}
	return matches
}

// 1等が当選するまでシミュレーションを実行する
func (se *SimulationEngine) RunSimulation(userNumbers []int) (*SimulationResult, error) {
	if len(userNumbers) != 7 {
		return nil, fmt.Errorf("ユーザーの数字は7個必要です")
	}

	// ソート
	sortedUserNumbers := make([]int, len(userNumbers))
	copy(sortedUserNumbers, userNumbers)
	sort.Ints(sortedUserNumbers)

	startTime := time.Now()
	drawCount := 0
	var winningNumbers []int

	// 1等が出るまでループ
	for {
		drawCount++

		// 抽選を実行
		draw, err := se.randomEngine.GenerateRandomCombination()
		if err != nil {
			return nil, fmt.Errorf("抽選の生成に失敗: %w", err)
		}

		// 一致チェック
		if se.checkMatch(sortedUserNumbers, draw) {
			winningNumbers = draw
			break
		}

		// 安全装置（無限ループ防止）
		if drawCount > 100000000 { // 1億回で打ち切り
			return nil, fmt.Errorf("1億回の抽選で1等が出ませんでした")
		}
	}

	duration := time.Since(startTime)

	return &SimulationResult{
		UserNumbers:    sortedUserNumbers,
		DrawCount:      drawCount,
		WinningNumbers: winningNumbers,
		Duration:       duration,
		Attempts:       int64(drawCount),
	}, nil
}

// 複数回のシミュレーションを実行して統計を取る
func (se *SimulationEngine) RunMultipleSimulations(count int, userNumbers []int) (*SimulationStats, error) {
	if count <= 0 {
		return nil, fmt.Errorf("シミュレーション回数は1以上である必要があります")
	}

	if count > 1000 {
		return nil, fmt.Errorf("シミュレーション回数は1000回以下に制限されています")
	}

	results := make([]SimulationResult, 0, count)
	totalDuration := time.Duration(0)
	minDraws := 0
	maxDraws := 0
	totalDraws := 0

	for i := 0; i < count; i++ {
		result, err := se.RunSimulation(userNumbers)
		if err != nil {
			return nil, fmt.Errorf("シミュレーション %d 回目で失敗: %w", i+1, err)
		}

		results = append(results, *result)
		totalDuration += result.Duration
		totalDraws += result.DrawCount

		if i == 0 || result.DrawCount < minDraws {
			minDraws = result.DrawCount
		}
		if result.DrawCount > maxDraws {
			maxDraws = result.DrawCount
		}
	}

	// 中央値の計算
	drawCounts := make([]int, len(results))
	for i, r := range results {
		drawCounts[i] = r.DrawCount
	}
	sort.Ints(drawCounts)
	median := drawCounts[len(drawCounts)/2]

	avgDraws := float64(totalDraws) / float64(count)

	return &SimulationStats{
		Simulations:   count,
		MinDraws:      minDraws,
		MaxDraws:      maxDraws,
		AvgDraws:      avgDraws,
		MedianDraws:   median,
		TotalDuration: totalDuration,
		Results:       results,
	}, nil
}

// ユーザーの数字をランダムに生成する
func (se *SimulationEngine) GenerateUserNumbers() ([]int, error) {
	return se.randomEngine.GenerateRandomCombination()
}
