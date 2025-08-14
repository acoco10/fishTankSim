package tutorial

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"reflect"
)

var typeRegistry = map[string]tasks.Event{
	"events.NewDay":         events.NewDay{},
	"tasks.TaskCompleted":   tasks.TaskCompleted{},
	"tasks.TaskCreated":     tasks.TaskCreated{},
	"events.UISpriteAction": events.UISpriteAction{},
	"events.UnFocus":        events.UnFocus{},
	"events.PHGuess":        events.PHGuess{},
}

const (
	tipx1 float64 = 960 / 4
	tipy1         = 10
)

type State uint8

const (
	notTriggered       = 0
	triggered    State = iota
	completed
	timeCompleted
)

type tip struct {
	state          State
	nextTip        *tip
	previousTip    *tip
	msg            string
	closeCondition string
	*util.Timer
}

type flags struct {
	FishClickable bool
}

type Manager struct {
	lastPublishedGraphicID int
	currCondition          string
	eventHub               *tasks.EventHub
	currenTip              tip
	tipMap                 map[string]*tip
	tipEventFilters        map[string]func(event tasks.Event) bool
	closeEventFilters      map[string]func(event tasks.Event) bool
	tipHead                *tip
	previousThread         *tip
	lastTip                bool
	flags                  flags
}

func InitData(m *Manager, hub *tasks.EventHub) {
	m.eventHub = hub

	eventMapper := make(map[string]*tip)
	tipMapFilter := make(map[string]func(event tasks.Event) bool)
	closeTipFilter := make(map[string]func(event tasks.Event) bool)

	condition1 := reflect.TypeOf(events.NewDay{}).String()

	lsTipMsgs1 := []string{
		"Welcome to Fish Fish Fish! Press Enter to advance tips",
		"press B to go back to previous tip",
		"The fish tank UI is defined by the objects on your desk",
		"Click the box of PH strips to take a PH reading of your tank",
	}

	lsTipMsgs1Timers := make(map[int]*util.Timer)
	lsTipMsgs1Timers[0] = util.NewTimer(2)
	lsTipMsgs1Timers[1] = util.NewTimer(2)
	lsTipMsgs1Timers[2] = util.NewTimer(2)

	//condition1Close := reflect.TypeOf(events.UISpriteAction{}).String()
	condition1Filter := func(event tasks.Event) bool {
		ev := event.(events.NewDay)
		return ev.Day == 1
	}

	tipMapFilter[condition1] = condition1Filter

	/*condition1CloseFilter := func(event tasks.Event) bool {
		ev := event.(events.UISpriteAction)
		return ev.UiSprite == "phreader" && ev.UiSpriteAction == "picked up"
	}*/
	eventMapper[condition1] = makeTipLinkedListFromStringList(lsTipMsgs1, lsTipMsgs1Timers, "")
	//closeTipFilter[condition1Close] = condition1CloseFilter

	lsTipMsgs1b := []string{
		"Press E to return any object to its place on the shelf",
	}
	lsTipMsgs1bTimers := make(map[int]*util.Timer)
	condition1b := reflect.TypeOf(events.PHGuess{}).String()
	closeCondition1b := reflect.TypeOf(events.UnFocus{}).String()

	eventMapper[condition1b] = makeTipLinkedListFromStringList(lsTipMsgs1b, lsTipMsgs1bTimers, closeCondition1b)

	condition2 := reflect.TypeOf(tasks.TaskCompleted{}).String()

	lsTipMsgs2 := []string{
		"Pick Your first structure for your fish tank",
		"Different structure effect the tank environment in different ways",
		"Position with cursor and click to place",
	}

	lsTipMsgs2Timers := make(map[int]*util.Timer)
	lsTipMsgs2Timers[0] = util.NewTimer(2)
	lsTipMsgs2Timers[1] = util.NewTimer(2)

	eventMapper[condition2] = makeTipLinkedListFromStringList(lsTipMsgs2, lsTipMsgs2Timers, "")

	condition2Filter := func(event tasks.Event) bool {
		ev := event.(tasks.TaskCompleted)
		return ev.Task.Text == "1. Take a ph reading of your tank"
	}

	tipMapFilter[condition2] = condition2Filter

	condition3 := reflect.TypeOf(tasks.TaskCreated{}).String()

	condition3Filter := func(event tasks.Event) bool {
		ev := event.(tasks.TaskCreated)
		return ev.Task.Text == "3. Feed your fish"
	}

	lstips3 := []string{
		"Click the Fish food to pick it up",
		"Click while holding it over the tank to dispense food",
		"press e to return to shelf",
		"Click on a fish to see their status and whether they're full"}

	eventMapper[condition3] = makeTipLinkedListFromStringList(lstips3, map[int]*util.Timer{}, "")

	tipMapFilter[condition3] = condition3Filter

	m.tipMap = eventMapper
	m.tipEventFilters = tipMapFilter
	m.closeEventFilters = closeTipFilter

	Subs(m.eventHub, m)
}

