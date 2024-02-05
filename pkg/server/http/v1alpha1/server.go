package v1alpha1

import (
	"context"

	"github.com/blend/go-sdk/web"
	core "github.com/mat285/boardgames/pkg/server/core/v1alpha1"
)

type Server struct {
	Config web.Config
	App    *web.App

	Router *core.EngineRouter
	Polls  map[string]*PollClient
}

func New(config web.Config) *Server {
	s := &Server{
		Config: config,
		Router: core.NewEngineRouter(),
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	app, err := web.New(web.OptConfig(s.Config))
	if err != nil {
		return err
	}
	s.App = app

	s.App.Register(s)
	return s.App.Start()
}
