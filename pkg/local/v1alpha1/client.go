package v1alpha1

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

var (
	_ connection.Client = new(Client)
)

type Client struct {
	ID      uuid.UUID
	Engine  uuid.UUID
	Server  *Server
	Handler connection.PacketHandler
}

func NewClient(server *Server) *Client {
	return &Client{
		Server: server,
	}
}

func (c *Client) Connect(ctx context.Context, _ connection.ConnectionInfo) error {
	if c.Server == nil {
		return fmt.Errorf("No Server")
	}
	id, err := c.Server.Connect(ctx, connection.ClientInfo{
		ID:       c.ID,
		Username: "",
		Sender:   connection.PipeReceiverToSender(c),
	})
	if err != nil {
		return err
	}
	c.ID = id
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
