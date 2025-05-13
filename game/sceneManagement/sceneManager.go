package sceneManagement

import (
	"github.com/acoco10/fishTankWebGame/game/gameEntities"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameLog struct {
	PlayerLoginId  string
	PreviousScene  SceneId
	Save           *gameEntities.SaveGameState
	GlobalEventHub *gameEntities.EventHub
	SongPlayer     *soundFX.SongPlayer
	SoundPlayer    *soundFX.SongPlayer
}

type SceneId uint

const (
	FishTank SceneId = iota
	StartScene
)

type GameMode uint

const (
	Standard GameMode = iota
)

type Scene interface {
	Update() (SceneId, error)
	Draw(screen *ebiten.Image)
	FirstLoad(gameLog *GameLog)
	OnEnter(gameLog *GameLog)
	OnExit()
	IsLoaded() bool
}
