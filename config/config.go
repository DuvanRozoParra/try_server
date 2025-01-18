package config

type ActionType int
type EventServer int

const (
	Address string = ":8080"
)

const (
	RayInteraction EventServer = iota
	MovePlayer
)
