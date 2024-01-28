package v1alpha1

import (
	"encoding/json"

	"github.com/blend/go-sdk/uuid"
)

type Message struct {
	Type MessageType
	Data []byte
}

func MessageFromError(err error) Message {
	return Message{
		Type: MessageTypeError,
		Data: []byte(err.Error()),
	}
}

func MessagePlayerMove(id uuid.UUID, move Move) (Message, error) {
	data := map[string]interface{}{
		"player": id,
		"move":   move,
	}
	bytes, err := json.Marshal(data)
	return Message{
		Type: MessageTypePlayerMove,
		Data: bytes,
	}, err
}

func MessageGameOver(winners []uuid.UUID) (Message, error) {
	data := map[string]interface{}{
		"winners": winners,
	}
	bytes, err := json.Marshal(data)
	return Message{
		Type: MessageTypePlayerMove,
		Data: bytes,
	}, err
}

func (m Message) String() string {
	return string(m.Data)
}

type MessageType int

const (
	MessageTypeUnknown     MessageType = 0
	MessageTypeError       MessageType = 1
	MessageTypePlayerMove  MessageType = 2
	MessageTypeStateUpdate MessageType = 3
	MessageTypeGameOver    MessageType = 4
	MessageTypeGameStopped MessageType = 5
)
