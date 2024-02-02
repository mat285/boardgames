package game

import (
	"encoding/json"
	"fmt"

	"github.com/mat285/boardgames/games/splendor/meta"
	"github.com/mat285/boardgames/games/splendor/pkg/items"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

var _ v1alpha1.Move = new(Move)

type Move struct {
	meta.Object
	Pass     *PassMove
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

// this way there's always at least one option
type PassMove struct{}

func NewMove() *Move {
	return &Move{}
}

func NewPassMove() *Move {
	return &Move{
		Pass: &PassMove{},
	}
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
	res.State, res.Valid, err = state.apply(*m)
	return res, err

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
	if m.Pass != nil {
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

func (c CollectMove) Validate(gems items.GemCount) error {
	take := c.Take.ToMap()
	total := 0
	double := false

	for k, v := range take {
		if v < 0 {
			return fmt.Errorf("Cannot take negative gems")
		}
		if v == 2 {
			double = true
			if gems.Get(k) < 4 {
				return fmt.Errorf("Cannot take 2 when fewer than 4 remain")
			}
		}
		total += v
	}

	if total > 3 {
		return fmt.Errorf("Can take at most 3 gems")
	}

	if double && total != 2 {
		return fmt.Errorf("Can only take one of each gem unless only taking 2 total")
	}
	return nil
}
