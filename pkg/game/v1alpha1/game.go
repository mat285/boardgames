package v1alpha1

import "github.com/blend/go-sdk/uuid"

type Game interface {
	Meta() Meta
	Serializer() Serializer
	Initialize([]uuid.UUID) (StateData, error)
	Load(StateData) error
}
