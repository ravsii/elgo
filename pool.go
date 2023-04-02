package elgo

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type Match struct {
	Player1 Player
	Player2 Player
}

type Pool struct {
	players   sync.Map
	inQueue   atomic.Int32
	rangeLock sync.Mutex

	playersCh chan Player
	matchCh   chan Match

	retrySearchIn time.Duration

	// increaseRatingBorders shows by how many points seaching borders will be
	// increased when no opponent was found
	increaseRatingBorders float64
}

// NewPool creates a new pool for players.
// Pools aren't connected to each other so creating multiple of them is safe.
func NewPool(options ...OptionFunc) *Pool {
	p := &Pool{
		playersCh: make(chan Player),
		matchCh:   make(chan Match),

		retrySearchIn:         5 * time.Second,
		increaseRatingBorders: 100,
	}

	for _, option := range options {
		option(p)
	}

	go p.acceptPlayers()

	return p
}

// Queue returns a queue channel to send new players to.
func (p *Pool) Queue() chan<- Player {
	return p.playersCh
}

// Matches returns a channel that sends found matches between players.
func (p *Pool) Matches() <-chan Match {
	return p.matchCh
}

// Size returns current amount of players in queue
func (p *Pool) Size() int32 {
	return p.inQueue.Load()
}

// Close closes the pool and return players that are still in the queue.
func (p *Pool) Close() map[string]Player {
	close(p.playersCh)
	close(p.matchCh)

	playersLeft := make(map[string]Player, 0)
	p.players.Range(func(key, value any) bool {
		playersLeft[key.(string)] = value.(Player)
		p.players.Delete(key)
		return true
	})

	p.inQueue.Store(0)

	return playersLeft
}

func (p *Pool) acceptPlayers() {
	for player := range p.playersCh {
		p.players.Store(player.Identify(), player)
		p.inQueue.Add(1)

		go p.findMatchFor(player)
	}
}

func (p *Pool) findMatchFor(player Player) {
	t := time.NewTicker(p.retrySearchIn)
	interval := p.increaseRatingBorders

	for {
		if ok := p._findMatch(player, interval); ok {
			break
		}

		interval += p.increaseRatingBorders
		<-t.C
	}

	t.Stop()
}

func (p *Pool) _findMatch(player Player, allowedInterval float64) bool {
	identifier := player.Identify()
	found := false

	p.rangeMap(func(opponentIdent string, opponent Player) bool {
		if identifier == opponentIdent {
			return true
		}

		if ok := couldMatch(player, opponent, allowedInterval); ok {
			p.createMatch(player, opponent)
			found = true
			return false
		}

		return true
	})

	return found
}

func (p *Pool) rangeMap(fn func(ident string, player Player) bool) {
	p.rangeLock.Lock()
	defer p.rangeLock.Unlock()

	p.players.Range(func(k, v any) bool {
		return fn(k.(string), v.(Player))
	})
}

func (p *Pool) createMatch(p1, p2 Player) {
	p.removePlayersFromQueue(p1, p2)

	p.matchCh <- Match{Player1: p1, Player2: p2}
}

func (p *Pool) removePlayersFromQueue(players ...Player) {
	for _, player := range players {
		p.players.Delete(player.Identify())
		p.inQueue.Add(-1)
	}
}

func couldMatch(p1, p2 Player, allowedInterval float64) bool {
	diff := math.Abs(p1.Rating() - p2.Rating())
	return allowedInterval > diff
}
