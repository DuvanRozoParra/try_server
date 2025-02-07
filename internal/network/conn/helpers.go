package conn

import (
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
			// start := time.Now()
			processMessage(s, msg)
			// log.Printf("[Shard %d] Mensaje MEDIUM procesado en %v | Evento: %v",
			// 	0, time.Since(start), msg.Event)

		case msg := <-s.lowPriorityChan:
			// start := time.Now()
			processMessage(s, msg)
			// log.Printf("[Shard %d] Mensaje LOW procesado en %v | Evento: %v",
			// 	0, time.Since(start), msg.Event)

		case cmd := <-s.commandChan:
			// start := time.Now()
			handleCommand(s, cmd)
			// log.Printf("[Shard %d] Comando '%s' ejecutado en %v",
			// 	0, cmd.Type, time.Since(start))
		}
	}
}

func processMessage(s *Shard, msg MessageObject) {
	// 1. Validar jugador
	s.mu.RLock()
	player, exists := s.players[msg.From]
	if !exists {
		log.Printf("Jugador no encontrado: %s", msg.From)
		return
	}
	s.mu.RUnlock()

	// 2. Procesar según tipo de evento
	switch msg.Event {
	case config.MovePlayer:
		handleMovement(s, player, msg.Data)

	case config.RayInteraction:
		// handleRayInteraction(s, player, msg.Data)

	default:
		log.Printf("Evento desconocido: %v", msg.Event)
	}
	// 3. Métricas de prioridad
	//log.Printf("Mensaje procesado - Prioridad: %v, Evento: %v", priority, msg.Event)
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
	for _, c := range s.connections {
		err := c.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			log.Printf("Error escribiendo en WebSocket: %v", err)
		}
	}
}
