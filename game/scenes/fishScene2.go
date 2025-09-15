package scenes

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/daySystem"
	"github.com/acoco10/fishTankWebGame/game/debug"
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/input"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/physics"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/system"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/ebitenui/ebitenui"
	eInput "github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"image"
	"log"
	"math/rand"
	"os"
	"strings"
)

type lightingState uint8

const (
	NightLight lightingState = iota
	Night
	Day
)

var backGroundImgShelfHeight = 150

type FishScene2 struct {
	zoomScreen              *ebiten.Image
	state                   FishSceneState
	loaded                  bool
	tankSize                image.Rectangle
	ui                      *ebitenui.UI
	gameLog                 *sceneManagement.GameLog
	timers                  map[string]*util.Timer
	returnScene             sceneManagement.SceneId
	collisionMap            map[string]image.Rectangle
	environment             *system.Environment
	playerState             *entities.Player
	images                  *loader.BackGroundImages
	debug                   *debug.DebugData
	currentTask             int
	testProp                *entities.StructureProp
	gameState               *entities.GameState
	renderBuffers           [10]*ebiten.Image
	lightingState           lightingState
	lightEntManager         *entities.LightingEntManager
	globalLightingShaderMap map[lightingState]*ebiten.Shader
	shaderParamsMap         map[lightingState]map[string]any
	smallerResolution       *ebiten.Image
	whiteBoardSprite        *entities.WhiteBoardSprite
	taskRelatedEventQueue   []tasks.Event
	zBounds                 [13]image.Rectangle
	propData                entities.PropData
	fid                     uint32
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
	propPts, err := util.LoadCoords()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i <= 12; i++ {
		layerName := fmt.Sprintf("z%d", i)

		rect := collisionMap["tank"]
		rect.Min.Y -= 5
		rect.Max.Y -= 8
		//top layer is obfuscated by tank front layer, dont worry about it for now
		//2 pixels per z layer (first iteration
		rect.Max.Y -= 12 * 2 //get us to the rear most layer
		rect.Max.Y += i * 2  //add up to the layer were on
		rect.Max.X -= 12 * 2
		rect.Min.X += 12 * 2
		rect.Max.X += i * 2
		rect.Min.X -= i * 2
		g.collisionMap[layerName] = rect
		log.Printf("setting z bounds: %d \n max y = %d", i, rect.Max.Y)
		g.zBounds[i] = rect
	}
	g.propData = entities.PropData{PtMap: propPts, CollisionMap: collisionMap, PlacementParams: make(map[string]any), ZBounds: g.zBounds}
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

	for buf, _ := range g.renderBuffers {
		g.renderBuffers[buf] = ebiten.NewImage(ScreenWidth, ScreenHeight)
	}

	g.globalLightingShaderMap = make(map[lightingState]*ebiten.Shader)

	g.globalLightingShaderMap[Day] = registry.ShaderMap["DayLight"]
	g.globalLightingShaderMap[NightLight] = registry.ShaderMap["OnePointLighting"]

	tankRect := g.collisionMap["tank"]

	g.shaderParamsMap = make(map[lightingState]map[string]any)

	nightLightParams := make(map[string]any)
	nightLightParams["LightPoint"] = [3]float64{
		float64(g.collisionMap["tank"].Min.X + 194),
		float64(g.collisionMap["tank"].Min.Y - 52),
		25.0}
	nightLightParams["LightWidth"] = 120.0
	nightLightParams["TankRect"] = [4]float64{
		float64(tankRect.Min.X),
		float64(tankRect.Min.Y),
		float64(tankRect.Max.X),
		float64(tankRect.Max.Y)}

	g.shaderParamsMap[NightLight] = nightLightParams
	g.shaderParamsMap[Day] = make(map[string]any) //no params as of now

	tankX := g.images.FishTank.Bounds().Dx()
	tankY := g.images.FishTank.Bounds().Dy()

	startingX := int(ScreenWidth * 0.2)
	startingY := ScreenHeight - backGroundImgShelfHeight - g.images.FishTank.Bounds().Dy()

	tankRect = image.Rect(startingX, startingY, tankX+startingX, tankY+startingY)
	g.tankSize = tankRect
	g.collisionMap["tankRect"] = tankRect

	fishtankFrontSpriteFlags := loader.SpriteEntFlags{Unfocusable: true, Zlayer: 12}
	flid := loader.MakeSpriteEntity(g.images.FishTankFrontLayerNoLightSmaller, float32(g.tankSize.Min.X), float32(g.tankSize.Min.Y), fishtankFrontSpriteFlags)
	g.fid = flid

	ftdl := entImportableLoaders.LoadEffect("fishTankDayLight")
	ftdl.Unfocusable = true
	ftdl.GetAnimation().SpeedInTPS = 32
	ftdl.X = float32(g.tankSize.Min.X)
	ftdl.Y = float32(g.tankSize.Min.Y)
	ftid := entities.RegisterEntity(&entities.Entity{Sprite: ftdl, Z: 1})

	ftn := entImportableLoaders.LoadEffect("fishTankNightLight")
	ftn.Unfocusable = true
	ftn.GetAnimation().SpeedInTPS = 32
	ftn.X = float32(g.tankSize.Min.X)
	ftn.Y = float32(g.tankSize.Min.Y)
	ftnId := entities.RegisterEntity(&entities.Entity{Sprite: ftn, Z: 1})

	fishTankFrontLayerDy := g.images.FishTankFrontLayerNoLightSmaller.Bounds().Dy()
	fishTankHeightDy := g.images.FishTank.Bounds().Dy()
	y := fishTankFrontLayerDy - fishTankHeightDy
	y = g.tankSize.Min.Y - y

	nightLight1 := entImportableLoaders.LoadEffect("fishTankNightLightTransparency1")
	nightLight1.Unfocusable = true
	nightLight1.GetAnimation().SpeedInTPS = 32
	nightLight1.X = float32(g.tankSize.Min.X)
	nightLight1.Y = float32(g.tankSize.Min.Y)
	nlId := entities.RegisterEntity(&entities.Entity{Sprite: nightLight1, Z: 4})

	nightLight2 := entImportableLoaders.LoadEffect("fishTankNightLightTransparency2")
	nightLight2.Unfocusable = true
	nightLight2.GetAnimation().SpeedInTPS = 32
	nightLight2.X = float32(g.tankSize.Min.X)
	nightLight2.Y = float32(g.tankSize.Min.Y)
	nl2Id := entities.RegisterEntity(&entities.Entity{Sprite: nightLight2, Z: 8})

	waveSprite := entImportableLoaders.LoadEffect("lightRay")
	waveSprite.GetAnimation().SpeedInTPS = 32
	waveSprite.X = float32(g.tankSize.Min.X)
	waveSprite.Y = float32(g.tankSize.Min.Y)
	dayWaveID := entities.RegisterEntity(&entities.Entity{Sprite: waveSprite, Z: 5})

	dither := entImportableLoaders.LoadEffect("ditheredCaustic")
	dither.X = float32(startingX + 5)
	dither.Y = float32(g.collisionMap["tank"].Min.Y - 25)
	entities.RegisterEntity(&entities.Entity{Sprite: dither, Z: 5})

	tImg := ebiten.NewImage(dither.GetAnimation().SpriteWidth, dither.GetAnimation().SpriteHeight)
	dopts := &ebiten.DrawImageOptions{}
	dopts.GeoM.Scale(1.4, 1.4)
	dopts.GeoM.Translate(-100, -100)
	tImg.DrawImage(dither.GetAnimation().GetFirstFrameAsStaticImage(), dopts)
	tSprite := &sprite.Sprite{Img: tImg, X: dither.X, Y: dither.Y}

	ent := &entities.Entity{Sprite: tSprite, Z: 1}
	entities.RegisterEntity(ent)

	g.lightEntManager = &entities.LightingEntManager{}
	g.lightEntManager.NightEnts = append(g.lightEntManager.NightEnts, nlId, nl2Id, ftnId)
	g.lightEntManager.DayEnts = append(g.lightEntManager.DayEnts, dayWaveID, ftid)
	g.lightEntManager.Subscribe(gameLog.GlobalEventHub)

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

	mainUI, err := ui.LoadMainFishMenu(ScreenWidth, ScreenHeight, gameLog.GlobalEventHub)
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

		//g.gameLog.Save.GoldFish = []entities.SavedFish{entities.SavedFish{FishType: "fish", Size: 2}, entities.SavedFish{FishType: "fish", Size: 3}}

	}
	g.zoomScreen = ebiten.NewImage(registry.Config.ScreenWidth, registry.Config.ScreenHeight)
	//store := system.NewStore(g.gameLog.GlobalEventHub)
	//g.store = store
	g.LoadStuff()

	ev2 := events.FishTankLayout{
		Rectangle: g.tankSize,
	}
	g.gameLog.GlobalEventHub.Publish(ev2)

	fishScene2UISprites := []entities.Label{
		entities.WhiteBoard,
		entities.PiggyBank,
		entities.Thermometer,
		entities.Magazine,
		entities.GrandpasJournal,
		entities.Phreader,
		entities.FishFood,
		entities.Skimmer,
		entities.LightSwitch,
	}

	_, wbSprite := loader.LoadAllEntities(fishScene2UISprites, g.gameLog.Save.Fish, g.environment, g.gameLog.GlobalEventHub, g.collisionMap)
	g.whiteBoardSprite = wbSprite
	g.loaded = true

	daySystem.LoadDaysTasks(g.gameLog)
	log.Println("----fishScene2.firstLoad() finished----")

}

