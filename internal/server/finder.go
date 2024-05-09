package server

import (
	"context"
	"errors"
	"github.com/Vaniog/Snaker/internal/game"
	"math/rand"
	"sync"
	"time"
)

const hubLifetime = 20 * time.Minute

type HubRepository struct {
	hubs map[HubID]*Hub
	lock sync.RWMutex
}

func NewHubRepository() *HubRepository {
	hr := HubRepository{
		hubs: make(map[HubID]*Hub),
	}
	return &hr
}

func (hr *HubRepository) FindHub(opts game.Options) HubID {
	hr.lock.Lock()
	defer hr.lock.Unlock()

	for id, h := range hr.hubs {
		if h.lobby.IsOpen() && h.lobby.Opts == opts {
			return id
		}
	}

	id := HubID(rand.Int63())
	h := newHub(opts)
	h.ID = id
	hr.hubs[id] = h

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(hubLifetime))
	go func() {
		time.Sleep(hubLifetime)
		hr.lock.Lock()
		defer hr.lock.Unlock()
		delete(hr.hubs, id)
		cancel()
	}()
	go h.Run(ctx)
	return id
}

func (hr *HubRepository) GetHubById(id HubID) (*Hub, error) {
	hr.lock.RLock()
	defer hr.lock.RUnlock()

	h, has := hr.hubs[id]
	if !has {
		return nil, errors.New("hub not found")
	}
	return h, nil
}
