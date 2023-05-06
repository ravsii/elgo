package socket

import (
	"bufio"
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

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
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
			size := s.pool.Size()
			writer.Write(createEvent(Size, size))
			writer.Flush()
			fmt.Print(1)

		default:
			conn.Write([]byte("test"))
		}
	}
}
