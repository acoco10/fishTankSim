package scenes

import (
	"github.com/acoco10/fishTankWebGame/game/daySystem"
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"log"
	"math"
	"math/rand"
)

type TransitionScene struct {
	isLoaded          bool
	nextSceneTrigger  *util.Timer
	dotTimer          *util.Timer
	gameLog           *sceneManagement.GameLog
	dots              int
	loadingMsg        string
	stars             []*sprite.Sprite
	smallerResolution *ebiten.Image
}

func LoadTransitionScene(gameLog *sceneManagement.GameLog) *TransitionScene {
	s := TransitionScene{}
	s.gameLog = gameLog
	s.dots = 1
	s.loadingMsg = "Next Day Loading"
	return &s
}

func (s *TransitionScene) Update() (sceneManagement.SceneId, error) {
	for _, star := range s.stars {
		star.Update()
	}
	if s.nextSceneTrigger.TimerState == util.Done {
		return sceneManagement.FishTank, nil
	}
	s.nextSceneTrigger.Update()
	s.dotTimer.Update()
	return sceneManagement.TransitionScene, nil
}

func (s *TransitionScene) Draw(screen *ebiten.Image) {

	/*	face, err := util.LoadFont(24, "nk57")
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
		text.Draw(screen, s.loadingMsg, face, dopts)*/
	s.smallerResolution.Fill(colornames.Midnightblue)

	for _, star := range s.stars {
		star.Draw(s.smallerResolution)
	}

	dopts := &ebiten.DrawImageOptions{}
	dopts.GeoM.Scale(registry.Config.ResolutionScalingF, registry.Config.ResolutionScalingF)

	screen.DrawImage(s.smallerResolution, dopts)

}

func (s *TransitionScene) FirstLoad() {
	s.isLoaded = true
	s.dotTimer = util.NewTimer(0.5)
	s.nextSceneTrigger = util.NewTimer(2)

	s.stars = CreateRandomSemiGridSprites(75, registry.Config.ScreenWidth, registry.Config.ScreenHeight, "data/animationData/glowInTheDarkStar.json")

	s.smallerResolution = ebiten.NewImage(registry.Config.ScreenWidth, registry.Config.ScreenHeight)
}

func (s *TransitionScene) OnEnter() {

	s.nextSceneTrigger.TurnOn()
	s.dotTimer.TurnOn()

	daySystem.LoadDaysTasks(s.gameLog)
}

func (s *TransitionScene) OnExit() {

	s.nextSceneTrigger.TurnOff()
	s.dotTimer.TurnOff()
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

func CreateRandomSemiGridSprites(numSprites int, screenWidth, screenHeight int, animationPath string) []*sprite.Sprite {
	var sprites []*sprite.Sprite

	// Calculate rough grid dimensions
	cols := int(math.Sqrt(float64(numSprites)) * 1.5) // Slightly more columns than rows
	rows := (numSprites + cols - 1) / cols            // Ceiling division

	cellWidth := float64(screenWidth) / float64(cols)
	cellHeight := float64(screenHeight) / float64(rows)

	// Randomization parameters
	jitterAmount := 0.4 // How much to randomly offset from grid positions (0-1)
	skipChance := 0.1   // Chance to skip a grid position (creates gaps)

	spriteIndex := 0

	for row := 0; row < rows && spriteIndex < numSprites; row++ {
		for col := 0; col < cols && spriteIndex < numSprites; col++ {
			// Skip some positions randomly to break up the grid
			if rand.Float64() < skipChance {
				continue
			}

			starImg, err := util.LoadImageAssetAsEbitenImage("effectSpriteSheets/glowInTheDarkStar")
			if err != nil {
				log.Fatal(err)
			}

			animation, err := entImportableLoaders.LoadAnimation("data/animationData/glowInTheDarkStar.json")
			if err != nil {
				log.Fatal(err)
			}
			animation.Img = starImg

			// Base grid position
			baseX := float64(col)*cellWidth + cellWidth/2
			baseY := float64(row)*cellHeight + cellHeight/2

			// Add random jitter within the cell
			jitterX := (rand.Float64() - 0.5) * cellWidth * jitterAmount
			jitterY := (rand.Float64() - 0.5) * cellHeight * jitterAmount

			finalX := baseX + jitterX
			finalY := baseY + jitterY

			// Clamp to screen boundaries
			finalX = math.Max(0, math.Min(float64(screenWidth), finalX))
			finalY = math.Max(0, math.Min(float64(screenHeight), finalY))

			sp := &sprite.Sprite{
				AnimationMap:     map[string]*sprite.Animation{"Default": animation},
				CurrentAnimation: "Default",
				X:                float32(finalX),
				Y:                float32(finalY),
			}

			sprites = append(sprites, sp)

			spriteIndex++
		}
	}

	return sprites
}
