package game

import (
	"encoding/json"
	"fishTankWebGame/game/debug"
	"fishTankWebGame/game/gameEntities"
	"fishTankWebGame/game/soundFx"
	"fishTankWebGame/game/ui"
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
	statMenu *ui.TextBoxUi
	sprites  []gameEntities.DrawableSprite
	*soundFX.SongPlayer
	debugRect *debug.Rect
}

const (
	ScreenWidth  = 940
	ScreenHeight = 593
)

func NewGame(fishes gameEntities.SaveGameState) *Game {

	backGroundImgShelfHeight := 124

	println("initiating game in ebiten NewGame()")

	g := &Game{}

	collisionMap, err := debug.LoadCollisions()
	if err != nil {
		log.Fatal(err)
	}
	err = g.loadBackground()
	if err != nil {
		log.Fatal(err)
	}

	for i, fish := range fishes.Fish {
		println("saved fish: ", i, "size: ", fish.Size)
	}

	if len(fishes.Fish) == 0 {

	}

	g.debugRect = &debug.Rect{}
	g.debugRect.RectState = debug.Off

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g.eventHub = gameEntities.NewEventHub()

	uiSprites, err := gameEntities.LoadUISprites(*g.eventHub, ScreenWidth, ScreenHeight)

	if err != nil {
		log.Fatal(err)
	}
	g.sprites = uiSprites

	tankX := g.fishTankImg.Bounds().Max.X
	tankY := g.fishTankImg.Bounds().Max.Y

	startingX := (ScreenWidth - tankX) / 2
	startingY := ScreenHeight - backGroundImgShelfHeight - g.fishTankImg.Bounds().Dy()

	tankRect := image.Rect(startingX, startingY, tankX+startingX, tankY+startingY)

	g.tankSize = tankRect
	g.loaded = true
	g.ui = ui.LoadMenu(ScreenWidth, ScreenHeight, g.eventHub)

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

	for _, fish := range fishes.Fish {
		loadedFish := gameEntities.NewFish(g.eventHub, collisionMap["tank"], fish)
		g.Creatures = append(g.Creatures, loadedFish)
	}

	g.ui.Update()
	g.gameMode = Normal
	songPlayer, err := soundFX.NewSongPlayer()
	if err != nil {
		log.Fatal(err)
	}
	songPlayer.Play(soundFX.WaterBubbles)

	EventSubsGame(g.eventHub, g, collisionMap)

	txtMenu, err := ui.NewTextBlocKMenu(g.eventHub)
	if err != nil {
		log.Fatal(err)
	}
	g.statMenu = txtMenu

	return g
}

func (g *Game) FirstLoad(state gameEntities.SaveGameState) {

}

func (g *Game) OnEntry() {

}
func (g *Game) OnExit() {

}

func (g *Game) IsLoaded() {

}

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

	g.statMenu.Update()

	if g.gameMode == Position {
		if ebiten.IsKeyPressed(ebiten.KeyM) {
			g.debugRect.Init("clickableFishFood")
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			g.saveUISpritePositions()
		}
		g.debugRect.Update()
	}

	g.ui.Update()
	return nil
}

func (g *Game) printGameMode(screen *ebiten.Image) {
	if g.gameMode == Position {
		DebugText("Position Mode", screen)
	}
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
		s.Draw(screen)
	}
	g.statMenu.Draw(screen)
	g.debugRect.Draw(screen)
	g.printGameMode(screen)
	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	x, y := ebiten.WindowSize()
	if x > 0 && y > 0 {
		return ebiten.WindowSize()
	}
	return 940, 593
}

func EventSubsGame(hub *gameEntities.EventHub, g *Game, colMap map[string]debug.Rect) {
	g.eventHub.Subscribe(gameEntities.MouseButtonPressed{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.MouseButtonPressed)
		xCheck := ev.Point.X > float32(g.tankSize.Min.X)+100 && ev.Point.X < float32(g.tankSize.Max.X)
		yCheck := ev.Point.Y < float32(g.tankSize.Min.Y)-20

		if xCheck && yCheck {
			x := rand.Float32() * 10
			ev.Point.X = ev.Point.X - 50 + x
			ev.Point.Y += 50
			p := gameEntities.NewParticle(ev.Point, colMap["tank"])
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

	g.eventHub.Subscribe(gameEntities.ButtonClickedEvent{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.ButtonClickedEvent)
		switch ev.ButtonText {
		case "Save":
			g.SaveGame()
		case "Mode":
			g.SwitchGameMode()
		}
	})
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

func (g *Game) loadBackground() error {

	background, err := gameEntities.LoadImageAssetAsEbitenImage("roomBackground")
	if err != nil {
		return err
	}

	fishTankImg, err := gameEntities.LoadImageAssetAsEbitenImage("fishTank")
	if err != nil {
		return err
	}

	frontLayer, err := gameEntities.LoadImageAssetAsEbitenImage("frontLayer")
	if err != nil {
		return err
	}

	g.frontLayer = frontLayer
	g.background = background
	g.fishTankImg = fishTankImg
	return nil
}

func (g *Game) saveUISpritePositions() {

	spMap := make(map[string]gameEntities.SavePositionData)

	for _, sprite := range g.sprites {
		uiSprite, ok := sprite.(*gameEntities.UiSprite)
		if !ok {
			continue
		}

		spData := uiSprite.SavePosition()
		spMap[uiSprite.Label] = spData

		ffSrite, ok := sprite.(*gameEntities.FishFoodSprite)
		if !ok {
			continue
		}

		ffData := ffSrite.SavePosition()
		spMap[ffSrite.Label] = ffData
	}

	outputSave, err := json.Marshal(spMap)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("../assets/data/spritePosition.json", outputSave, 999)
	if err != nil {
		log.Fatal(err)
	}
}

func DebugText(debugText string, screen *ebiten.Image) {
	face, err := ui.LoadFont(24.0)
	if err != nil {
		log.Fatal(err)
	}

	dopts := text.DrawOptions{}

	dopts.GeoM.Translate(ScreenWidth/2-float64(len(debugText)*6), ScreenHeight/10)
	text.Draw(screen, debugText, face, &dopts)
	dopts.GeoM.Reset()
}
