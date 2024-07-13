package http

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/games/splendor"
	splendorgame "github.com/mat285/boardgames/games/splendor/pkg/game"
	gamesclient "github.com/mat285/boardgames/pkg/client/http/v1alpha1"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
	messages "github.com/mat285/boardgames/pkg/messages/v1alpha1"
)

type Client struct {
	*gamesclient.Client
	Game splendor.Game
}

func New(c *gamesclient.Client) *Client {
	return &Client{
		Client: c,
		Game:   splendor.Game{},
	}
}

func (c *Client) GetState(ctx context.Context, id uuid.UUID) (*splendorgame.State, error) {
	packet, err := c.Client.GetState(ctx, id)
	if err != nil {
		return nil, err
	}
	untyped, err := c.Game.DeserializeState(&game.SerializedObject{
		ID:   id,
		Data: packet.Payload,
	})
	if err != nil {
		return nil, err
	}
	typed, ok := untyped.(splendorgame.State)
	if !ok {
		return nil, fmt.Errorf("wrong type for state")
	}
	return &typed, nil
}

func (c *Client) SendMove(ctx context.Context, game uuid.UUID, move splendorgame.Move) error {
	packet, err := messages.NewProvider(c.Game).MessagePlayerMove(&move, c.Client.UserID)
	if err != nil {
		return err
	}
	_, err = c.Client.SendPacket(ctx, game, c.Client.UserID, *packet)
	return err
}
