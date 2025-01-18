package main

import (
	"log"

	"github.com/DuvanRozoParra/try_server/config"
	"github.com/DuvanRozoParra/try_server/websocketgame"
)

func main() {
	ServerCbtic := websocketgame.ServerVR()
	log.Fatal(ServerCbtic.Listen(config.Address))
	/*
		manager := players.NewManagePlayers(5)

		playerData1 := `{
			"id": "player1",
			"head": {"position": {"x": 1, "y": 1, "z": 1}, "rotation": {"x": 0, "y": 0, "z": 0, "w": 1}},
			"body": {"position": {"x": 2, "y": 2, "z": 2}, "rotation": {"x": 0, "y": 0, "z": 0, "w": 1}},
			"handLeft": {"position": {"x": 3, "y": 3, "z": 3}, "rotation": {"x": 0, "y": 0, "z": 0, "w": 1}},
			"handRight": {"position": {"x": 4, "y": 4, "z": 4}, "rotation": {"x": 0, "y": 0, "z": 0, "w": 1}}
		}`

		playerData2 := `{
			"id": "player2",
			"head": {"position": {"x": 5, "y": 5, "z": 5}, "rotation": {"x": 0, "y": 0, "z": 0, "w": 1}},
			"body": {"position": {"x": 6, "y": 6, "z": 6}, "rotation": {"x": 0, "y": 0, "z": 0, "w": 1}},
			"handLeft": {"position": {"x": 7, "y": 7, "z": 7}, "rotation": {"x": 0, "y": 0, "z": 0, "w": 1}},
			"handRight": {"position": {"x": 8, "y": 8, "z": 8}, "rotation": {"x": 0, "y": 0, "z": 0, "w": 1}}
		}`

		// Agregar jugadores.
		err := manager.AddPlayer(playerData1)
		if err != nil {
			fmt.Println("Error:", err)
		}

		err = manager.AddPlayer(playerData2)
		if err != nil {
			fmt.Println("Error:", err)
		}

		// Obtener y mostrar todos los jugadores.
		allPlayers, _ := manager.GetAllPlayers()
		fmt.Printf("All players: %+v\n", allPlayers)
	*/
}
