package hub

import (
	"time"

	"github.com/gofiber/websocket/v2"
)

type usrMessage struct {
	id   string
	conn *websocket.Conn
}

type sndMessage struct {
	Ts  int64  `json:"ts"`
	Id  string `json:"id"`
	Msg string `json:"msg"`
}

type Hub interface {
	Add(id string, conn *websocket.Conn)
	Del(id string, conn *websocket.Conn)
	Snd(id, msg string)
	Run()
}

type hub struct {
	add   chan usrMessage
	del   chan usrMessage
	snd   chan sndMessage
	users map[string]*websocket.Conn
}

func NewHub() Hub {
	return &hub{
		add:   make(chan usrMessage, 16),
		del:   make(chan usrMessage, 16),
		snd:   make(chan sndMessage, 16),
		users: make(map[string]*websocket.Conn),
	}
}

func (h *hub) Add(id string, conn *websocket.Conn) {
	h.add <- usrMessage{id: id, conn: conn}
}

func (h *hub) Del(id string, conn *websocket.Conn) {
	h.del <- usrMessage{id: id, conn: conn}
}

func (h *hub) Snd(id, msg string) {
	h.snd <- sndMessage{
		Id:  id,
		Msg: msg,
		Ts:  time.Now().Unix(),
	}
}
func (h *hub) Run() {
	for {
		select {
		case mess := <-h.add:
			h.users[mess.id] = mess.conn
		case mess := <-h.del:
			delete(h.users, mess.id)
		case mess := <-h.snd:
			for id, conn := range h.users {
				if id == mess.Id {
					continue
				}

				conn.WriteJSON(mess)
			}
		}
	}
}
