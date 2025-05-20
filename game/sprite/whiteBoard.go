package sprite

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/loaders"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
)

type WhiteBoardSprite struct {
	*UiSprite
	completedTasks []events.Task
	tasks          []events.Task
}

func (w *WhiteBoardSprite) Update() {
	w.clicked = false

	w.updateState()

	if w.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && len(w.completedTasks) > 0 {
		var taskPosition int

		for i, task := range w.tasks {
			if task.Text == w.completedTasks[0].Text {
				taskPosition = i
			}
		}

		ev := graphics.DrawGraphic{}

		vlsGraphic := loaders.NewVlS(
			w.X+20, w.Y+55+float32(20*taskPosition), w.X+22, w.Y+57+float32(20*taskPosition), w.X+170,
			color.RGBA{R: 20, G: 100, B: 100, A: 200},
		)

		ev.Graphic = vlsGraphic
		w.EventHub.Publish(ev)

		ev2 := events.TaskCompleted{}
		w.EventHub.Publish(ev2)

		w.completedTasks = w.completedTasks[1:]

		if len(w.completedTasks) <= 0 {
			w.Shader = nil
		}
	}

	if w.XYUpdater != nil {
		w.XYUpdater.Update()
	}

	w.stateWas = w.state

	w.UpdateShader()
}

func (w *WhiteBoardSprite) Subscribe(hub *events.EventHub) {
	hub.Subscribe(events.TaskRequirementsCompleted{}, func(e events.Event) {
		ev := e.(events.TaskRequirementsCompleted)
		taskPublished := false
		for _, task := range w.completedTasks {
			if task.Text == ev.Task.Text {
				taskPublished = true
			}
		}
		if !taskPublished {
			println("appending completed task to whiteboard")
			w.completedTasks = append(w.completedTasks, ev.Task)
			ols := shaders.LoadPulseOutlineShader()
			w.Shader = ols
			w.ShaderParams["OutlineColor"] = [4]float64{0.2, 0.7, 0.2, 255}
			w.ShaderParams["OutlineColor2"] = [4]float64{0.1, 0.9, 0.1, 255}
			w.ShaderParams["Counter"] = 0
			w.UpdateShaderParams = shaders.UpdatePulse
		}
	})

	hub.Subscribe(events.TaskCreated{}, func(e events.Event) {
		ev := e.(events.TaskCreated)
		w.tasks = append(w.tasks, ev.Task)
	})
}
