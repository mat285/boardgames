package v1alpha1

import (
	"github.com/blend/go-sdk/uuid"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

const (
	PacketTypePlayerMoveInfo wire.PacketType = wire.PacketTypeGameData + 101
	PacketTypeStateUpdate    wire.PacketType = wire.PacketTypeGameData + 102
	PacketTypeGameOver       wire.PacketType = wire.PacketTypeGameData + 103
	PacketTypeGameStopped    wire.PacketType = wire.PacketTypeGameData + 104

	PacketTypeRequestMove wire.PacketType = wire.PacketTypeGameData + 201
	PacketTypePlayerMove  wire.PacketType = wire.PacketTypeGameData + 202
)

type MessageBodyPlayerMoveInfo struct {
	Player uuid.UUID
	Move   *game.SerializedObject
}

type MessageBodyWinners []uuid.UUID
