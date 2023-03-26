package elgo

import (
	"sync"
	"time"

	"golang.org/x/exp/slices"
)

type Match struct {
	Player1 Player
	Player2 Player
}

type Pool struct {
	players     []Player
	playersLock sync.RWMutex

	playersCh chan Player
	matchCh   chan Match
}

func NewPool() *Pool {
	p := &Pool{
		playersCh: make(chan Player),
		matchCh:   make(chan Match, 1),
	}

	go p.matchmaking()
	go p.findMatchLoop()

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

// Close closes the pool and return players that are still in the queue.
func (p *Pool) Close() []Player {
	close(p.playersCh)
	close(p.matchCh)

	p.playersLock.RLock()
	players := make([]Player, len(p.players))
	copy(players, p.players)
	p.playersLock.RUnlock()

	p.players = nil

	return players
}

func (p *Pool) matchmaking() {
	for player := range p.playersCh {
		p.playersLock.Lock()
		p.players = append(p.players, player)
		p.playersLock.Unlock()
	}
}

func (p *Pool) findMatchLoop() {
	for {
		p.findMatch()
		time.Sleep(5 * time.Second)
	}
}

func (p *Pool) findMatch() {
	p.playersLock.RLock()
	defer p.playersLock.RUnlock()

	if len(p.players) == 0 {
		return
	}

	player := p.players[0]

	for _, opponent := range p.players {
		if player.Identify() == opponent.Identify() {
			continue
		}

		p.createMatch(player, opponent)
		break
	}
}

func (p *Pool) createMatch(p1, p2 Player) bool {
	p.playersLock.Lock()
	p1i := slices.Index(p.players, p1)
	p.players = append(p.players[0:p1i], p.players[p1i+1:]...)

	p2i := slices.Index(p.players, p2)
	p.players = append(p.players[0:p2i], p.players[p2i+1:]...)

	p.players = slices.Clip(p.players)
	p.playersLock.Unlock()

	p.matchCh <- Match{Player1: p1, Player2: p2}

	return true
}
