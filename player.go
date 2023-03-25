package elgo

// Identifier is an interface that helps identify players.
type Identifier interface {
	// Identify should returns any kind of a unique identifier between players.
	Identify() string
}

// Ratinger is an interface to receive and change player's rating.
type Ratinger interface {
	// Rating should return player's rating.
	Rating() float64

	// SetRating should change the rating of a player.
	SetRating(float64)
}

// Player is an interface that implements Identifier and Ratinger.
type Player interface {
	Identifier
	Ratinger
}
