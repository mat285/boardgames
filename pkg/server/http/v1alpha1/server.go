package v1alpha1

import (
	"context"

	"github.com/blend/go-sdk/web"
	router "github.com/mat285/boardgames/pkg/router/v1alpha1"
)

type Server struct {
	Config web.Config
	App    *web.App

	Router *router.Router
	Polls  map[string]*PollClient
}

func New(config web.Config) *Server {
	s := &Server{
		Config: config,
		Router: router.NewRouter(),
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
