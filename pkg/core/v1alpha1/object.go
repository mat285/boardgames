package v1alpha1

import "github.com/blend/go-sdk/uuid"

type Object struct {
	ID uuid.UUID
}

func NewObject(id uuid.UUID) Object {
	return Object{
		ID: id,
	}
}

func (o Object) GetID() uuid.UUID {
	return o.ID
}
