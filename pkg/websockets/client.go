package websockets

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2 * 1024 * 1024
)

type Client struct {
	id       uuid.UUID
	username string
	lock     sync.Mutex

	conn *Conn

	outbound chan *Packet
	inbound  chan Packet

	rwCtx    context.Context
	rwCancel context.CancelFunc
	rwWG     *sync.WaitGroup
	running  bool
	stopped  bool
}

type Deregister func(context.Context, *Client) error

type Packet struct {
	Type int
	Data []byte
}

func NewClient(id uuid.UUID, username string, conn *websocket.Conn, c chan Packet) *Client {
	return &Client{
		id:       id,
		username: username,
		conn:     &Conn{conn: conn},
		outbound: make(chan *Packet, 32),
		inbound:  c,
	}
}

func (c *Client) GetID() uuid.UUID {
	return c.id
}

func (c *Client) GetUsername() string {
	return c.username
}

func (c *Client) Send(ctx context.Context, p Packet) error {
	if !c.isOpen() {
		return fmt.Errorf("Closed connection")
	}

	c.outbound <- &p
	return nil
}

func (c *Client) isOpen() bool {
	return c.getConn() != nil
}

func (c *Client) getConn() *Conn {
	return c.conn
}

func (c *Client) Start(ctx context.Context) error {
	log := logger.GetLogger(ctx)
	c.lock.Lock()
	if c.running {
		c.lock.Unlock()
		return fmt.Errorf("already running ws client")
	}
	if c.conn == nil {
		c.lock.Unlock()
		return fmt.Errorf("no websocket connection")
	}
	logger.MaybeDebugfContext(ctx, log, "Starting websockets client")

	c.rwCtx, c.rwCancel = context.WithCancel(ctx)
	errs := make(chan error)

	cancel := c.rwCancel
	var wg sync.WaitGroup
	wg.Add(3)
	go func(wg *sync.WaitGroup) {
		defer recover()
		defer wg.Done()
		errs <- c.pingLoop(c.rwCtx)
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer recover()
		defer wg.Done()
		errs <- c.write(c.rwCtx)
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer recover()
		defer wg.Done()
		errs <- c.read(c.rwCtx)
	}(&wg)

	c.rwWG = &wg
	c.lock.Unlock()
	errStrs := make([]string, 0, 3)
	err := <-errs
	cancel()
	wg.Wait()
	if err != nil {
		errStrs = append(errStrs, err.Error())
	}
	err = <-errs
	if err != nil {
		errStrs = append(errStrs, err.Error())
	}
	err = <-errs
	if err != nil {
		errStrs = append(errStrs, err.Error())
	}
	if len(errStrs) > 0 {
		return fmt.Errorf(strings.Join(errStrs, "\n"))
	}
	return nil
}

func (c *Client) Stop(ctx context.Context) error {
	c.lock.Lock()
	if !c.running {
		c.lock.Unlock()
		return nil
	}
	wg := c.rwWG
	conn := c.conn
	c.conn = nil
	cancel := c.rwCancel
	c.rwCancel = nil
	cancel()
	if conn != nil {
		return conn.Close()
	}
	c.lock.Unlock()
	if wg != nil {
		wg.Wait()
	}
	c.lock.Lock()
	c.stopped = true
	c.running = false
	c.lock.Unlock()
	return nil
}

func (c *Client) pingLoop(ctx context.Context) error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case _, ok := <-ticker.C:
			if !ok {
				return fmt.Errorf("ping ticker stopped")
			}
			if len(c.outbound) != 0 {
				continue
			}
			err := c.Send(ctx, Packet{Type: websocket.PingMessage})
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (c *Client) writeMessage(packet Packet) error {
	conn := c.getConn()
	if conn == nil {
		return fmt.Errorf("closed websocket connection")
	}
	err := conn.Write(packet)
	if err != nil {
		return err
	}
	if packet.Type == websocket.CloseMessage {
		return fmt.Errorf("closed websocket connection")
	}
	return nil
}

func (c *Client) write(ctx context.Context) error {
	for {
		select {
		case packet, ok := <-c.outbound:
			if !ok {
				return c.writeMessage(Packet{Type: websocket.CloseMessage, Data: []byte{}})
			}
			if packet == nil {
				continue
			}
			// write the packet
			err := c.writeMessage(*packet)
			if err != nil {
				return err
			}
			continue

		case <-ctx.Done():
			c.writeMessage(Packet{Type: websocket.CloseMessage, Data: []byte{}})
			return ctx.Err()
		}
	}
}

func (c *Client) readMessage(ctx context.Context) (int, []byte, error) {
	conn := c.getConn()
	if conn == nil {
		return -1, nil, fmt.Errorf("closed websocket connection")
	}
	return conn.Read()
}

func (c *Client) read(ctx context.Context) error {
	log := logger.GetLogger(ctx)

	conn := c.getConn()
	if conn == nil {
		return fmt.Errorf("no connection open to read")
	}

	for {
		// check ctx first then fall through
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		t, rawData, err := c.readMessage(ctx)
		if err != nil {
			logger.MaybeErrorContext(ctx, log, err)
			return err
		}

		switch t {
		case websocket.TextMessage, websocket.BinaryMessage:
			err := c.handleClientMessage(ctx, Packet{Type: t, Data: rawData})
			if err != nil {
				logger.MaybeErrorContext(ctx, log, err)
				// drop packet
				continue
			}
		case websocket.PingMessage:
			err := c.Send(ctx, Packet{Type: websocket.PongMessage})
			if err != nil {
				logger.MaybeErrorContext(ctx, log, err)
				continue
			}
		case websocket.PongMessage:
			continue
		default:
			err := fmt.Errorf("Unknown message type %d", t)
			if err != nil {
				logger.MaybeErrorContext(ctx, log, err)
				continue
			}
		}

	}
}

func (c *Client) handleClientMessage(ctx context.Context, p Packet) error {
	timeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	select {
	case <-timeout.Done():
		return timeout.Err()
	case c.inbound <- p:
		return nil
	}
}
