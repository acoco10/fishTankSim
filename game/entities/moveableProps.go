package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
	"image"
	"io/fs"
	"log"
	"math/rand"
	"strings"
)

type Decoration string

const (
	Castle    Decoration = "castle"
	ZenFriend Decoration = "zenFriend"
	ZenBridge Decoration = "zenBridge"
	Log       Decoration = "log"
	CoolRock  Decoration = "coolRock"
	HotRock   Decoration = "hotRock"
)

type PropState uint8

const (
	Moveable PropState = iota
	SettingInPlace
	SetInPlace
)

var PM PropManager

type PropManager struct {
	placedProps              []*StructureProp
	placementReticule        *ebiten.Image
	invalidPlacementReticule *ebiten.Image
}

type PropAssets struct {
	normalMap *ebiten.Image
	layers    []*ebiten.Image
	merged    *ebiten.Image
}

type PropData struct {
	CollisionMap            map[string]image.Rectangle
	PtMap                   map[string][]image.Point
	PlacementParams         map[string]any
	ZBounds                 [13]image.Rectangle
	SubscriptionForNextProp tasks.Event
}

type StructureProp struct {
	state    PropState
	stateWas PropState
	*sprite.Sprite
	Sprite2         *sprite.Sprite
	shadowPoint     image.Point
	boundaries      image.Rectangle
	StaticShadow    bool
	baseY           float32
	Tag             string
	alreadyPlaced   bool
	cornerOffsets   [4]image.Point
	baseCorners     [4]image.Point //0-4 leftUpperCorner rightUpperCorner leftLowerCorner rightLowerCorner
	particleSystems []*ParticleSystem
	OnLanding       func(ent *Entity)
}

func (pm *PropManager) loadPropManager() {
	valid, invalid := LoadPlacementImg()
	pm.placementReticule = valid
	pm.invalidPlacementReticule = invalid
}

func CheckNormalTags(match string) bool {
	normalTags := []string{"normal", "_n"}
	for _, tag := range normalTags {
		if strings.Contains(strings.ToLower(match), tag) {
			return true
		}
	}
	return false
}

func CheckMergedTags(match string) bool {
	normalTags := []string{"complete", "full", "merged", "all", "joined"}
	for _, tag := range normalTags {
		if strings.Contains(strings.ToLower(match), tag) {
			return true
		}
	}
	return false
}

func LoadPropImgs(propName string) (PropAssets, error) {

	var matches []string
	var ps PropAssets

	fs.WalkDir(assets.ImagesDir, "images/tankProps", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.Contains(strings.ToLower(path), strings.ToLower(propName)) {
			matches = append(matches, path)
		}
		return nil
	})

	for _, match := range matches {
		if CheckNormalTags(match) {
			log.Printf("Succesfully loaded Normal Map for %s path: %s", propName, match)
			img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, match)
			if err != nil {
				return ps, err
			}
			ps.normalMap = img

		} else if CheckMergedTags(match) {
			log.Printf("Succesfully loaded  merged for %s path: %s", propName, match)
			img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, match)
			if err != nil {
				return ps, err
			}

			ps.merged = img
		} else {
			img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, match)
			if err != nil {
				return ps, err
			}
			//layers should sort alphabetically on their own but this is an assumption for sure, try to name layers the same thing besides their number
			ps.layers = append(ps.layers, img)
		}
	}

	return ps, nil
}

func (p *StructureProp) State() PropState {
	return p.state
}

