package v1alpha1

import (
	"context"

	client "github.com/mat285/boardgames/pkg/client/v1alpha1"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
	messages "github.com/mat285/boardgames/pkg/messages/v1alpha1"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

type Bot struct {
	client.Player
	Client   connection.Client
	Strategy Strategy
}

type Strategy interface {
	ChooseMove(context.Context, game.StateData) (game.Move, error)
}

func NewBot(username string, g game.Game, conn connection.Client, strategy Strategy) *Bot {
	b := &Bot{
		Player:   *client.NewPlayer(username, g, conn),
		Strategy: strategy,
		Client:   conn,
	}
	return b
}

func (b *Bot) Start(ctx context.Context) error {
	return b.Listen(ctx, b.Handle)
}

func (b *Bot) Handle(ctx context.Context, packet wire.Packet) error {
	switch packet.Type {
	case messages.PacketTypeRequestMove:
		state, err := b.Message.ExtractState(packet)
		if err != nil {
			return err
		}
		move, err := b.Request(ctx, state)
		if err != nil {
			return err
		}
		packet, err := b.Message.MessagePlayerMove(move, packet.ID)
		if err != nil {
			return err
		}
		return b.Client.Send(ctx, *packet)
	}
	// dropped
	return nil
}

func (b *Bot) Request(ctx context.Context, state game.StateData) (game.Move, error) {
	return b.Strategy.ChooseMove(ctx, state)
}
