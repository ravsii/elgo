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
func (f kFactors) Less(i, j int) bool { return f[i].k < f[j].k }

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
func NewCalc(k float64, options ...calcOption) *Calculator {
	c := &Calculator{k: k}

	for _, option := range options {
		option(c)
	}

	sort.Sort(c.kFactors)

	return c
}

// CalcRating calculates rating for winner and loser and calls SetRating for
// both of them.
func CalcRating(winner, loser Ratinger) {
	if winner == nil || loser == nil {
		return
	}

	winnerProb := (1.0 / (1.0 + math.Pow(10.0, ((loser.Rating()-winner.Rating())/400.0))))
	loserProb := (1.0 / (1.0 + math.Pow(10.0, ((winner.Rating()-loser.Rating())/400.0))))

	K := 30.0

	newWinnerRating := winner.Rating() + K*(1-winnerProb)
	newLoserRating := loser.Rating() + K*(-loserProb)

	winner.SetRating(newWinnerRating)
	loser.SetRating(newLoserRating)
}
