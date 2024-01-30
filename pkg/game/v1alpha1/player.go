package v1alpha1

import (
	"github.com/blend/go-sdk/uuid"
)

type Player struct {
	ID       uuid.UUID
	Username string
}
