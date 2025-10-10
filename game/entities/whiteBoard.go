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
	"golang.org/x/image/colornames"
	"image"
	"log"
)

const (
	//whiteBoardTimers
	clickableAfterTimer = iota
)

const (
	dailyTasks = "Daily Tasks"
	allDone    = "All Done =)"
	CrossOut   = "crossOut"

	//state
	DisabledWB = iota
	IdleWB
	WritingTask
	WritingRequest
	Erasing
	CrossingOut
	ClickTransition
)

const (
	//task state
	NotReceivingTasks = iota
	ReadyToClick
	ReceivingTasks
	AllTaskCompleted
	ClickToCompleteAfterWriting
)

const (
	//day state
	JustStarted = iota
	DayTime
	NightTime
	DayOver
)

const (
	//positions
	Center = iota
	UpperCenter
	LowerCenter

	UpperLeft
	UpperRight

	LowerLeft
	LowerRight
)

type WhiteBoardSprite struct {
	currentTask      *tasks.Task
	textBeingWritten *graphics.TextWithShader
	DstImg           *sprite.Sprite
	NoEraseDst       *sprite.Sprite
	graphicManager   *graphics.GraphicManager

	wbState             uint8
	taskState           uint8
	dayState            uint8
	allTasksCompleted   bool
	timers              [7]*util.Timer
	numberOfTasksToday  int
	tasksCompletedToday int
	UiSprite            *Entity
	lastTaskID          int
	spacing             int
	header              bool

	textGraphicEntityIDs    map[uint32]struct{}
	crossOutGraphicEntityID map[uint32]struct{}
	noEraseTextGraphicEntID map[uint32]struct{}
	msgsToWriteLater        []WriteToWhiteBoard
	occupied                map[uint8]int
	crossOutTexture         *ebiten.Image

	drawingEnt uint32

	pubQueue                       []tasks.Event
	taskCreatedEventQueue          []tasks.TaskCreated
	taskRequirementsCompletedEvent *tasks.TaskRequirementsCompleted
	writeToWhiteBoardQueue         []WriteToWhiteBoard
	eraseRequest                   *EraseRequest

	unlockEvent     tasks.Event
	unlockEventType string

	day                       int
	YSelectOptionsForCrossOut []int
	wbStateMachine            *WhiteBoardStateMachine
}

func (wb *WhiteBoardSprite) Init() {
	wb.spacing = 20
	imgSizeWidth := wb.UiSprite.Sprite.Img.Bounds().Dx()
	imgSizeHeight := wb.UiSprite.Sprite.Img.Bounds().Dy()
	dstImg := ebiten.NewImage(imgSizeWidth, imgSizeHeight)

	wb.DstImg = &sprite.Sprite{Img: dstImg, X: wb.UiSprite.Sprite.X, Y: wb.UiSprite.Sprite.Y, IsBuffer: true}
	RegisterEntity(&Entity{Sprite: wb.DstImg, Z: 2})

	noEraseDstImg := ebiten.NewImage(imgSizeWidth, imgSizeHeight)
	wb.NoEraseDst = &sprite.Sprite{Img: noEraseDstImg, X: wb.UiSprite.Sprite.X, Y: wb.UiSprite.Sprite.Y, IsBuffer: true}
	RegisterEntity(&Entity{Sprite: wb.NoEraseDst, Z: 2})

	wb.DstImg.ShaderParams = make(map[string]any)

	wb.graphicManager = &graphics.GraphicManager{}
	wb.graphicManager.Params = make(map[string]any)
	wb.graphicManager.GMParams.Coordinates = [2]float32{wb.UiSprite.Sprite.X, wb.UiSprite.Sprite.Y}
	wb.graphicManager.DstImage = wb.DstImg.Img

	wb.graphicManager.Timers = make(map[string]*util.Timer)
	wb.graphicManager.Timers[graphics.DeInit] = util.NewTimer(2)
	wb.graphicManager.Timers[graphics.DeInit].TimerUpdater = graphics.DeInitGraphicTimerUpdater

	wb.graphicManager.Tag = "whiteBoard"
	wb.YSelectOptionsForCrossOut = []int{0, 1, 2, 3}

	wb.occupied = make(map[uint8]int)
	wb.crossOutGraphicEntityID = make(map[uint32]struct{})
	wb.textGraphicEntityIDs = make(map[uint32]struct{})
	wb.noEraseTextGraphicEntID = make(map[uint32]struct{})

	markerCrossOutTexture, err := util.LoadImageAssetAsEbitenImage("textures/crossOutTexture")
	if err != nil {
		log.Fatal(err)
	}

	wb.timers[clickableAfterTimer] = util.NewTimer(0.5)

	wb.crossOutTexture = markerCrossOutTexture

	wb.wbStateMachine = initWBStateMachine()

}

