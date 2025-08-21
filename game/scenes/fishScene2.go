package scenes

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/daySystem"
	"github.com/acoco10/fishTankWebGame/game/debug"
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/input"
	"github.com/acoco10/fishTankWebGame/game/loader"
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
	"golang.org/x/image/colornames"
	"image"
	"log"
	"math/rand"
	"os"
	"strings"
	time2 "time"
)

type lightingState uint8

const (
	NightLight lightingState = iota
	Night
	Day
)

var backGroundImgShelfHeight = 150

type FishScene2 struct {
	zoomScreen      *ebiten.Image
	state           FishSceneState
	loaded          bool
	tankSize        image.Rectangle
	ui              *ebitenui.UI
	gameLog         *sceneManagement.GameLog
	timers          map[string]*util.Timer
	returnScene     sceneManagement.SceneId
	tutorialManager *tutorial.Manager
	//store                 *system.Store
	collisionMap            map[string]image.Rectangle
	environment             *system.Environment
	playerState             *entities.Player
	images                  *loader.BackGroundImages
	debug                   *debug.DebugData
	currentTask             int
	activatedCollisions     []Collision
	testProp                *entities.StructureProp
	gameState               *entities.GameState
	lightingState           lightingState
	globalLightingShaderMap map[lightingState]*ebiten.Shader
	shaderParamsMap         map[lightingState]map[string]any
	smallerResolution       *ebiten.Image
	whiteBoardSprite        *entities.WhiteBoardSprite
	taskRelatedEventQueue   []tasks.Event
	cachedEntity            uint32 //random place to stash an entity i want to refocus or keep track of
}

type Collision struct {
	image.Rectangle
	z int
}

var creatureManager CreatureManager

type CreatureManager struct {
	allFishFed bool
}

func NewFishScene2(gameLog *sceneManagement.GameLog) *FishScene2 {
	g := &FishScene2{}
	collisionMap, err := util.LoadCollisions()
	if err != nil {
		log.Fatal(err)
	}

	_, exists := collisionMap["tank"]

	if !exists {
		log.Fatal()
	}

	g.collisionMap = collisionMap

	roomImages, err := loader.LoadAllRoomBackGroundImages("assets/images/roomImagesSmaller")
	if err != nil {
		log.Fatal("error while loading fish tank room assets:", err)
	}

	println("initiating game in ebiten NewFishScene2()")

	//render layers

	g.images = roomImages
	g.images.OffScreen = ebiten.NewImage(ScreenWidth, ScreenHeight)
	g.images.OffScreen2 = ebiten.NewImage(ScreenWidth, ScreenHeight)
	g.smallerResolution = ebiten.NewImage(ScreenWidth, ScreenHeight)

	g.globalLightingShaderMap = make(map[lightingState]*ebiten.Shader)

	g.globalLightingShaderMap[Day] = registry.ShaderMap["DayLight"]
	g.globalLightingShaderMap[NightLight] = registry.ShaderMap["OnePointLighting"]

	tankRect := g.collisionMap["tank"]

	g.shaderParamsMap = make(map[lightingState]map[string]any)

	nightLightParams := make(map[string]any)
	nightLightParams["LightPoint"] = [2]float64{
		float64(g.collisionMap["tank"].Min.X+194) * registry.Config.ResolutionScalingF,
		float64(g.collisionMap["tank"].Min.Y-52) * registry.Config.ResolutionScalingF}
	nightLightParams["LightWidth"] = 120.0 * registry.Config.ResolutionScalingF
	nightLightParams["TankRect"] = [4]float64{
		float64(tankRect.Min.X) * registry.Config.ResolutionScalingF,
		float64(tankRect.Min.Y) * registry.Config.ResolutionScalingF,
		float64(tankRect.Max.X) * registry.Config.ResolutionScalingF,
		float64(tankRect.Max.Y) * registry.Config.ResolutionScalingF}

	g.shaderParamsMap[NightLight] = nightLightParams

	g.shaderParamsMap[Day] = make(map[string]any) //no params as of now

	tankX := g.images.FishTank.Bounds().Dx()
	tankY := g.images.FishTank.Bounds().Dy()

	startingX := int(ScreenWidth * 0.2)
	startingY := ScreenHeight - backGroundImgShelfHeight - g.images.FishTank.Bounds().Dy()

	tankRect = image.Rect(startingX, startingY, tankX+startingX, tankY+startingY)
	g.tankSize = tankRect
	g.collisionMap["tankRect"] = tankRect

	fishtankSpriteFlags := loader.SpriteEntFlags{Unfocusable: true}
	loader.MakeSpriteEntity(g.images.FishTank, float32(g.tankSize.Min.X), float32(g.tankSize.Min.Y), fishtankSpriteFlags)

	fishtankFrontSpriteFlags := loader.SpriteEntFlags{Unfocusable: true, Zlayer: 2}
	loader.MakeSpriteEntity(g.images.FishTankFrontLayerNoLightSmaller, float32(g.tankSize.Min.X), float32(g.tankSize.Min.Y), fishtankFrontSpriteFlags)
	fishTankFrontLayerDy := g.images.FishTankFrontLayerNoLightSmaller.Bounds().Dy()
	fishTankHeightDy := g.images.FishTank.Bounds().Dy()

	y := fishTankFrontLayerDy - fishTankHeightDy
	y = g.tankSize.Min.Y - y

	//stuff that needs to exist before game publishes shit

	//g.gameMode = Position
	g.gameLog = gameLog

	g.LoadTimers()

	g.subs(collisionMap)

	//g.debugRect.RectState = geometry.Off

	//nighlight parameters

	g.environment = &system.Environment{}
	system.InitEnvironment(g.environment)
	g.environment.Subscribe(g.gameLog.GlobalEventHub)

	mainUI, _, err := ui.LoadMainFishMenu(ScreenWidth, ScreenHeight, gameLog.GlobalEventHub)
	if err != nil {
		log.Fatal("error loading scene")
	}

	g.ui = mainUI

	return g
}

