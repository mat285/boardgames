package v1alpha1

import (
	"encoding/json"

	"github.com/mat285/boardgames/pkg/websockets"
)

type Packet struct {
	Header
	Payload []byte
}

func FromWebsocket(p websockets.Packet) Packet {
	wp := Packet{
		Header:  NewHeader(),
		Payload: p.Data,
	}
	wp.Header.Type = PacketType(p.Type)
	return wp
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

func (p Packet) MustJSON() string {
	data, _ := json.MarshalIndent(p, "", " ")
	return string(data)
}

func (p Packet) Websocket() websockets.Packet {
	return websockets.Packet{
		Type: int(p.Header.Type),
		Data: p.Payload,
	}
}

func ErrorPacket(err error) Packet {
	return NewPacket(
		OptPacketType(PacketTypeError),
		OptPacketPayload([]byte(err.Error())),
	)
}
