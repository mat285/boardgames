package v1alpha1

import "github.com/blend/go-sdk/uuid"

type Game struct {
	ID      uuid.UUID
	Game    string
	Players []Player
	Status  GameStatus
}

type Player struct {
	Username string
}

type GameStatus string
