package drawables

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type DrawableSaveAbleSprite interface {
	Draw(screen *ebiten.Image)
	Update()
	SavePosition() SavePositionData
}

type SavePositionData struct {
	X    float32
	Y    float32
	Name string
}

type Drawable interface {
	Draw(screen *ebiten.Image)
	Update()
	SpriteHovered() bool
	Coord() (float32, float32)
}
