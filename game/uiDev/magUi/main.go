package main

import (
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
)

const (
	screenWidth  = 1000
	screenHeight = 1000
)

type Game struct {
	ui  *ui.Magazine
	hub *tasks.EventHub
}

func newGame() *Game {
	g := Game{}
	hub := tasks.NewEventHub()
	g.hub = hub
	magUi, err := ui.LoadMagazineUiMenu(hub, 1000, 1000)
	if err != nil {
		log.Fatal(err)
	}
	g.ui = magUi
	return &g

}

func (g *Game) Update() error {
	g.ui.Update()
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.ui.Trigger()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	g.ui.Draw(screen)

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
