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
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/props"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/acoco10/fishTankWebGame/game/system"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/tutorial"
	"github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"log"
	"math/rand"
	"os"
)

type DebugOption uint8

const (
	Normal DebugOption = iota
	Position
	Print
)

type GameMode uint8

const (
	User GameMode = iota
	Debug
)

var nextID int

type FishScene struct {
	background        *ebiten.Image
	backgroundShader  *ebiten.Shader
	backGroundParams  map[string]any
	propMap           map[string]*props.StructureProp
	loaded            bool
	ImageMap          map[string]*ebiten.Image
	particles         []*entities.Particle
	tankSize          image.Rectangle
	gameMode          GameMode
	debugParameter    map[DebugOption]bool
	sprites           []drawables.Drawable
	debugRect         *geometry.Rect
	ui                *ebitenui.UI
	gameLog           *sceneManagement.GameLog
	handledClick      bool
	timers            map[string]*entities.Timer
	graphicManagerMap map[string]*graphics.GraphicManager
	shader            *ebiten.Shader
	shaderParams      map[string]any
	returnScene       sceneManagement.SceneId
	cursorOccupied    bool
	debugText         string
	tutorialManager   *tutorial.Manager
	collisionMap      map[string]geometry.Rect
	store             system.Store
	playerState       *entities.Player
	images            gameImages
}

type gameImages struct {
	background *ebiten.Image
	offScreen  *ebiten.Image
	fishTank   *ebiten.Image
}

var backGroundImgShelfHeight = 248

func NewFishScene(gameLog *sceneManagement.GameLog) *FishScene {

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	//hardcoded ass parameter cus i aint about to read sum dark brown particles dawg

	println("initiating game in ebiten NewFishScene()")

	g := &FishScene{}

	//stuff that needs to exist before game publishes shit

	//g.gameMode = Position
	g.gameLog = gameLog
	g.ImageMap = make(map[string]*ebiten.Image)
	g.LoadGraphicManagerMap()

	g.shader = registry.ShaderMap["OnePointLighting"]
	g.shaderParams = make(map[string]any)

	g.ImageMap["offScreen"] = ebiten.NewImage(ScreenWidth, ScreenHeight)
	g.LoadTimers()

	collisionMap, err := geometry.LoadCollisions()
	if err != nil {
		log.Fatal(err)
	}

	g.collisionMap = collisionMap

	g.subs(collisionMap)

	err = g.loadBackground()
	if err != nil {
		log.Fatal(err)
	}

	g.debugRect = &geometry.Rect{}
	//g.debugRect.RectState = geometry.Off

	fishTankSizeX1 := float64(g.ImageMap["fishTank"].Bounds().Max.X)
	fishTankSizeY1 := float64(g.ImageMap["fishTank"].Bounds().Max.Y)

	g.shaderParams["ImgRect"] = [4]float64{0, 0, fishTankSizeX1, fishTankSizeY1}
	g.shaderParams["LightPoint"] = [2]float64{440, 250}

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
	fishSceneUISprites := []interactableUIObjects.Label{
		interactableUIObjects.FishBook,
		interactableUIObjects.Records,
		interactableUIObjects.FishFood,
		interactableUIObjects.WhiteBoard,
		interactableUIObjects.PiggyBank,
		interactableUIObjects.Pillow,
		interactableUIObjects.Thermometer,
		interactableUIObjects.Magazine,
	}

	uiSprites, err := loader.LoadUISprites(fishSceneUISprites, gameLog.GlobalEventHub, ScreenWidth, ScreenHeight)
	g.sprites = uiSprites

	return g
}

func LoadPurchasedSprite(inputName string, hub *tasks.EventHub, tankSize geometry.Rect) *entities.Creature {

	fData := entities.SavedFish{
		FishType: inputName,
		Size:     1,
	}

	creature := loader.NewFish(hub, tankSize, fData)

	return creature
}

func (g *FishScene) LoadStuff() {
	propMap := loader.LoadProps(g.gameLog.Save.TankObjects, g.tankSize, g.gameLog.GlobalEventHub)
	g.propMap = propMap

	store := system.NewStore(g.gameLog.GlobalEventHub)
	g.store = store

	g.playerState = &entities.Player{Money: 10, EventHub: g.gameLog.GlobalEventHub}
	g.playerState.Subscribe()

	loaderMan := loader.Manager{Hub: g.gameLog.GlobalEventHub}
	loaderMan.Subscriptions()

}

