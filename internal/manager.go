package network

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ServerVR() *fiber.App {

	app := fiber.New(fiber.Config{
		ServerHeader: "ServerUnimeta",
		AppName:      "Cbtic",
	})

	// Configurar CORS correctamente
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Permite todas las conexiones
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Initialize default config (Assign the middleware to /metrics)
	// app.Get("/metrics", monitor.New())
	// app.Get("/metrics", monitor.New(monitor.Config{Title: "Monitor ServerCbtic", Refresh: 1000}))
	app.Use("/ws", Middleware)
	app.Get("/ws/:id", websocket.New(HandleConnection))

	return app
}
