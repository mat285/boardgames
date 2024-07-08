package v1alpha1

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/blend/go-sdk/uuid"
// 	"github.com/gorilla/websocket"

// 	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
// 	core "github.com/mat285/boardgames/pkg/core/v1alpha1"
// 	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
// )

// var (
// 	_ connection.ClientInfo = new(Client)
// )

// const (
// 	// Time allowed to write a message to the peer.
// 	writeWait = 10 * time.Second

// 	// Time allowed to read the next pong message from the peer.
// 	pongWait = 60 * time.Second

// 	// Send pings to peer with this period. Must be less than pongWait.
// 	pingPeriod = (pongWait * 9) / 10

// 	// Maximum message size allowed from peer.
// 	maxMessageSize = 512
// )

// type Client struct {
// 	core.User
// 	Conn *websocket.Conn

// 	Handle connection.PacketHandler

// 	outbound chan *wire.Packet
// }

// func NewClient(id uuid.UUID, username string, conn *websocket.Conn, handle connection.PacketHandler) *Client {
// 	c := &Client{
// 		User:     core.NewUser(id, username),
// 		Conn:     conn,
// 		Handle:   handle,
// 		outbound: make(chan *wire.Packet, 10),
// 	}
// 	return c
// }

// func (c *Client) Start(ctx context.Context) {
// 	go c.write(ctx)
// 	go c.read(ctx)
// }

// func (c *Client) Send(ctx context.Context, packet wire.Packet) (err error) {
// 	defer func() {
// 		e := recover()
// 		if e != nil {
// 			err = fmt.Errorf("Connection closed")
// 		}
// 	}()
// 	if c.Conn == nil {
// 		return fmt.Errorf("No connection")
// 	}
// 	select {
// 	case <-ctx.Done():
// 		return ctx.Err()
// 	default:
// 	}
// 	c.outbound <- &packet
// 	return nil
// }

// func (c *Client) Stop(ctx context.Context) {
// 	if c.Conn == nil {
// 		return
// 	}

// 	c.Conn.Close()
// 	c.Conn = nil

// }

// func (c *Client) write(ctx context.Context) {
// 	ticker := time.NewTicker(pingPeriod)
// 	defer func() {
// 		ticker.Stop()
// 		c.Conn.Close()
// 	}()
// 	for {
// 		select {
// 		case packet, ok := <-c.outbound:
// 			if !ok {
// 				// The hub closed the channel.
// 				c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}

// 			if packet == nil {
// 				continue
// 			}
// 			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

// 			raw, err := packet.Serialize()
// 			if err != nil {
// 				// drop invalid packet
// 				continue
// 			}
// 			// TODO handle write errors
// 			c.Conn.WriteMessage(websocket.BinaryMessage, raw)

// 		case <-ticker.C:
// 			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				return
// 			}
// 		case <-ctx.Done():
// 			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
// 			return
// 		}
// 	}
// }

// func (c *Client) read(ctx context.Context) {
// 	defer func() {
// 		c.Conn.Close()
// 	}()

// 	c.Conn.SetReadLimit(maxMessageSize)
// 	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
// 	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

// 	for {
// 		_, rawData, err := c.Conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}
// 		packet, err := wire.DeserializePacket(rawData)
// 		if err != nil {
// 			_ = c.sendError(ctx, err)
// 		}
// 		if packet == nil {
// 			// todo
// 			continue
// 		}
// 		err = c.Handle(ctx, *packet)
// 		if err != nil {
// 			_ = c.sendError(ctx, err)
// 		}
// 	}
// }

// func (c *Client) sendError(ctx context.Context, err error) error {
// 	// TODO
// 	return nil
// }
