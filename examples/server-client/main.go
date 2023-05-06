package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/ravsii/elgo"
	"github.com/ravsii/elgo/examples/player"
	"github.com/ravsii/elgo/socket"
)

func main() {
	pool := elgo.NewPool()
	pool.AddPlayer(
		player.New("Example 1", 0),
		player.New("Example 2", 0))

	server := socket.NewServer(":3000", pool)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		log.Println("Server started")
		wg.Done()
		log.Fatal(server.Listen())
	}()

	wg.Wait()

	client, err := socket.NewClient(":3000")
	if err != nil {
		log.Fatal("unable to connect to server:", err)
	}

	for i := 0; i < 100; i++ {
		p := player.New(fmt.Sprint(i), rand.Float64())
		fmt.Println("client: adding new player", p)
		err := client.Add(p)
		if err != nil {
			log.Fatal("client: add", err)
		}

		time.Sleep(time.Second)
	}
}
