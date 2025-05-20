package ui

import (
	"github.com/acoco10/fishTankWebGame/game/geometry"
)

type ButtonEvent struct {
	ButtonText string
	EType      string
}

type ButtonClickedEvent struct {
	ButtonText string
}

type MouseButtonPressed struct {
	Point *geometry.Point
}
