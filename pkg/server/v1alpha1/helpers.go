package v1alpha1

import (
	engine "github.com/mat285/boardgames/pkg/engine/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

func GameFromEngine(e *engine.Engine) *Game {
	if e == nil {
		return nil
	}
	return &Game{
		ID:      e.ID,
		Game:    e.Game.Name(),
		Players: PlayersFromPlayers(e.State.Players),
	}
}

func PlayersFromPlayers(players []v1alpha1.Player) []Player {
	ret := make([]Player, len(players))
	for i := range players {
		ret[i] = Player{
			Username: players[i].Username,
		}
	}
	return ret
}
