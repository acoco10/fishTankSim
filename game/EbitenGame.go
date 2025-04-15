package game

import (
	"fishTankWebGame/assets"
	"fishTankWebGame/ebitenToJs"
	cursorUpdater "fishTankWebGame/game/cursor"
	"fishTankWebGame/game/events"
	"fishTankWebGame/game/ui"
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
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
	particles  []*Particle
	tankSize   image.Rectangle
	counter    int
}

const (
	screenWidth  = 600
	screenHeight = 500
)

func (g *Game) Update() error {
	g.counter++
	g.ui.Update()

	for _, creature := range g.Creatures {
		creature.Update()
	}
	for _, particle := range g.particles {
		particle.Update()
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if g.counter%2 == 0 {
			x, y := ebiten.CursorPosition()
			ev := events.MouseButtonPressed{
				Point: image.Point{x, y},
			}
			println("publishing event for mouse click")
			g.eventHub.Publish(ev)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	cursor := cursorUpdater.CreateCursorUpdater()
	input.SetCursorUpdater(cursor)
	screen.Fill(color.RGBA{50, 100, 100, 255})
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(g.tankSize.Min.X), float64(g.tankSize.Min.Y))
	screen.DrawImage(g.background, &opts)
	for _, particle := range g.particles {
		particle.Draw(screen)
	}
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
	g.eventHub.Subscribe(events.MouseButtonPressed{}, func(e events.Event) {
		ev := e.(events.MouseButtonPressed)
		p := NewParticle(ev.Point.X, ev.Point.Y)
		g.particles = append(g.particles, &p)
	})
	g.eventHub.Subscribe(events.CreatureReachedPoint{}, func(e events.Event) {
		//
	})

	g.background = LoadImageAssetAsEbitenImage("fishTank")

	tankX := g.background.Bounds().Max.X
	tankY := g.background.Bounds().Max.Y

	startingX := (screenWidth - tankX) / 2
	startingY := (screenHeight - tankY) / 2

	tankRect := image.Rect(startingX, startingY, tankX+startingX, tankY+startingY)

	g.tankSize = tankRect
	g.loaded = true
	g.ui = ui.LoadMenu(screenWidth, screenHeight, g.eventHub)

	ebiten.SetWindowSize(screenWidth, screenHeight)

	firstFish := NewFish(g.eventHub, g.tankSize)
	secondFish := NewFish(g.eventHub, g.tankSize)

	g.Creatures = append(g.Creatures, firstFish, secondFish)

	return g
}

func (g *Game) SaveGame() {
	println("save game event generated and recieved")
	var saveData string
	for _, creature := range g.Creatures {
		saveData = saveData + "current size =" + string(rune(creature.Size))
	}
	ebitenToJs.SaveToBackend(saveData)
}

func LoadImageAssetAsEbitenImage(assetName string) *ebiten.Image {
	imgPath := fmt.Sprintf("images/%s.png", assetName)
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, imgPath)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
