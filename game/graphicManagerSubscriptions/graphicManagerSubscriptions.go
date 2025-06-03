package graphicManagerSubscriptions

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/loaders"
	"image/color"
	"log"
)

func WhiteBoardGMSubs(hub *events.EventHub, manager *graphics.GraphicManager) {
	hub.Subscribe(events.TaskCompleted{}, func(e events.Event) {
		index := e.(events.TaskCompleted).Task.Index
		x0 := 721.0
		y0 := 271.0 + index*25
		MaxX := x0 + 200.0
		y1 := y0 + 2.0
		crossoutGraphic := graphics.NewVlS(float32(x0), float32(y0), float32(x0), float32(y1), float32(MaxX), color.Black)
		manager.QueueGraphic(crossoutGraphic)
	})

	hub.Subscribe(events.AllTasksCompleted{}, func(e events.Event) {

		cg, err := loaders.LoadClothGraphic()
		if err != nil {
			log.Fatal("error loading cloth Graphic\n", err)
		}
		manager.ResetVls()
		manager.QueueGraphic(cg)
	})
}
