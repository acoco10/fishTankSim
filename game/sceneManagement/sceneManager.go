package sceneManagement

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameLog struct {
	TaskManager        *tasks.TaskManager
	PlayerLoginId      string
	PreviousScene      SceneId
	Save               *entities.SaveGameState
	GlobalEventHub     *tasks.EventHub
	StartSceneEventHub *tasks.EventHub
	MowerSceneEventHub *tasks.EventHub
	SongPlayer         *soundFX.SoundPlayer
	SoundPlayer        *soundFX.SoundPlayer
	Day                int
	DayType            DayType
}

type DayType uint8

const (
	Free DayType = iota
	Chores
	Camp
)

func NewGameLog(state entities.SaveGameState, flag string) *GameLog {

	g := GameLog{}
	g.Save = &state
	eHub := tasks.NewEventHub()
	g.GlobalEventHub = eHub

	g.TaskManager = &tasks.TaskManager{EventHub: eHub}
	g.TaskManager.Subscribe()

	if flag == "w" {
		soundFX.LoadOggs()
		oggPlay := &soundFX.SoundPlayer{}
		oggPlay.LoadPlayer("ogg")
		g.SoundPlayer = oggPlay

		oggPlayMusic := &soundFX.SoundPlayer{}
		oggPlayMusic.LoadPlayer("oggMusic")
		g.SongPlayer = oggPlayMusic
	} else {
		soundFX.LoadSounds()
		songP := &soundFX.SoundPlayer{}
		soundP := &soundFX.SoundPlayer{}

		songP.LoadPlayer("music")
		soundP.LoadPlayer("sound")

		g.SongPlayer = songP
		g.SoundPlayer = soundP
	}

	g.Day = 1
	return &g
}

type SceneId uint

const (
	FishTank SceneId = iota
	StartScene
	TransitionScene
	MowingMiniGameScene
	CampScene
	FishSceneDev
	TitleScene
	Reset
)

type GameMode uint

const (
	Standard GameMode = iota
)

type Scene interface {
	Update() (SceneId, error)
	Draw(screen *ebiten.Image)
	FirstLoad()
	OnEnter()
	OnExit()
	IsLoaded() bool
}
