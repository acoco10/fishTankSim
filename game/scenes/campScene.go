package scenes

import (
	"github.com/acoco10/fishTankWebGame/game/daySystem"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

type CampScene struct {
	isLoaded          bool
	nextSceneTrigger  *util.Timer
	dotTimer          *util.Timer
	msgTime           *util.Timer
	gameLog           *sceneManagement.GameLog
	backGround        *ebiten.Image
	smallerResolution *ebiten.Image
	loadingMsg        string
	graphid           int
}

func LoadCampScene(gameLog *sceneManagement.GameLog) *CampScene {
	s := CampScene{}
	s.gameLog = gameLog
	background, err := util.LoadImageAssetAsEbitenImage("campScenes/campRockEvent")
	if err != nil {
		log.Fatal(err)
	}
	s.backGround = background
	return &s
}

func (s *CampScene) Update() (sceneManagement.SceneId, error) {

	if s.nextSceneTrigger.TimerState == util.Done {
		return sceneManagement.FishTank, nil
	}
	if s.msgTime.TimerState == util.Done {
		graphics.DeInitGraphicId(s.graphid)
		s.graphid = graphics.NewFadeInTextGraphic("You found a cool rock", 200.0, 100.0)
		s.msgTime.TurnOff()
	}
	s.nextSceneTrigger.Update()
	s.msgTime.Update()

	graphics.UpdateGraphics()

	return sceneManagement.CampScene, nil

}

func (s *CampScene) Draw(screen *ebiten.Image) {

	drawOpts := &ebiten.DrawImageOptions{}

	s.smallerResolution.DrawImage(s.backGround, drawOpts)

	drawOpts.GeoM.Reset()
	drawOpts.GeoM.Scale(registry.Config.ResolutionScalingF, registry.Config.ResolutionScalingF)
	drawOpts.GeoM.Translate(0, registry.Config.YOffsetF)
	screen.DrawImage(s.smallerResolution, drawOpts)
	graphics.DrawUnScaledGraphics(screen)

}

func (s *CampScene) FirstLoad() {
	s.isLoaded = true
	s.dotTimer = util.NewTimer(0.5)
	s.nextSceneTrigger = util.NewTimer(5)
	s.msgTime = util.NewTimer(2)
	s.smallerResolution = ebiten.NewImage(registry.Config.ScreenWidth, registry.Config.ScreenHeight)
}

func (s *CampScene) OnEnter() {
	println("------Entering Camp Scene---------")
	s.nextSceneTrigger.TurnOn()
	s.dotTimer.TurnOn()
	s.msgTime.TurnOn()
	s.graphid = graphics.NewFadeInTextGraphic("You had a great day at camp!", 200.0, 100.0)
	daySystem.LoadDaysTasks(s.gameLog)

}

func (s *CampScene) OnExit() {
	graphics.DeInitAllGraphics()
	s.nextSceneTrigger.TurnOff()
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
