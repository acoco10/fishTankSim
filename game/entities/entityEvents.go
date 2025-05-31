package entities

import (
	"github.com/acoco10/fishTankWebGame/game/geometry"
)

type CreatureReachedPoint struct {
	Point    *geometry.Point
	Creature *Creature
}

type PointGenerated struct {
	Point  *geometry.Point
	Source string
}

type SendData struct {
	DataFor string
	Data    string
}

type RequestData struct {
	DataType   string
	RequestFor any
}

type FishEvent struct {
	fish  *Creature
	event string
}

type FishLevelUp struct {
	fish Creature
}
