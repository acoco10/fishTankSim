package scenes

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/fishTankWebGame/game"
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/input"
	"github.com/acoco10/fishTankWebGame/game/interactableUIObjects"
	"github.com/acoco10/fishTankWebGame/game/loaders"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/acoco10/fishTankWebGame/game/sprite"
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

var nextID int

type FishScene struct {
	background          *ebiten.Image
	loaded              bool
	Creatures           []*entities.Creature
	eventHub            *events.EventHub
	particles           []*entities.Particle
	tankSize            image.Rectangle
	fishTankImg         *ebiten.Image
	frontLayer          *ebiten.Image
	fishTankFrontLayer  *ebiten.Image
	gameMode            gameMode
	sprites             []drawables.DrawableSprite
	debugRect           *geometry.Rect
	ui                  *ebitenui.UI
	gameLog             *sceneManagement.GameLog
	songTimer           *entities.Timer
	task                *events.Task
	graphics            []drawables.DrawableSprite
	handledClick        bool
	pointGeneratedTimer *entities.Timer
	graphicManager      *graphics.GraphicManager
}

const (
	ScreenWidth  = 940
	ScreenHeight = 593
)

func NewFishScene(gameLog *sceneManagement.GameLog) *FishScene {
	backGroundImgShelfHeight := 124

	println("initiating game in ebiten NewFishScene()")

	g := &FishScene{}
	g.pointGeneratedTimer = entities.NewTimer(0.2)
	g.pointGeneratedTimer.TurnOn()
	g.gameLog = gameLog
	g.Creatures = []*entities.Creature{}
	collisionMap, err := geometry.LoadCollisions()
	if err != nil {
		log.Fatal(err)
	}

	err = g.loadBackground()
	if err != nil {
		log.Fatal(err)
	}

	g.debugRect = &geometry.Rect{}
	g.debugRect.RectState = geometry.Off

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g.eventHub = gameLog.GlobalEventHub

	fishSceneUISprites := []interactableUIObjects.UISpriteLabel{interactableUIObjects.FishBook, interactableUIObjects.Records, interactableUIObjects.FishFood, interactableUIObjects.WhiteBoard, interactableUIObjects.Plant}

	uiSprites, err := loaders.LoadUISprites(fishSceneUISprites, g.eventHub, ScreenWidth, ScreenHeight)

	if err != nil {
		log.Fatal(err)
	}

	g.sprites = uiSprites

	tankX := g.fishTankImg.Bounds().Max.X
	tankY := g.fishTankImg.Bounds().Max.Y

	startingX := int(ScreenWidth * 0.2)
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

	g.songTimer = entities.NewTimer(15)

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

	return g
}

func (g *FishScene) FirstLoad(gameLog *sceneManagement.GameLog) {
	NewFishScene(gameLog)

}

func (g *FishScene) OnExit() {

}

func (g *FishScene) OnEnter(gameLog *sceneManagement.GameLog) {

	for n, task := range gameLog.Tasks {
		fmt.Printf("loading task: %d,:%s", n, task.Text)
		task.Publish(gameLog.GlobalEventHub)
	}

	g.graphicManager = graphics.NewGraphicManager(g.eventHub)

	collisionMap, err := geometry.LoadCollisions()
	if err != nil {
		log.Fatal(err)
	}

	g.gameLog = gameLog

	println("length of game log save = ", len(g.gameLog.Save.Fish))
	fishes := g.gameLog.Save

	for _, fish := range fishes.Fish {
		loadedFish := loaders.NewFish(g.eventHub, collisionMap["tank"], fish)
		g.Creatures = append(g.Creatures, loadedFish)
	}

	g.songTimer.TurnOn()
}

func (g *FishScene) IsLoaded() bool {
	return g.loaded
}

func (g *FishScene) Update() (sceneManagement.SceneId, error) {
	g.handledClick = false

	g.pointGeneratedTimer.Update()

	for _, creature := range g.Creatures {
		creature.Update()
	}

	for _, particle := range g.particles {
		particle.Update()
	}

	for _, sprite := range g.sprites {
		sprite.Update()
	}

	g.gameLog.SoundPlayer.Update()

	g.graphicManager.Update()

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

	if timerState == entities.Done {
		//g.gameLog.SoundPlayer.Play(soundFX.Lounge)
		g.songTimer.TurnOff()
	}

	for _, graph := range g.graphics {
		graph.Update()
	}

	return sceneManagement.FishTank, nil
}

func (g *FishScene) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{50, 100, 100, 255})
	opts := ebiten.DrawImageOptions{}
	screen.DrawImage(g.background, &opts)

	fisTankFrontLayerDy := g.fishTankFrontLayer.Bounds().Dy()
	fishTankHeightDy := g.fishTankImg.Bounds().Dy()

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
	opts.GeoM.Reset()

	y := fisTankFrontLayerDy - fishTankHeightDy
	y = g.tankSize.Min.Y - y

	opts.GeoM.Translate(float64(g.tankSize.Min.X), float64(y))

	screen.DrawImage(g.fishTankFrontLayer, &opts)

	for _, s := range g.sprites {
		s.Draw(screen)
	}

	for _, creature := range g.Creatures {
		if creature.Shader != nil {
			creature.Draw(screen)
		}
	}

	opts.GeoM.Reset()

	//g.debugRect.Draw(screen)
	//g.printGameMode(screen)

	g.ui.Draw(screen)
	g.graphicManager.Draw(screen)

}

