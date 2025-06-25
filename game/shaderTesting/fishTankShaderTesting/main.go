package main

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/loaders"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
)

type Direction uint8

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
	screenWidth  = 800
	screenHeight = 800
)

type Game struct {
	testSprite         *sprite.Sprite
	animatedTestSprite *sprite.AnimatedSprite
	offScreen          *ebiten.Image
	offScreenParams    map[string]any
	offScreenShader    *ebiten.Shader
	img                *ebiten.Image
	imgNormal          *ebiten.Image
}

func newGame() *Game {
	g := Game{}
	//collisionMap, err := geometry.LoadCollisions()

	shader := shaders.LoadOnePointLightingBlue()
	shaderParams := make(map[string]any)

	shaderParams["ImgRect"] = [2]float64{800, 800}
	shaderParams["LightPoint"] = [2]float64{150, 0}
	g.offScreenShader = shader
	g.offScreenParams = shaderParams
	//s.shader = outlineShader

	fishSprite := loader.LoadFishSprite(entities.Fish, 2)
	fishSprite.X = 150
	fishSprite.Y = 100

	g.animatedTestSprite = fishSprite
	ls := shaders.LoadSpriteLighting()
	g.animatedTestSprite.Shader = ls

	return &g
}

func (g *Game) Update() error {
	//g.testSprite.Update()
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.animatedTestSprite.Y += 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.animatedTestSprite.Y -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.animatedTestSprite.X += 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.animatedTestSprite.X -= 1.0
	}

	g.animatedTestSprite.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//g.testSprite.Draw(screen)

	offScreen := ebiten.NewImage(800, 800)
	offScreen.Fill(color.Black)
	dopts := ebiten.DrawImageOptions{}
	shaderOpts := ebiten.DrawRectShaderOptions{}
	shaderOpts.GeoM.Translate(float64(g.animatedTestSprite.X), float64(g.animatedTestSprite.Y))
	shaderOpts.GeoM.Scale(2, 2)
	dopts.GeoM.Translate(float64(g.animatedTestSprite.X), float64(g.animatedTestSprite.Y))
	dopts.GeoM.Scale(2, 2)

	g.animatedTestSprite.Draw(offScreen, &dopts, &shaderOpts)

	shaderOpts.GeoM.Reset()
	shaderOpts.Uniforms = g.offScreenParams
	shaderOpts.Images[0] = offScreen

	screen.DrawRectShader(800, 800, g.offScreenShader, &shaderOpts)
	//vector.StrokeRect(screen, 80.0, 80.0, 10.0, 10.0, 1.0, outlineColor2, false)
	//opts := &ebiten.DrawImageOptions{}
	//screen.DrawImage(g.img, opts)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 500, 500
}

func main() {
	g := newGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sprite Outline")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
