package v1alpha1

import (
	"fmt"
	"strings"
)

type Config struct {
	Scheme string
	Host   string
	Port   uint16
}

func (c Config) URL(route string) string {
	route = strings.TrimPrefix(route, "/")
	return fmt.Sprintf("%s/%s", c.Addr(), route)
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s://%s", c.Scheme, c.HostPort())
}

func (c Config) HostPort() string {
	if c.Port == 0 {
		return c.Host
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
