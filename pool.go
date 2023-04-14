package elgo

import (
	"errors"
	"math"
	"sync"
	"time"
)

var (
	ErrAlreadyExists = errors.New("player already exists")
)

type poolPlayer struct {
	player Player

	ratingBorders float64
	retryAt       time.Time
}

// Match is a struct that holds 2 players who should be matched.
type Match struct {
	Player1 Player
	Player2 Player
}

// Pool is a main struct for matchmaking pool.
// Use NewPool(options...) to create a new pool.
type Pool struct {
	players     map[string]*poolPlayer
	playersLock sync.RWMutex

	matchCh chan Match

	// retrySearchIn holds a duration that should be waited before iterations
	// if no match was found.
	retrySearchIn time.Duration

	// increaseRatingBorders shows by how many points seaching borders will be
	// increased when no opponent was found
	increaseRatingBorders float64
}

// NewPool creates a new pool for players.
// Pools aren't connected to each other so creating multiple of them is safe.
func NewPool(options ...OptionFunc) *Pool {
	p := &Pool{
		players: make(map[string]*poolPlayer),
		matchCh: make(chan Match),

		retrySearchIn:         5 * time.Second,
		increaseRatingBorders: 100,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

// AddPlayer returns a queue channel to send new players to.
// ErrAlreadyExists is returned if identifier is already taken.
func (p *Pool) AddPlayer(player Player) error {
	p.playersLock.Lock()
	defer p.playersLock.Unlock()

	id := player.Identify()
	if _, ok := p.players[id]; ok {
		return ErrAlreadyExists
	}

	p.players[id] = &poolPlayer{
		player:        player,
		ratingBorders: p.increaseRatingBorders,
		retryAt:       time.Now(),
	}

	return nil
}

// Matches returns a channel that sends found matches between players.
func (p *Pool) Matches() <-chan Match {
	return p.matchCh
}

// Size returns current amount of players in queue.
// It's concurrent-safe.
func (p *Pool) Size() int {
	p.playersLock.RLock()
	defer p.playersLock.RUnlock()
	return len(p.players)
}

// Close closes the pool and return players that are still in the queue.
func (p *Pool) Close() map[string]Player {
	close(p.matchCh)

	p.playersLock.Lock()
	playersLeft := make(map[string]Player, 0)
	for id, player := range p.players {
		playersLeft[id] = player.player
	}
	p.players = nil
	p.playersLock.Unlock()

	return playersLeft
}

func (p *Pool) Run() {
	ticker := time.NewTicker(p.retrySearchIn)

	for {
		if !p.iteration() {
			<-ticker.C
		}
	}
}

func (p *Pool) iteration() bool {
	if p.Size() < 2 {
		return false
	}

	p.playersLock.Lock()
	defer p.playersLock.Unlock()
	for id1, p1 := range p.players {
		for id2, p2 := range p.players {
			if id1 == id2 {
				continue
			}

			if ok := couldMatch(p1, p2); ok {
				p.createMatch(p1, p2)
				return true
			}

		}

		p1.ratingBorders += p.increaseRatingBorders
	}

	return false
}

func (p *Pool) createMatch(p1, p2 *poolPlayer) {
	p.removePlayersFromQueue(p1, p2)
	p.matchCh <- Match{Player1: p1.player, Player2: p2.player}
}

func (p *Pool) removePlayersFromQueue(players ...*poolPlayer) {
	for _, player := range players {
		delete(p.players, player.player.Identify())
	}
}

// couldMatch checks the difference between players' ratings
// AND their allowed interval of both.
//
// Difference (absolute value) in players' rating should be less or equals
// to the allowed rating borders of each player.
func couldMatch(p1, p2 *poolPlayer) bool {
	diff := math.Abs(p1.player.Rating() - p2.player.Rating())
	return p1.ratingBorders >= diff && p2.ratingBorders >= diff
}
