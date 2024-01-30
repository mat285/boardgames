package main

import (
	"context"
	"fmt"

	"github.com/mat285/boardgames/games/splendor"
	"github.com/mat285/boardgames/games/splendor/pkg/client/terminal"
	"github.com/mat285/boardgames/games/splendor/pkg/game"
	bot "github.com/mat285/boardgames/pkg/bot/v1alpha1"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
	local "github.com/mat285/boardgames/pkg/local/v1alpha1"
)

func main() {

	s := getLocalServer()
	g := getGame()
	players := getPlayers(g, s)

	e := s.NewEngine(g)

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
	p := terminal.NewTerminalPlayer("user", g, getLocalClient(s))
	b := bot.NewBot("bot", g, getLocalClient(s), bot.NewRandom())
	go p.Start(context.Background())
	go b.Start(context.Background())
	return []connection.Client{
		p,
		b,
	}
}

func getLocalClient(s *local.Server) *local.Client {
	return local.NewClient(s)
}

func getLocalServer() *local.Server {
	return local.NewServer()
}
