package players

import (
	"fmt"
)

var Mp PlayersManage = PlayersManage{
	Players:      make(map[string]*Players),
	LimitPlayers: 5,
}

func MovePlayer(playerID, data string) (string, error) {
	player, err := Mp.ConvertToJson(data)
	if err != nil {
		return "", fmt.Errorf("error al convertir datos del jugador %s: %v", playerID, err)
	}

	exist, _ := Mp.PlayerExists(playerID)
	if !exist {
		return "", fmt.Errorf("error al verificar existencia del jugador %s: %v", playerID, err)
	}

	Mp.Players[playerID] = player

	players, err := Mp.GetDataPlayers("aasd")
	if err != nil {
		return "", fmt.Errorf("error al obtener jugadores para %s: %v", playerID, err)
	}

	return players, nil
}
