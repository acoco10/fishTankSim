package entities

import (
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/input"
	"github.com/acoco10/fishTankWebGame/game/movement"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	box2d "github.com/oliverbestmann/box2d-go"
	"image"
	"log"
)

var currentEntID uint32

var EntityCatalogue = make(map[uint32]*Entity)

var LiveList []*Entity

var SpriteList [20][]*Entity

var ParticleList []*ParticleSystem

func RegisterEntity(ent *Entity) uint32 {
	ent.Id = currentEntID
	EntityCatalogue[currentEntID] = ent
	ent.Draw = true
	currentEntID++
	LiveList = append(LiveList, ent)
	ent.Parameters = make(map[string]any)

	if ent.Sprite != nil {
		if ent.Z > 20 {
			log.Fatal(ent)
		}
		SpriteList[ent.Z] = append(SpriteList[ent.Z], ent)
	}

	if ent.ParticleSystem != nil {
		ParticleList = append(ParticleList, ent.ParticleSystem)
	}

	ZSortEntities()
	return ent.Id
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
	Id                          uint32
	LinkedID                    uint32 //sometimes we want to manipulate entities from a central update entity like ph strip and box of ph strips
	Z                           int
	LayerIndex                  int // optional parameter for tightly defined layers
	Sprite                      *sprite.Sprite
	CreatureData                *CreatureData
	UiData                      *UiSpriteData
	PropData                    *StructureProp
	ParticleData                *FoodParticle
	TriggeredAnimation          string
	EventHub                    *tasks.EventHub
	Draw                        bool
	GraphicManager              *graphics.GraphicManager
	UpdateFunc                  func(entity *Entity)
	Parameters                  map[string]any
	physics                     box2d.Body
	TankMovement                *TankCharacter
	MovementState               *movement.State
	MovementSystem              *movement.System
	StateMachine                *StateMachine
	ParticleSystem              *ParticleSystem
	DeposeAfterNAnimationCycles int
	AnimationCycles             int
	NoZoom                      bool
	LifeTime                    float64
	EndAfter                    float64
}

type StateMachine struct {
	States       map[int]*StateHandler
	CurrentState int
}

func (s *StateMachine) Transition(ent *Entity) {
	if s.States[s.CurrentState].TransitionFunc != nil {
		s.States[s.CurrentState].TransitionFunc(ent)
	}
	s.CurrentState = s.States[s.CurrentState].TransitionTo

}

type StateHandler struct {
	Updater        func(entity *Entity, state GameState)
	TransitionTo   int
	TransitionFunc func(entity *Entity)
}

func (s *StateMachine) Update(ent *Entity, gs GameState) {
	s.States[s.CurrentState].Updater(ent, gs)
}

type GameState struct {
	MouseFlags           *input.MouseFlags
	FocusedEntity        *Entity
	HoveredUiSprite      *Entity
	WasHoveredUiSprite   *Entity
	ActiveCollisions     []FishCollision
	PreZoomFocusedEntity *Entity
	CursorUpdater        *CursorUpdater
	CollisionMap         map[string]image.Rectangle
	Zbounds              [13]image.Rectangle
	PhysicsObjects       []uint32
}
type FishCollision struct {
	image.Rectangle
	Z int
}

