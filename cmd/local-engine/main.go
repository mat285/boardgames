package main

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/games/splendor"
	localconnection "github.com/mat285/boardgames/pkg/connection/local/v1alpha1"
	engine "github.com/mat285/boardgames/pkg/engine/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

func main() {

	e := getEngine()
	fmt.Println(e.Start(context.Background()))
}

func getEngine() *engine.E2 {
	return engine.NewE2(getState(), getGame())
}

func getGame() v1alpha1.Game {
	return new(splendor.Game)
}

func getState() *v1alpha1.State {
	return v1alpha1.NewState(getPlayers())
}

func getPlayers() []v1alpha1.Player {
	p := localconnection.NewTerminalPlayer()
	go p.Run()
	return []v1alpha1.Player{v1alpha1.NewPlayer(uuid.V4(), p)}
}
