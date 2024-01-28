package v1alpha1

import (
	"context"

	"github.com/blend/go-sdk/uuid"
)

type PlayerConnection interface {
	Request(context.Context, MoveRequest) (Move, error)
	Accept(context.Context, Message) error
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
