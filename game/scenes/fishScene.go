package scenes

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/fishTankWebGame/game"
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/graphicManagerSubscriptions"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/input"
	"github.com/acoco10/fishTankWebGame/game/interactableUIObjects"
	"github.com/acoco10/fishTankWebGame/game/loaders"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
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
	background        *ebiten.Image
	loaded            bool
	ImageMap          map[string]*ebiten.Image
	Creatures         []*entities.Creature
	particles         []*entities.Particle
	tankSize          image.Rectangle
	gameMode          gameMode
	sprites           []drawables.DrawableSprite
	debugRect         *geometry.Rect
	ui                *ebitenui.UI
	gameLog           *sceneManagement.GameLog
	handledClick      bool
	timers            map[string]*entities.Timer
	graphicManagerMap map[string]*graphics.GraphicManager
	offScreen         *ebiten.Image
	shader            *ebiten.Shader
	shaderParams      map[string]any
}

const (
	ScreenWidth  = 940
	ScreenHeight = 593
)

func NewFishScene(gameLog *sceneManagement.GameLog) *FishScene {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	backGroundImgShelfHeight := 124

	println("initiating game in ebiten NewFishScene()")

	g := &FishScene{}
	g.gameMode = Normal
	g.gameLog = gameLog
	g.ImageMap = make(map[string]*ebiten.Image)

	shader := shaders.LoadOnePointLightingBlue()
	shaderParams := make(map[string]any)

	g.offScreen = ebiten.NewImage(ScreenWidth, ScreenHeight)
	g.shader = shader
	g.shaderParams = shaderParams

	g.LoadTimers()

	g.Creatures = []*entities.Creature{}

	collisionMap, err := geometry.LoadCollisions()
	if err != nil {
		log.Fatal(err)
	}

	subs(g, collisionMap)

	err = g.loadBackground()
	if err != nil {
		log.Fatal(err)
	}

	g.debugRect = &geometry.Rect{}
	g.debugRect.RectState = geometry.Off

	fishSceneUISprites := []interactableUIObjects.UISpriteLabel{interactableUIObjects.FishBook, interactableUIObjects.Records, interactableUIObjects.FishFood, interactableUIObjects.WhiteBoard, interactableUIObjects.Plant}

	uiSprites, err := loaders.LoadUISprites(fishSceneUISprites, gameLog.GlobalEventHub, ScreenWidth, ScreenHeight)

	if err != nil {
		log.Fatal(err)
	}

	g.sprites = uiSprites

	fishTankSizeX1 := float64(g.ImageMap["fishTank"].Bounds().Max.X)
	fishTankSizeY1 := float64(g.ImageMap["fishTank"].Bounds().Max.Y)

	g.shaderParams["ImgRect"] = [4]float64{0, 0, fishTankSizeX1, fishTankSizeY1}
	g.shaderParams["LightPoint"] = [2]float64{420, 150}

	tankX := g.ImageMap["fishTank"].Bounds().Max.X
	tankY := g.ImageMap["fishTank"].Bounds().Max.Y

	startingX := int(ScreenWidth * 0.2)
	startingY := ScreenHeight - backGroundImgShelfHeight - g.ImageMap["fishTank"].Bounds().Dy()

	tankRect := image.Rect(startingX, startingY, tankX+startingX, tankY+startingY)

	g.tankSize = tankRect
	g.loaded = true

	mainUI, _, err := ui.LoadMainFishMenu(ScreenWidth, ScreenHeight, gameLog.GlobalEventHub)
	if err != nil {
		log.Fatal("error loading scene")
	}

	g.ui = mainUI

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

	collisionMap, err := geometry.LoadCollisions()
	if err != nil {
		log.Fatal(err)
	}

	g.gameLog = gameLog

	println("length of game log save = ", len(g.gameLog.Save.Fish))
	fishes := g.gameLog.Save

	for _, fish := range fishes.Fish {
		loadedFish := loaders.NewFish(gameLog.GlobalEventHub, collisionMap["tank"], fish)
		g.Creatures = append(g.Creatures, loadedFish)
	}
	g.LoadGraphicManagerMap()
	g.timers["songTimer"].TurnOn()
}

func (g *FishScene) LoadTimers() {
	g.timers = make(map[string]*entities.Timer)

	g.timers["pointGeneratedTimer"] = entities.NewTimer(0.2)
	g.timers["pointGeneratedTimer"].TurnOn()

	g.timers["songTimer"] = entities.NewTimer(15)
}

func (g *FishScene) IsLoaded() bool {
	return g.loaded
}

