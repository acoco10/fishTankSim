package input

import (
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputState uint8

type MouseFlags struct {
	HandledClick   bool
	CursorOccupied bool
}

type InputManager struct {
}

func handleCursorClick(x, y int, flags MouseFlags, hub *tasks.EventHub) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if !flags.CursorOccupied {
			hub.Publish(CursorPressed{})
		}
	}
}
