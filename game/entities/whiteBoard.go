package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/stringConstants"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"log"
	"math/rand"
	"strings"
)

const (
	dailyTasks = "Daily Tasks"

	UpperCenter = "upperCenter"
	Center      = "center"

	CrossOut = "crossOut"

	//timers
	Defer = "defer"

	//state
	IdleWB            = 0
	Writing           = 1
	NotReceivingTasks = 2
)

type WhiteBoardSprite struct {
	currentTask               *tasks.Task
	textBeingWritten          *graphics.TextWithShader
	DstImg                    *sprite.Sprite
	NoEraseDst                *sprite.Sprite
	dstShader                 *ebiten.Shader
	state                     uint8
	allTasksCompleted         bool
	timers                    map[string]*util.Timer
	allTasksBufferDone        bool
	numberOfTasksToday        int
	dayOver                   bool
	graphicManager            *graphics.GraphicManager
	tasksCompletedToday       int
	UiSprite                  *Entity
	taskJustCompleted         bool
	lastTaskID                int
	spacing                   int
	clickme                   bool
	windowOpen                bool
	tIds                      []uint32
	msgsToWriteLater          []WriteToWhiteBoard
	occupied                  map[string]float64
	crossOutTexture           *ebiten.Image
	pubQueue                  []tasks.Event
	clickAfterWriting         bool
	day                       int
	YSelectOptionsForCrossOut []int
}

func (w *WhiteBoardSprite) Init(eventHub *tasks.EventHub) {
	w.spacing = 20
	imgSizeWidth := w.UiSprite.Sprite.Img.Bounds().Dx()
	imgSizeHeight := w.UiSprite.Sprite.Img.Bounds().Dy()
	dstImg := ebiten.NewImage(imgSizeWidth, imgSizeHeight)

	w.DstImg = &sprite.Sprite{Img: dstImg, X: w.UiSprite.Sprite.X, Y: w.UiSprite.Sprite.Y, IsBuffer: true}
	RegisterEntity(&Entity{Sprite: w.DstImg, Z: 2})

	noEraseDstImg := ebiten.NewImage(imgSizeWidth, imgSizeHeight)
	w.NoEraseDst = &sprite.Sprite{Img: noEraseDstImg, X: w.UiSprite.Sprite.X, Y: w.UiSprite.Sprite.Y, IsBuffer: true}
	RegisterEntity(&Entity{Sprite: w.NoEraseDst, Z: 2})

	w.timers = make(map[string]*util.Timer)
	w.timers["AllTasksCompletedBuffer"] = util.NewTimer(1)
	w.timers["EraseAnimationCompleted"] = util.NewTimer(2)
	w.timers["FirstTask"] = util.NewTimer(2.5)
	w.timers["UnderLine"] = util.NewTimer(1.5)
	w.timers["DayOver"] = util.NewTimer(2.5)
	w.timers["Later"] = util.NewTimer(2.0)
	w.timers["Publish"] = util.NewTimer(0.5)

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
	w.graphicManager.Timers[graphics.DeInit] = util.NewTimer(1)
	w.graphicManager.Timers[graphics.DeInit].TimerUpdater = graphics.DeInitGraphicTimerUpdater

	w.graphicManager.Tag = "whiteBoard"
	w.YSelectOptionsForCrossOut = []int{0, 1, 2, 3}
	w.occupied = make(map[string]float64)

	markerCrossOutTexture, err := util.LoadImageAssetAsEbitenImage("textures/crossOutTexture")
	if err != nil {
		log.Fatal(err)
	}

	w.crossOutTexture = markerCrossOutTexture

	graphics.WhiteBoardGMSubs(w.graphicManager, eventHub)
}

func (w *WhiteBoardSprite) ResetImg() {
	w.DstImg.Img.Clear()
}

