package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
)

type GraphicManager struct {
	VlsGraphicQueue    []*VectorLineGraphic
	SpriteGraphicQueue []*SpriteGraphic
	TextGraphicQueue   map[int]*FadeInText
}

func NewGraphicManager(hub *tasks.EventHub, subs func(eventHub *tasks.EventHub, manager *GraphicManager)) *GraphicManager {
	gq := GraphicManager{}
	tMap := make(map[int]*FadeInText)
	gq.TextGraphicQueue = tMap
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

	for _, graph := range GraphMap {
		graph.Update()

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
