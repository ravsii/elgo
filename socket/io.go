package socket

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

const Delimiter byte = '\n'

type safeIO struct {
	r *bufio.Reader
	w *bufio.Writer
	m sync.Mutex
}

func newSafeIO(c net.Conn) *safeIO {
	return &safeIO{
		r: bufio.NewReader(c),
		w: bufio.NewWriter(c),
	}
}

// Read accepts a string event and parses it, returning Event type and
// the rest of the string.
//
// Events are expected to be in such format:
//
//	ADD ...
//	MATCH ...
//	SIZE
//
// If no known event type was found, Unknown is returned.
func (c *safeIO) Read() (Event, string, error) {
	s, err := c.r.ReadString(Delimiter)
	if err != nil {
		return Unknown, "", fmt.Errorf("reader: %w", err)
	}

	s = strings.Trim(s, " \n\r")
	eventStr, args, _ := strings.Cut(s, " ")
	log.Println("read", eventStr, args)

	switch eventStr {
	case "ADD":
		return Add, args, nil
	case "MATCH":
		return Match, args, nil
	case "REMOVE":
		return Remove, args, nil
	case "SIZE":
		return Size, args, nil
	default:
		return Unknown, s, nil
	}
}

func (c *safeIO) Write(event Event, args ...any) error {
	c.m.Lock()
	defer c.m.Unlock()

	str := string(event)

	if len(args) > 0 {
		strs := make([]string, 0, len(args))
		for _, arg := range args {
			strs = append(strs, fmt.Sprint(arg))
		}

		str += " " + strings.Join(strs, " ")
	}

	str += string(Delimiter)

	n, err := c.w.WriteString(str)
	if err != nil {
		return fmt.Errorf("write (n %d, len %d): %w", len(str), n, err)
	}

	if err := c.w.Flush(); err != nil {
		return fmt.Errorf("flush (n %d, len %d): %w", len(str), n, err)
	}

	return nil
}