func (w *WhiteBoardSprite) Update() {

	sp := w.UiSprite.Sprite
	w.graphicManager.UpdateTimers()

	if w.state == Writing {
		if w.CheckDoneWriting() {
			w.state = 0
		}
	}

	if w.clickAfterWriting && w.state != Writing {
		AddClickme(w)
		w.clickAfterWriting = false
	}

	if sp.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && w.taskJustCompleted {
		if !w.windowOpen {
			w.tasksCompletedToday++
			w.taskJustCompleted = false
			ev2 := tasks.TaskCompleted{Task: *w.currentTask}
			w.UiSprite.EventHub.Publish(ev2)
			w.UiSprite.Sprite.DOptsUpdaterTag = ""
			w.UiSprite.Sprite.Shader = nil
		}

	}

	if w.allTasksCompleted {
		w.allTasksCompleted = false
		w.timers["AllTasksCompletedBuffer"].TurnOn()
	}

	if w.allTasksBufferDone && w.clickme {
		AddClickme(w)

	}

	if w.allTasksBufferDone && sp.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		w.allTasksBufferDone = false
		w.UiSprite.Sprite.DOptsUpdaterTag = ""
		w.UiSprite.Sprite.Shader = nil
		w.UiSprite.EventHub.Publish(tasks.AllTasksCompleted{})
	}

	if w.DstImg.Shader != nil {
		maxCounter, ok := w.DstImg.ShaderParams["MaxCounter"].(int)
		if ok {
			counter := w.DstImg.ShaderParams["Counter"].(int)
			if counter >= maxCounter {
				w.DstImg.Shader = nil
			}
		}
	}
	w.UpdateTimers()
}

func (w *WhiteBoardSprite) CheckDoneWriting() bool {
	allDone := true
	for _, tid := range w.tIds {
		tEnt, exists := GetEntity(tid)
		if !exists {
			continue
		}
		if tEnt.Graphic != nil && !tEnt.Graphic.FullyDrawn {
			allDone = false
		}
	}

	return allDone

}

func AddClickme(w *WhiteBoardSprite) {
	w.clickme = false
	w.UiSprite.Sprite.Shader = registry.ShaderMap[registry.Highlight]
	w.UiSprite.Sprite.DOptsUpdaterTag = stringConstants.Swirl
}

func (w *WhiteBoardSprite) UpdateTimers() {

	for key, timer := range w.timers {
		state := timer.Update()
		switch key {

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
				for _, id := range w.tIds {
					RemoveEntity(id)
				}
				w.DstImg.Img.Clear()
				timer.TurnOff()
				if !w.dayOver {
					w.appendTextToBestSpot(WriteToWhiteBoard{Msg: "All Done =)", PreferredPosition: "upperCenter"})
					w.timers["UnderLine"].TurnOn()
					w.timers["Later"].TurnOn()
				}
			}
		case "FirstTask":
			if state == util.Done {
				w.appendTextToOpenSlot(w.currentTask.Text)
				timer.TurnOff()
			}
		case "UnderLine":
			if state == util.Done {
				rect := image.Rect(60, 32, 200, 48)
				sParams := defaultCrossOutShaderParams()
				sParams["Speed"] = 5.0
				dParams := defaultCrossOutDrawOptParams()
				w.AddCrossOutGraphicEntity(20, 70, &rect, sParams, dParams)
				timer.TurnOff()
			}
		case "Later":
			if state == util.Done {
				if len(w.msgsToWriteLater) <= 0 {
					return
				}
				w.appendTextToBestSpot(w.msgsToWriteLater[0])
				if w.msgsToWriteLater[0].EventToPublish != nil {
					w.timers["Publish"].TurnOn()
					w.pubQueue = append(w.pubQueue, w.msgsToWriteLater[0].EventToPublish)
				}

				if len(w.msgsToWriteLater) > 0 {
					w.msgsToWriteLater = w.msgsToWriteLater[1:]
				}

				if len(w.msgsToWriteLater) > 0 {
					if w.msgsToWriteLater[0].EventDriven != nil {
						timer.TurnOff()
					}
				} else {
					timer.TurnOff()
				}

				if len(w.msgsToWriteLater) <= 0 {
					timer.TurnOff()
				}

			}
		case "Publish":
			if state == util.Done {
				w.UiSprite.EventHub.Publish(w.pubQueue[0])
				timer.TurnOff()
			}
		}
	}
}

func defaultCrossOutShaderParams() map[string]any {
	ShaderParams := make(map[string]any)
	ShaderParams["Speed"] = 15 + rand.Float64()*10
	ShaderParams["MaxCounter"] = 500
	ShaderParams["Counter"] = 0
	maxOp := rand.Float64()*0.1 + 0.9
	ShaderParams["MaxOpacity"] = maxOp
	return ShaderParams
}

func defaultCrossOutDrawOptParams() map[string]float64 {
	DOptsUpdaterParams := make(map[string]float64)
	DOptsUpdaterParams["degree"] = rand.Float64() * -0.01
	return DOptsUpdaterParams
}

