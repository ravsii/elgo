package socket

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ravsii/elgo"
)

type Client struct {
	conn    net.Conn
	addCh   chan elgo.Player
	sizeCh  chan sizeOp
	matchCh chan elgo.Match
}

type sizeOp struct {
	size int
	err  error
}

func NewClient(listenAddr string) (*Client, error) {
	conn, err := net.Dial("tcp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	c := &Client{
		conn:    conn,
		addCh:   make(chan elgo.Player),
		sizeCh:  make(chan sizeOp),
		matchCh: make(chan elgo.Match),
	}

	go c.listen()

	return c, nil
}

func (c *Client) Add(players ...elgo.Player) error {
	encoded := make([]any, 0, len(players))

	for _, p := range players {
		encoded = append(encoded, encodePlayer(p))
	}

	if err := writeEvent(c.conn, Add, encoded...); err != nil {
		return fmt.Errorf("can't add players to the queue: %w", err)
	}

	return nil
}

// ReceiveMatch wait for a match to appear and returns now.
//
// This is a blocking operation.
func (c *Client) ReceiveMatch() elgo.Match {
	return <-c.matchCh
}

// Size returns current amount of players in the pool.
func (c *Client) Size() (int, error) {
	if err := writeEvent(c.conn, Size); err != nil {
		return 0, fmt.Errorf("unable to write: %w", err)
	}

	select {
	case result := <-c.sizeCh:
		if result.err != nil {
			return 0, result.err
		}

		return result.size, nil
	case <-time.After(10 * time.Second):
		return 0, nil
	}
}

// Close closes c.conn and matches channel
func (c *Client) Close() error {
	close(c.matchCh)
	return c.conn.Close()
}

func (c *Client) listen() {
	for {
		event, args, err := parseEvent(c.conn)
		if err != nil {
			log.Println("err while read: ", err)
			continue
		}

		go c.handleEvent(event, args)
	}
}

func (c *Client) handleEvent(event Event, args string) {
	switch event {
	case Size:
		size, err := parseSize(args)
		if err != nil {
			c.sizeCh <- sizeOp{0, err}
			return
		}

		c.sizeCh <- sizeOp{size, nil}
	default:
		fmt.Println("default")
	}
}
