package v1alpha1

import connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"

type ConnectionOption func(*Connection)

func OptConnectionBackend(b connection.Interface) ConnectionOption {
	return func(c *Connection) {
		c.Backend = b
	}
}
