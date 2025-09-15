package scenes

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"log"
)

type UserType int8

const (
	NewUser UserType = iota
	ExistingUser
)

type State int8

const (
	Ready State = iota
	LoadingNextScene
)

type Game struct {
	sceneMap       map[sceneManagement.SceneId]sceneManagement.Scene
	activeSceneId  sceneManagement.SceneId
	pendingSceneId sceneManagement.SceneId
	gameLog        *sceneManagement.GameLog
	state          State
}

const (
	ScreenWidth  = 960
	ScreenHeight = 540
)

func NewGame(log *sceneManagement.GameLog, userType UserType) *Game {

	switch userType {
	case NewUser:

		activeSceneId := sceneManagement.StartScene
		ebiten.SetWindowSize(1920, 1080)
		ebiten.SetFullscreen(false)
		ConfigResolution()

		sceneMap := map[sceneManagement.SceneId]sceneManagement.Scene{
			sceneManagement.StartScene:      NewStartScene(log),
			sceneManagement.FishTank:        NewFishScene2(log),
			sceneManagement.TransitionScene: LoadTransitionScene(log),
		}

		game := &Game{
			sceneMap:      sceneMap,
			activeSceneId: activeSceneId,
			gameLog:       log,
		}

		sceneMap[activeSceneId].FirstLoad()
		sceneMap[activeSceneId].OnEnter()

		return game

	case ExistingUser:

		println("existing user save = ", log.Save.Fish[0].FishType)

		activeSceneId := sceneManagement.FishTank

		sceneMap := map[sceneManagement.SceneId]sceneManagement.Scene{
			sceneManagement.FishTank: NewFishScene2(log),
		}

		game := &Game{
			sceneMap:      sceneMap,
			activeSceneId: activeSceneId,
			gameLog:       log,
		}

		sceneMap[activeSceneId].FirstLoad()

		sceneMap[activeSceneId].OnEnter()

		return game
	}
	return nil
}

func ScaleScreenToResolution() float64 {

	screenWidth := float64(ScreenWidth)
	screenHeight := float64(ScreenHeight)

	// Calculate scaling factors for both dimensions
	scaleX := float64(1920) / screenWidth
	scaleY := float64(1080) / screenHeight

	// Use the smaller scale to ensure content fits within the window
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
		println("scaling resolution to accommodate height")
	}

	// Safety check for division by zero
	if scale <= 0 {
		log.Fatalf("Invalid scaling: screenWidth=%f, screenHeight=%f, windowWidth=%d, windowHeight=%d",
			screenWidth, screenHeight, screenWidth, screenHeight)
	}

	println("resolution scalar =", scale)
	return scale
}

func (g *Game) Update() error {

	nextSceneId, err := g.sceneMap[g.activeSceneId].Update()
	if err != nil {
		return err
	}

	if nextSceneId != g.activeSceneId {
		if nextSceneId == sceneManagement.Reset {
			g.gameLog.GlobalEventHub = tasks.NewEventHub()
			g.sceneMap[g.activeSceneId].OnExit()
			g.sceneMap[g.activeSceneId] = g.sceneMap[sceneManagement.Reset]
			g.sceneMap[g.activeSceneId].FirstLoad()
			g.sceneMap[g.activeSceneId].OnEnter()

			switch g.activeSceneId {
			}

			return nil
		}
		nextScene := g.sceneMap[nextSceneId]
		// if not loaded? then load in
		if !nextScene.IsLoaded() && g.state == Ready {
			nextScene.FirstLoad()
			return nil
		}

		g.state = LoadingNextScene

		if nextScene.IsLoaded() {
			g.sceneMap[g.activeSceneId].OnExit()
			g.state = Ready
			g.gameLog.PreviousScene = g.activeSceneId
			graphics.DeInitAllGraphics()
			g.activeSceneId = nextSceneId
			g.sceneMap[g.activeSceneId].OnEnter()

		}

	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneMap[g.activeSceneId].Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func NewTestScene(TestSceneID sceneManagement.SceneId, log *sceneManagement.GameLog) *Game {
	registry.Config.Set(registry.Debug, true)
	ebiten.SetWindowSize(1920, 1080)
	ConfigResolution()
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	activeSceneId := TestSceneID

	sceneMap := make(map[sceneManagement.SceneId]sceneManagement.Scene)

	switch TestSceneID {
	case sceneManagement.MowingMiniGameScene:
	case sceneManagement.FishTank:
		log.Save.Fish = append(log.Save.Fish, entities.SavedFish{FishType: string(entities.MollyFish), Size: 1})
		sceneMap[TestSceneID] = NewFishScene2(log)
		//sceneMap[sceneManagement.Reset] = NewFishScene2(log)
	}

	sceneMap[TestSceneID].FirstLoad()

	sceneMap[TestSceneID].OnEnter()

	game := &Game{
		sceneMap:      sceneMap,
		activeSceneId: activeSceneId,
		gameLog:       log,
	}

	return game
}

func ConfigResolution() {

	registry.Config.Set(registry.ScreenWidth, ScreenWidth)
	registry.Config.Set(registry.ScreenHeight, ScreenHeight)
	registry.Config.Set(registry.ResolutionWidth, 1920)
	registry.Config.Set(registry.ResolutionHeight, 1080)
	scaling := ScaleScreenToResolution()
	println("res scaling =", scaling)
	//this one needs to be set last to get the correct y/x offset
	registry.Config.Set(registry.ResolutionScaling, scaling)
	ConfigZoom(2.0, image.Point{X: -300, Y: -300})
}

func ConfigZoom(zoomFactor float64, zoomOffset image.Point) {
	registry.Config.Set(registry.ZoomFactor, zoomFactor)
	registry.Config.Set(registry.ZoomOffset, zoomOffset)
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
