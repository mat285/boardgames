package v1alpha1

type PacketType uint64

const (
	PacketTypeUnknown    = 0
	PacketTypeByteStream = 1
	PacketTypeString     = 2

	PacketTypeError = 100

	PacketTypeAPI = 1000

	PacketTypeGameData = 100000
)
