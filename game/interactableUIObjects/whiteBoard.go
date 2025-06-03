package interactableUIObjects

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/textEffects"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type WhiteBoardSprite struct {
	*UiSprite
	textMap map[int]textEffects.TextWithShader
	completedTaskQueue []events.Task
	tasks              []*events.Task
	allTasksCompleted  bool
	timers             map[string]*entities.Timer
	bufferDone         bool
}

func (w *WhiteBoardSprite) Init() {
	w.timers = make(map[string]*entities.Timer)
	w.timers["TasksCompletedBuffer"] = entities.NewTimer(1)
	w.textMap = make(map[int]textEffects.TextWithShader)
}

func (w *WhiteBoardSprite) UpdateTimers() {
	for key, timer := range w.timers {
		state := timer.Update()
		if key == "TasksCompletedBuffer" {
			if state == entities.Done {
				timer.TurnOff()
				w.bufferDone = true
			}
		}
	}
}

func (w *WhiteBoardSprite) Update() {
	w.clicked = false

	w.updateState()

	if w.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && len(w.completedTaskQueue) > 0 {

		ev2 := events.TaskCompleted{}

		ev2.Task = w.completedTaskQueue[0]

		w.EventHub.Publish(ev2)

		w.completedTaskQueue = w.completedTaskQueue[1:]

		if len(w.completedTaskQueue) <= 0 {
			w.Shader = nil
		}
	}

	if w.XYUpdater != nil {
		w.XYUpdater.Update()
	}

	w.stateWas = w.state

	w.checkAllTasksCompleted()

	w.UpdateShader()

	if w.allTasksCompleted {
		w.timers["TasksCompletedBuffer"].TurnOn()

	}

	if w.bufferDone {
		w.initClickMeShader()
	}

	if w.allTasksCompleted && w.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && w.bufferDone {
		w.EventHub.Publish(events.AllTasksCompleted{})
	}

	w.UpdateTimers()
}

func (w *WhiteBoardSprite) Subscribe(hub *events.EventHub) {
	hub.Subscribe(events.TaskRequirementsCompleted{}, func(e events.Event) {
		taskPublished := false //hack way to limit to one event for each task
		ev := e.(events.TaskRequirementsCompleted)

		for _, task := range w.completedTaskQueue {
			if task.Text == ev.Task.Text {
				taskPublished = true //hack way to limit to one event for each task
			}
		}

		if !taskPublished {
			println("appending completed task to whiteboard")
			w.completedTaskQueue = append(w.completedTaskQueue, ev.Task)
			w.initClickMeShader()
		}
	})

	hub.Subscribe(events.TaskCreated{}, func(e events.Event) {

		ev := e.(events.TaskCreated)
		w.tasks = append(w.tasks, ev.Task)
	})


	hub.Subscribe(events.TaskCreated{}, func(e events.Event) {
		ev := e.(events.TaskCreated)

	})

	hub.Subscribe(events.TaskCompleted{}, func(e events.Event) {
		//ev := e.(entities.TaskCompleted)
	})

	hub.Subscribe(events.AllTasksCompleted{}, func(e events.Event) {
		t.ReplaceTextArea("All Done =)")
	})
}
}

func (w *WhiteBoardSprite) initClickMeShader() {
	ols := shaders.LoadPulseOutlineShader()
	w.Shader = ols
	w.ShaderParams["OutlineColor"] = [4]float64{0.2, 0.7, 0.2, 255}
	w.ShaderParams["OutlineColor2"] = [4]float64{0.1, 0.9, 0.1, 255}
	w.ShaderParams["Counter"] = 0
	w.UpdateShaderParams = shaders.UpdateCounter
}

func (w *WhiteBoardSprite) checkAllTasksCompleted() {

	allTasksCompleted := true

	for _, task := range w.tasks {
		if !task.Completed {
			allTasksCompleted = false
			break
		}
	}

	w.allTasksCompleted = allTasksCompleted
}


func (w *WhiteBoardSprite) appendTextToOpenSlot(txt string){
	ts := textEffects.NewTextWithShader(txt, w.Img)
}