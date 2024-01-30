package v1alpha1

type PacketType uint64

const (
	PacketTypeUnknown = 0

	PacketTypeError = 100

	PacketTypeGameData = 1000
)
