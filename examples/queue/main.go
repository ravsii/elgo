package main

import (
	"fmt"

	"github.com/ravsii/elgo"
	"github.com/ravsii/elgo/examples/queue/player"
)

func main() {
	P1 := &player.Player{Name: "Alex", EloRating: 1000}
	P2 := &player.Player{Name: "Greg", EloRating: 1500}
	P3 := &player.Player{Name: "Kate", EloRating: 3000}
	P4 := &player.Player{Name: "Oleg", EloRating: 5000}

	pool := elgo.NewPool()
	playerChan := pool.Queue()

	playerChan <- P1
	playerChan <- P2
	playerChan <- P3
	playerChan <- P4

	for i := 0; i < 2; i++ {
		match, ok := <-pool.Matches()
		fmt.Println("match", match, ok)
	}

	fmt.Println(pool.Close())
}