func (w *WhiteBoardSprite) AddCrossOutGraphicEntity(yCoord, xCoord float32, graphicBounds *image.Rectangle, shaderParams map[string]any, drawOptParams map[string]float64) {

	w.UiSprite.EventHub.Publish(events.WritingToWhiteBoard{
		Msg: CrossOut,
	})

	ySelectOp := rand.Intn(len(w.YSelectOptionsForCrossOut) - 1)
	ySelect := w.YSelectOptionsForCrossOut[ySelectOp]
	lowerY := ySelect * 16
	upperY := (ySelect + 1) * 16
	SubRect := image.Rect(0, lowerY, w.crossOutTexture.Bounds().Dx(), upperY)
	if graphicBounds != nil {
		SubRect = *graphicBounds
	}
	singleCross := w.crossOutTexture.SubImage(SubRect).(*ebiten.Image)
	//localCoords!!
	crossOutSprite := &sprite.Sprite{Img: singleCross, X: xCoord, Y: yCoord}
	crossOutSprite.Shader = registry.ShaderMap["HandWriting"]
	crossOutSprite.BufferDst = w.DstImg.Img
	crossOutSprite.ShaderParams = shaderParams
	crossOutSprite.DOptsUpdaterParams = drawOptParams
	crossOutSprite.UpdateShaderParams = shaders.UpdateCounterOneShot

	id := RegisterEntity(&Entity{Sprite: crossOutSprite})
	w.tIds = append(w.tIds, id)
	w.YSelectOptionsForCrossOut = append(w.YSelectOptionsForCrossOut[:ySelectOp], w.YSelectOptionsForCrossOut[ySelectOp+1:]...)
	if len(w.YSelectOptionsForCrossOut) == 1 {
		w.YSelectOptionsForCrossOut = []int{0, 1, 2, 3}
	}
}

func (w *WhiteBoardSprite) Subscribe(hub *tasks.EventHub) {
	hub.Subscribe(tasks.TaskRequirementsCompleted{}, func(e tasks.Event) {
		fmt.Printf("task requirment event recieved")
		if w.state != Writing && w.state != NotReceivingTasks {
			AddClickme(w)
		} else {
			w.clickAfterWriting = true
		}
		w.taskJustCompleted = true
	})

	hub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
		log.Printf("new task recieved")
		ev := e.(tasks.TaskCreated)
		w.currentTask = ev.Task
		if w.tasksCompletedToday == 0 {
			w.timers["FirstTask"].TurnOn()
		} else {
			w.state = Writing
			w.appendTextToOpenSlot(ev.Task.Text)
		}
	})

	hub.Subscribe(tasks.TaskCompleted{}, func(e tasks.Event) {
		yCoord := float32(w.spacing * (w.tasksCompletedToday + 1))
		w.AddCrossOutGraphicEntity(yCoord-10, 0, nil, defaultCrossOutShaderParams(), defaultCrossOutDrawOptParams())
		w.checkAllTasksCompleted()

	})

	hub.Subscribe(tasks.AllTasksCompleted{}, func(e tasks.Event) {
		w.initErase("")
	})

	hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		w.initErase("secondBuffer")
		w.dayOver = true
	})

	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		ev := e.(events.NewDay)
		w.numberOfTasksToday = ev.NTasks
		w.tasksCompletedToday = 0
		w.dayOver = false
		w.day = ev.Day
		fmt.Printf("White board received %d tasks on day %d", ev.NTasks, ev.Day)
		if ev.Day == 1 {
			w.state = NotReceivingTasks
		}
		if ev.Day != 1 {
			w.WriteDailyTasksHeader()
		}
	})

	hub.Subscribe(WriteToWhiteBoard{}, func(event tasks.Event) {
		w.state = Writing
		ev := event.(WriteToWhiteBoard)
		if ev.Later {
			w.msgsToWriteLater = append(w.msgsToWriteLater, ev)
			if ev.EventDriven != nil {
				hub.Subscribe(ev.EventDriven, func(e tasks.Event) {
					w.timers["Later"].TurnOn()
				})
			}
		} else {
			w.appendTextToBestSpot(ev)
		}
	})
	hub.Subscribe(events.WindowOpened{}, func(event tasks.Event) {
		w.windowOpen = true
	})
	hub.Subscribe(events.WindowClosed{}, func(e tasks.Event) {
		w.windowOpen = false
		ev := e.(events.WindowClosed)
		if ev.Window == string(GrandpasJournal) && w.day == 1 {
			w.state = Writing
			w.WriteDailyTasksHeader()
		}

	})

}

