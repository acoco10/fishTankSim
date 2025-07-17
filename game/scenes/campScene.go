package scenes

import (
	"github.com/acoco10/fishTankWebGame/game/daySystem"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
	"log"
)

type CampScene struct {
	isLoaded         bool
	nextSceneTrigger *entities.Timer
	dotTimer         *entities.Timer
	gameLog          *sceneManagement.GameLog
	loadingMsg       string
}

func LoadCampScene(gameLog *sceneManagement.GameLog) *CampScene {
	s := CampScene{}
	s.gameLog = gameLog
	s.loadingMsg = "You had a great Day at Camp"
	return &s
}

func (s *CampScene) Update() (sceneManagement.SceneId, error) {

	if s.nextSceneTrigger.TimerState == entities.Done {
		return sceneManagement.FishTank, nil
	}
	s.nextSceneTrigger.Update()
	return sceneManagement.CampScene, nil
}

func (s *CampScene) Draw(screen *ebiten.Image) {

	face, err := util.LoadFont(24, "nk57")
	if err != nil {
		log.Fatal(err, "Cant load font in transition scene")
	}
	screen.Fill(colornames.Papayawhip)

	dopts := &text.DrawOptions{}

	dopts.ColorScale.Scale(0, 1, 0.5, 1)
	dopts.GeoM.Translate(ScreenWidth/2, ScreenHeight/2)

	text.Draw(screen, s.loadingMsg, face, dopts)
}

func (s *CampScene) FirstLoad() {
	s.isLoaded = true
	s.dotTimer = entities.NewTimer(0.5)
	s.nextSceneTrigger = entities.NewTimer(2)
}

func (s *CampScene) OnEnter() {

	s.nextSceneTrigger.TurnOn()
	s.dotTimer.TurnOn()

	daySystem.LoadDaysTasks(s.gameLog)

}

func (s *CampScene) OnExit() {
}

func (s *CampScene) IsLoaded() bool {
	return s.isLoaded
}

func (s *CampScene) subs(gameLog *sceneManagement.GameLog) {
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
