package v1alpha1

import (
	"context"

	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

type Interface interface {
	Sender
	Receiver
}

type Sender interface {
	Send(context.Context, wire.Packet) error
}

type Listener interface {
	Listen(context.Context, PacketHandler) error
}

type Receiver interface {
	Receive(context.Context, wire.Packet) error
}
