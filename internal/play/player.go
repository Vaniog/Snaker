package play

import (
	"context"
)

type Player struct {
	events chan<- PlayerEvent
	Input  chan []byte
	Output chan []byte
}

func newPlayer(events chan<- PlayerEvent) *Player {
	p := &Player{
		events: events,
		Input:  make(chan []byte),
		Output: make(chan []byte),
	}
	return p
}

func (p *Player) inputPump(ctx context.Context) {
	for {
		select {
		case data, ok := <-p.Input:
			if !ok {
				return
			}
			p.events <- PlayerEvent{p, data}
		case <-ctx.Done():
			return
		}
	}
}

func (p *Player) Disconnect(ctx context.Context) {
	go p.fakeReadPump(ctx)
}

func (p *Player) fakeReadPump(ctx context.Context) {
	for {
		select {
		case <-p.Output:
		case <-ctx.Done():
			return
		}
	}
}
