package game

import (
	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/games/splendor/pkg/items"
)

type Player struct {
	ID   uuid.UUID
	Hand items.Hand
}

func NewPlayer(id uuid.UUID) Player {
	return Player{
		ID:   id,
		Hand: items.NewHand(),
	}
}