func (g *FishScene) DrawProps(screen *ebiten.Image) {
	for _, prop := range g.propMap {
		prop.Draw(screen)
	}
}

func (g *FishScene) UpdateProps() {
	for _, prop := range g.propMap {
		prop.Update()
	}
}

func (g *FishScene) FirstLoad(gameLog *sceneManagement.GameLog) {

}

func (g *FishScene) OnExit() {
}

func (g *FishScene) OnEnter(gameLog *sceneManagement.GameLog) {

	g.backgroundShader = registry.ShaderMap["NormalMap"]

	g.backGroundParams = make(map[string]any)
	g.backGroundParams["Cursor"] = []float64{440, 200}

	collisionMap, err := geometry.LoadCollisions()
	if err != nil {
		log.Fatal(err)
	}

	g.gameLog = gameLog

	if g.gameLog.Day == 1 {
		println("length of game log save =", len(g.gameLog.Save.Fish))
		fishes := g.gameLog.Save

		for _, fish := range fishes.Fish {
			loadedFish := loader.NewFish(gameLog.GlobalEventHub, collisionMap["tank"], fish)
			g.sprites = append(g.sprites, loadedFish)
		}
		g.LoadStuff()
	}

	ev2 := events.FishTankLayout{
		Rectangle: g.tankSize,
	}

	g.gameLog.GlobalEventHub.Publish(ev2)

	//No music on the base level as of now
	//g.timers["songTimer"].TurnOn()
	///JUST FOR TESTING, NO CHORE=> ALLOWANCE FRAMEWORK YET
	g.gameLog.GlobalEventHub.Publish(events.MoneyAvailable{Amount: 1})
	g.returnScene = sceneManagement.FishTank

	/*tutMngr := tutorial.Manager{}
	tutorial.InitData(&tutMngr, gameLog.GlobalEventHub)
	g.tutorialManager = &tutMngr*/

	ev := events.NewDay{NTasks: len(g.gameLog.Tasks)}
	g.gameLog.GlobalEventHub.Publish(ev)

	g.gameLog.Tasks[0].Activate()
	g.gameLog.Tasks[0].Publish(g.gameLog.GlobalEventHub)

}

func (g *FishScene) LoadTimers() {

	g.timers = make(map[string]*entities.Timer)
	g.timers["pointGeneratedTimer"] = entities.NewTimer(0.2)
	g.timers["pointGeneratedTimer"].TurnOn()
	g.timers["songTimer"] = entities.NewTimer(15)
	g.timers["sceneTransition"] = entities.NewTimer(2.5)

}

func (g *FishScene) IsLoaded() bool {
	return g.loaded
}

func (g *FishScene) Update() (sceneManagement.SceneId, error) {

	g.handledClick = false

	g.UpdateProps()

	g.gameLog.SoundPlayer.Update()

	for _, particle := range g.particles {
		particle.Update()
	}

	for _, gSprite := range g.sprites {
		gSprite.Update()
	}

	for _, graphicMan := range g.graphicManagerMap {
		graphicMan.Update()
	}

	if g.CheckIfAllFishFed() {
		ev2 := entities.AllFishFed{}
		g.gameLog.GlobalEventHub.Publish(ev2)
	}

	g.ui.Update()

	g.updateTimers()

	g.updateInput()

	//g.tutorialManager.Update()

	if g.gameMode == Debug {
		if g.debugParameter[Position] {
			g.positionModeUpdate()
		}
		err := g.debugRect.Update()
		if err != nil {
			//debug rect could error when saving collision location
			return g.returnScene, err
		}
	}

	return g.returnScene, nil
}

