package v1alpha1

type Move interface {
	Serializable
	Apply(StateData) (*MoveResult, error)
}

type MoveRequest struct {
	State StateData
}

type MoveResult struct {
	Valid bool
	State StateData
}
