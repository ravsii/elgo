package elgo_test

import (
	"math"
	"testing"

	"github.com/ravsii/elgo"
	"github.com/ravsii/elgo/mocks"
	"github.com/stretchr/testify/mock"
)

func TestCalcRating(t *testing.T) {
	t.Parallel()

	winner := mocks.NewRatinger(t)
	loser := mocks.NewRatinger(t)

	winnerRating := 1200.0
	loserRating := 1000.0

	winner.On("Rating").Return(winnerRating)
	loser.On("Rating").Return(loserRating)

	winner.On("SetRating", mock.AnythingOfType("float64")).Run(func(args mock.Arguments) {
		winnerRating = args.Get(0).(float64)
	})

	loser.On("SetRating", mock.AnythingOfType("float64")).Run(func(args mock.Arguments) {
		loserRating = args.Get(0).(float64)
	})

	elgo.CalcRating(winner, loser)

	if math.Round(winnerRating) != 1207 || math.Round(loserRating) != 993 {
		t.Error("wrong results")
	}
}
