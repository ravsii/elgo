package elgo

import "time"

type OptionFunc func(*Pool)

// WithRetry sets a retry duration at which system would again
// try to find opponents for players.
func WithRetry(d time.Duration) OptionFunc {
	return func(p *Pool) {
		p.retrySearchIn = d
	}
}

// WithIncreaseInterval sets an amount of points that will be added
// on a new search, if no opponent was found previously
func WithIncreaseInterval(interval float64) OptionFunc {
	return func(p *Pool) {
		p.increaseRatingBorders = interval
	}
}
