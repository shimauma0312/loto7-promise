package test

import (
	"math"
	"testing"

	"github.com/shimauma0312/loto7-promise/internal/recommendation"
)

var bayesTestResults = [][]string{
	{"1", "5", "10", "15", "20", "25", "30"},
	{"3", "5", "12", "18", "22", "28", "35"},
	{"1", "7", "10", "16", "21", "27", "33"},
	{"2", "5", "11", "14", "23", "29", "37"},
	{"1", "5", "8", "17", "20", "31", "36"},
}

func TestBayesianEstimator_Build_Appearances(t *testing.T) {
	b := recommendation.NewBayesianEstimator(1.0)
	b.Build(bayesTestResults)

	// 1 は 3 回出現しているので、事前確率は 5 回中 3 回
	prob1 := b.ConditionalProb(1, nil)
	prob2 := b.ConditionalProb(2, nil)

	if prob1 <= prob2 {
		t.Errorf("P(1) = %v、P(2) = %v: 出現回数の多い 1 のほうが高いはず", prob1, prob2)
	}
}

func TestBayesianEstimator_ConditionalProb_WithSelected(t *testing.T) {
	b := recommendation.NewBayesianEstimator(1.0)
	b.Build(bayesTestResults)

	// 1 と 5 は頻繁に同時出現しているので P(5|{1}) > P(3|{1}) を期待
	probGiven1WithNum5 := b.ConditionalProb(5, []int{1})
	probGiven1WithNum3 := b.ConditionalProb(3, []int{1})

	if probGiven1WithNum5 <= probGiven1WithNum3 {
		t.Errorf("P(5|{1})=%v <= P(3|{1})=%v: 1と5は共起が多いので前者が高いはず", probGiven1WithNum5, probGiven1WithNum3)
	}
}

func TestBayesianEstimator_ConditionalProb_AlwaysPositive(t *testing.T) {
	b := recommendation.NewBayesianEstimator(1.0)
	b.Build(bayesTestResults)

	for num := 1; num <= 37; num++ {
		prob := b.ConditionalProb(num, []int{5, 10})
		if prob <= 0 {
			t.Errorf("ConditionalProb(%d, {5,10}) = %v, want > 0 (ラプラス平滑化で常に正)", num, prob)
		}
		if math.IsNaN(prob) || math.IsInf(prob, 0) {
			t.Errorf("ConditionalProb(%d) = %v, want finite value", num, prob)
		}
	}
}

func TestBayesianEstimator_EmptyResults(t *testing.T) {
	b := recommendation.NewBayesianEstimator(1.0)
	b.Build([][]string{})

	// データなし時はフォールバック値を返す
	prob := b.ConditionalProb(7, nil)
	if prob <= 0 {
		t.Errorf("空データ時の ConditionalProb = %v, want > 0", prob)
	}
}
