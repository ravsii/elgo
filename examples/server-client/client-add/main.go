package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ravsii/elgo/examples/player"
	"github.com/ravsii/elgo/socket"
)

func main() {
	client, err := socket.NewClient(":8080")
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

	for i := 0; i < matches; i++ {
		match := <-client.ReceiveMatch()
		log.Println(i, "match", match.Player1.Identify(), match.Player2.Identify())
	}

	time.Sleep(10 * time.Second)
}
