package elgo_test

import (
	"math"
	"testing"

	"github.com/ravsii/elgo"
)

func TestWinnerIsNil(t *testing.T) {
	t.Parallel()

	loserExpected := 1000.
	loser := CreatePlayerMock("1", loserExpected)

	elgo.CalcRating(nil, loser)

	if math.Ceil(loser.Rating()) != loserExpected {
		t.Errorf("loser rating: want %f got %f", loser.Rating(), loserExpected)
	}
}

func TestLoserIsNil(t *testing.T) {
	t.Parallel()

	winnerExpected := 1000.
	winner := CreatePlayerMock("1", winnerExpected)

	elgo.CalcRating(winner, nil)

	if math.Ceil(winner.Rating()) != winnerExpected {
		t.Errorf("winner rating: want %f got %f", winner.Rating(), winnerExpected)
	}
}

func TestCalcRating(t *testing.T) {
	t.Parallel()

	winner := CreatePlayerMock("1", 1200)
	loser := CreatePlayerMock("1", 1000)
	winnerExpected := 1208.
	loserExpected := 993.

	elgo.CalcRating(winner, loser)

	if math.Ceil(winner.Rating()) != winnerExpected {
		t.Errorf("winner rating: want %f got %f", winner.Rating(), winnerExpected)
	}

	if math.Ceil(loser.Rating()) != loserExpected {
		t.Errorf("loser rating: want %f got %f", loser.Rating(), loserExpected)
	}
}