func (g *FishScene2) LoadStuff() {

	g.playerState = &entities.Player{Money: 0, EventHub: g.gameLog.GlobalEventHub}
	g.playerState.Subscribe()

}

func (g *FishScene2) FirstLoad() {
	log.Println("----fishScene2.firstLoad() entered----")
	if g.gameLog.Day == 1 {
		println("length of game log save =", len(g.gameLog.Save.Fish))

		//g.gameLog.Save.Fish = []entities.SavedFish{entities.SavedFish{FishType: "fish", Size: 2}, entities.SavedFish{FishType: "fish", Size: 3}}

	}
	g.zoomScreen = ebiten.NewImage(registry.Config.ScreenWidth, registry.Config.ScreenHeight)
	//store := system.NewStore(g.gameLog.GlobalEventHub)
	//g.store = store
	g.LoadStuff()

	tutMngr := &tutorial.Manager{}
	tutorial.InitData(tutMngr, g.gameLog.GlobalEventHub)
	g.tutorialManager = tutMngr

	ev2 := events.FishTankLayout{
		Rectangle: g.tankSize,
	}
	g.gameLog.GlobalEventHub.Publish(ev2)

	fishScene2UISprites := []entities.Label{
		entities.WhiteBoard,
		entities.PiggyBank,
		entities.Pillow,
		entities.Thermometer,
		entities.Magazine,
		entities.Door,
		entities.Phreader,
		entities.FishFood,
	}

	_, wbSprite := loader.LoadAllEntities(fishScene2UISprites, g.gameLog.Save.Fish, g.environment, g.gameLog.GlobalEventHub, g.collisionMap)
	g.whiteBoardSprite = wbSprite
	g.loaded = true

	daySystem.LoadDaysTasks(g.gameLog)
	log.Println("----fishScene2.firstLoad() finished----")

}

