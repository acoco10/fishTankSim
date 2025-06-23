package main

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/acoco10/fishTankWebGame/game/loaders"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

type Direction uint8

type TextWithShader struct {
	sprites *sprite.Sprite
	face    text.Face
	erasing bool
}

func NewTextWithShader(text string, dst *ebiten.Image) *TextWithShader {

	ts := &TextWithShader{}

	return ts
}

func (t *TextWithShader) Update() {

}

const (
	screenWidth  = 640 * 2
	screenHeight = 480 * 2
)

type Game struct {
	tmap           []int
	shaderParams   map[string]any
	shader         *ebiten.Shader
	animatedSprite *sprite.AnimatedSprite
	normalSprite   *sprite.AnimatedSprite
}

func newGame() *Game {

	diffuseImg, err := loaders.LoadImageAssetAsEbitenImage("fishSpriteSheets/mollyFish3SpriteSheet")
	if err != nil {
		log.Fatal(err)
	}
	/*
		uniform vec2 u_LightPos;     // In screen space [0,1]
		uniform vec3 u_LightColor;
		uniform vec3 u_AmbientColor;*/
	normalImg, err := loaders.LoadImageAssetAsEbitenImage("fishSpriteSheets/mollyFish3NormalSpriteSheet")
	if err != nil {
		log.Fatal(err)
	}

	mollyAnim := animations.NewAnimation(0, 3, 1, 20)
	mollySpriteSheet := spritesheet.NewSpritesheet(4, 1, 65, 37)

	mollySprite := sprite.NewAnimatedSprite()
	mollyNormals := sprite.NewAnimatedSprite()

	mollySprite.Img = diffuseImg
	mollySprite.Animation = mollyAnim
	mollySprite.SpriteSheet = mollySpriteSheet

	mollyNormals.Img = normalImg
	mollyNormals.Animation = mollyAnim
	mollyNormals.SpriteSheet = mollySpriteSheet

	shader := shaders.LoadNormalMapShader()

	uniforms := make(map[string]any)

	x, y := ebiten.CursorPosition()
	u := float32(x) / float32(screenWidth)
	v := float32(y) / float32(screenHeight)

	g := Game{}
	g.shaderParams = uniforms
	g.shaderParams["Cursor"] = []float32{u, v}
	g.shader = shader
	g.animatedSprite = mollySprite
	g.normalSprite = mollyNormals

	g.animatedSprite.X = 100
	g.animatedSprite.Y = 100
	//shader := shaders
	//s.shader = outlineShader

	return &g
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	u := float32(x)
	v := float32(y)
	g.animatedSprite.Update()
	g.shaderParams["Cursor"] = []float32{u, v}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	shaderOpts := ebiten.DrawRectShaderOptions{}
	shaderOpts.GeoM.Scale(1, 1)
	shaderOpts.GeoM.Translate(300, 300)

	frame := g.animatedSprite.Frame()
	rect := g.animatedSprite.Rect(frame)

	spImg := g.animatedSprite.Img.SubImage(rect).(*ebiten.Image)
	normalImg := g.normalSprite.Img.SubImage(rect).(*ebiten.Image)

	shaderOpts.Uniforms = g.shaderParams
	shaderOpts.Images[0] = spImg
	shaderOpts.Images[1] = normalImg

	screen.DrawRectShader(rect.Dx(), rect.Dy(), g.shader, &shaderOpts)
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
