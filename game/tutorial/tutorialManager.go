package tutorial

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"reflect"
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
}

type Manager struct {
	lastPublishedGraphicID int
	eventHub               *tasks.EventHub
	currenTip              tip
	tipMap                 map[string]*tip
	tipHead                *tip
	previousThread         *tip
}

func InitData(m *Manager, hub *tasks.EventHub) {
	m.eventHub = hub

	eventMapper := make(map[string]*tip)
	condition1 := reflect.TypeOf(events.NewDay{}).String()

	lsTipMsgs1 := []string{
		"Press Enter to advance tips, press B to go back to previous tip",
		"Click to pick up fish food",
		"Don't feed your fish too much",
		"Press E to return Fish food to the shelf",
	}

	var tipHead *tip
	//var firstTip *tip

	for i, msg := range lsTipMsgs1 {
		newTip := &tip{state: notTriggered, previousTip: nil, nextTip: nil, msg: msg}
		if i == 0 {
			//firstTip = newTip
			eventMapper[condition1] = newTip
		}

		if tipHead != nil {
			tipHead.previousTip = tipHead
			tipHead = eventMapper[condition1]
		}

	}

	m.tipMap = eventMapper
	Subs(m.eventHub, m)

}

func (m *Manager) Update() {

	if m.tipHead != nil {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if m.tipHead.nextTip != nil {
				graphics.DeInitGraphicId(m.lastPublishedGraphicID)
				m.tipHead = m.tipHead.nextTip
				x := 500.0
				y := 100.0
				m.lastPublishedGraphicID = graphics.NewFadeInTextGraphic(m.tipHead.msg, x, y)
			}
			if m.tipHead.nextTip == nil {
				graphics.DeInitGraphicId(m.lastPublishedGraphicID)
				m.tipHead = m.previousThread
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyB) {
			if m.tipHead.previousTip != nil {
				graphics.DeInitGraphicId(m.lastPublishedGraphicID)
				m.tipHead = m.tipHead.previousTip
				x := 500.0
				y := 100.0
				m.lastPublishedGraphicID = graphics.NewFadeInTextGraphic(m.tipHead.msg, x, y)
			}
		}
	}
}

func Subs(hub *tasks.EventHub, m *Manager) {

	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		if m.tipHead == nil {
			eventType := reflect.TypeOf(e).String()
			m.tipHead = m.tipMap[eventType]
			m.previousThread = m.tipMap[eventType]
			x := 500.0
			y := 100.0
			// get the graphic id to de init the tip when its done
			log.Printf("initiating graphic for tip: %s", m.tipHead.msg)
			m.lastPublishedGraphicID = graphics.NewFadeInTextGraphic(m.tipHead.msg, x, y)
		}

	})

}