func (g *FishScene2) OnEnter() {
	//cu := customCursor.CreateCursorUpdater(g.gameLog.GlobalEventHub)
	//eInput.SetCursorUpdater(cu)
	g.debug = &debug.DebugData{DebugRect: &util.Rect{
		Rectangle: &image.Rectangle{Min: image.Point{}, Max: image.Point{}}},
		DebugParameter: make(map[debug.DebugOption]bool)}

	g.debug.DebugParameter[debug.Print] = false
	g.debug.DebugParameter[debug.Position] = false

	g.debug.DebugRect.Init("CastleCollisions", g.gameLog.GlobalEventHub)
	log.Println("----FishScene2 OnEnter() called----")

	mouseFlags := &input.MouseFlags{HandledClick: false, CursorOccupied: false}
	g.gameState = &entities.GameState{MouseFlags: mouseFlags}

	//No music on the base level as of now

	g.timers["songTimer"].TurnOn()

	g.returnScene = sceneManagement.FishTank

	println("made it to new day")
	if g.gameLog.Day != g.state.lastDayEntered {
		g.configureNewDay()
		entities.UpdateCursorForEntitiesWNormals([]float64{0, 0})
		log.Println("----FishScene2 OnEnter() finished----")
		return
	}

	g.configureExistingDay()

	log.Println("----FishScene2 OnEnter() finished----")
}

func (g *FishScene2) configureExistingDay() {
	g.lightingState = NightLight
}

func (g *FishScene2) configureNewDay() {
	ev := events.NewDay{NTasks: len(g.gameLog.Tasks), Day: g.gameLog.Day}
	g.lightingState = Day
	switch g.gameLog.DayType {
	case sceneManagement.Chores:
		ev.DayType = "Chores"
	case sceneManagement.Camp:
		ev.DayType = "Camp"
	default:
		ev.DayType = "Free"
	}

	g.gameLog.GlobalEventHub.Publish(ev)
	g.state.lastDayEntered = g.gameLog.Day
	if len(g.gameLog.Tasks) > 0 {
		g.gameLog.Tasks[0].Activate(g.gameLog.GlobalEventHub)
	}

	creatureManager.allFishFed = false
}

func (g *FishScene2) OnExit() {
	log.Println("----FishScene2 Exit----")
	g.gameLog.SongPlayer.Pause()
	g.gameLog.GlobalEventHub.Publish(events.LeavingFishScene{})
	graphics.DeInitAllGraphics()

}

func (g *FishScene2) LoadTimers() {
	g.timers = make(map[string]*util.Timer)
	g.timers["pointGeneratedTimer"] = util.NewTimer(0.2)
	g.timers["pointGeneratedTimer"].TurnOn()
	g.timers["songTimer"] = util.NewTimer(2)
	g.timers["sceneTransition"] = util.NewTimer(3.0)
	g.timers["publishNewTask"] = util.NewTimer(0.2)
	g.timers["leaveScene"] = util.NewTimer(10.0)
	g.timers["leaveScene"].TurnOn()
}

func (g *FishScene2) IsLoaded() bool {
	return g.loaded
}

func (g *FishScene2) Update() (sceneManagement.SceneId, error) {
	g.ui.Update()
	g.gameState.MouseFlags.HandledClick = false
	g.tutorialManager.Update()

	//g.updateInput()

	entities.UpdateEntities(g.gameState)
	g.whiteBoardSprite.Update()

	g.gameLog.SongPlayer.Update()
	g.gameLog.SoundPlayer.Update()

	g.updateTimers()

	if g.debug.GameMode == debug.Debug {
		err := g.debug.DebugRect.Update()
		if err != nil {
			//debug rect could error when saving collision location
			return g.returnScene, err
		}

		if g.debug.DebugParameter[debug.ShaderTest] {
			//ShaderSwapper(g)
		}
	}

	graphics.UpdateGraphics()

	if g.gameState.FocusedEntity != nil && !g.gameState.FocusedEntity.Sprite.Focused {
		g.gameState.FocusedEntity = nil
	}

	return g.returnScene, nil
}

func (g *FishScene2) DrawOffScreen() {
	//define our draw options
	opts := &ebiten.DrawImageOptions{}
	//shaderOpts := &ebiten.DrawRectShaderOptions{}

	//step1 draw room background to primary offscreen
	/*x, y := ebiten.WindowSize()
	if float64(x)/float64(y) != 16.0/9.0 {
		g.images.OffScreen.DrawImage(g.images.LaptopRoomBackground, opts)
	}*/
	g.images.OffScreen.DrawImage(g.images.RoomBackground, opts)

	opts.GeoM.Translate(float64(g.tankSize.Min.X), float64(g.tankSize.Min.Y))
	//g.images.OffScreen.DrawImage(g.images.FishTank, opts)

	//b := g.images.FishTank.Bounds()
	opts.GeoM.Reset()

	entities.DrawEntities(g.images.OffScreen)

	sopts := &ebiten.DrawRectShaderOptions{}

	if registry.Config.Zoom {
		sopts.GeoM.Scale(registry.Config.ZoomFactor, registry.Config.ZoomFactor)
		sopts.GeoM.Translate(registry.Config.ZoomOffSetX, registry.Config.ZoomOffSetY)
		sopts.Images[0] = g.images.OffScreen
		sopts.Uniforms = g.shaderParamsMap[g.lightingState]
		b := g.images.OffScreen.Bounds()
		g.zoomScreen.DrawRectShader(b.Dx(), b.Dy(), g.globalLightingShaderMap[g.lightingState], sopts)
		//g.zoomScreen.DrawImage(g.images.OffScreen, opts)
	}
}

