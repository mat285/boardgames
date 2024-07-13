package v1alpha1

import (
	"github.com/blend/go-sdk/uuid"
)

type Player struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func (p Player) GetID() uuid.UUID {
	return p.ID
}

func (p Player) GetUsername() string {
	return p.Username
}
