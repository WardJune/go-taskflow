package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

type Notification struct {
	Type      string `json:"type"`
	ProjectID int64  `json:"project_id"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
}

type Client struct {
	UserID    int64
	ProjectID int64
	Send      chan []byte
	hub       *Hub
	conn      *Conn
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan BroadCastMsg
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type BroadCastMsg struct {
	ProjectID    int64
	Notification Notification
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan BroadCastMsg),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client Connected: userID=%d projectID=%d", client.UserID, client.ProjectID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected: userID=%d", client.UserID)

		case msg := <-h.broadcast:
			data, err := json.Marshal(msg.Notification)
			if err != nil {
				continue
			}

			h.mu.RLock()
			for client := range h.clients {
				if client.ProjectID == msg.ProjectID {
					select {
					case client.Send <- data:
					default:
						close(client.Send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastToProject(projectID int64, notif Notification) {
	h.broadcast <- BroadCastMsg{
		ProjectID:    projectID,
		Notification: notif,
	}
}
