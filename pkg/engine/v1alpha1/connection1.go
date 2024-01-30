package v1alpha1

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"time"

// 	"github.com/blend/go-sdk/uuid"
// 	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
// 	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
// )

// const (
// 	MessageHeaderRequireResponse = "require-response"
// )

// var (
// 	_ game.Connection = new(Connection)
// )

// type Connection struct {
// 	Backend connection.Interface
// 	Game    game.Game

// 	game.Sender
// 	handler game.MessageHandler

// 	waitingLock sync.Mutex
// 	waiting     map[string]chan *game.Message
// }

// func NewConnection(g game.Game, backend connection.Interface) *Connection {
// 	c := &Connection{
// 		Backend:     backend,
// 		Game:        g,
// 		waitingLock: sync.Mutex{},
// 		waiting:     make(map[string]chan *game.Message),
// 	}
// 	return c
// }

// func (c *Connection) Receive(ctx context.Context, handler game.MessageHandler) error {
// 	return c.Backend.Listen(ctx, c.handleMessagePacket(handler))
// }

// func (c *Connection) handleMessagePacket(handler game.MessageHandler) connection.PacketHandler {
// 	return func(ctx context.Context, packet connection.Packet) error {
// 		mtype := packet.Type - connection.MinDataPacketType()
// 		resp, err := handler(ctx, game.Message{Type: game.MessageType(mtype), Data: packet.Payload})
// 		if err != nil {
// 			return err
// 		}
// 		return c.handleResponse(ctx, packet.Request, resp)
// 	}

// }

// func (c *Connection) Send(ctx context.Context, m game.Message) error {
// 	_, err := c.send(ctx, m, nil)
// 	return err
// }

// func (c *Connection) Request(ctx context.Context, m game.Message) (*game.Message, error) {
// 	id, err := c.send(ctx, m, map[string]string{
// 		MessageHeaderRequireResponse: "true",
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return c.awaitResponse(ctx, id)
// }

// func (c *Connection) send(ctx context.Context, m game.Message, headers connection.HeaderOptions) (uuid.UUID, error) {
// 	packet, err := c.serializeMessage(m)
// 	if err != nil {
// 		return nil, err
// 	}
// 	packet.Header.Options.Append(headers)
// 	return packet.ID, c.Backend.Send(ctx, packet)
// }

// func (c *Connection) handleResponse(ctx context.Context, pid uuid.UUID, msg *game.Message) error {
// 	c.waitingLock.Lock()
// 	defer c.waitingLock.Unlock()
// 	key := pid.ToFullString()
// 	if await, has := c.waiting[key]; has {
// 		if len(await) > 0 {
// 			close(await)
// 			delete(c.waiting, key)
// 			return fmt.Errorf("Already responded")
// 		}
// 		await <- msg
// 	}
// 	return nil
// }

// func (c *Connection) awaitResponse(ctx context.Context, pid uuid.UUID) (*game.Message, error) {
// 	await := make(chan *game.Message, 1)

// 	c.waitingLock.Lock()
// 	if _, has := c.waiting[pid.ToFullString()]; has {
// 		c.waitingLock.Unlock()
// 		return nil, fmt.Errorf("Duplicate Packet ID")
// 	}
// 	c.waiting[pid.ToFullString()] = await
// 	c.waitingLock.Unlock()
// 	timeout := time.After(30 * time.Second)

// 	select {
// 	case <-ctx.Done():
// 		return nil, ctx.Err()
// 	case <-timeout:
// 		c.waitingLock.Lock()
// 		respChan := c.waiting[pid.ToFullString()]
// 		if respChan != nil {
// 			close(respChan)
// 			delete(c.waiting, pid.ToFullString())
// 		}
// 		c.waitingLock.Unlock()
// 		return nil, fmt.Errorf("Timeout waiting for response")
// 	case resp := <-await:
// 		return resp, nil
// 	}
// }

// func (c *Connection) serializeMessage(m game.Message) (connection.Packet, error) {
// 	return connection.NewPacket(connection.OptPacketType(connection.MinDataPacketType()+connection.PacketType(m.Type)), connection.OptPacketPayload(m.Data)), nil
// }
