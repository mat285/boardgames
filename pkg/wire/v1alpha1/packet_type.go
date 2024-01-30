package v1alpha1

type PacketType uint64

const (
	PacketTypeUnknown    = 0
	PacketTypeByteStream = 1
	PacketTypeString     = 2

	PacketTypeError = 100

	PacketTypeGameData = 1000
)
