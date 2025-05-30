package scenes

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
)

type StartScene struct {
	ui               *ui.StartMenu
	isLoaded         bool
	fishOptions      []*sprite.AnimatedSprite
	nextSceneTrigger *entities.Timer
	gameLog          *sceneManagement.GameLog
}

func NewStartScene(gameLog *sceneManagement.GameLog) *StartScene {
	s := StartScene{}

	s.gameLog = gameLog

	sUI, err := ui.LoadStartMenu(gameLog.GlobalEventHub)

	if err != nil {
		log.Fatal(fmt.Errorf("error initiating start menu: %s", err))
	}

	taskCondition := func(e events.Event) bool {
		ev, ok := e.(entities.CreatureReachedPoint)
		return ok && ev.Point.PType == geometry.Food
	}

	gameTask := events.NewTask(entities.CreatureReachedPoint{}, "1. Feed your fish", taskCondition)
	gameTask.Subscribe(gameLog.GlobalEventHub)

	taskCondition2 := func(e events.Event) bool {
		ev, ok := e.(entities.SendData)
		return ok && ev.DataFor == "statsMenu"
	}

	gameTask2 := events.NewTask(entities.SendData{}, "2. Click your fish", taskCondition2)
	gameTask2.Subscribe(gameLog.GlobalEventHub)

	taskCondition3 := func(e events.Event) bool {
		ev, ok := e.(entities.CreatureReachedPoint)
		return ok && ev.Point.PType == geometry.Food && ev.Creature.Hunger <= 1.0
	}

	gameTask3 := events.NewTask(entities.CreatureReachedPoint{}, "3. Feed them until they're full", taskCondition3)
	gameTask3.Subscribe(gameLog.GlobalEventHub)

	s.gameLog.Tasks = append(s.gameLog.Tasks, gameTask, gameTask2, gameTask3)

	s.ui = sUI
	s.subs(gameLog)
	timer := entities.NewTimer(1)
	s.nextSceneTrigger = timer
	return &s
}

func (s *StartScene) Update() (sceneManagement.SceneId, error) {
	s.ui.UI.Update()

	for _, fish := range s.ui.SelectSpritesToDraw {
		fish.Update()
	}

	s.nextSceneTrigger.Update()
	if s.nextSceneTrigger.TimerState == entities.Done {
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
	s.gameLog.SongPlayer.Play(soundFX.BestAdventureEver)

}

func (s *StartScene) OnEnter(gameLog *sceneManagement.GameLog) {

}

func (s *StartScene) OnExit() {
	s.gameLog.SongPlayer.Pause()
}

func (s *StartScene) IsLoaded() bool {
	return s.isLoaded
}

func (s *StartScene) subs(gameLog *sceneManagement.GameLog) {
	gameLog.GlobalEventHub.Subscribe(ui.ButtonEvent{}, func(e events.Event) {
		ev := e.(ui.ButtonEvent)
		if ev.EType == "cursor exited" {
			if ev.ButtonText != "Select" {
				if len(s.ui.SelectSpritesToDraw) > 1 {
					s.ui.SelectSpriteOptions[ev.ButtonText].UnLoadShader()
				}
			}
		}
	})

	gameLog.GlobalEventHub.Subscribe(ui.ButtonClickedEvent{}, func(e events.Event) {
		ev := e.(ui.ButtonClickedEvent)
		ols := shaders.LoadOutlineShader()
		switch ev.ButtonText {
		case "Common Molly":
			gameLog.SoundPlayer.Play(soundFX.SelectSound2)
			saveFish := entities.SavedFish{FishType: "mollyFish", Progress: 0, Size: 1}
			gameLog.Save.Fish = append(gameLog.Save.Fish, saveFish)
			s.ui.TextInputContainer.GetWidget().Disabled = false
			s.ui.TextInputContainer.GetWidget().Visibility = widget.Visibility_Show
			s.ui.TextInput.Focus(true)
			s.ui.SelectSpriteOptions[ev.ButtonText].LoadShader(ols)
		case "Goldfish":
			gameLog.SoundPlayer.Play(soundFX.SelectSound2)
			saveFish := entities.SavedFish{FishType: "fish", Progress: 0, Size: 1}
			gameLog.Save.Fish = append(gameLog.Save.Fish, saveFish)
			s.ui.TextInputContainer.GetWidget().Disabled = false
			s.ui.TextInputContainer.GetWidget().Visibility = widget.Visibility_Show
			s.ui.TextInput.Focus(true)
			s.ui.SelectSpriteOptions[ev.ButtonText].LoadShader(ols)
		case "Submit":
			ev2 := entities.SendData{
				DataFor: "Name Input",
				Data:    s.ui.TextInput.GetText(),
			}
			gameLog.GlobalEventHub.Publish(ev2)
		}
	})

	gameLog.GlobalEventHub.Subscribe(entities.SendData{}, func(e events.Event) {
		ev := e.(entities.SendData)
		if ev.DataFor == "Name Input" {
			gameLog.Save.Fish[0].Name = ev.Data
			s.ui.TextInputButton.Press()
			s.nextSceneTrigger.TurnOn()
			s.ui.TextInput.Focus(false)
			s.gameLog.SoundPlayer.Play(soundFX.SelectSound)
		}
	})

	gameLog.GlobalEventHub.Subscribe(ui.ButtonEvent{}, func(e events.Event) {
		ev := e.(ui.ButtonEvent)
		if ev.EType == "cursor entered" {
			ols := shaders.LoadRotatingHighlightShader()
			println(ev.ButtonText)
			if ev.ButtonText != "Select" {
				println(ev.ButtonText)
				s.ui.SelectSpriteOptions[ev.ButtonText].LoadShader(ols)
			}
		}

	})
}

func AcceptTextAndTriggerNextScene() {

}
