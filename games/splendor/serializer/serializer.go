package serializer

import (
	"encoding/json"
	"fmt"

	"github.com/mat285/boardgames/games/splendor/meta"
	"github.com/mat285/boardgames/games/splendor/pkg/game"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

var (
	_ v1alpha1.Serializer = new(Serializer)
)

type Get struct {
}

func (Get) Serializer() v1alpha1.Serializer {
	return Serializer{}
}

type Serializer struct {
	meta.Object
}

func (s Serializer) SerializeMove(move v1alpha1.Move) (*v1alpha1.SerializedObject, error) {
	if !move.Meta().ID().Equal(s.Meta().ID()) {
		return nil, fmt.Errorf("Incorrect object metadata for serializer")
	}
	bytes, err := json.Marshal(move)
	if err != nil {
		return nil, err
	}
	return &v1alpha1.SerializedObject{
		ID:   move.Meta().ID(),
		Data: bytes,
	}, nil
}

func (s Serializer) DeserializeMove(obj *v1alpha1.SerializedObject) (v1alpha1.Move, error) {
	if obj == nil {
		return nil, nil
	}
	if !obj.ID.Equal(meta.Meta{}.ID()) {
		return nil, fmt.Errorf("Incorrect Serializer for Object")
	}
	var move game.Move
	return &move, json.Unmarshal(obj.Data, &move)
}

func (s Serializer) SerializeState(state v1alpha1.StateData) (*v1alpha1.SerializedObject, error) {
	if !state.Meta().ID().Equal(s.Meta().ID()) {
		return nil, fmt.Errorf("Incorrect object metadata for serializer")
	}
	bytes, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}
	return &v1alpha1.SerializedObject{
		ID:   state.Meta().ID(),
		Data: bytes,
	}, nil
}

func (s Serializer) DeserializeState(obj *v1alpha1.SerializedObject) (v1alpha1.StateData, error) {
	if obj == nil {
		return nil, nil
	}
	if !obj.ID.Equal(s.Meta().ID()) {
		return nil, fmt.Errorf("Incorrect Serializer for Object")
	}
	var state game.State
	return state, json.Unmarshal(obj.Data, &state)
}
