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
	ID     uuid.UUID
	Engine uuid.UUID
	Server *Server

	packets chan packet
}

type packet struct {
	p wire.Packet
	c *clientReceive
}

type clientReceive struct {
	c    *Client
	errs chan error
}

func NewClient(server *Server) *Client {
	return &Client{
		Server:  server,
		packets: make(chan packet, 10),
	}
}

func (c *Client) Connect(ctx context.Context, _ connection.ConnectionInfo) error {
	if c.Server == nil {
		return fmt.Errorf("No Server")
	}
	id, err := c.Server.Connect(ctx, connection.ClientInfo{
		ID:       c.ID,
		Username: "",
		Sender:   &clientReceive{c: c, errs: make(chan error, 1)},
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
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case p := <-c.packets:
			p.c.errs <- handler(ctx, p.p)
		}
	}
}

func (c *clientReceive) Send(ctx context.Context, p wire.Packet) error {
	fmt.Println("client pipe sending to client")
	c.c.packets <- packet{p: p, c: c}
	fmt.Println()
	return <-c.errs
}
