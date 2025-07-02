package interactableUIObjects

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/textEffects"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
)

type WhiteBoardSprite struct {
	*UiSprite
	textMap              map[int]*textEffects.TextWithShader
	dstImg               *sprite.Sprite
	dstShader            *ebiten.Shader
	textIndex            int
	textToBeWrittenQueue []string
	completedTaskQueue   []tasks.Task
	tasks                []*tasks.Task
	allTasksCompleted    bool
	timers               map[string]*entities.Timer
	allTasksBufferDone   bool
	oneTaskBufferDone    bool
	taskNumber           int
}

func (w *WhiteBoardSprite) Init() {
	dstImg := ebiten.NewImage(w.Img.Bounds().Dx(), w.Img.Bounds().Dy())
	w.dstImg = &sprite.Sprite{}
	w.dstImg.Img = dstImg
	w.timers = make(map[string]*entities.Timer)
	w.timers["TaskCreatedBuffer"] = entities.NewTimer(0.5)
	w.timers["AllTasksCompletedBuffer"] = entities.NewTimer(1)
	w.timers["EraseAnimationCompleted"] = entities.NewTimer(2)
	w.textMap = make(map[int]*textEffects.TextWithShader)
	w.dstImg.ShaderParams = make(map[string]any)
}

func (w *WhiteBoardSprite) ResetTextMap() {
	w.textMap = make(map[int]*textEffects.TextWithShader)
}

func (w *WhiteBoardSprite) ResetImg() {
	w.dstImg.Img.Clear()
}

func (w *WhiteBoardSprite) Update() {
	w.clicked = false

	w.updateState()

	if w.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && len(w.completedTaskQueue) > 0 {

		ev2 := tasks.TaskCompleted{}

		ev2.Task = w.completedTaskQueue[0]
		ev2.Slot = w.textIndex
		w.EventHub.Publish(ev2)

		w.completedTaskQueue = w.completedTaskQueue[1:]

		if len(w.completedTaskQueue) <= 0 {
			turnOffClickMeEffect(w.UiSprite)
		}
	}

	w.stateWas = w.state

	w.Sprite.Update()

	w.dstImg.Update()

	if w.allTasksCompleted {
		w.allTasksCompleted = false
		w.timers["AllTasksCompletedBuffer"].TurnOn()
	}

	if w.oneTaskBufferDone {
		w.updateTextToBeWritten()
		w.oneTaskBufferDone = false
	}

	if w.allTasksBufferDone {
		initClickMeEffect(w.UiSprite)
	}

	if w.allTasksBufferDone && w.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && w.Shader != nil {
		w.allTasksBufferDone = false
		w.EventHub.Publish(tasks.AllTasksCompleted{})
		w.timers["EraseAnimationCompleted"].TurnOn()
		w.reset()
	}

	if w.dstImg.Shader != nil {
		maxCounter, ok := w.dstImg.ShaderParams["MaxCounter"].(int)
		if ok {
			counter := w.dstImg.ShaderParams["Counter"].(int)
			if counter >= maxCounter {
				w.ResetImg()
				w.dstImg.UnLoadShader()
				w.dstImg.ShaderParams = make(map[string]any)
			}
		}
	}

	w.updateText()
	w.UpdateTimers()
}

func (w *WhiteBoardSprite) Draw(screen *ebiten.Image) {

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(w.X), float64(w.Y))

	w.DrawText()
	w.dstImg.Draw(w.Img)
	w.Sprite.Draw(screen)

}

func (w *WhiteBoardSprite) UpdateTimers() {

	for key, timer := range w.timers {
		state := timer.Update()
		switch key {
		case "TaskCreatedBuffer":
			if state == entities.Done {
				timer.TurnOff()
				w.oneTaskBufferDone = true
			}

		case "AllTasksCompletedBuffer":
			//Write out in marker effect/animation is completed
			if state == entities.Done {
				timer.TurnOff()
				w.allTasksBufferDone = true
			}

		case "EraseAnimationCompleted":
			//wipe away animation is completed
			if state == entities.Done {
				println("animation completed timer in whiteboard was triggered and completed")
				timer.TurnOff()
				w.ResetTextMap()
				w.textIndex = 0
				/*x := float64(w.Sprite.Img.Bounds().Dx())/2.0 - 45
				y := float64(w.Sprite.Img.Bounds().Dy()/2.0) - 25
				insets := [2]float64{x, y}*/
				w.ResetImg()
				w.appendTextToOpenSlot("All Done =)")
			}
		}
	}
}

