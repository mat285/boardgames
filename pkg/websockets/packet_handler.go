package websockets

import "context"

type PacketHandler func(context.Context, *Packet) error
