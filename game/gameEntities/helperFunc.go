package gameEntities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func LoadImageAssetAsEbitenImage(assetName string) (*ebiten.Image, error) {
	imgPath := fmt.Sprintf("images/%s.png", assetName)
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, imgPath)
	if err != nil {
		return &ebiten.Image{}, err
	}
	return img, nil
}

func DrawShader(sprite Sprite, sImg *ebiten.Image, s *ebiten.Shader, screen *ebiten.Image) {
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

type XYUpdater struct {
	offSetX float32
	offSetY float32
	Loaded  bool
	*Sprite
}

func NewUpdater(sprite *Sprite) *XYUpdater {
	x, y := ebiten.CursorPosition()
	difX := float32(x) - sprite.X
	difY := float32(y) - sprite.Y
	newUpdater := XYUpdater{difX, difY, false, sprite}
	return &newUpdater
}

func (up *XYUpdater) Update() {
	x, y := ebiten.CursorPosition()
	up.Sprite.X = float32(x) - up.offSetX
	up.Sprite.Y = float32(y) - up.offSetY
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
