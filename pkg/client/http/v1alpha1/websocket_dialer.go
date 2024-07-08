package v1alpha1

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/blend/go-sdk/uuid"
	"github.com/gorilla/websocket"
	connection "github.com/mat285/boardgames/pkg/connection/v1alpha1"
	"github.com/mat285/boardgames/pkg/errors"
)

type WebsocketDialer struct {
	*Websocket
	lock     sync.Mutex
	Addr     string
	Username string
	UserID   uuid.UUID
}

func NewWebsocketDialer(addr string, userID uuid.UUID, username string) *WebsocketDialer {
	return &WebsocketDialer{
		Addr:     addr,
		Username: username,
		UserID:   userID,
	}
}

func (w *WebsocketDialer) Listen(ctx context.Context, handler connection.PacketHandler) error {
	return w.listen(ctx, handler)
}

func (w *WebsocketDialer) ListenRetry(ctx context.Context, handler connection.PacketHandler, attempts int) error {
	return w.retryListen(ctx, handler, attempts)
}

func (w *WebsocketDialer) listen(ctx context.Context, handler connection.PacketHandler) error {
	w.lock.Lock()
	if w.Websocket != nil {
		w.lock.Unlock()
		return fmt.Errorf("already listening")
	}
	err := w.dial(ctx)
	if err != nil {
		w.lock.Unlock()
		return err
	}
	wg := errors.NewErrorWaitGroup(2)
	wg.Add(2)
	go func() {
		wg.PushDone(w.Websocket.Open(ctx))
	}()
	go func() {
		wg.PushDone(w.Websocket.Listen(ctx, handler))
	}()
	w.lock.Unlock()
	err = wg.Wait()
	w.lock.Lock()
	if w.Websocket != nil {
		w.Websocket.Close(ctx)
	}
	w.Websocket = nil
	w.lock.Unlock()
	return err
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
	if w.Websocket == nil {
		return nil
	}
	err := w.Websocket.Close(ctx)
	if err != nil {
		return err
	}
	w.Websocket = nil
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
	w.Websocket = NewWebsocket(uuid.V4(), w.Username, conn)
	return nil
}
