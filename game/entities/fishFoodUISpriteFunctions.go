package entities

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/input"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type FishFoodSprite struct {
	*UiSpriteData
	activationRect image.Rectangle
}

func FishFoodUpdate(ent *Entity) {
	ff := ent.UiData
	s := ent.Sprite

	if s.XYUpdater == nil {
		s.XYUpdater = sprite.NewUpdater(s)
	}

	s.Shader = registry.ShaderMap["Outline"]
	s.ShaderParams = make(map[string]any)
	s.ShaderParams["OutlineColor"] = [4]float64{1, 1, 0, 1}

	x, y := util.GetScaledCursorPosition()
	pt := image.Point{x, y}

	if ff.state == Activatable {

		s.ShaderParams["OutlineColor"] = [4]float64{0, 1, 0, 1}
		s.Scale = 0.95
		UpdateEntityZAndReSortEntitySlice(ent.Id, 0)

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			leftRect := ff.ActivationRect
			leftRect.Max.X = leftRect.Max.X / 2
			s.Img = ff.AltImg

			ev2 := input.MouseButtonPressedUISpriteActivity{
				Point: &util.Point{X: float32(x), Y: float32(y), PType: util.Food},
			}
			if pt.In(leftRect) {
				s.Img = ff.HoverImg
				ev2.Point.X += float32(s.Img.Bounds().Dx())
				ev2.Point.Tag = "left"
			}

			ff.EventHub.Publish(ev2)
		}

		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && s.Img != ff.MainImg {
			//change from pouring state if were not clicking
			s.Img = ff.MainImg
		}

	}

	if pt.In(ff.ActivationRect) && ff.state != Activatable {
		ff.state = Activatable
	}

	if !pt.In(ff.ActivationRect) && ff.state == Activatable {
		s.ShaderParams["OutlineColor"] = [4]float64{1, 1, 0, 1}
		ff.state = Selected
		s.Scale = 1.0
		// change from "pouring state" to holding state
		s.Img = ff.MainImg
		if ent.Z == 0 {
			UpdateEntityZAndReSortEntitySlice(ent.Id, 2)
		}
	}

}

func (ff *FishFoodSprite) Subscribe() {
	ff.EventHub.Subscribe(events.FishTankLayout{}, func(e tasks.Event) {
		ev := e.(events.FishTankLayout)
		ff.activationRect = ev.Rectangle
	})
}
