package socket

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Event string

const (
	// Add is a server-only event that adds a player to the pool
	Add Event = "ADD"
	// Match is a client-only event if a match was created.
	Match Event = "MATCH"
	// Remove is a server-only event if a certain players leaves the pool.
	Remove Event = "REMOVE"
	// Size is a double-size event:
	//	- When sent from the client to the server, it will ask for
	//	the current amount of players in queue.
	//	It should be sent without arguments (SIZE)
	// - When sent from the server to the client, it returns the current of
	//	players in queue (SIZE 10)
	Size Event = "SIZE"
	// Unknown event is returned if no other event prefix was detected.
	Unknown Event = ""
)

const Delimiter = "\r\n"

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
	// var buf bytes.Buffer
	// io.Copy(&buf, c)
	s, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		return Unknown, "", err
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

func writeEvent(w net.Conn, event Event, args ...any) error {
	buf := append([]byte(event), ' ')

	if len(args) > 0 {
		strs := make([]string, 0, len(args))
		for _, arg := range args {
			strs = append(strs, fmt.Sprint(arg))
		}

		str := strings.Join(strs, " ")
		buf = append(buf, []byte(str)...)
	}

	buf = append(buf, []byte("\r\n")...)

	n, err := w.Write(buf)
	if err != nil {
		return fmt.Errorf("unable to write to net.Conn: %w", err)
	}

	// -2 to remove \r\n
	log.Println("write", len(buf), n, string(buf[:len(buf)-2]))

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
