package v1alpha1

import "github.com/blend/go-sdk/uuid"

type ListGamesResponse struct {
	Games []string
}

type NewGameRequest struct {
	Name   string
	Config []byte
}

type GameResponse struct {
	ID uuid.UUID
}
