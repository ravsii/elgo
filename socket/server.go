package socket

import (
	"fmt"
	"log"
	"net"

	"github.com/ravsii/elgo"
)

type Server struct {
	addr string
	pool *elgo.Pool
}

// NewServer creates a server. Use
//
//	server.Listen()
//
// to run it.
func NewServer(listenAddr string, pool *elgo.Pool) *Server {
	return &Server{
		addr: listenAddr,
		pool: pool,
	}
}

// Listen starts listening for connections. It is a blocking function.
//
// Returned error is always a non-nil error.
func (s *Server) Listen() (err error) {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	defer func() { err = listen.Close() }()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return fmt.Errorf("accept connection: %w", err)
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	for done := false; !done; {
		event, args, err := parseEvent(conn)
		if err != nil {
			log.Println("parse: ", err)
			continue
		}

		go s.handleEvent(conn, event, args)
	}
}

func (s *Server) handleEvent(conn net.Conn, event Event, args string) {
	switch event {
	case "ADD":

	case "MATCH":

	case "REMOVE":

	case "SIZE":
		size := s.pool.Size()
		writeEvent(conn, Size, size)
	default:
		log.Println("Unknown event:", event, args)
	}
}
