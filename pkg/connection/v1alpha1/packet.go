package v1alpha1

type Packet struct {
	Header  PacketHeader
	Payload []byte
}

type PacketHeader struct {
}
