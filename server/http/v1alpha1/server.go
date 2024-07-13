package v1alpha1

import (
	"context"
	"sync"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
	obj "github.com/mat285/boardgames/pkg/core/v1alpha1"
	"github.com/mat285/boardgames/pkg/websockets"
	"github.com/mat285/boardgames/pkg/wire/v1alpha1"
	core "github.com/mat285/boardgames/server/core/v1alpha1"
)

// var (
// 	_ connection.ServerInfo = new(Server)
// )

type Server struct {
	obj.Object
	Ctx    context.Context
	Config Config
	App    *web.App

	Router         *core.EngineRouter
	InboundPackets chan websockets.Packet
	stop           chan struct{}
	Polls          map[string]*PollClient

	usersLock sync.Mutex
	Users     map[string]uuid.UUID
}

func New(ctx context.Context, config Config) *Server {
	s := &Server{
		// Object: obj.NewObject(id),
		Ctx:            ctx,
		Config:         config,
		Router:         core.NewEngineRouter(),
		InboundPackets: make(chan websockets.Packet, 16),
		stop:           make(chan struct{}),
		Users:          make(map[string]uuid.UUID),
	}
	return s
}

func (s *Server) Start() error {
	err := s.Router.ConnectServer(s.Ctx, core.PipeReceiver(s.ID, s))
	if err != nil {
		return err
	}
	app, err := web.New(
		web.OptConfig(s.Config.Web),
		web.OptBaseContext(s.Ctx),
		web.OptLog(logger.GetLogger(s.Ctx)),
	)
	if err != nil {
		return err
	}
	s.App = app

	s.App.Register(s)
	go s.receivePackets()
	return s.App.Start()
}

func (s *Server) receivePackets() {
	log := logger.GetLogger(s.Ctx)
	for {
		select {
		case <-s.Ctx.Done():
			return
		case <-s.stop:
			return
		case p := <-s.InboundPackets:
			err := s.Router.Receive(s.Ctx, v1alpha1.FromWebsocket(p))
			if err != nil {
				logger.MaybeError(log, err)
			}
		}
	}
}

func (s *Server) Stop() error {
	close(s.stop)
	return s.App.Stop()
}
