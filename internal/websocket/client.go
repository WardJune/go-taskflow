package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Conn struct {
	ws *websocket.Conn
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.ws.Close()
		c.hub.unregister <- c
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.conn.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.ws.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.ws.Close()
	}()

	c.conn.ws.SetReadLimit(maxMessageSize)
	c.conn.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.ws.SetPongHandler(func(string) error {
		c.conn.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.conn.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}
	}
}

func ServeWS(hub *Hub, userID, projectID int64, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade error:%v", err)
		return
	}

	client := &Client{
		UserID:    userID,
		ProjectID: projectID,
		Send:      make(chan []byte, 256),
		hub:       hub,
		conn:      &Conn{ws: ws},
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}
