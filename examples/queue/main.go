package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ravsii/elgo"
	"github.com/ravsii/elgo/examples/queue/player"
)

func main() {
	pool := elgo.NewPool(elgo.WithRetry(1*time.Second), elgo.WithIncreaseInterval(0.01))
	playerChan := pool.Queue()

	go func() {
		t := time.NewTicker(time.Second)
		for {
			fmt.Println("size", pool.Size())
			<-t.C
		}
	}()

	for i := 0; i < 10; i++ {
		playerChan <- &player.Player{Name: fmt.Sprint(i), EloRating: rand.Float64()}
	}

	for i := 0; i < 5; i++ {
		match, ok := <-pool.Matches()
		fmt.Println("match", match, ok)
	}

	fmt.Println(pool.Close())
}
