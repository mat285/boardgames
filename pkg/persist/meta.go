package persist

import "github.com/blend/go-sdk/uuid"

type Meta struct {
	ID            uuid.UUID
	APIVersion    string
	ObjectVersion uint64
}
