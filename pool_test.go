package elgo_test

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/ravsii/elgo"
)

func TestPool(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		poolSize           int
		expectedMatches    int
		expectedSizeClosed int
	}{
		{"2", 2, 1, 0},
		{"3", 3, 1, 1},
		{"100", 100, 50, 0},
		{"101", 101, 50, 1},
		{"500", 500, 250, 0},
		{"501", 501, 250, 1},
		{"1000", 1000, 500, 0},
		{"1001", 1001, 500, 1},
		{"10000", 10000, 5000, 0},
		{"10001", 10001, 5000, 1},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pool := elgo.NewPool(
				elgo.WithPlayerRetryInterval(100*time.Millisecond),
				elgo.WithGlobalRetryInterval(100*time.Millisecond),
				elgo.WithIncreasePlayerBorderBy(0.05))

			go pool.Run()
			for i := 0; i < tc.poolSize; i++ {
				go pool.AddPlayer(CreatePlayerMock(fmt.Sprint(i), rand.Float64()))
			}

			for i := 0; i < tc.expectedMatches; i++ {
				acceptMatch(pool, t)
			}

			got := len(pool.Close())
			if got != tc.expectedSizeClosed {
				t.Errorf("pool size on Close() %v, want %v", got, tc.expectedSizeClosed)
			}
		})
	}
}

func TestErrAlreadyExists(t *testing.T) {
	t.Parallel()

	var (
		pool   = elgo.NewPool()
		player = CreatePlayerMock("mock", 1000)
		_      = pool.AddPlayer(player)
		err    = pool.AddPlayer(player)
	)

	if err == nil || !errors.Is(err, elgo.ErrAlreadyExists) {
		t.Errorf("expected error %s, got %s", elgo.ErrAlreadyExists, err)
	}
}

func TestPlayerRetryInterval(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool(
		elgo.WithIncreasePlayerBorderBy(100),
		elgo.WithPlayerRetryInterval(200*time.Millisecond),
		elgo.WithGlobalRetryInterval(time.Millisecond))

	pool.AddPlayer(CreatePlayerMock("1", 100))
	pool.AddPlayer(CreatePlayerMock("2", 1000))

	go pool.Run()

	wg := sync.WaitGroup{}
	wg.Add(1)

	select {
	case <-pool.Matches():
		wg.Done()
	case <-time.After(10 * time.Second):
		t.Errorf("match took too long to create")
	}

	queue := pool.Close()
	if len(queue) != 0 {
		t.Errorf("test queue should be empty, got: %v", queue)
	}
}

func TestGlobalRetryInterval(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool(
		elgo.WithIncreasePlayerBorderBy(100),
		elgo.WithPlayerRetryInterval(time.Millisecond),
		elgo.WithGlobalRetryInterval(1*time.Second))

	pool.AddPlayer(CreatePlayerMock("1", 100))
	pool.AddPlayer(CreatePlayerMock("2", 1000))

	go pool.Run()

	wg := sync.WaitGroup{}
	wg.Add(1)

	select {
	case <-pool.Matches():
		wg.Done()
	case <-time.After(10 * time.Second):
		t.Errorf("match took too long to create")
	}

	queue := pool.Close()
	if len(queue) != 0 {
		t.Errorf("test queue should be empty, got: %v", queue)
	}
}

// acceptMatch tries to read from match channel, throws error otherwise.
func acceptMatch(pool *elgo.Pool, t *testing.T) {
	_, ok := <-pool.Matches()
	if !ok {
		t.Error("channel is closed, but it shouldn't be")
	}
}
