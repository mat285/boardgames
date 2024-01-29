package v1alpha1

import "github.com/blend/go-sdk/uuid"

type Meta interface {
	ID() uuid.UUID
	Name() string
}
