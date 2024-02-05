package v1alpha1

import (
	"context"

	"github.com/blend/go-sdk/uuid"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	engine "github.com/mat285/boardgames/pkg/engine/v1alpha1"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

type Pipe struct {
	ID       uuid.UUID
	Username string

	Receiver connection.Receiver
}

func PipeEngine(e *engine.Engine) connection.ServerInfo {
	return &Pipe{
		ID:       e.ID,
		Receiver: e,
	}
}

func (p *Pipe) GetID() uuid.UUID {
	return p.ID
}

func (p Pipe) GetUsername() string {
	return p.Username
}

func (p Pipe) Send(ctx context.Context, packet wire.Packet) error {
	return p.Receiver.Receive(ctx, packet)
}
