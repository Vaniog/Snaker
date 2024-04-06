package play

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Vaniog/Snaker/internal/game"
	"github.com/Vaniog/Snaker/internal/play/event"
	"log"
)

type Player struct {
	room  *Room
	snake *game.Snake

	Input  chan []byte
	Output chan []byte
}

func RegisterPlayer(room *Room) *Player {
	p := &Player{
		Input:  make(chan []byte),
		Output: make(chan []byte),
		room:   room,
	}
	room.Register <- p
	return p
}

func (p *Player) inputPump(ctx context.Context) {
	for {
		select {
		case data := <-p.Input:
			err := p.handleInputData(data)
			if err != nil {
				log.Printf("cant handle Input data: err: {%s} data: {%s}", err, data)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (p *Player) handleInputData(data []byte) error {
	var e struct {
		Type event.Type `json:"event"`
	}

	err := json.Unmarshal(data, &e)
	if err != nil {
		return err
	}

	switch e.Type {
	case event.TypeRotate:
		eRotate, err := event.Parse[event.Rotate](data)
		if err != nil {
			return err
		}
		p.handleRotate(eRotate)
	default:
		return fmt.Errorf("unknown event type: %s", e.Type)
	}
	return nil
}

func (p *Player) handleRotate(e event.Rotate) {
	p.room.gameLock.Lock()
	p.snake.Rotate(e.Drc)
	p.room.gameLock.Unlock()
}
