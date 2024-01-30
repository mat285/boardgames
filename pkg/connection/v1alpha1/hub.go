package v1alpha1

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"time"

// 	"github.com/blend/go-sdk/uuid"
// 	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
// )

// var (
// 	_ Listener = new(Hub)
// )

// type Hub struct {
// 	sync.Mutex
// 	clients map[string]Interface
// 	servers map[string]Server

// 	inbound chan wire.Packet
// }

// func NewHub() *Hub {
// 	h := &Hub{
// 		inbound: make(chan wire.Packet),
// 	}
// 	h.Lock()
// 	defer h.Unlock()

// 	return h
// }

// func (h *Hub) ConnectToServer(clientID uuid.UUID, client Interface, server uuid.UUID) (Interface, error) {
// 	h.Lock()
// 	defer h.Unlock()
// 	return h.servers[server.ToFullString()], nil
// }

// func (h *Hub) ClientConnect(id uuid.UUID, client Interface) error {
// 	h.Lock()
// 	defer h.Unlock()
// 	h.clients[id.ToFullString()] = client
// 	return nil
// }

// func (h *Hub) ServerConnect(id uuid.UUID, server Interface) error {
// 	h.Lock()
// 	defer h.Unlock()
// 	h.servers[id.ToFullString()] = server
// 	return nil
// }

// // func (h *Hub) Connect()

// func (h *Hub) Send(ctx context.Context, packet wire.Packet) error {
// 	timeout := time.After(time.Second * 30)
// 	select {
// 	case <-ctx.Done():
// 		return ctx.Err()
// 	case <-timeout:
// 		return fmt.Errorf("timeout")
// 	case h.inbound <- packet:
// 		return nil
// 	}
// }

// func (h *Hub) Listen(ctx context.Context, handler PacketHandler) error {
// 	return h.listen(ctx, handler)
// }

// func (h *Hub) listen(ctx context.Context, handler PacketHandler) error {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		case packet := <-h.inbound:
// 			err := handler(ctx, packet)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 		}
// 	}
// }
