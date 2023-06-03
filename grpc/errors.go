package grpc

import (
	"fmt"

	"github.com/ravsii/elgo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrCreateMatch = status.Error(codes.Aborted, "can't create match")

func NewAlreadyExistsErr(p elgo.Identifier) error {
	return status.Error(
		codes.AlreadyExists,
		fmt.Sprintf("%s already in the queue", p.Identify()))
}

func NewCantAddErr(p elgo.Identifier) error {
	return status.Error(
		codes.Aborted,
		fmt.Sprintf("can't add %s to the queue", p.Identify()))
}
