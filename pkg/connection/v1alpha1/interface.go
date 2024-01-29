package v1alpha1

import "context"

type Interface interface {
	Connect(ctx context.Context) error
	Sender
	Listener
}

type Sender interface {
	Send(context.Context, Packet) error
}

type Listener interface {
	Listen(context.Context, PacketHandler) error
}
