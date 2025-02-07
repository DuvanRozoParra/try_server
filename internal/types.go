package network

import "github.com/DuvanRozoParra/try_server/config"

type ServerMessage struct {
	Data   string             `json:"data"`
	From   string             `json:"from"`
	Events config.EventServer `json:"events"`
}
