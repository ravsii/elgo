package socket

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ravsii/elgo"
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
	socketPlayer

	ELO float64
}

func (p *socketRatingPlayer) Rating() float64 {
	return p.ELO
}

func encodePlayer(p elgo.Player) string {
	return fmt.Sprintf("%s;%f", p.Identify(), p.Rating())
}

func decodeRatingPlayer(s string) (*socketRatingPlayer, error) {
	id, ratingStr, found := strings.Cut(s, ";")
	if !found {
		return nil, fmt.Errorf("%w: %s", ErrBadInput, s)
	}

	r, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil {
		return nil, fmt.Errorf("parse rating: %w", err)
	}

	return &socketRatingPlayer{socketPlayer: socketPlayer{ID: id}, ELO: r}, nil
}

func decodeRatingPlayers(s string) ([]elgo.Player, error) {
	split := strings.Split(s, " ")
	players := make([]elgo.Player, 0, len(split))

	for _, playerStr := range split {
		decoded, err := decodeRatingPlayer(playerStr)
		if err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}

		players = append(players, decoded)
	}

	return players, nil
}
