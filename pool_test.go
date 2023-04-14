package elgo_test

import (
	"fmt"
	"testing"

	"github.com/ravsii/elgo"
)

func TestPool2(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool()
	pool.AddPlayer(CreatePlayerMock("Test1", 1000))
	pool.AddPlayer(CreatePlayerMock("Test2", 1000))

	go pool.Run()

	_, ok := <-pool.Matches()
	if !ok {
		t.Error("channel is closed, but it shouldn't be")
	}

	queue := pool.Close()
	if len(queue) != 0 {
		t.Errorf("test queue should be empty, got: %v", queue)
	}
}

func TestPool3(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool()
	pool.AddPlayer(CreatePlayerMock("Test1", 1000))
	pool.AddPlayer(CreatePlayerMock("Test2", 1000))
	pool.AddPlayer(CreatePlayerMock("Test3", 1000))

	go pool.Run()

	_, ok := <-pool.Matches()
	if !ok {
		t.Error("channel is closed, but it shouldn't be")
	}

	queue := pool.Close()
	if len(queue) != 1 {
		t.Errorf("test queue should be len 1")
	}

}

func TestPool1000(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool()

	for i := 0; i < 1000; i++ {
		pool.AddPlayer(CreatePlayerMock(fmt.Sprint(i), 1000))
	}

	go pool.Run()

	for i := 0; i < 500; i++ {
		_, ok := <-pool.Matches()
		if !ok {
			t.Error("channel is closed, but it shouldn't be")
		}
	}

	queue := pool.Close()
	if len(queue) != 0 {
		t.Errorf("test queue should be empty, got: %v", queue)
	}
}

func TestPool1001(t *testing.T) {
	t.Parallel()

	pool := elgo.NewPool()

	for i := 0; i < 1001; i++ {
		pool.AddPlayer(CreatePlayerMock(fmt.Sprint(i), 1000))
	}

	go pool.Run()

	for i := 0; i < 500; i++ {
		_, ok := <-pool.Matches()
		if !ok {
			t.Error("channel is closed, but it shouldn't be")
		}
	}

	queue := pool.Close()
	if len(queue) != 1 {
		t.Errorf("test queue should have length 1, got: %v", queue)
	}
}
