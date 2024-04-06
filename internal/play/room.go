package play

import (
	"context"
	"github.com/Vaniog/Snaker/internal/game"
	"sync"
	"time"
)

type Room struct {
	game     *game.Game
	gameLock sync.RWMutex

	players  []*Player
	Register chan *Player
}

func NewRoom(g *game.Game) *Room {
	return &Room{
		game:     g,
		Register: make(chan *Player),
	}
}

func (r *Room) Run(ctx context.Context) {
	for {
		select {
		case p := <-r.Register:
			r.players = append(r.players, p)
			snake := r.game.RegisterSnake()

			p.snake = snake
			go p.inputPump(ctx)
			if len(r.players) == 2 {
				goto PLAY
			}
		case <-ctx.Done():
			return
		}
	}

PLAY:
	gameTicker := time.NewTicker(r.game.Opts.FrameDuration)
	defer gameTicker.Stop()
	r.game.Start()
	for {
		select {
		case <-gameTicker.C:
			r.gameLock.Lock()
			r.game.Update()
			r.gameLock.Unlock()

			r.gameLock.RLock()
			data := r.game.JSON()
			r.broadcast(data)
			r.gameLock.RUnlock()
		case <-ctx.Done():
			return
		}
	}
}

func (r *Room) broadcast(data []byte) {
	for _, p := range r.players {
		p.Output <- data
	}
}

func (r *Room) IsOpen() bool {
	return len(r.players) < 2
}
