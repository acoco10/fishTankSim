package scenes

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/fishTankWebGame/game"
	"github.com/acoco10/fishTankWebGame/game/debug"
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
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"log"
	"os"
)

type lightingState uint8

const (
	NightLight lightingState = iota
	Night
	Day
)

type FishScene struct {
	backGroundParams  map[string]any
	propQueue         props.PropQueue
	loaded            bool
	tankSize          image.Rectangle
	sprites           [4][]drawables.Drawable
	ui                *ebitenui.UI
	gameLog           *sceneManagement.GameLog
	timers            map[string]*entities.Timer
	graphicManagerMap map[string]*graphics.GraphicManager
	returnScene       sceneManagement.SceneId
	tutorialManager   *tutorial.Manager
	collisionMap      map[string]geometry.Rect
	store             *system.Store
	environment       *system.Environment
	playerState       *entities.Player
	images            *loader.BackGroundImages
	mouseFlags        *input.MouseFlags
	debug             *debug.DebugData
	state             *FishSceneState
	currentTask       int
	lightingShader    *ebiten.Shader
	offScreenShader   *ebiten.Shader
	tankShader        *ebiten.Shader
	tankShaderUpdater func(map[string]any) map[string]any
	shaderParams      map[string]any
	tankShaderParams  map[string]any
	lightingState     lightingState
	shaderUpdater     func(map[string]any) map[string]any
	testProp          *props.StructureProp //isolating for debug to be removed
	wallShader        *props.StructureProp
	smallerResolution *ebiten.Image
	resolutionScaling int
}

var backGroundImgShelfHeight = 248

func NewFishScene(gameLog *sceneManagement.GameLog) *FishScene {

	roomImages, err := loader.LoadAllRoomBackGroundImages("assets/images/roomImages")
	if err != nil {
		log.Fatal("error while loading fish tank room assests:", err)
	}

	println("initiating game in ebiten NewFishScene()")

	//render layers
	g := &FishScene{}
	g.images = roomImages
	g.images.OffScreen = ebiten.NewImage(ScreenWidth, ScreenHeight)
	g.images.OffScreen2 = ebiten.NewImage(ScreenWidth, ScreenHeight)
	g.smallerResolution = ebiten.NewImage(ScreenWidth, ScreenHeight)

	//stuff that needs to exist before game publishes shit

	//g.gameMode = Position
	g.gameLog = gameLog
	g.LoadGraphicManagerMap()

	g.tankShaderParams = make(map[string]any)
	g.shaderParams = make(map[string]any)
	g.backGroundParams = make(map[string]any)

	g.lightingShader = registry.ShaderMap["OnePointLighting"]
	g.offScreenShader = registry.ShaderMap["NormalMap"]
	g.tankShader = registry.ShaderMap["Water"]

	g.LoadTimers()

	collisionMap, err := geometry.LoadCollisions()
	if err != nil {
		log.Fatal(err)
	}

	g.collisionMap = collisionMap

	g.subs(collisionMap)

	g.debug = &debug.DebugData{DebugRect: &geometry.Rect{}, DebugParameter: make(map[debug.DebugOption]bool)}
	//g.debugRect.RectState = geometry.Off

	g.shaderParams["LightPoint"] = [2]float64{440, 250}

	tankX := g.images.FishTank.Bounds().Dx()
	tankY := g.images.FishTank.Bounds().Dy()

	startingX := int(ScreenWidth * 0.2)
	startingY := ScreenHeight - backGroundImgShelfHeight - g.images.FishTank.Bounds().Dy()

	tankRect := image.Rect(startingX, startingY, tankX+startingX, tankY+startingY)
	g.tankSize = tankRect
	store := system.NewStore(g.gameLog.GlobalEventHub)
	g.store = &store

	g.tankShaderParams["TankRect"] = [4]float64{float64(collisionMap["WaterEffect"].X1), float64(collisionMap["WaterEffect"].Y1), float64(collisionMap["WaterEffect"].X2), float64(collisionMap["WaterEffect"].Y2)}
	g.tankShaderParams["Counter"] = 0

	g.tankShaderUpdater = shaders.UpdateCounter
	g.environment = &system.Environment{}
	system.InitEnvironment(g.environment)
	g.environment.Subscribe(g.gameLog.GlobalEventHub)

	mainUI, _, err := ui.LoadMainFishMenu(ScreenWidth, ScreenHeight, gameLog.GlobalEventHub)
	if err != nil {
		log.Fatal("error loading scene")
	}

	g.ui = mainUI

	fishSceneUISprites := []interactableUIObjects.Label{
		interactableUIObjects.FishBook,
		interactableUIObjects.Records,
		interactableUIObjects.FishFood,
		interactableUIObjects.WhiteBoard,
		interactableUIObjects.PiggyBank,
		interactableUIObjects.Pillow,
		interactableUIObjects.Thermometer,
		interactableUIObjects.Magazine,
		interactableUIObjects.Door,
		interactableUIObjects.Phreader,
	}

	g.sprites = [4][]drawables.Drawable{}

	uiSprites, err := loader.LoadUISprites(fishSceneUISprites, g.environment, gameLog.GlobalEventHub, ScreenWidth, ScreenHeight)
	g.sprites[0] = append(g.sprites[0], uiSprites...)

	g.state = &FishSceneState{}

	return g
}

