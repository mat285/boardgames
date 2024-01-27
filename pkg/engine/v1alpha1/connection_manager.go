package v1alpha1

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"time"

// 	"github.com/mat285/boardgames/pkg/game/v1alpha1"
// )

// type ConnectionManager struct {
// 	Engine *Engine

// 	startLock   sync.Mutex
// 	started     bool
// 	startEngine chan struct{}

// 	connectionLock sync.Mutex
// 	connections    map[string]*Connection

// 	waitLock       sync.Mutex
// 	waiting        bool
// 	connectionWait chan struct{}
// 	timeoutWait    time.Timer

// 	connect    chan *Connection
// 	disconnect chan *Connection

// 	runLock sync.Mutex
// 	running bool
// }

// func NewConnectionManager(g *Engine) *ConnectionManager {
// 	return &ConnectionManager{
// 		Engine:      g,
// 		connections: make(map[string]*Connection),
// 	}
// }

// func (c *ConnectionManager) Start(ctx context.Context) error {
// 	if c.running {
// 		return fmt.Errorf("Engine already running")
// 	}

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		case p, ok := <-c.connect:
// 			if !ok {
// 				return fmt.Errorf("Engine closed")
// 			}
// 			c.connectionLock.Lock()
// 			c.connections[p.Player.ID.ToFullString()] = p
// 			c.connectionLock.Unlock()
// 			c.checkAllReady(ctx)

// 		case p := <-c.disconnect:
// 			go c.wait(ctx, p)
// 		}
// 	}
// }

// func (c *ConnectionManager) Ready(ctx context.Context) error {
// 	if !c.running {
// 		return fmt.Errorf("Manager not running")
// 	}

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		case <-c.startEngine:
// 			return nil
// 		}
// 	}
// }

// func (c *ConnectionManager) checkAllReady(ctx context.Context) {
// 	c.connectionLock.Lock()
// 	defer c.connectionLock.Unlock()
// 	for _, p := range c.Engine.Players() {
// 		if c, has := c.connections[p.ID.ToFullString()]; !has || !c.Connected() {
// 			return
// 		}
// 	}

// 	close(c.startEngine)
// }

// func (c *ConnectionManager) wait(ctx context.Context, conn *Connection) {
// 	err := c.pause(ctx)
// 	if err != nil {

// 	}
// 	err = conn.WaitOrTimeout(ctx)
// 	if err != nil {
// 		c.timeout(ctx)
// 	}
// 	c.resume(ctx)
// }

// func (c *ConnectionManager) timeout(ctx context.Context) error {
// 	// TODO disconnect everyone and save Engine state
// 	return nil
// }

// func (c *ConnectionManager) pause(ctx context.Context) error {
// 	// TODO send message to all clients we are waiting with a timer
// 	return nil
// }

// func (c *ConnectionManager) resume(ctx context.Context) error {
// 	// TODO send message to all clients we are waiting with a timer
// 	return nil
// }

// func (c *ConnectionManager) RequestMove(ctx context.Context, player *v1alpha1.Player, req v1alpha1.MoveRequest) (v1alpha1.Move, error) {
// 	if !c.running {
// 		return nil, fmt.Errorf("Manager not running")
// 	}

// 	c.connectionLock.Lock()
// 	conn := c.connections[player.ID.ToFullString()]
// 	if conn == nil {
// 		return nil, fmt.Errorf("Player missing from pool")
// 	}

// 	move, err := conn.RequestMove(ctx, req)
// 	if IsError(err, ErrTimeout) {
// 		// return c.Engine.pause(ctx)
// 	}
// 	return move, err
// }
