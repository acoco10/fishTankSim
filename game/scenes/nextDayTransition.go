package scenes

import (
	"github.com/acoco10/fishTankWebGame/game/daySystem"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
)

type TransitionScene struct {
	isLoaded         bool
	nextSceneTrigger *entities.Timer
	dotTimer         *entities.Timer
	gameLog          *sceneManagement.GameLog
	dots             int
	loadingMsg       string
}

func LoadTransitionScene(gameLog *sceneManagement.GameLog) *TransitionScene {
	s := TransitionScene{}
	s.dots = 1
	s.loadingMsg = "Next Day Loading"
	return &s
}

func (s *TransitionScene) Update() (sceneManagement.SceneId, error) {

	if s.nextSceneTrigger.TimerState == entities.Done {
		return sceneManagement.FishTank, nil
	}
	s.nextSceneTrigger.Update()
	s.dotTimer.Update()
	return sceneManagement.TransitionScene, nil
}

func (s *TransitionScene) Draw(screen *ebiten.Image) {

	face, err := util.LoadFont(24, "nk57")
	if err != nil {
		log.Fatal(err, "Cant load font in transition scene")
	}
	screen.Fill(color.RGBA{R: 120, G: 170, B: 210, A: 255})

	dopts := &text.DrawOptions{}

	dopts.ColorScale.Scale(1, 1, 1, 1)
	dopts.GeoM.Translate(ScreenWidth/2, ScreenHeight/2)

	s.dots++
	if s.dots%17 == 0 {
		s.loadingMsg += "."
		if len(s.loadingMsg) > len("Next Day Loading")+4 {
			s.loadingMsg = "Next Day Loading"
		}
	}
	text.Draw(screen, s.loadingMsg, face, dopts)
}

func (s *TransitionScene) FirstLoad(gameLog *sceneManagement.GameLog) {

}

func (s *TransitionScene) OnEnter(gameLog *sceneManagement.GameLog) {
	s.dotTimer = entities.NewTimer(0.5)
	s.nextSceneTrigger = entities.NewTimer(2)

	s.nextSceneTrigger.TurnOn()
	s.dotTimer.TurnOn()

	s.gameLog = gameLog
	daySystem.LoadDaysTasks(gameLog)
	s.isLoaded = true
}

func (s *TransitionScene) OnExit() {
	nTasks := len(s.gameLog.Tasks)
	ev := events.NewDay{
		NTasks: nTasks,
	}
	s.gameLog.GlobalEventHub.Publish(ev)
}

func (s *TransitionScene) IsLoaded() bool {
	return s.isLoaded
}

func (s *TransitionScene) subs(gameLog *sceneManagement.GameLog) {
	/*	gameLog.GlobalEventHub.Subscribe(events.ButtonEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonEvent)
		if ev.EType == "cursor entered" {
			ols := shaders.LoadRotatingHighlightShader()
			println(ev.ButtonText)
			if ev.ButtonText != "Select" {
				println(ev.ButtonText)
				s.ui.SelectSpriteOptions[ev.ButtonText].LoadShader(ols)
			}
		}

	})*/
}