func (wb *WhiteBoardSprite) ResetImgAndClearGraphics() {
	wb.DstImg.Img.Clear()
	for entId, _ := range wb.textGraphicEntityIDs {
		RemoveEntity(entId)
	}
	for entId, _ := range wb.crossOutGraphicEntityID {
		RemoveEntity(entId)
	}

	wb.DstImg.Shader = nil
}

type WhiteBoardUpdater struct {
	state             uint8
	updater           func(wb *WhiteBoardSprite, gs GameState)
	transitionOutFunc func(wb *WhiteBoardSprite)
	transitionToFunc  func(wb *WhiteBoardSprite)
	autoTransition    bool
	requireClickFlag  bool
}

type WhiteBoardStateMachine struct {
	State         uint8
	onClickState  uint8
	StateUpdaters [10]WhiteBoardUpdater
}

func (wsb *WhiteBoardStateMachine) Update(wb *WhiteBoardSprite, gs GameState) {
	gameStateProcessor(wb, gs) //early returns if window is open (things we want all states to consider)
	state := wsb.StateUpdaters[wsb.State]
	state.updater(wb, gs)
}

func (wsb *WhiteBoardStateMachine) Transition(wb *WhiteBoardSprite) {
	state := wsb.StateUpdaters[wsb.State]
	if state.transitionOutFunc != nil {
		state.transitionOutFunc(wb) //at transitio func
	}

	newState := wb.MapStateBasedOnEventQueues()

	newUpdater := wsb.StateUpdaters[newState]

	if newUpdater.requireClickFlag {
		if newState == Erasing {
			if !wb.eraseRequest.onClick {
				wsb.State = newState
				if newUpdater.transitionToFunc != nil {
					newUpdater.transitionToFunc(wb)
					return
				}
			}
		}

		wsb.onClickState = newState
		wsb.State = ClickTransition
		wsb.StateUpdaters[ClickTransition].transitionToFunc(wb)
	} else {
		wsb.State = newState
		if newUpdater.transitionToFunc != nil {
			newUpdater.transitionToFunc(wb)
		}
	}

}

func initWBStateMachine() *WhiteBoardStateMachine {

	wsb := &WhiteBoardStateMachine{}

	writingUpdate := WhiteBoardUpdater{
		state:             WritingTask,
		updater:           WriteHeaderMonitor,
		transitionToFunc:  writeTask,
		requireClickFlag:  false,
		transitionOutFunc: nil}

	writeRequestUpdater := WhiteBoardUpdater{
		state:             WritingRequest,
		updater:           WriteHeaderMonitor,
		transitionToFunc:  writeRequest,
		transitionOutFunc: nil,
		requireClickFlag:  false,
	}

	crossOutUpdate := WhiteBoardUpdater{
		state:             CrossingOut,
		updater:           WriteHeaderMonitor,
		transitionToFunc:  intoCrossOutTransition,
		requireClickFlag:  true,
		transitionOutFunc: nil}

	eraseUpdate := WhiteBoardUpdater{
		state:             Erasing,
		updater:           monitorErase,
		transitionToFunc:  intoErasingTransition,
		requireClickFlag:  true,
		transitionOutFunc: outOfErase}

	clickUpdate := WhiteBoardUpdater{
		state:             ClickTransition,
		updater:           OnClickUpdater,
		transitionToFunc:  OnClickTransition,
		transitionOutFunc: nil,
		requireClickFlag:  false,
	}

	idleUpdate := WhiteBoardUpdater{
		state:             IdleWB,
		updater:           WBIdleQueueMonitor,
		transitionToFunc:  nil,
		transitionOutFunc: nil,
		requireClickFlag:  false,
	}

	wsb.State = IdleWB
	wsb.StateUpdaters = [10]WhiteBoardUpdater{}
	wsb.StateUpdaters[IdleWB] = idleUpdate
	wsb.StateUpdaters[WritingTask] = writingUpdate
	wsb.StateUpdaters[WritingRequest] = writeRequestUpdater
	wsb.StateUpdaters[CrossingOut] = crossOutUpdate
	wsb.StateUpdaters[Erasing] = eraseUpdate
	wsb.StateUpdaters[ClickTransition] = clickUpdate

	return wsb
}

