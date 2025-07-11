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

type FishScene struct {
	backGroundParams  map[string]any
	propQueue         *props.PropQueue
	loaded            bool
	tankSize          image.Rectangle
	sprites           [4][]drawables.Drawable
	ui                *ebitenui.UI
	gameLog           *sceneManagement.GameLog
	timers            map[string]*entities.Timer
	graphicManagerMap map[string]*graphics.GraphicManager
	shaderParams      map[string]any
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
}

var backGroundImgShelfHeight = 248

func NewFishScene(gameLog *sceneManagement.GameLog) *FishScene {

	roomImages, err := loader.LoadAllRoomBackGroundImages("assets/images/roomImages")
	if err != nil {
		log.Fatal("error while loading fish tank room assests:", err)
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	//hardcoded ass parameter cus i aint about to read sum dark brown particles dawg

	println("initiating game in ebiten NewFishScene()")

	g := &FishScene{}
	g.images = roomImages
	g.images.OffScreen = ebiten.NewImage(ScreenWidth, ScreenHeight)
	//stuff that needs to exist before game publishes shit

	//g.gameMode = Position
	g.gameLog = gameLog
	g.LoadGraphicManagerMap()

	g.shaderParams = make(map[string]any)
	g.backGroundParams = make(map[string]any)

	g.LoadTimers()

	collisionMap, err := geometry.LoadCollisions()
	if err != nil {
		log.Fatal(err)
	}

	g.collisionMap = collisionMap

	g.subs(collisionMap)

	g.debug = &debug.DebugData{DebugRect: &geometry.Rect{}, DebugParameter: make(map[debug.DebugOption]bool)}
	//g.debugRect.RectState = geometry.Off

	fishTankSizeX1 := float64(g.images.FishTank.Bounds().Max.X)
	fishTankSizeY1 := float64(g.images.FishTank.Bounds().Max.Y)

	g.shaderParams["ImgRect"] = [4]float64{0, 0, fishTankSizeX1, fishTankSizeY1}
	g.shaderParams["LightPoint"] = [2]float64{440, 250}

	tankX := g.images.FishTank.Bounds().Max.X
	tankY := g.images.FishTank.Bounds().Max.Y

	startingX := int(ScreenWidth * 0.2)
	startingY := ScreenHeight - backGroundImgShelfHeight - g.images.FishTank.Bounds().Dy()

	tankRect := image.Rect(startingX, startingY, tankX+startingX, tankY+startingY)
	g.tankSize = tankRect
	store := system.NewStore(g.gameLog.GlobalEventHub)
	g.store = &store

	g.environment = &system.Environment{}
	g.environment.Subscribe(g.gameLog.GlobalEventHub)

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
		interactableUIObjects.Door,
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
	log.Println("----FishScene OnEnter() called----")

	g.mouseFlags = &input.MouseFlags{HandledClick: false, CursorOccupied: false}
	g.backGroundParams["Cursor"] = []float64{440, 200}

	//No music on the base level as of now
	//g.timers["songTimer"].TurnOn()
	///JUST FOR TESTING, NO CHORE=> ALLOWANCE FRAMEWORK YET
	g.gameLog.GlobalEventHub.Publish(events.MoneyAvailable{Amount: 1})
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
		}
		g.gameLog.GlobalEventHub.Publish(ev)
		g.state.lastDayEntered = g.gameLog.Day
		g.gameLog.Tasks[0].Activate(g.gameLog.GlobalEventHub)
	}
	log.Println("----FishScene OnEnter() finished----")
}

func (g *FishScene) OnExit() {
	log.Println("----FishScene Exit----")
}

func (g *FishScene) LoadTimers() {

	g.timers = make(map[string]*entities.Timer)
	g.timers["pointGeneratedTimer"] = entities.NewTimer(0.2)
	g.timers["pointGeneratedTimer"].TurnOn()
	g.timers["songTimer"] = entities.NewTimer(15)
	g.timers["sceneTransition"] = entities.NewTimer(2.5)
	g.timers["publishNewTask"] = entities.NewTimer(0.2)

}

func (g *FishScene) IsLoaded() bool {
	return g.loaded
}

