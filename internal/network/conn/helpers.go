package conn

import (
	"encoding/json"
	"log"

	"github.com/DuvanRozoParra/try_server/config"
	"github.com/gofiber/contrib/websocket"
)

func shardWorker(s *Shard) {
	for {
		select {
		case msg := <-s.highPriorityChan:
			// start := time.Now()
			processMessage(s, msg)
			// log.Printf("[Shard %d] Mensaje HIGH procesado en %v | Evento: %v",
			// 	0, time.Since(start).Nanoseconds(), msg.Event)
		case msg := <-s.mediumPriorityChan:
			processMessage(s, msg)
		case msg := <-s.lowPriorityChan:
			processMessage(s, msg)
		case cmd := <-s.commandChan:
			handleCommand(s, cmd)
		}
	}
}

func processMessage(s *Shard, msg MessageObject) {
	s.mu.RLock()
	player, exists := s.players[msg.From]
	if !exists {
		log.Printf("Jugador no encontrado: %s", msg.From)
		return
	}
	s.mu.RUnlock()

	switch msg.Event {
	case config.MovePlayer:
		handleMovement(s, player, msg.Data)

	case config.RayInteraction:
		handleRayInteraction(s, player, msg.Data)

	case config.ActionHandsPlayer:
		handleActionsHandsAnimation(s, player.ID, msg.Data)

	default:
		log.Printf("Evento desconocido: %v", msg.Event)
	}
}

func handleCommand(s *Shard, cmd Command) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch cmd.Type {
	case "add":
		if cmd.Id != "" && cmd.Player != nil {
			s.players[cmd.Id] = cmd.Player
			s.connections[cmd.Id] = cmd.Wb
			// log.Printf("[Shard %d] Jugador ADDED: %s | Total: %d",
			// 	0, cmd.Id, len(s.players))
		}

	case "remove":
		if cmd.Id != "" {

			delete(s.players, cmd.Id)
			delete(s.connections, cmd.Id)
			aver := MessageObject{
				Data:     "",
				From:     cmd.Id,
				Priority: config.Low,
				Event:    config.DeletePlayer,
			}

			jsonBytes, err := json.Marshal(aver)
			if err != nil {
				log.Fatal("Error al convertir a JSON:", err)
			}
			log.Printf("REMOVE PLAYER")
			broadcastUpdate(s, jsonBytes)
			// log.Printf("[Shard %d] Jugador REMOVED: %s | Restantes: %d",
			// 	0, cmd.Id, len(s.players))
		}

	case "broadcast":
		broadcastUpdate(s, cmd.Message)
		// log.Printf("[Shard %d] Broadcast enviado a %d jugadores",
		// 	0, len(s.players))
	}
}

func broadcastUpdate(s *Shard, data []byte) {
	//s.mu.Lock()
	for _, c := range s.connections {
		err := c.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			log.Printf("Error escribiendo en WebSocket: %v", err)
		}
	}
	//s.mu.Unlock()
}
