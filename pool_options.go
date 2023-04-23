package elgo

import "time"

type PoolOpt func(*Pool)

// WithPlayerRetry sets a duration that a player should wait before a pool
// should try and find a match for them again.
func WithPlayerRetry(d time.Duration) PoolOpt {
	return func(p *Pool) {
		p.playerRetryInterval = d
	}
}

// WithGlobalRetry sets a duration that a pool should wait between iterations
// if no match was found.
func WithGlobalRetry(d time.Duration) PoolOpt {
	return func(p *Pool) {
		p.playerRetryInterval = d
	}
}

// WithIncreaseInterval sets an amount of points that will be added
// on a new search, if no opponent was found previously
func WithIncreaseInterval(interval float64) PoolOpt {
	return func(p *Pool) {
		p.playersBordersIncreaseBy = interval
	}
}