func NewStructureProp(ps PropAssets, hub *tasks.EventHub, bounds image.Rectangle, tag string) *StructureProp {
	if PM.placementReticule == nil {
		PM.loadPropManager()
	}
	p := StructureProp{}
	sp := &sprite.Sprite{Img: ps.layers[0], NormalMap: ps.normalMap, UnFocusable: true}
	if len(ps.layers) > 1 {
		log.Println("loading second image for structure prop as sprite")
		sp2 := &sprite.Sprite{Img: ps.layers[1], NormalMap: ps.layers[1], UnFocusable: true}
		p.Sprite2 = sp2
		if ps.normalMap != nil {
			normalMapShader := registry.ShaderMap["NormalMap"]
			sp2.Shader = normalMapShader
			sp2.ShaderParams = make(map[string]any)
			sp2.ShaderParams["Cursor"] = []float64{0, 0, 100}
		} else {
			sp2.ShaderParams = make(map[string]any)
		}
	}

	if ps.normalMap != nil {
		normalMapShader := registry.ShaderMap["NormalMap"]
		sp.Shader = normalMapShader
		sp.ShaderParams = make(map[string]any)
		sp.ShaderParams["Cursor"] = []float64{0, 0, 100}
	} else {
		sp.ShaderParams = make(map[string]any)
	}

	p.Sprite = sp
	p.state = Moveable
	sprite.LoadPulseOutlineNormalShader(p.Sprite)
	p.boundaries = bounds

	if p.Sprite2 != nil {
		p.Sprite2.Y = p.Y
	}

	p.Tag = tag

	return &p
}

func InitPropStateMachine() *StateMachine {
	propUpdater1 := &StateHandler{Updater: MoveableUpdater, TransitionTo: 2}
	propUpdater2 := &StateHandler{Updater: SettingInPlaceUpdaterNew, TransitionTo: 3}
	propUpdater3 := &StateHandler{Updater: PropFallingUpdater, TransitionTo: 4}
	propUpdater4 := &StateHandler{Updater: TransitionToSetInPlace, TransitionTo: 5}
	propUpdater5 := &StateHandler{Updater: SetInPlaceUpdater, TransitionTo: 6}
	propUpdater6 := &StateHandler{Updater: SetInPlaceSelected, TransitionTo: 1}

	states := map[int]*StateHandler{
		1: propUpdater1,
		2: propUpdater2,
		3: propUpdater3,
		4: propUpdater4,
		5: propUpdater5,
		6: propUpdater6,
	}

	PropStateMachine := &StateMachine{States: states, CurrentState: 1}
	return PropStateMachine

}

func MoveableUpdater(ent *Entity, gs GameState) {
	p := ent.PropData

	placementEnt, exists := ent.GetLinkedEnt()
	if !exists {
		log.Fatal("placement ent id for prop at moveable state not working")
	}

	ent.Sprite.X = placementEnt.Sprite.X
	ent.Sprite.Y = placementEnt.Sprite.Y - float32(placementEnt.Sprite.SpriteHeight()) - 100

	newCorners := updateCorners(p.baseCorners, p.cornerOffsets, ent.Sprite.X, ent.Sprite.Y)
	p.baseCorners = newCorners

	//ent.Z = min(z, 12)

	/*var lastSprite *sprite.Sprite
	if p.LinkedSprite == nil {
		for _, corner := range p.baseCorners {
			placeMentSprite := &sprite.Sprite{Img: PM.placementReticule, X: float32(corner.X), Y: float32(corner.Y)}
			if lastSprite != nil {
				lastSprite.LinkedSprite = placeMentSprite
				lastSprite = placeMentSprite
			} else {
				p.LinkedSprite = placeMentSprite
				lastSprite = p.LinkedSprite
			}
		}
	}*/

	good := true
	for _, corner := range p.baseCorners {

		for _, otherProp := range PM.placedProps {
			if pointInPolygon(corner, otherProp.baseCorners) {
				good = false
			}
		}

		if !good {
			p.ShaderParams["OutlineColor"] = [4]float32{0.7, 0.2, 0.2, 1.0}
		} else {
			p.ShaderParams["OutlineColor"] = [4]float32{0.2, 0.7, 0.2, 1.0}
		}

	}

	if registry.ClickCheck() && p.CheckPropPointsValid() {
		ent.Transition()
	}

	_, _, ent.Z = PositionPointOnZ(int(placementEnt.Sprite.X), int(placementEnt.Sprite.Y)+placementEnt.Sprite.SpriteHeight(), gs.Zbounds)

}

