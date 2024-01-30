package v1alpha1

import (
	"encoding/json"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

type Provider struct {
	game.Serializer
}

func NewProvider(s game.Serializer) Provider {
	return Provider{
		Serializer: s,
	}
}

func (mp Provider) NewPacket(t wire.PacketType, payload interface{}) (*wire.Packet, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	packet := wire.NewPacket(wire.OptPacketType(t), wire.OptPacketPayload(bytes))
	return &packet, nil
}

func (mp Provider) MessagePlayerMoveInfo(id uuid.UUID, move game.Move) (*wire.Packet, error) {
	so, err := mp.SerializeMove(move)
	if err != nil {
		return nil, err
	}
	data := MessageBodyPlayerMoveInfo{
		Player: id,
		Move:   so,
	}
	return mp.NewPacket(PacketTypePlayerMoveInfo, data)
}

func (mp Provider) ExtractPlayerMoveInfo(packet wire.Packet) (*MessageBodyPlayerMoveInfo, error) {
	var data MessageBodyPlayerMoveInfo
	err := json.Unmarshal(packet.Payload, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (mp Provider) MessagePlayerMove(move game.Move) (*wire.Packet, error) {
	so, err := mp.SerializeMove(move)
	if err != nil {
		return nil, err
	}
	return mp.NewPacket(PacketTypePlayerMove, so)
}

func (mp Provider) MessageGameOver(winners []uuid.UUID) (*wire.Packet, error) {
	return mp.NewPacket(PacketTypeGameOver, MessageBodyWinners(winners))
}

func (mp Provider) MessageRequestMove(state game.StateData) (*wire.Packet, error) {
	so, err := mp.SerializeState(state)
	if err != nil {
		return nil, err
	}
	return mp.NewPacket(PacketTypeRequestMove, so)
}

func (mp Provider) ExtractMove(packet wire.Packet) (game.Move, error) {
	if packet.Type != PacketTypePlayerMove {
		return nil, fmt.Errorf("Wrong Packet Type")
	}
	var so game.SerializedObject
	err := json.Unmarshal(packet.Payload, &so)
	if err != nil {
		return nil, err
	}
	return mp.DeserializeMove(&so)
}

func (mp Provider) ExtractState(packet wire.Packet) (game.StateData, error) {
	var so game.SerializedObject
	err := json.Unmarshal(packet.Payload, &so)
	if err != nil {
		return nil, err
	}
	return mp.DeserializeState(&so)
}
