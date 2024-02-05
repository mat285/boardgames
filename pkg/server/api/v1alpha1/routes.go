package v1alpha1

import (
	"github.com/mat285/boardgames/pkg/apiversions"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

const (
	RouteBase = "/api/" + apiversions.V1Alpha1
)

const (
	PacketTypeListGamesRequest  wire.PacketType = wire.PacketTypeAPI + 1
	PacketTypeListGamesResponse wire.PacketType = wire.PacketTypeAPI + 2
	PacketTypeNewGameRequest    wire.PacketType = wire.PacketTypeAPI + 3
	PacketTypeNewGameResponse   wire.PacketType = wire.PacketTypeAPI + 4
	PacketTypeJoinGameRequest   wire.PacketType = wire.PacketTypeAPI + 5
	PacketTypeJoinGameResponse  wire.PacketType = wire.PacketTypeAPI + 6
)
