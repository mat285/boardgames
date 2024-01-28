package v1alpha1

import "context"

type Sender interface {
	Send(context.Context, Packet) error
}

type Listener interface {
	Listen(context.Context, PacketHandler)
}
