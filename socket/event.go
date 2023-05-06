package socket

import (
	"bufio"
	"fmt"
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

const Delimiter = '\n'

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
	s, err := bufio.NewReader(c).ReadString(Delimiter)
	if err != nil {
		return Unknown, "", err
	}

	s = strings.TrimSpace(s)
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

func writeEvent(w net.Conn, event Event, args ...any) error {
	writer := bufio.NewWriter(w)

	if _, err := writer.WriteString(string(event)); err != nil {
		return fmt.Errorf("unable to write to net.Conn: %w", err)
	}

	if err := writer.WriteByte(' '); err != nil {
		return fmt.Errorf("unable to write to net.Conn: %w", err)
	}

	strs := make([]string, 0, len(args))

	for _, arg := range args {
		s, ok := arg.(string)
		if !ok {
			return fmt.Errorf("unable to convert %v to string", arg)
		}

		strs = append(strs, s)
	}

	if _, err := writer.WriteString(strings.Join(strs, " ")); err != nil {
		return fmt.Errorf("unable to write to net.Conn: %w", err)
	}

	if err := writer.WriteByte(Delimiter); err != nil {
		return fmt.Errorf("unable to write to net.Conn: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("unable to write to net.Conn: %w", err)
	}

	return nil
}

func parseSize(s string) (int, error) {
	s = strings.TrimSpace(s)
	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse size: %w", err)
	}

	return int(size), nil
}
