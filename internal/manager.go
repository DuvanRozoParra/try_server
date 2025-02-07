package network

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func ServerVR() *fiber.App {

	app := fiber.New(fiber.Config{
		ServerHeader: "ServerUnimeta",
		AppName:      "Cbtic",
	})

	// Initialize default config (Assign the middleware to /metrics)
	// app.Get("/metrics", monitor.New())
	// app.Get("/metrics", monitor.New(monitor.Config{Title: "Monitor ServerCbtic", Refresh: 1000}))
	app.Use("/ws", Middleware)
	app.Get("/ws/:id", websocket.New(HandleConnection))

	return app
}
