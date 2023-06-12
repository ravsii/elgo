package calc

import (
	"math"
	"testing"

	"github.com/ravsii/elgo"
)

func TestWinnerNil(t *testing.T) {
	t.Parallel()

	calc := New(1)
	loser := elgo.BaseRatingPlayer{ID: "loser", ELO: 1}
	winnerRating, loserRating := calc.Win(nil, &loser)

	if winnerRating != loserRating || loserRating != 0 {
		t.Errorf("loser rating: want 0 got %f", loser.Rating())
	}
}

func TestLoserNil(t *testing.T) {
	t.Parallel()

	calc := New(1)
	winner := elgo.BaseRatingPlayer{ID: "winner", ELO: 1}
	winnerRating, loserRating := calc.Win(nil, &winner)

	if winnerRating != loserRating || winnerRating != 0 {
		t.Errorf("winner rating: want 0 got %f", winner.Rating())
	}
}

func TestWin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		k              float64
		funcs          []CalcOpt
		p1, p2         float64
		p1want, p2want float64
	}{
		{"base 0, empty", 0, nil, 1000, 2000, 1000, 2000},
		{"base 30, empty", 30, nil, 1200, 1000, 1208, 993},
		{
			"base 0, 1000-20, 2000-40", 0,
			[]CalcOpt{
				WithKFactor(1000, 20),
				WithKFactor(2000, 40),
			},
			1000, 2000, 1020, 1961,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			calc := New(tt.k, tt.funcs...)
			p1 := elgo.BaseRatingPlayer{ID: "p1", ELO: tt.p1}
			p2 := elgo.BaseRatingPlayer{ID: "p2", ELO: tt.p2}
			p1got, p2got := calc.Win(&p1, &p2)

			if math.Ceil(p1got) != tt.p1want {
				t.Errorf("p1 rating: want %f got %f", tt.p1want, p1.Rating())
			}

			if math.Ceil(p2got) != tt.p2want {
				t.Errorf("p2 rating: want %f got %f", tt.p2want, p2.Rating())
			}
		})
	}
}

func TestDrawP1Nil(t *testing.T) {
	t.Parallel()

	calc := New(1)
	p1 := elgo.BaseRatingPlayer{ID: "p1", ELO: 1}
	p1Rating, p2Rating := calc.Draw(&p1, nil)

	if p1Rating != p2Rating || p2Rating != 0 {
		t.Errorf("p1 rating: want 0 got %f", p1.Rating())
	}
}

func TestDrawP2Nil(t *testing.T) {
	t.Parallel()

	calc := New(1)
	p2 := elgo.BaseRatingPlayer{ID: "p2", ELO: 1}
	p1Rating, p2Rating := calc.Draw(nil, &p2)

	if p1Rating != p2Rating || p2Rating != 0 {
		t.Errorf("p2 rating: want 0 got %f", p2.Rating())
	}
}

func TestDraw(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		k              float64
		funcs          []CalcOpt
		p1, p2         float64
		p1want, p2want float64
	}{
		{"base 0, empty", 0, nil, 1000, 2000, 1000, 2000},
		{"base 40, empty", 40, nil, 2000, 1000, 1981, 1020},
		{
			"base 0, 1000-1, 2000-40", 0,
			[]CalcOpt{
				WithKFactor(1000, 1),
				WithKFactor(2000, 40),
			},
			2000, 1000, 1981, 1001,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			calc := New(tt.k, tt.funcs...)
			p1 := elgo.BaseRatingPlayer{ID: "p1", ELO: tt.p1}
			p2 := elgo.BaseRatingPlayer{ID: "p2", ELO: tt.p2}
			p1got, p2got := calc.Draw(&p1, &p2)

			if math.Ceil(p1got) != tt.p1want {
				t.Errorf("p1 rating: want %f got %f", tt.p1want, p1.Rating())
			}

			if math.Ceil(p2got) != tt.p2want {
				t.Errorf("p2 rating: want %f got %f", tt.p2want, p2.Rating())
			}
		})
	}
}
