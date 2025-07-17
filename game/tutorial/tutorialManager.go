package tutorial

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"reflect"
)

const (
	tipx1 float64 = 500
	tipy1         = 100
)

type State uint8

const (
	notTriggered       = 0
	triggered    State = iota
	completed
	timeCompleted
)

type tip struct {
	state       State
	nextTip     *tip
	previousTip *tip
	msg         string
	*entities.Timer
}

type flags struct {
	FishClickable bool
}

type Manager struct {
	lastPublishedGraphicID int
	eventHub               *tasks.EventHub
	currenTip              tip
	tipMap                 map[string]*tip
	tipEventFilters        map[string]func(event tasks.Event) bool
	tipHead                *tip
	previousThread         *tip
	flags                  flags
}

func InitData(m *Manager, hub *tasks.EventHub) {
	m.eventHub = hub

	eventMapper := make(map[string]*tip)
	tipMapFilter := make(map[string]func(event tasks.Event) bool)

	condition1 := reflect.TypeOf(events.NewDay{}).String()

	lsTipMsgs1 := []string{
		"Press Enter to advance tips, press B to go back to previous tip",
		"The fish tank UI is defined by the objects on you desk",
		"Click the box of PH strips to take a PH reading of your tank",
		"Don't feed your fish too much",
		"Press E to return Fish food to the shelf",
	}
	lsTipMsgs1Timers := make(map[int]*entities.Timer)

	eventMapper[condition1] = makeTipLinkedListFromStringList(lsTipMsgs1, lsTipMsgs1Timers)

	condition2 := reflect.TypeOf(tasks.TaskCompleted{}).String()

	lsTipMsgs2 := []string{
		"Pick Your first structure for your fish tank",
		"Different structure effect the tank environment in different ways",
		"Position with map and click to place",
	}
	lsTipMsgs2Timers := make(map[int]*entities.Timer)
	lsTipMsgs2Timers[1] = entities.NewTimer(4)

	eventMapper[condition2] = makeTipLinkedListFromStringList(lsTipMsgs2, lsTipMsgs2Timers)

	condition2Filter := func(event tasks.Event) bool {
		ev := event.(tasks.TaskCompleted)
		return ev.Task.Text == "1. Take a ph reading of your tank"
	}

	tipMapFilter[condition2] = condition2Filter

	contidion3 := reflect.TypeOf(tasks.TaskCreated{}).String()

	contidion3Filter := func(event tasks.Event) bool {
		ev := event.(tasks.TaskCreated)
		return ev.Task.Text == "3. Feed your fish"
	}

	lstips3 := []string{"Click the Fish food to pick it up",
		"Click while holding it over the tank to dispense food",
		"press e to return to shelf",
		"Click on a fish to see their status and whether they're full"}

	eventMapper[contidion3] = makeTipLinkedListFromStringList(lstips3, map[int]*entities.Timer{})

	tipMapFilter[contidion3] = contidion3Filter

	m.tipMap = eventMapper
	m.tipEventFilters = tipMapFilter

	Subs(m.eventHub, m)
}

func makeTipLinkedListFromStringList(tipMsgs []string, timerMap map[int]*entities.Timer) *tip {
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

		prevTip = newTip // Update prevTip for next iteration
	}
	return tipHead
}

func (m *Manager) Update() {
	if m.tipHead != nil {
		if m.tipHead.Timer != nil {
			m.tipHead.Timer.TurnOn()
			state := m.tipHead.Timer.Update()
			if state == entities.Done {
				m.tipHead.Timer.TurnOff()
				m.transitionToTip(m.tipHead.nextTip)
			}
		}
		m.handleNextTip()
		m.handlePreviousTip()
	}
}

func (m *Manager) handleNextTip() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if m.tipHead.nextTip != nil {
			m.transitionToTip(m.tipHead.nextTip)
		}
		if m.tipHead.nextTip == nil {
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
	m.lastPublishedGraphicID = graphics.NewFadeInTextGraphic(m.tipHead.msg, tipx1, tipy1)
}

func Subs(hub *tasks.EventHub, m *Manager) {

	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		m.CheckForTip(e)
	})
	hub.Subscribe(tasks.TaskCompleted{}, func(e tasks.Event) {
		m.CheckForTip(e)
	})
	hub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
		m.CheckForTip(e)
	})
}

func (m *Manager) CheckForTip(e tasks.Event) {
	if m.tipHead != nil {
		graphics.DeInitGraphicId(m.lastPublishedGraphicID)
	}

	eventType := reflect.TypeOf(e).String()
	println("In tutorial checking event:", eventType)

	if m.tipEventFilters[eventType] != nil {
		if m.tipEventFilters[eventType](e) {
			m.tipHead = m.tipMap[eventType]
			m.lastPublishedGraphicID = graphics.NewFadeInTextGraphic(m.tipHead.msg, tipx1, tipy1)
			return
		}
		println("In tutorial checking event:", eventType, "task doesnt filter")
		return
	}

	m.tipHead = m.tipMap[eventType]
	m.lastPublishedGraphicID = graphics.NewFadeInTextGraphic(m.tipHead.msg, tipx1, tipy1)
	m.previousThread = m.tipMap[eventType]
}
