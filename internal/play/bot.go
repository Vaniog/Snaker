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
	ds := []game.Direction{
		game.Up,
		game.Right,
		game.Down,
		game.Left,
	}

	bestDrc := b.snake.Drc
	bestScore := -1
	for _, drc := range ds {
		score := 0

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
					break
				}
			}
		}

		if !good {
			continue
		}
		if drc == b.snake.Drc {
			score += 1
		}
		if drc == game.ApproxDrcBetween(b.snake.Head(), b.game.Food.Point) {
			score += 2
		}
		if score > bestScore {
			bestDrc = drc
			bestScore = score
		}
	}

	b.Player.Input <- event.Bytes(rotate{Type: typeRotate, Drc: bestDrc})
}
