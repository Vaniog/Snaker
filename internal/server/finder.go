package server

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

const hubLifetime = 10 * time.Minute

type Repository struct {
	hubs map[HubID]*Hub
	lock sync.RWMutex
}

func NewHubRepository() *Repository {
	hr := Repository{
		hubs: make(map[HubID]*Hub),
	}
	return &hr
}

func (hr *Repository) FindHub() HubID {
	hr.lock.Lock()
	defer hr.lock.Unlock()

	for id, h := range hr.hubs {
		if len(h.game.Snakes) != 2 {
			return id
		}
	}

	id := HubID(rand.Int63())
	h := newHub()
	h.ID = id
	hr.hubs[id] = h

	go func() {
		time.Sleep(hubLifetime)
		hr.lock.Lock()
		defer hr.lock.Unlock()
		delete(hr.hubs, id)
	}()

	go h.Run()
	return id
}

func (hr *Repository) GetHubById(id HubID) (*Hub, error) {
	hr.lock.RLock()
	defer hr.lock.RUnlock()

	h, ok := hr.hubs[id]
	if !ok {
		return nil, errors.New("hub not found")
	}
	return h, nil
}
