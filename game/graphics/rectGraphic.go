package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

type RectGraphic struct {
	image.Rectangle
	Color color.RGBA
}

func (rg *RectGraphic) Update() {
	return
}

func (rg *RectGraphic) Scaled() ScaledType {
	return NormalScaled
}

func NewRectGraphic(rectangle image.Rectangle, clr color.RGBA) int {
	rg := &RectGraphic{rectangle, clr}
	return AddGraphic(rg)
}

func (rg *RectGraphic) Draw(screen *ebiten.Image) {
	util.StrokeRectFromImageRect(rg.Rectangle, screen, rg.Color)
}
