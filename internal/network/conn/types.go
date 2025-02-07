package conn

import (
	"github.com/DuvanRozoParra/try_server/config"
)

type MessageObject struct {
	Data  string             `json:"data"`
	From  string             `json:"from"`
	Event config.EventServer `json:"events"`
}
