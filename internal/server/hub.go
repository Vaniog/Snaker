package server

import (
	"context"
	"github.com/Vaniog/Snaker/internal/game"
	"github.com/Vaniog/Snaker/internal/server/event"
	"log"
	"sync"
	"time"
)

const gameLifeTime = 10 * time.Minute

type HubID int64
type Hub struct {
	ID       HubID
	game     *game.Game
	gameLock sync.RWMutex
	players  map[*Player]struct{}

	register chan *Player
	events   chan event.Event
}

func newHub() *Hub {
	return &Hub{
		game:     game.NewGame(game.DefaultOptions),
		register: make(chan *Player),
		players:  make(map[*Player]struct{}),
	}
}

func (h *Hub) Run() {
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(gameLifeTime),
	)
	defer cancel()
	log.Printf("Start hub: %d", h.ID)

	for {
		select {
		case p := <-h.register:
			log.Printf("New player in hub %d", h.ID)
			h.players[p] = struct{}{}
			snake := h.game.RegisterSnake()
			p.snake = snake
			go p.readPump(ctx)
			go p.writePump(ctx)
			if len(h.players) == 2 {
				log.Printf("Start game in hub %d", h.ID)
				go h.game.Run(ctx, func() {
					h.gameLock.RLock()
					data, _ := h.game.JSON()
					h.gameLock.RUnlock()
					for p := range h.players {
						p.broadcast <- data
					}
				})
			}
		case <-ctx.Done():
			return
		}
	}
}
