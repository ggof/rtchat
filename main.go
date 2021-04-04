package main

import (
	"errors"
	"log"

	"github.com/ggof/rtchat/hub"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type WSHandler = func(*websocket.Conn)

func main() {
	hub := hub.NewHub("test", "gofg2301")

	r := fiber.New()

	r.Use("/ws", upgrade)

	r.Get("/ws/:id", websocket.New(handleWS(hub)))

	go hub.Run()

	r.Listen(":5000")
}

func upgrade(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	return c.Next()
}

func handleWS(h hub.Hub) WSHandler {
	return func(c *websocket.Conn) {
		id := c.Params("id")

		if err := <-h.Login(id, c); err != nil {
			log.Println(err.Error())
			c.Close()
		}

		log.Printf("New connection from id %s\n", id)

		for {
			mess, err := decode(c)

			if err != nil {
				log.Println(err.Error())
				h.Logout(id)
				return
			}

			mess.Handle(h)
		}
	}
}

func decode(c *websocket.Conn) (hub.Message, error) {
	t, msgBytes, err := c.ReadMessage()

	if err != nil {
		return nil, errors.New("ws: " + err.Error())
	}

	if t != websocket.TextMessage {
		return nil, errors.New("ws: not a text message")
	}

	for _, tryDecode := range hub.DefaultDecoders {
		if msg := tryDecode(string(msgBytes)); msg != nil {
			return msg, nil
		}
	}

	return nil, errors.New("encodable: no decoders to decode message")
}