func (w *WhiteBoardSprite) reset() {
	turnOffClickMeEffect(w.UiSprite)
	w.Shader = nil
	w.textIndex = 0
	w.allTasksCompleted = false
	w.tasks = []*tasks.Task{}
	w.completedTaskQueue = []tasks.Task{}
}

func (w *WhiteBoardSprite) DrawText() {

	for index, txt := range w.textMap {
		if txt != nil {
			if index == 0 && !txt.IsFullyDrawn() {
				//first one draw that shit
				txt.Draw(w.dstImg.Img)
				return
			}
			if index > 0 {
				// not first one, first one or greater finished, draw that shit
				if w.textMap[index-1].IsFullyDrawn() {
					txt.Draw(w.dstImg.Img)
				}
			}
		}
	}

}

func (w *WhiteBoardSprite) updateText() {
	for index, txt := range w.textMap {
		if index == 0 && !txt.IsFullyDrawn() {
			log.Printf("updating Text")
			txt.Update()
			return
		}

		if index > 0 {
			if w.textMap[index-1].IsFullyDrawn() {
				txt.Update()
			}
		}
	}
}

func (w *WhiteBoardSprite) Subscribe(hub *tasks.EventHub) {
	hub.Subscribe(tasks.TaskRequirementsCompleted{}, func(e tasks.Event) {

		taskPublished := false //hack way to limit to one event for each task
		ev := e.(tasks.TaskRequirementsCompleted)

		for _, task := range w.completedTaskQueue {
			if task.Text == ev.Task.Text {
				taskPublished = true //hack way to limit to one event for each task
			}
		}

		if !taskPublished {
			println("appending completed task to whiteboard")
			w.completedTaskQueue = append(w.completedTaskQueue, ev.Task)
			if w.textMap[w.textIndex] != nil {
				w.textMap[w.textIndex].FullyDrawn = true
			}
			initClickMeEffect(w.UiSprite)
		}
	})

	hub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
		ev := e.(tasks.TaskCreated)
		if len(w.tasks) == 0 {
			log.Println("Appending first task")
			w.appendTextToOpenSlot(ev.Task.Text)
			w.tasks = append(w.tasks, ev.Task)
			return
		}
		w.textToBeWrittenQueue = append(w.textToBeWrittenQueue, ev.Task.Text)
		w.timers["TaskCreatedBuffer"].TurnOn()
	})

	hub.Subscribe(tasks.TaskCompleted{}, func(e tasks.Event) {
		w.checkAllTasksCompleted()
	})

	hub.Subscribe(tasks.AllTasksCompleted{}, func(e tasks.Event) {
		w.initErase()
	})

	hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		w.initErase()
	})

	hub.Subscribe(events.DayOverTransitionComplete{}, func(e tasks.Event) {
		w.ResetTextMap()
		w.ResetImg()
		w.reset()
		w.textIndex = 0
	})

	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		ev := events.UISpriteLayedOut{}
		ev.Label = "Whiteboard"
		ev.X = w.X
		ev.Y = w.Y
		hub.Publish(ev)
	})

	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		ev := e.(events.NewDay)
		w.taskNumber = ev.NTasks
	})

}
func (w *WhiteBoardSprite) updateTextToBeWritten() {

	if len(w.textToBeWrittenQueue) > 0 {
		log.Printf("adding text: %s to queue", w.textToBeWrittenQueue[0])
		w.appendTextToOpenSlot(w.textToBeWrittenQueue[0])
		w.textToBeWrittenQueue = w.textToBeWrittenQueue[1:]
	}

}

func (w *WhiteBoardSprite) initErase() {
	w.dstImg.LoadShader(registry.ShaderMap["Erase"])
	w.dstImg.ShaderParams = make(map[string]any)
	w.dstImg.ShaderParams["Counter"] = 0
	w.dstImg.ShaderParams["MaxCounter"] = 100
	w.dstImg.UpdateShaderParams = shaders.UpdateCounterOneShot
}

func (w *WhiteBoardSprite) checkAllTasksCompleted() {

	if len(w.completedTaskQueue) > 0 {
		if w.completedTaskQueue[len(w.completedTaskQueue)-1].Index == w.taskNumber-1 {
			w.allTasksCompleted = true
		}
	}

}

func (w *WhiteBoardSprite) appendTextToOpenSlot(txt string) {
	insets := [2]float64{10, float64((w.textIndex + 1) * 20)}
	ts := textEffects.NewTextWithMarkerShader(txt, w.dstImg.Img.Bounds(), insets, ColorScaleSlice[0])
	log.Printf("whiteboard text index = %d, appending new text: %s at this slot", w.textIndex, txt)

	w.textMap[w.textIndex] = ts
	w.textIndex++
}
