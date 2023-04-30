package elgo

import (
	"testing"
	"time"
)

func TestWithIncreaseIntervals(t *testing.T) {
	t.Parallel()

	var (
		want = 1000.
		pool = NewPool(WithIncreasePlayerBorderBy(want))
	)

	t.Cleanup(func() { pool.Close() })

	if pool.playersBordersIncreaseBy != want {
		t.Errorf("want playersBordersIncreaseBy %f, got %f", want, pool.playersBordersIncreaseBy)
	}
}

func TestWithPlayerRetry(t *testing.T) {
	t.Parallel()

	var (
		want = 3 * time.Hour
		pool = NewPool(WithPlayerRetryInterval(want))
	)

	t.Cleanup(func() { pool.Close() })

	if pool.playerRetryInterval != want {
		t.Errorf("want playerRetryInterval %d, got %d", want, pool.playerRetryInterval)

	}
}

func TestWithGlobalRetry(t *testing.T) {
	t.Parallel()

	var (
		want = 3 * time.Hour
		pool = NewPool(WithGlobalRetryInterval(want))
	)

	t.Cleanup(func() { pool.Close() })

	if pool.globalRetryInterval != want {
		t.Errorf("want globalRetryInterval %d, got %d", want, pool.globalRetryInterval)
	}
}
