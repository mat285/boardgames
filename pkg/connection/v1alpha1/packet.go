package v1alpha1

import "encoding/json"

type Packet struct {
	Header  Header
	Payload []byte
}

func DeserializePacket(raw []byte) (*Packet, error) {
	var p Packet
	return &p, json.Unmarshal(raw, &p)
}

func (p Packet) Serialize() ([]byte, error) {
	return json.Marshal(p)
}
