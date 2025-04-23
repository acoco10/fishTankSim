package game

import (
	"encoding/json"
	cursorUpdater "fishTankWebGame/game/cursor"
	"fishTankWebGame/game/gameEntities"
	"fishTankWebGame/game/ui"
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/color"
	"math/rand/v2"
)

type cursorMode uint8

const (
	FishFood cursorMode = iota
	Normal
)

type Game struct {
	img         *ebiten.Image
	background  *ebiten.Image
	loaded      bool
	Creatures   []*gameEntities.Creature
	ui          *ebitenui.UI
	eventHub    *gameEntities.EventHub
	particles   []*gameEntities.Particle
	tankSize    image.Rectangle
	counter     int
	fishTankImg *ebiten.Image
	sprites     []*gameEntities.UiSprite
	cursorMode
	*gameEntities.XYUpdater
	ffCursor *cursorUpdater.CursorUpdater
}

const (
	screenWidth  = 784
	screenHeight = 520
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

	for _, sprite := range g.sprites {
		if sprite.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.XYUpdater = gameEntities.NewUpdater(sprite.Sprite)
		}
	}

	if g.XYUpdater != nil && g.XYUpdater.Sprite != nil {
		g.XYUpdater.Update()
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if g.counter%2 == 0 && g.cursorMode == FishFood {
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

	screen.Fill(color.RGBA{50, 100, 100, 255})
	opts := ebiten.DrawImageOptions{}
	screen.DrawImage(g.background, &opts)
	opts.GeoM.Reset()
	opts.GeoM.Translate(float64(g.tankSize.Min.X), float64(g.tankSize.Min.Y+6))
	screen.DrawImage(g.fishTankImg, &opts)
	for _, particle := range g.particles {
		particle.Draw(screen)
	}
	for _, creature := range g.Creatures {
		creature.Draw(screen)
	}
	for _, s := range g.sprites {
		s.Draw(screen)
	}

	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	x, y := ebiten.WindowSize()
	if x > 0 && y > 0 {
		return ebiten.WindowSize()
	}
	return 785, 520
}

func NewGame() *Game {
	println("inititating game in ebiten NewGame()")
	g := &Game{}

	println("loading data from hard coded value")
	saveData := gameEntities.LoadSaveJson("data/saveaidan.json")

	var fishes gameEntities.SaveGameState
	err := json.Unmarshal([]byte(saveData), &fishes)
	if err != nil {
		return nil
	}

	for i, sFish := range fishes.Fish {
		println("saved fish: ", i, "size: ", sFish.Size)
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	g.eventHub = gameEntities.NewEventHub()
	g.cursorMode = Normal
	g.eventHub.Subscribe(gameEntities.ButtonClickedEvent{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.ButtonClickedEvent)
		switch ev.ButtonText {
		case "Save":
			g.SaveGame()
		case "fish food":
			g.cursorMode = FishFood
			g.ffCursor = cursorUpdater.CreateCursorUpdater()
			input.SetCursorUpdater(g.ffCursor)
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

	g.background = gameEntities.LoadImageAssetAsEbitenImage("roomBackground")
	g.fishTankImg = gameEntities.LoadImageAssetAsEbitenImage("fishTank")
	fishFoodImg := gameEntities.LoadImageAssetAsEbitenImage("fishFoodCursor")

	ffSprite := gameEntities.UiSprite{&gameEntities.Sprite{fishFoodImg, 681, 405, 0, 0}}
	g.sprites = append(g.sprites, &ffSprite)
	tankX := g.fishTankImg.Bounds().Max.X
	tankY := g.fishTankImg.Bounds().Max.Y

	startingX := (screenWidth - tankX) / 2
	startingY := 74

	tankRect := image.Rect(startingX, startingY, tankX+startingX, tankY+startingY)

	g.tankSize = tankRect
	g.loaded = true
	g.ui = ui.LoadMenu(screenWidth, screenHeight, g.eventHub)

	ebiten.SetWindowSize(screenWidth, screenHeight)

	firstFish := gameEntities.NewFish(g.eventHub, g.tankSize, fishes.Fish[0])
	secondFish := gameEntities.NewFish(g.eventHub, g.tankSize, fishes.Fish[1])

	g.Creatures = append(g.Creatures, firstFish, secondFish)

	return g
}

func (g *Game) SaveGame() {
	println("save game event generated and received")
	jsonSaveData, err := json.Marshal(g.Creatures)
	if err != nil {
		fmt.Println("Error marshaling:", err)
		return
	}

	SaveToBackend(string(jsonSaveData))
}
