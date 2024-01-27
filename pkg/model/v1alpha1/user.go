package v1alpha1

import "github.com/blend/go-sdk/uuid"

type User struct {
	ID       uuid.UUID
	Username string
}
