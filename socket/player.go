package socket

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ravsii/elgo"
)

var ErrBadInput = errors.New("bad player input")

type SocketPlayer struct {
	ID  string
	ELO float64
}

func (p *SocketPlayer) Identify() string {
	return p.ID
}

func (p *SocketPlayer) Rating() float64 {
	return p.ELO
}

func encodePlayer(p elgo.Player) string {
	return fmt.Sprintf("%s;%f", p.Identify(), p.Rating())
}

func decodePlayer(s string) (elgo.Player, error) {
	data := strings.SplitN(s, ";", 1)
	if len(data) != 2 {
		return nil, ErrBadInput
	}

	r, err := strconv.ParseFloat(data[1], 64)
	if err != nil {
		return nil, fmt.Errorf("parse rating: %w", err)
	}

	return &SocketPlayer{ID: data[0], ELO: r}, nil
}
