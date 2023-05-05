package socket

import (
	"fmt"
	"net"

	"github.com/ravsii/elgo"
)

type Client struct {
	server  net.Conn
	matches chan elgo.Match
}

func NewClient(port int) (*Client, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	c := &Client{
		server:  conn,
		matches: make(chan elgo.Match),
	}

	go c.listenForMatches()

	return c, nil
}

func (c *Client) Add() {}

// ReceiveMatch wait for a match to appear and returns now.
//
// This is a blocking operation.
func (c *Client) ReceiveMatch() elgo.Match {
	return <-c.matches
}

// Size returns current amount of players in the pool
func (c *Client) Size() int {
	return 0
}

func (c *Client) listenForMatches() {
	// for {

	// }
}
