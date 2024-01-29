package v1alpha1

import (
	"github.com/blend/go-sdk/uuid"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type Player interface {
	game.Player
}

type ConnectedPlayer struct {
	game.BasePlayer
	*Connection
}

func NewConnectedPlayer(id uuid.UUID, username string, s game.Game, server connection.Server) ConnectedPlayer {
	return ConnectedPlayer{
		BasePlayer: game.NewBasePlayer(id, username),
		Connection: NewConnection(s, server.Serve()),
	}
}

func (cp ConnectedPlayer) GetConnection() game.Connection {
	return cp.Connection
}
