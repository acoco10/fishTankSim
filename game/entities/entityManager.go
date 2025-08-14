package entities

import (
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/input"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"log"
	"sort"
)

var currentEntID uint32

var EntityCatalogue = make(map[uint32]*Entity)

var LiveList []*Entity

func RegisterEntity(ent *Entity) {
	ent.Id = currentEntID
	EntityCatalogue[currentEntID] = ent
	ent.Draw = true
	currentEntID++
	LiveList = append(LiveList, ent)
	ZSortEntities()
}

func GetEntity(id uint32) (*Entity, bool) {
	entity, exists := EntityCatalogue[id]

	if exists {
		return entity, true
	}

	return nil, false
}

func RemoveEntity(id uint32) {
	delete(EntityCatalogue, id)
	for i, ent := range LiveList {
		if id == ent.Id {
			LiveList = append(LiveList[:i], LiveList[i+1:]...)
			return
		}
	}
}

type Entity struct {
	Id                 uint32
	LinkedID           uint32 //sometimes we want to manipulate entities from a central update entity like ph strip and box of ph strips
	Sprite             *sprite.Sprite
	CreatureData       *CreatureData
	UiData             *UiSpriteData
	PropData           *StructureProp
	ParticleData       *Particle
	AnimationMap       map[string]*sprite.Sprite
	TriggeredAnimation string
	EventHub           *tasks.EventHub
	Draw               bool
	GraphicManager     *graphics.GraphicManager
	UpdateFunc         func(entity *Entity)
	Parameters         map[string]any
}

type GameState struct {
	MouseFlags         *input.MouseFlags
	FocusedEntity      *Entity
	HoveredEntities    []*Entity
	HoveredUiSprite    *Entity
	WasHoveredUiSprite *Entity
}

func UpdateEntities(gs *GameState) {
	//if we are doing shit like update game state within this function then be mindful in which order you update
	//And what is utilizing game state.

	gs.HoveredEntities = []*Entity{}
	lastuiSpriteStillHovered := false

	if gs.HoveredUiSprite != nil {
		if gs.HoveredUiSprite.Sprite.SpriteHoveredWithBuffer(20) {
			lastuiSpriteStillHovered = true
		} else {
			gs.HoveredUiSprite = nil
		}
	}
	GetHoveredUISprite(gs, lastuiSpriteStillHovered)

	for _, ent := range LiveList {

		if ent.UpdateFunc != nil {
			ent.UpdateFunc(ent)
		}
		if ent.Sprite != nil {
			if ent.Sprite.Remove {
				RemoveEntity(ent.Id)
				return
			}
			ent.Sprite.Update()
		}

		if ent.UiData != nil && !gs.MouseFlags.WindowOpen && !registry.Config.Zoom {
			ent.UpdateUiSprite(gs)
		}

		if ent.CreatureData != nil {
			ent.FishUpdate()
		}
		if ent.PropData != nil {
			ent.PropData.Update()
		}
		if ent.ParticleData != nil {
			ent.ParticleData.Update()
		}
		if ent.GraphicManager != nil {
			ent.GraphicManager.Update()
		}
	}

}

func GetHoveredUISprite(gs *GameState, uiSpriteStillHovered bool) {
	for _, ent := range LiveList {
		if ent.Sprite != nil {
			if ent.Sprite.SpriteHovered() {
				gs.HoveredEntities = append(gs.HoveredEntities, ent)
				//only let the first ui sprite (highest in draw order/ printed last be hovered
				//confusing for player and programmer if more than one ui entity can be hovered at a time
				if ent.UiData != nil && !uiSpriteStillHovered && ent.UiData.state != Disabled {
					gs.HoveredUiSprite = ent
				}
			}
		}
	}
}

func DrawEntities(screen *ebiten.Image) {
	for _, ent := range LiveList {

		if !ent.Draw {
			continue
		}
		if ent.Sprite != nil {
			if ent.Sprite.Z < 4 {
				ent.Sprite.Draw(screen)
			}
		}
		if ent.ParticleData != nil {
			ent.ParticleData.Draw(screen)
		}
	}
}

func DrawNonZoomedEntities(screen *ebiten.Image) {
	for _, ent := range LiveList {

		if !ent.Draw {
			continue
		}

		if ent.Sprite != nil {
			if ent.Sprite.Z < 4 {
				continue
			}
			ent.Sprite.Scale = registry.Config.ZoomFactor
			ent.Sprite.Draw(screen)
		}

		if ent.ParticleData != nil {
			ent.ParticleData.Draw(screen)
		}
	}
}

