package hub

import (
	"time"

	"github.com/gofiber/websocket/v2"
)

// Hub interface representing every possible action on the hub.
type Hub interface {
	Join(by string, id string, conn *websocket.Conn) chan error
	Remove(by string, id string, conn *websocket.Conn) chan error
	Rename(by string, id string, name string) chan error
	Send(id string, msg string) chan error
	Read(id string, ts int64) chan error
	Login(id string, conn *websocket.Conn) chan error
	Logout(id string) chan error
	Run()
}

type event = func()

type hub struct {
	name    string           // Name of the hub.
	owner   string           // Owner of the hub. Only him can add / delete people
	users   map[string]entry // map of username to hub entry
	public  bool             // wether this hub is public or not
	mailbox chan event       // mailbox of the hub. Every event goes here
}

func NewHub(name, owner string) Hub {
	return &hub{
		name:    name,
		owner:   owner,
		users:   make(map[string]entry),
		mailbox: make(chan event, 16),
	}
}

// Add adds req.ID as a user in this hub.
// Users that aren't admins can only join public hubs
func (h *hub) Join(by string, id string, conn *websocket.Conn) chan error {
	err := make(chan error)

	h.mailbox <- func() {
		if _, exists := h.users[id]; exists {
			returnError(err, ErrDuplicate)
			return
		}

		if by != h.owner && !h.public {
			returnError(err, ErrPrivateHub)
			return
		}

		h.users[id] = entry{name: id, conn: conn}
		close(err)
	}

	return err
}

// Remove deletes id from this hub, closing it's connection if it is still opened.
// If conn is not null, then it must match
func (h *hub) Remove(by string, id string, conn *websocket.Conn) chan error {
	err := make(chan error)

	h.mailbox <- func() {
		user, ok := h.users[id]

		if !ok {
			returnError(err, ErrNotInHub)
			return
		}

		if by != h.owner && conn != user.conn {
			returnError(err, ErrNotAdmin)
			return
		}

		user.Close()
		delete(h.users, id)
		close(err)
	}

	return err
}

// Rename changes the nickname of id to name for this hub
func (h *hub) Rename(by string, id string, name string) chan error {
	err := make(chan error)
	defer close(err)
	return err
}

func (h *hub) Send(id, msg string) chan error {
	err := make(chan error)

	h.mailbox <- func() {
		for _, user := range h.users {
			user.Send(encodeSend(time.Now().Unix(), id, msg))
		}

		close(err)
	}

	return err
}

func (h *hub) Read(id string, ts int64) chan error {
	err := make(chan error)

	h.mailbox <- func() {
		for _, user := range h.users {
			user.Send(encodeRead(id, ts))
		}

		close(err)
	}

	return err
}

func (h *hub) Login(id string, conn *websocket.Conn) chan error {
	err := make(chan error)
	h.mailbox <- func() {
		user, ok := h.users[id]
		if !ok {
			returnError(err, ErrNotInHub)
			return
		}

		user.conn = conn
		h.users[id] = user
		close(err)
	}
	return err
}

func (h *hub) Logout(id string) chan error {
	err := make(chan error)
	h.mailbox <- func() {
		user, ok := h.users[id]
		
		if !ok {
			returnError(err, ErrNotInHub)
			return
		}

		user.conn = nil
		h.users[id] = user
		close(err)
	}

	return err
}

func (h *hub) Run() {
	for mess := range h.mailbox {
		mess()
	}
}

