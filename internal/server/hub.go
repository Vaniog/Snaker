package server

import (
	"context"
	"github.com/Vaniog/Snaker/internal/game"
	"github.com/Vaniog/Snaker/internal/play"
	"github.com/gorilla/websocket"
	"log"
)

type HubID int64
type Hub struct {
	ID   HubID
	Room *play.Room

	Register chan *Client
}

func newHub() *Hub {
	return &Hub{
		Room:     play.NewRoom(game.NewGame(game.DefaultOptions)),
		Register: make(chan *Client),
	}
}

func (h *Hub) RegisterClient(conn *websocket.Conn) *Client {
	c := &Client{conn: conn}
	h.Register <- c
	return c
}

func (h *Hub) Run(ctx context.Context) {
	log.Printf("StartTime new hub: %d", h.ID)
	go h.Room.Run(ctx)
	for {
		select {
		case c := <-h.Register:
			c.player = play.RegisterPlayer(h.Room)
			go c.readPump(ctx)
			go c.writePump(ctx)
		}
	}
}
