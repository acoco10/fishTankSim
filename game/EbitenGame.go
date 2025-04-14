package game

import (
	"fishTankWebGame/assets"
	"fishTankWebGame/game/events"
	"fishTankWebGame/game/ui"
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"log"
)

type Game struct {
	img        *ebiten.Image
	background *ebiten.Image
	loaded     bool
	Creatures  []*Creature
	ui         *ebitenui.UI
	eventHub   *events.EventHub
}

const (
	screenWidth  = 550
	screenHeight = 400
)

func (g *Game) Update() error {
	g.ui.Update()

	for _, creature := range g.Creatures {
		creature.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 100, 100, 255})
	opts := ebiten.DrawImageOptions{}
	screen.DrawImage(g.background, &opts)

	for _, creature := range g.Creatures {
		creature.Draw(screen)
	}
	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {

	g := &Game{}

	g.eventHub = events.NewEventHub()
	g.eventHub.Subscribe(events.ButtonClickedEvent{}, func(e events.Event) {
		ev := e.(events.ButtonClickedEvent)
		switch ev.ButtonText {
		case "save":
			g.SaveGame()
		}
	})
	g.background = LoadImageAssetAsEbitenImage("fishTank")
	g.loaded = true
	g.ui = ui.LoadMenu(screenWidth, screenHeight)

	ebiten.SetWindowSize(screenWidth, screenHeight)

	firstFish := NewFish()

	g.Creatures = append(g.Creatures, firstFish)

	return g
}

func (g *Game) SaveGame() {

}

func LoadImageAssetAsEbitenImage(assetName string) *ebiten.Image {
	imgPath := fmt.Sprintf("images/%s.png", assetName)
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, imgPath)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
