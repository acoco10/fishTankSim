package sceneManagement

import (
	"fishTankWebGame/game/gameEntities"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameLog struct {
	PlayerLoginId string
	PreviousScene SceneId
	Mode          GameMode
}

type SceneId uint

const (
	FishTank SceneId = iota
	StartingScreenNewSave
	StartingScreenExistingSave
)

type GameMode uint

const (
	Standard GameMode = iota
)

type Scene interface {
	Update() (SceneId, error)
	Draw(screen *ebiten.Image)
	FirstLoad(state gameEntities.SaveGameState)
	OnEnter()
	OnExit()
	IsLoaded() bool
}
