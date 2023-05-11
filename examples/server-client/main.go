package main

import (
	"context"
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

	defer pool.Close()

	go pool.Run()

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
			time.Sleep(500 * time.Millisecond)
			size, err := client.Size()
			if err != nil {
				log.Fatalf("client: size: %s", err)
			}

			fmt.Println("client: size", size)
		}
	}()

	matches := 10000

	for i := 0; i < matches*2; i++ {
		p := player.New(fmt.Sprint(i), rand.Float64())
		err := client.Add(p)
		if err != nil {
			log.Fatalf("client: add player: %s", err)
		}
	}

	ctx := context.Background()

	for i := 0; i < matches; i++ {
		p := client.ReceiveMatch(ctx)
		log.Println("match", p.Player1.Identify(), p.Player2.Identify())
	}

	time.Sleep(10 * time.Second)
}
