package socket

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/ravsii/elgo"
)

// Listen creates a new server on a port using a pool and listens for connections.
// Listen is a blocking function.
//
// Returned error is always a non-nil error.
func Listen(port int, pool *elgo.Pool) (err error) {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	defer func() { err = listen.Close() }()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return fmt.Errorf("accept connection: %w", err)
		}

		go handle(conn, pool)
	}
}

func handle(conn net.Conn, p *elgo.Pool) {
	defer conn.Close()

	reader := bufio.NewReaderSize(conn, 1024)
	writer := bufio.NewWriterSize(conn, 1024)

	for done := false; !done; {
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println("read all: ", err)
			return
		}

		log.Println("received", input)

		event, _ := parseEvent(input)
		switch event {
		case "ADD":

		case "MATCH":

		case "REMOVE":

		case "SIZE":
			size := p.Size()
			writer.Write(createEvent(Size, size))
			writer.Flush()
			fmt.Print(1)

		default:
			conn.Write([]byte("test"))
		}
	}
}
