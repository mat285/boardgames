package v1alpha1

import (
	"context"

	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

type PacketHandler func(context.Context, wire.Packet) error
