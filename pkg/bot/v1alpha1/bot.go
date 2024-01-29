package v1alpha1

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	engine "github.com/mat285/boardgames/pkg/engine/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

var _ engine.Player = new(Bot)

type Bot struct {
	engine.ConnectedPlayer
	Strategy Strategy
}

type Strategy interface {
	ChooseMove(context.Context, v1alpha1.StateData) (v1alpha1.Move, error)
}

func NewBot(id uuid.UUID, username string, conn connection.Interface, strategy Strategy) *Bot {
	b := &Bot{
		ConnectedPlayer: engine.NewConnectedPlayer(id, username, conn),
		Strategy:        strategy,
	}
	return b
}

func (b *Bot) Connect(ctx context.Context) error {
	go b.ConnectedPlayer.Receive(ctx, b.Handle)
	return nil
}

func (b *Bot) Handle(ctx context.Context, message v1alpha1.Message) (*v1alpha1.Message, error) {
	switch message.Type {
	case v1alpha1.MessageTypeRequestMove:
		req, err := v1alpha1.MoveRequestFromMessage(message)
		if err != nil {
			return nil, err
		}
		move, err := b.Request(ctx, req)
		if err != nil {
			return nil, err
		}
		msg := v1alpha1.MessagePlayerMove(move)
		return &msg, nil
	}
	// dropped
	return nil, nil
}

func (b *Bot) Request(ctx context.Context, req v1alpha1.MoveRequest) (v1alpha1.Move, error) {
	return b.Strategy.ChooseMove(ctx, req.State)
}

func (b *Bot) Send(ctx context.Context, messsage v1alpha1.Message) error {
	fmt.Println(messsage)
	return nil
}
