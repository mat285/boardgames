package v1alpha1

import (
	"context"

	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

type rsPipe struct {
	r Receiver
	s Sender
}

func PipeReceiverToSender(r Receiver) Sender {
	return &rsPipe{
		r: r,
	}
}

func PipeSendToReceiver(s Sender) Receiver {
	return &rsPipe{
		s: s,
	}
}

func (rs *rsPipe) Send(ctx context.Context, packet wire.Packet) error {
	return rs.r.Receive(ctx, packet)
}

func (rs *rsPipe) Receive(ctx context.Context, packet wire.Packet) error {
	return rs.s.Send(ctx, packet)
}
