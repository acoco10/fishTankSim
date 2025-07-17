package graphicManagerSubscriptions

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"golang.org/x/image/colornames"
	"log"
)

func WhiteBoardGMSubs(hub *tasks.EventHub, manager *graphics.GraphicManager) {
	var wbX float32
	var wbY float32

	hub.Subscribe(events.UISpriteLayedOut{}, func(e tasks.Event) {
		ev := e.(events.UISpriteLayedOut)
		println("uisprite label =", ev.Label)
		if ev.Label == "whiteBoard" {
			wbX = ev.X
			wbY = ev.Y
		}
	})

	hub.Subscribe(tasks.TaskCompleted{}, func(e tasks.Event) {
		index := e.(tasks.TaskCompleted).Slot - 1
		x0 := wbX + 10
		y0 := wbY + 55 + float32(index*20)
		MaxX := x0 + 200.0
		y1 := y0 + 2.0
		crossoutGraphic := graphics.NewVlS(x0, y0, x0, y1, MaxX, colornames.Orangered)
		manager.QueueGraphic(crossoutGraphic)
	})

	hub.Subscribe(tasks.AllTasksCompleted{}, func(e tasks.Event) {
		cg, err := loader.LoadClothGraphic([2]float32{wbX, wbY})
		if err != nil {
			log.Fatal("error loading cloth Graphic\n", err)
		}
		if manager != nil {
			manager.ResetVls()
			manager.QueueGraphic(cg)
		}
	})

	hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		cg, err := loader.LoadClothGraphic([2]float32{wbX, wbY})
		if err != nil {
			log.Fatal("error loading cloth Graphic\n", err)
		}
		if manager != nil {
			manager.ResetVls()
			manager.QueueGraphic(cg)
		}
	})
}
