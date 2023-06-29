package socket

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ravsii/elgo"
)

var ErrNoResponse = errors.New("server didn't respond")

type Client struct {
	conn       net.Conn
	readWriter *ReadWriter
	sizeCh     chan int
	matchCh    chan *elgo.Match
	closeCh    chan struct{}
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
		closeCh:    make(chan struct{}),
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
func (c *Client) Close() error {
	c.closeCh <- struct{}{}
	close(c.closeCh)
	close(c.matchCh)
	close(c.sizeCh)
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("conn close: %w", err)
	}

	return nil
}

// In this function we check that the error occurred is about server disconnection or not
// first check that error is about network operation or not
// second we check that it is about syscallErrors or not
// and at the end we check is it specifically ECONNRESET or not if it is, then it is a server disconnection
func isConnectionResetError(err error) bool {
	netErr, ok := err.(*net.OpError)
	if ok {
		syscallErr, ok := netErr.Err.(*os.SyscallError)
		if ok {
			errno, ok := syscallErr.Err.(syscall.Errno)
			if ok && errno == syscall.ECONNRESET {
				return true
			}
		}
	}
	return false
}

func (c *Client) listen() {
	for {
		select {
		case <-c.closeCh:
			return
		default:
			event, args, err := c.readWriter.Read()
			// If an error occurred when reading data, entered the if block
			if err != nil {
				// if each one of these conditions set true we will break the for and stop the client
				if isConnectionResetError(err) || strings.ContainsAny(err.Error(), "closed") ||
					err == io.EOF || strings.ContainsAny(err.Error(), "reset") || strings.ContainsAny(err.Error(), "EOF") {
					log.Println("server disconnected.")
					return
				}
				// but if none of them set true, it means server is still connected and there is a problem in readwrite buffer
				log.Println("err while read: ", err)
				continue
			}
			go c.handleEvent(event, args)
		}
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
			Player1: &elgo.BasePlayer{ID: p1Ident},
			Player2: &elgo.BasePlayer{ID: p2Ident},
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
