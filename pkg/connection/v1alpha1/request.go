package v1alpha1

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blend/go-sdk/uuid"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

const (
	PacketHeaderRequestResponse = "Request-Response"
	PacketHeaderRequestID       = "Request-ID"
)

type Requester interface {
	Request(context.Context, Sender, wire.Packet) (*wire.Packet, error)
	Receiver
}

type RequestManager struct {
	sync.Mutex
	responses map[string]packetErrorChannelPair

	upstream PacketHandler
}

type packetErrorChannelPair struct {
	p chan wire.Packet
	e chan error
}

func newPacketErrorChannelPair() packetErrorChannelPair {
	return packetErrorChannelPair{
		p: make(chan wire.Packet, 1),
		e: make(chan error, 1),
	}
}

func NewRequestManager(upstream PacketHandler) Requester {
	return &RequestManager{
		Mutex:     sync.Mutex{},
		responses: make(map[string]packetErrorChannelPair),

		upstream: upstream,
	}
}

func (rm *RequestManager) Receive(ctx context.Context, packet wire.Packet) error {
	reqID := packet.Options.Value(PacketHeaderRequestID)
	parsed, err := uuid.Parse(reqID)
	if err != nil {
		return rm.upstream(ctx, packet)
	}
	reqID = parsed.ToFullString()
	if len(reqID) == 0 {
		return rm.upstream(ctx, packet)
	}
	rm.Lock()
	if resp, has := rm.responses[reqID]; has && resp.p != nil && len(resp.p) == 0 {
		rm.Unlock()
		resp.p <- packet
		return nil
	}
	rm.Unlock()
	return rm.upstream(ctx, packet)
}

func (rm *RequestManager) Request(ctx context.Context, sender Sender, packet wire.Packet) (*wire.Packet, error) {
	key := packet.ID.ToFullString()
	rm.Lock()
	if _, has := rm.responses[key]; has {
		rm.Unlock()
		return nil, fmt.Errorf("Already Waiting for Response for %s", key)
	}
	resp := newPacketErrorChannelPair()
	rm.responses[packet.ID.ToFullString()] = resp
	rm.Unlock()
	packet.Options.Add(PacketHeaderRequestResponse, packet.ID.ToFullString())
	err := sender.Send(ctx, packet)
	if err != nil {
		rm.removeResponseChannel(key)
		return nil, err
	}

	tick := time.After(time.Second * 300)
	select {
	case <-ctx.Done():
		rm.removeResponseChannel(key)
		return nil, ctx.Err()
	case <-tick:
		rm.removeResponseChannel(key)
		return nil, fmt.Errorf("Timeout")
	case packet := <-resp.p:
		rm.removeResponseChannel(key)
		return &packet, nil
	}
}

func (rm *RequestManager) removeResponseChannel(key string) {
	rm.Lock()
	defer rm.Unlock()
	rm.closeResponseChannelUnsafe(key)
	delete(rm.responses, key)
}

func (rm *RequestManager) closeResponseChannel(key string) {
	rm.Lock()
	defer rm.Unlock()
	rm.closeResponseChannelUnsafe(key)
}

func (rm *RequestManager) closeResponseChannelUnsafe(key string) {
	if c, has := rm.responses[key]; has {
		if c.p != nil {
			close(c.p)
		}
		if c.e != nil {
			close(c.e)
		}
		delete(rm.responses, key)
	}
}
