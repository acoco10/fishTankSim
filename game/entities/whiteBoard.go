package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"log"
)

type WhiteBoardSprite struct {
	currentTask         *tasks.Task
	textBeingWritten    *graphics.TextWithShader
	DstImg              *sprite.Sprite
	dstShader           *ebiten.Shader
	allTasksCompleted   bool
	timers              map[string]*util.Timer
	allTasksBufferDone  bool
	oneTaskBufferDone   bool
	numberOfTasksToday  int
	graphicManager      *graphics.GraphicManager
	tasksCompletedToday int
	UiSprite            *Entity
	taskJustCompleted   bool
	lastTaskID          int
	spacing             int
	clickme             bool
	windowOpen          bool
}

func (w *WhiteBoardSprite) Init(eventHub *tasks.EventHub) {
	w.spacing = 30
	imgSizeWidth := int(float64(w.UiSprite.Sprite.Img.Bounds().Dx()) * registry.Config.ResolutionScalingF)
	imgSizeHeight := int(float64(w.UiSprite.Sprite.Img.Bounds().Dy()) * registry.Config.ResolutionScalingF)

	imgX := w.UiSprite.Sprite.X * registry.Config.ResolutionScalingf
	imgY := w.UiSprite.Sprite.Y*registry.Config.ResolutionScalingf + float32(registry.Config.YOffsetF)

	dstImg := ebiten.NewImage(imgSizeWidth, imgSizeHeight)
	w.DstImg = &sprite.Sprite{Img: dstImg, X: imgX, Y: imgY}

	w.timers = make(map[string]*util.Timer)

	w.timers["TaskCreatedBuffer"] = util.NewTimer(0.5)
	w.timers["AllTasksCompletedBuffer"] = util.NewTimer(1)
	w.timers["EraseAnimationCompleted"] = util.NewTimer(2)

	w.DstImg.ShaderParams = make(map[string]any)

	w.graphicManager = &graphics.GraphicManager{}
	w.graphicManager.Params = make(map[string]any)
	w.graphicManager.Params["Coordinates"] = image.Point{X: int(w.UiSprite.Sprite.X), Y: int(w.UiSprite.Sprite.Y)}
	w.graphicManager.Params["Width"] = float32(imgSizeWidth)
	w.graphicManager.Params["Height"] = imgSizeHeight
	w.graphicManager.Params["Index"] = 0
	w.graphicManager.Params["Spacing"] = w.spacing
	w.graphicManager.DstImage = w.DstImg.Img

	w.graphicManager.Timers = make(map[string]*util.Timer)
	w.graphicManager.Timers["DeInit"] = util.NewTimer(1)
	w.graphicManager.Tag = "whiteBoard"

	graphics.WhiteBoardGMSubs(w.graphicManager, eventHub)
}

func (w *WhiteBoardSprite) ResetImg() {
	w.DstImg.Img.Clear()
}

func (w *WhiteBoardSprite) Update() {

	sp := w.UiSprite.Sprite
	uiDat := w.UiSprite.UiData

	w.graphicManager.UpdateTimers()

	if sp.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && w.taskJustCompleted {
		if !w.windowOpen {
			w.tasksCompletedToday++
			w.taskJustCompleted = false
			ev2 := tasks.TaskCompleted{Task: *w.currentTask}
			w.UiSprite.EventHub.Publish(ev2)
			turnOffClickMeEffect(w.UiSprite.UiData)
		}

	}

	w.DstImg.Update()

	if w.allTasksCompleted {
		w.allTasksCompleted = false
		w.timers["AllTasksCompletedBuffer"].TurnOn()
	}

	if w.allTasksBufferDone && w.clickme {
		w.clickme = false
		initClickMeEffect(uiDat)
	}

	if w.allTasksBufferDone && sp.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		w.allTasksBufferDone = false
		turnOffClickMeEffect(w.UiSprite.UiData)
		w.UiSprite.EventHub.Publish(tasks.AllTasksCompleted{})
		w.timers["EraseAnimationCompleted"].TurnOn()
	}

	if w.DstImg.Shader != nil {
		maxCounter, ok := w.DstImg.ShaderParams["MaxCounter"].(int)
		if ok {
			counter := w.DstImg.ShaderParams["Counter"].(int)
			if counter >= maxCounter {
				graphics.DeInitGraphicId(w.lastTaskID)
				w.DstImg.Img.Clear()
				w.DstImg.Shader = nil
			}
		}
	}
	w.UpdateTimers()
}

