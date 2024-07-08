package v1alpha1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/blend/go-sdk/uuid"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

type Client struct {
	Ctx        context.Context
	Config     Config
	Username   string
	UserID     uuid.UUID
	HTTPClient http.Client
	Websocket  *WebsocketDialer
}

func New(opts ...Option) *Client {
	c := new(Client)
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) GetID() uuid.UUID {
	return c.UserID
}

func (c *Client) Send(ctx context.Context, p wire.Packet) error {
	if c.Websocket == nil {
		return fmt.Errorf("no websocket conn")
	}
	return c.Websocket.Websocket.Send(ctx, p)
}

func (c *Client) Response(ctx context.Context, r *http.Request) (*http.Response, error) {
	return c.HTTPClient.Do(r)
}

func (c *Client) ResponseBody(ctx context.Context, r *http.Request) (io.Reader, error) {
	resp, err := c.HTTPClient.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status code %d\n Body: %s", resp.StatusCode, string(b))
	}
	return resp.Body, nil
}

func (c *Client) JSON(ctx context.Context, r *http.Request, out interface{}) error {
	resp, err := c.ResponseBody(ctx, r)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}

func (c *Client) Do(ctx context.Context, r *http.Request) error {
	resp, err := c.HTTPClient.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("bad status code %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) NewJSONRequest(ctx context.Context, verb string, route string, params map[string]string, body interface{}, opts ...RequestOption) (*http.Request, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return c.NewRequest(ctx, verb, route, params, bytes.NewReader(data), opts...)
}

func (c *Client) NewRequest(ctx context.Context, verb string, route string, params map[string]string, body io.Reader, opts ...RequestOption) (*http.Request, error) {
	for k, v := range params {
		route = strings.ReplaceAll(route, k, v)
	}
	req, err := http.NewRequestWithContext(ctx, verb, c.Config.URL(route), body)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}
