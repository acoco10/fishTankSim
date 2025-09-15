package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

type GraphicManager struct {
	Tag                   string
	publishedGraphics     []int
	Params                map[string]any
	Timers                map[string]*util.Timer
	DstImage              *ebiten.Image
	GraphicsToBePublished []*SpriteGraphic
	texts                 []*TextWithShader
}

func (g *GraphicManager) UpdateTimers() {
	for key, timer := range g.Timers {
		switch key {
		case "DeInit":
			state := timer.Update()
			if state == util.Done {
				println("De-initiating graphics in:", g.Tag)
				DeInitGraphics(g.publishedGraphics)
				timer.TurnOff()
			}
		case "Trigger":
			state := timer.Update()
			if state == util.Done {
				if len(g.GraphicsToBePublished) != 0 {
					AddGraphic(g.GraphicsToBePublished[0])
					g.GraphicsToBePublished = g.GraphicsToBePublished[1:]
				}
				if len(g.GraphicsToBePublished) <= 0 {
					timer.TurnOff()
				}
			}
		}
	}
}

func DeInitGraphics(graphics []int) {
	for _, graphicId := range graphics {
		DeInitGraphicId(graphicId)
	}
}

func WhiteBoardGMSubs(manager *GraphicManager, hub *tasks.EventHub) {

	/*	hub.Subscribe(tasks.TaskCompleted{}, func(e tasks.Event) {

		index := manager.Params["Index"].(int)
		spacing := manager.Params["Spacing"].(int)

		x0 := float32(15)
		y0 := float32((index+1)*spacing) + 22.0
		MaxX := x0 + manager.Params["Width"].(float32) - 50 //arbitrary distance before end of whiteboard
		y1 := y0 + 2.0                                      //arbitrary height of triangle that forms the diagonal line
		id := NewVlS(x0, y0, x0, y1, MaxX, colornames.Navy, manager.DstImage)
		manager.publishedGraphics = append(manager.publishedGraphics, id)
		index++
		manager.Params["Index"] = index
		//manager.Timers["DeInit"].TurnOn()

	})*/

	hub.Subscribe(tasks.AllTasksCompleted{}, func(e tasks.Event) {
		AddClothGraphiC(manager)
		manager.Timers["DeInit"].TurnOn()
	})

	hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		AddClothGraphiC(manager)
		manager.Timers["DeInit"].TurnOn()
		manager.Params["Index"] = 0
	})

}

func AddClothGraphiC(manager *GraphicManager) {
	coord, ok := manager.Params["Coordinates"].(image.Point)
	if !ok {
		log.Fatal("White board graphic manager received invalid coordinates or no coordinates")
	}

	cg, err := LoadClothGraphic([2]float32{float32(coord.X), float32(coord.Y)})
	if err != nil {
		log.Fatal("error loading cloth Graphic\n", err)
	}

	id := AddGraphic(cg)
	manager.publishedGraphics = append(manager.publishedGraphics, id)
}

func (gm *GraphicManager) Update() {
	gm.UpdateTimers()
}
