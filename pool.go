package elgo

import (
	"math"
	"sync"
	"time"
)

type poolPlayer struct {
	retryAt       time.Time
	player        Player
	ratingBorders float64
}

// Match is a struct that holds 2 players who should be matched.
type Match struct {
	Player1 Identifier
	Player2 Identifier
}

type Players map[string]Player

// Pool is a main struct for matchmaking pool.
// Use NewPool(options...) to create a new pool.
type Pool struct {
	players     map[string]*poolPlayer
	matchCh     chan Match
	playersLock sync.RWMutex

	// playerRetryInterval holds a duration of how much time a player
	// should wait before the next try if no match was found.
	playerRetryInterval time.Duration

	// globalRetryInterval holds a duration how much time a pool
	// should wait if not a single match was found (created).
	globalRetryInterval time.Duration

	// playersBordersIncreaseBy shows by how many points seaching borders
	// will be increased when no opponent was found
	playersBordersIncreaseBy float64
}

// NewPool creates a new pool for players.
// Pools aren't connected to each other so creating multiple of them is safe.
// To close a pool use pool.Close().
func NewPool(opts ...PoolOpt) *Pool {
	p := &Pool{
		players: make(map[string]*poolPlayer),
		matchCh: make(chan Match),

		playerRetryInterval:      time.Second,
		globalRetryInterval:      time.Second,
		playersBordersIncreaseBy: 100,
	}

	for _, option := range opts {
		option(p)
	}

	return p
}

// AddPlayer returns a queue channel to send new players to.
// ErrAlreadyExists is returned if identifier is already taken.
// ErrPoolClosed is returned if the pool is closed.
func (p *Pool) AddPlayer(players ...Player) error {
	p.playersLock.Lock()
	defer p.playersLock.Unlock()

	select {
	case <-p.matchCh:
		return ErrPoolClosed
	default:
	}

	for _, player := range players {
		id := player.Identify()
		if _, ok := p.players[id]; ok {
			return NewAlreadyExistsErr(player)
		}

		p.players[id] = &poolPlayer{
			player:        player,
			ratingBorders: p.playersBordersIncreaseBy,
			retryAt:       time.Now(),
		}
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
// It's safe to call Close() multiple times, but in that case nil will be returned.
func (p *Pool) Close() Players {
	p.playersLock.Lock()
	defer p.playersLock.Unlock()

	select {
	case <-p.matchCh:
		return nil
	default:
		close(p.matchCh)
	}

	playersLeft := make(map[string]Player, 0)
	for id, player := range p.players {
		playersLeft[id] = player.player
	}
	p.players = nil

	return playersLeft
}

// Remove removes players from queue. It's concurrency-safe.
func (p *Pool) Remove(players ...Identifier) {
	p.playersLock.Lock()
	defer p.playersLock.Unlock()

	for _, player := range players {
		delete(p.players, player.Identify())
	}
}

// RemoveStrs is a copy of Remove but it accepts strings instead of Identifier.
// It's concurrency-safe.
func (p *Pool) RemoveStrs(players ...string) {
	p.playersLock.Lock()
	defer p.playersLock.Unlock()

	for _, player := range players {
		delete(p.players, player)
	}
}

// Run start an infinite loop for matchmaking. Usually it's a good idea to
// use it as a goroutine:
//
//	go pool.Run()
//
// And when you need to close it, use:
//
//	playersInQueue := pool.Close()
func (p *Pool) Run() {
	ticker := time.NewTicker(p.globalRetryInterval)

	for {
		select {
		case <-p.matchCh:
			return
		default:
			if !p.iteration() {
				<-ticker.C
			}
		}
	}
}

// iteration returns true if a match was found.
func (p *Pool) iteration() bool {
	if p.Size() < 2 {
		return false
	}

	p.playersLock.Lock()
	defer p.playersLock.Unlock()
	for id1, p1 := range p.players {
		// skipping a player if his retry time is still "on cooldown".
		if p1.retryAt.Compare(time.Now()) >= 0 {
			continue
		}

		for id2, p2 := range p.players {
			if id1 == id2 {
				continue
			}

			if ok := couldMatch(p1, p2); ok {
				p.createMatch(p1, p2)
				return true
			}
		}

		p1.ratingBorders += p.playersBordersIncreaseBy
		p1.retryAt = p1.retryAt.Add(p.playerRetryInterval)
	}

	return false
}

// createMatch removes two players from queue and sends them to the match channel.
func (p *Pool) createMatch(p1, p2 *poolPlayer) {
	select {
	case p.matchCh <- Match{Player1: p1.player, Player2: p2.player}:
		p.removePlayersFromQueue(p1, p2)
	default:
	}
}

// removePlayersFromQueue removes players from queue.
// It's a part of p.iteration(), which already locked the map, so no lock needed.
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
