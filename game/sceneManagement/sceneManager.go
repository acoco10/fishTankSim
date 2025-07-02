package sceneManagement

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

type GameLog struct {
	PlayerLoginId  string
	PreviousScene  SceneId
	Save           *entities.SaveGameState
	GlobalEventHub *tasks.EventHub
	SongPlayer     *soundFX.SoundPlayer
	SoundPlayer    *soundFX.SoundPlayer
	Day            int
	Tasks          []*tasks.Task
}

func NewGameLog(state entities.SaveGameState) *GameLog {
	g := GameLog{}
	g.Save = &entities.SaveGameState{}
	g.Save = &state
	eHub := tasks.NewEventHub()
	g.GlobalEventHub = eHub
	songP, err := soundFX.NewSoundPlayer()
	if err != nil {
		log.Fatal(err)
	}

	soundP, err := soundFX.NewSoundPlayer()
	if err != nil {
		log.Fatal(err)
	}

	g.SongPlayer = songP
	g.SoundPlayer = soundP
	g.Day = 1
	var tasks []*tasks.Task
	g.Tasks = tasks

	return &g
}

type SceneId uint

const (
	FishTank SceneId = iota
	StartScene
	TransitionScene
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
