package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/hajimehoshi/ebiten/v2"
)

type ScaledType uint8

const (
	NormalScaled ScaledType = iota
	UnScaled
	ScaledButTopLevel //weird ass way to draw scale graphics above text drawn unscaled like eraser
)

var GraphicId = 1

var GraphMap = make(map[int]Graphic)

type Graphic interface {
	Draw(screen *ebiten.Image)
	Update()
	Scaled() ScaledType
	//AutoDeinit() bool
}

func AssignAndIncrement(graphic Graphic) int {
	currentGraphid := GraphicId
	GraphMap[GraphicId] = graphic
	GraphicId++
	return currentGraphid
}

func DeInitGraphicId(id int) {
	//no op if key doesnt exist
	delete(GraphMap, id)
}

func DeInitAllGraphics() {
	GraphMap = make(map[int]Graphic)
}

func DrawScaledGraphics(screen *ebiten.Image) {
	for _, graph := range GraphMap {
		if graph.Scaled() == NormalScaled {
			graph.Draw(screen)
		}
	}
}

func DrawUnScaledGraphics(screen *ebiten.Image) {
	for _, graph := range GraphMap {
		if graph.Scaled() == UnScaled {
			graph.Draw(screen)
		}
	}
}

func DrawScaledGraphicsOnMainScreen(screen *ebiten.Image) {
	for _, graph := range GraphMap {
		if graph.Scaled() == ScaledButTopLevel {
			graph.Draw(screen)
		}
	}
}

func UpdateGraphics() {
	for _, graph := range GraphMap {
		graph.Update()
	}
}

func NewFadeInTextGraphicCentered(msg string, lifetime int) int {
	id := NewFadeInTextGraphic(msg, float64(registry.Config.ScreenWidth/2), float64(registry.Config.ScreenHeight/3), lifetime)
	return id
}

func NewTextGraphic(msg string, x, y float64, lifeTime int) int {
	cs := ebiten.ColorScale{}
	cs.SetR(0.9)
	cs.SetB(0.9)
	cs.SetG(0.9)
	cs.SetA(1.0)

	/*x = x * registry.Config.ResolutionScalingF
	y = (y + float64(registry.Config.YOffset)) * registry.Config.ResolutionScalingF*/
	id := newOutlineGraphicText(&msg, x, y, false, cs, 0, false, lifeTime)
	return id
}

func NewPulseGraphic(msg string, x, y float64, lifeTime int) int {
	cs := ebiten.ColorScale{}
	cs.SetR(0.9)
	cs.SetB(0.9)
	cs.SetG(0.9)
	cs.SetA(1.0)

	/*x = x * registry.Config.ResolutionScalingF
	y = (y + float64(registry.Config.YOffset)) * registry.Config.ResolutionScalingF*/
	id := newOutlineGraphicText(&msg, x, y, true, cs, 0, true, lifeTime)
	return id
}

func NewFadeInTextGraphic(msg string, x, y float64, lifeTime int) int {
	cs := ebiten.ColorScale{}
	cs.SetR(0.9)
	cs.SetB(0.9)
	cs.SetG(0.9)
	cs.SetA(1.0)

	/*x = x * registry.Config.ResolutionScalingF
	y = (y + float64(registry.Config.YOffset)) * registry.Config.ResolutionScalingF*/
	id := newOutlineGraphicText(&msg, x, y, false, cs, 0, true, lifeTime)
	return id
}

func NewFadeInTextGraphicSmall(msg string, x, y float64, lifetime int) int {
	cs := ebiten.ColorScale{}
	cs.SetR(0.9)
	cs.SetB(0.9)
	cs.SetG(0.9)
	cs.SetA(1.0)

	id := newOutlineGraphicText(&msg, x, y, false, cs, 0, true, lifetime)
	return id
}

func NewUpdateAbleTextGraphic(msg *string, x, y float64, lifetime int) int {
	cs := ebiten.ColorScale{}
	cs.SetR(0.9)
	cs.SetB(0.9)
	cs.SetG(0.9)
	cs.SetA(1.0)

	id := newOutlineGraphicText(msg, x, y, false, cs, 0, true, lifetime)
	return id
}

func AddGraphic(graphic Graphic) int {
	id := AssignAndIncrement(graphic)
	return id
}
