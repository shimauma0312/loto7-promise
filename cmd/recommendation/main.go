package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shimauma0312/loto7-promise/internal/recommendation"
)

//  1. フラグ解析: -history（分析回数）、-count（推薦数）などのオプションを受け取る
//  2. RecommendationEngine 初期化: Analyzer / WeightedScorer / ZoneBasedSelector / TextFormatter を組み立てる
//  3. LoadData: 過去N回の抽選結果を取得し、統計分析（出現頻度・ホット数字・復活候補）と
//     7種のバリデータ（ゾーン分布・奇偶比・合計値・近接ペアなど）を設定する
//  4. GenerateRecommendations: ゾーン別の重み付き選択で候補を生成し、全制約を満たす組み合わせを採用する
//  5. FormatRecommendations: 推薦番号・統計データをテキスト整形して標準出力へ出力する
func main() {
	var (
		history      = flag.Int("history", 100, "分析する過去の抽選回数")
		count        = flag.Int("count", 5, "生成する推薦組み合わせの数")
		showScores   = flag.Bool("scores", false, "各数字のスコア詳細を表示")
		help         = flag.Bool("help", false, "ヘルプを表示")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// 推薦エンジンの設定
	config := recommendation.DefaultConfig()
	config.MaxRecommendations = *count

	// エンジンを作成
	engine := recommendation.NewRecommendationEngine(config)

	// データを読み込み
	fmt.Printf("過去%d回分のデータを分析中...\n", *history)
	if err := engine.LoadData(*history); err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	// 推薦組み合わせを生成
	recommendations, err := engine.GenerateRecommendations()
	if err != nil {
		fmt.Printf("推薦生成エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を表示
	result := engine.FormatRecommendations(recommendations)
	fmt.Print(result)

	// スコア詳細を表示（オプション）
	if *showScores {
		fmt.Println("\n=== 数字別スコア詳細 ===")
		showDetailedScores(engine)
	}

	// 人間の感覚による推薦ロジックの説明
	showRecommendationLogic()
}

func showHelp() {
	fmt.Println("ロト7推薦ツール")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  recommendation [オプション]")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println("  -history int    分析する過去の抽選回数 (デフォルト: 100)")
	fmt.Println("  -count int      生成する推薦組み合わせの数 (デフォルト: 5)")
	fmt.Println("  -scores         各数字のスコア詳細を表示")
	fmt.Println("  -help           このヘルプを表示")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  recommendation -history 150 -count 3")
	fmt.Println("  recommendation -scores")
}

func showDetailedScores(engine *recommendation.RecommendationEngine) {
	// この関数は推薦エンジンからスコア情報を取得して表示します
	// 実装はrecommendationパッケージでスコア情報を公開する必要があります
	fmt.Println("（スコア詳細表示は今後の実装で対応予定）")
}

func showRecommendationLogic() {
	fmt.Println()
	fmt.Println("=== 推薦ロジック ===")
	fmt.Println("設定と統計に基づき候補を自動生成します。")
}
