package v1alpha1

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/uuid"
	"github.com/gorilla/websocket"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	websockets "github.com/mat285/boardgames/pkg/websockets"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

type Websocket struct {
	lock sync.Mutex

	client  *websockets.Client
	inbound chan websockets.Packet

	stop chan struct{}
}

func NewWebsocket(id uuid.UUID, username string, conn *websocket.Conn) *Websocket {
	in := make(chan websockets.Packet, 32)
	return &Websocket{
		lock:    sync.Mutex{},
		client:  websockets.NewClient(id, username, conn, in),
		inbound: in,
	}
}

func (w *Websocket) Open(ctx context.Context) error {
	w.lock.Lock()
	if w.client == nil {
		w.lock.Unlock()
		return fmt.Errorf("No websocket connection")
	}
	c := w.client
	w.lock.Unlock()
	return c.Start(ctx)
}

func (w *Websocket) Send(ctx context.Context, m wire.Packet) error {
	w.lock.Lock()
	client := w.client
	w.lock.Unlock()
	if client == nil {
		return fmt.Errorf("No websocket connection")
	}
	p, err := w.serializeMessage(m)
	if err != nil {
		return err
	}
	return client.Send(ctx, *p)
}

func (w *Websocket) Listen(ctx context.Context, handle connection.PacketHandler) error {
	log := logger.GetLogger(ctx)
	w.lock.Lock()
	if w.stop != nil {
		w.lock.Unlock()
		return fmt.Errorf("Already listening")
	}

	if w.client == nil {
		w.lock.Unlock()
		return fmt.Errorf("No websocket connection")
	}

	stop := make(chan struct{})
	w.stop = stop
	defer func() {
		w.lock.Lock()
		if w.stop == stop {
			w.stop = nil
		}
		w.lock.Unlock()
	}()

	// watch the underlying connection and make sure we
	// exit if it does
	wsExit := make(chan error, 1)
	go func() {
		err := w.client.Start(ctx)
		wsExit <- err
		close(wsExit)
	}()

	w.lock.Unlock()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-w.stop:
			// graceful stop
			return nil
		case err := <-wsExit:
			logger.MaybeErrorContext(ctx, log, err)
			// ws connection broke
			return err
		case p, ok := <-w.inbound:
			if !ok {
				// channel was closed
				return fmt.Errorf("Websocket channel closed")
			}
			m, err := w.deserializePacket(p)
			if err != nil {
				logger.MaybeErrorContext(ctx, log, err)
				// bad packet
				continue
			}
			go w.handle(ctx, *m, handle)
		}
	}
}

func (w *Websocket) Close(ctx context.Context) error {
	w.lock.Lock()
	if w.stop == nil {
		w.lock.Unlock()
		return nil
	}
	close(w.stop)
	if w.client != nil {
		err := w.client.Stop(ctx)
		w.lock.Unlock()
		return err
	}
	w.lock.Unlock()
	return nil
}

func (w *Websocket) handle(ctx context.Context, p wire.Packet, h connection.PacketHandler) {
	log := logger.GetLogger(ctx)
	timeout, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	err := h(timeout, p)
	if err != nil {
		logger.MaybeErrorContext(ctx, log, err)
	}
}

func (w *Websocket) serializeMessage(m wire.Packet) (*websockets.Packet, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return &websockets.Packet{Type: websocket.TextMessage, Data: data}, nil
}

func (w *Websocket) deserializePacket(p websockets.Packet) (*wire.Packet, error) {
	var m wire.Packet
	return &m, json.Unmarshal(p.Data, &m)
}
