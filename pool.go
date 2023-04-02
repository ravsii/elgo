package elgo

import (
	"fmt"
	"math"
	"sort"
	"sync/atomic"
	"time"
)

type Match struct {
	Player1 Player
	Player2 Player
}

type Pool struct {
	players *playerStore
	inQueue atomic.Int32
	queueSf chan int8

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
		players:   newStore(),
		playersCh: make(chan Player),
		matchCh:   make(chan Match),
		queueSf:   make(chan int8, 1),

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

	p.inQueue.Store(0)

	return playersLeft
}

func (p *Pool) acceptPlayers() {
	for player := range p.playersCh {
		p.players.Set(player.Identify(), player)
		p.inQueue.Add(1)

		go p.findMatchFor(player)
	}
}

func (p *Pool) findMatchFor(player Player) {
	t := time.NewTicker(p.retrySearchIn)
	interval := p.increaseRatingBorders

	for {
		p.queueSf <- 1
		if ok := p._findMatch(player, interval); ok {
			<-p.queueSf
			break
		}
		<-p.queueSf

		interval += p.increaseRatingBorders
		<-t.C
	}

	t.Stop()
}

func (p *Pool) _findMatch(player Player, allowedInterval float64) bool {
	identifier := player.Identify()

	for opponentIdent, opponent := range p.players.All() {
		if identifier == opponentIdent {
			continue
		}

		if ok := couldMatch(player, opponent, allowedInterval); ok {
			p.createMatch(player, opponent)
			return true
		}
	}

	return false
}

func (p *Pool) createMatch(p1, p2 Player) {
	p.removePlayersFromQueue(p1, p2)

	p.matchCh <- Match{Player1: p1, Player2: p2}
}

func (p *Pool) removePlayersFromQueue(players ...Player) {
	p.printMap()
	for _, player := range players {
		p.players.Delete(player.Identify())
		p.inQueue.Add(-1)
	}
	p.printMap()
}

func (p *Pool) printMap() {
	sl := make([]string, 0)

	for k := range p.players.All() {
		sl = append(sl, k)
	}
	sort.Slice(sl, func(i, j int) bool {
		return sl[i] < sl[j]
	})

	fmt.Println(p.Size(), ":", sl)
}

func couldMatch(p1, p2 Player, allowedInterval float64) bool {
	diff := math.Abs(p1.Rating() - p2.Rating())
	return allowedInterval > diff
}
