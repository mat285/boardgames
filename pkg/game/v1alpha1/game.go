package v1alpha1

import "github.com/blend/go-sdk/uuid"

type Game interface {
	Initialize([]uuid.UUID) (StateData, error)
	Load(StateData) error
	Serializer
}