func (g *FishScene) DrawOffScreen() {

	opts := ebiten.DrawImageOptions{}
	shaderOpts := &ebiten.DrawRectShaderOptions{}

	shaderOpts.Uniforms = g.backGroundParams
	shaderOpts.Images[0] = g.ImageMap["fishTank"]
	shaderOpts.Images[1] = g.ImageMap["fishTankNormal"]

	g.ImageMap["offScreen"].DrawImage(g.ImageMap["background"], &opts)
	b := g.ImageMap["fishTank"].Bounds()

	opts.GeoM.Reset()

	shaderOpts.GeoM.Translate(float64(g.tankSize.Min.X), float64(g.tankSize.Min.Y))

	g.images.offScreen.DrawRectShader(b.Dx(), b.Dy(), g.backgroundShader, shaderOpts)

	g.DrawProps(g.images.offScreen)

	for _, particle := range g.particles {
		particle.Draw(g.images.offScreen)
	}

	for _, s := range g.sprites {
		s.Draw(g.images.offScreen)
	}

	opts.GeoM.Reset()
	fishTankFrontLayerDy := g.ImageMap["fishTankFrontLayer"].Bounds().Dy()
	fishTankHeightDy := g.ImageMap["fishTank"].Bounds().Dy()

	y := fishTankFrontLayerDy - fishTankHeightDy
	y = g.tankSize.Min.Y - y

	opts.GeoM.Translate(float64(g.tankSize.Min.X), float64(y))
	g.images.offScreen.DrawImage(g.ImageMap["fishTankFrontLayer"], &opts)
	opts.GeoM.Reset()
	g.images.offScreen.DrawImage(g.ImageMap["frontLayer"], &opts)

}

func (g *FishScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
	//draw objects receiving global lighting shader
	g.DrawOffScreen()

	//apply shader to offscreen
	ShaderOpts := &ebiten.DrawRectShaderOptions{}
	ShaderOpts.Images[0] = g.images.offScreen
	ShaderOpts.Uniforms = g.shaderParams

	screen.DrawRectShader(ScreenWidth, ScreenHeight, g.shader, ShaderOpts)

	g.debugRect.Draw(screen)

	//vector.StrokeCircle(screen, 440, 170, 2, 1, colornames.Yellow, false)

	/*for _, s := range g.sprites {
		s.Draw(screen)
	}*/

	for _, graphicMan := range g.graphicManagerMap {
		graphicMan.Draw(screen)
	}

	DebugText(g.debugText, screen)

	graphics.DrawGraphics(screen)
	g.ui.Draw(screen)

}

func (g *FishScene) positionModeUpdate() {
	if ebiten.IsKeyPressed(ebiten.KeyM) {
		g.debugRect.Init("tank")
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.saveUISpritePositions()
	}

	err := g.debugRect.Update()

	if err != nil {
		log.Printf("Couldn't save positions or something with the debug rect got fucked up")
	}
}

func (g *FishScene) updateTimers() {

	for key, timer := range g.timers {
		state := timer.Update()

		if key == "songTimer" && state == entities.Done {
			timer.TurnOff()
		}

		if key == "sceneTransition" && state == entities.Done {
			timer.TurnOff()
			g.returnScene = sceneManagement.TransitionScene
			g.gameLog.GlobalEventHub.Publish(events.DayOverTransitionComplete{})
		}
	}
}

func (g *FishScene) updateInput() {
	//function for handling ebiten input directly in game mode mainly for convenience
	//or avoiding the event system (latency?)
	//not necessarily core game functions
	g.checkForFishSelected()
	g.debugInputCheck()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.gameLog.GlobalEventHub.Publish(events.CloseWindow{})
	}

}

func (g *FishScene) checkForFishSelected() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		xCheck := x > g.tankSize.Min.X && x < g.tankSize.Max.X
		yCheck := y > g.tankSize.Min.Y && y < g.tankSize.Max.Y

		if xCheck && yCheck {
			filterFunc := func(distance any) bool {
				return distance.(float64) < 100
			}
			cursorX, cursorY := ebiten.CursorPosition()
			closestCreature := util.ClosestDrawableToCursor(cursorX, cursorY, g.sprites, filterFunc, "*entities.Creature")

			if closestCreature != nil {
				cre, ok := closestCreature.(*entities.Creature)
				if ok {
					SelectCreature(cre)
				}
			}
		}
	}
}

func (g *FishScene) debugInputCheck() {

	if ebiten.IsKeyPressed(ebiten.KeyF) {
		switch g.gameMode {
		case Debug:
			if g.debugParameter == nil {
				g.debugParameter = make(map[DebugOption]bool)
			}
			g.gameMode = User
		case User:
			g.gameMode = Debug
			ev := events.ButtonClickedEvent{ButtonText: "Mode"}
			g.gameLog.GlobalEventHub.Publish(ev)
			g.debugModeParameterPrinterUpdater()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		g.debugParameter[Print] = true
	}

	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.debugParameter[Position] = true
		g.debugModeParameterPrinterUpdater()
	}

	for key, dbp := range g.debugParameter {
		switch key {
		case Position:
			if dbp && ebiten.IsKeyPressed(ebiten.KeyS) {
				g.saveUISpritePositions()
			}
		}
	}

}

