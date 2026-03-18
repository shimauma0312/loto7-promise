// スコアリング重みをバックテストで最適化し、weights.json に保存するコマンド
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/shimauma0312/loto7-promise/internal/recommendation"
	"github.com/shimauma0312/loto7-promise/internal/result"
)

func main() {
	var (
		historyCount = flag.Int("history", 300, "使用する過去抽選データ数")
		lookback     = flag.Int("lookback", 100, "1回の予測に使うトレーニング回数")
		iterations   = flag.Int("iter", 500, "最適化反復回数")
		trialCount   = flag.Int("trials", 10, "評価時に生成する推腐数")
		outputFile   = flag.String("out", recommendation.WeightsFilePath, "重みの保存先ファイル")
		loadFile     = flag.String("load", "", "初期重みを読み込む JSON ファイル（省略時はデフォルト重みから開始）")
	)
	flag.Parse()

	fmt.Println("=== スコアリング重みの最適化 ===")
	fmt.Println()

	fmt.Printf("データ取得中（直近%d回）...\n", *historyCount)
	results := result.FilterAbnormalResults(result.GetResult(*historyCount))
	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "エラー: データを取得できませんでした")
		os.Exit(1)
	}
	testCount := len(results) - *lookback
	if testCount <= 0 {
		fmt.Fprintf(os.Stderr, "エラー: データ数(%d)が lookback(%d)以下です\n", len(results), *lookback)
		os.Exit(1)
	}
	fmt.Printf("取得件数: %d 回  /  評価データ: %d 回（lookback=%d）\n\n", len(results), testCount, *lookback)

	// 初期重みの決定
	initial := recommendation.DefaultScorerWeights()
	if *loadFile != "" {
		wf, err := recommendation.LoadWeights(*loadFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "警告: %s を読み込めませんでした（デフォルト重みを使用）: %v\n", *loadFile, err)
		} else {
			initial = wf.Weights
			fmt.Printf("初期重みを %s から読み込みました\n", *loadFile)
		}
	}

	backtester := recommendation.NewBacktester(results, *lookback)
	backtester.SetTrialCount(*trialCount)

	evalRNG := rand.New(rand.NewSource(42))
	initialScore := backtester.Evaluate(initial, evalRNG)
	fmt.Printf("初期スコア（平均ベストヒット数）: %.4f\n\n", initialScore)

	optimizer := recommendation.NewHillClimbOptimizer(backtester, time.Now().UnixNano())

	fmt.Printf("最適化開始（iterations=%d, trials=%d）...\n", *iterations, *trialCount)
	reportEvery := *iterations / 10
	if reportEvery == 0 {
		reportEvery = 1
	}

	optimized := optimizer.Optimize(initial, *iterations, func(i int, best float64) {
		if (i+1)%reportEvery == 0 {
			fmt.Printf("  [%4d/%d] best = %.4f\n", i+1, *iterations, best)
		}
	})

	evalRNG = rand.New(rand.NewSource(42))
	finalScore := backtester.Evaluate(optimized, evalRNG)
	fmt.Printf("\n最適化後スコア: %.4f  （改善: %+.4f）\n\n", finalScore, finalScore-initialScore)

	fmt.Println("=== 最適化された重み ===")
	fmt.Printf("Recent1Penalty:      %9.4f  （default: %.4f）\n", optimized.Recent1Penalty, initial.Recent1Penalty)
	fmt.Printf("Recent2Penalty:      %9.4f  （default: %.4f）\n", optimized.Recent2Penalty, initial.Recent2Penalty)
	fmt.Printf("Within50Boost:       %9.4f  （default: %.4f）\n", optimized.Within50Boost, initial.Within50Boost)
	fmt.Printf("LowFrequencyPenalty: %9.4f  （default: %.4f）\n", optimized.LowFrequencyPenalty, initial.LowFrequencyPenalty)
	fmt.Printf("PositionFrequencyWeight: %7.4f  （default: %.4f）\n", optimized.PositionFrequencyWeight, initial.PositionFrequencyWeight)
	fmt.Println("PositionPenaltyScale:")
	for i, scale := range optimized.PositionPenaltyScale {
		fmt.Printf("  pos%d: %6.4f  （default: %.4f）\n", i, scale, initial.PositionPenaltyScale[i])
	}

	currentDraw := result.NewNumber()
	if err := recommendation.SaveWeights(optimized, currentDraw, *outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "\n警告: %s への保存に失敗しました: %v\n", *outputFile, err)
	} else {
		fmt.Printf("\n重みを保存しました: %s （抽選第%d回による学習）\n", *outputFile, currentDraw)
		fmt.Println("次回 recommendation を実行すると自動的にこの重みが使用されます。")
	}
}
