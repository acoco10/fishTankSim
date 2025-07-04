package system

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"math/rand"
)

type Environment struct {
	Temperature int
}

func (env *Environment) Subscribe(hub *tasks.EventHub) {
	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		env.Temperature = rand.Intn(10) + 67
	})
}
