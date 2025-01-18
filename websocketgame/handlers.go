package websocketgame

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/DuvanRozoParra/try_server/config"
	"github.com/DuvanRozoParra/try_server/players"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type MessageObject struct {
	Data  string             `json:"data"`
	From  string             `json:"from"`
	Event config.EventServer `json:"events"`
}

type Connections struct {
	sync.Mutex
	players map[string]*websocket.Conn
}

var connections Connections = Connections{
	players: make(map[string]*websocket.Conn),
}

func Middleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func HandleConnection(c *websocket.Conn) {
	var manager *players.PlayersManage = players.NewManagePlayers(5)
	playerId := c.Params("id")
	log.Printf("PLAYER CONNECT %s SUCCESS...", playerId)

	connections.Lock()
	connections.players[playerId] = c
	connections.Unlock()

	defer func() {
		// Eliminar la conexión al desconectarse
		connections.Lock()
		delete(connections.players, playerId)
		connections.Unlock()
		// Eliminar al jugador del manager
		if err := manager.RemovePlayer(playerId); err != nil {
			log.Printf("Error removing player %s: %v", playerId, err)
		}
		log.Printf("PLAYER DISCONNECT %s", playerId)
		c.Close()
	}()

	var (
		mt  int
		msg []byte
		err error
	)

	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}

		// Convertir mensaje a MessageObject
		var message MessageObject
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Error al deserializar el mensaje:", err)
			continue // Saltar al siguiente mensaje si hay error
		}

		switch message.Event {
		case config.RayInteraction:
			// Procesar RayInteraction
			log.Printf("Ray Interaction event received from %v => data: %+v\n", message.From, message)
			// Realiza las acciones relacionadas con RayInteraction

		case config.MovePlayer:
			// Procesar MovePlayer
			// log.Printf("✅ Move Player event received from %v => data: %+v\n", message.From, message)
			exists, _ := manager.PlayerExists(playerId)

			if exists {
				
				fmt.Printf("⚠️ El jugador ya existe: %+v\n", playerId)
				manager.ModifyPlayer(playerId, message.Data)
			} else {
				fmt.Printf("✅ Agregando nuevo jugador: %+v\n", playerId)
				err := manager.AddPlayer(message.Data)
				if err != nil {
					fmt.Printf("❌ Error al agregar jugador: %v\n", err)
				}
			}

			players, _ := manager.GetAllPlayers()
			message.Data = players
			data, err := json.Marshal(message)
			if err != nil {
				fmt.Println("Error al convertir a byte:", err)
				return
			}
			connections.Lock()
			Emits(&connections, mt, data, playerId)
			connections.Unlock()

		default:
			// Si el evento no es válido
			log.Printf("Invalid event received: %+v from %s", message.Event, message.From)
			continue // Si quieres continuar con el siguiente mensaje
		}
	}

}

func Emits(connections *Connections, mt int, msg []byte, playerID string) {
	for id, conn := range connections.players {
		if id != playerID {
			if err := conn.WriteMessage(mt, msg); err != nil {
				log.Printf("Error enviando mensaje a %s: %v", id, err)
			}
		}
	}
}
