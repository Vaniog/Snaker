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
	*game.Game
	players map[*Player]*game.Snake

	events chan PlayerEvent
}

func newGame(lobby *Lobby) *Game {
	g := &Game{
		Game:    game.NewGame(lobby.Opts),
		events:  lobby.events,
		players: make(map[*Player]*game.Snake, len(lobby.players)),
	}
	for _, p := range lobby.players {
		g.players[p] = g.RegisterSnake()
	}
	return g
}

func (g *Game) Run(ctx context.Context) {
	gameTicker := time.NewTicker(g.Opts.FrameDuration)
	defer gameTicker.Stop()
	g.Start()

	for range g.Opts.BotsAmount {
		g.addBot(ctx)
	}

	g.broadcast(event.Bytes(event.Event{Type: typeGameStart}))

	for {
		if g.allDead() {
			return
		}

		select {
		case <-gameTicker.C:
			g.Update()
			data := g.JSON()
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

func (g *Game) allDead() bool {
	for _, s := range g.players {
		if s.Alive {
			return false
		}
	}
	return true
}

func (g *Game) broadcast(data []byte) {
	var toDelete []*Player

	for p := range g.players {
		select {
		case p.Output <- data:
		default:
			toDelete = append(toDelete, p)
		}
	}

	for _, p := range toDelete {
		delete(g.players, p)
	}
}

func (g *Game) addBot(ctx context.Context) {
	bot := Bot{
		Player: newPlayer(g.events),
		game:   g,
		snake:  g.RegisterSnake(),
	}
	g.players[bot.Player] = bot.snake
	go bot.fakeReadPump(ctx)
	go bot.Player.inputPump(ctx)
	go bot.Run(ctx)
}
