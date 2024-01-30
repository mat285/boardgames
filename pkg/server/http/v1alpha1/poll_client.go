package v1alpha1

import (
	"context"
	"fmt"
	"sync"

	"github.com/blend/go-sdk/uuid"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

type PollClient struct {
	sync.Mutex
	ID      uuid.UUID
	Packets chan wire.Packet
}

func NewPollClient(id uuid.UUID) *PollClient {
	return &PollClient{
		ID:      id,
		Packets: make(chan wire.Packet, 10),
	}
}

func (pc *PollClient) Send(ctx context.Context, packet wire.Packet) error {
	pc.Lock()
	if len(pc.Packets) == cap(pc.Packets) {
		pc.Unlock()
		return fmt.Errorf("Client disconnected")
	}
	pc.Packets <- packet
	pc.Unlock()
	return nil
}

func (pc *PollClient) Poll(ctx context.Context) *wire.Packet {
	select {
	case <-ctx.Done():
		return nil
	case p, ok := <-pc.Packets:
		if !ok {
			return nil
		}
		return &p
	}
}
