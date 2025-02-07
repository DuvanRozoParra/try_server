package events

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/gofiber/contrib/websocket"
)

type Connections struct {
	players   sync.Map 
	limitUser int32
	count     int32 
}

func newConnection(limitUser int) *Connections {
	return &Connections{
		limitUser: int32(limitUser),
	}
}

func (conn *Connections) AddPlayer(id string, c *websocket.Conn) {
	if atomic.LoadInt32(&conn.count) >= conn.limitUser {
		log.Printf("LIMIT PASS: %s", id)
		return
	}

	if _, loaded := conn.players.LoadOrStore(id, c); !loaded {
		atomic.AddInt32(&conn.count, 1)
		log.Printf("PLAYER CONNECT %s", id)
	}
}

func (conn *Connections) DeletePlayer(id string) {
	if _, loaded := conn.players.LoadAndDelete(id); loaded {
		atomic.AddInt32(&conn.count, -1)
		log.Printf("PLAYER DISCONNECT %s", id)
	}
}

func (conn *Connections) Emit(data []byte) {
	conn.players.Range(func(key, value interface{}) bool {
		go func(id string, c *websocket.Conn) {
			if err := c.WriteMessage(websocket.BinaryMessage, data); err != nil {
				log.Printf("Error enviando a %s: %v", id, err)
				conn.DeletePlayer(id)
			}
		}(key.(string), value.(*websocket.Conn))
		return true
	})
}

var ManageConnections = newConnection(5)
