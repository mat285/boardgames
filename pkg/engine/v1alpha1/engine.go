package v1alpha1

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
	persist "github.com/mat285/boardgames/pkg/persist/v1alpha1"
)

type Engine struct {
	sync.Mutex
	ID uuid.UUID

	started bool

	State *v1alpha1.State
	Game  v1alpha1.Game

	Persist persist.Interface

	interrupt chan Event
	stop      chan struct{}
}

func NewEngine(state *v1alpha1.State, game v1alpha1.Game) *Engine {
	return &Engine{
		ID:        uuid.V4(),
		State:     state,
		Game:      game,
		interrupt: make(chan Event),
		stop:      make(chan struct{}),
	}
}

func (e *Engine) Join(player v1alpha1.Player) error {
	if e.started {
		return fmt.Errorf("Game Already Started")
	}
	e.Lock()
	defer e.Unlock()
	for i, p := range e.State.Players {
		if p.ID.Equal(player.ID) {
			e.State.Players[i] = player
			return nil
		}
	}
	e.State.Players = append(e.State.Players, player)
	return nil
}

func (e *Engine) Start(ctx context.Context) error {
	e.Lock()
	data, err := e.Game.Initialize(e.PlayerIDs())
	if err != nil {
		e.Unlock()
		return err
	}
	e.State.Data = data
	e.started = true
	e.Unlock()
	return e.gameLoop(ctx)
}

func (e *Engine) gameLoop(ctx context.Context) error {
	if e.stop == nil {
		e.stop = make(chan struct{})
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-e.stop:
			return nil
		case inter := <-e.interrupt:
			err := e.handleInterrupt(ctx, inter)
			if err != nil {
				return err
			}
		default:
		}

		if e.State.Data.IsDone() {
			msg, err := v1alpha1.MessageGameOver(e.State.Data.Winners())
			if err != nil {
				fmt.Println(err)
			}
			return e.Broadcast(ctx, msg)
		}

		pid, err := e.State.Data.CurrentPlayer()
		if err != nil {
			fmt.Println(err)
			continue
		}

		player := e.State.GetPlayer(pid)
		if player == nil {
			fmt.Println("No player for id", pid)
			continue
		}

		move, err := player.Connection.Request(ctx, v1alpha1.MoveRequest{State: e.State.Data})
		if err != nil {
			fmt.Println(err)
			continue
		}

		response, err := move.Apply(e.State.Data)
		if err != nil {
			fmt.Println(err)
			player.Connection.Accept(ctx, v1alpha1.MessageFromError(err))
			continue
		}

		if !response.Valid {
			fmt.Println("Invalid Move")
			continue
		}

		msg, err := v1alpha1.MessagePlayerMove(player.ID, move)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = e.Broadcast(ctx, msg, player.ID)
		if err != nil {
			fmt.Println(err)
			continue
		}

		e.State.Data = response.State
	}
}

func (e *Engine) handleInterrupt(ctx context.Context, event Event) error {
	switch event.Type {
	case EventTypeUnknown:
		fmt.Printf("Unknown interrupt %v\n", event.Body)
		return nil
	case EventTypeStop:
		return e.Stop(ctx, event.Body)
	case EventTypeSave:
		err := e.Save(ctx)
		if err != nil {
			fmt.Println(err)
		}
		return e.Stop(ctx, event.Body)
	}
	return nil
}

func (e *Engine) Stop(ctx context.Context, optional ...interface{}) error {
	close(e.stop)
	var body []byte
	if len(optional) > 0 && optional[0] != nil {
		body, _ = json.Marshal(optional[0])
	}
	msg := v1alpha1.Message{
		Type: v1alpha1.MessageTypeGameStopped,
		Data: body,
	}
	e.Broadcast(ctx, msg)
	return nil
}

func (e *Engine) Save(ctx context.Context) error {
	if e.Persist == nil {
		return nil
	}
	obj := persist.Object{
		Meta: persist.Meta{
			ID:            e.ID,
			APIVersion:    APIVersion,
			ObjectVersion: e.State.Version,
		},
		Data: e.State,
	}
	_, err := e.Persist.CheckAndSet(ctx, obj)
	return err
}

func (e *Engine) Broadcast(ctx context.Context, message v1alpha1.Message, exclude ...uuid.UUID) error {
	for _, player := range e.State.Players {
		if excludeUUID(player.ID, exclude...) {
			continue
		}
		err := player.Connection.Accept(ctx, message)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (e *Engine) PlayerIDs() []uuid.UUID {
	ids := make([]uuid.UUID, len(e.State.Players))
	for i := range e.State.Players {
		ids[i] = e.State.Players[i].ID
	}
	return ids
}
