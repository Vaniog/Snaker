package server

import (
	"context"
	"github.com/Vaniog/Snaker/internal/play"
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

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	player *play.Player
}

// readPump handles pongs and rotations
func (c *Client) readPump(ctx context.Context) {
	defer func() {
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)

	initPongs(c.conn)

	for {
		if ctx.Err() != nil {
			return
		}

		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}

		c.player.Input <- data
	}
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

// writePump handles pings and broadcasts messages from Client.broadcast
func (c *Client) writePump(ctx context.Context) {
	pongTicker := time.NewTicker(pingPeriod)
	defer func() {
		_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		pongTicker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.player.Output:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
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
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
