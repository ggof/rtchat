package hub

import "github.com/gofiber/websocket/v2"

type entry struct {
	name string
	conn *websocket.Conn
}

func (e *entry) Send(msg *Message) {
	if e.conn != nil {
		e.conn.WriteJSON(msg)
	}
}

func (e *entry) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}
