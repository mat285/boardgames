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

// var (
// 	_ connection.ServerInfo = new(Engine)
// )

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

	stop chan struct{}
}

func NewEngine(g game.Game) *Engine {
	e := &Engine{
		ID:              uuid.V4(),
		Players:         make(map[string]*Player),
		MessageProvider: messages.NewProvider(g),
		Game:            g,
		stop:            make(chan struct{}),
	}
	e.request = connection.NewRequestManager(e.receive)
	e.State = game.NewState(e.GamePlayers())
	return e
}

func (e *Engine) GetID() uuid.UUID {
	return e.ID
}

func (e *Engine) SetID(id uuid.UUID) {
	e.ID = id
}

func (e *Engine) Join(ctx context.Context, client connection.ClientInfo) error {
	if e.started {
		return fmt.Errorf("Game Already Started")
	}
	e.Lock()
	defer e.Unlock()
	player := NewPlayer(client.GetID(), client.GetUsername(), client)
	e.Players[player.ID.ToFullString()] = player
	return nil
}

func (e *Engine) Receive(ctx context.Context, packet wire.Packet) error {
	return e.request.Receive(ctx, packet)
}

func (e *Engine) ReceiveFromPlayer(ctx context.Context, id uuid.UUID, packet wire.Packet) error {
	packet.Header.Values().Add(connection.PacketHeaderRequestID, id.String())
	return e.request.Receive(ctx, packet)
}

func (e *Engine) receive(ctx context.Context, packet wire.Packet) error {

	return nil
}

func (e *Engine) Start(ctx context.Context) error {
	e.Lock()
	if e.started {
		e.Unlock()
		return fmt.Errorf("Game already started")
	}
	data, err := e.Game.Initialize(e.PlayerIDs())
	if err != nil {
		e.Unlock()
		return err
	}
	e.State.Data = data
	e.started = true
	e.stop = make(chan struct{})
	e.Unlock()
	return e.gameLoop(ctx)
}

func (e *Engine) gameLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-e.stop:
			return nil
		default:
		}

		e.Lock()
		err := e.gameTurnUnsafe(ctx)
		e.Unlock()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (e *Engine) gameTurnUnsafe(ctx context.Context) error {
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
		return err
	}

	player := e.GetPlayer(pid)
	if player == nil {
		return fmt.Errorf("No player for id %s", pid)
	}

	msg, err := e.MessageProvider.MessageRequestMove(e.State.Data)
	if err != nil {
		return err
	}
	msg.Destination = pid
	msg.Origin = e.ID
	msg.ID = pid

	fmt.Println("sending move request to player", pid)
	resp, err := e.request.Request(ctx, player, *msg)
	if err != nil {
		fmt.Println("error getting moves from player", pid, err)
		return err
	}
	fmt.Println("Got move from player", pid)
	move, err := e.MessageProvider.ExtractMove(*resp)
	if err != nil {
		return err
	}

	response, err := move.Apply(e.State.Data)
	if err != nil {
		player.Send(ctx, wire.ErrorPacket(err))
		return err
	}

	if !response.Valid {
		player.Send(ctx, wire.ErrorPacket(fmt.Errorf("Invalid Move")))
		return fmt.Errorf("Invalid Move")
	}

	msg, err = e.MessageProvider.MessagePlayerMoveInfo(player.ID, move)
	if err != nil {
		return err
	}
	err = e.Broadcast(ctx, msg, player.ID)
	if err != nil {
		return err
	}

	e.State.Data = response.State
	return nil
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
