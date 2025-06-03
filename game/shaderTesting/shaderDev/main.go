package main

import (
	"github.com/acoco10/fishTankWebGame/game/loaders"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
)

type Direction uint8

type TextWithShader struct {
	image        *ebiten.Image
	shader       *ebiten.Shader
	shaderParams map[string]any
	text         string
	face         text.Face
	updateFunc   func(map[string]any) map[string]any
}

func NewTextWithShader(text string, dst *ebiten.Image) *TextWithShader {
	face, err := ui.LoadFont(18, "rockSalt")
	if err != nil {
		log.Fatal(err)
	}

	ts := &TextWithShader{}
	ts.text = text
	ts.face = face

	ts.updateFunc = shaders.UpdateCounter
	shader := shaders.LoadHandWritingShader()
	ts.shader = shader

	ts.shaderParams = make(map[string]any)
	ts.shaderParams["Counter"] = 0

	ts.image = dst

	return ts
}

func (t *TextWithShader) Update() {
	t.shaderParams = t.updateFunc(t.shaderParams)
}

func (t *TextWithShader) Draw(dst *ebiten.Image) {

	dopts := &text.DrawOptions{}
	shaderOpts := ebiten.DrawRectShaderOptions{}
	dopts.ColorScale.Scale(0, 0, 0, 1)
	dopts.GeoM.Translate(float64(10), float64(10))

	text.Draw(t.image, t.text, t.face, dopts)
	shaderOpts.Uniforms = t.shaderParams
	shaderOpts.Images[0] = t.image

	dst.DrawRectShader(t.image.Bounds().Dx(), t.image.Bounds().Dy(), t.shader, &shaderOpts)
}

const (
	Left Direction = iota
	Right
	Down
	Up
)

// ... (rest of your import statements and setup)

var (
	outlineColor  = [4]float64{0.2, 0.1, 0.05, 255}        // Yellow outline
	outlineColor2 = color.RGBA{R: 1, G: 255, B: 1, A: 255} // Yellow outline
)

const (
	screenWidth  = 640 * 2
	screenHeight = 480 * 2
)

type Game struct {
	testSprite *sprite.Sprite
	ts         *TextWithShader
}

func newGame() *Game {
	g := Game{}

	whiteBoard, err := loaders.LoadImageAssetAsEbitenImage("uiSprites/whiteBoardMain")
	if err != nil {
		log.Fatal(err)
	}
	dst := ebiten.NewImage(whiteBoard.Bounds().Dx(), whiteBoard.Bounds().Dy())

	ts := NewTextWithShader("testing testing", dst)
	g.ts = ts
	testSprite := sprite.Sprite{Img: whiteBoard, X: 250, Y: 250}

	g.testSprite = &testSprite

	//shader := shaders
	//s.shader = outlineShader

	return &g
}

func (g *Game) Update() error {
	g.testSprite.Update()
	g.ts.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	dopts := text.DrawOptions{}

	dopts.ColorScale.Scale(0, 0, 0, 1)
	dopts.GeoM.Translate(float64(10), float64(10))

	g.ts.Draw(g.testSprite.Img)

	g.testSprite.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := newGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hand writing shader")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
