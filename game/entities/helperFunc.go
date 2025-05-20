package entities

import (
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func DrawShader(sprite sprite.Sprite, sImg *ebiten.Image, s *ebiten.Shader, screen *ebiten.Image) {
	vertices := [4]ebiten.Vertex{}
	bounds := sImg.Bounds()

	vertices[0].DstX = sprite.X                         // top-left
	vertices[0].DstY = sprite.Y                         // top-left
	vertices[1].DstX = sprite.X + float32(bounds.Max.X) // top-right
	vertices[1].DstY = sprite.Y                         // top-right
	vertices[2].DstX = sprite.X                         // bottom-left
	vertices[2].DstY = sprite.Y + float32(bounds.Max.Y) // bottom-left
	vertices[3].DstX = sprite.X + float32(bounds.Max.X) // bottom-right
	vertices[3].DstY = sprite.Y + float32(bounds.Max.Y) // bottom-right

	var shaderOpts ebiten.DrawTrianglesShaderOptions
	shaderOpts.Images[0] = sImg
	//draw shader
	indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	screen.DrawTrianglesShader(vertices[:], indices, s, &shaderOpts)
}

func ApplyShaderToText(screen *ebiten.Image, inputText string, face text.Face) {

	//just an idea/ expirement with offscreen rendering
	dOpts := text.DrawOptions{}
	//dOpts.GeoM.Translate(ScreenWidth/2-float64(len(debugText)*6), ScreenHeight/10)
	offScreen := ebiten.NewImage(400, 100)
	text.Draw(offScreen, inputText, face, &dOpts)
	shader := LoadOutlineShader()
	sOpts := ebiten.DrawRectShaderOptions{}
	sOpts.GeoM.Translate(100, 100)
	sOpts.Images[0] = offScreen
	var paramaMapa = make(map[string]any)
	clrArr := [4]float64{0.2, 0.6, 0.05, 255}
	paramaMapa["OutlineColor"] = clrArr
	sOpts.Uniforms = paramaMapa

	screen.DrawRectShader(400, 100, shader, &sOpts)
}
