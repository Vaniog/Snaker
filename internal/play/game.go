package play

import (
	"context"
	"github.com/Vaniog/Snaker/internal/game"
	"github.com/Vaniog/Snaker/internal/play/event"
	"time"
)

const typeRotate event.Type = "rotate"

type rotate struct {
	Type event.Type     `json:"event"`
	Drc  game.Direction `json:"direction"`
}

type PlayerEvent struct {
	player *Player
	bytes  []byte
}

type Game struct {
	players map[*Player]*game.Snake
	game    *game.Game

	events chan PlayerEvent
}

func newGame(lobby *Lobby) *Game {
	g := &Game{
		events:  lobby.events,
		players: make(map[*Player]*game.Snake, len(lobby.players)),
		game:    game.NewGame(lobby.opts),
	}
	for _, p := range lobby.players {
		g.players[p] = g.game.RegisterSnake()
	}
	return g
}

func (g *Game) Run(ctx context.Context) {
	gameTicker := time.NewTicker(g.game.Opts.FrameDuration)
	defer gameTicker.Stop()
	g.game.Start()

	// TODO move to better place
	for p := range g.players {
		go p.inputPump(ctx)
	}

	for {
		for {
			select {
			case <-gameTicker.C:
				g.game.Update()
				data := g.game.JSON()
				g.broadcast(data)
			case ep := <-g.events:
				p := ep.player
				data := ep.bytes

				switch event.ParseType(data) {
				case typeRotate:
					if eRotate, ok := event.Parse[rotate](data); ok {
						g.players[p].Rotate(eRotate.Drc)
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

func (g *Game) broadcast(data []byte) {
	for p := range g.players {
		p.Output <- data
	}
}
