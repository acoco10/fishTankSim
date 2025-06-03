package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/hajimehoshi/ebiten/v2"
)

type GraphicManager struct {
	VlsGraphicQueue    []*VectorLineGraphic
	SpriteGraphicQueue []*SpriteGraphic
}

func NewGraphicManager(hub *events.EventHub, subs func(eventHub *events.EventHub, manager *GraphicManager)) *GraphicManager {
	gq := GraphicManager{}
	subs(hub, &gq)
	return &gq
}

func (g *GraphicManager) QueueGraphic(graphic any) {
	if spriteGraphic, ok := graphic.(*SpriteGraphic); ok {
		g.SpriteGraphicQueue = append(g.SpriteGraphicQueue, spriteGraphic)
	}

	if vlsGraphic, ok := graphic.(*VectorLineGraphic); ok {
		g.VlsGraphicQueue = append(g.VlsGraphicQueue, vlsGraphic)
	}
}

func (g *GraphicManager) Update() {
	if len(g.VlsGraphicQueue) > 0 {
		for _, graphic := range g.VlsGraphicQueue {
			graphic.Update()
		}
	}

	if len(g.SpriteGraphicQueue) > 0 {
		for i, graphic := range g.SpriteGraphicQueue {
			if graphic.complete {
				g.SpriteGraphicQueue = append(g.SpriteGraphicQueue[0:i], g.SpriteGraphicQueue[i+1:]...)
			}
			graphic.Update()
		}
	}
}

func (g *GraphicManager) ResetVls() {
	g.VlsGraphicQueue = []*VectorLineGraphic{}
}

func (g *GraphicManager) Draw(screen *ebiten.Image) {
	if len(g.VlsGraphicQueue) > 0 {
		for _, graphic := range g.VlsGraphicQueue {
			graphic.Draw(screen)
		}
	}

	if len(g.SpriteGraphicQueue) > 0 {
		for _, graphic := range g.SpriteGraphicQueue {
			graphic.Draw(screen)
		}
	}
}
