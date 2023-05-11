package socket

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	for {
		event, args, err := parseEvent(r)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			log.Println("parse: ", err)
			continue
		}

		go s.handleEvent(w, event, args)
	}
}

func (s *Server) handleEvent(w *bufio.Writer, event Event, args string) {
	switch event {
	case Add:
		players, err := decodePlayers(args)
		if err != nil {
			log.Println(err)
			return
		}

		if err := s.pool.AddPlayer(players...); err != nil {
			log.Println("pool add:", err)
		}
	case Remove:

	case Size:
		size := s.pool.Size()
		if err := writeEvent(w, Size, size); err != nil {
			log.Println("size write:", err)
		}
	default:
		log.Println("Unknown event:", event, args)
	}
}
