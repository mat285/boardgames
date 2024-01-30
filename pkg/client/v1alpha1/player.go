package v1alpha1

import (
	"context"

	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
	logger "github.com/mat285/boardgames/pkg/logger/v1alpha1"
	messages "github.com/mat285/boardgames/pkg/messages/v1alpha1"
)

type Player struct {
	game.Player
	connection.Client
	logger.Interface
	game.Game
	connection.ConnectionInfo
	Message messages.Provider
}

func NewPlayer(username string, g game.Game, client connection.Client) *Player {
	return &Player{
		Player: game.Player{
			Username: username,
		},
		Game:    g,
		Message: messages.NewProvider(g),
		Client:  client,
	}
}

func (p *Player) Listen(ctx context.Context, handler connection.PacketHandler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		err := p.listen(ctx, handler)
		if err != nil {
			// TODO retry connection
			return err
		}
	}
}

func (p *Player) listen(ctx context.Context, handler connection.PacketHandler) error {
	err := p.Connect(ctx, p.ConnectionInfo)
	if err != nil {
		return err
	}
	return p.Client.Listen(ctx, handler)
}
