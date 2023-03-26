package player

type Player struct {
	Name      string
	EloRating float64
}

func (p *Player) Identify() string {
	return p.Name
}

func (p *Player) Rating() float64 {
	return p.EloRating
}

func (p *Player) SetRating(rating float64) {
	p.EloRating = rating
}

func (p *Player) String() string {
	return p.Name
}
