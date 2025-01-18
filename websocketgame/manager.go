package websocketgame

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func ServerVR() *fiber.App {

	app := fiber.New(fiber.Config{
		ServerHeader: "ServerUnimeta",
		AppName:      "Cbtic",
	})

	app.Use("/ws", Middleware)
	app.Get("/ws/:id", websocket.New(HandleConnection))

	return app
}