func UpdateEntityZAndReSortEntitySlice(id uint32, zLayer int) {
	ent, exists := GetEntity(id)
	if !exists {
		log.Println("tried to update an entityID Z that doesnt exist")
		return
	}
	if ent.Sprite != nil {
		ent.Sprite.Z = zLayer
	}
	ZSortEntities()

}

func ZSortEntities() {
	sort.SliceStable(LiveList, func(i, j int) bool {
		if LiveList[i].Sprite.Z != LiveList[j].Sprite.Z {
			return LiveList[i].Sprite.Z < LiveList[j].Sprite.Z
		} else {
			return LiveList[i].Sprite.LayerIndex < LiveList[j].Sprite.LayerIndex
		}
	})
}

func UnFocus(ID uint32) {
	e, exists := GetEntity(ID)
	if !exists {
		log.Fatal("making this a crash for now for detecting un focus events on ents that dont exist")
	}
	if e.UiData != nil {
		UiSpriteTurnOffEverything(e)
	}

	if e.Sprite != nil {
		e.Sprite.UnFocus()
		graphics.DeInitGraphics(e.Sprite.PublishedGraphicId)
		e.Sprite.PublishedGraphicId = []int{}
	}

	if e.CreatureData != nil {
		UpdateEntityZAndReSortEntitySlice(e.Id, 1)
	}

	if e.EventHub != nil {
		e.EventHub.Publish(events.UnFocus{EntID: ID})
	} else {
		log.Fatal("making this a crash for now for detecting un focus events on ents that dont have event hubs initiated", ID)
	}

	if e.LinkedID != 0 {
		DeInitLinkedEnts(e.LinkedID)
	}
}

func DeInitLinkedEnts(id uint32) {
	curEnt, exists := GetEntity(id)
	if !exists {
		return
	}
	for curEnt.LinkedID != 0 {
		linkedId := curEnt.LinkedID
		RemoveEntity(curEnt.Id)
		curEnt, exists = GetEntity(linkedId)
		if !exists {
			return
		}
	}
	RemoveEntity(curEnt.Id)
}

func Focus(ID uint32) {
	e, exists := GetEntity(ID)
	if !exists {
		log.Fatal("making this a crash for now for detecting un focus events on ents that dont exist")
	}

	if e.Sprite != nil {
		if e.UiData != nil && e.UiData.state == Disabled {
			return
		}
		e.Sprite.Focus()
	}

	if e.CreatureData != nil {
		e.Sprite.Z = 3
		eff := LoadCreatureEffect("Day", e)
		AddEffectToSprite(eff, e.Sprite)
		MakeFishMenu(e.Id)
	}

	if e.EventHub != nil {
		e.EventHub.Publish(events.Focus{EntID: ID})
	} else {
		log.Fatal("making this a crash for now for detecting un focus events on ents that dont have event hubs initiated", ID)
	}
}

func LoadCreatureEffect(state string, ent *Entity) *sprite.Sprite {

	switch state {
	case "Night":
		return entImportableLoaders.LoadEffect("zzz")
	case "Day":
		switch ent.CreatureData.HealthState {
		case Healthy:
			return entImportableLoaders.LoadEffect("happy")
		case Stressed:
			return entImportableLoaders.LoadEffect("stressed")
		}
	}
	return nil
}

func AddEffectToSprite(effect *sprite.Sprite, sp *sprite.Sprite) {
	id := graphics.NewTravelingEffect(effect, &sp.X, &sp.Y)
	sp.SavePublishedGraphicID(id)
}

func ReFocus(ID uint32) {
	e, exists := GetEntity(ID)
	if !exists {
		log.Fatal("making this a crash for now for detecting un focus events on ents that dont exist")
	}

	if e.UiData != nil {
		x, y := util.GetScaledCursorPosition()
		e.Sprite.X = float32(x)
		e.Sprite.Y = float32(y)
		e.Sprite.Focus()
	}

	if e.EventHub != nil {
		e.EventHub.Publish(events.Focus{EntID: ID})
	} else {
		log.Fatal("making this a crash for now for detecting un focus events on ents that dont have event hubs initiated", ID)
	}
}

func PHModifier(ent *Entity) {
	centerishPoint := image.Point{X: registry.Config.ScreenWidth/2 - 100, Y: registry.Config.ScreenHeight/2 - 100}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && util.DistanceBetweenPoints(image.Point{X: int(ent.Sprite.X), Y: int(ent.Sprite.Y)}, centerishPoint) < 100 {
		RemoveEntity(ent.Id)
		ent.EventHub.Publish(events.UISpriteAction{UiSprite: "PHModifier", UiSpriteAction: ent.Parameters["tag"].(string)})
	}

}
