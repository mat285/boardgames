package v1alpha1

import (
	"context"
	"fmt"
	"sync"
)

var (
	_ Interface = new(LocalConnection)
)

type LocalConnection struct {
	Interface

	handlerLock sync.Mutex
	handler     PacketHandler
	stop        chan struct{}
}

func NewLocalConnection() *LocalConnection {
	return &LocalConnection{}
}

func (lc *LocalConnection) Send(ctx context.Context, packet Packet) error {
	return lc.passthrough(ctx, packet)
}

func (lc *LocalConnection) Listen(ctx context.Context, handle PacketHandler) error {
	lc.handlerLock.Lock()
	if lc.handler != nil {
		lc.handlerLock.Unlock()
		return fmt.Errorf("Already listening")
	}
	lc.handler = handle
	lc.stop = make(chan struct{})
	stop := lc.stop
	lc.handlerLock.Unlock()
	select {
	case <-ctx.Done():
		lc.Stop()
		return ctx.Err()
	case _ = <-stop:
		lc.Stop()
		return nil
	}
}

func (lc *LocalConnection) Stop() {
	lc.handlerLock.Lock()
	defer lc.handlerLock.Unlock()
	if lc.stop != nil {
		close(lc.stop)
		lc.stop = nil
	}
	lc.handler = nil
}

func (lc *LocalConnection) passthrough(ctx context.Context, packet Packet) error {
	lc.handlerLock.Lock()
	defer lc.handlerLock.Unlock()
	if lc.handler == nil {
		return fmt.Errorf("Connection Refused")
	}
	return lc.handler(ctx, packet)
}