func LoadPurchasedSprite(environment *system.Environment, inputName string, hub *tasks.EventHub, tankSize geometry.Rect) *entities.Creature {
	fData := entities.SavedFish{
		FishType: inputName,
		Size:     1,
	}
	creature := loader.NewFish(environment, hub, tankSize, fData)
	return creature
}

func (g *FishScene) LoadStuff() {
	propq := loader.LoadProps(g.gameLog.Save.TankObjects, g.tankSize, g.gameLog.GlobalEventHub)
	g.propQueue = propq

	g.playerState = &entities.Player{Money: 10, EventHub: g.gameLog.GlobalEventHub}
	g.playerState.Subscribe()

	loaderMan := loader.Manager{Hub: g.gameLog.GlobalEventHub}
	loaderMan.Subscriptions()
}

func (g *FishScene) FirstLoad() {
	log.Println("----fishScene.firstLoad() entered----")
	if g.gameLog.Day == 1 {
		println("length of game log save =", len(g.gameLog.Save.Fish))
		fishes := g.gameLog.Save

		for _, fish := range fishes.Fish {
			loadedFish := loader.NewFish(g.environment, g.gameLog.GlobalEventHub, g.collisionMap["tank"], fish)
			g.sprites[0] = append(g.sprites[0], loadedFish)
		}
	}
	g.LoadStuff()
	g.loaded = true

	ev2 := events.FishTankLayout{
		Rectangle: g.tankSize,
	}

	g.gameLog.GlobalEventHub.Publish(ev2)
	log.Println("----fishScene.firstLoad() finished----")

}

func (g *FishScene) OnEnter() {
	g.resolutionScaling = ScaleScreenToResolution(g.smallerResolution)
	log.Println("----FishScene OnEnter() called----")

	g.mouseFlags = &input.MouseFlags{HandledClick: false, CursorOccupied: false}
	g.backGroundParams["Cursor"] = []float64{440, 200}

	//No music on the base level as of now
	//g.timers["songTimer"].TurnOn()
	///JUST FOR TESTING, NO CHORE=> ALLOWANCE FRAMEWORK YET
	g.returnScene = sceneManagement.FishTank

	tutMngr := tutorial.Manager{}
	tutorial.InitData(&tutMngr, g.gameLog.GlobalEventHub)
	g.tutorialManager = &tutMngr

	if g.gameLog.Day != g.state.lastDayEntered {
		ev := events.NewDay{NTasks: len(g.gameLog.Tasks)}
		switch g.gameLog.DayType {
		case sceneManagement.Chores:
			ev.Type = "Chores"
		case sceneManagement.Free:
			ev.Type = "Free"
		case sceneManagement.Camp:
			ev.Type = "Camp"
		}
		g.gameLog.GlobalEventHub.Publish(ev)
		g.state.lastDayEntered = g.gameLog.Day
		g.gameLog.Tasks[0].Activate(g.gameLog.GlobalEventHub)
		g.lightingShader = registry.ShaderMap["DayLight"]
		g.backGroundParams["Cursor"] = [2]float64{0, 0}
		g.shaderUpdater = nil
		g.lightingState = Day
		return
	}
	g.SetNightLight()
	g.gameLog.SongPlayer.Play(soundFX.TropicalHouse)
	log.Println("----FishScene OnEnter() finished----")
}

