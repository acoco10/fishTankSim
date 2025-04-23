package gameEntities

import (
	"fishTankWebGame/assets"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

func LoadImageAssetAsEbitenImage(assetName string) *ebiten.Image {
	imgPath := fmt.Sprintf("images/%s.png", assetName)
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, imgPath)
	if err != nil {
		log.Fatal(err)
	}
	return img
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
	println("sprite x y =", up.Sprite.X, up.Sprite.Y)
	x, y := ebiten.CursorPosition()
	up.Sprite.X = float32(x) - up.offSetX
	up.Sprite.Y = float32(y) - up.offSetY
}
