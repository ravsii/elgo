package socket

import (
	"fmt"
	"io"
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

	for done := false; !done; {
		b, err := io.ReadAll(conn)
		if err != nil {
			log.Println("read all: ", err)
			return
		}

		event, _ := parseEvent(string(b))
		switch event {
		case "ADD":

		case "MATCH":

		case "REMOVE":

		case "SIZE":
			size := p.Size()
			conn.Write(createEvent(Size, size))

			// default:
			// 	continue
		}

		conn.Close()
	}

}
