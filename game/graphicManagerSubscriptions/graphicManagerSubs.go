package graphicManagerSubscriptions

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
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
		y0 := wbY + 33 + float32(index*20)
		MaxX := x0 + 200.0
		y1 := y0 + 2.0
		crossoutGraphic := graphics.NewVlS(x0, y0, x0, y1, MaxX, color.Black)
		manager.QueueGraphic(crossoutGraphic)
	})

	hub.Subscribe(tasks.AllTasksCompleted{}, func(e tasks.Event) {
		cg, err := loader.LoadClothGraphic()
		if err != nil {
			log.Fatal("error loading cloth Graphic\n", err)
		}
		if manager != nil {
			manager.ResetVls()
			manager.QueueGraphic(cg)
		}
	})

	hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		cg, err := loader.LoadClothGraphic()
		if err != nil {
			log.Fatal("error loading cloth Graphic\n", err)
		}
		if manager != nil {
			manager.ResetVls()
			manager.QueueGraphic(cg)
		}
	})
}

func ScreenGMSubs(hub *tasks.EventHub, manager *graphics.GraphicManager) {

	hub.Subscribe(events.ClickMeGraphicEvent{}, func(e tasks.Event) {
		ev := e.(events.ClickMeGraphicEvent)
		cs := ebiten.ColorScale{}
		cs.SetR(0.1)
		cs.SetB(0.2)
		cs.SetG(1.0)
		cs.SetA(1.0)
		graphics.NewGraphicText("Click Me", 24, ev.X, ev.Y, true, cs, ev.SpriteWidth)

	})

	hub.Subscribe(events.TurnOffGraphic{}, func(e tasks.Event) {
		ev := e.(events.TurnOffGraphic)
		graphics.DeInitGraphicBasedOnCoords(ev.X, ev.Y)
	})

	hub.Subscribe(events.FadeInTextEvent{}, func(e tasks.Event) {
	})
}
