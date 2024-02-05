package v1alpha1

import (
	"context"
	"fmt"
	"sync"

	"github.com/blend/go-sdk/uuid"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

var (
	_ connection.Router = new(Router)
)

type Router struct {
	sync.Mutex
	clients map[string]connection.ClientInfo
	servers map[string]connection.ServerInfo
}

func NewRouter() *Router {
	s := &Router{
		clients: make(map[string]connection.ClientInfo),
		servers: make(map[string]connection.ServerInfo),
	}
	return s
}

func (s *Router) ConnectClient(ctx context.Context, client connection.ClientInfo) error {
	s.Lock()
	defer s.Unlock()
	s.clients[client.GetID().ToFullString()] = client
	return nil
}

func (s *Router) ConnectServer(ctx context.Context, server connection.ServerInfo) error {
	s.Lock()
	defer s.Unlock()
	s.servers[server.GetID().ToFullString()] = server
	return nil
}

func (s *Router) Receive(ctx context.Context, packet wire.Packet) error {
	if s.GetClient(packet.Origin) != nil {
		s := s.GetServer(packet.Destination)
		if s == nil {
			return fmt.Errorf("No server")
		}
		return s.Send(ctx, packet)
	} else if s.GetServer(packet.Origin) != nil {
		client := s.GetClient(packet.Destination)
		server := s.GetServer(packet.Destination)
		if client != nil {
			return client.Send(ctx, packet)
		} else if server != nil {
			return server.Send(ctx, packet)
		} else {
			return fmt.Errorf("Destination unreachable")
		}
	} else {
		return fmt.Errorf("Unknown origin")
	}
}

func (s *Router) GetServer(id uuid.UUID) connection.ServerInfo {
	s.Lock()
	defer s.Unlock()
	return s.servers[id.ToFullString()]
}

func (s *Router) GetClient(id uuid.UUID) connection.ClientInfo {
	s.Lock()
	defer s.Unlock()
	return s.clients[id.ToFullString()]
}
