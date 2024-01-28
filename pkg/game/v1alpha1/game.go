package v1alpha1

import "github.com/blend/go-sdk/uuid"

type Game interface {
	Name() string
	Initialize([]uuid.UUID) (StateData, error)
	Load(StateData) error
}
