package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

type Hub struct {
	mu         sync.RWMutex
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan RoomMessage
}

type RoomMessage struct {
	Room    string `json:"room"`
	Payload []byte `json:"payload"`
}

func NewHub() *Hub {
	h := &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan RoomMessage, 256),
	}
	go h.run()
	return h
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.rooms[client.Room]; !ok {
				h.rooms[client.Room] = make(map[*Client]bool)
			}
			h.rooms[client.Room][client] = true
			h.mu.Unlock()
			log.Printf("[WebSocket] Client joined room [%s], total clients in room: %d\n", client.Room, len(h.rooms[client.Room]))

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.Room]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.rooms, client.Room)
					}
					log.Printf("[WebSocket] Client left room [%s]\n", client.Room)
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			if clients, ok := h.rooms[msg.Room]; ok {
				for client := range clients {
					select {
					case client.Send <- msg.Payload:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToRoom sends a message to all connected clients in a specific room
func (h *Hub) BroadcastToRoom(room string, payload any) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	h.broadcast <- RoomMessage{
		Room:    room,
		Payload: bytes,
	}
	return nil
}

// BroadcastRaw sends raw bytes to a room
func (h *Hub) BroadcastRaw(room string, payload []byte) {
	h.broadcast <- RoomMessage{
		Room:    room,
		Payload: payload,
	}
}