func WBIdleQueueMonitor(wb *WhiteBoardSprite, gs GameState) {
	if wb.MapStateBasedOnEventQueues() != IdleWB {
		wb.wbStateMachine.Transition(wb)
	}
}

func OnClickTransition(wb *WhiteBoardSprite) {
	wb.timers[clickableAfterTimer].TurnOn()
}

func OnClickUpdater(wb *WhiteBoardSprite, gs GameState) {
	if registry.ClickCheck() && wb.UiSprite.Sprite.SpriteHovered() && wb.timers[clickableAfterTimer].On == false {
		RemoveClickMe(wb.UiSprite)
		onClickState := wb.wbStateMachine.onClickState
		wb.wbStateMachine.State = onClickState
		wb.wbStateMachine.StateUpdaters[onClickState].transitionToFunc(wb)
	}
}

func DrawingMonitor(wb *WhiteBoardSprite, gs GameState) {
	if CheckDoneDrawing(wb.drawingEnt) {
		if wb.header {
			wb.header = false
			wb.underLineHeader()
		}
		wb.wbStateMachine.Transition(wb)
	}
}

func WriteHeaderMonitor(wb *WhiteBoardSprite, gs GameState) {
	if CheckDoneDrawing(wb.drawingEnt) {
		if wb.header {
			wb.header = false
			wb.underLineHeader()
		}
		wb.wbStateMachine.Transition(wb)
	}
}

func CrossOutMonitor(wb *WhiteBoardSprite, gs GameState) {
	CheckDoneDrawing(wb.drawingEnt)
}

func gameStateProcessor(wb *WhiteBoardSprite, gs GameState) {
	//non-state specific game state processing
	if gs.MouseFlags.WindowOpen {
		//don't update if UI window is open
		return
	}
}

func (wb *WhiteBoardSprite) MapStateBasedOnEventQueues() uint8 {
	//order of priority

	if wb.unlockEvent != nil {
		return IdleWB
	}

	if len(wb.writeToWhiteBoardQueue) > 0 {
		return WritingRequest
	}

	if len(wb.taskCreatedEventQueue) > 0 {
		return WritingTask
	}

	if wb.taskRequirementsCompletedEvent != nil {
		return CrossingOut
	}

	if wb.eraseRequest != nil {
		return Erasing
	}

	return IdleWB
}

func (wbs *WhiteBoardStateMachine) DisableUpdate(wb *WhiteBoardSprite, gs GameState) {
	if wb.unlockEvent == nil {
		wbs.Transition(wb)
	}
}

func (wb *WhiteBoardSprite) Update(gs GameState) {
	wb.UpdateTimers()
	wb.wbStateMachine.Update(wb, gs)
	wb.graphicManager.Update()
}

func writeTask(wb *WhiteBoardSprite) {
	wb.processQueuedWriteEvent()
}

func writeRequest(wb *WhiteBoardSprite) {
	wb.processQueuedWriteRequest()
}

func intoErasingTransition(wb *WhiteBoardSprite) {
	wb.initErase("")
	graphics.AddClothGraphic(wb.graphicManager)
	if wb.eraseRequest.time == NightTime {
		wb.WriteHeader(allDone)
		wb.header = true
	}
}

func monitorErase(wb *WhiteBoardSprite, gs GameState) {
	if wb.EraseFinished() {
		wb.wbStateMachine.Transition(wb)
	}
}

func outOfErase(wb *WhiteBoardSprite) {
	wb.ResetImgAndClearGraphics()
	wb.UiSprite.EventHub.Publish(WhiteBoardErased{When: wb.eraseRequest.time})
	wb.eraseRequest = nil
}

