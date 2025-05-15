package gameEntities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type VectorLineGraphic struct {
	X0, Y0 float32
	X1, Y1 float32
	maxX1  float32
	clr    color.Color
}

func (v *VectorLineGraphic) Update() {
	if v.X1 < v.maxX1 {
		v.X1 += 10
	}
}

func (v *VectorLineGraphic) Draw(screen *ebiten.Image) {
	vector.StrokeLine(screen, v.X0, v.Y0, v.X1, v.Y1, 4, v.clr, false)
}
