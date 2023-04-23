package elgo

import "time"

type PoolOpt func(*Pool)

// WithPlayerRetryInterval sets a duration of how much time a player
// should wait before a pool would try and find a match for them again.
func WithPlayerRetryInterval(d time.Duration) PoolOpt {
	return func(p *Pool) {
		p.playerRetryInterval = d
	}
}

// WithGlobalRetryInterval sets a duration of how much time a pool
// should wait between iterations if not a single match was found.
func WithGlobalRetryInterval(d time.Duration) PoolOpt {
	return func(p *Pool) {
		p.globalRetryInterval = d
	}
}

// WithIncreasePlayerBorderBy sets an amount of points that will be added
// on a new search, if no opponent was found previously.
func WithIncreasePlayerBorderBy(interval float64) PoolOpt {
	return func(p *Pool) {
		p.playersBordersIncreaseBy = interval
	}
}
