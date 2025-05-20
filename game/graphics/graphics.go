package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type VectorLineGraphic struct {
	X0, Y0 float32
	X1, Y1 float32
	MaxX1  float32
	Color  color.Color
}

func (v *VectorLineGraphic) Update() {
	if v.X1 < v.MaxX1 {
		v.X1 += 10
	}
}

func (v *VectorLineGraphic) Draw(screen *ebiten.Image) {
	vector.StrokeLine(screen, v.X0, v.Y0, v.X1, v.Y1, 4, v.Color, false)
}

func (v *VectorLineGraphic) SavePosition() drawables.SavePositionData {
	s := drawables.SavePositionData{}
	s.Name = "vector line graphic"
	s.X = v.X0
	s.Y = v.Y0

	return s
}
