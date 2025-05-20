package drawables

import "github.com/hajimehoshi/ebiten/v2"

type DrawableSprite interface {
	Draw(screen *ebiten.Image)
	Update()
	SavePosition() SavePositionData
}

type SavePositionData struct {
	X    float32
	Y    float32
	Name string
}
