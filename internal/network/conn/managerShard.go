package conn

import (
	"encoding/json"
	"hash/fnv"
	"log"

	"github.com/DuvanRozoParra/try_server/config"
	"github.com/DuvanRozoParra/try_server/internal/game/players"
	"github.com/gofiber/contrib/websocket"
)

const (
	Shards            = 1   // Número de shards
	MsgBufferPerShard = 500 // Buffer para picos de mensajes
	WorkerPerShard    = 20  // Workers concurrentes por shard
)

type ShardManager struct {
	shards []*Shard // Slice de shards
}

type Command struct {
	Type    string
	Id      string
	Player  *players.Players
	Wb      *websocket.Conn
	Message []byte
}

var eventPriorities = map[config.EventServer]config.MessagePriority{
	config.MovePlayer:     config.High,
	config.RayInteraction: config.Medium,
}

func getEventPriority(event config.EventServer) config.MessagePriority {
	if priority, exists := eventPriorities[event]; exists {
		return priority
	}
	return config.Low
}

func NewShardManager() *ShardManager {
	sm := &ShardManager{
		shards: make([]*Shard, Shards),
	}

	for i := 0; i < Shards; i++ {
		// Inicializar cada shardr
		shard := &Shard{
			sm:                 sm,
			connections:        make(map[string]*websocket.Conn),
			players:            make(map[string]*players.Players),
			highPriorityChan:   make(chan MessageObject, MsgBufferPerShard/3),
			mediumPriorityChan: make(chan MessageObject, MsgBufferPerShard/3),
			lowPriorityChan:    make(chan MessageObject, MsgBufferPerShard/3),
			commandChan:        make(chan Command, 10),
		}
		// shard.StartMetricsLogger()

		sm.shards[i] = shard

		// Iniciar workers para el shard
		for w := 0; w < WorkerPerShard; w++ {
			go shardWorker(shard)
		}
	}

	return sm
}

func (sm *ShardManager) EnquequeMessage(msg MessageObject) {
	shard := sm.getShard(msg.From)
	priority := getEventPriority(msg.Event)
	shard.Enqueue(msg, priority)
}

func (sm *ShardManager) GlobalBroadcast(data []byte) {
	for _, shard := range sm.shards {
		shard.commandChan <- Command{
			Type:    "broadcast",
			Message: data,
		}
	}
}

func (sm *ShardManager) AddPlayer(id string, player *players.Players, wb *websocket.Conn) {
	shard := sm.getShard(id)
	shard.commandChan <- Command{
		Type:   "add",
		Id:     id,
		Player: player,
		Wb:     wb,
	}
}

func (sm *ShardManager) RemovePlayer(id string) {
	shard := sm.getShard(id)
	shard.commandChan <- Command{
		Type: "remove",
		Id:   id,
	}
}

func (sm *ShardManager) getShard(userID string) *Shard {
	h := fnv.New32a()
	h.Write([]byte(userID))
	return sm.shards[h.Sum32()%Shards]
}

func handleMovement(s *Shard, player *players.Players, dataPlayer string) {
	s.mu.Lock()
	playersCopy := make([]players.Players, 0, len(s.players))
	for _, py := range s.players {
		if py.ID != player.ID {
			modifiedPlayer := *py
			modifiedPlayer.Head.Position.Y += 0.3
			playersCopy = append(playersCopy, modifiedPlayer)
		}
	}
	s.mu.Unlock()

	dataPlayerMarshal, _ := players.ConvertToJson(dataPlayer)

	s.mu.Lock()
	s.players[player.ID] = dataPlayerMarshal
	s.mu.Unlock()

	if len(playersCopy) == 0 {
		return
	}

	allPlayersJSON, _ := json.Marshal(players.PlayersWrapper{Players: playersCopy})

	data := MessageObject{
		Data:  string(allPlayersJSON),
		From:  player.ID,
		Event: config.MovePlayer,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	s.mu.Lock()
	err = s.connections[player.ID].WriteMessage(websocket.BinaryMessage, jsonData)
	s.mu.Unlock()

	if err != nil {
		log.Printf("Error escribiendo en WebSocket: %v", err)
	}
}

func handleRayInteraction(s *Shard, player *players.Players, eventData string) {
	s.mu.RLock()
	data := MessageObject{
		Data:  eventData,
		From:  player.ID,
		Event: config.RayInteraction,
	}
	s.mu.RUnlock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	broadcastUpdate(s, jsonData)
}

var ManagerShading = NewShardManager()
