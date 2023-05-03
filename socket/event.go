package socket

import (
	"bytes"
	"fmt"
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
func parseEvent(s string) (Event, string) {
	split := strings.SplitN(strings.TrimSpace(s), " ", 1)
	switch split[0] {
	case "ADD":
		return Add, split[1]
	case "MATCH":
		return Match, split[1]
	case "REMOVE":
		return Remove, split[1]
	case "SIZE":
		return Size, ""
	default:
		return Unknown, split[1]
	}
}

func createEvent(e Event, args ...any) []byte {
	buf := bytes.Buffer{}
	buf.WriteString(string(e))
	for _, arg := range args {
		buf.WriteString(fmt.Sprint(arg))
	}

	return buf.Bytes()
}