func MakeBaseCornerPlacementSprite(baseCorners [4]image.Point, propImgSize image.Rectangle) *sprite.Sprite {
	top, bottom, left, right := findTopBottomCornersAndLeftRightCornerIndexes(baseCorners)

	img := ebiten.NewImage(baseCorners[right].X-baseCorners[left].X+(PM.placementReticule.Bounds().Dx()*2), baseCorners[top].Y-baseCorners[bottom].Y+PM.placementReticule.Bounds().Dy())
	dopts := &ebiten.DrawImageOptions{}

	for _, corner := range baseCorners {
		dopts.GeoM.Reset()
		x := float64(corner.X)
		y := float64(corner.Y - baseCorners[bottom].Y)

		dopts.GeoM.Translate(x, y)
		img.DrawImage(PM.placementReticule, dopts)
	}

	sp := &sprite.Sprite{Img: img}
	sp.DOptsUpdaterTag = "offset"
	sp.DOptsUpdaterParams = make(map[string]float64)
	sp.DOptsUpdaterParams["offSetX"] = -6
	sp.DOptsUpdaterParams["offSetY"] = -12

	return sp
}

func findTopBottomCornersAndLeftRightCornerIndexes(baseCorners [4]image.Point) (top int, bottom int, left int, right int) {

	topMostCornerIndex := 0
	bottomMostCornerIndex := 0

	leftMostCornerIndex := 0
	rightMostCornerIndex := 0

	for i := 1; i < 4; i++ {
		if baseCorners[i].X < baseCorners[leftMostCornerIndex].X {
			leftMostCornerIndex = i
		}
		if baseCorners[i].X > baseCorners[rightMostCornerIndex].X {
			rightMostCornerIndex = i
		}

		if baseCorners[i].Y < baseCorners[topMostCornerIndex].Y {
			bottomMostCornerIndex = i
		}
		if baseCorners[i].Y > baseCorners[bottomMostCornerIndex].Y {
			topMostCornerIndex = i
		}
	}

	return topMostCornerIndex, bottomMostCornerIndex, leftMostCornerIndex, rightMostCornerIndex
}

func SettingInPlaceUpdaterNew(ent *Entity, gs GameState) {
	//this is really a transition func but we need game state for z bounds, feels weird to break either pattern:
	//making entity store global state z bound data or having update func do one off transition stuff
	p := ent.PropData
	p.Sprite.UnLoadShader()
	p.Sprite.Shader = registry.ShaderMap["NormalMap"]

	p.LinkedSprite = nil
	spriteMidPoint := ent.Sprite.GetSpriteRect().Dx() / 2

	ent.AddDeInitHandler(DrawSpotLight(ent.Sprite.X+float32(spriteMidPoint), EffectParams{}))
	p.state = SettingInPlace

	placementEnt, exists := ent.GetLinkedEnt()
	if !exists {
		log.Fatal("placement ent id for prop at moveable state not working")
	}

	placementEnt.Draw = false

	for i, corner := range p.baseCorners {

		psZ := ent.Z
		psZ = min(psZ, len(gs.Zbounds)-2)
		psZ = max(psZ, 1)
		if i > 1 {
			psZ = min(psZ+2, len(gs.Zbounds)-2)
		}

		nps := NewBubbleSystem(float64(corner.X-gs.Zbounds[psZ].Min.X), float64(corner.Y-gs.Zbounds[psZ].Min.Y), gs.Zbounds[psZ])
		nps.SpawnRate = 5000
		nps.On = true
		ent.PropData.particleSystems = append(ent.PropData.particleSystems, nps)
		RegisterEntity(&Entity{ParticleSystem: nps, Z: psZ, EndAfter: 10.0, Sprite: nps.Sprite})

	}

	psZ := ent.Z
	psZ = min(psZ, len(gs.Zbounds)-2)
	psZ = max(psZ-2, 1)
	nps := NewBubbleSystem(float64(p.baseCorners[0].X-gs.Zbounds[psZ].Min.X+p.Sprite.GetSpriteRect().Dx()/2+rand.Intn(15)), float64(p.baseCorners[0].Y-gs.Zbounds[psZ].Min.Y), gs.Zbounds[psZ])
	nps.SpawnRate = 5000
	nps.On = true
	ent.PropData.particleSystems = append(ent.PropData.particleSystems, nps)
	RegisterEntity(&Entity{ParticleSystem: nps, Z: psZ, EndAfter: 10.0, Sprite: nps.Sprite})

	nps2 := NewBubbleSystem(float64(p.baseCorners[0].X-gs.Zbounds[psZ].Min.X+p.Sprite.GetSpriteRect().Dx()/4-rand.Intn(15)), float64(p.baseCorners[0].Y-gs.Zbounds[psZ].Min.Y), gs.Zbounds[psZ])
	nps2.SpawnRate = 5000
	nps2.On = true
	ent.PropData.particleSystems = append(ent.PropData.particleSystems, nps2)
	RegisterEntity(&Entity{ParticleSystem: nps, Z: psZ, EndAfter: 10.0, Sprite: nps2.Sprite})
	ent.Transition()
}

