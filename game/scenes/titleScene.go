//go:build old

package scenes

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

type TitleScene struct {
	isLoaded bool
	gameLog  *sceneManagement.GameLog
	hub      *tasks.EventHub
	ui       *ebitenui.UI
}

func NewTitleScene(gameLog *sceneManagement.GameLog) *CampScene {

	return &s
}

func (s *TitleScene) Update() (sceneManagement.SceneId, error) {
	s.ui.Update()
	return sceneManagement.TitleScene, nil

}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	s.ui.Draw(screen)
}

func (s *TitleScene) FirstLoad() {

}

func (s *TitleScene) OnEnter() {
	s.hub = tasks.NewEventHub()
	buttonText := []string{"Laptop", "Monitor"}
	menuUI, err := ui.LoadNextOptionsMenuUI("Choose Resolution", buttonText, s.hub)
	if err != nil {
		return
	}
	s.ui.AddWindow(menuUI)

	s.subs()
}

func (s *TitleScene) OnExit() {

}

func (s *TitleScene) IsLoaded() bool {
	return s.isLoaded
}

func (s *TitleScene) subs() {
	s.hub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		switch ev.ButtonText {
		case "Laptop":
			return
		case "Choose Resolution: Monitor":
			println("setting window size to monitor")
			ConfigResolution()
			s.ui.Container.RequestRelayout()
		}
	})
}