func (g *FishScene2) OnEnter() {
	g.images.OffScreen.Clear()
	g.smallerResolution.Clear()
	if g.zoomScreen != nil {
		g.zoomScreen.Clear()
	}

	cu := entities.CreateCursorUpdater(g.gameLog.GlobalEventHub)
	eInput.SetCursorUpdater(cu)

	entities.LoadEffects(g.zBounds)

	mouseFlags := &input.MouseFlags{HandledClick: false, CursorOccupied: false}
	g.gameState = &entities.GameState{MouseFlags: mouseFlags, CollisionMap: g.collisionMap, Zbounds: g.zBounds, CursorUpdater: cu}

	g.debug = &debug.DebugData{DebugParameter: make(map[debug.DebugOption]bool), EventHub: g.gameLog.GlobalEventHub, GameState: *g.gameState, PropData: g.propData}

	g.gameLog.SongPlayer.Play(soundFX.ElectricBuzz)

	g.debug.DebugParameter[debug.Print] = false
	g.debug.DebugParameter[debug.Position] = false

	log.Println("----FishScene2 OnEnter() called----")

	//No music on the base level as of now

	g.timers["songTimer"].TurnOn()

	bs := entities.NewBubbleSystem(300.0, float64(g.zBounds[1].Dy()), g.zBounds[1])
	entities.RegisterEntity(&entities.Entity{ParticleSystem: bs, Z: 1, Sprite: bs.Sprite})
	bs.BubbleSubscriptions(g.gameLog.GlobalEventHub)

	bs2 := entities.NewBubbleSystem(300.0, float64(g.zBounds[3].Dy()), g.zBounds[3])
	entities.RegisterEntity(&entities.Entity{ParticleSystem: bs2, Z: 3, Sprite: bs2.Sprite})
	bs2.BubbleSubscriptions(g.gameLog.GlobalEventHub)

	bs3 := entities.NewBubbleSystem(300.0, float64(g.zBounds[5].Dy()), g.zBounds[5])
	entities.RegisterEntity(&entities.Entity{ParticleSystem: bs3, Z: 5, Sprite: bs3.Sprite})

	bs3.BubbleSubscriptions(g.gameLog.GlobalEventHub)

	totallyRandBubbles := entities.NewGenericParticleSystem(300.0, float64(g.zBounds[1].Min.Y), g.zBounds[5], 7)
	totallyRandBubbles.PConfig = &entities.ParticleConfig{XVariance: 80, YVariance: 80, XVelocityVariance: 20, YVelocityVariance: 20, MaxLife: 0}
	totallyRandBubbles.SpawnRate = 20
	totallyRandBubbles.MaxParticles = 75
	totallyRandBubbles.EndAfter = 100
	entities.RegisterEntity(&entities.Entity{ParticleSystem: totallyRandBubbles, Z: 10, Sprite: totallyRandBubbles.Sprite})
	totallyRandBubbles2 := entities.NewGenericParticleSystem(300.0, float64(g.zBounds[1].Min.Y), g.zBounds[10], 6)
	totallyRandBubbles2.PConfig = &entities.ParticleConfig{XVariance: 80, YVariance: 80, XVelocityVariance: 12, YVelocityVariance: 12, MaxLife: 0}
	totallyRandBubbles2.SpawnRate = 10
	totallyRandBubbles2.MaxParticles = 50
	totallyRandBubbles2.EndAfter = 100
	entities.RegisterEntity(&entities.Entity{ParticleSystem: totallyRandBubbles2, Z: 5, Sprite: totallyRandBubbles2.Sprite})
	g.returnScene = sceneManagement.FishTank

	println("made it to new day")
	if g.gameLog.Day != g.state.lastDayEntered {
		g.configureNewDay()
		entities.UpdateCursorForEntitiesWNormals([]float64{0, 0, 100})
		log.Println("----FishScene2 OnEnter() finished----")
		return
	}

	g.configureExistingDay()

	log.Println("----FishScene2 OnEnter() finished----")
}

