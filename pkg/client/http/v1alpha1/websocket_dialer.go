package v1alpha1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/blend/go-sdk/uuid"
	"github.com/gorilla/websocket"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	"github.com/mat285/boardgames/pkg/websockets"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

type WebsocketDialer struct {
	ws       *websockets.Client
	stop     chan struct{}
	inbound  chan websockets.Packet
	lock     sync.Mutex
	Addr     string
	Username string
	UserID   uuid.UUID

	listening bool
}

func NewWebsocketDialer(addr string, userID uuid.UUID, username string) *WebsocketDialer {
	return &WebsocketDialer{
		Addr:     addr,
		Username: username,
		UserID:   userID,
		inbound:  make(chan websockets.Packet, 16),
	}
}

func (w *WebsocketDialer) Open(ctx context.Context) error {
	return nil
}

func (w *WebsocketDialer) Send(ctx context.Context, m wire.Packet) error {
	w.lock.Lock()
	client := w.ws
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

func (w *WebsocketDialer) Listen(ctx context.Context, handler connection.PacketHandler) error {
	return w.listen(ctx, handler)
}

func (w *WebsocketDialer) ListenRetry(ctx context.Context, handler connection.PacketHandler, attempts int) error {
	return w.retryListen(ctx, handler, attempts)
}

func (w *WebsocketDialer) listen(ctx context.Context, handler connection.PacketHandler) error {
	if w.listening {
		return fmt.Errorf("already listening")
	}
	w.lock.Lock()
	if w.listening {
		w.lock.Unlock()
		return fmt.Errorf("already listening")
	}
	w.listening = true
	err := w.dial(ctx)
	if err != nil {
		return err
	}
	stop := make(chan struct{})
	w.stop = stop
	w.lock.Unlock()

	errs := make(chan error)
	go func() {
		errs <- w.ws.Start(ctx)
	}()

	defer func() {
		w.lock.Lock()
		defer w.lock.Unlock()
		w.ws.Stop(ctx)
		w.listening = false
		if w.stop != nil {
			w.stop = nil
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errs:
			return err
		case <-stop:
			return nil
		case packet, ok := <-w.inbound:
			if !ok {
				return fmt.Errorf("ws channel closed")
			}
			typed, err := w.deserializePacket(packet)
			if err != nil {
				continue
			}
			go func() {
				err := handler(ctx, *typed)
				if err != nil {

				}
			}()
			continue
		}
	}
}

func (w *WebsocketDialer) retryListen(ctx context.Context, handler connection.PacketHandler, attempts int) error {
	errs := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		err := w.listen(ctx, handler)
		if err == nil {
			return nil
		}
		errs++
		if errs >= attempts {
			return fmt.Errorf("retries exceeded")
		}
		time.Sleep(5 * time.Second)
	}
}

func (w *WebsocketDialer) Close(ctx context.Context) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	if !w.listening || w.stop == nil {
		return nil
	}
	close(w.stop)
	w.stop = nil
	return nil
}

func (w *WebsocketDialer) dial(ctx context.Context) error {
	header := http.Header{}
	conn, resp, err := websocket.DefaultDialer.DialContext(context.Background(), w.Addr, header)
	if err != nil && (resp == nil || resp.StatusCode != 307) {
		return err
	}

	if resp.StatusCode == 307 {
		loc := resp.Header.Get("Location")
		if len(loc) == 0 {
			return fmt.Errorf("ws dial: no redirct location")
		}

		loc = strings.Replace(loc, "http", "ws", 1)

		conn, _, err = websocket.DefaultDialer.DialContext(context.Background(), loc, header)
		if err != nil {
			return err
		}
	}
	ws := websockets.NewClient(w.UserID, w.Username, conn, w.inbound)
	w.ws = ws
	return nil
}

func (w *WebsocketDialer) deserializePacket(p websockets.Packet) (*wire.Packet, error) {
	var m wire.Packet
	return &m, json.Unmarshal(p.Data, &m)
}

func (w *WebsocketDialer) serializeMessage(m wire.Packet) (*websockets.Packet, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return &websockets.Packet{Type: websocket.TextMessage, Data: data}, nil
}
