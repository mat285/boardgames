package v1alpha1

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/blend/go-sdk/uuid"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
	persist "github.com/mat285/boardgames/pkg/persist/v1alpha1"
)

type Engine struct {
	sync.Mutex
	ID uuid.UUID

	started bool

	Players []Player
	State   *game.State
	Game    game.Game

	Persist persist.Interface

	Server *Connection

	interrupt chan Event
	stop      chan struct{}
}

func NewEngine(players []Player, g game.Game) *Engine {
	e := &Engine{
		ID:        uuid.V4(),
		Players:   players,
		Game:      g,
		interrupt: make(chan Event),
		stop:      make(chan struct{}),
	}
	e.State = game.NewState(e.GamePlayers())
	return e
}

func (e *Engine) Join(ctx context.Context, player Player) error {
	if e.started {
		return fmt.Errorf("Game Already Started")
	}
	e.Lock()
	defer e.Unlock()
	for i, p := range e.Players {
		if p.GetID().Equal(player.GetID()) {
			e.Players[i] = player
			return nil
		}
	}
	e.Players = append(e.Players, player)
	e.State.Players = append(e.State.Players, player)
	return nil
}

func (e *Engine) ConnectPlayers(ctx context.Context) error {
	// for _, player := range e.Players {
	// 	err := player.Connect(ctx)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
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
	err = e.ConnectPlayers(ctx)
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
			msg, err := game.MessageGameOver(e.State.Data.Winners())
			if err != nil {
				return err
			}
			return e.Broadcast(ctx, *msg)
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

		msg, err := game.MessageRequestMove(e.State.Data, e.Game.Serializer())
		if err != nil {
			fmt.Println(err)
			continue
		}

		resp, err := player.GetConnection().Request(ctx, *msg)
		if err != nil {
			fmt.Println(err)
			continue
		}
		move, err := game.MoveFromMessage(*resp, e.Game.Serializer())
		if err != nil {
			fmt.Println(err)
			continue
		}

		response, err := move.Apply(e.State.Data)
		if err != nil {
			fmt.Println(err)
			player.GetConnection().Send(ctx, game.MessageFromError(err))
			continue
		}

		if !response.Valid {
			fmt.Println("Invalid Move")
			continue
		}

		msg, err = game.MessagePlayerMoveInfo(player.GetID(), move)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = e.Broadcast(ctx, *msg, player.GetID())
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
	msg := game.Message{
		Type: game.MessageTypeGameStopped,
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

func (e *Engine) Broadcast(ctx context.Context, message game.Message, exclude ...uuid.UUID) error {
	for _, player := range e.State.Players {
		if excludeUUID(player.GetID(), exclude...) {
			continue
		}
		err := player.GetConnection().Send(ctx, message)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (e *Engine) PlayerIDs() []uuid.UUID {
	ids := make([]uuid.UUID, len(e.State.Players))
	for i := range e.Players {
		ids[i] = e.Players[i].GetID()
	}
	return ids
}

func (e *Engine) GamePlayers() []game.Player {
	players := make([]game.Player, len(e.Players))
	for i := range e.Players {
		players[i] = e.Players[i]
	}
	return players
}