func UpdateEntities(gs *GameState) {
	//if we are doing shit like update game state within this function then be mindful in which order you update
	//And what is utilizing game state.

	lastUISpriteStillHovered := false
	if !gs.MouseFlags.WindowOpen && gs.FocusedEntity == nil {
		if gs.HoveredUiSprite != nil {
			if gs.HoveredUiSprite.Sprite.SpriteHoveredWithBuffer(20) {
				lastUISpriteStillHovered = true
			} else {
				gs.HoveredUiSprite = nil
			}
		}
		GetHoveredUISprite(gs, lastUISpriteStillHovered)
	} else {
		gs.HoveredUiSprite = nil
	}
	for _, ent := range LiveList {

		ent.LifeTime += 0.016
		if ent.EndAfter != 0 {
			ent.CheckAndSelfRemove()
		}
		if ent.UpdateFunc != nil {
			ent.UpdateFunc(ent)
		}
		if ent.StateMachine != nil {
			ent.StateMachine.Update(ent, *gs)
		}

		if ent.Sprite != nil {
			if ent.Sprite.Remove {
				RemoveEntity(ent.Id)
				return
			}
			if ent.Sprite.CurrentAnimation != "" {
				animation := ent.Sprite.GetAnimation()
				if animation.Frame() == animation.LastF {
					if animation.FrameCounter == 1 {
						ent.AnimationCycles++
					}
				}
				if ent.DeposeAfterNAnimationCycles != 0 && ent.DeposeAfterNAnimationCycles == ent.AnimationCycles {
					RemoveEntity(ent.Id)
				}
			}
			ent.Sprite.Update()
		}

		if ent.UiData != nil && !gs.MouseFlags.WindowOpen {
			if registry.Config.Zoom {
				if ent.Sprite.X != ent.UiData.baseX {
					ent.Sprite.Shader = nil
					ent.Sprite.X = ent.UiData.baseX
					ent.Sprite.Y = ent.UiData.baseY
					ent.Z = 0
					ZSortEntities()
				}
				continue
			}
			ent.UpdateUiSprite(gs)

		}

		if ent.CreatureData != nil {
			ent.FishUpdate(gs)
		}
		if ent.PropData != nil {
			ent.UpdateProp(gs.Zbounds)
		}
		if ent.ParticleData != nil {
			ent.ParticleData.Update()
		}
		if ent.GraphicManager != nil {
			ent.GraphicManager.Update()
		}
		if ent.ParticleSystem != nil {
			ent.ParticleSystem.Update()
		}
	}

}

func GetHoveredUISprite(gs *GameState, uiSpriteStillHovered bool) {
	for _, ent := range LiveList {
		if ent.UiData == nil {
			continue
		}
		if ent.UiData.Label == string(WhiteBoard) {
			continue
		}
		if ent.Sprite == nil {
			continue
		}
		if ent.UiData.state == Disabled {
			continue
		}
		if ent.Sprite.SpriteHovered() {
			//only let the first ui sprite (highest in draw order/ printed last be hovered
			//confusing for player and programmer if more than one ui entity can be hovered at a time
			gs.HoveredUiSprite = ent
		}
	}
}

func (ent *Entity) CheckAndSelfRemove() {
	if ent.LifeTime >= ent.EndAfter {
		RemoveEntity(ent.Id)
	}
}

func DrawEntities(screen *ebiten.Image, gs *GameState) {
	for _, ps := range ParticleList {
		ps.Draw() // draws to "sprite buffer" all sprites are drawn together below
	}

	for z := 0; z < 20; z++ {
		for _, ent := range SpriteList[z] {
			if ent.Draw && ent.Sprite != nil {
				ent.Sprite.Draw(screen)
			}
		}
	}

	/*	if !ent.Draw || ent.NoZoom {
			continue
		}
		if ent.Sprite != nil {
			{
				ent.Sprite.Draw(screen)
			}
		}
		if ent.PropData != nil {
			ent.PropData.Draw(screen)
		}
		if ent.ParticleData != nil {
			ent.ParticleData.Draw(screen)
		}
		if ent.ParticleSystem != nil {
			ent.ParticleSystem.Draw(screen)
		}*/
}

