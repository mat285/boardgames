package v1alpha1

import (
	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/pkg/apiversions"
)

type Header struct {
	Metadata
	Options HeaderOptions
}

type Metadata struct {
	ID         uuid.UUID
	Request    uuid.UUID
	Type       PacketType
	APIVersion string
}

type HeaderOptions map[string]string

func (ho *HeaderOptions) Append(other HeaderOptions) {
	if other == nil {
		return
	}
	for k, v := range other {
		(*ho)[k] = v
	}
}

func NewHeader() Header {
	return Header{
		Metadata: Metadata{
			ID:         uuid.V4(),
			APIVersion: apiversions.V1Alpha1,
		},
		Options: make(HeaderOptions),
	}
}

type PacketType uint64

const (
	PacketTypeUnknown = 0

	PacketTypeClientData = 1000
)

func MinDataPacketType() PacketType {
	return PacketTypeClientData
}
