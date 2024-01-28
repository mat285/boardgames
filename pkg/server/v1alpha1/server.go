package v1alpha1

import (
	"context"
	"sync"

	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
	engine "github.com/mat285/boardgames/pkg/engine/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type Server struct {
	Config web.Config
	App    *web.App

	enginesLock sync.Mutex
	engines     map[string]*engine.Engine
}

func New(config web.Config) *Server {
	s := &Server{
		Config: config,

		enginesLock: sync.Mutex{},
		engines:     make(map[string]*engine.Engine),
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

func (s *Server) GetEngine(id uuid.UUID) *engine.Engine {
	return s.engines[id.ToFullString()]
}

func (s *Server) NewEngine(g v1alpha1.Game) *engine.Engine {
	s.enginesLock.Lock()
	defer s.enginesLock.Unlock()
	e := engine.NewEngine(v1alpha1.NewState([]v1alpha1.Player{}), g)
	s.engines[e.ID.ToFullString()] = e
	return e
}
