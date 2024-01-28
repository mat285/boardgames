package v1alpha1

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type Bot struct {
	ID       uuid.UUID
	Strategy Strategy
}

type Strategy interface {
	ChooseMove(context.Context, v1alpha1.StateData) (v1alpha1.Move, error)
}

func NewBot(strategy Strategy) *Bot {
	b := &Bot{
		Strategy: strategy,
	}
	return b
}

func (b *Bot) Request(ctx context.Context, req v1alpha1.MoveRequest) (v1alpha1.Move, error) {
	return b.Strategy.ChooseMove(ctx, req.State)
}

func (b *Bot) Accept(ctx context.Context, messsage v1alpha1.Message) error {
	fmt.Println(messsage)
	return nil
}
