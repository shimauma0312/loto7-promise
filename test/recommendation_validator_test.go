package test

import (
	"testing"

	"github.com/shimauma0312/loto7-promise/internal/recommendation"
)

func TestZoneDistributionValidator_Valid(t *testing.T) {
	v := &recommendation.ZoneDistributionValidator{}
	// zone1: 3個 (1,2,3), zone2: 2個 (14,15), zone3: 2個 (27,28)
	if !v.Validate([]int{1, 2, 3, 14, 15, 27, 28}) {
		t.Error("valid combination rejected")
	}
}

func TestZoneDistributionValidator_Zone1TooFew(t *testing.T) {
	v := &recommendation.ZoneDistributionValidator{}
	// zone1: 2個 → 不正
	if v.Validate([]int{1, 2, 14, 15, 16, 27, 28}) {
		t.Error("zone1=2 should be invalid")
	}
}

func TestZoneDistributionValidator_Zone3TooFew(t *testing.T) {
	v := &recommendation.ZoneDistributionValidator{}
	// zone3: 1個 → 不正
	if v.Validate([]int{1, 2, 3, 14, 15, 16, 27}) {
		t.Error("zone3=1 should be invalid")
	}
}

func TestZoneDistributionValidator_ZoneEmpty(t *testing.T) {
	v := &recommendation.ZoneDistributionValidator{}
	// zone2が空 → 不正
	if v.Validate([]int{1, 2, 3, 4, 27, 28, 29}) {
		t.Error("empty zone2 should be invalid")
	}
}

func TestOddEvenRatioValidator_Valid43(t *testing.T) {
	v := &recommendation.OddEvenRatioValidator{}
	// 奇4:偶3
	if !v.Validate([]int{1, 3, 5, 7, 2, 4, 6}) {
		t.Error("odd=4 even=3 should be valid")
	}
}

func TestOddEvenRatioValidator_Valid34(t *testing.T) {
	v := &recommendation.OddEvenRatioValidator{}
	// 奇3:偶4
	if !v.Validate([]int{1, 3, 5, 2, 4, 6, 8}) {
		t.Error("odd=3 even=4 should be valid")
	}
}

func TestOddEvenRatioValidator_Invalid55(t *testing.T) {
	v := &recommendation.OddEvenRatioValidator{}
	// 奇5:偶2 → 不正（7まで行くと注意: 1+3+5+7+9 = 5奇、2+4 = 2偶）
	// 実際: 5+2=7
	if v.Validate([]int{1, 3, 5, 7, 9, 2, 4}) {
		t.Error("odd=5 even=2 should be invalid")
	}
}

func TestTotalSumValidator_Valid(t *testing.T) {
	v := &recommendation.TotalSumValidator{}
	// 合計 = 1+14+15+20+21+27+28 = 126
	if !v.Validate([]int{1, 14, 15, 20, 21, 27, 28}) {
		t.Error("sum=126 should be valid")
	}
}

func TestTotalSumValidator_TooLow(t *testing.T) {
	v := &recommendation.TotalSumValidator{}
	// 合計 = 1+2+3+4+5+6+7 = 28
	if v.Validate([]int{1, 2, 3, 4, 5, 6, 7}) {
		t.Error("sum=28 should be invalid (< TotalSumMin)")
	}
}

func TestTotalSumValidator_TooHigh(t *testing.T) {
	v := &recommendation.TotalSumValidator{}
	// 合計 = 31+32+33+34+35+36+37 = 238
	if v.Validate([]int{31, 32, 33, 34, 35, 36, 37}) {
		t.Error("sum=238 should be invalid (> TotalSumMax)")
	}
}

func TestAverageValueValidator_Valid(t *testing.T) {
	v := &recommendation.AverageValueValidator{}
	// 合計 = 126, 平均 = 18.0
	if !v.Validate([]int{1, 14, 15, 20, 21, 27, 28}) {
		t.Error("avg=18.0 should be valid")
	}
}

func TestSmoothnessValidator_Valid(t *testing.T) {
	v := &recommendation.SmoothnessValidator{}
	// 均等な差分で滑らか
	if !v.Validate([]int{5, 10, 15, 20, 25, 30, 35}) {
		t.Error("evenly spaced combination should be valid")
	}
}

func TestCloseNumberPairValidator_HasPair(t *testing.T) {
	v := &recommendation.CloseNumberPairValidator{}
	// 差1のペア (1,2) あり
	if !v.Validate([]int{1, 2, 10, 15, 20, 25, 30}) {
		t.Error("combination with close pair should be valid")
	}
}

func TestCloseNumberPairValidator_NoPair(t *testing.T) {
	v := &recommendation.CloseNumberPairValidator{}
	// すべての差が4以上
	if v.Validate([]int{1, 5, 10, 15, 20, 25, 30}) {
		t.Error("combination without close pair should be invalid")
	}
}

func TestCompositeValidator_AllPass(t *testing.T) {
	v := recommendation.NewCompositeValidator(
		&recommendation.TotalSumValidator{},
		&recommendation.OddEvenRatioValidator{},
	)
	// 合計=126, 奇3偶4 → 両方通過
	if !v.Validate([]int{1, 14, 15, 20, 21, 27, 28}) {
		t.Error("composite validator should pass")
	}
}

func TestCompositeValidator_OneFails(t *testing.T) {
	v := recommendation.NewCompositeValidator(
		&recommendation.TotalSumValidator{},
		&recommendation.OddEvenRatioValidator{},
	)
	// 合計=28 (低すぎ) → TotalSumValidator が弾く
	if v.Validate([]int{1, 2, 3, 4, 5, 6, 7}) {
		t.Error("composite validator should fail when one validator rejects")
	}
}
