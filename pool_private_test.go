package elgo

import "testing"

func TestWithIntervals(t *testing.T) {
	expected := 1000.

	pool := NewPool(WithIncreaseInterval(expected))
	if pool.increaseRatingBorders != expected {
		t.Errorf("expected interval %f, got %f", expected, pool.increaseRatingBorders)

	}
}
