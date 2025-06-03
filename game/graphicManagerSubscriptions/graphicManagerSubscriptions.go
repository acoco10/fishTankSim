package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"image/color"
)

func WhiteBoardGMSubs(hub *events.EventHub, manager *GraphicManager) {
	hub.Subscribe(events.TaskCompleted{}, func(e events.Event) {
		index := e.(events.TaskCompleted).Task.Index
		x0 := 721.0
		y0 := 271.0 + index*25
		MaxX := x0 + 200.0
		y1 := y0 + 2.0
		var crossoutGraphic = NewVlS(float32(x0), float32(y0), float32(x0), float32(y1), float32(MaxX), color.Black)
		manager.vlsGraphicQueue = append(manager.vlsGraphicQueue, crossoutGraphic)
	})

	hub.Subscribe(events.AllTasksCompleted{}, func(e events.Event) {

	})
}
