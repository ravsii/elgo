package calc

import (
	"math"
	"sort"

	"github.com/ravsii/elgo"
)

type kFactor struct {
	k        float64
	startsAt float64
}

// kFactors is a slice of k factors. Slice should be always sorted from lowest
// factors to highest.
type kFactors []kFactor

func (f kFactors) Len() int           { return len(f) }
func (f kFactors) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f kFactors) Less(i, j int) bool { return f[i].startsAt < f[j].startsAt }

// Calculator holds options for calculations. Use
//
//	calc.New(...)
//
// to create a new calculator.
type Calculator struct {
	kFactors kFactors
	k        float64
}

// New creates a new calculator for rating changes.
// k is used as a default option and you can pass multiple k factors that
// will be applied depending on players' rating. Thus:
//
//	New(5, WithKFactor(1000, 20), WithKFactor(1500, 30))
//
// will be used as:
//
//	rating > 0 && rating < 1000, k = 5
//	rating >= 1000 && rating < 1500, k = 20
//	rating >= 1500, k = 30
//
// K-factors are sorted automatically.
//
// There's no option to add max value for rating ranges,
// it's either infinite or next factor's min value.
func New(k float64, opts ...CalcOpt) *Calculator {
	c := &Calculator{k: k}

	for _, option := range opts {
		option(c)
	}

	sort.Sort(c.kFactors)

	return c
}

// Win calculates new rating for winner and loser.
// If one if the players is nil, then `0, 0` is returned.
func (c *Calculator) Win(winner, loser elgo.Ratinger) (winnerNew, loserNew float64) {
	if winner == nil || loser == nil {
		return 0, 0
	}

	return c.WinFloat(winner.Rating(), loser.Rating())
}

// WinFloat acts like Win() but accepts floats instead of interfaces.
func (c *Calculator) WinFloat(winner, loser float64) (winnerNew, loserNew float64) {
	winnerNew = winner + c.getK(winner)*(1-probability(winner, loser))
	loserNew = loser + c.getK(loser)*(-probability(loser, winner))

	return winnerNew, loserNew
}

// Draw calculates rating for both players using a 0.5 co-efficient.
// If one if the players is nil, then `0, 0` is returned.
func (c *Calculator) Draw(player1, player2 elgo.Ratinger) (player1New, player2New float64) {
	if player1 == nil || player2 == nil {
		return 0, 0
	}

	return c.DrawFloat(player1.Rating(), player2.Rating())
}

// DrawFloat acts like Draw() but accepts floats instead of interfaces.
func (c *Calculator) DrawFloat(player1, player2 float64) (player1New, player2New float64) {
	p1New := player1 + c.getK(player1)*(0.5-probability(player1, player2))
	p2New := player2 + c.getK(player2)*(0.5-probability(player2, player1))

	return p1New, p2New
}

func (c *Calculator) getK(rating float64) float64 {
	for _, bracket := range c.kFactors {
		if bracket.startsAt >= rating {
			return bracket.k
		}
	}

	return c.k
}

// probability counts probability of winning of a player against opponent
// using this formula: https://en.wikipedia.org/wiki/Elo_rating_system#Mathematical_details
func probability(player, opponent float64) float64 {
	// note: can't use Pow10 here because Rating() returns float64
	return 1.0 / (1.0 + math.Pow(10.0, (opponent-player)/400.0))
}
