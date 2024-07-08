package v1alpha1

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/blend/go-sdk/uuid"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

var (
	_ connection.Client = new(Client)
)

func (c *Client) Connect(ctx context.Context, _ connection.ConnectionInfo) error {
	_, err := c.getUserID(ctx)
	return err
}

func (c *Client) GetUser(ctx context.Context, name string) (uuid.UUID, error) {
	req, err := c.NewRequest(
		ctx,
		http.MethodGet,
		"/api/v1alpha1/user/:name",
		map[string]string{
			":name": name,
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	var id uuid.UUID
	return id, c.JSON(ctx, req, &id)
}

func (c *Client) getUserID(ctx context.Context) (uuid.UUID, error) {
	if c.UserID.IsZero() {
		id, err := c.GetUser(ctx, c.Username)
		if err != nil {
			return nil, err
		}
		c.UserID = id
	}
	return c.UserID, nil
}

func (c *Client) Listen(ctx context.Context, handler connection.PacketHandler) error {
	userID, err := c.getUserID(ctx)
	if err != nil {
		return err
	}
	dialer := NewWebsocketDialer(c.webSocketsAddress(c.Username), userID, c.Username)
	err = dialer.Listen(ctx, handler)
	fmt.Println(err)
	return err
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
	uid, err := c.getUserID(ctx)
	if err != nil {
		return err
	}
	p := v1alpha1.Player{
		ID:       uid,
		Username: c.Username,
	}
	req, err := c.NewJSONRequest(
		ctx,
		http.MethodPost,
		"/api/v1alpha1/game/:id/join",
		map[string]string{
			":id": id.ToFullString(),
		},
		p,
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

func (c *Client) webSocketsAddress(user string) string {
	url := fmt.Sprintf("ws://%s", path.Join(c.Config.HostPort(), "/api/v1alpha1/websockets/"+user))
	return url
}