func PropFallingUpdater(ent *Entity, gs GameState) {

	p := ent.PropData
	p.Y++

	newCorners := updateCorners(p.baseCorners, p.cornerOffsets, ent.Sprite.X, ent.Sprite.Y)
	p.baseCorners = newCorners

	placementEnt, exists := ent.GetLinkedEnt()
	if !exists {
		log.Fatal("placement ent id for prop at moveable state not working")
	}
	//ent.Sprite.Y = placementEnt.Sprite.Y - float32(placementEnt.Sprite.SpriteHeight()) - 100

	if ent.Sprite.Y+float32(ent.Sprite.SpriteHeight()) > placementEnt.Sprite.Y-float32(placementEnt.Sprite.DOptsUpdaterParams["offSetY"]) {
		RemoveEntity(ent.LinkedID)
		ent.Transition()
	}

	_, _, ent.Z = PositionPointOnZ(int(placementEnt.Sprite.X), int(placementEnt.Sprite.Y)+placementEnt.Sprite.SpriteHeight(), gs.Zbounds)

}

func TransitionToSetInPlace(ent *Entity, gs GameState) {

	p := ent.PropData

	for _, ps := range ent.PropData.particleSystems {
		ps.On = false
	}

	if !p.alreadyPlaced {
		ent.EventHub.Publish(events.NewProp{PropId: ent.Id, Name: p.Tag})
	}

	ent.DeInitEffects() // turn off spotlight effect

	psz := max(ent.Z, 0)
	psz = min(psz, 12) // Fix thent.min function

	for _, corner := range p.baseCorners {
		nps := NewGenericParticleSystem(float64(corner.X), float64(corner.Y), gs.Zbounds[psz], 0)
		nps.PConfig = &ParticleConfig{
			XVariance:         10,
			YVariance:         10,
			XVelocityVariance: 10,
			YVelocityVariance: 10,
			MaxLife:           0,
			BaseYVelocity:     -40,
			RotationSpeed:     rand.NormFloat64() * .01,
			Scale:             1.5,
			AlphaDecay:        0.5}
		RegisterEntity(&Entity{ParticleSystem: nps, Z: psz, EndAfter: 10.0, Sprite: nps.Sprite})
	}
	ZSortEntities()
	p.state = SetInPlace
	if p.OnLanding != nil {
		p.OnLanding(ent)
	}

	PM.placedProps = append(PM.placedProps, ent.PropData)
	p.alreadyPlaced = true
	ent.Transition()
}

func SetInPlaceUpdater(ent *Entity, gs GameState) {
	if gs.Debug == "DebugOn" {
		return
	}

	if gs.ZoomedFor == NotZoomed {
		return
	}

	if ent.Parameters.Ints[IndexCounter] < 120 {
		ent.Parameters.Ints[IndexCounter]++
		return
	}

	if registry.ClickCheck() && ent.Sprite.SpriteHovered() {
		ent.Sprite.AddColoredOutlineShader(colornames.Yellow)
		ent.Transition()
	}
}

func SetInPlaceSelected(ent *Entity, gs GameState) {
	if registry.ClickCheck() && ent.Sprite.SpriteHovered() {

		filtered := PM.placedProps[:0]

		for _, prop := range PM.placedProps {
			if prop != ent.PropData {
				filtered = append(filtered, prop)
			}
		}

		PM.placedProps = filtered

		sp := MakeBaseCornerPlacementSprite(ent.PropData.cornerOffsets, ent.Sprite.Img.Bounds())
		sp.X = ent.Sprite.X
		sp.Y = ent.Sprite.Y + float32(ent.Sprite.SpriteHeight()) + 100
		ent3 := &Entity{Sprite: sp, Z: 10}
		ent3.UpdateFunc = MoveRectEnt
		RegisterEntity(ent3)
		ent.LinkedID = ent3.Id
		ent.Transition()
	}
}

