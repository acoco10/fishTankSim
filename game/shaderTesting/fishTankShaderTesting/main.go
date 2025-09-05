package main

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/scenes"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
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
	screenWidth  = 1920
	screenHeight = 1080
)

type Game struct {
	testSprite         *sprite.Sprite
	animatedTestSprite *sprite.AnimatedSprite
	offScreen          *ebiten.Image
	offScreenParams    map[string]any
	shader             *ebiten.Shader
	img                *ebiten.Image
	imgNormal          *ebiten.Image
	vertices           [4]ebiten.Vertex
	smallerResolution  *ebiten.Image
	resolutionScalar   float64
}

func newGame() *Game {
	g := Game{}
	//collisionMap, err := geometry.LoadCollisions()
	g.offScreen = ebiten.NewImage(960, 540)
	g.smallerResolution = ebiten.NewImage(960, 540)

	shaderParams := make(map[string]any)
	shaderParams["LightPoint"] = [2]float64{150, 150}
	shaderParams["LightWidth"] = 10.0
	shaderParams["TankRect"] = [4]float64{100, 100, 200, 200}

	//g.offScreenShader = shader
	g.offScreenParams = shaderParams
	g.resolutionScalar = scenes.ScaleScreenToResolution()

	g.shader = shaders.LoadOnePointLightingBlue()

	fishSprite, err := loader.LoadFishSprite(entities.GoldFish, 2)
	if err != nil {
		log.Fatal(err)
	}
	fishSprite.Shader = nil

	fishSprite.X = 150
	fishSprite.Y = 100
	g.animatedTestSprite = fishSprite
	return &g
}

func (g *Game) Update() error {
	//g.testSprite.CharUpdate()
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

	g.offScreen.Fill(colornames.Orange)
	dopts := ebiten.DrawImageOptions{}

	dopts.GeoM.Translate(float64(g.animatedTestSprite.X), float64(g.animatedTestSprite.Y))
	dopts.GeoM.Scale(2, 2)

	g.animatedTestSprite.UpdateOpts(dopts)
	g.animatedTestSprite.Draw(g.offScreen)

	bounds := g.smallerResolution.Bounds()

	g.vertices[0].DstX = float32(bounds.Min.X) // top-left
	g.vertices[0].DstY = float32(bounds.Min.Y) // top-left
	g.vertices[1].DstX = float32(bounds.Max.X) // top-right
	g.vertices[1].DstY = float32(bounds.Min.Y) // top-right
	g.vertices[2].DstX = float32(bounds.Min.X) // bottom-left
	g.vertices[2].DstY = float32(bounds.Max.Y) // bottom-left
	g.vertices[3].DstX = float32(bounds.Max.X) // bottom-right
	g.vertices[3].DstY = float32(bounds.Max.Y) // bottom-right

	srcBounds := g.offScreen.Bounds()

	g.vertices[0].SrcX = float32(srcBounds.Min.X) // top-left
	g.vertices[0].SrcY = float32(srcBounds.Min.Y) // top-left
	g.vertices[1].SrcX = float32(srcBounds.Max.X) // top-right
	g.vertices[1].SrcY = float32(srcBounds.Min.Y) // top-right
	g.vertices[2].SrcX = float32(srcBounds.Min.X) // bottom-left
	g.vertices[2].SrcY = float32(srcBounds.Max.Y) // bottom-left
	g.vertices[3].SrcX = float32(srcBounds.Max.X) // bottom-right
	g.vertices[3].SrcY = float32(srcBounds.Max.Y) // bottom-right

	shadOpts := &ebiten.DrawTrianglesShaderOptions{}
	shadOpts.Images[0] = g.offScreen
	shadOpts.Uniforms = g.offScreenParams

	shaderOpts := &ebiten.DrawRectShaderOptions{}
	shaderOpts.Uniforms = g.offScreenParams
	shaderOpts.Images[0] = g.offScreen

	//indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	g.smallerResolution.DrawRectShader(srcBounds.Dx(), srcBounds.Dy(), g.shader, shaderOpts)

	vector.StrokeCircle(g.smallerResolution, 150, 150, 10, 10, colornames.Burlywood, false)
	dopts.GeoM.Reset()
	dopts.GeoM.Scale(float64(g.resolutionScalar), float64(g.resolutionScalar))

	screen.DrawImage(g.smallerResolution, &dopts)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 800
}

func main() {
	g := newGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sprite Outline")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
