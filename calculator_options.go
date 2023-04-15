package elgo

type calcOption func(*Calculator)

func WithKFactor(startsAt float64, k float64) calcOption {
	return func(c *Calculator) {
		c.kFactors = append(c.kFactors, kFactor{k: k, startsAt: startsAt})
	}
}