func AddEffectOnLanding(ent *Entity) {
	DrawSmoke(ent.Sprite.X, ent.Sprite.Y, EffectParams{Cycles: 10})
}

func LoadProp(propName string, pd PropData, eventhub *tasks.EventHub, event tasks.Event, zbounds [13]image.Rectangle) uint32 {
	propAssets, err := LoadPropImgs(propName)
	if err != nil {
		log.Fatal(err)
	}

	var castleCollisions []image.Rectangle
	var points []image.Point

	normalizedName := strings.ToLower(propName)

	for key, collision := range pd.CollisionMap {
		normalizedKey := strings.ToLower(key)
		if strings.Contains(normalizedKey, normalizedName) {
			castleCollisions = append(castleCollisions, collision)
			println("collision found:", key)
		}
	}

	for key, pts := range pd.PtMap {
		normalizedKey := strings.ToLower(key)
		if strings.Contains(normalizedKey, normalizedName) {
			points = pts
			println("point found:", key)
		}
	}

	prop := NewStructureProp(propAssets, eventhub, zbounds[0], propName)
	if len(points) < 4 {
		points = pd.PtMap["Castle"]
	}

	prop.cornerOffsets = [4]image.Point(points[0:4])

	x, y, _ := positionPointBasedOnCursorOnZslice(zbounds)
	for i := 0; i <= 3; i++ {
		prop.baseCorners[i].X = prop.cornerOffsets[i].X + x
	}
	for i := 0; i <= 3; i++ {
		prop.baseCorners[i].Y = prop.cornerOffsets[i].Y + y
	}

	prop.Sprite.X = float32(x)
	prop.Sprite.Y = float32(prop.baseCorners[0].Y - 45 - prop.Img.Bounds().Dy())

	if prop.Sprite2 != nil {
		prop.Sprite2.X = float32(x)
		prop.Sprite2.Y = prop.Sprite.Y
	}

	if propName == string(HotRock) {
		prop.OnLanding = AddEffectOnLanding
	}

	ent := &Entity{PropData: prop, Sprite: prop.Sprite, EventHub: eventhub}
	ent.StateMachine = InitPropStateMachine()
	ent.EventHub = eventhub
	ent.Z = 6
	RegisterEntity(ent)

	if prop.Sprite2 != nil {
		ent2 := &Entity{Sprite: prop.Sprite2, EventHub: eventhub}
		ent2.Z = 0
		RegisterEntity(ent2)
	}

	eventhub.Publish(events.PlacementMode{})

	sp := MakeBaseCornerPlacementSprite(ent.PropData.cornerOffsets, ent.Sprite.Img.Bounds())
	sp.X = 400
	sp.Y = 400
	ent3 := &Entity{Sprite: sp, Z: 10}
	ent3.UpdateFunc = MoveRectEnt
	RegisterEntity(ent3)

	selectRepositionChain := func(ent *Entity, gs GameState) {
		if registry.ClickCheck() {
			RemoveEntity(ent.Id)
			RemoveEntity(ent.Parameters.EntIds[LinkedGraphic1])
			DrawControlEffect(float64(ent.Sprite.X), float64(ent.Sprite.Y), EffectParams{Zoom: true}, ClickHere, "Reposition", ClearEnterGraphicOnClick)
		}
	}

	if pd.SubscriptionForNextProp != nil {
		ent.SubWUnsubAfterCompletion(pd.SubscriptionForNextProp, func(event tasks.Event) {
			effectX := float64(ent.Sprite.X) + float64(ent.Sprite.SpriteWidth()/2)
			effectY := float64(ent.Sprite.Y) + float64(ent.Sprite.SpriteHeight()/2)
			DrawControlEffect(effectX, effectY, EffectParams{Zoom: true}, ClickHere, "Select", selectRepositionChain)
		})
	}

	ent.LinkedID = ent3.Id

	return ent.Id
}

