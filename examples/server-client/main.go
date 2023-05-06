package main

import (
	"fmt"
	"log"
	"math/rand"
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

	g := make(chan bool)

	go func() {
		log.Println("Server started")
		g <- true
		log.Fatal(server.Listen())
	}()

	<-g

	client, err := socket.NewClient(":3000")
	if err != nil {
		log.Fatal("unable to connect to server:", err)
	}

	go func() {
		for {
			time.Sleep(time.Second)
			size, err := client.Size()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("client: size", size)
		}
	}()

	for i := 0; i < 10; i++ {
		p := player.New(fmt.Sprint(i), rand.Float64())
		err := client.Add(p)
		if err != nil {
			log.Fatal("client: add", err)
		}

		// time.Sleep(time.Second)
	}

	time.Sleep(2 * time.Second)
}