func intoCrossOutTransition(w *WhiteBoardSprite) {

	w.UiSprite.EventHub.Publish(tasks.TaskCompleted{Task: w.taskRequirementsCompletedEvent.Task})

	w.tasksCompletedToday += 1
	yCoord := float32(w.spacing * (w.tasksCompletedToday + 1))
	//defineAndRegisterCrossOutGraphicEntity
	w.AddCrossOutGraphicEntity(
		yCoord-10,
		0,
		nil,
		defaultCrossOutShaderParams(),
		defaultCrossOutDrawOptParams())
	w.taskRequirementsCompletedEvent = nil
}

func (wb *WhiteBoardSprite) EraseFinished() bool {
	return sprite.CheckSpriteWithShaderCounterFinished(wb.DstImg)
}

func (wb *WhiteBoardSprite) WriteHeader(msg string) {
	insets := getWBPreferredPositionInsets(msg, UpperCenter, wb.occupied, wb.DstImg.GetSpriteRect())
	wb.queueEvent(WriteToWhiteBoard{Msg: msg, PreferredPosition: UpperCenter, Insets: insets})

}

func (wb *WhiteBoardSprite) Subscriptions(hub *tasks.EventHub) {
	hub.Subscribe(tasks.TaskRequirementsCompleted{}, func(e tasks.Event) {
		wb.queueEvent(e)
	})

	hub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
		wb.queueEvent(e)
	})

	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		ev := e.(events.NewDay)
		fmt.Printf("White board received %d tasks on day %d", ev.NTasks, ev.Day)

		wb.wbState = IdleWB
		wb.tasksCompletedToday = 0
		wb.numberOfTasksToday = ev.NTasks
		wb.day = ev.Day

		wb.NoEraseDst.Shader = nil
		wb.NoEraseDst.Img.Clear()
		wb.header = true

		wb.WriteHeader(dailyTasks)

	})

	hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		wb.eraseRequest = &EraseRequest{onClick: false}
	})

	hub.Subscribe(WriteToWhiteBoard{}, func(e tasks.Event) {
		ev := e.(WriteToWhiteBoard)
		insets := getWBPreferredPositionInsets(ev.Msg, ev.PreferredPosition, wb.occupied, wb.DstImg.GetSpriteRect())
		ev.Insets = insets
		wb.writeToWhiteBoardQueue = append(wb.writeToWhiteBoardQueue, ev)
	})

	hub.Subscribe(DisableWhiteBoard{}, func(e tasks.Event) {
		ev := e.(DisableWhiteBoard)
		wb.unlockEvent = ev.UnLockEvent
		hub.Subscribe(ev.UnLockEvent, func(e tasks.Event) {
			if ev.Condition(e) {
				wb.unlockEvent = nil
			}
		})
	})

	hub.Subscribe(tasks.AllTasksCompleted{}, func(e tasks.Event) {
		wb.allTasksCompleted = true
		wb.eraseRequest = &EraseRequest{time: NightTime, onClick: true}
	})

}

func (wb *WhiteBoardSprite) queueEvent(event tasks.Event) {
	evTc, ok := event.(tasks.TaskCreated)
	if ok {
		wb.taskCreatedEventQueue = append(wb.taskCreatedEventQueue, evTc)
	}

	evReq, ok := event.(tasks.TaskRequirementsCompleted)
	if ok {
		wb.taskRequirementsCompletedEvent = &evReq
	}

	writeReq, ok := event.(WriteToWhiteBoard)
	if ok {
		wb.writeToWhiteBoardQueue = append(wb.writeToWhiteBoardQueue, writeReq)
	}

}

