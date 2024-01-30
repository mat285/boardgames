package v1alpha1

import (
	"github.com/blend/go-sdk/uuid"
)

type Player struct {
	ID       uuid.UUID
	Username string
}

// type BasePlayer struct {
// 	ID       uuid.UUID
// 	Username string
// }

// type Connection interface {
// 	Sender
// 	Receiver
// }

// type Sender interface {
// 	Send(context.Context, Message) error
// 	Request(context.Context, Message) (*Message, error)
// }

// type Receiver interface {
// 	Receive(context.Context, MessageHandler) error
// }

// func NewBasePlayer(id uuid.UUID, username string) BasePlayer {
// 	return BasePlayer{
// 		ID:       id,
// 		Username: username,
// 	}
// }

// func (bp BasePlayer) GetID() uuid.UUID {
// 	return bp.ID
// }

// func (bp BasePlayer) GetUsername() string {
// 	return bp.Username
// }
