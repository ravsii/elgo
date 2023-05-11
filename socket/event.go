package socket

import (
	"fmt"
	"strconv"
	"strings"
)

type Event string

const (
	// Add is a server-only event that adds a player to the pool.
	Add Event = "ADD"
	// Match is a client-only event if a match was created.
	Match Event = "MATCH"
	// Remove is a server-only event if a certain players leaves the pool.
	Remove Event = "REMOVE"
	// Size is a double-side event:
	//	- When sent from the client to the server, it will ask for
	//	the current amount of players in queue.
	//	It should be sent without arguments (SIZE)
	// - When sent from the server to the client, it returns the current of
	//	players in queue (SIZE 10)
	Size Event = "SIZE"
	// Unknown event is returned if no other event prefix was detected.
	Unknown Event = ""
)

func parseSize(s string) (int, error) {
	s = strings.TrimSpace(s)
	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse size: %w", err)
	}

	return int(size), nil
}
