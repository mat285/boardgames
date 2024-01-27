package v1alpha1

import (
	"context"
	"sync"
	"time"

	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type Connection struct {
	sync.Mutex
	Player v1alpha1.Player

	drain     sync.WaitGroup
	connected bool
	wait      chan struct{}
	timer     *time.Timer

	alert chan *Connection
}

func (c *Connection) Connected() bool {
	return c.connected
}

func (c *Connection) Reconnect(ctx context.Context) {
	c.Lock()
	defer c.Unlock()
	close(c.wait)
	if c.timer != nil {
		c.timer.Stop()
	}
	c.connected = true
}

func (c *Connection) Disconnect(ctx context.Context) {
	c.Lock()
	defer c.Unlock()
	c.drain.Wait()
	c.connected = false
	c.wait = make(chan struct{})
	c.timer = time.NewTimer(time.Minute)
	c.alert <- c
}

func (c *Connection) RequestMove(ctx context.Context, req v1alpha1.MoveRequest) (v1alpha1.Move, error) {
	err := c.WaitOrTimeout(ctx)
	if err != nil {
		return nil, err
	}

	// do the request and handle response
	return nil, nil
}

func (c *Connection) WaitOrTimeout(ctx context.Context) error {
	c.Lock()
	if c.connected {
		c.Unlock()
		return nil
	}
	c.Unlock()
	c.drain.Add(1)
	defer c.drain.Done()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.wait:
			return nil
		case <-c.timer.C:
			c.timer.Stop()
			return ErrTimeout
		}
	}
}
