// ランダムフォレストによる多特徴量スコアリングを提供する
//
// 集計統計系の特徴量（出現頻度比・直近トレンド・出現間隔等）を使い、
// BayesianEstimator のペア相関系と役割を分担してアンサンブルする。
package recommendation

import (
	"math"
	"math/rand"
	"sort"
)

const (
	rfTreeCount   = 15 // 決定木の数
	rfMaxDepth    = 4  // 木の最大深さ
	rfMaxWindows  = 50 // 訓練ウィンドウの最大数（速度と精度のバランス）
	rfMaxThresholds = 10 // 分割候補閾値の最大数（速度最適化）
	featuresCount = 7  // 特徴量の数
)

// extractFeatures は数字の集計統計特徴ベクトルを抽出する。
//
// ペア共起情報は含まず、BayesianEstimator との特徴量重複を避ける。
func extractFeatures(num int, results [][]string, lookback int) [featuresCount]float64 {
	var f [featuresCount]float64
	if lookback > len(results) {
		lookback = len(results)
	}
	if lookback == 0 {
		return f
	}

	// f[0]: freq100 - 直近 lookback 回の出現頻度（期待値で正規化）
	freq := 0
	for i := 0; i < lookback; i++ {
		for _, n := range parseDraw(results[i]) {
			if n == num {
				freq++
				break
			}
		}
	}
	expectedFreq := float64(lookback) * float64(lotoDrawCount) / float64(lotoTotalNumbers)
	if expectedFreq > 0 {
		f[0] = float64(freq) / expectedFreq
	}

	// f[1]: freq20 - 直近 20 回の出現頻度（短期傾向）
	shortLookback := 20
	if shortLookback > lookback {
		shortLookback = lookback
	}
	freq20 := 0
	for i := 0; i < shortLookback; i++ {
		for _, n := range parseDraw(results[i]) {
			if n == num {
				freq20++
				break
			}
		}
	}
	expectedFreq20 := float64(shortLookback) * float64(lotoDrawCount) / float64(lotoTotalNumbers)
	if expectedFreq20 > 0 {
		f[1] = float64(freq20) / expectedFreq20
	}

	// f[2]: gap - 最後に出現してからの抽選数を 0-1 に正規化（大きいほど長期未出現）
	gap := lookback
	for i := 0; i < lookback; i++ {
		found := false
		for _, n := range parseDraw(results[i]) {
			if n == num {
				gap = i
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	f[2] = float64(gap) / float64(lookback)

	// f[3]: isOdd - 奇数なら 1.0
	if num%2 == 1 {
		f[3] = 1.0
	}

	// f[4]: zone - ゾーン (1-13=0.33, 14-26=0.67, 27-37=1.0)
	if num <= 13 {
		f[4] = 1.0 / 3.0
	} else if num <= 26 {
		f[4] = 2.0 / 3.0
	} else {
		f[4] = 1.0
	}

	// f[5]: numNorm - 番号を 0-1 に正規化
	f[5] = float64(num) / float64(lotoTotalNumbers)

	// f[6]: intervalMean - 出現間隔の平均の逆数（出現周期が短いほど高スコア）
	lastIdx := -1
	totalInterval := 0
	intervalCount := 0
	for i := 0; i < lookback; i++ {
		for _, n := range parseDraw(results[i]) {
			if n == num {
				if lastIdx >= 0 {
					totalInterval += i - lastIdx
					intervalCount++
				}
				lastIdx = i
				break
			}
		}
	}
	if intervalCount > 0 {
		meanInterval := float64(totalInterval) / float64(intervalCount)
		f[6] = 1.0 / (1.0 + meanInterval)
	}

	return f
}

// trainSample は特徴量と正解ラベルのペアを表す
type trainSample struct {
	features [featuresCount]float64
	label    float64
}

// treeNode は決定木のノードを表す
type treeNode struct {
	featureIdx int
	threshold  float64
	left       *treeNode
	right      *treeNode
	prediction float64
	isLeaf     bool
}

// predict は特徴ベクトルの出現確率予測値を返す
func (n *treeNode) predict(f [featuresCount]float64) float64 {
	if n.isLeaf {
		return n.prediction
	}
	if f[n.featureIdx] < n.threshold {
		return n.left.predict(f)
	}
	return n.right.predict(f)
}

// buildDecisionTree は CART アルゴリズム（Gini 不純度最小化）で決定木を構築する
func buildDecisionTree(samples []trainSample, depth, maxDepth, featureSubsetSize int, rng *rand.Rand) *treeNode {
	if depth >= maxDepth || len(samples) < 2 {
		return makeLeafNode(samples)
	}
	parentGini := calcGini(samples)
	if parentGini < 1e-9 {
		return makeLeafNode(samples)
	}

	featureIndices := selectRandomFeatures(featureSubsetSize, featuresCount, rng)
	bestGain := 0.0
	bestFeatIdx := -1
	bestThreshold := 0.0
	var bestLeft, bestRight []trainSample

	for _, fi := range featureIndices {
		for _, thresh := range calcThresholds(samples, fi) {
			left, right := splitByThreshold(samples, fi, thresh)
			if len(left) == 0 || len(right) == 0 {
				continue
			}
			n := float64(len(samples))
			gain := parentGini -
				float64(len(left))/n*calcGini(left) -
				float64(len(right))/n*calcGini(right)
			if gain > bestGain {
				bestGain = gain
				bestFeatIdx = fi
				bestThreshold = thresh
				bestLeft = left
				bestRight = right
			}
		}
	}

	if bestFeatIdx == -1 {
		return makeLeafNode(samples)
	}

	return &treeNode{
		featureIdx: bestFeatIdx,
		threshold:  bestThreshold,
		left:       buildDecisionTree(bestLeft, depth+1, maxDepth, featureSubsetSize, rng),
		right:      buildDecisionTree(bestRight, depth+1, maxDepth, featureSubsetSize, rng),
	}
}

func makeLeafNode(samples []trainSample) *treeNode {
	pred := 0.0
	if len(samples) > 0 {
		for _, s := range samples {
			pred += s.label
		}
		pred /= float64(len(samples))
	}
	return &treeNode{isLeaf: true, prediction: pred}
}

func calcGini(samples []trainSample) float64 {
	if len(samples) == 0 {
		return 0
	}
	pos := 0.0
	for _, s := range samples {
		pos += s.label
	}
	p := pos / float64(len(samples))
	return 2 * p * (1 - p)
}

func selectRandomFeatures(subsetSize, total int, rng *rand.Rand) []int {
	perm := rng.Perm(total)
	if subsetSize > total {
		subsetSize = total
	}
	return perm[:subsetSize]
}

// calcThresholds は quantile サブサンプリングで rfMaxThresholds 個の分割候補閾値を返す
func calcThresholds(samples []trainSample, fi int) []float64 {
	vals := make([]float64, len(samples))
	for i, s := range samples {
		vals[i] = s.features[fi]
	}
	sort.Float64s(vals)

	// 重複排除
	uniq := vals[:0]
	for i, v := range vals {
		if i == 0 || v != vals[i-1] {
			uniq = append(uniq, v)
		}
	}
	if len(uniq) <= 1 {
		return nil
	}

	// 均等間隔でサブサンプリング（速度最適化）
	step := len(uniq) / rfMaxThresholds
	if step < 1 {
		step = 1
	}
	var thresholds []float64
	for i := 0; i < len(uniq)-1; i += step {
		mid := (uniq[i] + uniq[i+1]) / 2.0
		thresholds = append(thresholds, mid)
		if len(thresholds) >= rfMaxThresholds {
			break
		}
	}
	return thresholds
}

func splitByThreshold(samples []trainSample, fi int, threshold float64) (left, right []trainSample) {
	for _, s := range samples {
		if s.features[fi] < threshold {
			left = append(left, s)
		} else {
			right = append(right, s)
		}
	}
	return
}

// RandomForest は複数の決定木アンサンブルで数字の出現確率を予測する
type RandomForest struct {
	trees     []*treeNode
	treeCount int
	maxDepth  int
	rng       *rand.Rand
}

// NewRandomForest は RandomForest を作成する
func NewRandomForest(treeCount, maxDepth int, rng *rand.Rand) *RandomForest {
	return &RandomForest{
		trees:     make([]*treeNode, 0, treeCount),
		treeCount: treeCount,
		maxDepth:  maxDepth,
		rng:       rng,
	}
}

// Train は過去のロト7抽選データでランダムフォレストを訓練する。
//
// 各 draw をターゲットとし、直前 lookback 回の集計統計特徴量で訓練サンプルを生成する。
// 直近 rfMaxWindows ウィンドウのみを使用して訓練速度を確保する。
func (rf *RandomForest) Train(results [][]string, lookback int) {
	if len(results) <= lookback {
		return
	}

	maxWindowCount := len(results) - lookback - 1
	if maxWindowCount > rfMaxWindows {
		maxWindowCount = rfMaxWindows
	}

	allSamples := make([]trainSample, 0, maxWindowCount*lotoTotalNumbers)
	for i := 0; i < maxWindowCount; i++ {
		target := parseDraw(results[i])
		train := results[i+1 : i+1+lookback]
		targetSet := make(map[int]bool, len(target))
		for _, n := range target {
			targetSet[n] = true
		}
		for num := 1; num <= lotoTotalNumbers; num++ {
			label := 0.0
			if targetSet[num] {
				label = 1.0
			}
			allSamples = append(allSamples, trainSample{
				features: extractFeatures(num, train, lookback),
				label:    label,
			})
		}
	}

	if len(allSamples) == 0 {
		return
	}

	featureSubsetSize := int(math.Sqrt(float64(featuresCount))) + 1
	rf.trees = make([]*treeNode, 0, rf.treeCount)
	n := len(allSamples)
	for t := 0; t < rf.treeCount; t++ {
		// ブートストラップサンプリング
		bootstrap := make([]trainSample, n)
		for i := range bootstrap {
			bootstrap[i] = allSamples[rf.rng.Intn(n)]
		}
		tree := buildDecisionTree(bootstrap, 0, rf.maxDepth, featureSubsetSize, rf.rng)
		rf.trees = append(rf.trees, tree)
	}
}

// Predict は num が次の抽選に出現する予測確率を返す
func (rf *RandomForest) Predict(num int, results [][]string, lookback int) float64 {
	if len(rf.trees) == 0 {
		return float64(lotoDrawCount) / float64(lotoTotalNumbers)
	}
	f := extractFeatures(num, results, lookback)
	sum := 0.0
	for _, t := range rf.trees {
		sum += t.predict(f)
	}
	return sum / float64(len(rf.trees))
}

// IsTrained は訓練済みかどうかを返す
func (rf *RandomForest) IsTrained() bool {
	return len(rf.trees) > 0
}
