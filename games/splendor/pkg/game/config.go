package game

import "github.com/mat285/boardgames/games/splendor/pkg/items"

type Config struct {
	StartingPlayer int
	VictoryPoints  int
}

func StandardConfig() Config {
	return Config{
		StartingPlayer: 0,
		VictoryPoints:  items.StandardVictoryPoints,
	}
}