func (wb *WhiteBoardSprite) initErase(tag string) {
	wb.DstImg.Shader = registry.ShaderMap[registry.Erase]
	wb.DstImg.ShaderParams = make(map[string]any)
	wb.DstImg.ShaderParams["Counter"] = 0
	wb.DstImg.ShaderParams["MaxCounter"] = 100
	wb.DstImg.UpdateShaderParams = shaders.UpdateCounterOneShot
	wb.DstImg.ShaderTexture = wb.UiSprite.Sprite.Img
	wb.occupied = make(map[uint8]int)

	if tag == "secondBuffer" {
		wb.NoEraseDst.Shader = registry.ShaderMap["Erase"]
		wb.NoEraseDst.ShaderParams = make(map[string]any)
		wb.NoEraseDst.ShaderParams["Counter"] = 0
		wb.NoEraseDst.ShaderParams["MaxCounter"] = 100
		wb.NoEraseDst.UpdateShaderParams = shaders.UpdateCounterOneShot
		wb.NoEraseDst.ShaderTexture = wb.UiSprite.Sprite.Img
	}

}

func (wb *WhiteBoardSprite) checkAllTasksCompleted() bool {
	return wb.numberOfTasksToday == wb.tasksCompletedToday
}

func (wb *WhiteBoardSprite) WriteTextToOpenSlot(txt string) {
	//w.wbState = WritingTask no state in non state funcs plz
	wb.UiSprite.EventHub.Publish(events.WritingToWhiteBoard{
		Msg: txt,
	})
	yinset := float64(wb.spacing*(wb.tasksCompletedToday+1)) + 10.0
	insets := [2]float64{float64(wb.UiSprite.Sprite.X), float64(wb.UiSprite.Sprite.Y)}

	cs := util.ConvertRGBAtoEbitenCS(colornames.Crimson)
	ts := graphics.NewTextWithMarkerShader(txt, wb.DstImg.Img, insets, cs, 10, yinset)
	ent := &Entity{ShaderTextGraphic: ts}
	id := RegisterEntity(ent)
	wb.drawingEnt = id
	wb.textGraphicEntityIDs[id] = struct{}{}
}

func (wb *WhiteBoardSprite) appendTextToBestSpot(ev WriteToWhiteBoard) {
	//w.wbState = WritingTask
	wb.UiSprite.EventHub.Publish(events.WritingToWhiteBoard{
		Msg: ev.Msg,
	})

	insets := [2]float64{
		float64(wb.UiSprite.Sprite.X),
		float64(wb.UiSprite.Sprite.Y),
	}

	cs := util.ConvertRGBAtoEbitenCS(colornames.Crimson)

	dst := wb.DstImg
	if ev.NoErase {
		dst = wb.NoEraseDst
	}

	ts := graphics.NewTextWithMarkerShader(
		ev.Msg,
		dst.Img,
		insets,
		cs,
		ev.Insets[0],
		ev.Insets[1])

	ent := &Entity{ShaderTextGraphic: ts}
	id := RegisterEntity(ent)
	if ev.NoErase {
		wb.noEraseTextGraphicEntID[id] = struct{}{}
	} else {
		wb.drawingEnt = id
		wb.textGraphicEntityIDs[id] = struct{}{}
	}
}

func (wb *WhiteBoardSprite) UpdateTimers() {

	for key, timer := range wb.timers {
		if timer == nil {
			continue
		}
		state := timer.Update()
		if key == clickableAfterTimer && state == util.Done {
			timer.TurnOff()
			AddClickme(wb.UiSprite)
		}
	}
	return
}

func getWBPreferredPositionInsets(
	msg string,
	positionRequest uint8,
	occupied map[uint8]int,
	bounds image.Rectangle) [2]float64 {

	var x float64
	var y float64
	var lineSpacing = 5.0

	width, height := util.MeasureText(msg, 16, "RockSalt_12")

	switch positionRequest {

	case UpperCenter:
		x = float64(bounds.Dx()/2) - width/2
		y = height
		occupied[UpperCenter] += 1
		return [2]float64{x, y}
	case LowerLeft:
		y = float64(bounds.Dy() - 30)
		x = 12
		occupied[LowerLeft] += 1
		return [2]float64{x, y}
	case LowerRight:
		y = float64(bounds.Dy() - 30)
		x = float64(bounds.Dx()-30) - width
		occupied[LowerRight] += 1
		return [2]float64{x, y}
	case Center:
		x = float64(bounds.Dx()/2) - width/2
		y = float64(bounds.Dy()/2) - height/2 + ((lineSpacing + height) * float64(occupied[Center])) - 20
		occupied[Center] += 1
		return [2]float64{x, y}
	default:
		return [2]float64{0, 0}

	}

}
