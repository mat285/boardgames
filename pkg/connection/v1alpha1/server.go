package v1alpha1

import (
	"context"

	"github.com/blend/go-sdk/uuid"
)

type Router interface {
	Receiver
	ConnectServer(context.Context, ServerInfo) error
	ConnectClient(context.Context, ClientInfo) error
	GetClient(uuid.UUID) ClientInfo
	GetServer(uuid.UUID) ServerInfo
}

type ClientInfo interface {
	GetID() uuid.UUID
	GetUsername() string
	Sender
}

type ServerInfo interface {
	GetID() uuid.UUID
	Sender
}
