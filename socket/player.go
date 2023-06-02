package socket

import (
	"errors"
)

var ErrBadInput = errors.New("bad player input")

// socketRatingPlayer is a VERY basic implementation of a player (Identifier).
type socketPlayer struct {
	ID string
}

func (p *socketPlayer) Identify() string {
	return p.ID
}

// socketRatingPlayer is a VERY basic implementation of a player with ELO rating.
type socketRatingPlayer struct {
	ID  string
	ELO float64
}

func (p *socketRatingPlayer) Identify() string {
	return p.ID
}

func (p *socketRatingPlayer) Rating() float64 {
	return p.ELO
}
