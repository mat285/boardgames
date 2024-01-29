package v1alpha1

import (
	"encoding/json"

	"github.com/blend/go-sdk/uuid"
)

type Serializer interface {
	StateSerializer
	MoveSerializer
}

type StateSerializer interface {
	SerializeState(StateData) (*SerializedObject, error)
	DeserializeState(*SerializedObject) (StateData, error)
}

type MoveSerializer interface {
	SerializeMove(Move) (*SerializedObject, error)
	DeserializeMove(*SerializedObject) (Move, error)
}

type Serializable interface {
	Serialize(Serializer) (*SerializedObject, error)
	Deserialize(Serializer) (*SerializedObject, error)
}

type SerializedObject struct {
	ID   uuid.UUID
	Data []byte
}

func (so SerializedObject) Serialize() ([]byte, error) {
	return json.Marshal(so)
}
