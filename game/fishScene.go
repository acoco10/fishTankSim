package game

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/debug"
	"github.com/acoco10/fishTankWebGame/game/gameEntities"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/acoco10/fishTankWebGame/game/ui"
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
	Normal gameMode = iota
	Position
)

type FishScene struct {
	background  *ebiten.Image
	loaded      bool
	Creatures   []*gameEntities.Creature
	eventHub    *gameEntities.EventHub
	particles   []*gameEntities.Particle
	tankSize    image.Rectangle
	counter     int
	fishTankImg *ebiten.Image
	frontLayer  *ebiten.Image
	gameMode
	sprites   []gameEntities.DrawableSprite
	debugRect *debug.Rect
	ui        *ebitenui.UI
	*sceneManagement.GameLog
	songTimer *gameEntities.Timer
}

const (
	ScreenWidth  = 940
	ScreenHeight = 593
)

func NewFishScene(gameLog *sceneManagement.GameLog) *FishScene {
	backGroundImgShelfHeight := 124

	println("initiating game in ebiten NewFishScene()")

	g := &FishScene{}
	g.GameLog = gameLog
	g.Creatures = []*gameEntities.Creature{}
	collisionMap, err := debug.LoadCollisions()
	if err != nil {
		log.Fatal(err)
	}

	err = g.loadBackground()
	if err != nil {
		log.Fatal(err)
	}

	g.debugRect = &debug.Rect{}
	g.debugRect.RectState = debug.Off

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g.eventHub = gameLog.GlobalEventHub

	fishSceneUI := []gameEntities.UISpriteLabel{gameEntities.FishBook, gameEntities.Records, gameEntities.FishFood}

	uiSprites, err := gameEntities.LoadUISprites(fishSceneUI, g.eventHub, ScreenWidth, ScreenHeight)

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

	g.gameMode = Normal

	subs(g, collisionMap)

	mainUI, _, err := ui.LoadMainFishMenu(ScreenWidth, ScreenHeight, g.eventHub)
	if err != nil {
		log.Fatal("error loading scene")
	}

	g.ui = mainUI

	g.songTimer = gameEntities.NewTimer(15)

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

	return g
}

func (g *FishScene) FirstLoad(gameLog *sceneManagement.GameLog) {
	NewFishScene(gameLog)

}

func (g *FishScene) OnExit() {

}

func (g *FishScene) OnEnter(gameLog *sceneManagement.GameLog) {

	collisionMap, err := debug.LoadCollisions()
	if err != nil {
		log.Fatal(err)
	}
	g.GameLog = gameLog
	println("length of game log save = ", g.GameLog.Save)
	fishes := g.GameLog.Save

	for _, fish := range fishes.Fish {
		loadedFish := gameEntities.NewFish(g.eventHub, collisionMap["tank"], fish)
		g.Creatures = append(g.Creatures, loadedFish)
	}
	g.songTimer.TurnOn()
}

func (g *FishScene) IsLoaded() bool {
	return g.loaded
}

func (g *FishScene) Update() (sceneManagement.SceneId, error) {
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

	if g.gameMode == Position {
		if ebiten.IsKeyPressed(ebiten.KeyM) {
			g.debugRect.Init("tank")
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			g.saveUISpritePositions()
		}
		g.debugRect.Update()
	}

	g.ui.Update()
	timerState := g.songTimer.Update()

	if timerState == gameEntities.Done {
		g.GameLog.SongPlayer.Play(soundFX.Lounge)
		g.songTimer.TurnOff()
	}

	return sceneManagement.FishTank, nil
}

func (g *FishScene) Draw(screen *ebiten.Image) {

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

	g.debugRect.Draw(screen)
	g.printGameMode(screen)

	g.ui.Draw(screen)

}

func (g *FishScene) Layout(outsideWidth, outsideHeight int) (int, int) {
	x, y := ebiten.WindowSize()
	if x > 0 && y > 0 {
		return ebiten.WindowSize()
	}
	return 940, 593
}

func subs(g *FishScene, colMap map[string]debug.Rect) {
	g.eventHub.Subscribe(gameEntities.MouseButtonPressed{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.MouseButtonPressed)
		xCheck := ev.Point.X > float32(g.tankSize.Min.X)+100 && ev.Point.X < float32(g.tankSize.Max.X)
		yCheck := ev.Point.Y < float32(g.tankSize.Min.Y)-20

		if xCheck && yCheck && g.counter%19 == 0 {
			x := rand.Float32() * 10
			ev.Point.X = ev.Point.X - 50 + x
			ev.Point.Y += 50
			p := gameEntities.NewParticle(ev.Point, colMap["tank"], g.eventHub)
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
		println(ev.ButtonText, "button event received")
		switch ev.ButtonText {
		case "Save":
			g.SaveGame()
		case "Mode":
			println("Mode button event received")
			g.SwitchGameMode()
		}
	})

	g.eventHub.Subscribe(gameEntities.SendData{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.SendData)
		if ev.DataFor == "soundFx" && ev.Data == "particle entered water" {
			g.GameLog.SoundPlayer.Play(soundFX.PlopSound)
		}
	})

	g.eventHub.Subscribe(gameEntities.UISpriteAction{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.UISpriteAction)
		if ev.UiSprite == "fishFood" && ev.UiSpriteAction == "put back" {
			g.GameLog.SoundPlayer.Play(soundFX.PickUpOne)
		}
		if ev.UiSprite == "fishFood" && ev.UiSpriteAction == "picked up" {
			g.GameLog.SoundPlayer.Play(soundFX.SelectSound2)
		}
	})
}

func (g *FishScene) SaveGame() {
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

func (g *FishScene) SwitchGameMode() {
	switch g.gameMode {
	case Normal:
		g.gameMode = Position
	case Position:
		g.gameMode = Normal
	}
}

func (g *FishScene) printGameMode(screen *ebiten.Image) {
	switch g.gameMode {
	case Position:
		DebugText("Position Mode", screen)
	case Normal:
		DebugText("Normal  Mode", screen)
	}
}

func DebugText(debugText string, screen *ebiten.Image) {
	face, err := ui.LoadFont(24.0)
	if err != nil {
		log.Fatal(err)
	}

	dOpts := text.DrawOptions{}
	dOpts.GeoM.Translate(ScreenWidth/2-float64(len(debugText)*6), ScreenHeight/10)
	text.Draw(screen, debugText, face, &dOpts)
	dOpts.GeoM.Reset()
}

func (g *FishScene) loadBackground() error {

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

func (g *FishScene) saveUISpritePositions() {

	spMap := make(map[string]gameEntities.SavePositionData)

	for _, sprite := range g.sprites {

		uiSprite, ok := sprite.(*gameEntities.UiSprite)
		if !ok {
			continue
		}

		spData := uiSprite.SavePosition()
		spMap[uiSprite.Label] = spData

		ffSprite, ok := sprite.(*gameEntities.FishFoodSprite)
		if !ok {
			continue
		}

		ffData := ffSprite.SavePosition()
		spMap[ffSprite.Label] = ffData
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
