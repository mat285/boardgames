package v1alpha1

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type Engine struct {
	State *v1alpha1.State
	Game  v1alpha1.Game
}

func NewEngine(state *v1alpha1.State, game v1alpha1.Game) *Engine {
	return &Engine{
		State: state,
		Game:  game,
	}
}

func (e *Engine) Start(ctx context.Context) error {
	data, err := e.Game.Initialize(e.PlayerIDs())
	if err != nil {
		return err
	}
	e.State.Data = data

	return e.gameLoop(ctx)
}

func (e *Engine) gameLoop(ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
		default:
		}

		if e.State.Data.IsDone() {
			// todo winners
			fmt.Println(e.State.Data.Winners())
			return nil
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

func excludeUUID(id uuid.UUID, exclude ...uuid.UUID) bool {
	for _, e := range exclude {
		if id.Equal(e) {
			return true
		}
	}
	return false
}
