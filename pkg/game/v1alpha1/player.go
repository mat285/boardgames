package v1alpha1

import (
	"context"

	"github.com/blend/go-sdk/uuid"
)

type PlayerConnection interface {
	RequestMove(context.Context, MoveRequest) (Move, error)
}

type Player struct {
	ID       uuid.UUID
	Username string

	Connection PlayerConnection
}

func NewPlayer(id uuid.UUID, conn PlayerConnection) Player {
	return Player{
		ID:         id,
		Connection: conn,
	}
}
