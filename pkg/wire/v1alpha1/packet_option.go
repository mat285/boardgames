package v1alpha1

import "github.com/blend/go-sdk/uuid"

type PacketOption func(*Packet)

func OptPacketType(t PacketType) PacketOption {
	return func(p *Packet) {
		p.Header.Type = t
	}
}

func OptPacketPayload(data []byte) PacketOption {
	return func(p *Packet) {
		p.Payload = data
	}
}

func OptPacketOrigin(id uuid.UUID) PacketOption {
	return func(p *Packet) {
		p.Origin = id
	}
}

func OptPacketDestination(id uuid.UUID) PacketOption {
	return func(p *Packet) {
		p.Destination = id
	}
}

func OptPacketReference(id uuid.UUID) PacketOption {
	return func(p *Packet) {
		p.Reference = id
	}
}

func OptPacketHeaderValue(k, v string) PacketOption {
	return func(p *Packet) {
		p.Header.Options.Add(k, v)
	}
}
