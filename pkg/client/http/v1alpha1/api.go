package v1alpha1

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/blend/go-sdk/uuid"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
	server "github.com/mat285/boardgames/server/http/v1alpha1"
)

var (
	_ connection.Client = new(Client)
)

func (c *Client) Connect(ctx context.Context, _ connection.ConnectionInfo) error {
	if c.Websocket != nil {
		return nil
	}
	if c.UserID.IsZero() {
		err := c.Login(ctx)
		if err != nil {
			return err
		}
	}
	c.Websocket = NewWebsocketDialer(c.webSocketsAddress(c.Username), c.UserID, c.Username)
	return nil
}

func (c *Client) Listen(ctx context.Context, handler connection.PacketHandler) error {
	if c.Websocket == nil {
		return fmt.Errorf("please connect first")
	}
	return c.Websocket.Listen(ctx, handler)
}

func (c *Client) Close(ctx context.Context) error {
	if c.Websocket == nil {
		return nil
	}
	err := c.Websocket.Close(ctx)
	c.Websocket = nil
	return err
}

func (c *Client) Send(ctx context.Context, packet wire.Packet) error {
	if c.Websocket == nil {
		return fmt.Errorf("please connect first")
	}
	return c.Websocket.Send(ctx, packet)
}

func (c *Client) Login(ctx context.Context) error {
	body := game.Player{
		Username: c.Username,
	}
	req, err := c.NewJSONRequest(
		ctx,
		http.MethodPost,
		"/api/v1alpha1/user/login",
		nil,
		body,
	)
	if err != nil {
		return err
	}
	err = c.JSON(ctx, req, &body)
	if err != nil {
		return err
	}
	c.UserID = body.ID
	return nil
}

func (c *Client) GetUserGames(ctx context.Context) ([]server.UserGame, error) {
	req, err := c.NewRequest(
		ctx,
		http.MethodGet,
		"/api/v1alpha1/user/games",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	var res []server.UserGame
	return res, c.JSON(ctx, req, &res)
}

func (c *Client) NewGame(ctx context.Context, name string, config interface{}) (uuid.UUID, error) {
	req, err := c.NewJSONRequest(
		ctx,
		http.MethodPost,
		"/api/v1alpha1/games/:name/new",
		map[string]string{
			":name": name,
		},
		config,
	)
	if err != nil {
		return nil, err
	}
	var id uuid.UUID
	return id, c.JSON(ctx, req, &id)
}

func (c *Client) Join(ctx context.Context, id uuid.UUID) error {
	req, err := c.NewRequest(
		ctx,
		http.MethodPost,
		"/api/v1alpha1/game/:id/join",
		map[string]string{
			":id": id.ToFullString(),
		},
		nil,
	)
	if err != nil {
		return err
	}

	return c.Do(ctx, req)
}

func (c *Client) Start(ctx context.Context, id uuid.UUID) error {
	req, err := c.NewJSONRequest(
		ctx,
		http.MethodPost,
		"/api/v1alpha1/game/:id/start",
		map[string]string{
			":id": id.ToFullString(),
		},
		nil,
	)
	if err != nil {
		return err
	}

	return c.Do(ctx, req)
}

type GameResponse struct {
	ID    uuid.UUID           `json:"id"`
	State *v1alpha1.StateData `json:"state"`
}

func (c *Client) GetState(ctx context.Context, id uuid.UUID) (*wire.Packet, error) {
	req, err := c.NewJSONRequest(
		ctx,
		http.MethodGet,
		"/api/v1alpha1/game/:id/state",
		map[string]string{
			":id": id.ToFullString(),
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	var res wire.Packet
	return &res, c.JSON(ctx, req, &res)
}

func (c *Client) SendPacket(ctx context.Context, id, player uuid.UUID, move wire.Packet) (*wire.Packet, error) {
	req, err := c.NewJSONRequest(
		ctx,
		http.MethodPost,
		"/api/v1alpha1/game/:id/packet",
		map[string]string{
			":id": id.ToFullString(),
		},
		move,
	)
	if err != nil {
		return nil, err
	}
	var res wire.Packet
	return &res, c.JSON(ctx, req, &res)
}

func (c *Client) webSocketsAddress(user string) string {
	url := fmt.Sprintf("ws://%s", path.Join(c.Config.HostPort(), "/api/v1alpha1/websockets/"+user))
	return url
}
