package v1alpha1

// import (
// 	"context"
// 	"fmt"
// 	"sync"

// 	"github.com/blend/go-sdk/uuid"
// 	"github.com/mat285/boardgames/pkg/game/v1alpha1"
// 	"github.com/mat285/boardgames/pkg/persist"
// )

// type Engine struct {
// 	ID uuid.UUID

// 	Connections *ConnectionManager

// 	State *v1alpha1.State
// 	Game  v1alpha1.Game

// 	Persist persist.Interface

// 	interrupt chan struct{}

// 	runningLock sync.Mutex
// 	running     bool
// }

// func New(game v1alpha1.Game, players []v1alpha1.Player) *Engine {
// 	e := &Engine{
// 		State: &v1alpha1.State{
// 			Players: players,
// 		},
// 		Game:      game,
// 		interrupt: make(chan struct{}),
// 	}
// 	e.Connections = NewConnectionManager(e)
// 	return e
// }

// func (e *Engine) Players() []v1alpha1.Player {
// 	if e.State == nil {
// 		return nil
// 	}
// 	return e.State.Players
// }

// func (e *Engine) PlayerIDs() []uuid.UUID {
// 	ids := make([]uuid.UUID, len(e.State.Players))
// 	for i := range e.State.Players {
// 		ids[i] = e.State.Players[i].ID
// 	}
// 	return ids
// }

// func (e *Engine) Initialize() error {
// 	if e.State != nil {
// 		return fmt.Errorf("Game already initialized")
// 	}

// 	data, err := e.Game.Initialize(e.PlayerIDs())
// 	if err != nil {
// 		return err
// 	}

// 	e.State = &v1alpha1.State{
// 		Version: 1,
// 		Data:    data,
// 	}
// 	return nil
// }

// func (e *Engine) Start(ctx context.Context) error {
// 	go e.Connections.Start(ctx)
// 	err := e.Connections.Ready(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	// all players connected lets play
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		case <-e.interrupt:
// 			return e.pause(ctx)
// 		default:
// 			// fall through
// 		}

// 		err = e.step(ctx)
// 		if err != nil {
// 			// TODO retry
// 		}
// 	}
// }

// func (e *Engine) step(ctx context.Context) error {
// 	gameState := e.State.Data
// 	p := gameState.CurrentPlayer()

// 	move, err := e.Connections.RequestMove(ctx, &p, v1alpha1.MoveRequest{gameState})
// 	if err != nil {
// 		return err
// 	}

// 	if !move.IsValid(gameState) {
// 		e.State.Attempts++
// 		if e.State.Attempts > 5 {
// 			// kick player, end game just spamming
// 		}
// 		return nil // retry
// 	}

// 	newState, err := move.Apply(gameState)
// 	if err != nil {
// 		return err
// 	}

// 	e.persist(ctx, newState)
// 	if err != nil {
// 		// retry with new reconciled state if need be
// 		// broadcast current state
// 	}
// 	return nil
// }

// func (e *Engine) persist(ctx context.Context, data v1alpha1.StateData) error {
// 	state := v1alpha1.State{
// 		Players:  e.State.Players,
// 		Version:  e.State.Version + 1,
// 		Attempts: 0,
// 		Data:     data,
// 	}
// 	// send upstream and get back current version

// 	obj := persist.Object{
// 		Meta: persist.Meta{
// 			ID:            e.ID,
// 			APIVersion:    "",
// 			ObjectVersion: state.Version,
// 		},
// 		Data: state,
// 	}

// 	curr, err := e.Persist.CheckAndSet(ctx, obj)
// 	if err != nil {
// 		return err
// 	}

// 	currState, ok := curr.Data.(*v1alpha1.State)
// 	if !ok {
// 		return fmt.Errorf("Bad game state")
// 	}
// 	e.State = currState
// 	e.State.Version = curr.ObjectVersion
// 	return nil
// }

// func (e *Engine) pause(ctx context.Context) error {
// 	// save game state
// 	// alert all players to rejoin
// 	return nil
// }