func plantUpdateFunc(ent *Entity, gs GameState) {
	if ent.Sprite.CurrentAnimation == "StartUp" && ent.AnimationCycles == 1 {
		ani := ent.Sprite.GetAnimation()
		ent.Sprite.Img = ani.GetLastFrameAsStaticImage()
		ent.Sprite.NormalMap = ani.GetLastFrameNormalAsStaticImage()
		ent.Sprite.CurrentAnimation = ""
		ent.UpdateFunc = movePlantUpdateFunc
		sprite.InitSwayAnimation(ent.Sprite, 15)
		ent.EventHub.Publish(events.NewProp{PropId: ent.Id})
	}
}

func movePlantUpdateFunc(ent *Entity, gs GameState) {
	if registry.ClickCheck() && ent.Sprite.SpriteHovered() && gs.ZoomedFor == PlayerZoomed {
		ent.Sprite.AddColoredOutlineShader(colornames.Yellow)
		ent.Parameters.Ints[lastZ] = ent.Z
		ent.UpdateFunc = MoveEnt
	}
}

func MoveEnt(ent *Entity, gs GameState) {
	x := ent.Sprite.X
	y := ent.Sprite.Y
	spriteBounds := ent.Sprite.GetSpriteRect()

	newX, newY, newZ := MovePoint(x, y, spriteBounds, gs, ent.Parameters.Ints[lastZ])

	ent.Sprite.X = newX
	ent.Sprite.Y = newY
	ent.Parameters.Ints[lastZ] = newZ
	ent.Z = newZ
}

func MoveRectEnt(ent *Entity, gs GameState) {
	x := ent.Sprite.X
	y := ent.Sprite.Y

	spriteBounds := ent.Sprite.GetSpriteRect()
	if ent.Sprite.DOptsUpdaterTag == "offset" {
		spriteBounds.Max.Y += int(ent.Sprite.DOptsUpdaterParams["offSetY"] * 1.5) // arbitrary scalar based on what looks good with position reticule
	}

	entslastZ := ent.Z

	newX, newY, newZ := MovePoint(x, y, spriteBounds, gs, entslastZ)

	ent.Sprite.X = newX
	ent.Sprite.Y = newY
	ent.Parameters.Ints[lastZ] = newZ

	if ent.Sprite.Y+6 < float32(gs.Zbounds[0].Max.Y) {
		ent.Sprite.Y = float32(gs.Zbounds[0].Max.Y) - 6
	}

}

func MovePoint(x float32, y float32, objectBeingMovedBounds image.Rectangle, gs GameState, lastZ int) (float32, float32, int) {
	Dy := 0.0
	Dx := 0.0

	width := objectBeingMovedBounds.Dx()
	height := objectBeingMovedBounds.Dy()

	if registry.LeftCheck() {
		Dx -= 1
	}
	if registry.RightCheck() {
		Dx += 1
	}
	if registry.UpCheck() {
		Dy -= 1
	}
	if registry.DownCheck() {
		Dy += 1
	}

	_, _, z := PositionPointOnZ(int(x), int(y)+height, gs.Zbounds)

	if z != lastZ {
		slices := SliceRectangle(gs.Zbounds[z], 5)
		positionIndex := FindSliceIndex(x, slices)
		switch positionIndex {
		case 0:
			if Dy > 0 {
				Dx -= 1
			} else if Dy < 0 {
				Dx += 1
			}
		case 1:
			if Dy > 0 {
				Dx -= 0.5
			} else if Dy < 0 {
				Dx += 0.5
			}
		case 2:
			//keep change
		case 3:
			if Dy > 0 {
				Dx += 0.5
			} else if Dy < 0 {
				Dx -= 0.5
			}
		case 4:
			if Dy > 0 {
				Dx += 1
			} else if Dy < 0 {
				Dx -= 1
			}
		}
	}

	x += float32(Dx)
	y += float32(Dy)

	upperBounds := float32(gs.Zbounds[12].Max.Y - height)
	leftBounds := float32(gs.Zbounds[z].Min.X)
	rightBounds := float32(gs.Zbounds[z].Max.X - width)

	if y > upperBounds {
		y = upperBounds
	}
	/*if ent.Sprite.Y < lowerBounds {
		ent.Sprite.Y = lowerBounds
	}*/
	if x < leftBounds {
		x = leftBounds
	}
	if x > rightBounds {
		x = rightBounds
	}
	if z == 0 {
		lowerBounds := float32(gs.Zbounds[0].Max.Y - height)
		if y < lowerBounds {
			y = lowerBounds
		}
	}

	return x, y, z

}

