package elgo

import "testing"

func TestWithKFactor(t *testing.T) {
	t.Parallel()

	t.Run("empty k", func(t *testing.T) {
		t.Parallel()

		expected := 10.
		calc := NewCalc(expected)
		if calc.k != expected {
			t.Errorf("expected k = %f, got %f", expected, calc.k)
		}
	})

	t.Run("k sorted 2", func(t *testing.T) {
		t.Parallel()

		calc := NewCalc(0, WithKFactor(2000, 10), WithKFactor(1000, 20))
		expected := kFactors{
			{20., 1000.},
			{10., 2000.},
		}
		if !factorsEqual(calc.kFactors, expected) {
			t.Errorf("factors do not match: want %v, got %v", expected, calc.kFactors)
		}
	})

	t.Run("k sorted 5", func(t *testing.T) {
		t.Parallel()

		calc := NewCalc(0,
			WithKFactor(2000, 10),
			WithKFactor(1000, 20),
			WithKFactor(1500, 15),
			WithKFactor(1501, 15),
			WithKFactor(1700, 17))
		expected := kFactors{
			{20., 1000.},
			{15., 1500.},
			{15., 1501.},
			{17., 1700.},
			{10., 2000.},
		}
		if !factorsEqual(calc.kFactors, expected) {
			t.Errorf("factors do not match: want %v, got %v", expected, calc.kFactors)
		}
	})
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
