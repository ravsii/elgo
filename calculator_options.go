package elgo

type calcOption func(*Calculator)

func WithKFactor(k float64, startsAt float64) calcOption {
	return func(c *Calculator) {
		c.kFactors = append(c.kFactors, kFactor{k: k, startsAt: startsAt})
	}
}