func (g *FishScene2) configureExistingDay() {
	g.lightingState = NightLight
	ev := events.LightEvent{Day: false}
	g.gameLog.GlobalEventHub.Publish(ev)
}

func (g *FishScene2) configureNewDay() {
	ev := events.NewDay{NTasks: len(g.gameLog.TaskManager.Tasks), Day: g.gameLog.Day}
	g.lightingState = Day
	g.gameLog.GlobalEventHub.Publish(ev)
	g.state.lastDayEntered = g.gameLog.Day

	if len(g.gameLog.TaskManager.Tasks) > 0 {
		g.gameLog.TaskManager.Activate()
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
	g.timers["bubbles"] = util.NewTimer(20.0)
	g.timers["bubbles"].TurnOn()
	g.timers["turnOffBubbles"] = util.NewTimer(10.0)
	g.timers["placementCooldown"] = util.NewTimer(1.5)
}

func (g *FishScene2) IsLoaded() bool {
	return g.loaded
}

func (g *FishScene2) Update() (sceneManagement.SceneId, error) {
	g.ui.Update()
	g.gameState.MouseFlags.HandledClick = false

	if !g.gameState.MouseFlags.WindowOpen {
		g.updateInput()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.gameLog.GlobalEventHub.Publish(events.CloseWindow{})
	}

	entities.UpdateEntities(g.gameState)
	g.whiteBoardSprite.Update()

	g.gameLog.SongPlayer.Update()
	g.gameLog.SoundPlayer.Update()

	g.updateTimers()

	if g.debug.GameMode == debug.Debug {
		g.debug.Update()
		g.debugInputCheck()
	}

	graphics.UpdateGraphics()

	if g.gameState.FocusedEntity != nil && !g.gameState.FocusedEntity.Sprite.Focused {
		g.gameState.FocusedEntity = nil
	}

	return g.returnScene, nil
}

func (g *FishScene2) DrawOffScreen() {

	g.images.OffScreen.Clear()

	opts := &ebiten.DrawImageOptions{}

	switch g.lightingState {
	case NightLight:
		g.images.OffScreen.DrawImage(g.images.RoomBackGroundNight, opts)
	case Day:
		g.images.OffScreen.DrawImage(g.images.RoomBackground, opts)
	}
	opts.GeoM.Reset()
	entities.DrawEntities(g.images.OffScreen, g.gameState)
	graphics.DrawScaledGraphics(g.images.OffScreen)

}

func (g *FishScene2) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.zoomScreen.Clear()
	//call offscreen function to render everything we want to apply global shaders to
	g.DrawOffScreen()

	dOpts := &ebiten.DrawImageOptions{}

	if g.debug.GameMode == debug.Debug {
		g.debug.Draw(g.images.OffScreen)
		g.debug.Draw(g.zoomScreen)
		for _, rect := range g.gameState.ActiveCollisions {
			util.StrokeRectFromImageRect(rect.Rectangle, g.smallerResolution, colornames.Greenyellow, false)
		}
	}

	if registry.Config.Zoom {

		dOpts.GeoM.Reset()

		dOpts.GeoM.Scale(registry.Config.ZoomFactor, registry.Config.ZoomFactor)
		dOpts.GeoM.Translate(registry.Config.ZoomOffSetX, registry.Config.ZoomOffSetY)
		g.zoomScreen.DrawImage(g.images.OffScreen, dOpts)

		dOpts.GeoM.Reset()
		dOpts.GeoM.Translate(0, registry.Config.YOffsetF)
		dOpts.GeoM.Scale(registry.Config.ResolutionScalingF, registry.Config.ResolutionScalingF)

		screen.DrawImage(g.zoomScreen, dOpts)
		entities.DrawNonZoomedEntities(screen)
		graphics.DrawScaledGraphicsOnMainScreen(screen)
		graphics.DrawUnScaledGraphics(screen)
		g.ui.Draw(screen) //also drawing our custom cursor
		return
	} else {
		dOpts.GeoM.Reset()

		dOpts.GeoM.Translate(0, registry.Config.YOffsetF)
		dOpts.GeoM.Scale(registry.Config.ResolutionScalingF, registry.Config.ResolutionScalingF)

		ebitenutil.DebugPrintAt(g.images.OffScreen, fmt.Sprintf("TPS: %0.2f\n", ebiten.ActualTPS()), 30, 0)
		ebitenutil.DebugPrintAt(g.images.OffScreen, fmt.Sprintf("Window Open:%t, Focused Entity = %v, HoveredUiSprite = %v", g.gameState.MouseFlags.WindowOpen, g.gameState.FocusedEntity, g.gameState.HoveredUiSprite), 30, 20)
		ebitenutil.DebugPrintAt(g.images.OffScreen, fmt.Sprintf("Environment Temperatre: %f, Environment natural PH: %f, Environment Modified PH: %f", g.environment.Temperature, g.environment.NaturalPHLevel, g.environment.ModifiedPHLevel), 30, 40)

		screen.DrawImage(g.images.OffScreen, dOpts)

		graphics.DrawUnScaledGraphics(screen)
		//g.whiteBoardSprite.DstImg.Draw(screen)
		graphics.DrawScaledGraphicsOnMainScreen(screen)

		g.ui.Draw(screen)

	}

}

