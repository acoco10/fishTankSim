package graphics

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

const (
	DeInit  = "deInit"
	Trigger = "trigger"
	Trig1   = "trig1"
	Trig2   = "trig2"
)

type GmState uint8

const (
	InProgress GmState = iota
	Finished
)

type GraphicManager struct {
	GmState
	Tag                       string
	UpdateableText            *string
	UpdateableFloat           float64
	UpdateableFloatStop       float64
	publishedGraphics         []int
	Params                    map[string]any
	Timers                    map[string]*util.Timer
	DstImage                  *ebiten.Image
	TextGraphicsToBePublished []*FadeInText
	GraphicsToBePublished     []*SpriteGraphic
	texts                     []*TextWithShader
	TextPosition              []float64
	Strings                   []string
	textID                    int
	FinishedFunc              func(any) //flexible cleanup func based on gm state
}

func DeInitGraphicTimerUpdater(timer *util.Timer, graphicManager any) {
	gm, ok := graphicManager.(*GraphicManager)
	if !ok {
		log.Fatal("graphic manager timer update func got non graphic manager")
	}

	state := timer.Update()
	if state == util.Done {
		println("De-initiating graphics in:", gm.Tag)
		DeInitGraphics(gm.publishedGraphics)
		timer.TurnOff()
	}
}

func CashUpdaterTimerUpdater(timer *util.Timer, graphicManager any) {
	gm, ok := graphicManager.(*GraphicManager)
	if !ok {
		log.Fatal("graphic manager timer update func got non graphic manager")
	}
	state := timer.Update()
	if state == util.Done {
		if gm.UpdateableFloat <= gm.UpdateableFloatStop {
			newString := fmt.Sprintf("$%0.2f", gm.UpdateableFloat)
			NewTextGraphic(newString, gm.TextPosition[0], gm.TextPosition[1], 10)
			gm.UpdateableFloat += 0.25
		} else {
			newString := fmt.Sprintf("$%0.2f", gm.UpdateableFloatStop)
			NewPulseGraphic(newString, gm.TextPosition[0], gm.TextPosition[1], 180)
			gm.GmState = Finished
			timer.TurnOff()
		}
	}
}

func TriggerTimerUpdater(timer *util.Timer, graphicManager any) {
	gm, ok := graphicManager.(*GraphicManager)
	if !ok {
		log.Fatal("graphic manager timer update func got non graphic manager")
	}
	state := timer.Update()
	if state == util.Done {
		if len(gm.Strings) != 0 {
			NewFadeInTextGraphicCentered(gm.Strings[0], 100)
			gm.Strings = gm.Strings[1:]
		}
		if len(gm.GraphicsToBePublished) != 0 {
			AddGraphic(gm.GraphicsToBePublished[0])
			gm.GraphicsToBePublished = gm.GraphicsToBePublished[1:]
		} else if len(gm.GraphicsToBePublished) <= 0 {
			if gm.UpdateableFloat != 0 {
				newString := fmt.Sprintf("$%0.2f", gm.UpdateableFloat)
				NewFadeInTextGraphic(newString, gm.TextPosition[0], gm.TextPosition[1], 2)
				secondTime := gm.Params[Trig2].(string)
				gm.Timers[secondTime].TurnOn()
				timer.TurnOff()
			}
		}
	}
}

func (g *GraphicManager) UpdateTimers() {
	for _, timer := range g.Timers {
		if timer.TimerUpdater != nil {
			timer.TimerUpdater(timer, g)
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
		manager.Timers[DeInit].TurnOn()
	})

	hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		AddClothGraphiC(manager)
		manager.Timers[DeInit].TurnOn()
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
