package socket

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type Event string

const (
	Add     Event = "ADD"
	Match   Event = "MATCH"
	Remove  Event = "REMOVE"
	Size    Event = "SIZE"
	Unknown Event = ""
)

// parseEvent accepts a string event and parses it, returning Event type and
// the rest of the string.
//
// Events are expected to be in such format:
//
//	ADD ...
//	MATCH ...
//	SIZE
//
// If no known event type was found, Unknown is returned.
func parseEvent(c net.Conn) (Event, string, error) {
	buf := make([]byte, 0, 1024)
	tmp := make([]byte, 0, 1024)

	for {
		n, err := c.Read(tmp)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return Unknown, "", err
		}

		buf = append(buf, tmp[:n]...)
	}

	s := strings.TrimSpace(string(buf))
	split := strings.SplitN(s, " ", 1)
	switch split[0] {
	case "ADD":
		return Add, s[3:], nil
	case "MATCH":
		return Match, s[5:], nil
	case "REMOVE":
		return Remove, s[6:], nil
	case "SIZE":
		return Size, s[4:], nil
	default:
		return Unknown, s, nil
	}
}

func writeEvent(w net.Conn, e Event, args ...any) error {
	buf := bytes.Buffer{}
	buf.WriteString(string(e))
	for i, arg := range args {
		if i == 0 {
			buf.WriteByte(' ')
		}

		buf.WriteString(fmt.Sprint(arg))
	}

	_, err := w.Write(buf.Bytes())
	return err
}

func parseSize(s string) (int, error) {
	s = strings.TrimSpace(s)
	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse size: %w", err)
	}

	return int(size), nil
}
