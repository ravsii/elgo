package calc

type CalcOpt func(*Calculator)

// WithKFactor adds one more range that uses a different K-factor.
// startsAt is inclusive.
func WithKFactor(startsAt float64, k float64) CalcOpt {
	return func(c *Calculator) {
		c.kFactors = append(c.kFactors, kFactor{k: k, startsAt: startsAt})
	}
}
