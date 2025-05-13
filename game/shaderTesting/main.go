package main

import (
	"github.com/acoco10/fishTankWebGame/game/gameEntities"
	"github.com/acoco10/fishTankWebGame/game/shaderHelpers"
	ui2 "github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/ebitenui/ebitenui"
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
	screenWidth  = 500
	screenHeight = 500
)

type Game struct {
	img *ebiten.Image
	s   *ebiten.Shader
	*gameEntities.Timer
	pulseDuration *gameEntities.Timer
	pulse         bool
	drawPointX    float32
	drawPointY    float32
	direction     Direction
	shaderParams  map[string]any
	ui            *ebitenui.UI
	ehub          *gameEntities.EventHub
}

func newGame() *Game {
	g := Game{}
	img, err := gameEntities.LoadImageAssetAsEbitenImage("uiSprites/fishFoodMain")
	if err != nil {
		log.Fatal(err)
	}
	g.ehub = gameEntities.NewEventHub()
	g.ui, _, err = ui2.LoadMainFishMenu(10, 0, g.ehub)

	g.img = img
	g.drawPointX = 0
	g.drawPointY = 0

	outlineShader := gameEntities.LoadSolidColorShader()
	g.s = outlineShader
	g.direction = Right

	g.Timer = gameEntities.NewTimer(1)
	g.pulseDuration = gameEntities.NewTimer(0.3)
	g.Timer.TurnOn()

	clrArr := outlineColor
	clrArr2 := shaderHelpers.ClrToArr(outlineColor2)

	var paramaMapa = make(map[string]any)
	paramaMapa["OutlineColor"] = clrArr
	paramaMapa["OutlineColor2"] = clrArr2
	paramaMapa["DrawPointX"] = 0
	paramaMapa["DrawPointY"] = 0
	g.shaderParams = paramaMapa

	return &g

}

func UpdateDPoint(g *Game) {
	if g.direction == Right {
		g.drawPointX++
	}
	if g.direction == Left {
		g.drawPointX--
	}
	if g.direction == Down {
		g.drawPointY++
	}
	if g.direction == Up {
		g.drawPointY--
	}
	g.ui.Update()
}

func UpdateDirection(currentD Direction, x, y float32, xBound, yBound float32) Direction {
	if x == 0 && y == 0 {
		return Right
	}

	if x >= xBound && y <= 0 {
		println("changing direction to down")
		return Down
	}

	if x >= xBound && y >= yBound {
		println("changing direction to left")
		return Left
	}

	if x <= 0 && y >= yBound {
		println("changing direction to up")
		return Up
	}
	return currentD
}

func (g *Game) Update() error {
	b := g.img.Bounds()

	g.direction = UpdateDirection(g.direction, g.drawPointX, g.drawPointY, float32(b.Max.X), float32(b.Max.Y))

	UpdateDPoint(g)

	println(g.drawPointY)

	g.shaderParams["DrawPointX"] = g.drawPointX
	g.shaderParams["DrawPointY"] = g.drawPointY

	return nil

}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	g.ui.Draw(screen)
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = g.img
	opts.Uniforms = g.shaderParams
	opts.GeoM.Translate(100, 100)
	opts.GeoM.Scale(2, 2)

	screen.DrawRectShader(g.img.Bounds().Dx(), g.img.Bounds().Dy(), g.s, opts) // Use DrawImage with shader options
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
