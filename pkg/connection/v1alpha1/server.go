package v1alpha1

import (
	"github.com/blend/go-sdk/uuid"
)

// type Server interface {
// 	Receiver
// 	Connect(context.Context, ClientInfo) (uuid.UUID, error)
// }

type ClientInfo struct {
	ID       uuid.UUID
	Username string
	Sender   Sender
}
