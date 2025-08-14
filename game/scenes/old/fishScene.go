//go:build old

package old

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/fishTankWebGame/game"
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
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image"
	"log"
)

type lightingState uint8

const (
	NightLight lightingState = iota
	Night
	Day
)

type FishScene struct {
	backGroundParams     map[string]any
	loaded               bool
	tankSize             image.Rectangle
	sprites              [4][]drawables.Drawable
	allSprites           []drawables.Drawable
	ui                   *ebitenui.UI
	gameLog              *sceneManagement.GameLog
	timers               map[string]*util.Timer
	returnScene          sceneManagement.SceneId
	tutorialManager      *tutorial.Manager
	collisionMap         map[string]image.Rectangle
	store                *system.Store
	environment          *system.Environment
	playerState          *entities.Player
	images               *loader.BackGroundImages
	mouseFlags           *input.MouseFlags
	debug                *debug.DebugData
	state                *FishSceneState
	currentTask          int
	lightingShader       *ebiten.Shader
	offScreenShader      *ebiten.Shader
	tankShader           *ebiten.Shader
	tankShaderUpdater    func(map[string]any) map[string]any
	lightingShaderParams map[string]any
	tankShaderParams     map[string]any
	lightingState        lightingState
	shaderUpdater        func(map[string]any) map[string]any
	testProp             *entities.StructureProp //isolating for debug to be removed
	smallerResolution    *ebiten.Image
	resolutionScaling    int
	zoomFishMode         bool
	currentSpriteHovered drawables.Drawable
}

