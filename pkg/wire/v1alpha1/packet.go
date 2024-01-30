package v1alpha1

import "encoding/json"

type Packet struct {
	Header
	Payload []byte
}

func NewPacket(opts ...PacketOption) Packet {
	p := Packet{
		Header: NewHeader(),
	}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

func DeserializePacket(raw []byte) (*Packet, error) {
	var p Packet
	return &p, json.Unmarshal(raw, &p)
}

func (p Packet) Serialize() ([]byte, error) {
	return json.Marshal(p)
}

func ErrorPacket(err error) Packet {
	return NewPacket(
		OptPacketType(PacketTypeError),
		OptPacketPayload([]byte(err.Error())),
	)
}
