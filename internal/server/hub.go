package server

import (
	"context"
	"github.com/Vaniog/Snaker/internal/play"
	"github.com/gorilla/websocket"
	"log"
)

type HubID int64
type Hub struct {
	ID    HubID
	lobby *play.Lobby

	register chan *Client
}

func newHub() *Hub {
	return &Hub{
		lobby:    play.NewLobby(),
		register: make(chan *Client),
	}
}

func (h *Hub) RegisterClient(conn *websocket.Conn) *Client {
	c := &Client{conn: conn}
	h.register <- c
	return c
}

func (h *Hub) Run(ctx context.Context) {
	log.Printf("Start new hub: %d", h.ID)
	go h.lobby.Run(ctx)
	for {
		select {
		case c := <-h.register:
			c.player = h.lobby.RegisterPlayer(ctx)
			go c.readPump(ctx)
			go c.writePump(ctx)
		case <-ctx.Done():
			return
		}
	}
}
