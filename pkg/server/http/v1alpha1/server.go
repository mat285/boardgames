package v1alpha1

import (
	"context"

	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
	obj "github.com/mat285/boardgames/pkg/core/v1alpha1"
	core "github.com/mat285/boardgames/pkg/server/core/v1alpha1"
)

// var (
// 	_ connection.ServerInfo = new(Server)
// )

type Server struct {
	obj.Object
	Config web.Config
	App    *web.App

	Router *core.EngineRouter
	Polls  map[string]*PollClient
}

func New(id uuid.UUID, config web.Config) *Server {
	s := &Server{
		Object: obj.NewObject(id),
		Config: config,
		Router: core.NewEngineRouter(),
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	err := s.Router.ConnectServer(ctx, core.PipeReceiver(s.ID, s))
	if err != nil {
		return err
	}
	app, err := web.New(web.OptConfig(s.Config))
	if err != nil {
		return err
	}
	s.App = app

	s.App.Register(s)
	return s.App.Start()
}
