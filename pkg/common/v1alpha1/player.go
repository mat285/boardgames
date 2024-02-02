package v1alpha1

import "github.com/blend/go-sdk/uuid"

type Player interface {
	GetID() uuid.UUID
	GetUsername() string
}
