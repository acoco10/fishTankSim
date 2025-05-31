package geometry

type Direction uint8

const (
	// Right  = increasing x
	//Down = increasing y
	//Left = decreasing x
	//Up = decreasing y

	Right Direction = iota
	Left
	Up
	Down
)
