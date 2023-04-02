package elgo

import "sync"

type playerMap map[string]Player

type playerStore struct {
	store playerMap
	lock  sync.RWMutex
}

func newStore() *playerStore {
	return &playerStore{
		store: make(map[string]Player, 0),
	}
}

func (s *playerStore) All() playerMap {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.store
}

func (s *playerStore) Get(key string) (Player, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	v, ok := s.store[key]
	return v, ok
}

func (s *playerStore) Set(key string, value Player) {
	s.lock.Lock()
	s.store[key] = value
	s.lock.Unlock()
}

func (s *playerStore) Delete(key string) {
	s.lock.Lock()
	delete(s.store, key)
	s.lock.Unlock()
}
