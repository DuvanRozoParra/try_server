package websocketgame

import (
	"encoding/json"
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

var playersManage players.PlayersManage = players.PlayersManage{
	Players:      make(map[string]*players.Players),
	LimitPlayers: 5,
}

func Middleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("playersManage", &playersManage)
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func HandleConnection(c *websocket.Conn) {
	playerId := c.Params("id")
	manager, ok := c.Locals("playersManage").(*players.PlayersManage)
	if !ok {
		log.Println("Error: no se pudo obtener playersManage del contexto")
		c.Close()
		return
	}

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
		manager.Lock()
		delete(manager.Players, playerId)
		manager.Unlock()
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
			player, err := manager.ConvertToJson(message.Data)
			if err != nil {
				log.Printf("Error al convertir datos del jugador %s: %v", playerId, err)
				break
			}

			manager.Lock()
			exists, _ := manager.PlayerExists(playerId)
			if !exists {
				// Actualizar información del jugador existente
				log.Printf("✅ Agregando jugador: %s. Usuarios conectados: %d", playerId, len(manager.Players)+1)
				manager.Players[playerId] = player
			} else {
				// Agregar un nuevo jugador
				log.Printf("🔄 Actualizar jugador: %s.", playerId)
				manager.Players[playerId] = player
			}

			// Obtener información de otros jugadores (excluyendo al actual)
			players, err := manager.GetDataPlayers(playerId)
			if err != nil {
				log.Printf("Error al obtener jugadores para %s: %v", playerId, err)
				break
			}
			manager.Unlock()

			// Actualizar el mensaje con la lista de jugadores
			message.Data = players
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error al serializar mensaje para %s: %v", playerId, err)
				break
			}

			// Enviar mensaje a todos los jugadores excepto al actual
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

func Emits(connections *Connections, messageType int, data []byte, excludedID string) {
	for id, conn := range connections.players {
		if id != excludedID { // Excluir al jugador actual
			if err := conn.WriteMessage(messageType, data); err != nil {
				log.Printf("Error al enviar mensaje a %s: %v", id, err)
			}
		}
	}
}
