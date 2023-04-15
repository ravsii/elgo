package elgo

import (
	"testing"
	"time"
)

func TestWithIncreaseIntervals(t *testing.T) {
	t.Parallel()

	expected := 1000.

	pool := NewPool(WithIncreaseInterval(expected))
	if pool.increaseRatingBorders != expected {
		t.Errorf("expected interval %f, got %f", expected, pool.increaseRatingBorders)

	}

	_ = pool.Close()
}

func TestWithPlayerRetry(t *testing.T) {
	t.Parallel()

	expectedDuration := 3 * time.Hour

	pool := NewPool(WithPlayerRetry(expectedDuration))
	if pool.retryPlayerSearch != expectedDuration {
		t.Errorf("expected player retry %d, got %d", expectedDuration, pool.retryPlayerSearch)

	}

	_ = pool.Close()
}

func TestWithGlobalRetry(t *testing.T) {
	t.Parallel()

	expectedDuration := 3 * time.Hour

	pool := NewPool(WithGlobalRetry(expectedDuration))
	if pool.retryPlayerSearch != expectedDuration {
		t.Errorf("expected global retry %d, got %d", expectedDuration, pool.retryGlobalSearch)

	}

	_ = pool.Close()
}
