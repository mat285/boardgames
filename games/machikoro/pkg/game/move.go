package game

import (
	"fmt"

	"github.com/mat285/boardgames/games/machikoro/meta"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

var (
	_ v1alpha1.Move = new(Move)
)

type Move struct {
	meta.Object
}

func (m *Move) Apply(raw v1alpha1.StateData) (*v1alpha1.MoveResult, error) {
	state, ok := raw.(State)
	if !ok {
		return nil, fmt.Errorf("Invalid State Type")
	}
	state = state
	return nil, nil

}
