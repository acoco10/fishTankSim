package bank

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
)

type Bank struct {
	Amount float32
}

func NewBank(hub *tasks.EventHub, amount float32) *Bank {
	b := Bank{Amount: amount}
	return &b
}

func Subscribe(hub *tasks.EventHub, bank Bank) {

	hub.Subscribe(events.MoneySpent{}, func(e tasks.Event) {
		ev := e.(events.MoneySpent)
		bank.Amount -= ev.Amount
	})

	hub.Subscribe(events.MoneyAdded{}, func(e tasks.Event) {
		ev := e.(events.MoneySpent)
		bank.Amount += ev.Amount
	})
}
