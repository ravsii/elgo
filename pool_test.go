package elgo_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ravsii/elgo"
)

func TestPool2(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool()
	_ = pool.AddPlayer(CreatePlayerMock("Test1", 1000))
	_ = pool.AddPlayer(CreatePlayerMock("Test2", 1000))

	go pool.Run()

	acceptMatch(pool, t)

	queue := pool.Close()
	if len(queue) != 0 {
		t.Errorf("test queue should be empty, got: %v", queue)
	}
}

func TestPool3(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool()
	_ = pool.AddPlayer(CreatePlayerMock("Test1", 1000))
	_ = pool.AddPlayer(CreatePlayerMock("Test2", 1000))
	_ = pool.AddPlayer(CreatePlayerMock("Test3", 1000))

	go pool.Run()

	acceptMatch(pool, t)

	queue := pool.Close()
	if len(queue) != 1 {
		t.Errorf("test queue should be len 1")
	}

}

func TestPool1000(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool(elgo.WithGlobalRetry(time.Millisecond), elgo.WithPlayerRetry(time.Millisecond))

	for i := 0; i < 1000; i++ {
		_ = pool.AddPlayer(CreatePlayerMock(fmt.Sprint(i), 1000))
	}

	go pool.Run()

	for i := 0; i < 500; i++ {
		acceptMatch(pool, t)
	}

	queue := pool.Close()
	if len(queue) != 0 {
		t.Errorf("test queue should be empty, got: %v", queue)
	}
}

func TestPool1001(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool(elgo.WithGlobalRetry(time.Millisecond), elgo.WithPlayerRetry(time.Millisecond))

	for i := 0; i < 1001; i++ {
		_ = pool.AddPlayer(CreatePlayerMock(fmt.Sprint(i), 1000))
	}

	go pool.Run()

	for i := 0; i < 500; i++ {
		acceptMatch(pool, t)
	}

	queue := pool.Close()
	if len(queue) != 1 {
		t.Errorf("test queue should have length 1, got: %v", queue)
	}
}

func TestErrAlreadyExists(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool()

	player := CreatePlayerMock("mock", 1000)

	_ = pool.AddPlayer(player)
	err := pool.AddPlayer(player)
	if err == nil || !errors.Is(err, elgo.ErrAlreadyExists) {
		t.Errorf("expected error %s, got %s", elgo.ErrAlreadyExists, err)
	}
}

func TestMatchLongWait(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool(elgo.WithIncreaseInterval(100), elgo.WithPlayerRetry(100*time.Millisecond), elgo.WithGlobalRetry(100*time.Millisecond))
	_ = pool.AddPlayer(CreatePlayerMock("1", 100))
	_ = pool.AddPlayer(CreatePlayerMock("2", 500))

	go pool.Run()

	acceptMatch(pool, t)

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
