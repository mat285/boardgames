package game

import (
	"encoding/json"
	"fmt"

	"github.com/mat285/boardgames/games/splendor/pkg/items"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

var _ v1alpha1.Move = new(Move)

type Move struct {
	Collect  *CollectMove
	Purchase *CardMove
	Reserve  *CardMove
}

type CollectMove struct {
	Take   items.GemCount
	Return items.GemCount
}

type CardMove struct {
	Card items.Card
}

func NewMove() *Move {
	return &Move{}
}

func (m *Move) Apply(raw v1alpha1.StateData) (*v1alpha1.MoveResult, error) {
	state, ok := raw.(State)
	if !ok {
		return nil, fmt.Errorf("Invalid State Type")
	}

	valid, err := m.Validate()
	if err != nil {
		return nil, err
	}
	if !valid {
		return &v1alpha1.MoveResult{
			Valid: false,
			State: state,
		}, nil
	}
	res := &v1alpha1.MoveResult{}
	res.State, res.Valid = m.apply(state)
	return res, nil

}

func (m *Move) apply(state State) (State, bool) {
	if m.Collect != nil {
		return state.applyCollect(*m.Collect)
	} else if m.Purchase != nil {
		return state.applyPurchase(*m.Purchase)
	} else if m.Reserve != nil {
		return state.applyReserve(*m.Reserve)
	} else {
		return state, false
	}
}

func (m *Move) Validate() (bool, error) {
	nonNil := 0
	if m.Collect != nil {
		nonNil++
	}
	if m.Purchase != nil {
		nonNil++
	}
	if m.Reserve != nil {
		nonNil++
	}
	if nonNil != 1 {
		return false, fmt.Errorf("Invalid Move")
	}

	return true, nil
}

func (m *Move) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Move) Deserialize(data []byte) error {
	return json.Unmarshal(data, m)
}

func MoveSliceToMoveSlice(moves []*Move) []v1alpha1.Move {
	ret := make([]v1alpha1.Move, len(moves))
	for i := range moves {
		ret[i] = moves[i]
	}
	return ret
}