func (g *FishScene) Layout(outsideWidth, outsideHeight int) (int, int) {
	x, y := ebiten.WindowSize()
	if x > 0 && y > 0 {
		return ebiten.WindowSize()
	}
	return 940, 593
}

func subs(g *FishScene, colMap map[string]geometry.Rect) {
	g.eventHub.Subscribe(input.MouseButtonPressed{}, func(e events.Event) {
		ev := e.(input.MouseButtonPressed)
		xCheck := ev.Point.X > float32(g.tankSize.Min.X)+100 && ev.Point.X < float32(g.tankSize.Max.X)
		yCheck := ev.Point.Y < float32(g.tankSize.Min.Y)-20

		if xCheck && yCheck && g.pointGeneratedTimer.TimerState == entities.Done && !g.handledClick {
			g.handledClick = true
			pt := ev.Point.Clone()
			pt.X = pt.X - 50 + rand.Float32()*10
			pt.Y += 50
			p := entities.NewParticle(pt, colMap["tank"], g.eventHub)
			g.particles = append(g.particles, p)
		}
	})

	g.eventHub.Subscribe(entities.CreatureReachedPoint{}, func(e events.Event) {
		ev := e.(entities.CreatureReachedPoint)
		for i, p := range g.particles {
			if p.Point == ev.Point {
				g.particles = append(g.particles[:i], g.particles[i+1:]...)
			}
		}
	})

	g.eventHub.Subscribe(ui.ButtonClickedEvent{}, func(e events.Event) {
		ev := e.(ui.ButtonClickedEvent)
		println(ev.ButtonText, "button event received")
		switch ev.ButtonText {
		case "Save":
			g.SaveGame()
		case "Mode":
			println("Mode button event received")
			g.SwitchGameMode()
		}
	})

	g.eventHub.Subscribe(entities.SendData{}, func(e events.Event) {
		ev := e.(entities.SendData)
		if ev.DataFor == "soundFx" && ev.Data == "particle entered water" {
			g.gameLog.SoundPlayer.AddToQueue(soundFX.PlopSound)
		}
	})

	g.eventHub.Subscribe(sprite.UISpriteAction{}, func(e events.Event) {
		ev := e.(sprite.UISpriteAction)
		if ev.UiSprite == "fishFood" && ev.UiSpriteAction == "put back" {
			g.gameLog.SoundPlayer.AddToQueue(soundFX.PickUpOne)
		}
		if ev.UiSprite == "fishFood" && ev.UiSpriteAction == "picked up" {
			g.gameLog.SoundPlayer.AddToQueue(soundFX.SelectSound2)
		}
	})

	g.eventHub.Subscribe(graphics.DrawGraphic{}, func(e events.Event) {
		ev := e.(graphics.DrawGraphic)
		println("adding draw graphic to game struct")
		g.graphics = append(g.graphics, ev.Graphic)
	})

	g.eventHub.Subscribe(events.TaskCompleted{}, func(e events.Event) {
		g.gameLog.SoundPlayer.AddToQueue(soundFX.SuccessMusic)
	})
}

func (g *FishScene) SaveGame() {
	println("save game event generated and received")
	var savedFish []entities.SavedFish

	for _, creature := range g.Creatures {
		f := entities.GameFishToSaveFish(creature)
		savedFish = append(savedFish, f)
	}

	save := entities.SaveGameState{}
	save.Fish = savedFish

	save.Tasks = g.gameLog.Tasks

	jsonSaveData, err := json.Marshal(save)

	println("save data before sending to js:", string(jsonSaveData))

	if err != nil {
		fmt.Println("Error marshalling:", err)
		return
	}

	/*err = os.WriteFile("../assets/data/saveWithTasks.json", jsonSaveData, 999)
	if err != nil {
		log.Fatal(err)
	}*/

	game.SaveToBackend(jsonSaveData)
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
	face, err := ui.LoadFont(24.0, "nk57")
	if err != nil {
		log.Fatal(err)
	}
	dOpts := text.DrawOptions{}
	dOpts.GeoM.Translate(ScreenWidth/2-float64(len(debugText)*6), ScreenHeight/10)
	text.Draw(screen, debugText, face, &dOpts)
	dOpts.GeoM.Reset()
}

func (g *FishScene) loadBackground() error {

	background, err := loaders.LoadImageAssetAsEbitenImage("roomBackground")
	if err != nil {
		return err
	}

	fishTankImg, err := loaders.LoadImageAssetAsEbitenImage("fishTank")
	if err != nil {
		return err
	}

	frontLayer, err := loaders.LoadImageAssetAsEbitenImage("frontLayer")
	if err != nil {
		return err
	}

	fishTankFrontLayer, err := loaders.LoadImageAssetAsEbitenImage("fishTankFrontLayer")
	if err != nil {
		return err
	}

	g.fishTankFrontLayer = fishTankFrontLayer
	g.frontLayer = frontLayer
	g.background = background
	g.fishTankImg = fishTankImg
	return nil
}

func (g *FishScene) saveUISpritePositions() {

	spMap := make(map[string]drawables.SavePositionData)

	for _, sprite := range g.sprites {

		spData := sprite.SavePosition()
		spMap[spData.Name] = spData

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
