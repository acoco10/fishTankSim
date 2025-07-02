package scenes

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/daySystem"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
)

type StartScene struct {
	ui               *ui.StartMenu
	isLoaded         bool
	fishOptions      []*sprite.AnimatedSprite
	nextSceneTrigger *entities.Timer
	gameLog          *sceneManagement.GameLog
	selectedFish     entities.SavedFish
	selectedProp     entities.TankObject
}

func NewStartScene(gameLog *sceneManagement.GameLog) *StartScene {
	s := StartScene{}

	s.gameLog = gameLog
	sUI, err := ui.LoadStartMenu(gameLog.GlobalEventHub, ScreenWidth, ScreenHeight)
	if err != nil {
		log.Fatal(fmt.Errorf("error initiating start menu: %s", err))
	}
	daySystem.LoadDaysTasks(gameLog)
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

	if s.ui.DrawOptions["Back"].SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.ui.Back()
	}

	return sceneManagement.StartScene, nil
}

func (s *StartScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 120, G: 170, B: 210, A: 255})
	opts := &ebiten.DrawImageOptions{}
	sopts := &ebiten.DrawRectShaderOptions{}

	s.ui.Draw(screen)

	for _, fish := range s.ui.SelectSpritesToDraw {
		fish.Draw(screen)
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
	log.Printf("Leaving Start Scene")
	s.gameLog.SongPlayer.Pause()
}

func (s *StartScene) IsLoaded() bool {
	return s.isLoaded
}

func (s *StartScene) subs(gameLog *sceneManagement.GameLog) {
	gameLog.GlobalEventHub.Subscribe(events.ButtonEvent{}, func(e tasks.Event) {

		ev := e.(events.ButtonEvent)
		if ev.EType == "cursor exited" {
			if ev.ButtonText != "Select" {
				if len(s.ui.SelectSpritesToDraw) > 1 {
					//this filter logic doesnt follow from anything about the code
					s.ui.DrawOptions[ev.ButtonText].(*sprite.AnimatedSprite).UnLoadShader()
				}
			}
		}

	})

	gameLog.GlobalEventHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)

		switch ev.ButtonText {

		case "Common Molly":
			gameLog.SoundPlayer.Play(soundFX.SelectSound2)
			saveFish := entities.SavedFish{FishType: "mollyFish", Progress: 0, Size: 1}
			s.selectedFish = saveFish

		case "Goldfish":
			gameLog.SoundPlayer.Play(soundFX.SelectSound2)
			saveFish := entities.SavedFish{FishType: "fish", Progress: 0, Size: 1}
			s.selectedFish = saveFish

		case "Castle", "Log":
			prop := entities.TankObject{Name: ev.ButtonText}
			s.selectedProp = prop

		case "Submit":
			s.nextSceneTrigger.TurnOn()
			s.gameLog.SoundPlayer.Play(soundFX.SelectSound)
			//s.selectedFish.Name = s.ui.TextInput.GetText()
			gameLog.Save.Fish = append(gameLog.Save.Fish, s.selectedFish)
			gameLog.Save.TankObjects = append(gameLog.Save.TankObjects, s.selectedProp)
		}
	})

	gameLog.GlobalEventHub.Subscribe(events.ButtonEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonEvent)
		if ev.EType == "cursor entered" {
			if ev.ButtonText != "Select" {
				s.ui.DrawOptions[ev.ButtonText].(*sprite.AnimatedSprite).LoadShader(registry.ShaderMap["Outline"])
			}
		}

	})
}
