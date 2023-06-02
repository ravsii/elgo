package elgo

// Identifier is an interface that helps identify players.
type Identifier interface {
	// Identify should return any kind of a (unique) identifier for a player.
	Identify() string
}

// Ratinger is an interface to receive and change player's rating.
type Ratinger interface {
	// Rating should return player's ELO rating.
	Rating() float64
}

// Player is an interface that implements Identifier and Ratinger.
type Player interface {
	Identifier
	Ratinger
}

type BasePlayer struct {
	ID string
}

func (p *BasePlayer) Identify() string {
	return p.ID
}

type BaseRatingPlayer struct {
	ID  string
	ELO float64
}

func (p *BaseRatingPlayer) Identify() string {
	return p.ID
}

func (p *BaseRatingPlayer) Rating() float64 {
	return p.ELO
}
