package debug

import (
	"github.com/acoco10/fishTankWebGame/game/util"
)

type DebugOption uint8

const (
	Normal DebugOption = iota
	Position
	Print
	ShaderTest
)

type GameMode uint8

const (
	User GameMode = iota
	Debug
)

type DebugData struct {
	DebugText      string
	DebugRect      *util.Rect
	DebugParameter map[DebugOption]bool
	GameMode       GameMode
}

func (d *DebugData) Update() {

}
