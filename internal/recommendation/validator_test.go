package recommendation

import "testing"

// TestZoneDistributionValidator_ValidCombination 有効なゾーン分布をテスト
func TestZoneDistributionValidator_ValidCombination(t *testing.T) {
	validator := &ZoneDistributionValidator{}

	tests := []struct {
		name        string
		combination []int
		want        bool
	}{
		{
			name:        "ゾーン1:3個, ゾーン2:2個, ゾーン3:2個",
			combination: []int{1, 5, 10, 15, 20, 30, 35},
			want:        true,
		},
		{
			name:        "ゾーン1:4個, ゾーン2:1個, ゾーン3:2個",
			combination: []int{1, 5, 10, 13, 20, 30, 35},
			want:        true,
		},
		{
			name:        "ゾーン1:3個, ゾーン2:1個, ゾーン3:3個",
			combination: []int{1, 5, 10, 20, 30, 35, 37},
			want:        true,
		},
		{
			name:        "ゾーン1が空",
			combination: []int{14, 15, 20, 25, 30, 35, 37},
			want:        false,
		},
		{
			name:        "ゾーン3が空",
			combination: []int{1, 5, 10, 15, 20, 22, 25},
			want:        false,
		},
		{
			name:        "ゾーン1が2個（少なすぎ）",
			combination: []int{1, 5, 15, 20, 22, 30, 35},
			want:        false,
		},
		{
			name:        "ゾーン1が5個（多すぎ）",
			combination: []int{1, 2, 3, 4, 5, 30, 35},
			want:        false,
		},
		{
			name:        "ゾーン3が1個（少なすぎ）",
			combination: []int{1, 5, 10, 15, 20, 22, 30},
			want:        false,
		},
		{
			name:        "ゾーン3が4個（多すぎ）",
			combination: []int{1, 5, 10, 30, 32, 35, 37},
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.Validate(tt.combination)
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestOddEvenRatioValidator_ValidCombination 有効な奇数偶数比率をテスト
func TestOddEvenRatioValidator_ValidCombination(t *testing.T) {
	validator := &OddEvenRatioValidator{}

	tests := []struct {
		name        string
		combination []int
		want        bool
	}{
		{
			name:        "奇4:偶3",
			combination: []int{1, 3, 5, 7, 2, 4, 6},
			want:        true,
		},
		{
			name:        "奇3:偶4",
			combination: []int{1, 3, 5, 2, 4, 6, 8},
			want:        true,
		},
		{
			name:        "奇5:偶2",
			combination: []int{1, 3, 5, 7, 9, 2, 4},
			want:        false,
		},
		{
			name:        "奇2:偶5",
			combination: []int{1, 3, 2, 4, 6, 8, 10},
			want:        false,
		},
		{
			name:        "全て奇数",
			combination: []int{1, 3, 5, 7, 9, 11, 13},
			want:        false,
		},
		{
			name:        "全て偶数",
			combination: []int{2, 4, 6, 8, 10, 12, 14},
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.Validate(tt.combination)
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTotalSumValidator_ValidCombination 有効な合計値をテスト
func TestTotalSumValidator_ValidCombination(t *testing.T) {
	validator := &TotalSumValidator{}

	tests := []struct {
		name        string
		combination []int
		want        bool
	}{
		{
			name:        "合計100（下限）",
			combination: []int{1, 2, 3, 4, 20, 30, 40}, // 合計100
			want:        true,
		},
		{
			name:        "合計160（上限）",
			combination: []int{10, 15, 20, 25, 30, 30, 30}, // 合計160
			want:        true,
		},
		{
			name:        "合計130（範囲内）",
			combination: []int{5, 10, 15, 20, 25, 27, 28}, // 合計130
			want:        true,
		},
		{
			name:        "合計99（下限未満）",
			combination: []int{1, 2, 3, 4, 5, 6, 78}, // 合計99
			want:        false,
		},
		{
			name:        "合計161（上限超過）",
			combination: []int{20, 21, 22, 23, 24, 25, 26}, // 合計161
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.Validate(tt.combination)
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAverageValueValidator_ValidCombination 有効な平均値をテスト
func TestAverageValueValidator_ValidCombination(t *testing.T) {
	validator := &AverageValueValidator{}

	tests := []struct {
		name        string
		combination []int
		want        bool
	}{
		{
			name:        "平均16.0（下限）",
			combination: []int{16, 16, 16, 16, 16, 16, 16}, // 平均16.0
			want:        true,
		},
		{
			name:        "平均23.0（上限）",
			combination: []int{23, 23, 23, 23, 23, 23, 23}, // 平均23.0
			want:        true,
		},
		{
			name:        "平均19.5（範囲内）",
			combination: []int{10, 15, 18, 20, 22, 25, 27}, // 平均約19.57
			want:        true,
		},
		{
			name:        "平均15.9（下限未満）",
			combination: []int{10, 11, 12, 13, 20, 21, 24}, // 平均約15.86
			want:        false,
		},
		{
			name:        "平均23.1（上限超過）",
			combination: []int{20, 21, 22, 23, 24, 25, 27}, // 平均約23.14
			want:        false,
		},
		{
			name:        "空の組み合わせ",
			combination: []int{},
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.Validate(tt.combination)
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSmoothnessValidator_ValidCombination 数字の並び滑らかさをテスト
func TestSmoothnessValidator_ValidCombination(t *testing.T) {
	validator := &SmoothnessValidator{}

	tests := []struct {
		name        string
		combination []int
		want        bool
	}{
		{
			name:        "滑らかな並び",
			combination: []int{1, 5, 10, 15, 20, 25, 30},
			want:        true,
		},
		{
			name:        "やや滑らか",
			combination: []int{2, 8, 12, 18, 24, 29, 35},
			want:        true,
		},
		{
			name:        "ジグザグが激しい",
			combination: []int{1, 2, 20, 21, 22, 36, 37},
			want:        false,
		},
		{
			name:        "1個のみ",
			combination: []int{10},
			want:        true,
		},
		{
			name:        "2個のみ",
			combination: []int{10, 20},
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.Validate(tt.combination)
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCloseNumberPairValidator_ValidCombination 近距離ペアの存在をテスト
func TestCloseNumberPairValidator_ValidCombination(t *testing.T) {
	validator := &CloseNumberPairValidator{}

	tests := []struct {
		name        string
		combination []int
		want        bool
	}{
		{
			name:        "差1のペアあり",
			combination: []int{1, 2, 10, 15, 20, 25, 30},
			want:        true,
		},
		{
			name:        "差2のペアあり",
			combination: []int{1, 3, 10, 15, 20, 25, 30},
			want:        true,
		},
		{
			name:        "差3のペアあり",
			combination: []int{1, 4, 10, 15, 20, 25, 30},
			want:        true,
		},
		{
			name:        "近距離ペアなし",
			combination: []int{1, 5, 10, 15, 20, 25, 30},
			want:        false,
		},
		{
			name:        "複数の近距離ペア",
			combination: []int{1, 2, 3, 10, 20, 30, 37},
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.Validate(tt.combination)
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCompositeValidator_AllValidators 複合バリデータのテスト
func TestCompositeValidator_AllValidators(t *testing.T) {
	validator := NewCompositeValidator(
		&ZoneDistributionValidator{},
		&OddEvenRatioValidator{},
		&TotalSumValidator{},
		&AverageValueValidator{},
		&CloseNumberPairValidator{},
	)

	tests := []struct {
		name        string
		combination []int
		want        bool
	}{
		{
			name:        "全ての制約を満たす",
			combination: []int{2, 9, 10, 11, 23, 29, 36}, // 実際の推薦例
			want:        true,
		},
		{
			name:        "ゾーン分布が無効",
			combination: []int{1, 2, 3, 4, 5, 6, 7},
			want:        false,
		},
		{
			name:        "奇数偶数比率が無効",
			combination: []int{1, 3, 5, 7, 9, 11, 30},
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.Validate(tt.combination)
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
