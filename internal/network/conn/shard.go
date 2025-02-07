package conn

import (
	"log"
	"sync"
	"time"

	"github.com/DuvanRozoParra/try_server/config"
	"github.com/DuvanRozoParra/try_server/internal/game/players"
	"github.com/gofiber/contrib/websocket"
)

type Shard struct {
	sm                 *ShardManager
	connections        map[string]*websocket.Conn
	players            map[string]*players.Players
	highPriorityChan   chan MessageObject
	mediumPriorityChan chan MessageObject
	lowPriorityChan    chan MessageObject
	commandChan        chan Command
	mu                 sync.RWMutex
}

func (s *Shard) Enqueue(msg MessageObject, priority config.MessagePriority) {
	switch priority {
	case config.High:
		s.highPriorityChan <- msg
		// fmt.Printf("Aver msg => %+v\n", msg)
		// log.Printf("[Shard %d] Encolado HIGH | Cola: %d",
		// 	0, len(s.highPriorityChan))
	case config.Medium:
		s.mediumPriorityChan <- msg
		// log.Printf("[Shard %d] Encolado MEDIUM | Cola: %d",
		// 	0, len(s.mediumPriorityChan))
	default:
		s.lowPriorityChan <- msg
		// log.Printf("[Shard %d] Encolado LOW | Cola: %d",
		// 	0, len(s.lowPriorityChan))
	}
}

func (s *Shard) StartMetricsLogger() {
	go func() {
		for {
			time.Sleep(5 * time.Second) // Log cada 30 segundos
			s.mu.RLock()
			log.Printf("[Shard %d] Métricas | Jugadores: %d | Colas: H=%d, M=%d, L=%d",
				0,
				len(s.players),
				len(s.highPriorityChan),
				len(s.mediumPriorityChan),
				len(s.lowPriorityChan),
			)
			s.mu.RUnlock()
		}
	}()
}
