package gameEntities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
)

type WhiteBoardSprite struct {
	*UiSprite
	TaskQueue []*Task
}

func (w *WhiteBoardSprite) Update() {
	w.clicked = false

	w.updateState()

	if w.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && w.Shader != nil {
		w.Shader = nil
		ev := DrawGraphic{}
		graphic := VectorLineGraphic{}

		graphic.X0 = w.X + 20
		graphic.Y0 = w.Y + 74

		graphic.X1 = w.X + 22
		graphic.Y1 = w.Y + 72
		graphic.clr = color.RGBA{20, 100, 100, 200}
		graphic.maxX1 = w.X + 170

		ev.Graphic = &graphic
		w.EventHub.Publish(ev)

		ev2 := TaskCompleted{}
		w.EventHub.Publish(ev2)
	}

	if w.XYUpdater != nil {
		w.XYUpdater.Update()
	}

	w.stateWas = w.state

	w.UpdateShader()
}
