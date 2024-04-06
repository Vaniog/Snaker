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
	const infTime = 1000000 * time.Hour
	gameTicker := time.NewTicker(infTime)
	defer gameTicker.Stop()

	for {
		select {
		case p := <-r.Register:
			r.players = append(r.players, p)
			snake := r.game.RegisterSnake()

			p.snake = snake
			go p.inputPump(ctx)
			if len(r.players) == 2 {
				gameTicker.Reset(r.game.Opts.FrameDuration)
				r.game.Start()
			}
		case <-gameTicker.C:
			r.gameLock.Lock()
			r.game.Update()
			r.gameLock.Unlock()

			r.gameLock.RLock()
			data := r.game.JSON()
			for _, p := range r.players {
				p.broadcast <- data
			}
			r.gameLock.RUnlock()
		case <-ctx.Done():
			return
		}
	}
}

func (r *Room) IsOpen() bool {
	return len(r.players) < 2
}