func (g *FishScene) Update() (sceneManagement.SceneId, error) {
	g.handledClick = false

	for _, creature := range g.Creatures {
		creature.Update()
	}

	for _, particle := range g.particles {
		particle.Update()
	}

	for _, gSprite := range g.sprites {
		gSprite.Update()
	}

	g.gameLog.SoundPlayer.Update()

	for _, graphicMan := range g.graphicManagerMap {
		graphicMan.Update()
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
	g.UpdateTimers()
	return sceneManagement.FishTank, nil
}

func (g *FishScene) UpdateTimers() {
	for key, timer := range g.timers {
		state := timer.Update()
		if key == "songTimer" && state == entities.Done {
			timer.TurnOff()
		}
	}
}

func (g *FishScene) DrawOffScreen() {
	opts := ebiten.DrawImageOptions{}

	g.offScreen.DrawImage(g.ImageMap["background"], &opts)

	opts.GeoM.Reset()
	opts.GeoM.Translate(float64(g.tankSize.Min.X), float64(g.tankSize.Min.Y))
	g.offScreen.DrawImage(g.ImageMap["fishTank"], &opts)

	for _, particle := range g.particles {
		particle.Draw(g.offScreen)
	}
	for _, creature := range g.Creatures {
		creature.Draw(g.offScreen)
	}
	for _, s := range g.sprites {
		s.Draw(g.offScreen)
	}

	opts.GeoM.Reset()

	fishTankFrontLayerDy := g.ImageMap["fishTankFrontLayer"].Bounds().Dy()
	fishTankHeightDy := g.ImageMap["fishTank"].Bounds().Dy()

	y := fishTankFrontLayerDy - fishTankHeightDy
	y = g.tankSize.Min.Y - y

	opts.GeoM.Translate(float64(g.tankSize.Min.X), float64(y))
	g.offScreen.DrawImage(g.ImageMap["fishTankFrontLayer"], &opts)
	opts.GeoM.Reset()

	g.offScreen.DrawImage(g.ImageMap["frontLayer"], &opts)

	//g.offScreen.DrawImage(g.fishTankFrontLayer, &opts)
}

func (g *FishScene) Draw(screen *ebiten.Image) {
	g.DrawOffScreen()
	ShaderOpts := &ebiten.DrawRectShaderOptions{}
	ShaderOpts.Images[0] = g.offScreen
	ShaderOpts.Uniforms = g.shaderParams
	opts := ebiten.DrawImageOptions{}
	screen.DrawRectShader(ScreenWidth, ScreenHeight, g.shader, ShaderOpts)
	opts.GeoM.Reset()
	//g.debugRect.Draw(screen)
	//g.printGameMode(screen)

	for _, s := range g.sprites {
		s.Draw(screen)
	}

	for _, graphicMan := range g.graphicManagerMap {
		graphicMan.Draw(screen)
	}

	g.ui.Draw(screen)
	for _, creature := range g.Creatures {
		if creature.Shader != nil {
			creature.Draw(screen)
		}
	}
}

func (g *FishScene) LoadGraphicManagerMap() {
	g.graphicManagerMap = make(map[string]*graphics.GraphicManager)
	WhiteBoardGraphicManager := graphics.NewGraphicManager(g.gameLog.GlobalEventHub, graphicManagerSubscriptions.WhiteBoardGMSubs)
	g.graphicManagerMap["whiteBoard"] = WhiteBoardGraphicManager

}

func (g *FishScene) Layout(outsideWidth, outsideHeight int) (int, int) {
	x, y := ebiten.WindowSize()
	if x > 0 && y > 0 {
		return ebiten.WindowSize()
	}
	return 940, 593
}

func subs(g *FishScene, colMap map[string]geometry.Rect) {

	g.gameLog.GlobalEventHub.Subscribe(ui.ButtonClickedEvent{}, func(e events.Event) {
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

	g.gameLog.GlobalEventHub.Subscribe(graphics.DrawGraphic{}, func(e events.Event) {
		ev := e.(graphics.DrawGraphic)
		println("adding draw graphic to game struct")
		g.sprites = append(g.sprites, ev.Graphic)
	})

	g.soundSubs()
	g.creatureSubs(colMap)

}

func (g *FishScene) soundSubs() {
	g.gameLog.GlobalEventHub.Subscribe(entities.SendData{}, func(e events.Event) {
		ev := e.(entities.SendData)
		if ev.DataFor == "soundFx" && ev.Data == "particle entered water" {
			g.gameLog.SoundPlayer.AddToQueue(soundFX.PlopSound)
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(sprite.UISpriteAction{}, func(e events.Event) {
		ev := e.(sprite.UISpriteAction)
		if ev.UiSprite == "fishFood" && ev.UiSpriteAction == "put back" {
			g.gameLog.SoundPlayer.AddToQueue(soundFX.PickUpOne)
		}
		if ev.UiSprite == "fishFood" && ev.UiSpriteAction == "picked up" {
			g.gameLog.SoundPlayer.AddToQueue(soundFX.SelectSound2)
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(events.TaskCompleted{}, func(e events.Event) {
		g.gameLog.SoundPlayer.AddToQueue(soundFX.SuccessMusic)
	})
}

func (g *FishScene) creatureSubs(colMap map[string]geometry.Rect) {
	g.gameLog.GlobalEventHub.Subscribe(input.MouseButtonPressed{}, func(e events.Event) {
		ev := e.(input.MouseButtonPressed)
		xCheck := ev.Point.X > float32(g.tankSize.Min.X)+100 && ev.Point.X < float32(g.tankSize.Max.X)
		yCheck := ev.Point.Y < float32(g.tankSize.Min.Y)-20

		if xCheck && yCheck && g.timers["pointGeneratedTimer"].TimerState == entities.Done && !g.handledClick {
			g.handledClick = true
			pt := ev.Point.Clone()
			pt.X = pt.X - 50 + rand.Float32()*10
			pt.Y += 50
			p := entities.NewParticle(pt, colMap["tank"], g.gameLog.GlobalEventHub)
			g.particles = append(g.particles, p)
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(entities.CreatureReachedPoint{}, func(e events.Event) {
		ev := e.(entities.CreatureReachedPoint)
		for i, p := range g.particles {
			if p.Point == ev.Point {
				g.particles = append(g.particles[:i], g.particles[i+1:]...)
			}
		}
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

	g.ImageMap["fishTankFrontLayer"] = fishTankFrontLayer
	g.ImageMap["fishTank"] = fishTankImg
	g.ImageMap["frontLayer"] = frontLayer
	g.ImageMap["background"] = background

	g.background = background
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
