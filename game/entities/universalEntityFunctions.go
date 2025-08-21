package entities

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

func ClosestHoveredEntityToCursor(cursorX int, cursorY int, entSlice []*Entity) (*Entity, float64) {
	//cursor x and y taking as inputs in case cursor position changes

	cursorPoint := image.Point{X: cursorX, Y: cursorY}

	var distMap = make(map[float64]*Entity)
	var closestDistance float64

	closestDistance = 1000

	for _, ent := range entSlice {

		if ent.Sprite == nil {
			continue
		}

		if !ent.Sprite.SpriteHovered() {
			continue
		}

		if ent.Sprite.Unfocusable {
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

	if closestDistance == 1000 {
		return nil, 0
	}
	_, exists := distMap[closestDistance]
	if !exists {
		return nil, 0
	}

	if len(distMap) == 1 {
		return distMap[closestDistance], closestDistance
	} else {
		//if more than one sprite is "very close" return the first ui sprite rather then a creature
		var retSprite *Entity
		for _, ent := range distMap {
			if ent.UiData != nil {
				retSprite = ent
			}
		}
		return retSprite, closestDistance
	}
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
	if closestDistance < 200 {
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
		}
		return
	}

	if gs.MouseFlags.WindowOpen {
		return
	}
	xAtClick, yAtClick := util.GetScaledCursorPosition()
	closestToCursor, distance := ClosestHoveredEntityToCursor(xAtClick, yAtClick, LiveList)
	if closestToCursor == nil {
		log.Println("|||function for finding closest entity is returning nil|||")
		return
	}

	println("closest to cursor = entity ID", closestToCursor.Id, "distance =", distance)
	if gs.FocusedEntity == closestToCursor {
		return //if the focused entity is the one we would switch to just return
	}
	if gs.FocusedEntity != nil && gs.FocusedEntity.Sprite.AbleToBeUnfocusedAutomatically {
		// if we made it past the entity not being the currently focused one, unfocus whatever the focused one is
		println("unfocusing entity that is able to be auto unfocused")
		if gs.FocusedEntity.Sprite != nil {
			UnFocus(gs.FocusedEntity.Id)
			gs.FocusedEntity = nil
		}
	}
	if !gs.MouseFlags.WindowOpen {
		//if a sprite is focused they can do focused stuff in update rather then have input detection in each sprite
		if gs.FocusedEntity == nil {
			if closestToCursor.Sprite != nil {
				Focus(closestToCursor.Id)
				gs.FocusedEntity = closestToCursor
			}
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
