package entities

import (
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/input"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	physics "github.com/acoco10/fishTankWebGame/game/testPhysics"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image"
	"log"
)

var currentEntID uint32 = 1

var EntityCatalogue = make(map[uint32]*Entity)

var LiveList []*Entity

var SpriteList [20][]*Entity

var UiEffectSpriteList [10][]*Entity

var OverZoomSpriteList [10][]*Entity

var ParticleList []*ParticleSystem

var ParticleEntList []*EntityParticle

var CreatureList []*Entity

var TextGraphicList []*Entity

var FadeTextGraphicList []*Entity

var GraphicList []*Entity

var SpriteWithBufferDst []*sprite.Sprite

var UiData []*UiSpriteData

var GEntityManager *EntityManager

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
	Player               *Player
	Debug                string
	ZoomedFor            ZoomState
}

type FishCollision struct {
	image.Rectangle
	Z int
}
type EntityManager struct {
	hub                 *tasks.EventHub
	EntityCatalogue     map[uint32]*Entity
	LiveList            []*Entity
	SpriteList          [20][]*Entity
	UiEffectSpriteList  [10][]*Entity
	ParticleList        []*ParticleSystem
	ParticleEntList     []*EntityParticle
	CreatureList        []*Entity
	GraphicList         []*Entity
	SpriteWithBufferDst []*sprite.Sprite
	UiData              []*UiSpriteData
}

type Entity struct {
	labelForDebugging           string
	flags                       entFlags
	Id                          uint32
	TempLinkedID                uint32 //sometimes we want to manipulate entities from a central update entity like ph strip and box of ph strips
	LinkedID                    uint32
	Z                           int
	LayerIndex                  int // optional parameter for tightly defined layers
	Sprite                      *sprite.Sprite
	CreatureData                *CreatureData
	UiData                      *UiSpriteData
	DoAt                        map[string]func(entity *Entity, gs GameState) //subState one off functions (smaller state one offs that would be cumbersome to populate into a state machine for testing
	PropData                    *StructureProp
	ParticleData                *EntityParticle
	TriggeredAnimation          string
	EventHub                    *tasks.EventHub
	Draw                        bool
	PreZoomDraw                 bool
	GraphicManager              *graphics.GraphicManager
	UpdateFunc                  func(entity *Entity, state GameState)
	DrawTimeFunc                func(entity *Entity, state GameState)
	Parameters                  EntityParameters
	Flags                       map[string]bool
	StateMachine                *StateMachine
	ParticleSystem              *ParticleSystem
	EndOfDayNUnSubscribeEvents  [365][]tasks.CreatedEvent
	DeposeAfterNAnimationCycles int
	AnimationCycles             int
	NoZoom                      bool
	LifeTime                    float64
	EndAfter                    float64
	ShaderTextGraphic           *graphics.TextWithShader
	RectGraphic                 *graphics.RectGraphic
	PublishedGraphicIDs         []int
	effectDeInitHandlers        []DeInitFunc
	effectUpdateHandler         func(param float64)
	UIWidget                    *widget.Container
	updateEntities              map[uint32]struct{}
	graphicText                 *graphics.FadeInText
}

func InitGEntManager(hub *tasks.EventHub) {
	GEntityManager = &EntityManager{}
	GEntityManager.hub = hub
	GEntityManager.EntityCatalogue = make(map[uint32]*Entity)
}

func RegisterEntity(ent *Entity) uint32 {
	ent.Id = currentEntID
	EntityCatalogue[currentEntID] = ent
	ent.Draw = true

	currentEntID++

	hasBuffer := ent.Sprite != nil && ent.Sprite.BufferDst != nil

	ent.EventHub = GEntityManager.hub
	ent.universalEntitySubscriptions()

	if !hasBuffer && !ent.IsOverUI() && !ent.HasOverZoom() {
		if ent.Z > 20 {
			log.Fatal("Z is out of bounds in register Entity function", ent)
		}
		SpriteList[ent.Z] = append(SpriteList[ent.Z], ent)
	}

	if ent.HasOverZoom() {
		if ent.Z > 10 {
			log.Fatal("Z is out of bounds in register Entity function", ent)
		}
		OverZoomSpriteList[ent.Z] = append(OverZoomSpriteList[ent.Z], ent)
	}

	if ent.HasOverUi() {
		if ent.Z > 10 {
			log.Fatal("Z is out of bounds in register Entity function", ent)
		}
		UiEffectSpriteList[ent.Z] = append(UiEffectSpriteList[ent.Z], ent)
	}

	if ent.RectGraphic != nil {
		GraphicList = append(GraphicList, ent)
	}

	if ent.graphicText != nil {
		FadeTextGraphicList = append(FadeTextGraphicList, ent)
	}

	if hasBuffer {
		SpriteWithBufferDst = append(SpriteWithBufferDst, ent.Sprite)
		return ent.Id
	}

	if ent.ParticleSystem != nil {
		ParticleList = append(ParticleList, ent.ParticleSystem)
	}

	if ent.UiData != nil {
		UiData = append(UiData, ent.UiData)
	}

	if ent.ShaderTextGraphic != nil {
		TextGraphicList = append(TextGraphicList, ent)
		return ent.Id
	}
	if ent.CreatureData != nil {
		CreatureList = append(CreatureList, ent)
	}

	if ent.ParticleData != nil {
		ParticleEntList = append(ParticleEntList, ent.ParticleData)
	}

	LiveList = append(LiveList, ent)
	ZSortEntities()
	return ent.Id
}