func (g *FishScene2) Draw(screen *ebiten.Image) {

	//call offscreen function to render everything we want to apply global shaders to
	g.DrawOffScreen()

	dOpts := &ebiten.DrawImageOptions{}
	shaderOpts := &ebiten.DrawRectShaderOptions{}
	shaderOpts.Images[0] = g.images.OffScreen
	shaderOpts.Uniforms = g.shaderParamsMap[g.lightingState]
	b := g.images.OffScreen.Bounds()
	g.smallerResolution.DrawRectShader(b.Dx(), b.Dy(), g.globalLightingShaderMap[g.lightingState], shaderOpts)

	//debug rect needs to be scaled to base resolution
	//scale our draw opts to resolution

	if g.debug.GameMode == debug.Debug {

		g.debug.DebugRect.Draw(g.smallerResolution)
		/*for key, rect := range g.collisionMap {
			util.StrokeRectFromImageRect(rect, g.smallerResolution, colornames.Orangered, false)
			cs := ebiten.ColorScale{}
			cs.SetB(0.0)
			cs.SetG(0.0)
			cs.SetA(1.0)
			face := registry.FontMap["nk57"]
			tOpts := &text.DrawOptions{}
			tOpts.ColorScale = cs
			tOpts.GeoM.Translate(float64(rect.Min.X+rect.Dx())/registry.Config.ResolutionScalingF, float64(rect.Min.Y)*registry.Config.ResolutionScalingF)
			text.Draw(screen, key, face, tOpts)
		}*/

		for _, rect := range g.activatedCollisions {
			util.StrokeRectFromImageRect(rect.Rectangle, g.smallerResolution, colornames.Greenyellow, false)
		}
	}

	graphics.DrawScaledGraphics(g.smallerResolution)

	dOpts.GeoM.Reset()

	dOpts.GeoM.Translate(0, registry.Config.YOffsetF)
	dOpts.GeoM.Scale(registry.Config.ResolutionScalingF, registry.Config.ResolutionScalingF)

	if registry.Config.Zoom {
		screen.DrawImage(g.zoomScreen, dOpts)
		entities.DrawNonZoomedEntities(screen)
		graphics.DrawScaledGraphicsOnMainScreen(screen)
		graphics.DrawUnScaledGraphics(screen)
		g.ui.Draw(screen) //also drawing our custom cursor
		return
	} else {
		screen.DrawImage(g.smallerResolution, dOpts)
	}
	//custom debug text
	//tps monitor
	//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f\n", ebiten.ActualTPS()), 100, 100)
	//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Window Open:%t, Focused Entity = %v, HoveredUiSprite = %v", g.gameState.MouseFlags.WindowOpen, g.gameState.FocusedEntity, g.gameState.HoveredUiSprite), 100, 120)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Environment Temperatre: %f, Environment natural PH: %f, Environment Modified PH: %f", g.environment.Temperature, g.environment.NaturalPHLevel, g.environment.ModifiedPHLevel), 100, 140)

	//draw graphics

	graphics.DrawUnScaledGraphics(screen)
	g.whiteBoardSprite.DstImg.Draw(screen)
	graphics.DrawScaledGraphicsOnMainScreen(screen)

	//draw ui last
	g.ui.Draw(screen)

}

func (g *FishScene2) positionModeUpdate() {

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.saveUISpritePositions()
	}

	err := g.debug.DebugRect.Update()

	if err != nil {
		log.Printf("Couldn't save positions or something with the debug rect got fucked up")
	}
}

