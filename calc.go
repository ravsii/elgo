package elgo

import "math"

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
