package sceneManagement

import (
	"github.com/acoco10/fishTankWebGame/game/eventSytem"
	"github.com/acoco10/fishTankWebGame/game/gameEntities"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

type GameLog struct {
	PlayerLoginId  string
	PreviousScene  SceneId
	Save           *entities.SaveGameState
	GlobalEventHub *events.EventHub
	SongPlayer     *soundFX.SongPlayer
	SoundPlayer    *soundFX.SongPlayer
	Tasks          []*entities.Task
}

func NewGameLog(state entities.SaveGameState) *GameLog {
	g := GameLog{}
	g.Save = &entities.SaveGameState{Fish: []entities.SavedFish{}}
	g.Save = &state
	eHub := events.NewEventHub()
	g.GlobalEventHub = eHub
	songP, err := soundFX.NewSongPlayer()
	if err != nil {
		log.Fatal(err)
	}

	soundP, err := soundFX.NewSongPlayer()
	if err != nil {
		log.Fatal(err)
	}

	g.SongPlayer = songP
	g.SoundPlayer = soundP

	var tasks []*entities.Task
	g.Tasks = tasks

	return &g
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
