package v1alpha1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/games"
	api "github.com/mat285/boardgames/pkg/server/api/v1alpha1"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

func (s *Server) Receive(ctx context.Context, packet wire.Packet) error {
	if packet.Origin.IsZero() {
		return nil // drop
	}
	if packet.Origin.Equal(s.ID) {
		return nil // drop
	}
	switch packet.Type {
	case api.PacketTypeListGamesRequest:
		resp := api.ListGamesResponse{Games: games.ListGames()}
		return s.respondJSON(ctx, packet.ID, packet.Origin, api.PacketTypeListGamesResponse, resp)
	case api.PacketTypeNewGameRequest:
	case api.PacketTypeJoinGameRequest:
	default:
		return fmt.Errorf("Unsupported packet type %d", packet.Type)
	}
	return nil
}

func (s *Server) respondJSON(ctx context.Context, ref, dst uuid.UUID, t wire.PacketType, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	packet := s.newResponsePacket(ref, dst, t, data)
	// have the router send the package to the right client
	return s.Router.Receive(ctx, packet)
}

func (s *Server) newResponsePacket(ref, dst uuid.UUID, t wire.PacketType, body []byte) wire.Packet {
	return wire.NewPacket(
		wire.OptPacketOrigin(s.ID),
		wire.OptPacketDestination(dst),
		wire.OptPacketReference(ref),
		wire.OptPacketPayload(body),
		wire.OptPacketType(t),
	)
}
