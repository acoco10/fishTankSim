package entities

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

func ClosestHoveredEntityToCursor(cursorX int, cursorY int, entSlice []*Entity) *Entity {
	//cursor x and y taking as inputs in case cursor position changes

	for _, ent := range entSlice {

		if ent.UiData == nil {
			continue
		}
		if !ent.Sprite.SpriteHoveredWithBuffer(20) {
			continue
		}

		if ent.Sprite.UnFocusable {
			continue
		}
		if ent.UiData.state == Disabled {
			continue
		}

		return ent
	}
	return nil
}

func (ent *Entity) Subscribe(eventType tasks.Event, handler func(event tasks.Event)) int {
	return ent.EventHub.Subscribe(eventType, handler)
}

func (ent *Entity) UnSubscribe(eventType tasks.Event, id int) {
	ent.EventHub.Unsubscribe(eventType, id)
}

func (ent *Entity) SubWUnsubAfterCompletion(eventType tasks.Event, handler func(event tasks.Event)) {
	ent.Parameters.Ints[UnsubId] = ent.Subscribe(eventType, handler)
	ent.Subscribe(eventType, ent.UnsubFromUnsubID)
}

func (ent *Entity) UnsubFromUnsubID(event tasks.Event) {
	ent.UnSubscribe(event, ent.Parameters.Ints[UnsubId])
}

func (ent *Entity) Publish(event tasks.Event) {
	ent.EventHub.Publish(event)
}

func ClosestCreatureZoomed(cursorX int, cursorY int, entSlice []*Entity) *Entity {
	cursorPoint := image.Point{X: cursorX, Y: cursorY}

	var distMap = make(map[float64]*Entity)
	var closestDistance float64

	closestDistance = 1000

	for _, ent := range entSlice {

		if ent.CreatureData == nil {
			continue
		}

		s := ent.Sprite

		entPt := image.Point{X: int(s.X), Y: int(s.Y)}

		if entPt.X != 0 && entPt.Y != 0 {
			dis := util.DistanceBetweenPoints(cursorPoint, entPt)
			distMap[dis] = ent
			if dis < closestDistance {
				closestDistance = dis
			}
		}
	}
	if closestDistance < 100 {
		return distMap[closestDistance]
	}
	return nil
}

func MakeSpriteEntity(img *ebiten.Image, x float32, y float32) *Entity {
	sp := &sprite.Sprite{Img: img, X: x, Y: y}
	ent := &Entity{Sprite: sp}
	RegisterEntity(ent)
	return ent
}

func FocusOnClickedEntity(gs *GameState) {

	if registry.Config.Zoom {
		xAtClick, yAtClick := util.GetScaledCursorPosition()
		ent := ClosestCreatureZoomed(xAtClick, yAtClick, LiveList)
		if ent != nil {
			if gs.FocusedEntity != nil {
				UnFocus(gs.FocusedEntity.Id)
			}
			gs.FocusedEntity = ent
			Focus(gs.FocusedEntity.Id)
			return
		}
	}
}

func UpdateCursorForEntitiesWNormals(cursor []float64) {
	for _, ent := range LiveList {
		if ent.Sprite != nil {
			if ent.Sprite.NormalMap != nil {
				ent.Sprite.ShaderParams["Cursor"] = cursor
			}
		}
	}
}

func (e *Entity) AddUpdateEnt(entId uint32) {
	if e.updateEntities == nil {
		e.updateEntities = make(map[uint32]struct{})
	}
	e.updateEntities[entId] = struct{}{}
}

func PublishAtTimerUpdater(timer *util.Timer, e any) {
	ent := e.(*Entity)
	state := timer.Update()
	if state == util.Done {
		ent.EventHub.Publish(ent.Parameters.Events[EventAtTime])
		timer.TurnOff()
	}
}

func (e *Entity) DeInitEntityGraphics() {
	graphics.DeInitGraphics(e.PublishedGraphicIDs)

}

func (ent *Entity) universalEntitySubscriptions() {
	ent.EventHub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		ev := e.(events.DayOver)
		for _, event := range ent.EndOfDayNUnSubscribeEvents[ev.Day] {
			ent.EventHub.Unsubscribe(event.Ev, event.Id)
		}
	})
}

func (ent *Entity) GetLinkedEnt() (*Entity, bool) {
	return GetEntity(ent.LinkedID)
}
