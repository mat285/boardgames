package v1alpha1

import "github.com/blend/go-sdk/uuid"

type StateData interface {
	Meta() Meta
	CurrentPlayer() (uuid.UUID, error)
	IsDone() bool
	Winners() []uuid.UUID
	ValidMoves() ([]Move, error)
}
