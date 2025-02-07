package main

import (
	"log"

	"github.com/DuvanRozoParra/try_server/config"
	network "github.com/DuvanRozoParra/try_server/internal"
)

func main() {
	ServerCbtic := network.ServerVR()
	log.Fatal(ServerCbtic.Listen(config.Address))
}
