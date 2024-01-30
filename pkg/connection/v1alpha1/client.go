package v1alpha1

import (
	"context"

	"github.com/blend/go-sdk/uuid"
)

type Client interface {
	Connect(context.Context, ConnectionInfo) error
	Join(context.Context, uuid.UUID) error
	Sender
	Listener
}
