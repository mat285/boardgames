package main

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/games/splendor"
	"github.com/mat285/boardgames/games/splendor/pkg/client/terminal"
	"github.com/mat285/boardgames/games/splendor/pkg/game"
	bot "github.com/mat285/boardgames/pkg/bot/v1alpha1"
	engine "github.com/mat285/boardgames/pkg/engine/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

func main() {

	e := getEngine()
	fmt.Println(e.Start(context.Background()))
}

func getEngine() *engine.Engine {
	return engine.NewEngine(getState(), getGame())
}

func getGame() v1alpha1.Game {
	return splendor.New(game.Config{VictoryPoints: 7})
}

func getState() *v1alpha1.State {
	return v1alpha1.NewState(getPlayers())
}

func getPlayers() []v1alpha1.Player {
	p := terminal.NewTerminalPlayer()
	go p.Run(context.Background())
	return []v1alpha1.Player{
		v1alpha1.NewPlayer(uuid.V4(), p),
		v1alpha1.NewPlayer(uuid.V4(), bot.NewBot(bot.NewRandom())),
	}
}
