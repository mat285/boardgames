package v1alpha1

import (
	"context"
	"fmt"
)

func PushPacket(ctx context.Context, c chan Packet, p Packet) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c <- p:
		return nil
	}
}

func PushPacketAsync(ctx context.Context, c chan Packet, p Packet) {
	go PushPacket(ctx, c, p)
}

func PollPacket(ctx context.Context, c chan Packet) (Packet, error) {
	select {
	case <-ctx.Done():
		return Packet{}, ctx.Err()
	case p, ok := <-c:
		if !ok {
			return Packet{}, fmt.Errorf("channel closed")
		}
		return p, nil
	}
}
