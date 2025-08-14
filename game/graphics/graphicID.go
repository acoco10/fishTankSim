package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
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

func NewFadeInTextGraphic(msg string, x, y float64) int {
	cs := ebiten.ColorScale{}
	cs.SetR(0.9)
	cs.SetB(0.9)
	cs.SetG(0.9)
	cs.SetA(1.0)

	x = x * registry.Config.ResolutionScalingF
	y = (y + float64(registry.Config.YOffset)) * registry.Config.ResolutionScalingF
	id := NewOutlineGraphicText(&msg, 48, x, y, false, cs, 0, true)
	return id
}

func NewFadeInTextGraphicSmall(msg string, x, y float64) int {
	cs := ebiten.ColorScale{}
	cs.SetR(0.9)
	cs.SetB(0.9)
	cs.SetG(0.9)
	cs.SetA(1.0)

	x = x * registry.Config.ResolutionScalingF
	y = (y + float64(registry.Config.YOffset)) * registry.Config.ResolutionScalingF
	id := NewOutlineGraphicText(&msg, 24, x, y, false, cs, 0, true)
	return id
}

func NewUpdateAbleTextGraphic(msg *string, x, y float64) int {
	cs := ebiten.ColorScale{}
	cs.SetR(0.9)
	cs.SetB(0.9)
	cs.SetG(0.9)
	cs.SetA(1.0)
	x = x * registry.Config.ResolutionScalingF
	y = (y)*registry.Config.ResolutionScalingF + float64(registry.Config.YOffset)
	id := NewOutlineGraphicText(msg, 32, x, y, false, cs, 0, false)
	return id
}

func NewTravelingEffect(sprite *sprite.Sprite, x, y *float32) int {
	eff := &TaggedTravelingGraphic{sprite: sprite, x: x, y: y}
	GraphicId++
	GraphMap[GraphicId] = eff
	return GraphicId
}

type TaggedTravelingGraphic struct {
	sprite *sprite.Sprite
	x      *float32
	y      *float32
	id     int
}

func (t *TaggedTravelingGraphic) Update() {
	if t.x == nil {
		DeInitGraphicId(t.id)
	}
	t.sprite.X = *t.x
	t.sprite.Y = *t.y - float32(t.sprite.SpriteHeight)

	t.sprite.Update()
}

func (t *TaggedTravelingGraphic) Draw(screen *ebiten.Image) {
	dopts := &ebiten.DrawImageOptions{}
	dopts.GeoM.Translate(
		float64(t.sprite.X),
		float64(t.sprite.Y))
	t.sprite.UpdateOpts(dopts)
	t.sprite.Draw(screen)

}
func (t *TaggedTravelingGraphic) Scaled() ScaledType {
	return NormalScaled
}

func AddGraphic(graphic Graphic) int {
	id := AssignAndIncrement(graphic)
	return id
}

func AddHandwritingGraphic(txt string, buff *ebiten.Image, insets [2]float64, yInset float64, xInset float64) int {
	cs := &ebiten.ColorScale{}
	cs.SetA(1.0)
	cs.SetR(0.0)
	cs.SetB(0.0)
	cs.SetG(0.0)

	ts := NewTextWithMarkerShader(txt, buff, insets, *cs, yInset, xInset)
	id := AssignAndIncrement(ts)
	ts.Id = id

	return id
}
