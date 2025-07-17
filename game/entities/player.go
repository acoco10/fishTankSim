package entities

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
)

type Player struct {
	Money            float64
	EventHub         *tasks.EventHub
	ActionsTaken     int
	ActionsAvailable int
}

func (p *Player) AddMoney(money float64) {
	p.Money += money
}

func (p *Player) SpendMoney(money float64) {
	p.Money -= money
}

func (p *Player) Subscribe() {
	p.EventHub.Subscribe(events.MoneyAdded{}, func(e tasks.Event) {
		ev := e.(events.MoneyAdded)
		p.AddMoney(ev.Amount)
	})

	p.EventHub.Subscribe(events.MoneySpent{}, func(e tasks.Event) {
		ev := e.(events.MoneySpent)
		if ev.Amount < p.Money {
			p.SpendMoney(ev.Amount)
			ev2 := events.PurchaseSuccessful{}
			p.EventHub.Publish(ev2)
		} else {
			ev2 := events.InsufficientFunds{}
			p.EventHub.Publish(ev2)
		}

	})

}
