package socket

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/ravsii/elgo"
)

var ErrNoResponse = errors.New("server didn't respond")

type Client struct {
	conn       net.Conn
	readWriter *ReadWriter

	sizeCh  chan int
	matchCh chan *elgo.Match
}

func NewClient(listenAddr string) (*Client, error) {
	conn, err := net.Dial("tcp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	c := &Client{
		conn:       conn,
		readWriter: newReadWriter(conn),
		sizeCh:     make(chan int),
		matchCh:    make(chan *elgo.Match),
	}

	go c.listen()

	return c, nil
}

func (c *Client) Add(players ...elgo.Player) error {
	encoded := make([]any, 0, len(players))

	for _, p := range players {
		encoded = append(encoded, fmt.Sprintf("%s;%f", p.Identify(), p.Rating()))
	}

	if err := c.readWriter.Write(Add, encoded...); err != nil {
		return fmt.Errorf("can't add players to the queue: %w", err)
	}

	return nil
}

// ReceiveMatch returns a match channel to listen to.
func (c *Client) ReceiveMatch() <-chan *elgo.Match {
	return c.matchCh
}

// Remove removes players from the pool.
func (c *Client) Remove(players ...elgo.Identifier) error {
	s := make([]string, 0)
	for _, p := range players {
		s = append(s, p.Identify())
	}

	return c.RemoveStrs(s...)
}

// RemoveStrs removes players from the pool by their Identifiers.
func (c *Client) RemoveStrs(identifiers ...string) error {
	if err := c.readWriter.Write(Remove, strings.Join(identifiers, " ")); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

// Size returns current amount of players in the pool.
func (c *Client) Size() (int, error) {
	if err := c.readWriter.Write(Size); err != nil {
		return 0, fmt.Errorf("unable to write: %w", err)
	}

	select {
	case size := <-c.sizeCh:
		return size, nil
	case <-time.After(10 * time.Second):
		return 0, ErrNoResponse
	}
}

// Close closes c.conn and matches channel.
func (c *Client) Close() (err error) {
	close(c.matchCh)
	defer func() {
		if err = c.conn.Close(); err != nil {
			err = fmt.Errorf("conn close: %w", err)
		}
	}()

	return nil
}

func (c *Client) listen() {
	for {
		event, args, err := c.readWriter.Read()
		if err != nil {
			log.Println("err while read: ", err)
			continue
		}

		go c.handleEvent(event, args)
	}
}

func (c *Client) handleEvent(event Event, args string) {
	switch event {
	case Match:
		s := strings.TrimSpace(args)
		p1Ident, p2Ident, found := strings.Cut(s, ";")
		if !found {
			log.Printf("cut not found")
		}

		c.matchCh <- &elgo.Match{
			Player1: &socketPlayer{ID: p1Ident},
			Player2: &socketPlayer{ID: p2Ident},
		}
	case Size:
		s := strings.TrimSpace(args)
		size, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Println("parse size:", err)
		}

		c.sizeCh <- int(size)
	default:
		log.Printf("got unknown event %s %s, ignoring", event, args)
	}
}
