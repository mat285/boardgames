package v1alpha1

import (
	"context"

	core "github.com/mat285/boardgames/pkg/server/core/v1alpha1"
)

type Server struct {
	core.EngineRouter
}

func NewServer() *Server {
	s := &Server{
		EngineRouter: *core.NewEngineRouter(),
	}
	return s
}

func (s *Server) ConnectClient(ctx context.Context, client *Client) error {
	return s.EngineRouter.ConnectClient(ctx, client.Pipe())
}
