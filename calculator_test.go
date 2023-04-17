package elgo_test

import (
	"math"
	"testing"

	"github.com/ravsii/elgo"
)

func TestCalculations(t *testing.T) {
	t.Parallel()

	t.Run("winner is nil", func(t *testing.T) {
		t.Parallel()

		calc := elgo.NewCalc(1)
		loser := CreatePlayerMock("loser", 1)
		winnerRating, loserRating := calc.Win(nil, loser)

		if winnerRating != loserRating || loserRating != 0 {
			t.Errorf("loser rating: want 0 got %f", loser.Rating())
		}
	})

	t.Run("loser is nil", func(t *testing.T) {
		t.Parallel()

		calc := elgo.NewCalc(1)
		winner := CreatePlayerMock("winner", 1)
		winnerRating, loserRating := calc.Win(nil, winner)

		if winnerRating != loserRating || winnerRating != 0 {
			t.Errorf("winner rating: want 0 got %f", winner.Rating())
		}
	})

	t.Run("win", func(t *testing.T) {
		t.Parallel()

		calc := elgo.NewCalc(30)
		winner := CreatePlayerMock("1", 1200)
		loser := CreatePlayerMock("2", 1000)
		winnerExpected := 1208.
		loserExpected := 993.

		winnerRating, loserRating := calc.Win(winner, loser)

		if math.Ceil(winnerRating) != winnerExpected {
			t.Errorf("winner rating: want %f got %f", winnerExpected, winner.Rating())
		}

		if math.Ceil(loserRating) != loserExpected {
			t.Errorf("loser rating: want %f got %f", loserExpected, loser.Rating())
		}
	})

	t.Run("win brackets", func(t *testing.T) {
		t.Parallel()

		calc := elgo.NewCalc(0, elgo.WithKFactor(1000, 20), elgo.WithKFactor(2000, 40))
		winner := CreatePlayerMock("1", 1000)
		loser := CreatePlayerMock("2", 2000)
		winnerExpected := 1020.
		loserExpected := 1961.

		winnerRating, loserRating := calc.Win(winner, loser)

		if math.Ceil(winnerRating) != winnerExpected {
			t.Errorf("winner rating: want %f got %f", winnerExpected, winner.Rating())
		}

		if math.Ceil(loserRating) != loserExpected {
			t.Errorf("loser rating: want %f got %f", loserExpected, loser.Rating())
		}
	})

	t.Run("draw", func(t *testing.T) {
		t.Parallel()

		calc := elgo.NewCalc(40)
		p1 := CreatePlayerMock("1", 2000)
		p2 := CreatePlayerMock("2", 1000)
		p1Expected := 1981.
		p2Expected := 1020.

		p1Rating, p2Rating := calc.Draw(p1, p2)

		if math.Ceil(p1Rating) != p1Expected {
			t.Errorf("p1 rating: want %f got %f", p1Expected, p1.Rating())
		}

		if math.Ceil(p2Rating) != p2Expected {
			t.Errorf("p2 rating: want %f got %f", p2Expected, p2.Rating())
		}
	})

	t.Run("draw p1 is nil", func(t *testing.T) {
		t.Parallel()

		calc := elgo.NewCalc(1)
		p2 := CreatePlayerMock("p1", 1)
		p1Rating, p2Rating := calc.Win(nil, p2)

		if p1Rating != p2Rating || p2Rating != 0 {
			t.Errorf("p2 rating: want 0 got %f", p2.Rating())
		}
	})

	t.Run("draw p2 is nil", func(t *testing.T) {
		t.Parallel()

		calc := elgo.NewCalc(1)
		p1 := CreatePlayerMock("p1", 1)
		p1Rating, p2Rating := calc.Win(p1, nil)

		if p1Rating != p2Rating || p2Rating != 0 {
			t.Errorf("p1 rating: want 0 got %f", p1.Rating())
		}
	})

	t.Run("draw brackets", func(t *testing.T) {
		t.Parallel()

		calc := elgo.NewCalc(0, elgo.WithKFactor(1000, 1), elgo.WithKFactor(2000, 40))
		p1 := CreatePlayerMock("1", 2000)
		p2 := CreatePlayerMock("2", 1000)
		p1Expected := 1981.
		p2Expected := 1001.

		p1Rating, p2Rating := calc.Draw(p1, p2)

		if math.Ceil(p1Rating) != p1Expected {
			t.Errorf("p1 rating: want %f got %f", p1Expected, p1.Rating())
		}

		if math.Ceil(p2Rating) != p2Expected {
			t.Errorf("p2 rating: want %f got %f", p2Expected, p2.Rating())
		}
	})
}
