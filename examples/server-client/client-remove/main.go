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

	for i := 0; i < 10; i++ {
		p := player.New(fmt.Sprint(i), rand.Float64())
		if err := client.Add(p); err != nil {
			log.Fatalf("client: add player: %s", err)
		}
	}

	go client.RemoveStrs("1", "3")

	go time.AfterFunc(5*time.Second, func() {
		client.RemoveStrs("5", "6", "7")
	})

	time.Sleep(10 * time.Second)
}
