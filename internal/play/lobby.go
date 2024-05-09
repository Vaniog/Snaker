package play

import (
	"context"
	"github.com/Vaniog/Snaker/internal/game"
	"github.com/Vaniog/Snaker/internal/play/event"
	"slices"
)

const typeUpdateOptions event.Type = "update_options"
const typeGameStart event.Type = "game_start"

type updateOptions struct {
	Type event.Type   `json:"event"`
	Opts game.Options `json:"options"`
}

type Lobby struct {
	Opts    game.Options
	players []*Player

	register   chan *Player
	unregister chan *Player
	events     chan PlayerEvent
}

func NewLobby(opts game.Options) *Lobby {
	return &Lobby{
		events:   make(chan PlayerEvent),
		Opts:     opts,
		players:  nil,
		register: make(chan *Player),
	}
}

func (lb *Lobby) RegisterPlayer(ctx context.Context) *Player {
	p := newPlayer(lb.events)
	lb.register <- p
	go p.inputPump(ctx)
	return p
}

func (lb *Lobby) Run(ctx context.Context) {
	for {
		select {
		case p := <-lb.register:
			lb.players = append(lb.players, p)
			// TODO make better
			if len(lb.players) == lb.Opts.SnakesAmount {
				g := newGame(lb)
				go g.Run(ctx)
				return
			}
		case p := <-lb.unregister:
			lb.players = slices.DeleteFunc(lb.players, func(i *Player) bool {
				return p == i
			})
		case ep := <-lb.events:
			data := ep.bytes

			switch event.ParseType(data) {
			case typeUpdateOptions:
				if eUpdateOpts, ok := event.Parse[updateOptions](data); ok {
					lb.Opts = eUpdateOpts.Opts
				}
			case typeGameStart:
				g := newGame(lb)
				go g.Run(ctx)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (lb *Lobby) IsOpen() bool {
	return len(lb.players) != lb.Opts.SnakesAmount
}
