package config

type ActionType int
type EventServer int
type MessagePriority int

const (
	Address string = ":8080"
)

const (
	RayInteraction EventServer = iota
	MovePlayer
	DeletePlayer
	ActionHandsPlayer
)

const (
	High MessagePriority = iota
	Medium
	Low
)
