package elgo

import (
	"errors"
	"fmt"
)

var (
	ErrAlreadyExists = errors.New("player already exists")
	ErrPoolClosed    = errors.New("pool is closed")
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
