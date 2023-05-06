package socket

import (
	"fmt"
	"net"

	"github.com/ravsii/elgo"
)

type Client struct {
	conn    net.Conn
	matches chan elgo.Match
}

func NewClient(listenAddr string) (*Client, error) {
	conn, err := net.Dial("tcp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	c := &Client{
		conn:    conn,
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

// Close closes c.conn and matches channel
func (c *Client) Close() error {
	close(c.matches)
	return c.conn.Close()
}

func (c *Client) listenForMatches() {
	// for {

	// }
}
