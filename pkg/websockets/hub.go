package websockets

import (
	"context"
	"fmt"
	"sync"

	"github.com/blend/go-sdk/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	Upgrader websocket.Upgrader

	lock    sync.Mutex
	clients map[string]*Client

	register   chan *Client
	unregister chan *Client
}

func NewHub() (*Hub, error) {
	return &Hub{}, nil
}

func (h *Hub) Register(client *Client) error {
	if client.ID.IsZero() {
		return fmt.Errorf("Client ID required")
	}
	h.lock.Lock()
	h.clients[client.ID.ToFullString()] = client
	h.lock.Unlock()
	return nil
}

func (h *Hub) Unregister(client *Client) error {
	if client.ID.IsZero() {
		return fmt.Errorf("Client ID required")
	}
	h.lock.Lock()
	delete(h.clients, client.ID.ToFullString())
	h.lock.Unlock()
	return nil
}

func (h *Hub) Send(ctx context.Context, id uuid.UUID, packet *Packet) error {
	client := h.clients[id.ToFullString()]
	if client == nil {
		// drop packet
		return nil
	}
	return client.pushOutbound(packet)
}
