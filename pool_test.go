package elgo_test

import (
	"testing"

	"github.com/ravsii/elgo"
)

func TestPool2(t *testing.T) {
	pool := elgo.NewPool()
	pool.Queue() <- CreatePlayerMock("Test1", 1000)
	pool.Queue() <- CreatePlayerMock("Test2", 1000)

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
	pool := elgo.NewPool()
	pool.Queue() <- CreatePlayerMock("Test1", 1000)
	pool.Queue() <- CreatePlayerMock("Test2", 1000)
	pool.Queue() <- CreatePlayerMock("Test3", 1000)

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

	pool := elgo.NewPool(elgo.WithIncreaseInterval(100))

	for i := 0; i < 1000; i++ {
		pool.Queue() <- CreatePlayerMock("Test1", 1000)
	}

	for i := 0; i < 500; i++ {
		<-pool.Matches()
		t.Log(i)
	}

	queue := pool.Close()
	if len(queue) != 0 {
		t.Errorf("test queue should be empty, got: %v", queue)
	}
}
