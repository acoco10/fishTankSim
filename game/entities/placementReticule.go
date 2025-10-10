package entities

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"image"
)

func placemenReticuleUpdateFunc(ent *Entity, gs GameState) {
	MoveEnt(ent, gs)

	valid := true
	for _, prop := range PM.placedProps {
		valid = !pointInPolygon(image.Point{int(ent.Sprite.X), int(ent.Sprite.Y)}, prop.baseCorners)
	}

	if registry.ClickCheck() && valid {
		tag := ent.Parameters.Strings[Tag]
		ent.EventHub.Publish(PlacementPicked{PlacementFor: tag, X: ent.Sprite.X, Y: ent.Sprite.Y + float32(ent.Sprite.SpriteHeight()), Z: ent.Parameters.Ints[lastZ]})
		RemoveEntity(ent.Id)
	}

}

func LoadPlaceMentReticule(zBounds [13]image.Rectangle, tag string, hub *tasks.EventHub) {
	img, _ := LoadPlacementImg()
	ev := events.PlacementMode{}
	hub.Publish(ev)
	x, _, currentZ := positionPointBasedOnCursorOnZslice(zBounds)

	sp := &sprite.Sprite{Img: img, Y: float32(zBounds[currentZ].Max.Y), X: float32(x), UnFocusable: true}
	ent := &Entity{Sprite: sp}
	ent.Z = 13
	ent.UpdateFunc = placemenReticuleUpdateFunc
	ent.Sprite.DOptsUpdaterParams = make(map[string]float64)
	ent.Sprite.DOptsUpdaterTag = "offset"
	ent.Sprite.DOptsUpdaterParams["offSetX"] = -6
	ent.Sprite.DOptsUpdaterParams["offSetY"] = 6
	RegisterEntity(ent)
	ent.Parameters.Strings[Tag] = tag
	ent.EventHub = hub

}
