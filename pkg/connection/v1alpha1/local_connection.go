package v1alpha1

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// )

// var (
// 	_ Interface = new(Local)
// 	_ Server    = new(LocalConnector)
// )

// type LocalConnector struct {
// 	Server *Local
// 	Client *Local
// }

// type Local struct {
// 	Partner *Local

// 	handlerLock sync.Mutex
// 	handler     PacketHandler
// 	stop        chan struct{}
// }

// func NewLocalConnecter() *LocalConnector {
// 	lc := &LocalConnector{
// 		Server: new(Local),
// 		Client: new(Local),
// 	}
// 	lc.Server.Partner = lc.Client
// 	lc.Client.Partner = lc.Server
// 	return lc
// }

// func (lc *LocalConnector) Connect(ctx context.Context) (Interface, error) {
// 	return lc.Client, nil
// }

// func (lc *LocalConnector) Serve() Interface {
// 	return lc.Server
// }

// func (lc *Local) Send(ctx context.Context, packet Packet) error {
// 	return lc.Partner.Receive(ctx, packet)
// }

// func (lc *Local) Listen(ctx context.Context, handle PacketHandler) error {
// 	lc.handlerLock.Lock()
// 	if lc.handler != nil {
// 		lc.handlerLock.Unlock()
// 		return fmt.Errorf("Already listening")
// 	}
// 	lc.handler = handle
// 	lc.stop = make(chan struct{})
// 	stop := lc.stop
// 	lc.handlerLock.Unlock()
// 	select {
// 	case <-ctx.Done():
// 		lc.Stop()
// 		return ctx.Err()
// 	case _ = <-stop:
// 		lc.Stop()
// 		return nil
// 	}
// }

// func (lc *Local) Stop() {
// 	lc.handlerLock.Lock()
// 	defer lc.handlerLock.Unlock()
// 	if lc.stop != nil {
// 		close(lc.stop)
// 		lc.stop = nil
// 	}
// 	lc.handler = nil
// }

// func (lc *Local) Receive(ctx context.Context, packet Packet) error {
// 	lc.handlerLock.Lock()
// 	defer lc.handlerLock.Unlock()
// 	if lc.handler == nil {
// 		return fmt.Errorf("Connection Refused")
// 	}
// 	return lc.handler(ctx, packet)
// }
