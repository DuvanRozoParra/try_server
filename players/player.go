package players

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type Players struct {
	ID        string   `json:"id"` // Cambiado a mayúscula para ser exportado.
	Head      BodyPart `json:"head"`
	Body      BodyPart `json:"body"`
	HandLeft  BodyPart `json:"handLeft"`
	HandRight BodyPart `json:"handRight"`
}

type PlayersManage struct {
	Players      map[string]*Players
	LimitPlayers int
	sync.Mutex
}

func NewManagePlayers(limitUser int) *PlayersManage {
	return &PlayersManage{
		Players:      make(map[string]*Players),
		LimitPlayers: limitUser,
	}
}

func (p *PlayersManage) ConvertToJson(data string) (*Players, error) {
	var temp struct {
		ID        string `json:"id"`
		Head      string `json:"head"`
		Body      string `json:"body"`
		HandLeft  string `json:"handLeft"`
		HandRight string `json:"handRight"`
	}

	// Deserializar el JSON inicial en la estructura temporal
	err := json.Unmarshal([]byte(data), &temp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse player data: %w", err)
	}

	// Crear el objeto final de tipo Players
	var player Players
	player.ID = temp.ID

	// Deserializar las propiedades JSON embebidas
	if err := json.Unmarshal([]byte(temp.Head), &player.Head); err != nil {
		return nil, fmt.Errorf("failed to parse head data: %w", err)
	}
	if err := json.Unmarshal([]byte(temp.Body), &player.Body); err != nil {
		return nil, fmt.Errorf("failed to parse body data: %w", err)
	}
	if err := json.Unmarshal([]byte(temp.HandLeft), &player.HandLeft); err != nil {
		return nil, fmt.Errorf("failed to parse handLeft data: %w", err)
	}
	if err := json.Unmarshal([]byte(temp.HandRight), &player.HandRight); err != nil {
		return nil, fmt.Errorf("failed to parse handRight data: %w", err)
	}

	return &player, nil
}

func (p *PlayersManage) AddPlayer(data string) error {
	player, err := p.ConvertToJson(data)
	if err != nil {
		return err
	}

	if len(p.Players) >= p.LimitPlayers {
		return errors.New("player limit reached")
	}

	p.Players[player.ID] = player
	return nil
}

func (p *PlayersManage) PlayerExists(id string) (bool, *Players) {
	p.Lock()
	defer p.Unlock()

	player, exists := p.Players[id]
	return exists, player
}

func (p *PlayersManage) ModifyPlayer(id string, data string) {
	//p.Lock()

	player, _ := p.ConvertToJson(data)

	exists, _ := p.PlayerExists(id)
	if exists {
		//return fmt.Errorf("player with ID '%s' does not exist", id)
		p.Players[id] = player
	}

	//p.Unlock()

}

func (p *PlayersManage) GetAllPlayers() (string, error) {

	allPlayers := make([]Players, 0, len(p.Players))

	for _, player := range p.Players {
		allPlayers = append(allPlayers, *player)
	}

	players := PlayersWrapper{
		Players: allPlayers,
	}

	data, err := json.Marshal(players)
	if err != nil {
		return "", fmt.Errorf("failed to marshal players: %w", err)
	}

	return string(data), nil
}

func (p *PlayersManage) RemovePlayer(id string) error {
	p.Lock()
	defer p.Unlock()

	exists, _ := p.PlayerExists(id)
	if !exists {
		return fmt.Errorf("player with ID '%s' does not exist", id)
	}

	delete(p.Players, id)
	fmt.Printf("Player with ID '%s' has been removed.\n", id)
	return nil
}
