package game

import (
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/games/machikoro/meta"
	common "github.com/mat285/boardgames/pkg/common/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

var (
	_ v1alpha1.StateData = new(State)
)

type State struct {
	meta.Object
	Config

	Players []Player
	Turn    common.TurnCounter
}

func NewState(players []Player, config Config) State {
	return State{
		Players: players,
		Config:  config,
		Turn:    common.NewTurnCounter(len(players), 0),
		// Board:   items.NewBoard(),
	}
}

func (s State) CurrentPlayer() (uuid.UUID, error) {
	player, err := s.GetCurrentPlayer()
	if err != nil {
		return nil, err
	}
	return player.ID, nil
}

func (s State) GetCurrentPlayer() (Player, error) {
	idx := s.Turn.CurrentPlayer()
	if len(s.Players) <= idx {
		return Player{}, fmt.Errorf("Invalid Player Index %d", idx)
	}
	return s.Players[idx], nil
}

func (s State) setCurrentPlayer(p Player) State {
	idx := s.Turn.CurrentPlayer()
	if len(s.Players) > idx {
		s.Players[idx] = p
	}
	return s
}

func (s State) IsDone() bool {
	return len(s.Winners()) > 0
}

func (s State) Winners() []uuid.UUID {
	return nil
}
