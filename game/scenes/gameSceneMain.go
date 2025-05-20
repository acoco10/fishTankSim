package scenes

import (
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/hajimehoshi/ebiten/v2"
)

type UserType int8

const (
	NewUser UserType = iota
	ExistingUser
)

type Game struct {
	sceneMap      map[sceneManagement.SceneId]sceneManagement.Scene
	activeSceneId sceneManagement.SceneId
	gameLog       *sceneManagement.GameLog
}

func NewGame(log *sceneManagement.GameLog, userType UserType) *Game {

	switch userType {
	case NewUser:
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

		activeSceneId := sceneManagement.StartScene

		ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

		sceneMap := map[sceneManagement.SceneId]sceneManagement.Scene{
			sceneManagement.StartScene: NewStartScene(log),
			sceneManagement.FishTank:   NewFishScene(log),
		}

		game := &Game{
			sceneMap,
			activeSceneId,
			log,
		}

		sceneMap[activeSceneId].FirstLoad(game.gameLog)

		return game
	case ExistingUser:

		println("existing user save = ", log.Save.Fish[0].FishType)
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

		activeSceneId := sceneManagement.FishTank

		ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

		sceneMap := map[sceneManagement.SceneId]sceneManagement.Scene{
			sceneManagement.FishTank: NewFishScene(log),
		}

		game := &Game{
			sceneMap,
			activeSceneId,
			log,
		}

		sceneMap[activeSceneId].FirstLoad(game.gameLog)

		sceneMap[activeSceneId].OnEnter(log)

		return game
	}

	return nil
}

func (g *Game) Update() error {
	nextSceneId, err := g.sceneMap[g.activeSceneId].Update()
	if err != nil {
		return err
	}
	// switched scenes
	if nextSceneId != g.activeSceneId {
		g.gameLog.PreviousScene = g.activeSceneId
		g.sceneMap[g.activeSceneId].OnExit()
		nextScene := g.sceneMap[nextSceneId]
		// if not loaded? then load in
		if !nextScene.IsLoaded() {
			nextScene.FirstLoad(g.gameLog)
		}
		nextScene.OnEnter(g.gameLog)
	}
	g.activeSceneId = nextSceneId
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneMap[g.activeSceneId].Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	x, y := ebiten.WindowSize()
	if x > 0 && y > 0 {
		return ebiten.WindowSize()
	}
	return 940, 593
}
