package games

import (
	"sort"

	"github.com/mat285/boardgames/games/splendor"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type RegisteredGame struct {
	Name   string
	New    func(interface{}) (v1alpha1.Game, error)
	Config func() interface{}
}

func RegisteredGames() map[string]RegisteredGame {
	return map[string]RegisteredGame{
		splendor.Name: {
			Name:   splendor.Name,
			New:    splendor.New,
			Config: splendor.NewConfig,
		},
	}
}

func ListGames() []string {
	rgs := RegisteredGames()
	names := make([]string, 0, len(rgs))
	for name := range rgs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
