package game

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/gameEntities"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
)

type StartScene struct {
	ui               *ui.StartMenu
	isLoaded         bool
	fishOptions      []*gameEntities.AnimatedSprite
	nextSceneTrigger *gameEntities.Timer
	*sceneManagement.GameLog
}

func NewStartScene(gameLog *sceneManagement.GameLog) *StartScene {
	s := StartScene{}

	s.GameLog = gameLog

	sUI, err := ui.LoadStartMenu(gameLog.GlobalEventHub)

	if err != nil {
		log.Fatal(fmt.Errorf("error initiating start menu: %s", err))
	}

	s.ui = sUI
	s.subs(gameLog)
	timer := gameEntities.NewTimer(1)
	s.nextSceneTrigger = timer
	return &s
}

func (s *StartScene) Update() (sceneManagement.SceneId, error) {
	s.ui.UI.Update()

	for _, fish := range s.ui.SelectSpritesToDraw {
		fish.Update()
	}

	s.nextSceneTrigger.Update()
	if s.nextSceneTrigger.TimerState == gameEntities.Done {
		s.ui.UI.ClearFocus()
		return sceneManagement.FishTank, nil
	}

	return sceneManagement.StartScene, nil
}

func (s *StartScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 120, G: 170, B: 210, A: 255})
	opts := &ebiten.DrawImageOptions{}
	sopts := &ebiten.DrawRectShaderOptions{}

	s.ui.Draw(screen)

	for _, fish := range s.ui.SelectSpritesToDraw {
		opts.GeoM.Scale(4, 4)
		sopts.GeoM.Scale(4, 4)

		opts.GeoM.Translate(float64(fish.X), float64(fish.Y))
		sopts.GeoM.Translate(float64(fish.X), float64(fish.Y))

		fish.Draw(screen, opts, sopts)
		opts.GeoM.Reset()
		sopts.GeoM.Reset()
	}
}

func (s *StartScene) FirstLoad(gameLog *sceneManagement.GameLog) {
	s.isLoaded = true
	s.GameLog.SongPlayer.Play(soundFX.BestAdventureEver)

}

func (s *StartScene) OnEnter(gameLog *sceneManagement.GameLog) {

}

func (s *StartScene) OnExit() {
	s.GameLog.SongPlayer.Pause()
}

func (s *StartScene) IsLoaded() bool {
	return s.isLoaded
}

func (s *StartScene) subs(gameLog *sceneManagement.GameLog) {
	gameLog.GlobalEventHub.Subscribe(gameEntities.ButtonEvent{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.ButtonEvent)
		if ev.EType == "cursor exited" {
			if ev.ButtonText != "Select" {
				if len(s.ui.SelectSpritesToDraw) > 1 {
					s.ui.SelectSpriteOptions[ev.ButtonText].UnLoadShader()
				}
			}
		}

	})

	gameLog.GlobalEventHub.Subscribe(gameEntities.ButtonClickedEvent{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.ButtonClickedEvent)
		ols := gameEntities.LoadOutlineShader()
		switch ev.ButtonText {
		case "Common Molly":
			gameLog.SoundPlayer.Play(soundFX.SelectSound2)
			saveFish := gameEntities.SavedFish{FishType: "mollyFish", Progress: 0, Size: 1}
			gameLog.Save.Fish = append(gameLog.Save.Fish, saveFish)
			s.ui.TextInputContainer.GetWidget().Disabled = false
			s.ui.TextInputContainer.GetWidget().Visibility = widget.Visibility_Show
			s.ui.TextInput.Focus(true)
			s.ui.SelectSpriteOptions[ev.ButtonText].LoadShader(ols)
		case "Goldfish":
			gameLog.SoundPlayer.Play(soundFX.SelectSound2)
			saveFish := gameEntities.SavedFish{FishType: "fish", Progress: 0, Size: 1}
			gameLog.Save.Fish = append(gameLog.Save.Fish, saveFish)
			s.ui.TextInputContainer.GetWidget().Disabled = false
			s.ui.TextInputContainer.GetWidget().Visibility = widget.Visibility_Show
			s.ui.TextInput.Focus(true)
			s.ui.SelectSpriteOptions[ev.ButtonText].LoadShader(ols)
		}
	})

	gameLog.GlobalEventHub.Subscribe(gameEntities.SendData{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.SendData)
		if ev.DataFor == "Name Input" {
			gameLog.Save.Fish[0].Name = ev.Data
			s.ui.TextInputButton.Press()
			s.nextSceneTrigger.TurnOn()
			s.ui.TextInput.Focus(false)
			s.GameLog.SoundPlayer.Play(soundFX.SelectSound)
		}
	})

	gameLog.GlobalEventHub.Subscribe(gameEntities.ButtonEvent{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.ButtonEvent)
		if ev.EType == "cursor entered" {
			ols := gameEntities.LoadOutlineShader()
			println(ev.ButtonText)
			if ev.ButtonText != "Select" {
				println(ev.ButtonText)
				s.ui.SelectSpriteOptions[ev.ButtonText].LoadShader(ols)
			}
		}

	})
}
