package hub

import (
	"github.com/gofiber/websocket/v2"
)

type entry struct {
	name string
	conn *websocket.Conn
}

func (e *entry) Send(value string) {
	if e.conn != nil {
		e.conn.WriteMessage(websocket.TextMessage, []byte(value))
	}
}

func (e *entry) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}
