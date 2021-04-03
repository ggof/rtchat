package hub

import (
	"github.com/gofiber/websocket/v2"
)

type Message struct {
	Ts     int64    `json:"ts"`
	Id     string   `json:"id"`
	Msg    string   `json:"msg"`
	SentBy string   `json:"sent_by"`
	RecvBy []string `json:"recv_by"`
}

type Invite struct {
	By string
	ID string
}

// Hub interface representing every possible action on the hub.
type Hub interface {
	Join(id string, conn *websocket.Conn) chan error
	Remove(id string, conn *websocket.Conn) chan error
	Rename(id, name string) chan error
	Send(id string, msg Message) chan error
	Receive(id string, msg Message) chan error
	Login(id string) chan error
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

// Add adds id as a user in this hub. if added by an admin, conn will be nil
// Users that aren't admins can only join public hubs
func (h *hub) Join(id string, conn *websocket.Conn) chan error {
	err := make(chan error)

	h.mailbox <- func() {
		if _, exists := h.users[id]; exists {
			returnError(err, ErrDuplicate)
			return
		}

		if conn != nil && !h.public {
			returnError(err, ErrPrivateHub)
			return
		}

		h.users[id] = entry{name: id, conn: conn}
		close(err)
	}

	return err
}

// Remove deletes id from this hub, closing it's connection if it is still opened.
func (h *hub) Remove(id string, conn *websocket.Conn) chan error {
	err := make(chan error)

	h.mailbox <- func() {
		user, ok := h.users[id]

		if !ok {
			returnError(err, ErrNotInHub)
			return
		}

		user.Close()
		delete(h.users, id)
		close(err)
	}

	return err
}

// Rename changes the nickname of id to name for this hub
func (h *hub) Rename(id, name string) chan error {
	err := make(chan error)

	h.mailbox <- func() {
		user, ok := h.users[id]

		if !ok {
			returnError(err, ErrNotInHub)
			return
		}

		user.name = name
		h.users[id] = user

		close(err)
	}

	return err
}

func (h *hub) Send(id string, msg Message) chan error {
	err := make(chan error)

	h.mailbox <- func() {
		for userId, user := range h.users {
			if id == userId {
				continue
			}

			user.Send(&msg)
		}
	}

	return err
}

func (h *hub) Receive(id string, msg Message) chan error {
	err := make(chan error)

	return err
}

func (h *hub) Login(id string) chan error {
	err := make(chan error)

	return err
}

func (h *hub) Logout(id string) chan error {
	err := make(chan error)

	return err
}

func (h *hub) Run() {
	for mess := range h.mailbox {
		mess()
	}
}
