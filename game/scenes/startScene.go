package scenes

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/daySystem"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"log"
)

type StartScene struct {
	i                      int
	ui                     *ui.StartMenu
	isLoaded               bool
	fishOptions            []*sprite.AnimatedSprite
	nextSceneTrigger       *entities.Timer
	gameLog                *sceneManagement.GameLog
	selectedFish           entities.SavedFish
	smallerResolution      *ebiten.Image
	shaderUpdater          func(map[string]any) map[string]any
	selectedProp           entities.TankObject
	backGroundShader       *ebiten.Shader
	offScreen1             *ebiten.Image
	backGroundShaderParams map[string]any
	resolutionScaling      int
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
	s.smallerResolution = ebiten.NewImage(640, 360)
	s.resolutionScaling = ScaleScreenToResolution(s.smallerResolution)

	s.offScreen1, err = loader.LoadImageAssetAsEbitenImage("roomImages/startBackGround")
	if err != nil {
		log.Fatal("bad import path ya dingbat")
	}

	s.backGroundShader = registry.ShaderMap["Water"]
	s.backGroundShaderParams = make(map[string]any)
	s.backGroundShaderParams["TankRect"] = [4]float64{0, 0, ScreenWidth, ScreenHeight}
	s.backGroundShaderParams["TankSize"] = [2]float64{ScreenWidth, ScreenHeight}
	s.backGroundShaderParams["Counter"] = 0

	s.shaderUpdater = shaders.UpdateCounter
	return &s
}

func (s *StartScene) Update() (sceneManagement.SceneId, error) {
	s.i++
	s.ui.UI.Update()

	s.backGroundShaderParams = s.shaderUpdater(s.backGroundShaderParams)

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
	s.offScreen1.Fill(colornames.Aliceblue)
	shaderOpts := &ebiten.DrawRectShaderOptions{}

	shaderOpts.Uniforms = s.backGroundShaderParams
	shaderOpts.Images[0] = s.offScreen1

	s.smallerResolution.DrawRectShader(ScreenWidth, ScreenHeight, s.backGroundShader, shaderOpts)

	opts := &ebiten.DrawImageOptions{}
	if s.i > 2 {
	}

	for _, fish := range s.ui.SelectSpritesToDraw {
		if s.smallerResolution != nil {
			fish.Draw(s.smallerResolution)
		}
	}

	opts.GeoM.Scale(float64(s.resolutionScaling), float64(s.resolutionScaling))

	screen.DrawImage(s.smallerResolution, opts)
	opts.GeoM.Reset()

	s.ui.Draw(screen)

}

func (s *StartScene) FirstLoad() {
	s.isLoaded = true
	s.gameLog.SongPlayer.Play(soundFX.BestAdventureEver)
}

func (s *StartScene) OnEnter() {

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
			gameLog.SoundPlayer.Play(soundFX.SelectSound2)
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