func NewFishScene(gameLog *sceneManagement.GameLog) *FishScene {

	roomImages, err := loader.LoadAllRoomBackGroundImages("assets/images/roomImagesSmaller")
	if err != nil {
		log.Fatal("error while loading fish tank room assets:", err)
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

	g.tankShaderParams = make(map[string]any)
	g.lightingShaderParams = make(map[string]any)
	g.backGroundParams = make(map[string]any)

	g.lightingShader = registry.ShaderMap["OnePointLighting"]
	g.offScreenShader = registry.ShaderMap["NormalMap"]
	g.tankShader = registry.ShaderMap["Water"]

	g.LoadTimers()

	collisionMap, err := util.LoadCollisions()
	if err != nil {
		log.Fatal(err)
	}
	_, exists := collisionMap["tank"]
	if !exists {
		log.Fatal()
	}

	g.collisionMap = collisionMap

	g.subs(collisionMap)

	//g.debugRect.RectState = geometry.Off

	//nighlight parameters
	tankRect := g.collisionMap["Tank"]

	g.lightingShaderParams["LightPoint"] = [2]float64{ScreenWidth / 2, ScreenHeight / 5}
	g.lightingShaderParams["LightWidth"] = 120.0
	g.lightingShaderParams["TankRect"] = [4]float64{
		float64(tankRect.Min.X),
		float64(tankRect.Min.Y),
		float64(tankRect.Max.X),
		float64(tankRect.Max.Y)}

	tankX := g.images.FishTank.Bounds().Dx()
	tankY := g.images.FishTank.Bounds().Dy()

	startingX := int(ScreenWidth * 0.2)
	startingY := ScreenHeight - backGroundImgShelfHeight - g.images.FishTank.Bounds().Dy()

	tankRect = image.Rect(startingX, startingY, tankX+startingX, tankY+startingY)
	g.tankSize = tankRect
	store := system.NewStore(g.gameLog.GlobalEventHub)
	g.store = &store

	g.tankShaderParams["TankRect"] = [4]float64{
		float64(collisionMap["tank"].Min.X),
		float64(collisionMap["tank"].Min.Y),
		float64(collisionMap["tank"].Max.X),
		float64(collisionMap["tank"].Max.Y)}

	g.tankShaderParams["Counter"] = 0

	g.tankShaderUpdater = shaders.UpdateCounter
	g.tankShader = registry.ShaderMap["Water"]

	g.environment = &system.Environment{}
	system.InitEnvironment(g.environment)
	g.environment.Subscribe(g.gameLog.GlobalEventHub)

	mainUI, _, err := ui.LoadMainFishMenu(ScreenWidth, ScreenHeight, gameLog.GlobalEventHub)
	if err != nil {
		log.Fatal("error loading scene")
	}

	g.ui = mainUI

	/*fishSceneUISprites := []entities.Label{
		entities.FishBook,
		entities.Records,
		entities.FishFood,
		entities.WhiteBoard,
		entities.PiggyBank,
		entities.Pillow,
		entities.Thermometer,
		entities.Magazine,
		entities.Door,
		entities.Phreader,
	}*/

	g.sprites = [4][]drawables.Drawable{}

	/*	uiSprites, err := loader.LoadUISprites(fishSceneUISprites, g.environment, gameLog.GlobalEventHub)
		g.sprites[0] = append(g.sprites[0], uiSprites...)
		g.allSprites = append(g.sprites[0], uiSprites...)*/

	g.state = &FishSceneState{}

	return g
}

func LoadPurchasedSprite(environment *system.Environment, inputName string, hub *tasks.EventHub, tankSize image.Rectangle) *entities.CreatureData {
	fData := entities.SavedFish{
		FishType: inputName,
		Size:     1,
	}
	creature := loader.NewFish(environment, hub, tankSize, fData)
	return creature
}

func (g *FishScene) LoadStuff() {

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
			println("tank Collision map x =", g.collisionMap["tank"].Min.X)
			g.sprites[0] = append(g.sprites[0], loadedFish)
			g.allSprites = append(g.allSprites, loadedFish)
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

	//some of this stuff can be moved to new day config just didnt feel like figuring it out for now
	g.debug = &debug.DebugData{DebugRect: &util.Rect{
		Rectangle: &image.Rectangle{Min: image.Point{}, Max: image.Point{}}},
		DebugParameter: make(map[debug.DebugOption]bool)}

	g.debug.DebugParameter[debug.Print] = false
	g.debug.DebugParameter[debug.Position] = false

	g.debug.DebugRect.Init("tank")
	log.Println("----FishScene OnEnter() called----")

	g.mouseFlags = &input.MouseFlags{HandledClick: false, CursorOccupied: false}
	g.backGroundParams["Cursor"] = []float64{0, 0}

	//No music on the base level as of now
	//g.timers["songTimer"].TurnOn()

	g.returnScene = sceneManagement.FishTank

	tutMngr := tutorial.Manager{}
	tutorial.InitData(&tutMngr, g.gameLog.GlobalEventHub)
	g.tutorialManager = &tutMngr

	if g.gameLog.Day != g.state.lastDayEntered {

		g.configureNewDay()
		log.Println("----FishScene OnEnter() finished----")
		return
	}

	g.configureExistingDay()

	log.Println("----FishScene OnEnter() finished----")
}

func (g *FishScene) configureExistingDay() {
	g.SetNightLight()
	g.gameLog.SongPlayer.Play(soundFX.TropicalHouse)
}

func (g *FishScene) configureNewDay() {
	ev := events.NewDay{NTasks: len(g.gameLog.Tasks), Day: g.gameLog.Day}

	switch g.gameLog.DayType {
	case sceneManagement.Chores:
		ev.Type = "Chores"
	case sceneManagement.Camp:
		ev.Type = "Camp"
	default:
		ev.Type = "Free"
	}

	g.gameLog.GlobalEventHub.Publish(ev)
	g.state.lastDayEntered = g.gameLog.Day
	g.gameLog.Tasks[0].Activate(g.gameLog.GlobalEventHub)
	g.lightingShader = registry.ShaderMap["DayLight"]
	g.backGroundParams["Cursor"] = [2]float64{0, 0}
	g.shaderUpdater = nil

	creatureManager.allFishFed = false
}

func (g *FishScene) SetNightLight() {
	g.lightingShader = registry.ShaderMap["OnePointLighting"]
	tankRect := g.collisionMap["tank"]
	g.lightingShaderParams["LightPoint"] = [2]float64{float64(tankRect.Min.X + tankRect.Dx()/2), ScreenHeight/5 + 20}
	g.lightingShaderParams["LightWidth"] = 120.0
	g.lightingShaderParams["TankRect"] = [4]float64{
		float64(tankRect.Min.X),
		float64(tankRect.Min.Y),
		float64(tankRect.Max.X),
		float64(tankRect.Max.Y)}
}

func (g *FishScene) OnExit() {
	log.Println("----FishScene Exit----")
	g.gameLog.GlobalEventHub.Publish(events.LeavingFishScene{})
	graphics.DeInitAllGraphics()

}

func (g *FishScene) LoadTimers() {
	g.timers = make(map[string]*util.Timer)
	g.timers["pointGeneratedTimer"] = util.NewTimer(0.2)
	g.timers["pointGeneratedTimer"].TurnOn()
	g.timers["songTimer"] = util.NewTimer(15)
	g.timers["sceneTransition"] = util.NewTimer(5.0)
	g.timers["publishNewTask"] = util.NewTimer(0.2)
	g.timers["leaveScene"] = util.NewTimer(10.0)
	g.timers["leaveScene"].TurnOn()
}

func (g *FishScene) IsLoaded() bool {
	return g.loaded
}

func (g *FishScene) Update() (sceneManagement.SceneId, error) {
	g.mouseFlags.HandledClick = false

	if g.shaderUpdater != nil {
		g.lightingShaderParams = g.shaderUpdater(g.lightingShaderParams)
	}

	if g.tankShaderUpdater != nil {
		g.tankShaderParams = g.tankShaderUpdater(g.tankShaderParams)
	}

	g.tutorialManager.Update()

	g.gameLog.SoundPlayer.Update()

	for _, sprite := range g.sprites[0] {
		sprite.Update()
		/*if sprite.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && g.currentSpriteHovered == nil {
			g.currentSpriteHovered = sprite
			g.mouseFlags.CursorOccupied = true
		}*/
	}

	g.ManageLayers()
	//update second layer after managing layers?
	for _, s := range g.sprites[1] {
		s.Update()
	}

	g.ui.Update()

	g.updateTimers()

	g.updateInput()

	//g.tutorialManager.CharUpdate()

	if g.debug.GameMode == debug.Debug {
		if g.debug.DebugParameter[debug.Position] {
			g.positionModeUpdate()
			//ev := events.ButtonClickedEvent{ButtonText: "Mode"}
			//g.gameLog.GlobalEventHub.Publish(ev)
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
	if g.ui.HasFocus() {
		g.mouseFlags.CursorOccupied = true
	} else {
		g.mouseFlags.CursorOccupied = false
	}
	graphics.UpdateGraphics()

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

	//define our draw options
	opts := &ebiten.DrawImageOptions{}
	shaderOpts := &ebiten.DrawRectShaderOptions{}

	//step1 draw room background to primary offscreen
	/*x, y := ebiten.WindowSize()
	if float64(x)/float64(y) != 16.0/9.0 {
		g.images.OffScreen.DrawImage(g.images.LaptopRoomBackground, opts)
	}*/
	g.images.OffScreen.DrawImage(g.images.RoomBackground, opts)

	b := g.images.FishTank.Bounds()
	opts.GeoM.Reset()

	//draw the fishtank depending on lighting conditions with normal map shader

	shaderOpts.Uniforms = g.backGroundParams
	shaderOpts.Images[0] = g.images.FishTank
	shaderOpts.Images[1] = g.images.FishTank_n

	if g.lightingState == NightLight {
		shaderOpts.Images[0] = g.images.FishTankNight
	}

	//apply normal map if its initiated
	if g.offScreenShader != nil {
		shaderOpts.GeoM.Translate(float64(g.tankSize.Min.X), float64(g.tankSize.Min.Y))
		g.images.OffScreen.DrawRectShader(b.Dx(), b.Dy(), g.offScreenShader, shaderOpts)
	} else {
		opts.GeoM.Translate(float64(g.tankSize.Min.X), float64(g.tankSize.Min.Y))
		g.images.OffScreen.DrawImage(g.images.FishTank, opts)
	}

	//draw our tank props
	if g.testProp != nil {
		g.testProp.Draw(g.images.OffScreen)
	}
	//props.DrawProps(g.propQueue, g.images.OffScreen)

	//draw our non-highlighted("below global lighting shader") sprites(these will handle their applied shaders internally)
	for _, s := range g.sprites[0] {
		s.Draw(g.images.OffScreen)
	}

	opts.GeoM.Reset()
	fishTankFrontLayerDy := g.images.FishTankFrontLayerNoLightSmaller.Bounds().Dy()
	fishTankHeightDy := g.images.FishTank.Bounds().Dy()

	// draw the front tank of the fishlayer over everything
	y := fishTankFrontLayerDy - fishTankHeightDy
	y = g.tankSize.Min.Y - y
	opts.GeoM.Translate(float64(g.tankSize.Min.X), float64(y))
	g.images.OffScreen.DrawImage(g.images.FishTankFrontLayerNoLightSmaller, opts)

	//apply water affect if it's on
	if g.tankShader != nil {
		shaderOpts.GeoM.Reset()
		shaderOpts.Uniforms = g.tankShaderParams
		shaderOpts.Images[0] = g.images.OffScreen
		shaderOpts.Images[1] = nil
		g.images.OffScreen2.DrawRectShader(ScreenWidth, ScreenHeight, g.tankShader, shaderOpts)
	}

	//stupid desk thing
	//g.images.OffScreen2.DrawImage(g.images.FrontLayer, opts)

}

func (g *FishScene) Draw(screen *ebiten.Image) {

	//call offscreen function to render everything we want to apply global shaders to
	g.DrawOffScreen()

	dOpts := &ebiten.DrawImageOptions{}
	ShaderOpts := &ebiten.DrawRectShaderOptions{}

	//apply global lighting shader
	ShaderOpts.Images[0] = g.images.OffScreen2
	ShaderOpts.Uniforms = g.lightingShaderParams

	if g.lightingShader != nil {
		g.smallerResolution.DrawRectShader(ScreenWidth, ScreenHeight, g.lightingShader, ShaderOpts)
	} else {
		g.smallerResolution.DrawImage(g.images.OffScreen, dOpts)
	}

	for _, sp := range g.sprites[1] {
		//draw highlighted sprites (won't be affected by lighting shaders)
		sp.Draw(g.smallerResolution)
	}

	//debug rect needs to be scaled to base resolution

	//scale our draw opts to resolution

	if g.debug.GameMode == debug.Debug {

		g.debug.DebugRect.Draw(g.smallerResolution)

		for key, rect := range g.collisionMap {
			util.StrokeRectFromImageRect(rect, g.smallerResolution)
			cs := ebiten.ColorScale{}
			cs.SetB(0.0)
			cs.SetG(0.0)
			cs.SetA(1.0)
			face := registry.FontMap["nk57"]
			tOpts := &text.DrawOptions{}
			tOpts.ColorScale = cs
			tOpts.GeoM.Translate(float64(rect.Min.X+rect.Dx()/2*g.resolutionScaling), float64(rect.Min.Y*g.resolutionScaling))
			text.Draw(screen, key, face, tOpts)
		}

		cursor, ok := g.backGroundParams["Cursor"].([2]float64)
		if ok {
			vector.StrokeCircle(g.smallerResolution, float32(cursor[0]), float32(cursor[1]), 10, 5, colornames.Orange, false)
		}

		lightPoint, ok := g.lightingShaderParams["LightPoint"].([2]float64)
		if ok {
			vector.StrokeCircle(g.smallerResolution, float32(lightPoint[0]), float32(lightPoint[1]), 10, 5, colornames.Greenyellow, false)
		}

	}

	graphics.DrawScaledGraphics(g.smallerResolution)

	dOpts.GeoM.Reset()
	x, y := ebiten.WindowSize()
	if x%ScreenWidth > 0 {
		yTranslate := float64(y-ScreenHeight) / 8.0
		dOpts.GeoM.Translate(0, yTranslate)
	}

	dOpts.GeoM.Scale(registry.Config.ResolutionScalingF, registry.Config.ResolutionScalingF)
	screen.DrawImage(g.smallerResolution, dOpts)

	//custom debug text
	DebugText(g.debug.DebugText, screen)

	//tps monitor
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))

	//draw graphics
	graphics.DrawUnScaledGraphics(screen)

	//draw ui last
	g.ui.Draw(screen)

}

func (g *FishScene) positionModeUpdate() {
	if ebiten.IsKeyPressed(ebiten.KeyM) {
		g.debug.DebugRect.Init("WaterEffect")
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

		if key == "songTimer" && state == util.Done {
			timer.TurnOff()

		}

		if key == "sceneTransition" && state == util.Done {
			timer.TurnOff()
			g.returnScene = sceneManagement.TransitionScene
			g.lightingState = Night
			g.gameLog.GlobalEventHub.Publish(events.DayOverTransitionComplete{})
		}

		if key == "publishNewTask" && state == util.Done {
			if len(g.gameLog.Tasks) > g.currentTask {
				timer.TurnOff()
				g.gameLog.Tasks[g.currentTask].Activate(g.gameLog.GlobalEventHub)
				if g.gameLog.Tasks[g.currentTask].Type == tasks.FishFed {
					if creatureManager.allFishFed {
						g.gameLog.GlobalEventHub.Publish(entities.AllFishFed{})
					}
				}

			} else {
				timer.TurnOff()
				g.currentTask = 0
			}
		}
		if key == "leaveScene" && state == util.Done {
			timer.TurnOff()
			if g.gameLog.DayType == sceneManagement.Camp {
				eff := loader.LoadStaticEffect("timeForCamp", 100, 85)
				g.sprites[1] = append(g.sprites[1], eff)
			}

		}

	}
}

func (g *FishScene) updateInput() {

	//OLD COMMENTS
	//function for handling ebiten input directly in game mode mainly for convenience
	//or avoiding the event system (latency?)
	//not necessarily core game functions
	//OLD COMMENTS

	//central hub for game input

	//focus on correct sprite at click
	//or equivalent of mouse button press on controller

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		xAtClick, yAtClick := ebiten.CursorPosition()
		closestToCursor, distance := drawables.ClosestDrawableToCursor(xAtClick, yAtClick, g.allSprites)
		//focus on it if its less then 10 pixels away from click (arbitrary, smarter sorting coming)
		if distance < 10 {
			//if a sprite is focused they can do focused stuff in update rather then have input detection in each sprite
			closestToCursor.Focus()
		}
	}
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
			ev := events.ButtonClickedEvent{ButtonText: "butt"}
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

func (g *FishScene) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth > 0 {
		return outsideWidth, outsideHeight
	}
	return ScreenWidth, ScreenHeight
}

func (g *FishScene) subs(colMap map[string]image.Rectangle) {

	g.uiSubs()
	g.soundSubs()
	//g.creatureSubs(colMap)

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
			g.lightingShaderParams["LightIntensity"] = 1.0
			g.lightingShaderParams["Counter"] = 0.0
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
		creature, ok := draw.(*entities.CreatureData)
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

		saveTask["Name"] = task.Type
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

func (g *FishScene) saveUISpritePositions() {

	/*spMap := make(map[string]drawables.SavePositionData)
	//THIS IS AWFUL BECAUSE I MADE TOO MANY UI SPRITE TYPES DEAL WITH SOME DAY
	//check for ui sprites
	for _, layer := range g.sprites {
		for _, sprite := range layer {
			uiSprite, ok := sprite.(*entities.UiSprite)
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
	}*/
}
