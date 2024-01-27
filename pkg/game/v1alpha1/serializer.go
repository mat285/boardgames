package v1alpha1

type Serializer interface {
	SerializeState(StateData) ([]byte, error)
	DeserializeState([]byte) (StateData, error)

	SerializeMove(Move) ([]byte, error)
	DeserializeMove([]byte) (Move, error)
}

type Serializable interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}
