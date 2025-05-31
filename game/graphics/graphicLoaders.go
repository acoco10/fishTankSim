package graphics

import (
	"image/color"
)

func NewVlS(x0, y0, x1, y1, maxX float32, clr color.Color) *VectorLineGraphic {
	vls := VectorLineGraphic{
		X0:    x0,
		Y0:    y0,
		X1:    x1,
		Y1:    y1,
		Color: clr,
		MaxX1: maxX,
	}
	return &vls
}