func (g *FishScene2) updateTimers() {

	for key, timer := range g.timers {
		state := timer.Update()

		if key == "songTimer" && state == util.Done {
			if g.lightingState == NightLight {
				g.gameLog.SongPlayer.FadeIn(soundFX.TropicalHouse)
			}
			if g.lightingState == Day {
				g.gameLog.SongPlayer.FadeIn(soundFX.DayTimeJazz)
			}
			timer.TurnOff()
		}

		if key == "sceneTransition" && state == util.Done {
			timer.TurnOff()
			g.returnScene = sceneManagement.TransitionScene
			//g.lightingState = Night
			g.gameLog.GlobalEventHub.Publish(events.DayOverTransitionComplete{})
		}

		if key == "publishNewTask" && state == util.Done {
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

func (g *FishScene2) updateInput() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		entities.FocusOnClickedEntity(g.gameState)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		if g.gameState.FocusedEntity != nil && !g.gameState.MouseFlags.WindowOpen {
			entities.UnFocus(g.gameState.FocusedEntity.Id)
			g.gameState.FocusedEntity = nil
		}
	}

	//g.debugInputCheck()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.gameLog.GlobalEventHub.Publish(events.CloseWindow{})
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if !registry.Config.Zoom {
			registry.Config.Set(registry.Zoom, true)
			if g.gameState.FocusedEntity != nil {
				g.cachedEntity = g.gameState.FocusedEntity.Id
				entities.UnFocus(g.gameState.FocusedEntity.Id)
			}

		} else {

			registry.Config.Set(registry.Zoom, false)

			if g.cachedEntity != 0 {
				entities.ReFocus(g.cachedEntity)

				ent, exists := entities.GetEntity(g.cachedEntity)
				if !exists {
					log.Println("tried to refocus an entity in zoom input updater that doesnt exist")
					return
				}
				g.gameState.FocusedEntity = ent
			}
			g.cachedEntity = 0
		}
	}

}