func SelectCreature(creature *entities.Creature) {
	println("selecting creature")
	creature.Selected = true
	creature.Shader = registry.ShaderMap["Outline"]
	loader.LoadRotatingHighlightOutlineAnimated(creature.AnimatedSprite)
}

func (g *FishScene) debugModeParameterPrinterUpdater() {
	g.debugText = "Debug Mode Activated| Parameters:"
	for key, dbp := range g.debugParameter {
		if dbp {
			switch key {
			case Position:
				g.debugText += "Position"
			}
		}
	}
}

func (g *FishScene) LoadGraphicManagerMap() {

	g.graphicManagerMap = make(map[string]*graphics.GraphicManager)
	WhiteBoardGraphicManager := graphics.NewGraphicManager(g.gameLog.GlobalEventHub, graphicManagerSubscriptions.WhiteBoardGMSubs)
	ScreenGraphicManager := graphics.NewGraphicManager(g.gameLog.GlobalEventHub, graphicManagerSubscriptions.ScreenGMSubs)

	g.graphicManagerMap["whiteBoard"] = WhiteBoardGraphicManager
	g.graphicManagerMap["screen"] = ScreenGraphicManager

}

func (g *FishScene) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth > 0 {
		return outsideWidth, outsideHeight
	}
	return ScreenWidth, ScreenHeight
}

func (g *FishScene) subs(colMap map[string]geometry.Rect) {

	g.uiSubs()
	g.soundSubs()
	g.creatureSubs(colMap)

}

func (g *FishScene) uiSubs() {
	g.gameLog.GlobalEventHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		println(ev.ButtonText, "button event received")
		switch ev.ButtonText {
		case "Save":
			g.SaveGame()
		case "Mode":
			println("Mode button event received")
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(graphics.DrawGraphic{}, func(e tasks.Event) {
		ev := e.(graphics.DrawGraphic)
		println("adding draw graphic to game struct")
		g.sprites = append(g.sprites, ev.Graphic)
	})

	g.gameLog.GlobalEventHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		if ev.ButtonText == "Go To Bed" {
			g.timers["sceneTransition"].TurnOn()
			ev2 := events.DayOver{}
			g.gameLog.GlobalEventHub.Publish(ev2)
			g.gameLog.Day++
		}
	})
}