func GetEntity(id uint32) (*Entity, bool) {
	if id == 0 {
		return nil, false
	}
	entity, exists := EntityCatalogue[id]
	if exists {
		return entity, true
	}
	return nil, false
}

func RemoveEntity(id uint32) {
	oEnt, exist := GetEntity(id)
	if !exist {
		return
	}

	oEnt.DeInitEffects()

	delete(EntityCatalogue, id)

	for i, spEnt := range SpriteList[oEnt.Z] {
		if spEnt.Id == id {
			SpriteList[oEnt.Z] = append(SpriteList[oEnt.Z][:i], SpriteList[oEnt.Z][i+1:]...)
		}
	}

	if oEnt.HasOverZoom() {
		for i, spZent := range OverZoomSpriteList[oEnt.Z] {
			if spZent.Id == id {
				OverZoomSpriteList[oEnt.Z] = append(OverZoomSpriteList[oEnt.Z][:i], OverZoomSpriteList[oEnt.Z][i+1:]...)
			}
		}
	}

	for i, ent := range LiveList {
		if id == ent.Id {
			LiveList = append(LiveList[:i], LiveList[i+1:]...)
			continue
		}
	}

	for i, ent := range GraphicList {
		if id == ent.Id {
			GraphicList = append(GraphicList[:i], GraphicList[i+1:]...)
			continue
		}
	}

	for i, ent := range TextGraphicList {
		if id == ent.Id {
			TextGraphicList = append(TextGraphicList[:i], TextGraphicList[i+1:]...)
			continue
		}
	}

	for i, ent := range FadeTextGraphicList {
		if id == ent.Id {
			FadeTextGraphicList = append(FadeTextGraphicList[:i], FadeTextGraphicList[i+1:]...)
			continue
		}
	}

	filtered := ParticleList[:0]
	for _, ps := range ParticleList {
		if oEnt.ParticleSystem != ps {
			filtered = append(filtered, ps)
		}
	}
	ParticleList = filtered

	sFiltered := SpriteWithBufferDst[:0]
	for _, s := range SpriteWithBufferDst {
		if oEnt.Sprite != s {
			sFiltered = append(sFiltered, s)
		}
	}

	SpriteWithBufferDst = sFiltered

}

type ZoomState uint8

const (
	NotZoomed ZoomState = iota
	ZoomedByPlayer
	ZoomedForFeeding
	ZoomedForPlacement
	PlayerZoomed
)

func UpdateEntities(gs *GameState, TankPhysics *physics.TankPhysics) {
	//if we are doing shit like update game state within this function then be mindful in which order you update
	//And what is utilizing game state.

	lastUISpriteStillHovered := false
	if !gs.MouseFlags.WindowOpen {
		if gs.HoveredUiSprite != nil {
			if gs.HoveredUiSprite.Sprite.SpriteHoveredWithBuffer(20) {
				lastUISpriteStillHovered = true
			} else {
				gs.HoveredUiSprite = nil
			}
		}
		if !lastUISpriteStillHovered {
			GetHoveredUISprite(gs)
		}
	}

	for _, ent := range LiveList {

		ent.LifeTime += 0.016
		if ent.EndAfter != 0 {
			ent.CheckAndSelfRemove()
		}
		if ent.UpdateFunc != nil {
			ent.UpdateFunc(ent, *gs) //script
		}
		if ent.StateMachine != nil {
			ent.StateMachine.Update(ent, *gs) //more crystallized state machine
		}

		if ent.Sprite != nil {
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

		if ent.UiData != nil {
			ent.UpdateUiSprite(gs)
		}

		if ent.PropData != nil {
			//ent.UpdateProp(gs.Zbounds)
		}
		if ent.ParticleData != nil {
			if ent.ParticleData.body != nil {
				TankPhysics.Water.ApplyWaterForces(ent.ParticleData.body)
			}
			ent.ParticleData.Update()
		}
		if ent.GraphicManager != nil {
			ent.GraphicManager.Update()
			if ent.GraphicManager.GmState == graphics.Finished {
				if ent.GraphicManager.FinishedFunc != nil {
					ent.GraphicManager.FinishedFunc(ent)
				}
				ent.GraphicManager = nil
			}
		}
		filtered := ent.PublishedGraphicIDs[:0] // reuse underlying array
		for _, graphicId := range ent.PublishedGraphicIDs {
			if _, exists := graphics.GraphMap[graphicId]; exists {
				filtered = append(filtered, graphicId)
			}
		}
		ent.PublishedGraphicIDs = filtered

	}

	for _, ps := range ParticleList {
		ps.Update()
	}

	for _, graphic := range TextGraphicList {
		graphic.ShaderTextGraphic.Update()
	}

	for _, graphic := range GraphicList {
		graphic.RectGraphic.Update()
	}

	for _, textGraphic := range FadeTextGraphicList {
		textGraphic.graphicText.Update()
	}

	for _, sp := range SpriteWithBufferDst {
		sp.Update()
	}

}

func GetHoveredUISprite(gs *GameState) {
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
		if ent.UiData.state == Disabled || ent.Draw == false {
			continue
		}
		if ent.Sprite.SpriteHoveredWithBuffer(20) {
			//only let the first ui sprite (highest in draw order/ printed last be hovered
			//confusing for player and programmer if more than one ui entity can be hovered at a time
			gs.HoveredUiSprite = ent
			return
		}
	}
}

