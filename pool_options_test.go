package elgo

import (
	"testing"
	"time"
)

func TestWithIncreaseIntervals(t *testing.T) {
	t.Parallel()

	var (
		expected = 1000.
		pool     = NewPool(WithIncreaseInterval(expected))
	)

	defer pool.Close()
	if pool.playersBordersIncreaseBy != expected {
		t.Errorf("expected interval %f, got %f", expected, pool.playersBordersIncreaseBy)

	}
}

func TestWithPlayerRetry(t *testing.T) {
	t.Parallel()

	var (
		expectedDuration = 3 * time.Hour
		pool             = NewPool(WithPlayerRetry(expectedDuration))
	)

	defer pool.Close()
	if pool.playerRetryInterval != expectedDuration {
		t.Errorf("expected player retry %d, got %d", expectedDuration, pool.playerRetryInterval)

	}
}

func TestWithGlobalRetry(t *testing.T) {
	t.Parallel()

	var (
		expectedDuration = 3 * time.Hour
		pool             = NewPool(WithGlobalRetry(expectedDuration))
	)

	defer pool.Close()
	if pool.playerRetryInterval != expectedDuration {
		t.Errorf("expected global retry %d, got %d", expectedDuration, pool.globalRetryInterval)

	}

}
