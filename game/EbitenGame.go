package game

import (
	"encoding/json"
	cursorUpdater "fishTankWebGame/game/cursor"
	"fishTankWebGame/game/gameEntities"
	"fishTankWebGame/game/soundFx"
	"fishTankWebGame/game/ui"
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
)

type gameMode uint8

const (
	Position gameMode = iota
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
	frontLayer  *ebiten.Image
	gameMode
	sprites []*gameEntities.UiSprite
	*gameEntities.XYUpdater
	ffCursor *cursorUpdater.CursorUpdater
	*soundFX.SongPlayer
}

const (
	screenWidth  = 940
	screenHeight = 593
)

func (g *Game) Update() error {
	g.counter++

	for _, creature := range g.Creatures {
		creature.Update()
	}
	for _, particle := range g.particles {
		particle.Update()
	}

	for _, sprite := range g.sprites {
		sprite.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{50, 100, 100, 255})
	opts := ebiten.DrawImageOptions{}
	screen.DrawImage(g.background, &opts)
	opts.GeoM.Reset()
	opts.GeoM.Translate(float64(g.tankSize.Min.X), float64(g.tankSize.Min.Y))
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
	opts.GeoM.Reset()
	screen.DrawImage(g.frontLayer, &opts)
	for _, s := range g.sprites {
		if s.SpriteHovered() {
			s.Draw(screen)
		}
	}

	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	x, y := ebiten.WindowSize()
	if x > 0 && y > 0 {
		return ebiten.WindowSize()
	}
	return 940, 593
}

func NewGame(fishes gameEntities.SaveGameState) *Game {

	var positions gameEntities.SavePositionData
	spritePosition, err := os.ReadFile("../assets/data/spritePosition.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(spritePosition, &positions)

	backGroundImgShelfHeigth := 124

	println("initiating game in ebiten NewGame()")
	g := &Game{}

	for i, fish := range fishes.Fish {
		println("saved fish: ", i, "size: ", fish.Size)
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g.eventHub = gameEntities.NewEventHub()

	g.background = gameEntities.LoadImageAssetAsEbitenImage("roomBackground")
	g.fishTankImg = gameEntities.LoadImageAssetAsEbitenImage("fishTank")
	g.frontLayer = gameEntities.LoadImageAssetAsEbitenImage("frontLayer")
	var ffImgs []*ebiten.Image

	fishFoodImg := gameEntities.LoadImageAssetAsEbitenImage("fishFoodCursor")
	outlineFishFoodImg := gameEntities.LoadImageAssetAsEbitenImage("fishFoodOutline")
	altFishFoodImg := gameEntities.LoadImageAssetAsEbitenImage("fishFoodAlt")

	ffImgs = append(ffImgs, fishFoodImg, outlineFishFoodImg, altFishFoodImg)
	ffSprite := gameEntities.NewUiSprite(ffImgs, g.eventHub, positions.X, positions.Y, "fishFood")

	g.sprites = append(g.sprites, ffSprite)
	tankX := g.fishTankImg.Bounds().Max.X
	tankY := g.fishTankImg.Bounds().Max.Y

	startingX := (screenWidth - tankX) / 2
	startingY := screenHeight - backGroundImgShelfHeigth - g.fishTankImg.Bounds().Dy()

	tankRect := image.Rect(startingX, startingY, tankX+startingX, tankY+startingY)

	g.tankSize = tankRect
	g.loaded = true
	g.ui = ui.LoadMenu(screenWidth, screenHeight, g.eventHub)

	ebiten.SetWindowSize(screenWidth, screenHeight)

	firstFish := gameEntities.NewFish(g.eventHub, g.tankSize, fishes.Fish[0])
	secondFish := gameEntities.NewFish(g.eventHub, g.tankSize, fishes.Fish[1])

	g.Creatures = append(g.Creatures, firstFish, secondFish)

	g.eventHub.Subscribe(gameEntities.ButtonClickedEvent{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.ButtonClickedEvent)
		switch ev.ButtonText {
		case "Save":
			g.SaveGame()
		case "Mode":
			g.SwitchGameMode()
		}
	})

	//eventhub subscriptions
	g.eventHub.Subscribe(gameEntities.MouseButtonPressed{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.MouseButtonPressed)
		xCheck := ev.Point.X > float32(g.tankSize.Min.X) && ev.Point.X < float32(g.tankSize.Max.X)
		yCheck := ev.Point.Y < float32(g.tankSize.Min.Y)-20

		if xCheck && yCheck {
			x := rand.Float32() * 100
			ev.Point.X = ev.Point.X + x
			p := gameEntities.NewParticle(ev.Point)
			g.particles = append(g.particles, &p)

			pointEvent := gameEntities.PointGenerated{Point: ev.Point}
			g.eventHub.Publish(pointEvent)
		}
	})

	g.eventHub.Subscribe(gameEntities.CreatureReachedPoint{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.CreatureReachedPoint)
		for i, p := range g.particles {
			if p.Point == ev.Point {
				g.particles = append(g.particles[:i], g.particles[i+1:]...)
			}
		}
	})

	g.ui.Update()
	songPlayer, err := soundFX.NewSongPlayer()
	if err != nil {
		log.Fatal(err)
	}
	songPlayer.Play(soundFX.JazzE)

	return g
}

func (g *Game) SaveGame() {
	println("save game event generated and received")
	var savedFish []gameEntities.SavedFish

	for _, creature := range g.Creatures {
		f := gameEntities.GameFishToSaveFish(creature)
		savedFish = append(savedFish, f)
	}
	jsonSaveData, err := json.Marshal(savedFish)
	if err != nil {
		fmt.Println("Error marshaling:", err)
		return
	}

	SaveToBackend(string(jsonSaveData))
}

func (g *Game) SwitchGameMode() {
	switch g.gameMode {
	case Normal:
		g.gameMode = Position
	case Position:
		g.gameMode = Normal
	}
}
