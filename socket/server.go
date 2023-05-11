package socket

import (
	"context"
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

	safeIO := newSafeIO(conn)
	ctx, cancel := context.WithCancel(context.Background())

	go s.sendMatches(ctx, safeIO)

	for {
		event, args, err := safeIO.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			log.Println("parse: ", err)
			continue
		}

		go s.handleEvent(safeIO, event, args)
	}

	cancel()
}

func (s *Server) sendMatches(ctx context.Context, safeWriter *safeIO) {
	for {
		select {
		case match, ok := <-s.pool.Matches():
			if !ok {
				return
			}

			s := fmt.Sprintf("%s;%s", match.Player1.Identify(), match.Player2.Identify())
			if err := safeWriter.Write(Match, s); err != nil {
				log.Println("write:", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *Server) handleEvent(safeWriter *safeIO, event Event, args string) {
	switch event {
	case Match: // Match is handled in sendMatches()

	case Add:
		players, err := decodeRatingPlayers(args)
		if err != nil {
			log.Println(err)
			return
		}

		if err := s.pool.AddPlayer(players...); err != nil {
			log.Println("pool add:", err)
		}
	case Remove:
		s.pool.Remove()
	case Size:
		size := s.pool.Size()
		if err := safeWriter.Write(Size, size); err != nil {
			log.Println("size write:", err)
		}
	case Unknown:
		fallthrough
	default:
		log.Println("Unknown event:", event, args)
	}
}
