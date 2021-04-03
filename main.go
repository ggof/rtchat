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
	hub := hub.NewHub()

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

		err := <-h.Add(id, c)

		for {
			mess, err := ReadString(c)

			if err != nil {
				log.Println(err.Error())
				c.Close()
				h.Del(id, c)
				return
			}

			h.Snd(id, mess)
		}
	}
}

func ReadString(c *websocket.Conn) (string, error) {
			t, messBytes, err := c.ReadMessage()

			if err != nil {
				return "", errors.New("ws: " + err.Error())
			}

			if t != websocket.TextMessage {
				return "", errors.New("ws: not a text message")
			}

			return string(messBytes), nil
}
