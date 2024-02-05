package v1alpha1

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	core "github.com/mat285/boardgames/pkg/server/core/v1alpha1"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

var (
	_ connection.Client = new(Client)
)

type Client struct {
	ID       uuid.UUID
	Username string
	Engine   uuid.UUID
	Server   *Server
	Handler  connection.PacketHandler
}

func NewClient(id uuid.UUID, server *Server) *Client {
	return &Client{
		ID:     id,
		Server: server,
	}
}

func (c *Client) GetID() uuid.UUID {
	return c.ID
}

func (c *Client) GetUsername() string {
	return c.Username
}

func (c *Client) Connect(ctx context.Context, _ connection.ConnectionInfo) error {
	if c.Server == nil {
		return fmt.Errorf("No Server")
	}
	err := c.Server.ConnectClient(ctx, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Send(ctx context.Context, packet wire.Packet) error {
	packet.Origin = c.ID
	packet.Destination = c.Engine
	return c.Server.Receive(ctx, packet)
}

func (c *Client) Join(ctx context.Context, id uuid.UUID) error {
	c.Engine = id
	return c.Server.Join(ctx, c.ID, id)
}

func (c *Client) Listen(ctx context.Context, handler connection.PacketHandler) error {
	c.Handler = handler
	return nil
}

func (c *Client) Receive(ctx context.Context, packet wire.Packet) error {
	return c.Handler(ctx, packet)
}

func (c *Client) Pipe() connection.ClientInfo {
	return &core.Pipe{
		ID:       c.ID,
		Username: c.Username,
		Receiver: c,
	}
}
