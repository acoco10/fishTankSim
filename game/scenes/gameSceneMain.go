package scenes

import (
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/hajimehoshi/ebiten/v2"
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
	ScreenHeight = 720
)

func NewGame(log *sceneManagement.GameLog, userType UserType) *Game {

	switch userType {
	case NewUser:
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

		activeSceneId := sceneManagement.StartScene

		ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

		sceneMap := map[sceneManagement.SceneId]sceneManagement.Scene{
			sceneManagement.StartScene:          NewStartScene(log),
			sceneManagement.FishTank:            NewFishScene(log),
			sceneManagement.TransitionScene:     LoadTransitionScene(log),
			sceneManagement.MowingMiniGameScene: NewMowingScene(log),
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
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		ebiten.SetFullscreen(true)
		activeSceneId := sceneManagement.FishTank

		sceneMap := map[sceneManagement.SceneId]sceneManagement.Scene{
			sceneManagement.FishTank: NewFishScene(log),
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

func (g *Game) Update() error {

	nextSceneId, err := g.sceneMap[g.activeSceneId].Update()
	if err != nil {
		return err
	}

	if nextSceneId != g.activeSceneId {

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
	return ScreenWidth, ScreenHeight
}

func NewMiniGameTest(miniGameName string, log *sceneManagement.GameLog) *Game {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	activeSceneId := sceneManagement.MowingMiniGameScene

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

	sceneMap := make(map[sceneManagement.SceneId]sceneManagement.Scene)

	switch miniGameName {
	case "MowingGame":
		sceneMap[sceneManagement.MowingMiniGameScene] = NewMowingScene(log)

		sceneMap[sceneManagement.MowingMiniGameScene].OnEnter()
	}

	game := &Game{
		sceneMap:      sceneMap,
		activeSceneId: activeSceneId,
		gameLog:       log,
	}

	return game
}
