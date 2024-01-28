package v1alpha1

import "github.com/blend/go-sdk/uuid"

type GameData struct {
	EngineID   uuid.UUID
	APIVersion string
	Data       []byte
}