func SliceRectangle(rect image.Rectangle, numSlices int) []float32 {
	if numSlices <= 0 {
		return []float32{}
	}

	width := float32(rect.Dx())
	startPoint := float32(rect.Min.X)
	sliceWidth := width / float32(numSlices)

	slices := make([]float32, numSlices+1)
	slices[0] = startPoint

	for i := 1; i <= numSlices; i++ {
		slices[i] = startPoint + sliceWidth*float32(i)
	}

	return slices
}

func FindSliceIndex(x float32, sliceBounds []float32) int {
	if len(sliceBounds) == 0 {
		return -1
	}

	// Before first boundary
	if x < sliceBounds[0] {
		return -1
	}

	// After last boundary
	if x >= sliceBounds[len(sliceBounds)-1] {
		return len(sliceBounds) - 1
	}

	// Find which slice the x falls into
	for i := 0; i < len(sliceBounds)-1; i++ {
		if x >= sliceBounds[i] && x < sliceBounds[i+1] {
			return i
		}
	}

	return -1
}

func LoadPlant(event tasks.Event, data PropData, hub *tasks.EventHub) uint32 {
	ev := event.(PlacementPicked)
	var plantName string
	pchance := rand.Intn(3)
	if pchance < 1 {
		plantName = "leafyPlant"
	} else if pchance < 2 {
		plantName = "grassyPlant"
	} else {
		plantName = "fern"
	}

	rareChance := rand.Intn(10)

	fert, ok := data.PlacementParams["fertilizer"].(bool)
	if ok {
		if fert {
			if rareChance < 7 {
				plantName += "Rare"
			}
		}
	} else {
		if rareChance < 2 {
			plantName += "Rare"
		}
	}

	pAni := LoadPlantAnimation(plantName)
	pImg := pAni.Img
	x := +ev.X - float32(pImg.Bounds().Dx()/(pAni.LastF+1))/2
	y := ev.Y - float32(pImg.Bounds().Dy())

	sp := &sprite.Sprite{AnimationMap: map[string]*sprite.Animation{"StartUp": pAni}, CurrentAnimation: "StartUp", X: x, Y: y}
	sp.UnFocusable = true
	sp.ShaderParams = make(map[string]any)
	sp.ShaderParams["Cursor"] = []float64{0, 0, 100}
	if plantName[len(plantName)-4:] == "Rare" {
		sp.Shader = registry.ShaderMap["NormalMapOutline"]
		switch plantName[:len(plantName)-4] {
		case "fern":
			sp.ShaderParams["OutlineColor"] = [4]float32{0.7, 0.1, 0.5, 0.07}
		case "leafyPlant":
			sp.ShaderParams["OutlineColor"] = [4]float32{0.7, 0.1, 0.2, 0.07}
		case "grassyPlant":
			sp.ShaderParams["OutlineColor"] = [4]float32{0.1, 0.6, 0.7, 0.07}
		}

	}

	ent := &Entity{Sprite: sp}
	ent.Z = max(1, ev.Z)
	ent.EventHub = hub
	ent.UpdateFunc = plantUpdateFunc
	println("Z picked for plant=", ent.Z)

	RegisterEntity(ent)
	pps := NewGenericParticleSystem(float64(ev.X), float64(ev.Y), data.ZBounds[ev.Z], 0)
	pps.PConfig = &ParticleConfig{
		XVariance:         10,
		YVariance:         10,
		XVelocityVariance: 10,
		YVelocityVariance: 10,
		MaxLife:           0,
		BaseYVelocity:     -50,
		RotationSpeed:     rand.NormFloat64() * .01,
		Scale:             1.0,
		AlphaDecay:        0.5}
	RegisterEntity(&Entity{ParticleSystem: pps, EndAfter: 10.0, Z: ev.Z, Sprite: pps.Sprite})

	pps2 := NewGenericParticleSystem(float64(ev.X), float64(ev.Y), data.ZBounds[ev.Z], 0)
	pps2.PConfig = &ParticleConfig{
		XVariance:         10,
		YVariance:         10,
		XVelocityVariance: 10,
		YVelocityVariance: 10,
		MaxLife:           0,
		BaseYVelocity:     -50,
		RotationSpeed:     rand.NormFloat64() * .01,
		Scale:             1.0,
		AlphaDecay:        0.5}
	RegisterEntity(&Entity{ParticleSystem: pps2, EndAfter: 10.0, Z: ev.Z, Sprite: pps2.Sprite})

	return ent.Id
}

