package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type GraphicManager struct {
	graphicQueue []*VectorLineGraphic
}

func NewGraphicManager(hub *events.EventHub) *GraphicManager {
	gq := GraphicManager{}
	gq.subscribe(hub)
	return &gq

}

func (g *GraphicManager) subscribe(hub *events.EventHub) {
	hub.Subscribe(events.TaskCompleted{}, func(e events.Event) {
		index := e.(events.TaskCompleted).Task.Index
		println("graphic manager recieving task:", index)
		x0 := 721.0
		y0 := 271.0 + index*25
		MaxX := x0 + 200.0
		y1 := y0 + 2.0
		var crossoutGraphic = NewVlS(float32(x0), float32(y0), float32(x0), float32(y1), float32(MaxX), color.Black)
		g.graphicQueue = append(g.graphicQueue, crossoutGraphic)
	})
}

func (g *GraphicManager) Update() {
	if len(g.graphicQueue) > 0 {
		for _, graphic := range g.graphicQueue {
			graphic.Update()
		}
	}
}

func (g *GraphicManager) Draw(screen *ebiten.Image) {
	if len(g.graphicQueue) > 0 {
		for _, graphic := range g.graphicQueue {
			graphic.Draw(screen)
		}
	}
}
