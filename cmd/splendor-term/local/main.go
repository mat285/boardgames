package main

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/games/splendor"
	"github.com/mat285/boardgames/games/splendor/pkg/client/terminal"
	"github.com/mat285/boardgames/games/splendor/pkg/game"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
	local "github.com/mat285/boardgames/server/local/v1alpha1"
)

func main() {

	ctx := context.Background()

	s := getLocalServer()
	g := getGame()
	players := getPlayers(g, s)

	e, err := s.NewEngine(ctx, g, nil)
	if err != nil {
		panic(err)
	}

	for _, player := range players {
		err := player.Connect(context.Background(), nil)
		if err != nil {
			panic(err)
		}
		err = player.Join(context.Background(), e.ID)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(s.StartEngine(context.Background(), e.ID))
}

func getGame() v1alpha1.Game {
	return splendor.NewGameWithConfig(game.Config{VictoryPoints: 7})
}

func getPlayers(g v1alpha1.Game, s *local.Server) []connection.Client {
	p := terminal.NewTerminalPlayer("user", g, getLocalClient(uuid.V4(), s))
	// b := bot.NewBot("bot", g, getLocalClient(s), bot.NewRandom())
	go p.Start(context.Background())
	// go b.Start(context.Background())
	return []connection.Client{
		p,
		// b,
	}
}

func getLocalClient(id uuid.UUID, s *local.Server) *local.Client {
	return local.NewClient(id, s)
}

func getLocalServer() *local.Server {
	return local.NewServer()
}
