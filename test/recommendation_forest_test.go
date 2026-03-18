package test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/shimauma0312/loto7-promise/internal/recommendation"
)

// buildRFTestResults generates n dummy draw results for RF testing
func buildRFTestResults(n int, seed int64) [][]string {
	rng := rand.New(rand.NewSource(seed))
	results := make([][]string, n)
	for i := range results {
		picked := rng.Perm(37)[:7]
		draw := make([]string, 7)
		for j, p := range picked {
			draw[j] = fmt.Sprintf("%d", p+1)
		}
		results[i] = draw
	}
	return results
}

func TestRandomForest_Train_IsTrained(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	rf := recommendation.NewRandomForest(5, 3, rng)

	if rf.IsTrained() {
		t.Error("train: IsTrained() should be false before training")
	}

	results := buildRFTestResults(30, 99)
	rf.Train(results, 20)

	if !rf.IsTrained() {
		t.Error("train: IsTrained() should be true after Train()")
	}
}

func TestRandomForest_Predict_ValidRange(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	rf := recommendation.NewRandomForest(5, 3, rng)
	results := buildRFTestResults(30, 99)
	rf.Train(results, 20)

	for num := 1; num <= 37; num++ {
		prob := rf.Predict(num, results, 20)
		if prob < 0 || prob > 1 {
			t.Errorf("Predict(%d) = %v, want in [0, 1]", num, prob)
		}
		if math.IsNaN(prob) || math.IsInf(prob, 0) {
			t.Errorf("Predict(%d) = %v, want finite", num, prob)
		}
	}
}

func TestRandomForest_Predict_Untrained_ReturnsDefault(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	rf := recommendation.NewRandomForest(5, 3, rng)
	results := buildRFTestResults(30, 99)

	// untrained: should return prior probability 7/37
	prob := rf.Predict(10, results, 20)
	expected := float64(7) / float64(37)
	if math.Abs(prob-expected) > 1e-9 {
		t.Errorf("untrained Predict = %v, want %v", prob, expected)
	}
}

func TestRandomForest_Train_InsufficientData_NoOp(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	rf := recommendation.NewRandomForest(5, 3, rng)

	// len(results) <= lookback: Train should be a no-op
	results := buildRFTestResults(20, 99)
	rf.Train(results, 20)

	if rf.IsTrained() {
		t.Error("insufficient data: IsTrained() should stay false")
	}
}