func (ent *Entity) CheckAndSelfRemove() {
	if ent.LifeTime >= ent.EndAfter {
		RemoveEntity(ent.Id)
	}
}

func (ent *Entity) buffCheck() bool {
	return ent.Draw && ent.Sprite != nil && ent.Sprite.IsBuffer
}

func DrawEntities(screen *ebiten.Image, gs GameState) {

	//clear sprite buffers before redraw graphics/particles to them
	for z := 0; z < 20; z++ {
		for _, ent := range SpriteList[z] {
			if ent.buffCheck() {
				ent.Sprite.Img.Clear()

			}
		}
	}

	for _, ent := range LiveList {
		if ent.DrawTimeFunc != nil {
			ent.DrawTimeFunc(ent, gs)
		}
	}

	for _, s := range SpriteWithBufferDst {
		s.Draw(s.BufferDst)
	}

	//draw to sprite buffers
	for _, ps := range ParticleList {
		ps.Draw()
	}

	for _, graphic := range TextGraphicList {
		graphic.ShaderTextGraphic.Draw()
	}

	for _, graphicText := range FadeTextGraphicList {
		graphicText.graphicText.DrawToDst()
	}

	for _, graphic := range GraphicList {
		graphic.RectGraphic.DrawWDst()
	}

	//draw all sprites in correct order

	for z := 0; z < 20; z++ {
		for _, ent := range SpriteList[z] {
			if ent.Draw && ent.Sprite != nil {
				if gs.Debug == "DebugOn" {
					if ent.Z <= 12 {
						ent.Sprite.DebugDraw(screen, ent.Z, gs.Zbounds[ent.Z])
					} else {
						ent.Sprite.Draw(screen)
					}
				} else {
					ent.Sprite.Draw(screen)
				}
			}
		}
	}

	if gs.Debug == "DebugOn" {
		for _, ent := range CreatureList {

			vector.StrokeCircle(screen, ent.CreatureData.TargetPoint.X, ent.CreatureData.TargetPoint.Y, 3, 1, colornames.Lightgreen, false)
			if ent.CreatureData.inBetweenPoint != nil {
				vector.StrokeCircle(screen, ent.CreatureData.inBetweenPoint.X, ent.CreatureData.inBetweenPoint.Y, 3, 1, colornames.Darkorange, false)

			}
		}
		for _, uidat := range UiData {
			if gs.HoveredUiSprite != nil && gs.HoveredUiSprite.UiData == uidat {
				util.StrokeRectFromImageRect(uidat.ActivationRect, screen, colornames.Darkorange, false)
			}
		}

	}

}

func DrawUIfx(screen *ebiten.Image) {
	for z := 0; z < 10; z++ {
		for _, ent := range UiEffectSpriteList[z] {
			if ent.Draw && ent.Sprite != nil {
				ent.Sprite.Draw(screen)
			}
		}
	}
}

