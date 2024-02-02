package game

import (
	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/games/splendor/pkg/items"
	common "github.com/mat285/boardgames/pkg/common/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type Player struct {
	v1alpha1.Player
	Hand items.Hand
}

func NewPlayer(id uuid.UUID) Player {
	return Player{
		Player: v1alpha1.Player{
			ID: id,
		},
		Hand: items.NewHand(),
	}
}

func ToCommonPlayerSlice(players ...Player) []common.Player {
	ret := make([]common.Player, len(players))
	for i := range players {
		ret[i] = players[i]
	}
	return ret
}
