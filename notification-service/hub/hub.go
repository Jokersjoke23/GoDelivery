package hub

import "sync"

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	notify     chan *Message
	mu         sync.RWMutex
}

type Message struct {
	UserID  string
	Payload []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		notify:     make(chan *Message, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()

		case msg := <-h.notify:
			h.mu.RLock()
			if client, ok := h.clients[msg.UserID]; ok {
				select {
				case client.Send <- msg.Payload:
				default:
					close(client.Send)
					delete(h.clients, msg.UserID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Notify(userID string, payload []byte) {
	h.notify <- &Message{UserID: userID, Payload: payload}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}
