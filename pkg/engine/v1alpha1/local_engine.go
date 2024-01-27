package v1alpha1

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type E2 struct {
	State *v1alpha1.State
	Game  v1alpha1.Game
}

func NewE2(state *v1alpha1.State, game v1alpha1.Game) *E2 {
	return &E2{
		State: state,
		Game:  game,
	}
}

func (e *E2) Start(ctx context.Context) error {
	data, err := e.Game.Initialize(e.PlayerIDs())
	if err != nil {
		return err
	}
	e.State.Data = data

	return e.gameLoop(ctx)
}

func (e *E2) gameLoop(ctx context.Context) error {

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

		move, err := player.Connection.RequestMove(ctx, v1alpha1.MoveRequest{State: e.State.Data})
		if err != nil {
			fmt.Println(err)
			continue
		}

		response, err := move.Apply(e.State.Data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if !response.Valid {
			fmt.Println("Invalid Move")
			continue
		}

		e.State.Data = response.State
	}
}

func (e *E2) PlayerIDs() []uuid.UUID {
	ids := make([]uuid.UUID, len(e.State.Players))
	for i := range e.State.Players {
		ids[i] = e.State.Players[i].ID
	}
	return ids
}
