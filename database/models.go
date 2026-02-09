package database

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client
type Client struct {
	ID       string
	Conn     *websocket.Conn
	RoomID   string
	Send     chan []byte
	Username string
}

func NewClient(id string, conn *websocket.Conn, username string) *Client {
	return &Client{
		ID:       id,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Username: username,
	}
}

// Room represents a chat/signaling room
type Room struct {
	ID      string
	Clients map[string]*Client
	Mutex   sync.RWMutex
}

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		Clients: make(map[string]*Client),
	}
}

func (r *Room) AddClient(client *Client) bool {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if len(r.Clients) >= 2 {
		return false
	}

	r.Clients[client.ID] = client
	client.RoomID = r.ID
	return true
}

func (r *Room) RemoveClient(clientID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if client, exists := r.Clients[clientID]; exists {
		close(client.Send)
		delete(r.Clients, clientID)
	}
}

func (r *Room) GetOtherClient(clientID string) *Client {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	for id, client := range r.Clients {
		if id != clientID {
			return client
		}
	}
	return nil
}

func (r *Room) IsFull() bool {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()
	return len(r.Clients) >= 2
}

func (r *Room) IsEmpty() bool {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()
	return len(r.Clients) == 0
}

func (r *Room) GetClientCount() int {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()
	return len(r.Clients)
}

// MessageType defines the type of WebSocket message
type MessageType string

const (
	MessageTypeOffer     MessageType = "offer"
	MessageTypeAnswer    MessageType = "answer"
	MessageTypeCandidate MessageType = "candidate"
	MessageTypeJoin      MessageType = "join"
	MessageTypeLeave     MessageType = "leave"
	MessageTypeReady     MessageType = "ready"
	MessageTypeError     MessageType = "error"
)
