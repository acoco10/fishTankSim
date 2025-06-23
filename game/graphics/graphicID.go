package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

var GraphicId int

var GraphMap = make(map[int]Graphic)

type Graphic interface {
	Draw(screen *ebiten.Image)
	Update()
	Coords() (float64, float64)
}

func AssignAndIncrement(graphic Graphic) int {
	currentGraphid := GraphicId
	GraphMap[GraphicId] = graphic
	GraphicId++
	return currentGraphid
}

func DeInitGraphicId(id int) {
	//no op if key doesnt exist
	log.Printf("deInitiating graphic with graphic id: %d", id)
	delete(GraphMap, id)

}

func DeInitGraphicBasedOnCoords(x, y float64) {
	//may need to tighten later (make sure input is the end coordinates if its animated graphic
	//typing could be added
	//can be faster with location-based partitions
	for id, graph := range GraphMap {
		x1, y1 := graph.Coords()
		dis := entities.DistanceFunc(float32(x1), float32(x), float32(y), float32(y1))
		if dis < 100 {
			DeInitGraphicId(int(id))
		}
	}
}

func DrawGraphics(screen *ebiten.Image) {
	for _, graph := range GraphMap {
		graph.Draw(screen)
	}
}

func NewFadeInTextGraphic(msg string, x, y float64) int {
	cs := ebiten.ColorScale{}
	cs.SetR(0.9)
	cs.SetB(0.9)
	cs.SetG(0.9)
	cs.SetA(1.0)
	id := NewGraphicText(msg, 24, x, y, false, cs, 0)
	return id
}