func (g *FishScene) Update() (sceneManagement.SceneId, error) {

	g.mouseFlags.HandledClick = false

	props.UpdateProps(g.propQueue)

	g.tutorialManager.Update()

	g.gameLog.SoundPlayer.Update()

	for _, sprite := range g.sprites[0] {
		sprite.Update()
	}

	for i := len(g.sprites[0]) - 1; i >= 0; i-- {
		if g.sprites[0][i].ShouldRemove() {
			g.sprites[0] = append(g.sprites[0][:i], g.sprites[0][i+1:]...)
		}
	}

	// Fix: iterate backwards when moving from sprites[0] to sprites[1]
	for i := len(g.sprites[0]) - 1; i >= 0; i-- {
		if g.sprites[0][i].Highlighted() {
			g.sprites[1] = append(g.sprites[1], g.sprites[0][i])
			g.sprites[0] = append(g.sprites[0][:i], g.sprites[0][i+1:]...)
		}
	}

	// Fix: iterate backwards when moving from sprites[1] to sprites[0]
	for i := len(g.sprites[1]) - 1; i >= 0; i-- {
		if !g.sprites[1][i].Highlighted() {
			g.sprites[0] = append(g.sprites[0], g.sprites[1][i])
			g.sprites[1] = append(g.sprites[1][:i], g.sprites[1][i+1:]...)
		}
	}

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
		}
		err := g.debug.DebugRect.Update()
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
	shaderOpts.Images[0] = g.images.FishTank
	shaderOpts.Images[1] = g.images.FishTank_n

	g.images.OffScreen.DrawImage(g.images.RoomBackground, &opts)
	b := g.images.FishTank.Bounds()

	opts.GeoM.Reset()

	shaderOpts.GeoM.Translate(float64(g.tankSize.Min.X), float64(g.tankSize.Min.Y))

	g.images.OffScreen.DrawRectShader(b.Dx(), b.Dy(), registry.ShaderMap["NormalMap"], shaderOpts)

	props.DrawProps(g.propQueue, g.images.OffScreen)

	for _, s := range g.sprites[0] {
		s.Draw(g.images.OffScreen)
	}

	opts.GeoM.Reset()
	fishTankFrontLayerDy := g.images.FishTankFrontLayer.Bounds().Dy()
	fishTankHeightDy := g.images.FishTank.Bounds().Dy()

	y := fishTankFrontLayerDy - fishTankHeightDy
	y = g.tankSize.Min.Y - y

	opts.GeoM.Translate(float64(g.tankSize.Min.X), float64(y))

	g.images.OffScreen.DrawImage(g.images.FishTankFrontLayer, &opts)

	opts.GeoM.Reset()

	g.images.OffScreen.DrawImage(g.images.FrontLayer, &opts)

}

func (g *FishScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
	//draw objects receiving global lighting shader
	g.DrawOffScreen()

	//apply shader to offscreen
	ShaderOpts := &ebiten.DrawRectShaderOptions{}
	ShaderOpts.Images[0] = g.images.OffScreen
	ShaderOpts.Uniforms = g.shaderParams

	screen.DrawRectShader(ScreenWidth, ScreenHeight, registry.ShaderMap["OnePointLighting"], ShaderOpts)

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
		g.debug.DebugRect.Init("tank")
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
	if g.currentTask > 0 || g.gameLog.Day > 1 {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			xCheck := x > g.tankSize.Min.X && x < g.tankSize.Max.X
			yCheck := y > g.tankSize.Min.Y && y < g.tankSize.Max.Y

			if xCheck && yCheck {
				filterFunc := func(distance any) bool {
					return distance.(float64) < 50
				}
				cursorX, cursorY := ebiten.CursorPosition()
				closestCreature := util.ClosestDrawableToCursor(cursorX, cursorY, g.sprites[0], filterFunc, "*entities.Creature")

				if closestCreature != nil {
					cre, ok := closestCreature.(*entities.Creature)
					if ok {
						SelectCreature(cre)
					}
				}
			}
		}
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

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		g.debug.DebugParameter[debug.Print] = true
	}

	if ebiten.IsKeyPressed(ebiten.Key1) {
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

func SelectCreature(creature *entities.Creature) {
	creature.Selected = true
	creature.Shader = registry.ShaderMap["Outline"]
	loader.LoadRotatingHighlightOutlineAnimated(creature.AnimatedSprite)
}

func (g *FishScene) debugModeParameterPrinterUpdater() {
	g.debug.DebugText = "Debug Mode Activated| Parameters:"
	for key, dbp := range g.debug.DebugParameter {
		if dbp {
			switch key {
			case debug.Position:
				g.debug.DebugText += "Position"
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
		case "Chores":
			g.returnScene = sceneManagement.MowingMiniGameScene
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(graphics.DrawGraphic{}, func(e tasks.Event) {
		ev := e.(graphics.DrawGraphic)
		println("adding draw graphic to game struct")
		g.sprites[1] = append(g.sprites[0], ev.Graphic)
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
			g.currentTask++ // this is zero indexed but slot is not so the current index is the just finished slot // FIX
			g.timers["publishNewTask"].TurnOn()
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
		g.gameLog.SoundPlayer.AddToQueue(soundFX.WhiteBoardMarker2, 2)
	})

}

func (g *FishScene) creatureSubs(colMap map[string]geometry.Rect) {
	g.gameLog.GlobalEventHub.Subscribe(input.MouseButtonPressedUISpriteActivity{}, func(e tasks.Event) {
		ev := e.(input.MouseButtonPressedUISpriteActivity)
		if g.timers["pointGeneratedTimer"].TimerState == entities.Done && !g.mouseFlags.HandledClick {
			g.mouseFlags.HandledClick = true
			pt := ev.Point.Clone()
			pt.X = pt.X - 50 + rand.Float32()*10
			pt.Y += 50
			p := entities.NewParticle(pt, colMap["tank"], g.gameLog.GlobalEventHub)
			g.sprites[0] = append(g.sprites[0], p)
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(events.NewPurchase{}, func(e tasks.Event) {
		ev := e.(events.NewPurchase)
		log.Printf("New Purchase:%s ", ev.Purchase)
		creature := LoadPurchasedSprite(g.environment, ev.Purchase, g.gameLog.GlobalEventHub, g.collisionMap["tank"])
		g.sprites[0] = append(g.sprites[0], creature)
	})

}

func (g *FishScene) CheckIfAllFishFed() bool {
	fed := true

	for _, draw := range g.sprites[0] {
		creature, ok := draw.(*entities.Creature)
		if ok && creature.Hunger != 0 {
			fed = false
		}
	}
	return fed

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
