package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conn, err := net.Dial("tcp", ":25561")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		reader := bufio.NewReader(conn)

		for {
			b, err := reader.ReadString('\n')
			if err != nil {
				log.Println("ERR:", err)
				continue
			}

			fmt.Println("RECV >", string(b))
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-c
		_ = conn.Close()
		os.Exit(1)
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println("ERR:", err)
			continue
		}

		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Println("ERR:", err)
			continue
		}
	}
}
