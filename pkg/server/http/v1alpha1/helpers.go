package v1alpha1

import (
	engine "github.com/mat285/boardgames/pkg/engine/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
	model "github.com/mat285/boardgames/pkg/model/v1alpha1"
)

func GameFromEngine(e *engine.Engine) *model.Game {
	if e == nil {
		return nil
	}
	return &model.Game{
		ID:      e.ID,
		Game:    e.Game.Name(),
		Players: PlayersFromPlayers(e.State.Players),
	}
}

func PlayersFromPlayers(players []v1alpha1.Player) []model.Player {
	ret := make([]model.Player, len(players))
	for i := range players {
		ret[i] = model.Player{
			Username: players[i].Username,
		}
	}
	return ret
}
