package sprite

import (
	"github.com/acoco10/fishTankWebGame/game/util"
)

type XYUpdater struct {
	offSetX float32
	offSetY float32
}

func (x *XYUpdater) SetCustomOffset(X, Y int) {
	if X != 0 {
		x.offSetX = float32(X)
	}
	if Y != 0 {
		x.offSetY = float32(Y)
	}
}

func NewUpdater(sprite *Sprite) *XYUpdater {
	x, y := util.GetScaledCursorPosition()
	difX := float32(x) - sprite.X
	difY := float32(y) - sprite.Y
	newUpdater := XYUpdater{difX, difY}
	return &newUpdater
}

func (up *XYUpdater) Update(sprite *Sprite) {
	x, y := util.GetScaledCursorPosition()
	sprite.X = float32(x) - up.offSetX
	sprite.Y = float32(y) - up.offSetY
}
