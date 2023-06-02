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

			t.Cleanup(func() { pool.Close() })
			go pool.Run()

			wg := sync.WaitGroup{}
			wg.Add(tc.expectedMatches)

			for i := 0; i < tc.poolSize; i++ {
				err := pool.AddPlayer(CreatePlayerMock(fmt.Sprint(i), rand.Float64())) //nolint
				if err != nil {
					t.Errorf("pool add player: %v", err)
				}
			}

			for i := 0; i < tc.expectedMatches; i++ {
				go func(t *testing.T, wg *sync.WaitGroup, pool *elgo.Pool) {
					t.Helper()
					acceptMatch(t, pool)
					wg.Done()
				}(t, &wg, pool)
			}

			wg.Wait()

			got := len(pool.Close())
			if got != tc.expectedSizeClosed {
				t.Errorf("pool size on Close() %v, want %v", got, tc.expectedSizeClosed)
			}
		})
	}
}

func TestPoolPrematureClose(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		poolSize          int
		closeAfterMatches int
		wantPlayersLeft   int
	}{
		{"2", 2, 1, 0},
		{"3", 3, 1, 1},
		{"100", 100, 49, 2},
		{"101", 101, 50, 1},
		{"500", 500, 200, 100},
		{"501", 500, 200, 100},
		{"1000", 1000, 200, 600},
		{"1001", 1001, 201, 599},
		{"10000", 10000, 100, 9800},
		{"10001", 10001, 4000, 2001},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			pool := elgo.NewPool(
				elgo.WithPlayerRetryInterval(100*time.Millisecond),
				elgo.WithGlobalRetryInterval(100*time.Millisecond),
				elgo.WithIncreasePlayerBorderBy(0.05))

			t.Cleanup(func() { pool.Close() })

			go pool.Run()
			for i := 0; i < testCase.poolSize; i++ {
				go pool.AddPlayer(CreatePlayerMock(fmt.Sprint(i), rand.Float64())) //nolint
			}

			for i := 0; i < testCase.closeAfterMatches; i++ {
				acceptMatch(t, pool)
			}

			got := len(pool.Close())
			if got != testCase.wantPlayersLeft {
				t.Errorf("pool size on Close() %v, want %v", got, testCase.wantPlayersLeft)
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

	t.Cleanup(func() { pool.Close() })

	if err == nil || !errors.Is(err, elgo.ErrAlreadyExists) {
		t.Errorf("expected error %s, got %s", elgo.ErrAlreadyExists, err)
	}
}

func TestErrPoolClosed(t *testing.T) {
	t.Parallel()

	var (
		pool   = elgo.NewPool()
		player = CreatePlayerMock("mock", 1000)
	)

	_ = pool.Close()

	if err := pool.AddPlayer(player); err == nil || !errors.Is(err, elgo.ErrPoolClosed) {
		t.Errorf("expected error %s, got %s", elgo.ErrPoolClosed, err)
	}
}

func TestPlayerRetryInterval(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool(
		elgo.WithIncreasePlayerBorderBy(100),
		elgo.WithPlayerRetryInterval(200*time.Millisecond),
		elgo.WithGlobalRetryInterval(time.Millisecond))

	t.Cleanup(func() { pool.Close() })

	if err := pool.AddPlayer(CreatePlayerMock("1", 100)); err != nil {
		t.Errorf("pool add player: %s", err)
	}

	if err := pool.AddPlayer(CreatePlayerMock("2", 1000)); err != nil {
		t.Errorf("pool add player: %s", err)
	}

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

	t.Cleanup(func() { pool.Close() })

	if err := pool.AddPlayer(CreatePlayerMock("1", 100)); err != nil {
		t.Errorf("pool add player: %s", err)
	}

	if err := pool.AddPlayer(CreatePlayerMock("2", 1000)); err != nil {
		t.Errorf("pool add player: %s", err)
	}

	go pool.Run()

	wg := sync.WaitGroup{}
	wg.Add(1)

	select {
	case <-pool.Matches():
		wg.Done()
	case <-time.After(10 * time.Second):
		t.Errorf("match took too long to create")
	}

	if queue := pool.Close(); len(queue) != 0 {
		t.Errorf("test queue should be empty, got: %v", queue)
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool()
	t.Cleanup(func() { pool.Close() })

	p1 := CreatePlayerMock("p1", 1000)
	p2 := CreatePlayerMock("p2", 1000)

	_ = pool.AddPlayer(p1, p2)
	pool.Remove(p1, p2)

	if players := pool.Close(); len(players) != 0 {
		t.Errorf("RemoveStrs did not remove all players.")
	}
}

func TestRemoveStrs(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool()
	t.Cleanup(func() { pool.Close() })

	p1 := CreatePlayerMock("p1", 1000)
	p2 := CreatePlayerMock("p2", 1000)

	_ = pool.AddPlayer(p1, p2)
	pool.RemoveStrs("p1", "p2")

	if players := pool.Close(); len(players) != 0 {
		t.Errorf("RemoveStrs did not remove all players.")
	}
}

// acceptMatch tries to read from match channel, throws error otherwise.
func acceptMatch(t *testing.T, pool *elgo.Pool) {
	t.Helper()

	if _, ok := <-pool.Matches(); !ok {
		t.Error("channel is closed, but it shouldn't be")
	}
}
