package v1alpha1

import (
	"github.com/blend/go-sdk/uuid"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type Player struct {
	game.Player
	connection.Sender
}

func NewPlayer(id uuid.UUID, username string, conn connection.Sender) *Player {
	return &Player{
		Player: game.Player{ID: id, Username: username},
		Sender: conn,
	}
}
