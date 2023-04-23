package elgo

import (
	"math"
	"sort"
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
//	elgo.NewCalc(...)
//
// to create a new calculator
type Calculator struct {
	k        float64
	kFactors kFactors
}

// NewCalc creates a new calculator for rating changes.
// k is used as a default option and you can pass multiple k factors that
// will be applied depending on players' rating. Thus:
//
//	NewCalc(5, WithKFactor(1000, 20), WithKFactor(1500, 30))
//
// will be used as:
//
//	rating > 0 && rating < 1000, k = 5
//	rating >= 1000 && rating < 1500, k = 20
//	rating >= 1500, k = 30
//
// There's no option to add max value for rating ranges,
// it's either infinite or next factor's min value.
func NewCalc(k float64, opts ...CalcOpt) *Calculator {
	c := &Calculator{k: k}

	for _, option := range opts {
		option(c)
	}

	sort.Sort(c.kFactors)

	return c
}

// Win calculates rating for winner and loser.
func (c *Calculator) Win(winner, loser Ratinger) (newWinnerRating, newLoserRating float64) {
	if winner == nil || loser == nil {
		return 0, 0
	}

	newWinnerRating = winner.Rating() + c.getK(winner)*(1-probability(winner, loser))
	newLoserRating = loser.Rating() + c.getK(loser)*(-probability(loser, winner))

	return newWinnerRating, newLoserRating
}

// Draw calculates rating for both players using a 0.5 co-efficient.
func (c *Calculator) Draw(p1, p2 Ratinger) (newP1Rating, newP2Rating float64) {
	if p1 == nil || p2 == nil {
		return 0, 0
	}

	newP1Rating = p1.Rating() + c.getK(p1)*(0.5-probability(p1, p2))
	newP2Rating = p2.Rating() + c.getK(p2)*(0.5-probability(p2, p1))

	return newP1Rating, newP2Rating
}

func (c *Calculator) getK(r Ratinger) float64 {
	for _, bracket := range c.kFactors {
		if bracket.startsAt >= r.Rating() {
			return bracket.k
		}
	}

	return c.k
}

// probability counts probability of winning of a player against opponent
// using this formula: https://en.wikipedia.org/wiki/Elo_rating_system#Mathematical_details
func probability(player, opponent Ratinger) float64 {
	// note: can't use Pow10 here because Rating() returns float64
	return 1.0 / (1.0 + math.Pow(10.0, (opponent.Rating()-player.Rating())/400.0))
}
