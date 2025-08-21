package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

type RectGraphic struct {
	id int
	image.Rectangle
	Color             color.RGBA
	filled            bool
	lastUpdate        int
	updater           *float32
	maxProgress       float32
	lastUpdateCounter int
	extraRect         image.Rectangle
	extraRectColor    color.RGBA
	Kind              string
}

func (rg *RectGraphic) Update() {

	if rg.updater != nil {
		rg.lastUpdateCounter++
		if rg.lastUpdateCounter == 30 {
			if int(*rg.updater) > rg.lastUpdate {
				rg.lastUpdate = int(*rg.updater)
				dif := int(*rg.updater/rg.maxProgress*50) + rg.Min.X - rg.Max.X
				rg.extraRect = image.Rect(rg.Rectangle.Max.X, rg.Rectangle.Min.Y, rg.Rectangle.Max.X+dif, rg.Rectangle.Max.Y)
			}
		}
		if rg.lastUpdateCounter > 30 && rg.lastUpdateCounter < 120 {
			rg.extraRectColor.A += uint8(255 / 60)
		}
		if rg.lastUpdateCounter > 60 {
			rg.extraRectColor.A = 0
			rg.Rectangle.Max.X += rg.extraRect.Dx()
			rg.extraRect = image.Rect(0, 0, 0, 0)
			rg.lastUpdateCounter = 0
		}

	}
}

func (rg *RectGraphic) Scaled() ScaledType {
	return ScaledButTopLevel
}

func NewRectGraphic(rectangle image.Rectangle, clr color.RGBA) int {
	rg := &RectGraphic{Rectangle: rectangle, Color: clr}
	return AddGraphic(rg)
}

func NewFilledRectGraphic(rectangle image.Rectangle, clr color.RGBA) int {
	rg := &RectGraphic{Rectangle: rectangle, Color: clr}
	rg.filled = true
	rg.id = AddGraphic(rg)
	return rg.id
}

func NewFilledRectGraphicWithPointerWidth(rectangle image.Rectangle, clr color.RGBA, ref *float32, max float32) int {
	rg := &RectGraphic{Rectangle: rectangle, Color: clr}
	rg.filled = true
	rg.maxProgress = max
	rg.updater = ref
	rg.id = AddGraphic(rg)
	rg.extraRectColor = clr
	rg.extraRectColor.A = 0
	return rg.id
}

func (rg *RectGraphic) Draw(screen *ebiten.Image) {
	if rg.Rectangle.Dx() > 0 {

		if rg.extraRect.Dx() > 0 {
			util.StrokeRectFromImageRect(rg.extraRect, screen, rg.extraRectColor, rg.filled)
		}
		util.StrokeRectFromImageRect(rg.Rectangle, screen, rg.Color, rg.filled)
	}
}
