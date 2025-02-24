package network

import (
	"encoding/json"
	"log"

	"github.com/DuvanRozoParra/try_server/config"
	"github.com/DuvanRozoParra/try_server/internal/game/players"
	"github.com/DuvanRozoParra/try_server/internal/network/conn"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func Middleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func HandleConnection(c *websocket.Conn) {
	playerId := c.Params("id")

	conn.ManagerShading.AddPlayer(playerId, players.NewPlayer(playerId), c)
	log.Printf("PLAYER CONNECT %s", playerId)

	defer func() {
		// Eliminar la conexión al desconectarse
		conn.ManagerShading.RemovePlayer(playerId)
		log.Printf("PLAYER DISCONNECT %s", playerId)
		c.Close()
	}()

	var (
		//mt  int
		msg []byte
		err error
	)

	for {
		if _, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}

		var message conn.MessageObject
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Error al deserializar el mensaje:", err)
			continue
		}

		// log.Printf("msg.Event => %+v\n", message.Event == 3)

		conn.ManagerShading.EnquequeMessage(message)
	}
}

func getMessagePriority(event config.EventServer) config.MessagePriority {
	switch event {
	case config.MovePlayer:
		return config.High // Prioridad alta
	case config.RayInteraction:
		return config.Medium // Prioridad media
	default:
		return config.Low // Prioridad baja
	}
}