func (g *FishScene2) positionModeUpdate() {

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
			if len(g.gameLog.TaskManager.Tasks) > g.gameLog.TaskManager.CurrentTask {
				timer.TurnOff()
				g.gameLog.TaskManager.Activate()
			} else {
				timer.TurnOff()
				g.currentTask = 0
			}
		}

		if key == "bubbles" && state == util.Done {
			g.gameLog.GlobalEventHub.Publish(entities.TurnOnBubbles{})
			g.timers["turnOffBubbles"].TurnOn()
		}

		if key == "turnOffBubbles" && state == util.Done {
			g.gameLog.GlobalEventHub.Publish(entities.TurnOffBubbles{})
			g.timers["turnOffBubbles"].TurnOff()
		}

		if key == "placementCooldown" && state == util.Done {
			pickedSoFar, ok := g.propData.PlacementParams["quantity"].(int)
			if !ok {
				g.timers["placementCooldown"].TurnOff()
				return
			}

			pickedSoFar -= 1
			g.propData.PlacementParams["quantity"] = pickedSoFar
			if pickedSoFar > 0 {
				entities.LoadPlaceMentReticule(g.zBounds, "Plant", g.gameLog.GlobalEventHub)
			}
			g.propData.PlacementParams["fertilizer"] = false
			g.timers["placementCooldown"].TurnOff()
		}
	}
}

