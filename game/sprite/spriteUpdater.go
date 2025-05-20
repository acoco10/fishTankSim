package sprite

import "github.com/hajimehoshi/ebiten/v2"

type XYUpdater struct {
	offSetX float32
	offSetY float32
	Loaded  bool
	*Sprite
}

func NewUpdater(sprite *Sprite) *XYUpdater {
	x, y := ebiten.CursorPosition()
	difX := float32(x) - sprite.X
	difY := float32(y) - sprite.Y
	newUpdater := XYUpdater{difX, difY, false, sprite}
	return &newUpdater
}

func (up *XYUpdater) Update() {
	x, y := ebiten.CursorPosition()
	up.Sprite.X = float32(x) - up.offSetX
	up.Sprite.Y = float32(y) - up.offSetY
}
