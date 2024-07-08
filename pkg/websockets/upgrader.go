package websockets

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(_ *http.Request) bool { return true },
}

func Upgrader() *websocket.Upgrader {
	return upgrader
}

func UpgradeRequestToClient(rw http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	// Upgrade conn to websockets
	return Upgrader().Upgrade(rw, r, nil)
}
