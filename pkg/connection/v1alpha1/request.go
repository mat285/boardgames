package v1alpha1

import (
	"context"
	"fmt"
	"sync"
	"time"

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
	fmt.Println("req manager got packet", packet.ID, reqID)
	if len(reqID) == 0 {
		return rm.upstream(ctx, packet)
	}
	rm.Lock()
	defer rm.Unlock()
	if resp, has := rm.responses[reqID]; has && resp.p != nil && len(resp.p) == 0 {
		fmt.Println("rm sent packet on channel", packet.ID)
		resp.p <- packet
		rm.closeResponseChannelUnsafe(reqID)
		return nil
	}
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
	fmt.Println("fm sending packet to cli")
	err := sender.Send(ctx, packet)
	if err != nil {
		fmt.Println("req manager exiting", err)
		rm.closeResponseChannel(key)
		return nil, err
	}

	fmt.Println("rm listening on channel")
	tick := time.After(time.Second * 300)
	select {
	case <-ctx.Done():
		rm.closeResponseChannel(key)
		return nil, ctx.Err()
	case <-tick:
		rm.closeResponseChannel(key)
		return nil, fmt.Errorf("Timeout")
	case packet := <-resp.p:
		fmt.Println("req manager got response packet", packet.ID)
		return &packet, nil
	}
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