func (g *FishScene) SetNightLight() {
	g.lightingState = NightLight
	g.lightingShader = registry.ShaderMap["OnePointLighting"]
	g.shaderParams["LightPoint"] = [2]float64{440, 250}
}

func (g *FishScene) OnExit() {
	log.Println("----FishScene Exit----")
}

func (g *FishScene) LoadTimers() {
	g.timers = make(map[string]*entities.Timer)
	g.timers["pointGeneratedTimer"] = entities.NewTimer(0.2)
	g.timers["pointGeneratedTimer"].TurnOn()
	g.timers["songTimer"] = entities.NewTimer(15)
	g.timers["sceneTransition"] = entities.NewTimer(5.0)
	g.timers["publishNewTask"] = entities.NewTimer(0.2)
	g.timers["leaveScene"] = entities.NewTimer(10.0)
	g.timers["leaveScene"].TurnOn()
}

func (g *FishScene) IsLoaded() bool {
	return g.loaded
}

func (g *FishScene) Update() (sceneManagement.SceneId, error) {
	if g.shaderUpdater != nil {
		g.shaderParams = g.shaderUpdater(g.shaderParams)
	}
	if g.tankShaderUpdater != nil {
		g.tankShaderParams = g.tankShaderUpdater(g.tankShaderParams)
	}

	g.mouseFlags.HandledClick = false

	props.UpdateProps(g.propQueue)

	g.tutorialManager.Update()

	g.gameLog.SoundPlayer.Update()

	for _, sprite := range g.sprites[0] {
		sprite.Update()
	}

	g.ManageLayers()
	//update second layer after managing layers?
	for _, s := range g.sprites[1] {
		s.Update()
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

	//g.tutorialManager.CharUpdate()

	if g.debug.GameMode == debug.Debug {
		if g.debug.DebugParameter[debug.Position] {
			g.positionModeUpdate()
			ev := events.ButtonClickedEvent{ButtonText: "Mode"}
			g.gameLog.GlobalEventHub.Publish(ev)
		}
		err := g.debug.DebugRect.Update()
		if err != nil {
			//debug rect could error when saving collision location
			return g.returnScene, err
		}

		if g.debug.DebugParameter[debug.ShaderTest] {
			ShaderSwapper(g)
		}
	}
	if g.testProp != nil {
		g.testProp.Update()
	}
	return g.returnScene, nil
}

func (g *FishScene) ManageLayers() {
	for i := len(g.sprites[0]) - 1; i >= 0; i-- {
		if g.sprites[0][i].ShouldRemove() {
			g.sprites[0] = append(g.sprites[0][:i], g.sprites[0][i+1:]...)
		}
	}

	// iterate backwards when moving from sprites[0] to sprites[1]
	for i := len(g.sprites[0]) - 1; i >= 0; i-- {
		if g.sprites[0][i].Highlighted() {
			g.sprites[1] = append(g.sprites[1], g.sprites[0][i])
			g.sprites[0] = append(g.sprites[0][:i], g.sprites[0][i+1:]...)
		}
	}

	// iterate backwards when moving from sprites[1] to sprites[0]
	for i := len(g.sprites[1]) - 1; i >= 0; i-- {
		if !g.sprites[1][i].Highlighted() {
			g.sprites[0] = append(g.sprites[0], g.sprites[1][i])
			g.sprites[1] = append(g.sprites[1][:i], g.sprites[1][i+1:]...)
		}
	}
}

func (g *FishScene) DrawOffScreen() {

	opts := &ebiten.DrawImageOptions{}
	shaderOpts := &ebiten.DrawRectShaderOptions{}

	shaderOpts.Uniforms = g.backGroundParams
	shaderOpts.Images[0] = g.images.RoomBackground

	if g.wallShader != nil {
		g.images.OffScreen.DrawRectShader(ScreenWidth, ScreenHeight, registry.ShaderMap["Wall"], shaderOpts)
	}

	g.images.OffScreen.DrawImage(g.images.RoomBackground, opts)

	b := g.images.FishTank.Bounds()
	opts.GeoM.Reset()

	shaderOpts.Uniforms = g.backGroundParams
	shaderOpts.Images[0] = g.images.FishTank
	shaderOpts.Images[1] = g.images.FishTank_n

	if g.offScreenShader != nil {
		shaderOpts.GeoM.Translate(float64(g.tankSize.Min.X), float64(g.tankSize.Min.Y))
		g.images.OffScreen.DrawRectShader(b.Dx(), b.Dy(), g.offScreenShader, shaderOpts)
	} else {
		opts.GeoM.Translate(float64(g.tankSize.Min.X), float64(g.tankSize.Min.Y))
		g.images.OffScreen.DrawImage(g.images.FishTank, opts)
	}

	if g.testProp != nil {
		g.testProp.Draw(g.images.OffScreen)
	}
	//props.DrawProps(g.propQueue, g.images.OffScreen)

	for _, s := range g.sprites[0] {
		s.Draw(g.images.OffScreen)
	}

	opts.GeoM.Reset()
	fishTankFrontLayerDy := g.images.FishTankFrontLayer.Bounds().Dy()
	fishTankHeightDy := g.images.FishTank.Bounds().Dy()

	y := fishTankFrontLayerDy - fishTankHeightDy
	y = g.tankSize.Min.Y - y

	opts.GeoM.Translate(float64(g.tankSize.Min.X), float64(y))
	switch g.lightingState {
	case NightLight:
		g.images.OffScreen.DrawImage(g.images.FishTankFrontLayerNoLight, opts)
	case Night:
		g.images.OffScreen.DrawImage(g.images.FishTankFrontLayerNoLight, opts)
	case Day:
		g.images.OffScreen.DrawImage(g.images.FishTankFrontLayerDayLight, opts)
	}

	opts.GeoM.Reset()

	if g.tankShader != nil {
		shaderOpts.GeoM.Reset()
		shaderOpts.Uniforms = g.tankShaderParams
		shaderOpts.Images[0] = g.images.OffScreen
		shaderOpts.Images[1] = nil
		g.images.OffScreen2.DrawRectShader(ScreenWidth, ScreenHeight, g.tankShader, shaderOpts)
	} else {
		g.images.OffScreen2.DrawImage(g.images.OffScreen, opts)
	}

	g.images.OffScreen2.DrawImage(g.images.FrontLayer, opts)

}

func (g *FishScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
	//draw objects receiving global lighting shader
	g.DrawOffScreen()

	//apply shader to offscreen
	ShaderOpts := &ebiten.DrawRectShaderOptions{}
	ShaderOpts.Images[0] = g.images.OffScreen2
	ShaderOpts.Uniforms = g.shaderParams

	g.smallerResolution.DrawRectShader(ScreenWidth, ScreenHeight, g.lightingShader, ShaderOpts)

	dopts := &ebiten.DrawImageOptions{}
	dopts.GeoM.Scale(float64(g.resolutionScaling), float64(g.resolutionScaling))

	screen.DrawImage(g.smallerResolution, dopts)

	g.debug.DebugRect.Draw(screen)

	//vector.StrokeCircle(screen, 440, 170, 2, 1, colornames.Yellow, false)

	/*for _, s := range g.sprites {
		s.Draw(screen)
	}*/

	for _, sp := range g.sprites[1] {
		sp.Draw(screen)
	}

	for _, graphicMan := range g.graphicManagerMap {
		graphicMan.Draw(screen)
	}

	DebugText(g.debug.DebugText, screen)

	graphics.DrawGraphics(screen)
	g.ui.Draw(screen)

}

func (g *FishScene) positionModeUpdate() {
	if ebiten.IsKeyPressed(ebiten.KeyM) {
		g.debug.DebugRect.Init("")
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.saveUISpritePositions()
	}

	err := g.debug.DebugRect.Update()

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
			g.lightingState = Night
			g.gameLog.GlobalEventHub.Publish(events.DayOverTransitionComplete{})
		}

		if key == "publishNewTask" && state == entities.Done {
			if len(g.gameLog.Tasks) > g.currentTask {
				timer.TurnOff()
				g.gameLog.Tasks[g.currentTask].Activate(g.gameLog.GlobalEventHub)
			} else {
				timer.TurnOff()
				g.currentTask = 0
			}
		}
		if key == "leaveScene" && state == entities.Done {
			timer.TurnOff()
			if g.gameLog.DayType == sceneManagement.Camp {
				eff := loader.LoadStaticEffect("timeForCamp", 100, 85)
				g.sprites[1] = append(g.sprites[1], eff)
			}

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

func (g *FishScene) debugInputCheck() {

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		switch g.debug.GameMode {
		case debug.Debug:
			g.debug.GameMode = debug.User
			g.debug.DebugText = ""
			g.debug.DebugParameter[debug.Print] = false
		case debug.User:
			g.debug.GameMode = debug.Debug
			ev := events.ButtonClickedEvent{ButtonText: "Mode"}
			g.gameLog.GlobalEventHub.Publish(ev)
			g.debugModeParameterPrinterUpdater()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.debug.DebugParameter[debug.Print] = true
	}

	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		switch g.debug.DebugParameter[debug.ShaderTest] {
		case true:
			g.debug.DebugParameter[debug.ShaderTest] = false
		case false:
			g.debug.DebugParameter[debug.ShaderTest] = true
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.debug.DebugParameter[debug.Position] = true
		g.debugModeParameterPrinterUpdater()
	}

	for key, dbp := range g.debug.DebugParameter {
		switch key {
		case debug.Position:
			if dbp && ebiten.IsKeyPressed(ebiten.KeyS) {
				g.saveUISpritePositions()
			}
		}
	}

}

func (g *FishScene) debugModeParameterPrinterUpdater() {
	g.debug.DebugText = "Debug Mode Activated| Parameters:"
	for key, dbp := range g.debug.DebugParameter {
		if dbp {
			switch key {
			case debug.Position:
				g.debug.DebugText += "Position"
			case debug.ShaderTest:
				g.debug.DebugText += "Lighting"
			}
		}
	}
}

func (g *FishScene) LoadGraphicManagerMap() {

	g.graphicManagerMap = make(map[string]*graphics.GraphicManager)
	WhiteBoardGraphicManager := graphics.NewGraphicManager(g.gameLog.GlobalEventHub, graphicManagerSubscriptions.WhiteBoardGMSubs)

	g.graphicManagerMap["whiteBoard"] = WhiteBoardGraphicManager

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

	var propPicked string

	g.gameLog.GlobalEventHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		println(ev.ButtonText, "button event received")
		switch ev.ButtonText {
		case "Save":
			g.SaveGame()
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(graphics.DrawGraphic{}, func(e tasks.Event) {
		ev := e.(graphics.DrawGraphic)
		println("adding draw graphic to game struct")
		g.sprites[1] = append(g.sprites[0], ev.Graphic)
	})

	g.gameLog.GlobalEventHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		if ev.ButtonText == "Go to Bed?: Yes" {
			g.lightingShader = registry.ShaderMap["TurnOff"]
			g.shaderUpdater = shaders.FadeLightIntensityForTurnOff
			g.shaderParams["LightIntensity"] = 1.0
			g.shaderParams["Counter"] = 0.0
			g.timers["sceneTransition"].TurnOn()
			ev2 := events.DayOver{}
			//change state to let the game know to draw unlit art
			g.lightingState = Night
			//turn of normal maps:
			g.offScreenShader = nil
			g.gameLog.GlobalEventHub.Publish(ev2)
			g.gameLog.Day++
			return
		}
		if ev.ButtonText == "Go do your Chores?: Yes" {
			g.returnScene = sceneManagement.MowingMiniGameScene
			return
		}
		if ev.ButtonText == "Go to Camp?: Yes" {
			g.returnScene = sceneManagement.CampScene
			return
		}
		if ev.ButtonText == "Castle" || ev.ButtonText == "Log" {
			propPicked = ev.ButtonText
			return
		}
		if ev.ButtonText == "Confirm for prop select" {
			if propPicked != "" {
				ev2 := events.CloseWindow{OverRide: true}
				g.gameLog.GlobalEventHub.Publish(ev2)
				g.environment.AddTankModifier(propPicked)
				p := loader.LoadProp(propPicked, g.tankSize, g.gameLog.GlobalEventHub)
				g.testProp = p
			}
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
			return
		}
		if ev.UiSprite == "fishFood" && ev.UiSpriteAction == "picked up" {
			g.gameLog.SoundPlayer.AddToQueue(soundFX.SelectSound2, 1)
			return
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(tasks.TaskCompleted{}, func(e tasks.Event) {
		g.gameLog.SoundPlayer.AddToQueue(soundFX.WhiteBoardMarker1, 2)
		g.gameLog.SoundPlayer.AddToQueue(soundFX.SuccessMusic, 1)
		ev := e.(tasks.TaskCompleted)
		if len(g.gameLog.Tasks) > ev.Slot {
			g.currentTask++ // this is zero indexed but slot is not so the current index is the just finished slot // FIX
			g.timers["publishNewTask"].TurnOn()
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
		g.gameLog.SoundPlayer.AddToQueue(soundFX.WhiteBoardMarker2, 2)
	})

}

func (g *FishScene) SaveGame() {
	println("save game event generated and received")
	var savedFish []entities.SavedFish

	for _, draw := range g.sprites[0] {
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
	DebugText(g.debug.DebugText, screen)
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

func (g *FishScene) saveUISpritePositions() {

	spMap := make(map[string]drawables.SavePositionData)
	//THIS IS AWFUL BECAUSE I MADE TOO MANY UI SPRITE TYPES DEAL WITH SOME DAY
	//check for ui sprites
	for _, layer := range g.sprites {
		for _, sprite := range layer {
			uiSprite, ok := sprite.(*interactableUIObjects.UiSprite)
			if !ok {
				continue
			}
			println("saving uiSprite", uiSprite.Label)
			spData := uiSprite.SavePosition()
			spMap[spData.Name] = spData
		}

	}
	//unforunately we made like 8 different types of these things so we re read to maintain their current positiona atleast
	currentPos, err := loader.LoadSpritePositionData()
	if err != nil {
		log.Fatal(err)
	}

	for _, sprite := range currentPos {
		if _, ok := spMap[sprite.Name]; !ok {
			spMap[sprite.Name] = *sprite
		}
	}

	outputSave, err := json.Marshal(spMap)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("assets/data/spritePosition.json", outputSave, 999)
	if err != nil {
		log.Fatal(err)
	}
}
