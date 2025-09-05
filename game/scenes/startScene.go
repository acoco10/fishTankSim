package scenes

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
)

type StartScene struct {
	ui                     *ui.StartMenu
	isLoaded               bool
	fishOptions            []*sprite.AnimatedSprite
	nextSceneTrigger       *util.Timer
	gameLog                *sceneManagement.GameLog
	selectedFish           entities.SavedFish
	smallerResolution      *ebiten.Image
	shaderUpdater          func(map[string]any) map[string]any
	selectedProp           entities.TankObject
	backGroundShaderParams map[string]any
	titleText              string
	resolutionScaling      int
	usedGraphId            int
}

func NewStartScene(gameLog *sceneManagement.GameLog) *StartScene {
	s := StartScene{}
	s.smallerResolution = ebiten.NewImage(ScreenWidth, ScreenHeight)

	s.gameLog = gameLog

	sUI, err := ui.LoadStartMenu(gameLog.GlobalEventHub, registry.Config.ResolutionScalingF)
	if err != nil {
		log.Fatal(fmt.Errorf("error initiating start menu: %s", err))
	}

	s.ui = sUI
	s.subs(gameLog)

	timer := util.NewTimer(1)
	s.nextSceneTrigger = timer

	return &s
}

func (s *StartScene) Update() (sceneManagement.SceneId, error) {
	s.ui.UI.Update()
	//graphics.UpdateGraphics()

	s.nextSceneTrigger.Update()

	if s.nextSceneTrigger.TimerState == util.Done {
		s.ui.UI.ClearFocus()
		return sceneManagement.FishTank, nil
	}

	/*if s.ui.DrawOptions["Back"].SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.ui.Back()
	}*/

	return sceneManagement.StartScene, nil
}

func (s *StartScene) Draw(screen *ebiten.Image) {

	s.smallerResolution.Fill(color.RGBA{R: 120, G: 170, B: 210, A: 255})

	dOpts := &ebiten.DrawImageOptions{}
	dOpts.GeoM.Translate(0, registry.Config.YOffsetF)
	dOpts.GeoM.Scale(registry.Config.ResolutionScalingF, registry.Config.ResolutionScalingF)
	screen.DrawImage(s.smallerResolution, dOpts)

	s.ui.Draw(screen)
	graphics.DrawUnScaledGraphics(screen)
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

	gameLog.GlobalEventHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		switch ev.ButtonText {
		case "Common Molly":
			gameLog.SoundPlayer.Play(soundFX.SelectSound2)
			saveFish := entities.SavedFish{FishType: "mollyFish", Progress: 0, Size: 1}
			s.selectedFish = saveFish

		case "Goldfish":
			gameLog.SoundPlayer.Play(soundFX.SelectSound2)
			saveFish := entities.SavedFish{FishType: "goldFish", Progress: 0, Size: 1}
			s.selectedFish = saveFish

		case "Submit":
			s.nextSceneTrigger.TurnOn()
			s.gameLog.SoundPlayer.Play(soundFX.SelectSound)
			s.selectedFish.Name = s.ui.TextInput.GetText()
			gameLog.Save.Fish = append(gameLog.Save.Fish, s.selectedFish)
		}
	})
}
