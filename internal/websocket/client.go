package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Lec7ral/WithWebSocket/internal/domain"
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	ID       string
	Username string
	RoomID   string

	conn    *websocket.Conn
	send    chan []byte
	limiter *rate.Limiter
}

// readPump pumps messages from the WebSocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close() // Close the connection on exit.
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		if err := c.limiter.Wait(context.Background()); err != nil {
			slog.Error("Rate limiter wait error", "error", err, "clientID", c.ID)
			break
		}

		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Warn("Unexpected WebSocket close error", "error", err, "clientID", c.ID)
			}
			break
		}

		var msg domain.Message
		if err := json.Unmarshal(rawMessage, &msg); err == nil {
			switch msg.Type {
			case "direct_message":
				var dmPayload domain.DirectMessagePayload
				payloadBytes, _ := json.Marshal(msg.Payload)
				if err := json.Unmarshal(payloadBytes, &dmPayload); err == nil {
					msg.Sender = c.ID
					msg.Payload = dmPayload
					c.hub.broadcast <- &msg
				}
			case "draw_start", "draw_move", "draw_end", "clear_board", "typing_start", "typing_stop":
				msg.Sender = c.ID
				msg.RoomID = c.RoomID
				c.hub.broadcast <- &msg
			default:
				// If the type is unknown but it's valid JSON, we assume it's a text message.
				// This handles the case where the client sends `{"type":"text_message", "payload":"..."}`
				if textPayload, ok := msg.Payload.(string); ok {
					c.sendRoomMessage([]byte(textPayload))
				}
			}
		} else {
			// If it's not valid JSON, treat as a plain text message for the room.
			c.sendRoomMessage(rawMessage)
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return // Exit on write error
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return // Exit on ping error
			}
		}
	}
}

// sendRoomMessage is a helper to create and send a standard text message to the hub.
func (c *Client) sendRoomMessage(rawMessage []byte) {
	roomMsg := &domain.Message{
		Type:    "text_message",
		Payload: string(rawMessage),
		Sender:  c.ID,
		RoomID:  c.RoomID,
	}
	c.hub.broadcast <- roomMsg
}