func (g *FishScene2) updateInput() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		entities.FocusOnClickedEntity(g.gameState)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		entities.LoadPlaceMentReticule(g.zBounds, "Plant", g.gameLog.GlobalEventHub)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		if g.gameState.FocusedEntity != nil && !g.gameState.MouseFlags.WindowOpen {
			entities.UnFocus(g.gameState.FocusedEntity.Id)
			g.gameState.FocusedEntity = nil
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.HandleZoom()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		//registry.Config.Set(registry.Debug, true)
		g.debug.GameMode = debug.Debug
		ev := events.ButtonClickedEvent{ButtonText: "Mode"}
		g.gameLog.GlobalEventHub.Publish(ev)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		/*	img, err := util.LoadImageAssetAsEbitenImage("textures/debris1")
			if err != nil {
				log.Fatal(err)
			}*/
		//sp := &sprite.Sprite{Y: 400, X: 400, Img: img, Unfocusable: true}
		rx := rand.Float64() * 3
		phys := physics.NewPhysicsBody(400+rx, 10)
		phys.Bounds = &physics.Rect{float64(g.zBounds[0].Min.X), float64(g.zBounds[0].Min.Y), float64(g.zBounds[0].Dx()), float64(g.zBounds[0].Dy())}
		//ent := &entities.Entity{Sprite: sp, Physics: phys, Z: 6}
		//entities.RegisterEntity(ent)
		//g.gameState.PhysicsObjects = append(g.gameState.PhysicsObjects, ent.Id)
	}

}

