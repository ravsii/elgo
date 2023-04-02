package elgo_test

type TestPlayer struct {
	Name      string
	EloRating float64
}

func (p *TestPlayer) Identify() string {
	return p.Name
}

func (p *TestPlayer) Rating() float64 {
	return p.EloRating
}

func (p *TestPlayer) SetRating(rating float64) {
	p.EloRating = rating
}

func CreatePlayerMock(name string, rating float64) *TestPlayer {
	return &TestPlayer{Name: name, EloRating: rating}
}
