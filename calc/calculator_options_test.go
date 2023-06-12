package calc

import "testing"

func TestWithoutKFactors(t *testing.T) {
	t.Parallel()

	if calc := New(10); calc.k != 10 {
		t.Errorf("NewCalc() k = %f, got %f", 10., calc.k)
	}
}

func TestWithKFactor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		k     float64
		funcs []CalcOpt
		want  kFactors
	}{
		{"empty", 10, nil, nil},
		{
			"2 sorted",
			10,
			[]CalcOpt{
				WithKFactor(1000, 20),
				WithKFactor(1500, 30),
			},
			kFactors{
				{20, 1000},
				{30, 1500},
			},
		},
		{
			"2 not sorted",
			10,
			[]CalcOpt{
				WithKFactor(1500, 30),
				WithKFactor(1000, 20),
			},
			kFactors{
				{20, 1000},
				{30, 1500},
			},
		},
		{
			"5 not sorted", 10,
			[]CalcOpt{
				WithKFactor(2000, 10),
				WithKFactor(1000, 20),
				WithKFactor(1500, 15),
				WithKFactor(1501, 15),
				WithKFactor(1700, 17),
			},
			kFactors{
				{20, 1000},
				{15, 1500},
				{15, 1501},
				{17, 1700},
				{10, 2000},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			calc := New(tt.k, tt.funcs...)
			if !factorsEqual(calc.kFactors, tt.want) {
				t.Errorf("factors do not match: want %v, got %v", tt.want, calc.kFactors)
			}
		})
	}
}

func factorsEqual(f1, f2 kFactors) bool {
	if len(f1) != len(f2) {
		return false
	}

	for i := 0; i < len(f1); i++ {
		if f1[i].k != f2[i].k || f1[i].startsAt != f2[i].startsAt {
			return false
		}
	}

	return true
}