func DrawOverZoom(screen *ebiten.Image) {
	for z := 0; z < 10; z++ {
		for _, ent := range OverZoomSpriteList[z] {
			if ent.Draw && ent.Sprite != nil {
				ent.Sprite.Draw(screen)
			}
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
		if ent.Sprite != nil && !ent.IsOverUI() && !ent.HasOverZoom() {
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

	/*	if e.UiData != nil {
		//	e.UiData.Timers["unFocusBuffer"].TurnOn()
		//e.UiData.state = Disabled
	}*/

	if e.CreatureData != nil {
		if len(e.effectDeInitHandlers) > 0 {
			e.DeInitEffects()
		}
		UpdateEntityZAndReSortEntitySlice(e.Id, 1)
	}

	if e.EventHub != nil {
		e.EventHub.Publish(events.UnFocusEvent{EntID: ID})
	} else {
		println(e.Id)
	}

	if e.TempLinkedID != 0 {
		DeInitLinkedEnts(e.TempLinkedID)
	}
	ZSortEntities()
}

func (ent *Entity) DeInitEffects() {
	for _, handler := range ent.effectDeInitHandlers {
		handler()
	}
}

func DeInitLinkedEnts(id uint32) {
	curEnt, exists := GetEntity(id)
	if !exists {
		return
	}
	for curEnt.TempLinkedID != 0 {
		if curEnt.Sprite != nil {
			if len(curEnt.Sprite.PublishedGraphicId) > 0 {
				graphics.DeInitGraphics(curEnt.Sprite.PublishedGraphicId)
			}
		}
		linkedId := curEnt.TempLinkedID
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
		LoadFishFactorsMenu(*e.CreatureData, e.Id)
		LoadCreatureEffect("Day", e)
		//MakeFishMenu(e.Id)
	}

	if e.EventHub != nil {
		e.EventHub.Publish(events.FocusEvent{EntID: ID})
	}
	ZSortEntities()

}

func (ent *Entity) AddDeInitHandler(deInitFunc DeInitFunc) {
	ent.effectDeInitHandlers = append(ent.effectDeInitHandlers, deInitFunc)
}

func LoadCreatureEffect(state string, ent *Entity) {

	switch state {
	case "Night":
		ent.DeInitEffects()
		ent.AddDeInitHandler(LoadFollowEffectAsEnt("zzz", ent.Id, ent.EventHub, EntityParameters{}))
	case "Day":
		switch ent.CreatureData.HealthState {
		case Healthy:
			ent.DeInitEffects()
			ent.AddDeInitHandler(LoadFollowEffectAsEnt("happy", ent.Id, ent.EventHub, EntityParameters{}))
		case Stressed:
			ent.DeInitEffects()
			ent.AddDeInitHandler(LoadFollowEffectAsEnt("stressed", ent.Id, ent.EventHub, EntityParameters{}))
		case ReallyStressed:
			ent.DeInitEffects()
			ent.AddDeInitHandler(LoadFollowEffectAsEnt("reallyStressed", ent.Id, ent.EventHub, EntityParameters{}))
		case Sick:
			ent.DeInitEffects()
			ent.AddDeInitHandler(LoadFollowEffectAsEnt("sick", ent.Id, ent.EventHub, EntityParameters{}))
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
		e.EventHub.Publish(events.FocusEvent{EntID: ID})

	} else {
		//log.Fatal("making this a crash for now for detecting un focus events on ents that dont have event hubs initiated", ID)
	}
}

type DeInitFunc func()

func LoadFollowEffectAsEnt(eff string, targID uint32, hub *tasks.EventHub, params EntityParameters) DeInitFunc {
	effect := entImportableLoaders.LoadEffect(eff)
	effect.UnFocusable = true
	effEnt := &Entity{Sprite: effect, UpdateFunc: FollowEnt, TempLinkedID: targID, Parameters: params}
	effEnt.Z = 13
	effect.AnimationMap[effect.CurrentAnimation].SpeedInTPS = 8
	RegisterEntity(effEnt)

	return func() { RemoveEntity(effEnt.Id) }
}

func FollowEnt(ent *Entity, gs GameState) {
	targetEnt, exists := GetEntity(ent.TempLinkedID)
	if !exists {
		log.Println("Follow effect lined to de initiated sprite, deinitiating effect")
		RemoveEntity(ent.Id)
	}

	yOff := ent.Parameters.Floats[YEffectOffset]
	xOff := ent.Parameters.Floats[XEffectOffSet]

	x := targetEnt.Sprite.X - float32(targetEnt.Sprite.GetSpriteRect().Dx())/2 + float32(xOff)
	if targetEnt.Sprite.Flip {
		x += float32(targetEnt.Sprite.GetSpriteRect().Dx())
	}
	y := targetEnt.Sprite.Y - float32(ent.Sprite.GetSpriteRect().Dy())*1.5 + float32(yOff)

	ent.Sprite.X = x
	ent.Sprite.Y = y

	pos := ent.Parameters.Strings[Position]
	if pos == "center" {
		ent.Sprite.X = targetEnt.Sprite.X + float32(targetEnt.Sprite.GetSpriteRect().Dx()/3+ent.Sprite.GetSpriteRect().Dx()/2)
		ent.Sprite.Y = targetEnt.Sprite.Y
	}
	if pos == "bottom" {
		ent.Sprite.Y += float32(ent.Sprite.SpriteHeight()) + 10
	}

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
