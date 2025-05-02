package shaderTesting

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	op := &ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]any{
		"Resolution":   []float32{float32(objectImage.Bounds().Dx()), float32(objectImage.Bounds().Dy())},
		"OutlineColor": []float32{float32(outlineColor.R) / 255, float32(outlineColor.G) / 255, float32(outlineColor.B) / 255, float32(outlineColor.A) / 255},
		"ImageSrc0":    objectImage,
	}
	screen.DrawImage(objectImage, op) // Use DrawImage with shader options