func makeTipLinkedListFromStringList(tipMsgs []string, timerMap map[int]*util.Timer, closeCondition string) *tip {
	var tipHead *tip
	var prevTip *tip
	for i, msg := range tipMsgs {
		newTip := &tip{state: notTriggered, previousTip: prevTip, nextTip: nil, msg: msg}
		if timerMap[i] != nil {
			newTip.Timer = timerMap[i]
		}
		if i == 0 {
			tipHead = newTip
		}

		if prevTip != nil {
			prevTip.nextTip = newTip // Link previous tip to current
		}
		if i == len(tipMsgs)-1 {
			newTip.closeCondition = closeCondition
		}

		prevTip = newTip // Update prevTip for next iteration
	}
	return tipHead
}

func (m *Manager) Update() {
	if m.tipHead != nil {
		if m.tipHead.Timer != nil {
			m.tipHead.Timer.TurnOn()
			state := m.tipHead.Timer.Update()
			if state == util.Done {
				m.tipHead.Timer.TurnOff()
				m.transitionToTip(m.tipHead.nextTip)
			}
		}
		m.handleNextTip()
		m.handlePreviousTip()
	}

}

func (m *Manager) handleNextTip() {
	if m.tipHead.nextTip == nil {
		if m.tipHead.closeCondition != "" { //m. CheckForClose(e)
		}
		m.lastTip = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if m.tipHead.nextTip != nil {
			m.transitionToTip(m.tipHead.nextTip)
		}
		if m.tipHead.nextTip == nil && m.lastTip {
			m.lastTip = false
			graphics.DeInitGraphicId(m.lastPublishedGraphicID)
		}

	}
}

func (m *Manager) handlePreviousTip() {
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		if m.tipHead.previousTip != nil {
			m.transitionToTip(m.tipHead.previousTip)
		}
	}
}

func (m *Manager) transitionToTip(newTip *tip) {
	graphics.DeInitGraphicId(m.lastPublishedGraphicID)
	m.tipHead = newTip
	m.lastPublishedGraphicID = graphics.NewFadeInTextGraphicSmall(m.tipHead.msg, tipx1, tipy1)
}

func Subs(hub *tasks.EventHub, m *Manager) {

	for key, _ := range m.tipMap {

		hub.Subscribe(typeRegistry[key], func(e tasks.Event) {
			m.CheckForTip(e)
		})
	}
}

func (m *Manager) CheckForClose(e tasks.Event) {
	eventType := reflect.TypeOf(e).String()
	if m.closeEventFilters[eventType] == nil {
		graphics.DeInitGraphicId(m.lastPublishedGraphicID)
	}
	if m.closeEventFilters[eventType](e) {
		graphics.DeInitGraphicId(m.lastPublishedGraphicID)
	}
}

func (m *Manager) CheckForTip(e tasks.Event) {

	eventType := reflect.TypeOf(e).String()
	println("In tutorial checking event:", eventType)

	if m.tipEventFilters[eventType] != nil {
		if m.tipEventFilters[eventType](e) {
			if m.lastPublishedGraphicID != 0 {
				graphics.DeInitGraphicId(m.lastPublishedGraphicID)
			}
			m.tipHead = m.tipMap[eventType]
			m.lastPublishedGraphicID = graphics.NewFadeInTextGraphicSmall(m.tipHead.msg, tipx1, tipy1)
			return
		} else {
			return
		}
	}
	if m.lastPublishedGraphicID != 0 {
		graphics.DeInitGraphicId(m.lastPublishedGraphicID)
	}
	m.currCondition = eventType
	m.tipHead = m.tipMap[eventType]
	m.lastPublishedGraphicID = graphics.NewFadeInTextGraphicSmall(m.tipHead.msg, tipx1, tipy1)
	m.previousThread = m.tipMap[eventType]

}
