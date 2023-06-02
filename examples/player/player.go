package player

import (
	"fmt"

	"github.com/ravsii/elgo"
)

var compileTest elgo.Player = (*Player)(nil)

type Player struct {
	Name      string
	EloRating float64
}

func (p *Player) Identify() string {
	return p.Name
}

func (p *Player) Rating() float64 {
	return p.EloRating
}

func (p *Player) SetRating(rating float64) {
	p.EloRating = rating
}

// Functions below are for other examples, they aren't mandatory

func (p *Player) String() string {
	return fmt.Sprintf("%s (%.2f)", p.Name, p.EloRating)
}

func New(name string, rating float64) *Player {
	return &Player{name, rating}
}
