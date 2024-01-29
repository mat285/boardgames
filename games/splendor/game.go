package splendor

import (
	"encoding/json"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/games/splendor/meta"
	"github.com/mat285/boardgames/games/splendor/pkg/game"
	"github.com/mat285/boardgames/games/splendor/serializer"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

const (
	Name = "splendor"
)

var (
	ID = uuid.V4()
)

var (
	_ v1alpha1.Game = new(Game)
)

type Game struct {
	meta.Object
	serializer.Get

	Config game.Config
}

func NewGameWithConfig(config game.Config) *Game {
	return &Game{
		Config: config,
	}
}

func New(config interface{}) (v1alpha1.Game, error) {
	if config == nil {
		return NewGameWithConfig(game.StandardConfig()), nil
	}
	typed, ok := config.(game.Config)
	if !ok {
		typed = game.StandardConfig()
	}
	return NewGameWithConfig(typed), nil
}

func NewConfig() interface{} {
	return game.Config{}
}

func (g Game) Name() string {
	return Name
}

func (g *Game) Initialize(pids []uuid.UUID) (v1alpha1.StateData, error) {
	players := make([]game.Player, len(pids))
	for i := range pids {
		players[i] = game.NewPlayer(pids[i])
	}
	return game.NewState(players, g.Config), nil
}

func (g *Game) Load(state v1alpha1.StateData) error {
	typed, ok := state.(game.State)
	if !ok {
		return fmt.Errorf("Invalid State for Game")
	}
	g.Config = typed.Config
	return nil
}

func (g Game) SerializeMove(move v1alpha1.Move) ([]byte, error) {
	return json.Marshal(move)
}

func (g Game) DeserializeMove(data []byte) (v1alpha1.Move, error) {
	var move game.Move
	return &move, json.Unmarshal(data, &move)
}

func (g Game) SerializeState(state v1alpha1.StateData) ([]byte, error) {
	typed, ok := state.(game.State)
	if !ok {
		return nil, fmt.Errorf("Wrong state for this game")
	}
	return json.Marshal(typed)
}

func (g Game) DeserializeState(data []byte) (v1alpha1.StateData, error) {
	var state game.State
	err := json.Unmarshal(data, &state)
	return state, err
}
