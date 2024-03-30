package server

import (
	"context"
	"encoding/json"
	"github.com/Vaniog/Snaker/internal/game"
	"github.com/Vaniog/Snaker/internal/server/event"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Player struct {
	hub  *Hub
	conn *websocket.Conn

	snake     *game.Snake
	broadcast chan []byte
}

func newPlayer(hub *Hub, conn *websocket.Conn) *Player {
	return &Player{
		hub:       hub,
		conn:      conn,
		snake:     nil,
		broadcast: make(chan []byte),
	}
}

// readPump handles pongs and rotations
func (p *Player) readPump(ctx context.Context) {
	defer func() {
		_ = p.conn.Close()
	}()

	p.conn.SetReadLimit(maxMessageSize)

	initPongs(p.conn)

	for {
		if ctx.Err() != nil {
			return
		}

		_, data, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}
		var e struct {
			Type event.Type `json:"event"`
		}

		err = json.Unmarshal(data, &e)
		if err != nil {
			return
		}

		switch e.Type {
		case event.TypeRotate:
			var eRotate event.Rotate
			err := json.Unmarshal(data, &eRotate)
			if err != nil {
				return
			}
			p.handleEventRotate(eRotate)
			continue
		default:
			log.Printf("illegal incoming data: %s", data)
			return
		}
	}
}

func (p *Player) handleEventRotate(e event.Rotate) {
	p.hub.gameLock.Lock()
	defer p.hub.gameLock.Unlock()

	p.snake.Rotate(e.Drc)
}

func initPongs(conn *websocket.Conn) {
	err := conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Printf("can't SetReadDeadline: %s", err)
		return
	}

	conn.SetPongHandler(func(string) error {
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Printf("can't SetReadDeadline: %s", err)
			return err
		}
		return nil
	})
}

// writePump handles pings and broadcasts messages from Player.broadcast
func (p *Player) writePump(ctx context.Context) {
	pongTicker := time.NewTicker(pingPeriod)
	defer func() {
		_ = p.conn.WriteMessage(websocket.CloseMessage, []byte{})
		pongTicker.Stop()
		_ = p.conn.Close()
	}()

	for {
		select {
		case message, ok := <-p.broadcast:
			_ = p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				return
			}

			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err = w.Write(message); err != nil {
				return
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-pongTicker.C:
			_ = p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