func (g *FishScene) soundSubs() {
	g.gameLog.GlobalEventHub.Subscribe(entities.SendData{}, func(e tasks.Event) {
		ev := e.(entities.SendData)
		if ev.DataFor == "soundFx" && ev.Data == "particle entered water" {
			g.gameLog.SoundPlayer.AddToQueue(soundFX.PlopSound, 1)
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(events.UISpriteAction{}, func(e tasks.Event) {
		ev := e.(events.UISpriteAction)
		if ev.UiSprite == "fishFood" && ev.UiSpriteAction == "put back" {
			g.gameLog.SoundPlayer.AddToQueue(soundFX.PickUpOne, 1)
		}
		if ev.UiSprite == "fishFood" && ev.UiSpriteAction == "picked up" {
			g.gameLog.SoundPlayer.AddToQueue(soundFX.SelectSound2, 1)
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(tasks.TaskCompleted{}, func(e tasks.Event) {
		g.gameLog.SoundPlayer.AddToQueue(soundFX.WhiteBoardMarker1, 2)
		g.gameLog.SoundPlayer.AddToQueue(soundFX.SuccessMusic, 1)
		ev := e.(tasks.TaskCompleted)
		if len(g.gameLog.Tasks) > ev.Slot {
			g.gameLog.Tasks[ev.Slot].Activate()
			g.gameLog.Tasks[ev.Slot].Publish(g.gameLog.GlobalEventHub)
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
		g.gameLog.SoundPlayer.AddToQueue(soundFX.WhiteBoardMarker2, 2)
	})

}

func (g *FishScene) creatureSubs(colMap map[string]geometry.Rect) {
	g.gameLog.GlobalEventHub.Subscribe(input.MouseButtonPressedUISpriteActivity{}, func(e tasks.Event) {

		ev := e.(input.MouseButtonPressedUISpriteActivity)

		if g.timers["pointGeneratedTimer"].TimerState == entities.Done && !g.handledClick {
			g.handledClick = true
			pt := ev.Point.Clone()
			pt.X = pt.X - 50 + rand.Float32()*10
			pt.Y += 50
			p := entities.NewParticle(pt, colMap["tank"], g.gameLog.GlobalEventHub)
			g.sprites = append(g.sprites, p)
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(events.NewPurchase{}, func(e tasks.Event) {
		ev := e.(events.NewPurchase)
		log.Printf("New Purchase:%s ", ev.Purchase)
		creature := LoadPurchasedSprite(ev.Purchase, g.gameLog.GlobalEventHub, g.collisionMap["tank"])
		g.sprites = append(g.sprites, creature)
	})

	g.gameLog.GlobalEventHub.Subscribe(entities.CreatureReachedPoint{}, func(e tasks.Event) {
		ev := e.(entities.CreatureReachedPoint)
		if ev.Point != nil {
			for i, p := range g.particles {
				if p.Point == ev.Point {
					println("removing point to particle in g.particles")
					g.particles = append(g.particles[:i], g.particles[i+1:]...)
				}
			}
		}
	})

}

func (g *FishScene) CheckIfAllFishFed() bool {
	fed := true

	for _, draw := range g.sprites {
		creature, ok := draw.(*entities.Creature)
		if ok && creature.Hunger == 0 {
			fed = false
		}
	}
	return fed

}

func (g *FishScene) SaveGame() {
	println("save game event generated and received")
	var savedFish []entities.SavedFish

	for _, draw := range g.sprites {
		creature, ok := draw.(*entities.Creature)
		if ok && creature.Hunger == 0 {
			f := entities.GameFishToSaveFish(creature)
			savedFish = append(savedFish, f)
		}
	}

	save := entities.SaveGameState{}
	save.Fish = savedFish

	var savedTasks []map[string]any

	for _, task := range g.gameLog.Tasks {
		saveTask := make(map[string]any)

		saveTask["Name"] = task.Name
		saveTask["Completed"] = task.Completed
		saveTask["Text"] = task.Text

		savedTasks = append(savedTasks, saveTask)
	}

	save.Tasks = savedTasks

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

func (g *FishScene) printGameMode(screen *ebiten.Image) {
	DebugText(g.debugText, screen)
}

func DebugText(debugText string, screen *ebiten.Image) {
	face, err := util.LoadFont(24.0, "nk57")
	if err != nil {
		log.Fatal("Couldnt Load font for debug text", err)
	}
	dOpts := text.DrawOptions{}
	dOpts.GeoM.Translate(ScreenWidth/2-float64(len(debugText)*6), ScreenHeight/10)
	text.Draw(screen, debugText, face, &dOpts)
	dOpts.GeoM.Reset()
}

func (g *FishScene) loadBackground() error {

	background, err := loader.LoadImageAssetAsEbitenImage("roomBackground")
	if err != nil {
		return err
	}

	fishTankImg, err := loader.LoadImageAssetAsEbitenImage("fishTank")
	if err != nil {
		return err
	}

	frontLayer, err := loader.LoadImageAssetAsEbitenImage("frontLayer")
	if err != nil {
		return err
	}

	fishTankFrontLayer, err := loader.LoadImageAssetAsEbitenImage("fishTankFrontLayer")
	if err != nil {
		return err
	}

	fishTankNormal, err := loader.LoadImageAssetAsEbitenImage("fishTank_n")

	g.ImageMap["fishTankFrontLayer"] = fishTankFrontLayer
	g.ImageMap["fishTank"] = fishTankImg
	g.ImageMap["frontLayer"] = frontLayer
	g.ImageMap["background"] = background
	g.ImageMap["fishTankNormal"] = fishTankNormal
	g.ImageMap["offScreen"] = ebiten.NewImage(ScreenWidth, ScreenHeight)

	g.background = background
	return nil
}

func (g *FishScene) saveUISpritePositions() {

	spMap := make(map[string]drawables.SavePositionData)

	/*	for _, sprite := range g.sprites {

		//spData := sprite.SavePosition()
		spMap[spData.Name] = spData

	}*/

	outputSave, err := json.Marshal(spMap)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("assets/data/spritePosition.json", outputSave, 999)
	if err != nil {
		log.Fatal(err)
	}
}