func (g *FishScene2) HandleZoom() {
	if !registry.Config.Zoom {
		g.publish(events.Zoom{})
		registry.Config.Set(registry.Zoom, true)

		if g.gameState.FocusedEntity != nil {
			g.gameState.PreZoomFocusedEntity = g.gameState.FocusedEntity
			entities.UnFocus(g.gameState.FocusedEntity.Id)
		}

		flidEnt, exists := entities.GetEntity(g.fid)
		if !exists {
			log.Fatal("flid Ent does not exist")
		}

		flidEnt.Draw = false

	} else {
		g.publish(events.UnZoom{})
		registry.Config.Set(registry.Zoom, false)

		if g.gameState.PreZoomFocusedEntity != nil {
			entities.ReFocus(g.gameState.PreZoomFocusedEntity.Id)
			g.gameState.FocusedEntity = g.gameState.PreZoomFocusedEntity
		}
		g.gameState.PreZoomFocusedEntity = nil

		flidEnt, exists := entities.GetEntity(g.fid)
		if !exists {
			log.Fatal("flid Ent does not exist")
		}

		flidEnt.Draw = true
	}
}

func (g *FishScene2) debugInputCheck() {

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
	if g.debug.DebugOption == debug.FishSpawn {
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			e := entities.SavedFish{FishType: string(entities.GoldFish), Size: 3}
			loader.InitFish(e, g.environment, g.gameLog.GlobalEventHub, g.collisionMap)
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			e := entities.SavedFish{FishType: string(entities.AngelFish), Size: 2}
			loader.InitFish(e, g.environment, g.gameLog.GlobalEventHub, g.collisionMap)
		}
		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			e := entities.SavedFish{FishType: string(entities.Kirbensis), Size: 2}
			loader.InitFish(e, g.environment, g.gameLog.GlobalEventHub, g.collisionMap)
		}

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		switch g.lightingState {
		case NightLight:
			entities.UpdateCursorForEntitiesWNormals([]float64{0, 0, 100})
			g.lightingState = Day
			g.gameLog.GlobalEventHub.Publish(events.LightEvent{Day: true})
		case Day:
			lightingCursor := g.shaderParamsMap[NightLight]["LightPoint"].([3]float64)
			entities.UpdateCursorForEntitiesWNormals([]float64{lightingCursor[0], lightingCursor[1], lightingCursor[2]})
			g.gameLog.GlobalEventHub.Publish(events.LightEvent{Day: false})
			g.lightingState = NightLight
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) && g.debug.DebugOption == debug.Position {
		g.saveUISpritePositions()
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
	g.entitySubs()
	g.soundSubs()
	g.creatureSubs(colMap)
	g.PropSubs()

}

func (g *FishScene2) subscribe(eventType tasks.Event, handler func(e tasks.Event)) {
	g.gameLog.GlobalEventHub.Subscribe(eventType, handler)
}

func (g *FishScene2) publish(event tasks.Event) {
	g.gameLog.GlobalEventHub.Publish(event)
}

func (g *FishScene2) entitySubs() {
	g.subscribe(entities.PlacementPicked{}, func(e tasks.Event) {
		ev := e.(entities.PlacementPicked)
		entities.LoadPlant(ev, g.propData, g.gameLog.GlobalEventHub)
		g.timers["placementCooldown"].TurnOn()
	})
	g.subscribe(events.ItemUsed{}, func(e tasks.Event) {
		ev := e.(events.ItemUsed)
		if ev.Name == string(entities.Fertilizer) {
			g.propData.PlacementParams["fertilizer"] = true
		}
	})

}

