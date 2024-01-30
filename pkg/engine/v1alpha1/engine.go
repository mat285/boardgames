package v1alpha1

import (
	"context"
	"fmt"
	"sync"

	"github.com/blend/go-sdk/uuid"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
	messages "github.com/mat285/boardgames/pkg/messages/v1alpha1"
	persist "github.com/mat285/boardgames/pkg/persist/v1alpha1"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

var (
// _ connection.Server = new(Engine)
)

type Engine struct {
	sync.Mutex
	ID uuid.UUID

	started bool

	Players map[string]*Player

	request connection.Requester

	MessageProvider messages.Provider

	State *game.State
	Game  game.Game

	Persist persist.Interface

	interrupt chan Event
	stop      chan struct{}
}

func NewEngine(g game.Game) *Engine {
	e := &Engine{
		ID:              uuid.V4(),
		Players:         make(map[string]*Player),
		MessageProvider: messages.NewProvider(g.Serializer()),
		Game:            g,
		interrupt:       make(chan Event),
		stop:            make(chan struct{}),
	}
	e.request = connection.NewRequestManager(e.receive)
	e.State = game.NewState(e.GamePlayers())
	return e
}

func (e *Engine) Join(ctx context.Context, client connection.ClientInfo) error {
	if e.started {
		return fmt.Errorf("Game Already Started")
	}
	e.Lock()
	defer e.Unlock()
	player := NewPlayer(client.ID, client.Username, client.Sender)
	e.Players[player.ID.ToFullString()] = player
	return nil
}

func (e *Engine) Receive(ctx context.Context, packet wire.Packet) error {
	fmt.Println("engine got packet", packet.ID)
	return e.request.Receive(ctx, packet)
}

func (e *Engine) receive(ctx context.Context, packet wire.Packet) error {

	return nil
}

func (e *Engine) Start(ctx context.Context) error {
	e.Lock()
	data, err := e.Game.Initialize(e.PlayerIDs())
	if err != nil {
		e.Unlock()
		return err
	}
	e.State.Data = data
	e.started = true
	e.Unlock()
	return e.gameLoop(ctx)
}

func (e *Engine) gameLoop(ctx context.Context) error {
	if e.stop == nil {
		e.stop = make(chan struct{})
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-e.stop:
			return nil
		case inter := <-e.interrupt:
			err := e.handleInterrupt(ctx, inter)
			if err != nil {
				return err
			}
		default:
		}

		if e.State.Data.IsDone() {
			msg, err := e.MessageProvider.MessageGameOver(e.State.Data.Winners())
			if err != nil {
				return err
			}
			return e.Broadcast(ctx, msg)
		}

		pid, err := e.State.Data.CurrentPlayer()
		if err != nil {
			fmt.Println(err)
			continue
		}

		player := e.GetPlayer(pid)
		if player == nil {
			fmt.Println("No player for id", pid)
			continue
		}

		msg, err := e.MessageProvider.MessageRequestMove(e.State.Data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("engine sending move request", msg.ID)
		resp, err := e.request.Request(ctx, player, *msg)
		if err != nil {
			fmt.Println(err)
			continue
		}

		move, err := e.MessageProvider.ExtractMove(*resp)
		if err != nil {
			fmt.Println(err)
			continue
		}

		response, err := move.Apply(e.State.Data)
		if err != nil {
			fmt.Println(err)
			player.Send(ctx, wire.ErrorPacket(err))
			continue
		}

		if !response.Valid {
			fmt.Println("Invalid Move")
			continue
		}

		msg, err = e.MessageProvider.MessagePlayerMoveInfo(player.ID, move)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = e.Broadcast(ctx, msg, player.ID)
		if err != nil {
			fmt.Println(err)
			continue
		}

		e.State.Data = response.State
	}
}

func (e *Engine) handleInterrupt(ctx context.Context, event Event) error {
	switch event.Type {
	case EventTypeUnknown:
		fmt.Printf("Unknown interrupt %v\n", event.Body)
		return nil
	case EventTypeStop:
		return e.Stop(ctx, event.Body)
	case EventTypeSave:
		err := e.Save(ctx)
		if err != nil {
			fmt.Println(err)
		}
		return e.Stop(ctx, event.Body)
	}
	return nil
}

func (e *Engine) Stop(ctx context.Context, optional ...interface{}) error {
	close(e.stop)
	// var body []byte
	// if len(optional) > 0 && optional[0] != nil {
	// 	body, _ = json.Marshal(optional[0])
	// }

	// msg := game.Message{
	// 	Type: game.MessageTypeGameStopped,
	// 	Data: body,
	// }
	// e.Broadcast(ctx, msg)
	return nil
}

func (e *Engine) Save(ctx context.Context) error {
	if e.Persist == nil {
		return nil
	}
	obj := persist.Object{
		Meta: persist.Meta{
			ID:            e.ID,
			APIVersion:    APIVersion,
			ObjectVersion: e.State.Version,
		},
		Data: e.State,
	}
	_, err := e.Persist.CheckAndSet(ctx, obj)
	return err
}

func (e *Engine) Broadcast(ctx context.Context, packet *wire.Packet, exclude ...uuid.UUID) error {
	if packet == nil {
		return nil
	}
	for _, player := range e.Players {
		if excludeUUID(player.ID, exclude...) {
			continue
		}
		err := player.Send(ctx, *packet)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (e *Engine) GetPlayer(id uuid.UUID) *Player {
	for i := range e.Players {
		if e.Players[i].ID.Equal(id) {
			return e.Players[i]
		}
	}
	return nil
}

func (e *Engine) PlayerIDs() []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(e.State.Players))
	for i := range e.Players {
		id, err := uuid.Parse(i)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

func (e *Engine) GamePlayers() []game.Player {
	players := make([]game.Player, 0, len(e.Players))
	for i := range e.Players {
		players = append(players, e.Players[i].Player)
	}
	return players
}
