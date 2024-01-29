package v1alpha1

import "github.com/blend/go-sdk/uuid"

type PacketOption func(*Packet)

func OptPacketType(t PacketType) PacketOption {
	return func(p *Packet) {
		p.Header.Type = t
	}
}

func OptPacketRequestID(id uuid.UUID) PacketOption {
	return func(p *Packet) {
		p.Header.Request = id
	}
}

func OptPacketPayload(data []byte) PacketOption {
	return func(p *Packet) {
		p.Payload = data
	}
}
