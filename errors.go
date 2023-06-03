package elgo

import (
	"errors"
	"fmt"
)

var (
	ErrPoolClosed   = errors.New("pool is closed")
	ErrNoMatchFound = errors.New("no match found")
)

type PlayerAlreadyExistsError struct {
	player Identifier
}

func NewAlreadyExistsErr(p Identifier) *PlayerAlreadyExistsError {
	return &PlayerAlreadyExistsError{player: p}
}

func (e *PlayerAlreadyExistsError) Error() string {
	return fmt.Sprintf("%s already in the queue", e.player.Identify())
}
