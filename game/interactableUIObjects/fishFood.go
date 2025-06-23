package interactableUIObjects

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/input"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"math"
)

type FishFoodSprite struct {
	*UiSprite
	activationRect image.Rectangle
}

func (ff *FishFoodSprite) Draw(screen *ebiten.Image) {

	opts := ebiten.DrawImageOptions{}

	if ff.state == Idle {
		opts.GeoM.Translate(float64(ff.X), float64(ff.Y))
		screen.DrawImage(ff.Img, &opts)
		opts.GeoM.Reset()
	} else if ff.state == Selected || ff.state == HoveredOver {
		if ff.HoverImg != nil {
			opts.GeoM.Translate(float64(ff.X), float64(ff.Y))
			opts.GeoM.Translate(float64(ff.AltOffsetX), float64(ff.AltOffsetY))
			screen.DrawImage(ff.HoverImg, &opts)
			opts.GeoM.Reset()
		}
	} else if ff.state == ClickedWhileBeingSelected {
		opts.GeoM.Translate(float64(ff.X), float64(ff.Y))
		screen.DrawImage(ff.AltImg, &opts)
		opts.GeoM.Reset()
	}

}

func (ff *FishFoodSprite) Update() {
	ff.clicked = false

	ff.updateState()

	if ff.state == HoveredOver && ff.stateWas != HoveredOver {
		ols := shaders.LoadOutlineShader()
		ff.Shader = ols
		ff.ShaderParams = make(map[string]any)
		ff.ShaderParams["OutlineColor"] = [4]float64{1, 1, 0, 1}

	}

	if ff.Shader != nil && (ff.state != HoveredOver && ff.state != Selected) {
		ff.Shader = nil
	}

	if ff.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && ff.state == Selected {

		if ff.XYUpdater == nil {
			ev := events.UISpriteAction{}
			ev.UiSprite = ff.Label
			ev.UiSpriteAction = "picked up"
			ff.EventHub.Publish(ev)
		}

		ff.XYUpdater = sprite.NewUpdater(ff.Sprite)
	}

	if ff.XYUpdater != nil {
		ff.XYUpdater.Update()
	}

	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && ff.state == ClickedWhileBeingSelected && !ff.clicked {
		ff.clicked = true
		ev := input.MouseButtonPressedUISpriteActivity{
			Point: &geometry.Point{X: float32(x), Y: float32(y), PType: geometry.Food},
		}
		ff.EventHub.Publish(ev)
	}

	baseDis := math.Hypot(float64(ff.X-ff.baseX), float64(ff.Y-ff.baseY))

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if ff.state == Selected && baseDis < 100 && ff.stateWas == Selected {
			ff.returnToBase()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyE) {
		if baseDis != 0 {
			ff.returnToBase()
		}
	}

	ff.stateWas = ff.state
}

func (ff *FishFoodSprite) updateState() {
	if !ff.SpriteHovered() {
		ff.state = Idle
	}

	if ff.SpriteHovered() && (ff.state != ClickedWhileBeingSelected && ff.state != Selected) {
		ff.state = HoveredOver
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && ff.state == HoveredOver {
		ff.state = Selected
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && ff.state == Selected && ff.stateWas == Selected {
		xCheck := ff.X > float32(ff.activationRect.Min.X)+100 && ff.X < float32(ff.activationRect.Max.X)
		yCheck := float32(ff.activationRect.Max.Y) > ff.Y
		if xCheck && yCheck {
			ff.state = ClickedWhileBeingSelected
		}
	}

	if ff.state == ClickedWhileBeingSelected && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		ff.state = Selected
	}
}

func (ff *FishFoodSprite) Subscribe() {
	ff.EventHub.Subscribe(events.FishTankLayout{}, func(e tasks.Event) {
		ev := e.(events.FishTankLayout)
		ff.activationRect = ev.Rectangle
	})
}
