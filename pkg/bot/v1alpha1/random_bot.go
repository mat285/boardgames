package v1alpha1

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type RandomStrategy struct {
}

func NewRandom() Strategy {
	return &RandomStrategy{}
}

func (r *RandomStrategy) ChooseMove(ctx context.Context, state v1alpha1.StateData) (v1alpha1.Move, error) {
	possible, err := state.ValidMoves()
	if err != nil {
		return nil, err
	}
	if len(possible) == 0 {
		return nil, fmt.Errorf("No possible moves")
	}
	return possible[rand.Intn(len(possible))], nil
}
