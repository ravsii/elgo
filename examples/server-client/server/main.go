package main

import (
	"log"

	"github.com/ravsii/elgo"
	"github.com/ravsii/elgo/examples/player"
	"github.com/ravsii/elgo/socket"
)

func main() {
	pool := elgo.NewPool()
	err := pool.AddPlayer(
		player.New("Example 1", 0),
		player.New("Example 2", 0))
	if err != nil {
		return
	}

	defer pool.Close()

	go pool.Run()

	server := socket.NewServer(pool)
	log.Fatal(server.Listen("tcp", ":8080"))
}