func (w *WhiteBoardSprite) UpdateTimers() {

	for key, timer := range w.timers {
		state := timer.Update()
		switch key {

		case "TaskCreatedBuffer":
			if state == util.Done {
				timer.TurnOff()
				w.oneTaskBufferDone = true
			}

		case "AllTasksCompletedBuffer":
			//Write out in marker effect/animation is completed
			if state == util.Done {
				timer.TurnOff()
				w.allTasksBufferDone = true
				w.clickme = true
			}

		case "EraseAnimationCompleted":
			//wipe away animation is completed
			if state == util.Done {
				println("animation completed timer in whiteboard was triggered and completed")
				timer.TurnOff()
				w.appendTextToOpenSlot("All Done =)")
			}
		}
	}
}

func (w *WhiteBoardSprite) reset() {
	sp := w.UiSprite.Sprite
	turnOffClickMeEffect(w.UiSprite.UiData)
	sp.Shader = nil
	w.allTasksCompleted = false
}

func (w *WhiteBoardSprite) Subscribe(hub *tasks.EventHub) {
	hub.Subscribe(tasks.TaskRequirementsCompleted{}, func(e tasks.Event) {
		fmt.Printf("task requirment event recieved")
		initClickMeEffect(w.UiSprite.UiData)
		w.taskJustCompleted = true
	})

	hub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
		log.Printf("new task recieved")
		ev := e.(tasks.TaskCreated)
		w.currentTask = ev.Task
		w.appendTextToOpenSlot(ev.Task.Text)
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

	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		ev := e.(events.NewDay)
		w.numberOfTasksToday = ev.NTasks
		w.tasksCompletedToday = 0
		fmt.Printf("White board received %d tasks on day %d", ev.NTasks, ev.Day)
	})

	hub.Subscribe(events.WriteToWhiteBoard{}, func(event tasks.Event) {
		ev := event.(events.WriteToWhiteBoard)
		w.appendTextToBestSpot(ev)
	})
	hub.Subscribe(events.WindowOpened{}, func(event tasks.Event) {
		w.windowOpen = true
	})
	hub.Subscribe(events.WindowClosed{}, func(event tasks.Event) {
		w.windowOpen = false
	})

}

func (w *WhiteBoardSprite) initErase() {
	w.DstImg.LoadShader(registry.ShaderMap["Erase"])
	w.DstImg.ShaderParams = make(map[string]any)
	w.DstImg.ShaderParams["Counter"] = 0
	w.DstImg.ShaderParams["MaxCounter"] = 100
	w.DstImg.UpdateShaderParams = shaders.UpdateCounterOneShot
}

func (w *WhiteBoardSprite) checkAllTasksCompleted() {

	if w.numberOfTasksToday == w.tasksCompletedToday {
		w.allTasksCompleted = true
	}

}

func (w *WhiteBoardSprite) appendTextToOpenSlot(txt string) {
	if w.lastTaskID != 0 {
		graphics.DeInitGraphicId(w.lastTaskID)
	}
	yinset := float64(w.spacing*(w.tasksCompletedToday+1)) + 10.0
	yOffSet := registry.Config.ScaledYOffsetF
	insets := [2]float64{float64(w.UiSprite.Sprite.X * registry.Config.ResolutionScalingf), float64(w.UiSprite.Sprite.Y*registry.Config.ResolutionScalingf) + yOffSet}
	w.lastTaskID = graphics.AddHandwritingGraphic(txt, w.DstImg.Img, insets, yinset, 20.0)

}

func (w *WhiteBoardSprite) appendTextToBestSpot(ev events.WriteToWhiteBoard) {
	if w.lastTaskID != 0 {
		graphics.DeInitGraphicId(w.lastTaskID)
	}
	switch ev.PreferredPosition {
	case "bottomLeft":
		yOffSet := registry.Config.ScaledYOffsetF
		insets := [2]float64{
			float64(w.UiSprite.Sprite.X) * registry.Config.ResolutionScalingF,
			float64(w.UiSprite.Sprite.Y)*registry.Config.ResolutionScalingF + yOffSet,
		}
		w.lastTaskID = graphics.AddHandwritingGraphic(ev.Msg, w.DstImg.Img, insets, 130*registry.Config.ResolutionScalingF, 20)
	}

}
