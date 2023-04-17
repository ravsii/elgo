package elgo

type CalcOpt func(*Calculator)

// WithKFactor adds one more range that uses a different K-factor.
func WithKFactor(startsAt float64, k float64) CalcOpt {
	return func(c *Calculator) {
		c.kFactors = append(c.kFactors, kFactor{k: k, startsAt: startsAt})
	}
}
