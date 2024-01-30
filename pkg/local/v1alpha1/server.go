package v1alpha1

import (
	"context"
	"fmt"
	"sync"

	"github.com/blend/go-sdk/uuid"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	engine "github.com/mat285/boardgames/pkg/engine/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

var (
	_ connection.Server = new(Server)
)

type Server struct {
	sync.Mutex
	clients map[string]*connection.ClientInfo
	engines map[string]*engine.Engine
}

func NewServer() *Server {
	s := &Server{
		clients: make(map[string]*connection.ClientInfo),
		engines: make(map[string]*engine.Engine),
	}
	return s
}

func (s *Server) Connect(ctx context.Context, client connection.ClientInfo) (uuid.UUID, error) {
	id := uuid.V4()
	client.ID = id
	s.Lock()
	defer s.Unlock()
	s.clients[client.ID.ToFullString()] = &client
	return id, nil
}

func (s *Server) Receive(ctx context.Context, packet wire.Packet) error {
	origin := packet.Origin
	dst := packet.Destination
	if s.GetClient(origin) != nil {
		e := s.GetEngine(dst)
		if e == nil {
			return fmt.Errorf("No engine")
		}
		return e.Receive(ctx, packet)
	} else if s.GetEngine(origin) != nil {
		client := s.GetClient(dst)
		if client == nil {
			return fmt.Errorf("No Client")
		}
		return client.Sender.Send(ctx, packet)
	} else {
		return fmt.Errorf("Unknown origin")
	}
}

func (s *Server) GetEngine(id uuid.UUID) *engine.Engine {
	s.Lock()
	defer s.Unlock()
	return s.engines[id.ToFullString()]
}

func (s *Server) GetClient(id uuid.UUID) *connection.ClientInfo {
	s.Lock()
	defer s.Unlock()
	return s.clients[id.ToFullString()]
}

func (s *Server) NewEngine(g v1alpha1.Game) *engine.Engine {
	s.Lock()
	defer s.Unlock()
	e := engine.NewEngine(g)
	s.engines[e.ID.ToFullString()] = e
	return e
}

func (s *Server) StartEngine(ctx context.Context, id uuid.UUID) error {
	e := s.GetEngine(id)
	if e == nil {
		return fmt.Errorf("No engine to start")
	}
	return e.Start(ctx)
}

func (s *Server) Join(ctx context.Context, clientID uuid.UUID, engine uuid.UUID) error {
	client := s.GetClient(clientID)
	if client == nil {
		return fmt.Errorf("Unknown client")
	}
	e := s.GetEngine(engine)
	if e == nil {
		return fmt.Errorf("No Engine")
	}
	return e.Join(ctx, *client)
}
