package websocketgame

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/DuvanRozoParra/try_server/config"
	"github.com/DuvanRozoParra/try_server/websocketgame/components/players"
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
	playerId := c.Params("id")

	log.Printf("PLAYER CONNECT %s", playerId)

	connections.Lock()
	connections.players[playerId] = c
	connections.Unlock()

	defer func() {
		// Eliminar la conexión al desconectarse
		connections.Lock()
		delete(connections.players, playerId)
		connections.Unlock()
		// Eliminar al jugador del manager
		players.Mp.Lock()
		delete(players.Mp.Players, playerId)
		players.Mp.Unlock()
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
			players.Mp.Lock()
			allPlayers, err := players.MovePlayer(playerId, message.Data)
			players.Mp.Unlock()

			message.Data = allPlayers
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error al serializar mensaje para %s: %v", playerId, err)
				break
			}

			if err := c.WriteMessage(mt, data); err != nil {
				log.Printf("Error al enviar mensaje a %s: %v", playerId, err)
			}

		default:
			// Si el evento no es válido
			log.Printf("Invalid event received: %+v from %s", message.Event, message.From)
			continue // Si quieres continuar con el siguiente mensaje
		}

	}

}
