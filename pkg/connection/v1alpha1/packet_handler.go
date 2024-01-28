package v1alpha1

import "context"

type PacketHandler func(context.Context, Packet) error
