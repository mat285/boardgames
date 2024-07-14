package websockets

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Conn struct {
	readLock  sync.Mutex
	writeLock sync.Mutex
	conn      *websocket.Conn
}

func (c *Conn) Write(p Packet) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(p.Type, p.Data)
}

func (c *Conn) Read() (int, []byte, error) {
	c.readLock.Lock()
	defer c.readLock.Unlock()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	return c.conn.ReadMessage()
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
