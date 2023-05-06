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
	id, ratingStr, found := strings.Cut(s, ";")
	if !found {
		return nil, fmt.Errorf("%w: %s", ErrBadInput, s)
	}

	r, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil {
		return nil, fmt.Errorf("parse rating: %w", err)
	}

	return &SocketPlayer{ID: id, ELO: r}, nil
}

func decodePlayers(s string) ([]elgo.Player, error) {
	split := strings.Split(s, " ")
	players := make([]elgo.Player, 0, len(split))

	for _, playerStr := range split {
		decoded, err := decodePlayer(playerStr)
		if err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}

		players = append(players, decoded)
	}

	return players, nil
}