func (w *WhiteBoardSprite) WriteDailyTasksHeader() {
	w.appendTextToBestSpot(WriteToWhiteBoard{Msg: dailyTasks, PreferredPosition: "upperCenter"})
	w.timers["UnderLine"].TurnOn()
}

func (w *WhiteBoardSprite) initErase(tag string) {
	w.DstImg.Shader = registry.ShaderMap["Erase"]
	w.DstImg.ShaderParams = make(map[string]any)
	w.DstImg.ShaderParams["Counter"] = 0
	w.DstImg.ShaderParams["MaxCounter"] = 121
	w.DstImg.UpdateShaderParams = shaders.UpdateCounterOneShot
	w.DstImg.ShaderTexture = w.UiSprite.Sprite.Img
	w.timers["EraseAnimationCompleted"].TurnOn()
	w.occupied = make(map[string]float64)

	if tag == "secondBuffer" {
		w.NoEraseDst.Shader = registry.ShaderMap["Erase"]
		w.NoEraseDst.ShaderParams = make(map[string]any)
		w.NoEraseDst.ShaderParams["Counter"] = 0
		w.NoEraseDst.ShaderParams["MaxCounter"] = 121
		w.NoEraseDst.UpdateShaderParams = shaders.UpdateCounterOneShot
	}

}

func (w *WhiteBoardSprite) checkAllTasksCompleted() {

	if w.numberOfTasksToday == w.tasksCompletedToday {
		w.allTasksCompleted = true
	}

}

func (w *WhiteBoardSprite) appendTextToOpenSlot(txt string) {
	w.UiSprite.EventHub.Publish(events.WritingToWhiteBoard{
		Msg: txt,
	})
	yinset := float64(w.spacing*(w.tasksCompletedToday+1)) + 10.0
	insets := [2]float64{float64(w.UiSprite.Sprite.X), float64(w.UiSprite.Sprite.Y)}
	cs := &ebiten.ColorScale{}
	cs.SetA(1.0)
	cs.SetR(0.0)
	cs.SetB(0.0)
	cs.SetG(0.0)
	ts := graphics.NewTextWithMarkerShader(txt, w.DstImg.Img, insets, *cs, yinset, 10)
	ent := &Entity{Graphic: ts}
	id := RegisterEntity(ent)
	w.tIds = append(w.tIds, id)

}

func (w *WhiteBoardSprite) appendTextToBestSpot(ev WriteToWhiteBoard) {
	w.UiSprite.EventHub.Publish(events.WritingToWhiteBoard{
		Msg: ev.Msg,
	})
	var x float64
	var y float64

	if strings.Contains(ev.Msg, "Happy") {
		w.UiSprite.EventHub.Publish(events.BedTime{})
	}
	width, height := util.MeasureText(ev.Msg, 16, "RockSalt_12")

	switch ev.PreferredPosition {

	case "upperCenter":
		x = float64(w.DstImg.Img.Bounds().Dx()/2) - width/2
		y = height
		w.occupied["upperCenter"] += 1
	case "bottomLeft":
		y = float64(w.DstImg.Img.Bounds().Dy() - 30)
		x = 12
		w.occupied["upperCenter"] += 1
	case "bottomRight":
		y = float64(w.DstImg.Img.Bounds().Dy() - 30)
		x = float64(w.DstImg.Img.Bounds().Dx()-30) - width
		w.occupied["upperCenter"] += 1
	case "center":
		currentMsgs := float64(len(w.msgsToWriteLater))
		x = float64(w.DstImg.Img.Bounds().Dx()/2) - width/2
		y = float64(w.DstImg.Img.Bounds().Dy()/2) - height/2 - (height+5)*currentMsgs
		w.occupied["center"] += 1
	}

	insets := [2]float64{
		float64(w.UiSprite.Sprite.X),
		float64(w.UiSprite.Sprite.Y),
	}
	cs := &ebiten.ColorScale{}
	cs.SetA(1.0)
	cs.SetR(0.0)
	cs.SetB(0.0)
	cs.SetG(0.0)
	dst := w.DstImg
	if ev.NoErase {
		dst = w.NoEraseDst
	}
	ts := graphics.NewTextWithMarkerShader(
		ev.Msg,
		dst.Img,
		insets,
		*cs,
		y,
		x)
	ent := &Entity{Graphic: ts}
	id := RegisterEntity(ent)
	if !ev.NoErase {
		w.tIds = append(w.tIds, id)
	}
}
