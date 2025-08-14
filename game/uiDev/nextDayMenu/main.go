package main

import (
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"log"
)

const (
	screenWidth  = 1000
	screenHeight = 1000
)

type Game struct {
	ui  *ebitenui.UI
	hub *tasks.EventHub
}

func newGame() *Game {

	loader.LoadFontRegistry()

	g := Game{}
	hub := tasks.NewEventHub()
	g.hub = hub
	ndui, _, err := ui.LoadMainFishMenu(1000, 1000, hub)
	if err != nil {
		log.Fatal(err)
	}
	g.ui = ndui
	return &g
}

func (g *Game) Update() error {
	g.ui.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ui.TriggerTextWindow(g.hub, g.ui, "How To Play",
			"1. Press Space and U to start your mower \n"+
				" 2. Hold Space to keep it running\n "+
				" 3. Mow as much grass as possible")
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Lightpink)
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
