package v1alpha1

import (
	"context"
	"net/http"
)

type Option func(*Client)

type RequestOption func(req *http.Request)

func OptUsername(name string) Option {
	return func(c *Client) {
		c.Username = name
	}
}

func OptContext(ctx context.Context) Option {
	return func(c *Client) {
		c.Ctx = ctx
	}
}

func OptConfig(config Config) Option {
	return func(c *Client) {
		c.Config = config
	}
}

func OptScheme(scheme string) Option {
	return func(c *Client) {
		c.Config.Scheme = scheme
	}
}

func OptHost(host string) Option {
	return func(c *Client) {
		c.Config.Host = host
	}
}

func OptPort(port uint16) Option {
	return func(c *Client) {
		c.Config.Port = port
	}
}
