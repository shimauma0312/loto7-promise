// ロト7の統計分析予測機能を提供する
//
// 統計分析（ホット/コールド分析、パリティ、大小バランス、合計値、
// 十の位グループ、一の位、引っ張り数字、ボーナス周辺）を組み合わせ、
// フィルタリングとスコアリングで最適な組み合わせを生成する
package prediction

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/shimauma0312/loto7-promise/internal/result"
)

const (
	// 候補プール生成数（多いほど精度が上がるが処理時間増）
	defaultPoolSize = 30000
	// 直近ウィンドウのデフォルト（ホット/コールド判定）
	defaultRecentWindow = 10
)

// Engine は予測生成の中核エンジン
type Engine struct {
	rng      *rand.Rand
	analyzer *Analyzer
	results  [][]string
}

// NewEngine は新しい予測エンジンを作成する
func NewEngine() *Engine {
	return &Engine{
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
		analyzer: &Analyzer{},
	}
}

// LoadData は過去の抽選データを読み込む
func (e *Engine) LoadData(count int) error {
	e.results = result.GetResult(count)
	if len(e.results) == 0 {
		return fmt.Errorf("データを取得できませんでした: count=%d", count)
	}
	return nil
}

// Generate は設定に従って予測組み合わせを生成する
func (e *Engine) Generate(cfg PredictionConfig) (PredictionResult, error) {
	if len(e.results) == 0 {
		return PredictionResult{}, fmt.Errorf("データが読み込まれていません")
	}

	// 直近ウィンドウを決定
	window := defaultRecentWindow
	if cfg.History < window {
		window = cfg.History
	}

	// 統計分析
	data := e.analyzer.Analyze(e.results, window)

	// 候補プール生成
	pool := generateCandidatePool(e.rng, defaultPoolSize)

	// ハードフィルタリング
	filtered := applyHardFilters(pool, data)
	if len(filtered) == 0 {
		// 全フィルター失敗の場合はプール全体を使用
		filtered = pool
	}

	// スコアリングと降順ソート
	type scoredEntry struct {
		combo  []int
		total  float64
		detail ScoreDetail
	}

	scoredList := make([]scoredEntry, len(filtered))
	for i, combo := range filtered {
		total, detail := Score(combo, data, cfg.Weights)
		scoredList[i] = scoredEntry{combo: combo, total: total, detail: detail}
	}

	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].total > scoredList[j].total
	})

	// 上位 N 件を重複なしで抽出
	need := cfg.Count
	if need > len(scoredList) {
		need = len(scoredList)
	}

	combinations := make([]ScoredCombination, 0, need)
	usedKeys := make(map[string]bool, need)

	for _, s := range scoredList {
		if len(combinations) >= need {
			break
		}
		key := comboKey(s.combo)
		if !usedKeys[key] {
			usedKeys[key] = true
			combinations = append(combinations, ScoredCombination{
				ID:          len(combinations) + 1,
				Numbers:     s.combo,
				TotalScore:  s.total,
				ScoreDetail: s.detail,
			})
		}
	}

	return PredictionResult{
		Combinations: combinations,
		AnalysisInfo: PredictionAnalysisInfo{
			AnalyzedDraws:    data.AnalyzedDraws,
			HotNumbers:       data.HotNumbers,
			ColdNumbers:      data.ColdNumbers,
			LastDrawNumbers:  data.LastDrawNumbers,
			LastBonusNumbers: data.LastBonusNumbers,
			LastDrawSum:      data.LastDrawSum,
			SumTrend:         sumTrendLabel(data.SumTrend),
			GeneratedAt:      time.Now(),
		},
		Config: cfg,
	}, nil
}

// sumTrendLabel はSumTrendを日本語ラベルに変換する
func sumTrendLabel(trend int) string {
	switch trend {
	case -1:
		return "2回連続上昇中（下降傾向を優先）"
	case 1:
		return "2回連続下降中（上昇傾向を優先）"
	default:
		return "変動なし"
	}
}