func (g *FishScene2) debugInputCheck() {
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		switch g.debug.GameMode {
		case debug.Debug:
			g.debug.GameMode = debug.User
			g.debug.DebugText = ""
			g.debug.DebugParameter[debug.Print] = false
			ev := events.ButtonClickedEvent{ButtonText: "Mode"}
			g.gameLog.GlobalEventHub.Publish(ev)
		case debug.User:
			g.debug.GameMode = debug.Debug
			ev := events.ButtonClickedEvent{ButtonText: "Mode"}
			g.gameLog.GlobalEventHub.Publish(ev)
			g.debugModeParameterPrinterUpdater()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		colMap, err := util.LoadCollisions()
		if err != nil {
			log.Fatal("error while trying to hot reload collisions", err)
		}
		g.collisionMap = colMap
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		testEvent := events.MoneyAvailable{Amount: 10.0}
		g.gameLog.GlobalEventHub.Publish(testEvent)
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

	if inpututil.IsKeyJustPressed(ebiten.Key9) {
		e := entities.SavedFish{FishType: string(entities.Fish), Size: 3}
		loader.InitFish(e, g.environment, g.gameLog.GlobalEventHub, g.collisionMap)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		switch g.lightingState {
		case NightLight:
			g.lightingState = Day
		case Day:
			g.lightingState = NightLight
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

func (g *FishScene2) debugModeParameterPrinterUpdater() {
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

func (g *FishScene2) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth > 0 {
		return outsideWidth, outsideHeight
	}
	return ScreenWidth, ScreenHeight
}

func (g *FishScene2) subs(colMap map[string]image.Rectangle) {

	g.uiSubs()
	g.soundSubs()
	g.creatureSubs(colMap)

}

func (g *FishScene2) uiSubs() {

	var propPicked string

	g.gameLog.GlobalEventHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		println(ev.ButtonText, "button event received")
		switch ev.ButtonText {
		case "Save":
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(events.WindowClosed{}, func(e tasks.Event) {
		g.gameState.MouseFlags.WindowOpen = false
	})

	g.gameLog.GlobalEventHub.Subscribe(events.WindowOpened{}, func(e tasks.Event) {
		g.gameState.MouseFlags.WindowOpen = true
	})

	g.gameLog.GlobalEventHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		if ev.ButtonText == "Go to Bed?: Yes" {
			g.timers["sceneTransition"].TurnOn()
			ev2 := events.DayOver{}
			//change state to let the game know to draw unlit art
			//g.lightingState = Night
			//turn of normal maps:
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
				p := entities.LoadProp(propPicked, g.collisionMap, g.gameLog.GlobalEventHub)
				ent := &entities.Entity{PropData: p, Sprite: p.Sprite, EventHub: g.gameLog.GlobalEventHub}
				ent.Sprite.Z = 1
				entities.RegisterEntity(ent)
				if p.Sprite2 != nil {
					ent2 := &entities.Entity{PropData: p, Sprite: p.Sprite2, EventHub: g.gameLog.GlobalEventHub}
					ent2.Sprite.Z = 0
					entities.RegisterEntity(ent2)
				}

				println("castle or whatever added")
			}
		}

		g.gameLog.GlobalEventHub.Subscribe(input.MouseButtonPressedUISpriteActivity{}, func(e tasks.Event) {
			ev := e.(input.MouseButtonPressedUISpriteActivity)
			//handled click relates to only proccessing one event per game tick
			if g.timers["pointGeneratedTimer"].TimerState == util.Done && !g.gameState.MouseFlags.HandledClick {
				g.gameState.MouseFlags.HandledClick = true
				pt := ev.Point.Clone()
				if pt.Tag == "left" {
					pt.X = pt.X + 50 - rand.Float32()*10
					pt.Y += 50
				} else {
					pt.X = pt.X - 50 + rand.Float32()*10
					pt.Y += 50
				}
				entities.NewParticle(pt, g.collisionMap["tank"], g.gameLog.GlobalEventHub)

			}
		})
	})
}

func LoadPHEffect(kind string, hub *tasks.EventHub) {
	img, err := util.LoadImageAssetAsEbitenImage("uiSprites/" + kind)
	if err != nil {
		log.Fatal(err)
	}

	x, y := util.GetScaledCursorPosition()
	parameters := make(map[string]any)
	parameters["tag"] = kind
	loader.MakeSpriteEntity(img, float32(x), float32(y), loader.SpriteEntFlags{Unfocusable: true, Updater: true,
		UpdateFunc: entities.PHModifier, Zlayer: 3, Parameters: parameters, EventHub: hub})
}

func (g *FishScene2) soundSubs() {
	g.gameLog.GlobalEventHub.Subscribe(entities.SendData{}, func(e tasks.Event) {
		ev := e.(entities.SendData)
		if ev.DataFor == "soundFx" && ev.Data == "particle entered water" {
			g.gameLog.SoundPlayer.Play(soundFX.PlopSound)
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(events.UISpriteAction{}, func(e tasks.Event) {
		ev := e.(events.UISpriteAction)
		if ev.UiSprite == "fishFood" && ev.UiSpriteAction == "put back" {
			g.gameLog.SoundPlayer.Play(soundFX.PickUpOne)
			return
		}
		if ev.UiSprite == "fishFood" && ev.UiSpriteAction == "picked up" {
			g.gameLog.SoundPlayer.Play(soundFX.PlopSound)
			return
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(tasks.TaskCompleted{}, func(e tasks.Event) {
		g.gameLog.SoundPlayer.Play(soundFX.WhiteBoardMarker1)
		g.gameLog.SoundPlayer.Play(soundFX.SuccessMusic)
		ev := e.(tasks.TaskCompleted)
		if len(g.gameLog.Tasks) > ev.Slot {
			g.currentTask++ // this is zero indexed but slot is not so the current index is the just finished slot // FIX
			g.timers["publishNewTask"].TurnOn()
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(events.UISpriteAction{}, func(e tasks.Event) {
		ev := e.(events.UISpriteAction)
		if ev.UiSprite == string(entities.PiggyBank) && ev.UiSpriteAction == "clicked" {
			g.gameLog.SoundPlayer.Play(soundFX.Coins1)
		}
	})
	g.gameLog.GlobalEventHub.Subscribe(events.UISpriteAction{}, func(e tasks.Event) {
		ev := e.(events.UISpriteAction)
		if ev.UiSprite == string(entities.Phreader) && ev.UiSpriteAction == "clicked" {
			g.gameLog.SoundPlayer.Play(soundFX.CardBoard)
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
		g.gameLog.SoundPlayer.Play(soundFX.WhiteBoardMarker2)
	})

}

func (g *FishScene2) printGameMode(screen *ebiten.Image) {
	DebugText(g.debug.DebugText, screen)
}

func (g *FishScene2) saveUISpritePositions() {

	spMap := make(map[string]drawables.SavePositionData)
	//THIS IS AWFUL BECAUSE I MADE TOO MANY UI SPRITE TYPES DEAL WITH SOME DAY
	//check for ui sprites

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

func positionCollisionBaseOnSprite(col image.Rectangle, prop *entities.Entity) image.Rectangle {
	width := col.Dx()
	height := col.Dy()
	minPointX := int(prop.Sprite.X) + col.Min.X
	minPointY := int(prop.Sprite.Y) + col.Min.Y
	maxPointX := minPointX + width
	maxPointY := minPointY + height
	newPos := image.Rect(minPointX, minPointY, maxPointX, maxPointY)
	return newPos
}

func (g *FishScene2) creatureSubs(colMap map[string]image.Rectangle) {
	g.gameLog.GlobalEventHub.Subscribe(input.MouseButtonPressedUISpriteActivity{}, func(e tasks.Event) {
		ev := e.(input.MouseButtonPressedUISpriteActivity)
		//handled click relates to only proccessing one event per game tick
		if g.timers["pointGeneratedTimer"].TimerState == util.Done && !g.gameState.MouseFlags.HandledClick {
			g.gameState.MouseFlags.HandledClick = true
			pt := ev.Point.Clone()
			pt.X = pt.X - 50 + rand.Float32()*10
			pt.Y += 50
			entities.NewParticle(pt, colMap["tank"], g.gameLog.GlobalEventHub)

		}
	})

	g.gameLog.GlobalEventHub.Subscribe(events.NewProp{}, func(e tasks.Event) {
		ev := e.(events.NewProp)
		prop, exists := entities.GetEntity(ev.PropId)
		if !exists {
			log.Fatal("new prop event returning prop that doesnt exist ")
		}

		g.debug.DebugRect.GivePoint(image.Point{int(prop.Sprite.X), int(prop.Sprite.Y)})

		println("got new prop event", ev.Name)
		for key, col := range g.collisionMap {
			if strings.Contains(key, ev.Name) {
				log.Println("adding active collision:", key)
				newPos := positionCollisionBaseOnSprite(col, prop)
				collision := Collision{z: prop.Sprite.Z, Rectangle: newPos}
				g.activatedCollisions = append(g.activatedCollisions, collision)
			}
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(events.PurchaseSuccessful{}, func(e tasks.Event) {
		ev := e.(events.PurchaseSuccessful)
		log.Printf("New Purchase:%s ", ev.Purchase)
		if ev.Purchase == "ph+" || ev.Purchase == "ph-" {
			LoadPHEffect(ev.Purchase, g.gameLog.GlobalEventHub)
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(events.InsufficientFunds{}, func(e tasks.Event) {
		id := graphics.NewFadeInTextGraphic("Get a job Kid!", float64(200), float64(100))
		time2.AfterFunc(4*time2.Second, func() { graphics.DeInitGraphicId(id) })
	})

	g.gameLog.GlobalEventHub.Subscribe(entities.CreatureReachedPoint{}, func(e tasks.Event) {
		if !creatureManager.allFishFed && entities.CheckIfAllFishFed() {
			creatureManager.allFishFed = true
			g.gameLog.GlobalEventHub.Publish(entities.AllFishFed{})
		}
	})
	g.gameLog.GlobalEventHub.Subscribe(events.PurchaseSuccessful{}, func(e tasks.Event) {
		ev := e.(events.PurchaseSuccessful)
		if entities.FishList(ev.Purchase) != "" {
			log.Printf("New Purchase:%s ", ev.Purchase)
			ent := loader.InitFish(entities.SavedFish{FishType: ev.Purchase, Size: 1}, g.environment, g.gameLog.GlobalEventHub, g.collisionMap)
			ev2 := events.CloseWindow{}
			entities.RegisterEntity(ent)
			entities.UnFocus(g.gameState.FocusedEntity.Id)
			g.gameLog.GlobalEventHub.Publish(ev2)
		}
	})
}
