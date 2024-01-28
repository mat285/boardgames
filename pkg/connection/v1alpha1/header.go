package v1alpha1

import "github.com/blend/go-sdk/uuid"

type Header struct {
	ID         uuid.UUID
	Type       PacketType
	APIVersion string
}

type PacketType uint64

const (
	PacketTypeUnknown = 0

	PacketTypeGameData = 1000
)