func (g *FishScene2) uiSubs() {
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
				entities.NewParticle(pt, g.collisionMap["tank"], g.gameLog.GlobalEventHub, g.zBounds)

			}
		})
	})

	g.gameLog.GlobalEventHub.Subscribe(tasks.AllTasksCompleted{}, func(e tasks.Event) {
		g.lightingState = NightLight
		ev := events.LightEvent{Day: false}
		g.gameLog.GlobalEventHub.Publish(ev)
		fishCount := 0.0
		happyFish := 0
		for _, ent := range entities.LiveList {
			if ent.CreatureData != nil {
				fishCount += 1
				if ent.CreatureData.Stress < 2 {
					happyFish += 1
				}
			}
		}
		income := 1.00 * fishCount
		happyFishBonus := .25 * float64(happyFish)
		mev := events.MoneyAvailable{Amount: income}

		msg1 := fmt.Sprintf("Today's Income:$%0.2f", income)
		ev1 := entities.WriteToWhiteBoard{Msg: msg1, PreferredPosition: "center", Later: true, EventToPublish: mev}

		msg2 := fmt.Sprintf("Happy Fish Bonus:$%0.2f", happyFishBonus)
		ev2 := entities.WriteToWhiteBoard{Msg: msg2, PreferredPosition: "center", Later: true, EventDriven: events.MoneyAdded{}}
		if happyFishBonus > 0 {
			mev2 := events.MoneyAvailable{Amount: happyFishBonus}
			ev2.EventToPublish = mev2
		}

		g.publish(ev1)
		g.publish(ev2)

	})
}

