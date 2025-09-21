package util

import "image/color"

func ConvertRGBAUintToDecimal(clr color.RGBA) [4]float64 {

	r := float64(clr.R) / 255
	g := float64(clr.G) / 255
	b := float64(clr.B) / 255
	a := float64(clr.A) / 255

	return [4]float64{r, g, b, a}
}
