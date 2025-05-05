package shaderHelpers

import (
	"image/color"
)

func ClrToArr(clr color.RGBA) [4]float64 {
	r := float64(clr.R)
	g := float64(clr.G)
	b := float64(clr.B)
	a := float64(clr.A)

	return [4]float64{r, g, b, a}
}
