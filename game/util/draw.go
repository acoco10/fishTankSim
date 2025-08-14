package util

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
)

func StrokeRectFromImageRect(rectangle image.Rectangle, screen *ebiten.Image, clr color.RGBA) {

	x1 := float32(rectangle.Min.X)
	y1 := float32(rectangle.Min.Y)
	width := float32(rectangle.Bounds().Dx())
	height := float32(rectangle.Bounds().Dy())

	vector.StrokeRect(screen, x1, y1, width, height, 1, clr, false)
}
