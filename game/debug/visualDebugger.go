package debug

import "github.com/acoco10/fishTankWebGame/game/geometry"

type DebugOption uint8

const (
	Normal DebugOption = iota
	Position
	Print
)

type GameMode uint8

const (
	User GameMode = iota
	Debug
)

type DebugData struct {
	DebugText      string
	DebugRect      *geometry.Rect
	DebugParameter map[DebugOption]bool
	GameMode       GameMode
}

func (d *DebugData) Update() {

}
