package splendor

import (
	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/games/splendor/pkg/game"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

var (
	_ v1alpha1.Game = new(Game)
)

type Game struct {
	v1alpha1.Serializer
	Config game.Config
}

func New(config game.Config) *Game {
	return &Game{
		Config: config,
	}
}

func (g *Game) Initialize(pids []uuid.UUID) (v1alpha1.StateData, error) {
	players := make([]game.Player, len(pids))
	for i := range pids {
		players[i] = game.NewPlayer(pids[i])
	}
	return game.NewState(players, g.Config), nil
}

func (g *Game) Load(state v1alpha1.StateData) error {
	return nil
}

// func (g *Game) SerializeMove(m v1alpha1.Move) ([]byte,error)
