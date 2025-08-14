package input

import (
	"github.com/acoco10/fishTankWebGame/game/util"
)

type MouseButtonPressedUISpriteActivity struct {
	Point *util.Point
}

func (m MouseButtonPressedUISpriteActivity) Type() string {
	return "MouseButtonPressUISpriteActivity"
}

type CursorOccupied struct {
}

func (c CursorOccupied) Type() string {
	return "CursorOccupied"
}