func (g *FishScene2) PropSubs() {
	zoomedBefore := false
	g.subscribe(events.PlacementMode{}, func(e tasks.Event) {
		g.gameState.MouseFlags.WindowOpen = true
		zoomedBefore = registry.Config.Zoom
		g.HandleZoom()
		g.publish(events.Zoom{})
	})

	g.subscribe(events.NewProp{}, func(e tasks.Event) {
		if !zoomedBefore {
			g.HandleZoom()
			g.publish(events.Zoom{})
		}
		g.gameState.MouseFlags.WindowOpen = false
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
	g.gameLog.GlobalEventHub.Subscribe(entities.TurnOnBubbles{}, func(e tasks.Event) {
		g.gameLog.SoundPlayer.Play(soundFX.WaterBubbles)
	})
	g.gameLog.GlobalEventHub.Subscribe(entities.TurnOffBubbles{}, func(e tasks.Event) {
		g.gameLog.SoundPlayer.Pause()
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
		if len(g.gameLog.TaskManager.Tasks) > ev.Slot {
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

	currentPos, err := loader.LoadSpritePositionData()
	if err != nil {
		log.Fatal(err)
	}

	for _, sprite := range currentPos {
		if _, ok := spMap[sprite.Name]; !ok {
			spMap[sprite.Name] = *sprite
		}
	}

	for _, ent := range entities.LiveList {
		if ent.UiData != nil {
			spMap[ent.UiData.Label] = drawables.SavePositionData{X: ent.Sprite.X, Y: ent.Sprite.Y, Name: ent.UiData.Label}
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
			entities.NewParticle(pt, colMap["tank"], g.gameLog.GlobalEventHub, g.zBounds)

		}
	})

	g.gameLog.GlobalEventHub.Subscribe(events.NewProp{}, func(e tasks.Event) {
		ev := e.(events.NewProp)
		prop, exists := entities.GetEntity(ev.PropId)
		if !exists {
			log.Fatal("new prop event returning prop that doesnt exist ")
		}

		println("got new prop event", ev.Name)

		for key, col := range g.collisionMap {
			if strings.Contains(strings.ToLower(key), strings.ToLower(ev.Name)) {
				log.Println("adding active collision:", key)
				newPos := positionCollisionBaseOnSprite(col, prop)
				collision := entities.FishCollision{Z: prop.Z, Rectangle: newPos}
				g.gameState.ActiveCollisions = append(g.gameState.ActiveCollisions, collision)
			}
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(events.PurchaseSuccessful{}, func(e tasks.Event) {
		ev := e.(events.PurchaseSuccessful)
		switch ui.PurchaseType(ev.PurchaseType) {
		case ui.Decoration:
			if entities.Decoration(ev.Purchase) == "" {
				log.Fatal("Recieved decoration purchase that has no const declaration ")
			}
			g.environment.AddTankModifier(ev.Purchase)
			entities.LoadProp(ev.Purchase, g.propData, g.gameLog.GlobalEventHub, entities.PlacementPicked{}, g.zBounds)

		case ui.Item:
			if ev.Purchase == "phBoost" || ev.Purchase == "phReduce" {
				var bought entities.Label
				switch ev.Purchase {
				case "phBoost":
					bought = entities.PhBoost
				case "phReduce":
					bought = entities.PhReduce
				}

				imgs, err := loader.LoadUiSpritesImgs(bought)
				if err != nil {
					log.Fatal(err)
				}

				g.gameLog.GlobalEventHub.Publish(events.CloseWindow{})
				x, y := util.GetScaledCursorPosition()
				uis := entities.NewUiSprite(g.environment, imgs, g.gameLog.GlobalEventHub, float32(x), float32(y), string(bought))
				ent := &entities.Entity{UiData: uis, Sprite: uis.Sprite}
				ent.Z = 13
				ent.UiData.ActivationRect = g.collisionMap["tank"]
				entities.AddUiSpriteXYUpdater(ent)
				ent.StateMachine = entities.InitStateMachine(entities.UpdateUseOnTank, nil, entities.UseOnTank, nil)
				ent.UiData.Flags["unDraw"] = true
				ent.UiData.Flags["clickForTime"] = true
				ent.UiData.Flags["oneOff"] = true
				ent.EventHub = g.gameLog.GlobalEventHub
				entities.RegisterEntity(ent)
			} else if ev.Purchase == "fertilizer" {
				imgs, err := loader.LoadUiSpritesImgs(entities.Fertilizer)
				if err != nil {
					log.Fatal(err)
				}
				uis := entities.NewUiSprite(g.environment, imgs, g.gameLog.GlobalEventHub, 735, 360, string(entities.Fertilizer))
				ent := &entities.Entity{UiData: uis, Sprite: uis.Sprite}
				ent.UiData.ActivationRect = g.collisionMap["tank"]
				ent.StateMachine = entities.InitStateMachine(nil, entities.UpdateUseOnTank, entities.AddUiSpriteXYUpdater, entities.UseOnTank)
				ent.UiData.Flags["clickForTime"] = true
				ent.UiData.Flags["resort"] = true
				ent.UiData.Flags["oneOff"] = true
				ent.EventHub = g.gameLog.GlobalEventHub
				entities.RegisterEntity(ent)
			}

		case ui.Plant:
			if ev.Purchase == "plantpack1" {
				g.propData.PlacementParams["quantity"] = 3
				entities.LoadPlaceMentReticule(g.zBounds, "Plant", g.gameLog.GlobalEventHub)
			}

		case ui.Fish:
			if entities.FishList(ev.Purchase) == "" {
				log.Fatal("Recieved decoration purchase that has no const declaration ")
			}

			log.Printf("New Purchase:%s ", ev.Purchase)
			ent := loader.InitFish(entities.SavedFish{FishType: ev.Purchase, Size: 1}, g.environment, g.gameLog.GlobalEventHub, g.collisionMap)
			ev2 := events.CloseWindow{}
			entities.RegisterEntity(ent)
			if g.gameState.FocusedEntity != nil {
				entities.UnFocus(g.gameState.FocusedEntity.Id)
			}
			g.gameLog.GlobalEventHub.Publish(ev2)
		}

		log.Printf("New Purchase:%s ", ev.Purchase)
	})

	g.gameLog.GlobalEventHub.Subscribe(events.InsufficientFunds{}, func(e tasks.Event) {
		graphics.NewFadeInTextGraphicCentered("Get a job Kid!", 120)
	})

	g.gameLog.GlobalEventHub.Subscribe(entities.CreatureReachedPoint{}, func(e tasks.Event) {
		if !creatureManager.allFishFed && entities.CheckIfAllFishFed() {
			creatureManager.allFishFed = true
			g.gameLog.GlobalEventHub.Publish(entities.AllFishFed{})
		}
	})

}
