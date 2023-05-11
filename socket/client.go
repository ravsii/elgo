package socket

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ravsii/elgo"
)

var ErrNoResponse = errors.New("server didn't respond")

type Client struct {
	conn net.Conn
	io   *safeIO

	addCh   chan elgo.Player
	sizeCh  chan int
	matchCh chan elgo.Match
}

func NewClient(listenAddr string) (*Client, error) {
	conn, err := net.Dial("tcp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	c := &Client{
		conn: conn,
		io:   newSafeIO(conn),

		addCh:   make(chan elgo.Player),
		sizeCh:  make(chan int),
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

	if err := c.io.Write(Add, encoded...); err != nil {
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
	if err := c.io.Write(Size); err != nil {
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
		event, args, err := c.io.Read()
		if err != nil {
			log.Println("err while read: ", err)
			continue
		}

		c.handleEvent(event, args)
	}
}

func (c *Client) handleEvent(event Event, args string) {
	switch event {
	case Size:
		size, err := parseSize(args)
		c.sizeCh <- size
		if err != nil {
			log.Println("size:", err)
		}
	case Unknown:
		fallthrough
	default:
		log.Printf("got unknown event %s %s, ignoring", event, args)
	}
}