func LoadPlantAnimation(plantName string) *sprite.Animation {

	var pImg *ebiten.Image
	var baseName = plantName

	if strings.Contains(plantName, "Rare") {
		baseName = plantName[0 : len(plantName)-4]
		img, err := util.LoadImageAssetAsEbitenImage(fmt.Sprintf("tankProps/plants/%s", plantName))
		if err != nil {
			println(" no rare found for:", plantName)
			img, err = util.LoadImageAssetAsEbitenImage(fmt.Sprintf("tankProps/plants/%s", baseName))
			if err != nil {
				log.Fatal(err)
			}
		}
		pImg = img
	} else {
		img, err := util.LoadImageAssetAsEbitenImage(fmt.Sprintf("tankProps/plants/%s", plantName))
		if err != nil {
			log.Fatal(err)
		}
		pImg = img
	}

	pAni, err := entImportableLoaders.LoadAnimation(fmt.Sprintf("data/animationData/%sAnimation.json", baseName))
	if err != nil {
		log.Fatal(err)
	}

	pNormal, err := util.LoadImageAssetAsEbitenImage(fmt.Sprintf("tankProps/plants/%s_n", baseName))
	if err != nil {
		log.Printf("no normal map loaded for:%s error:%s", plantName, err)
	}
	pAni.Img = pImg
	if pNormal != nil {
		pAni.NormalImg = pNormal
	}
	return pAni
}

func LoadPlacementImg() (*ebiten.Image, *ebiten.Image) {
	img, err := util.LoadImageAssetAsEbitenImage("uiSprites/placementReticule")
	if err != nil {
		log.Fatal("placement reticule path is bad", err)
	}
	img2, err2 := util.LoadImageAssetAsEbitenImage("uiSprites/placementReticuleInvalid")
	if err2 != nil {
		log.Fatal("placement reticule path is bad", err)
	}
	return img, img2
}

func positionPointBasedOnCursorOnZslice(zBounds [13]image.Rectangle) (int, int, int) {
	x, y := util.GetScaledCursorPosition()
	x -= 6
	y -= 5
	xMod, yMod, z := PositionPointOnZ(x, y, zBounds)
	return xMod, yMod, z
}

func PositionPointOnZ(x, y int, zBounds [13]image.Rectangle) (int, int, int) {
	currentZ := 0
	maxY := zBounds[12].Max.Y
	minY := zBounds[0].Max.Y

	if y >= minY && y <= maxY {
		rangept := (y - minY) / 2
		currentZ = rangept
	}

	if y > maxY {
		currentZ = 10
		y = maxY
	}

	if y < minY {
		currentZ = 0
		y = minY
	}

	maxX := zBounds[currentZ].Max.X
	if x > maxX {
		x = maxX
	}

	minX := zBounds[currentZ].Min.X
	if x < minX {
		x = minX
	}

	currentZ = max(currentZ, 11)
	return x, y, currentZ
}

func updateCorners(baseCorners [4]image.Point, cornerOffsets [4]image.Point, x, y float32) [4]image.Point {
	for i, _ := range baseCorners {
		baseCorners[i].X = int(x) + cornerOffsets[i].X
		baseCorners[i].Y = int(y) + cornerOffsets[i].Y
	}

	return baseCorners
}

func pointInPolygon(point image.Point, polygon [4]image.Point) bool {
	x, y := point.X, point.Y
	n := len(polygon)
	inside := false

	j := n - 1
	for i := 0; i < n; i++ {
		xi, yi := polygon[i].X, polygon[i].Y
		xj, yj := polygon[j].X, polygon[j].Y

		if ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}
	return inside
}

func (p *StructureProp) CheckPropPointsValid() bool {
	for _, point := range p.baseCorners {
		for _, otherProp := range PM.placedProps {
			if pointInPolygon(point, otherProp.baseCorners) {
				return false
			}
		}
	}
	return true
}
