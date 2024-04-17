package play

import (
	"context"
	"github.com/Vaniog/Snaker/internal/game"
	"github.com/Vaniog/Snaker/internal/play/event"
	"time"
)

type Bot struct {
	*Player
	game  *Game
	snake *game.Snake
}

func (b *Bot) Run(ctx context.Context) {
	ticker := time.NewTicker(b.game.Opts.FrameDuration / 10)
	go b.fakeReadPump(ctx)
	for {
		select {
		case <-ticker.C:
			b.Update()
		case <-ctx.Done():
			return
		}
	}
}

func (b *Bot) Update() {
	if len(b.snake.Body) < 3 {
		return
	}
	goodDrc := b.snake.Drc
	ds := map[game.Direction]int{
		game.Up:    0,
		game.Right: 0,
		game.Down:  0,
		game.Left:  0,
	}

	for drc := range ds {
		goodDrc = drc
		var newHead = b.game.Field.ToBounds(b.snake.Head().Move(drc))
		good := true
		for _, s := range b.game.Snakes {
			handleLen := len(s.Body)
			if s == b.snake {
				handleLen--
			}
			for _, p := range s.Body[0:handleLen] {
				if p == newHead {
					good = false
				}
			}
		}
		if good {
			break
		}
	}

	b.Player.Input <- event.Bytes(rotate{Type: typeRotate, Drc: goodDrc})
}
