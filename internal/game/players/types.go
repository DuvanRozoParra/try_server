package players

type RolePlayer = int

const (
	Admin RolePlayer = iota
	Player
	Spectator
)

type AnimationHand struct {
	Thumb      bool
	Pointer    bool
	RingFinger bool
}

type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Quaternion struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	W float64 `json:"w"`
}

type BodyPart struct {
	Position Vector3    `json:"position"`
	Rotation Quaternion `json:"rotation"`
}

type PlayersWrapper struct {
	Players []Players `json:"players"`
}
