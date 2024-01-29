package main

import (
	"context"
	"fmt"

	"github.com/mat285/boardgames/games/splendor"
	"github.com/mat285/boardgames/games/splendor/pkg/client/terminal"
	"github.com/mat285/boardgames/games/splendor/pkg/game"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	engine "github.com/mat285/boardgames/pkg/engine/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

func main() {

	e := getEngine()
	fmt.Println(e.Start(context.Background()))
}

func getEngine() *engine.Engine {
	g := getGame()
	return engine.NewEngine(getPlayers(g), g)
}

func getGame() v1alpha1.Game {
	return splendor.NewGameWithConfig(game.Config{VictoryPoints: 7})
}

func getPlayers(g v1alpha1.Game) []engine.Player {
	p := terminal.NewTerminalPlayer(g, connection.NewLocalConnection())
	go p.Run(context.Background())
	return []engine.Player{
		p,
		// bot.NewBot(uuid.V4(), "bot", connection.NewLocalConnection(), bot.NewRandom()),
	}
}
