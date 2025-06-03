package main

import (
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/loaders"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
)

const (
	screenWidth  = 1000
	screenHeight = 1000
)

type Game struct {
	testGraphic *graphics.SpriteGraphic
}

func newGame() *Game {
	g := Game{}

	whiteBoardGraphic, err := loaders.LoadClothGraphic()
	if err != nil {
		log.Fatal(err)
	}
	g.testGraphic = whiteBoardGraphic

	return &g

}

func (g *Game) Update() error {
	g.testGraphic.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	g.testGraphic.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1000, 1000
}

func main() {
	g := newGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Graphic Testing")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
