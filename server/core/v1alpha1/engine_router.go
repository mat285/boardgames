package v1alpha1

import (
	"context"
	"fmt"
	"sync"

	"github.com/blend/go-sdk/uuid"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	engine "github.com/mat285/boardgames/pkg/engine/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
	router "github.com/mat285/boardgames/pkg/router/v1alpha1"
)

type EngineRouter struct {
	connection.Router

	clientEnginesLock sync.Mutex
	clientEngines     map[string]map[string]bool
}

func NewEngineRouter() *EngineRouter {
	s := &EngineRouter{
		Router: router.NewRouter(),

		clientEngines: make(map[string]map[string]bool),
	}
	return s
}

func (r *EngineRouter) GetEngine(id uuid.UUID) *engine.Engine {
	e := r.GetServer(id)
	// fmt.Println("getting engine", id, e)
	if e == nil {
		return nil
	}
	pipe, ok := e.(*Pipe)
	if !ok {
		return nil
	}
	typed, ok := pipe.Receiver.(*engine.Engine)
	if !ok {
		return nil
	}
	return typed
}

func (r *EngineRouter) NewEngine(ctx context.Context, g v1alpha1.Game, host *engine.Player) (*engine.Engine, error) {
	e := engine.NewEngine(g, host)
	pipe := PipeEngine(e)
	err := r.ConnectServer(ctx, pipe)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *EngineRouter) StartEngine(ctx context.Context, id uuid.UUID) error {
	e := r.GetEngine(id)
	if e == nil {
		return fmt.Errorf("No engine to start")
	}
	return e.Start(ctx)
}

func (r *EngineRouter) Join(ctx context.Context, clientID uuid.UUID, engine uuid.UUID) error {
	client := r.GetClient(clientID)
	if client == nil {
		return fmt.Errorf("Unknown client")
	}
	e := r.GetEngine(engine)
	if e == nil {
		return fmt.Errorf("No Engine")
	}
	err := e.Join(ctx, client)
	if err != nil {
		return err
	}
	r.clientEnginesLock.Lock()
	if _, has := r.clientEngines[clientID.ToFullString()]; !has {
		r.clientEngines[clientID.ToFullString()] = make(map[string]bool)
	}
	r.clientEngines[clientID.ToFullString()][engine.ToFullString()] = true
	r.clientEnginesLock.Unlock()

	return nil
}

func (r *EngineRouter) ClientEngines(ctx context.Context, client uuid.UUID) []*engine.Engine {
	r.clientEnginesLock.Lock()
	engines := r.clientEngines[client.ToFullString()]
	r.clientEnginesLock.Unlock()
	ret := make([]*engine.Engine, 0, len(engines))
	for e, has := range engines {
		if !has {
			continue
		}
		id, err := uuid.Parse(e)
		if err != nil {
			continue
		}
		ret = append(ret, r.GetEngine(id))
	}
	return ret
}
