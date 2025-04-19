package game

import (
	cursorUpdater "fishTankWebGame/game/cursor"
	"fishTankWebGame/game/gameEntities"
	"fishTankWebGame/game/helperFunc"
	"fishTankWebGame/game/ui"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"math/rand/v2"
)

type Game struct {
	img        *ebiten.Image
	background *ebiten.Image
	loaded     bool
	Creatures  []*gameEntities.Creature
	ui         *ebitenui.UI
	eventHub   *gameEntities.EventHub
	particles  []*gameEntities.Particle
	tankSize   image.Rectangle
	counter    int
}

const (
	screenWidth  = 640
	screenHeight = 480
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
			ev := gameEntities.MouseButtonPressed{
				Point: &gameEntities.Point{X: float32(x), Y: float32(y), PType: gameEntities.Food},
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
	return ebiten.WindowSize()
}

func NewGame(gameState int) *Game {

	g := &Game{}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	g.eventHub = gameEntities.NewEventHub()

	g.eventHub.Subscribe(gameEntities.ButtonClickedEvent{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.ButtonClickedEvent)
		switch ev.ButtonText {
		case "save":
			g.SaveGame()
		}
	})

	g.eventHub.Subscribe(gameEntities.MouseButtonPressed{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.MouseButtonPressed)
		x := rand.Float32() * 100
		ev.Point.X = ev.Point.X + x
		p := gameEntities.NewParticle(ev.Point)
		g.particles = append(g.particles, &p)
	})

	g.eventHub.Subscribe(gameEntities.CreatureReachedPoint{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.CreatureReachedPoint)
		for i, p := range g.particles {
			if p.Point == ev.Point {
				g.particles = append(g.particles[:i], g.particles[i+1:]...)
			}
		}
	})

	g.background = helperFunc.LoadImageAssetAsEbitenImage("fishTank")

	tankX := g.background.Bounds().Max.X
	tankY := g.background.Bounds().Max.Y

	startingX := (screenWidth - tankX) / 2
	startingY := (screenHeight - tankY) / 2

	tankRect := image.Rect(startingX, startingY, tankX+startingX, tankY+startingY)

	g.tankSize = tankRect
	g.loaded = true
	g.ui = ui.LoadMenu(screenWidth, screenHeight, g.eventHub)

	ebiten.SetWindowSize(screenWidth, screenHeight)

	firstFish := gameEntities.NewFish(g.eventHub, g.tankSize)
	secondFish := gameEntities.NewFish(g.eventHub, g.tankSize)

	g.Creatures = append(g.Creatures, firstFish, secondFish)

	return g
}

func (g *Game) SaveGame() {
	println("save game event generated and recieved")
	var saveData string
	for _, creature := range g.Creatures {
		saveData = saveData + "current size =" + string(rune(creature.Size))
	}
	SaveToBackend(saveData)
}
