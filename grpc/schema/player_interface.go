package schema

import "github.com/ravsii/elgo"

// This file simply implements player interface for gRPC player for convinience.

var _ elgo.Player = (*Player)(nil)

func (p *Player) Identify() string {
	return p.GetId()
}

func (p *Player) Rating() float64 {
	return float64(p.GetElo())
}
