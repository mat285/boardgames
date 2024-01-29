package v1alpha1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/blend/go-sdk/uuid"
)

type Message struct {
	Type   MessageType
	Sender uuid.UUID
	Data   []byte
}

type MessageHandler func(context.Context, Message) (*Message, error)

func NewMessage(t MessageType, data *SerializedObject) (*Message, error) {
	bytes, err := data.Serialize()
	if err != nil {
		return nil, err
	}
	return &Message{
		Type: t,
		Data: bytes,
	}, nil
}

func (m Message) DeserializeToObject() (*SerializedObject, error) {
	var so SerializedObject
	return &so, json.Unmarshal(m.Data, &so)
}

// DataForMessageType returns an empty object to deserialize the message into
func DataForMessageType(t MessageType) interface{} {
	switch t {
	case MessageTypeRequestMove:
		return &MoveRequest{}
	}
	return map[string]interface{}{}
}

func MessageFromError(err error) Message {
	return Message{
		Type: MessageTypeError,
		Data: []byte(err.Error()),
	}
}

func MessagePlayerMoveInfo(id uuid.UUID, move Move) (*Message, error) {
	data := map[string]interface{}{
		"player": id,
		"move":   move,
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &Message{
		Type: MessageTypePlayerMoveInfo,
		Data: bytes,
	}, nil
}

func MessageGameOver(winners []uuid.UUID) (*Message, error) {
	serial, err := json.Marshal(winners)
	if err != nil {
		return nil, err
	}
	data := SerializedObject{Data: serial}
	bytes, err := data.Serialize()
	if err != nil {
		return nil, err
	}
	return &Message{
		Type: MessageTypeGameOver,
		Data: bytes,
	}, nil
}

func MessageRequestMove(state StateData, serializer Serializer) (*Message, error) {
	so, err := serializer.SerializeState(state)
	if err != nil {
		return nil, err
	}
	bytes, err := so.Serialize()
	if err != nil {
		return nil, err
	}
	return &Message{
		Type: MessageTypeRequestMove,
		Data: bytes,
	}, nil
}

// func MessagePlayerMove(move Move) (*Message, error) {
// 	return NewMessage(MessageTypePlayerMove, move)
// }

// func MoveRequestFromMessage1(message Message, provider EmptyProvider) (StateData, error) {
// 	if message.Type != MessageTypeRequestMove {
// 		return nil, fmt.Errorf("Wrong Message Type")
// 	}
// 	res, err := message.DeserializeMessageBody(provider)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req, ok := res.(StateData)
// 	if !ok {
// 		return nil, fmt.Errorf("No move request")
// 	}
// 	return req, nil
// }

func MoveFromMessage(message Message, serializer Serializer) (Move, error) {
	if message.Type != MessageTypePlayerMove {
		return nil, fmt.Errorf("Wrong Message Type")
	}
	res, err := message.DeserializeMessageBody(serializer)
	if err != nil {
		return nil, err
	}
	move, ok := res.(Move)
	if !ok {
		return nil, fmt.Errorf("No move request")
	}
	return move, nil
}

func (m Message) DeserializeMessageBody(serializer Serializer) (interface{}, error) {
	// switch m.Type {
	// case MessageTypeRequestMove:
	// 	serializer.DeserializeState(error)

	// case MessageTypePlayerMove:
	// 	move := provider.EmptyMove()
	// 	err := move.Deserialize(m.Data)
	// 	return move, err
	// }
	return nil, nil
}

func (m Message) String() string {
	bytes, _ := json.Marshal(m.Data)
	return string(bytes)
}

type MessageType uint64

const (
	MessageTypeUnknown MessageType = 0
	MessageTypeError   MessageType = 1

	MessageTypePlayerMoveInfo MessageType = 101
	MessageTypeStateUpdate    MessageType = 102
	MessageTypeGameOver       MessageType = 103
	MessageTypeGameStopped    MessageType = 104

	MessageTypeRequestMove MessageType = 201
	MessageTypePlayerMove  MessageType = 202
)
