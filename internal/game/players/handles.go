package players

import (
	"encoding/json"
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
	players map[string]*Players
	rw      sync.RWMutex
}

func NewPlayer(id string) *Players {
	return &Players{
		ID:        id,
		Head:      BodyPart{Position: Vector3{X: 0, Y: 0, Z: 0}, Rotation: Quaternion{X: 0, Y: 0, Z: 0, W: 0}},
		Body:      BodyPart{Position: Vector3{X: 0, Y: 0, Z: 0}, Rotation: Quaternion{X: 0, Y: 0, Z: 0, W: 0}},
		HandLeft:  BodyPart{Position: Vector3{X: 0, Y: 0, Z: 0}, Rotation: Quaternion{X: 0, Y: 0, Z: 0, W: 0}},
		HandRight: BodyPart{Position: Vector3{X: 0, Y: 0, Z: 0}, Rotation: Quaternion{X: 0, Y: 0, Z: 0, W: 0}},
	}
}

func ConvertToJson(data string) (*Players, error) {
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