func DrawFocusedEntityNoLightingShader(screen *ebiten.Image, gs *GameState) {
	if gs.FocusedEntity != nil {
		gs.FocusedEntity.Sprite.Draw(screen)
	}
}
func DrawNonZoomedEntities(screen *ebiten.Image) {
	for _, ent := range LiveList {

		if !ent.Draw || !ent.NoZoom {
			continue
		}

		if ent.Sprite != nil {
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
		ent.Z = zLayer
	}
	ZSortEntities()

}

func ZSortEntities() {
	// Clear all layers first
	for i := range SpriteList {
		SpriteList[i] = SpriteList[i][:0] // Keep capacity, reset length
	}

	// Re-populate from LiveList
	for _, ent := range LiveList {
		if ent.Sprite != nil {
			if ent.Z > 20 {
				log.Fatal(ent)
			}
			SpriteList[ent.Z] = append(SpriteList[ent.Z], ent)
		}
	}
}

func UnFocus(ID uint32) {
	e, exists := GetEntity(ID)
	if !exists {
		log.Fatal("making this a crash for now for detecting un focus events on ents that dont exist")
	}

	if e.Sprite != nil {
		e.Sprite.UnFocus()
		graphics.DeInitGraphics(e.Sprite.PublishedGraphicId)
		e.Sprite.PublishedGraphicId = []int{}
	}

	if e.UiData != nil {
		//UiSpriteTurnOffEverything(e)
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
	ZSortEntities()
}

func DeInitLinkedEnts(id uint32) {
	curEnt, exists := GetEntity(id)
	if !exists {
		return
	}
	for curEnt.LinkedID != 0 {
		if curEnt.Sprite != nil {
			if len(curEnt.Sprite.PublishedGraphicId) > 0 {
				graphics.DeInitGraphics(curEnt.Sprite.PublishedGraphicId)
			}
		}
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
		e.CreatureData.TargetZ = 10
		LoadCreatureEffect("Day", e)
		MakeFishMenu(e.Id)
	}

	if e.EventHub != nil {
		e.EventHub.Publish(events.Focus{EntID: ID})
	} else {
		log.Fatal("making this a crash for now for detecting un focus events on ents that dont have event hubs initiated", ID)
	}
	ZSortEntities()

}

func LoadCreatureEffect(state string, ent *Entity) {

	switch state {
	case "Night":
		LoadFollowEffectAsEnt("zzz", ent.Id, ent.EventHub)
	case "Day":
		switch ent.CreatureData.HealthState {
		case Healthy:
			LoadFollowEffectAsEnt("happy", ent.Id, ent.EventHub)
		case Stressed:
			LoadFollowEffectAsEnt("stressed", ent.Id, ent.EventHub)
		}
	}
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
		ent.EventHub.Publish(events.UISpriteAction{UiSprite: "PHModifier", UiSpriteAction: ent.Parameters["Tag"].(string)})
	}

}

func LoadFollowEffectAsEnt(eff string, targID uint32, hub *tasks.EventHub) {
	effect := entImportableLoaders.LoadEffect(eff)
	effect.Unfocusable = true
	effEnt := &Entity{Sprite: effect, UpdateFunc: FollowEnt, LinkedID: targID}
	effEnt.Z = 13
	effEnt.DeposeAfterNAnimationCycles = 10
	effect.AnimationMap[effect.CurrentAnimation].SpeedInTPS = 8
	RegisterEntity(effEnt)
	hub.Subscribe(events.UnFocus{}, func(e tasks.Event) {
		RemoveEntity(effEnt.Id)
	})
}

func FollowEnt(ent *Entity) {
	targetEnt, exists := GetEntity(ent.LinkedID)

	if !exists {
		log.Println("Follow effect lined to de initated sprite, deinitiating effect")
		RemoveEntity(ent.Id)
	}

	ent.Sprite.X = targetEnt.Sprite.X
	ent.Sprite.Y = targetEnt.Sprite.Y - float32(ent.Sprite.SpriteHeight()+10)
	if targetEnt.CreatureData != nil {
		if targetEnt.CreatureData.Flip {
			ent.Sprite.X -= float32(targetEnt.Sprite.GetAnimation().SpriteWidth)
		}
	}
}

type LightingEntManager struct {
	NightEnts []uint32
	DayEnts   []uint32
}

func (LE *LightingEntManager) Subscribe(hub *tasks.EventHub) {
	hub.Subscribe(events.LightEvent{}, func(e tasks.Event) {
		ev := e.(events.LightEvent)
		if ev.Day {
			for _, id := range LE.NightEnts {
				ent, exists := GetEntity(id)
				if exists {
					ent.Draw = false
				}
			}
			for _, id := range LE.DayEnts {
				ent, exists := GetEntity(id)
				if exists {
					ent.Draw = true
				}
			}
		} else {
			for _, id := range LE.NightEnts {
				ent, exists := GetEntity(id)
				if exists {
					ent.Draw = true
				}
			}
			for _, id := range LE.DayEnts {
				ent, exists := GetEntity(id)
				if exists {
					ent.Draw = false
				}
			}
		}

	})
	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		for _, id := range LE.NightEnts {
			ent, exists := GetEntity(id)
			if exists {
				ent.Draw = false
			}
		}
		for _, id := range LE.DayEnts {
			ent, exists := GetEntity(id)
			if exists {
				ent.Draw = true
			}
		}
	})
}
