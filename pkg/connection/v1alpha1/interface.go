package v1alpha1

import "context"

type Interface interface {
	Sender
	Listener
}

type Sender interface {
	Send(context.Context, Packet) error
}

type Listener interface {
	Listen(context.Context, PacketHandler) error
}

type Server interface {
	Connect(context.Context) (Interface, error)
	Serve() Interface
}
