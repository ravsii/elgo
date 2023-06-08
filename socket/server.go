package socket

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/ravsii/elgo"
)

type Server struct {
	pool     *elgo.Pool
	listener net.Listener
}

// NewServer creates a server. Use
//
//	server.Listen(network, addr)
//
// to run it.
func NewServer(pool *elgo.Pool) *Server {
	return &Server{
		pool: pool,
	}
}

// Listen starts listening for connections. It is a blocking operation.
//
// If you want to perform a graceful shutdown, use s.Close()
//
// Returned error is always a non-nil error.
func (s *Server) Listen(network, addr string) error {
	var err error
	s.listener, err = net.Listen(network, addr)
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return fmt.Errorf("accept connection: %w", err)
		}

		go s.handleConn(conn)
	}
}

// Close only closes the underlying connection, the pool remains open.
func (s *Server) Close() error {
	if s.listener == nil {
		return nil
	}

	return s.listener.Close()
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	readWriter := newReadWriter(conn)
	ctx, cancel := context.WithCancel(context.Background())

	go s.sendMatches(ctx, readWriter)

	for {
		event, args, err := readWriter.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			log.Println("parse: ", err)
			continue
		}

		go s.handleEvent(readWriter, event, args)
	}

	cancel()
}

func (s *Server) sendMatches(ctx context.Context, readWriter *ReadWriter) {
	for {
		select {
		case match, ok := <-s.pool.Matches():
			if !ok {
				return
			}

			s := fmt.Sprintf("%s;%s", match.Player1.Identify(), match.Player2.Identify())
			if err := readWriter.Write(Match, s); err != nil {
				log.Println("write:", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *Server) handleEvent(readWriter *ReadWriter, event Event, args string) {
	switch event {
	case Add:
		split := strings.Split(args, " ")
		players := make([]elgo.Player, 0, len(split))

		for _, playerStr := range split {
			id, ratingStr, found := strings.Cut(playerStr, ";")
			if !found {
				log.Println(ErrBadInput, ":", playerStr)
				return
			}

			r, err := strconv.ParseFloat(ratingStr, 64)
			if err != nil {
				log.Println("parse rating:", err)
				return
			}

			players = append(players, &elgo.BaseRatingPlayer{ID: id, ELO: r})
		}

		if err := s.pool.AddPlayer(players...); err != nil {
			log.Println("pool add:", err)
		}
	case Remove:
		toRemove := strings.Split(args, " ")
		s.pool.RemoveStrs(toRemove...)
	case Size:
		size := s.pool.Size()
		if err := readWriter.Write(Size, size); err != nil {
			log.Println("size write:", err)
		}
	default:
		log.Println("Unknown event:", event, args)
	}
}
