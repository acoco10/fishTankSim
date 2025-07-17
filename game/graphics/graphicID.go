package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

var GraphicId = 1

var GraphMap = make(map[int]Graphic)

type Graphic interface {
	Draw(screen *ebiten.Image)
	Update()
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

func DeInitAllGraphics() {
	GraphMap = make(map[int]Graphic)
}

func DrawGraphics(screen *ebiten.Image) {
	for _, graph := range GraphMap {
		graph.Draw(screen)
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
	id := NewGraphicText(&msg, 24, x, y, false, cs, 0, true)
	return id
}

func NewUpdateAbleTextGraphic(msg *string, x, y float64) int {
	cs := ebiten.ColorScale{}
	cs.SetR(0.9)
	cs.SetB(0.9)
	cs.SetG(0.9)
	cs.SetA(1.0)
	id := NewGraphicText(msg, 24, x, y, false, cs, 0, false)
	return id
}

func NewTravelingEffect(sprite *sprite.AnimatedSprite, x, y *float32) int {
	eff := &TaggedTravelingGraphic{sprite: sprite, x: x, y: y}
	GraphicId++
	GraphMap[GraphicId] = eff
	return GraphicId
}

type TaggedTravelingGraphic struct {
	sprite *sprite.AnimatedSprite
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
	t.sprite.Draw(screen)
}

func AddSpriteGraphic(graphic *SpriteGraphic) int {
	id := AssignAndIncrement(graphic)
	return id
}
