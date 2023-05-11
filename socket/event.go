package socket

type Event string

const (
	// Add is an event sent from the client to the server,
	// that adds a player to the pool.
	Add Event = "ADD"
	// Match is an event sent from the server to the client,
	// if a match was created.
	Match Event = "MATCH"
	// Remove is an event sent from the client to the server,
	// if a certain players leaves the pool.
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
