package main

import (
	"fmt"
	"log"

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

	go func() {
		log.Println("Server started")
		log.Fatal(server.Listen())
	}()

	client, err := socket.NewClient(":3000")
	if err != nil {
		log.Fatal("unable to connect to server:")
	}

	fmt.Println(client.Size())
}
