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
	placedProps              []StructureProp
	placementReticule        *ebiten.Image
	invalidPlacementReticule *ebiten.Image
}

type PropAssets struct {
	normalMap *ebiten.Image
	layers    []*ebiten.Image
	merged    *ebiten.Image
}

type PropData struct {
	CollisionMap    map[string]image.Rectangle
	PtMap           map[string][]image.Point
	PlacementParams map[string]any
	ZBounds         [13]image.Rectangle
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
	cornerOffsets   [4]image.Point
	baseCorners     [4]image.Point //0-4 leftUpperCorner rightUpperCorner leftLowerCorner rightLowerCorner
	particleSystems []*ParticleSystem
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

func (e *Entity) UpdateProp(zBounds [13]image.Rectangle) {
	p := e.PropData

	if p.state == Moveable {
		newCorners, z := updateCorners(p.baseCorners, zBounds)
		p.baseCorners = newCorners
		p.Sprite.Y = float32(p.baseCorners[0].Y-p.cornerOffsets[0].Y) - 45
		p.Sprite.X = float32(p.baseCorners[0].X - p.cornerOffsets[0].X)
		e.Z = min(z, 12)
		e.Z = max(1, z)

		var lastSprite *sprite.Sprite
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
		}
		ls := p.LinkedSprite
		good := true
		for _, corner := range p.baseCorners {
			ls.X = float32(corner.X)
			ls.Y = float32(corner.Y)
			for _, otherProp := range PM.placedProps {
				if pointInPolygon(corner, otherProp.baseCorners) {
					good = false
					ls.Img = PM.invalidPlacementReticule
				} else {
					ls.Img = PM.placementReticule
				}
			}
			if !good {
				p.ShaderParams["OutlineColor"] = [4]float32{0.7, 0.2, 0.2, 1.0}
			} else {
				p.ShaderParams["OutlineColor"] = [4]float32{0.2, 0.7, 0.2, 1.0}
			}

			ls = ls.LinkedSprite

		}

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && p.CheckPropPointsValid() {

			p.LinkedSprite = nil

			spriteMidPoint := e.Sprite.GetSpriteRect().Dx() / 2

			DrawSpotLight(e.Sprite.X + float32(spriteMidPoint))
			p.state = SettingInPlace

			for i, corner := range p.baseCorners {

				psZ := e.Z
				psZ = min(psZ, len(zBounds)-1)
				psZ = max(psZ, 1)
				if i > 1 {
					psZ = min(psZ+2, len(zBounds)-1)
				}

				nps := NewBubbleSystem(float64(corner.X-zBounds[psZ].Min.X), float64(corner.Y-zBounds[psZ].Min.Y), zBounds[psZ])
				nps.SpawnRate = 5000
				nps.On = true
				e.PropData.particleSystems = append(e.PropData.particleSystems, nps)
				RegisterEntity(&Entity{ParticleSystem: nps, Z: psZ, EndAfter: 10.0, Sprite: nps.Sprite})

			}
			psZ := e.Z
			psZ = min(psZ, len(zBounds)-1)
			psZ = max(psZ-2, 1)
			nps := NewBubbleSystem(float64(p.baseCorners[0].X-zBounds[psZ].Min.X+p.Sprite.GetSpriteRect().Dx()/2+rand.Intn(15)), float64(p.baseCorners[0].Y-zBounds[psZ].Min.Y), zBounds[psZ])
			nps.SpawnRate = 5000
			nps.On = true
			e.PropData.particleSystems = append(e.PropData.particleSystems, nps)
			RegisterEntity(&Entity{ParticleSystem: nps, Z: psZ, EndAfter: 10.0, Sprite: nps.Sprite})

			nps2 := NewBubbleSystem(float64(p.baseCorners[0].X-zBounds[psZ].Min.X+p.Sprite.GetSpriteRect().Dx()/4-rand.Intn(15)), float64(p.baseCorners[0].Y-zBounds[psZ].Min.Y), zBounds[psZ])
			nps2.SpawnRate = 5000
			nps2.On = true
			e.PropData.particleSystems = append(e.PropData.particleSystems, nps2)
			RegisterEntity(&Entity{ParticleSystem: nps, Z: psZ, EndAfter: 10.0, Sprite: nps2.Sprite})

		}
	}

	if p.Sprite2 != nil {
		p.Sprite2.X = p.X
		p.Sprite2.Y = p.Y
	}

	//p.baseCorners = updateCorners(p.baseCorners, zBounds)
	if p.state == SettingInPlace {
		p.Sprite.UnLoadShader()
		p.Sprite.Shader = registry.ShaderMap["NormalMap"]
		p.Y++
	}
	if p.Y == float32(p.baseCorners[0].Y-p.cornerOffsets[0].Y) {
		if p.stateWas == SettingInPlace {
			for _, ps := range e.PropData.particleSystems {
				ps.On = false
			}
			e.EventHub.Publish(events.NewProp{PropId: e.Id, Name: p.Tag})
			TurnOffSpotLight()
			maxY := max(e.PropData.baseCorners[3].Y, e.PropData.baseCorners[2].Y)
			_, _, z := PositionPointOnZ(e.PropData.baseCorners[3].X, maxY, zBounds)
			e.Z = z

			psz := max(e.Z, 0)
			psz = min(psz, 12) // Fix the min function

			for _, corner := range p.baseCorners {
				nps := NewGenericParticleSystem(float64(corner.X), float64(corner.Y), zBounds[psz], 0)
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
			PM.placedProps = append(PM.placedProps, *e.PropData)
		}

	}

	p.stateWas = p.state

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
	if len(points) == 0 {
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

	ent := &Entity{PropData: prop, Sprite: prop.Sprite, EventHub: eventhub}
	ent.EventHub = eventhub
	ent.Z = 6
	RegisterEntity(ent)
	if prop.Sprite2 != nil {
		ent2 := &Entity{Sprite: prop.Sprite2, EventHub: eventhub}
		ent2.Z = 0
		ent.EventHub = eventhub
		RegisterEntity(ent2)
	}
	eventhub.Publish(events.PlacementMode{})

	return ent.Id
}

func plantUpdateFunc(ent *Entity) {
	if ent.Sprite.CurrentAnimation == "StartUp" && ent.AnimationCycles == 1 {
		ani := ent.Sprite.GetAnimation()
		ent.Sprite.Img = ani.GetLastFrameAsStaticImage()
		ent.Sprite.NormalMap = ani.GetLastFrameNormalAsStaticImage()
		ent.Sprite.CurrentAnimation = ""
		ent.UpdateFunc = nil
		sprite.InitSwayAnimation(ent.Sprite, 15)
		ent.EventHub.Publish(events.NewProp{PropId: ent.Id})

	}
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
	x := +ev.X + 2 - float32(pImg.Bounds().Dx()/(pAni.LastF+1))/2
	y := ev.Y - float32(pImg.Bounds().Dy()) + 2

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
	ent.Z = ev.Z
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

func LoadPlaceMentReticule(zBounds [13]image.Rectangle, tag string, hub *tasks.EventHub) {
	img, _ := LoadPlacementImg()
	ev := events.PlacementMode{}
	hub.Publish(ev)
	x, _, currentZ := positionPointBasedOnCursorOnZslice(zBounds)

	sp := &sprite.Sprite{Img: img, Y: float32(zBounds[currentZ].Max.Y), X: float32(x), UnFocusable: true}
	ent := &Entity{Sprite: sp}
	ent.Z = 13
	ent.UpdateFunc = placemenReticuleUpdateFunc
	RegisterEntity(ent)
	ent.Parameters["zBounds"] = zBounds
	ent.Parameters["Tag"] = tag
	ent.EventHub = hub

}

func positionPointBasedOnCursorOnZslice(zBounds [13]image.Rectangle) (int, int, int) {
	x, y := util.GetScaledCursorPosition()
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
	return x, y, currentZ
}

func placemenReticuleUpdateFunc(ent *Entity) {
	zBounds, ok := ent.Parameters["zBounds"].([13]image.Rectangle)
	if !ok {
		log.Fatal("incorrect params passed to reticule")
	}
	x, y, currentZ := positionPointBasedOnCursorOnZslice(zBounds)

	ent.Sprite.X = float32(x)
	ent.Sprite.Y = float32(y)
	ent.Sprite.DOptsUpdaterParams = make(map[string]float64)
	ent.Sprite.DOptsUpdaterTag = "offSet"
	ent.Sprite.DOptsUpdaterParams["offSetX"] = -2
	ent.Sprite.DOptsUpdaterParams["offSetY"] = -4

	valid := true
	for _, prop := range PM.placedProps {
		valid = !pointInPolygon(image.Point{int(ent.Sprite.X), int(ent.Sprite.Y)}, prop.baseCorners)
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && valid {
		tag, ok := ent.Parameters["Tag"].(string)
		if !ok {
			log.Fatal("placement reticule got non string Tag")
		}
		ent.EventHub.Publish(PlacementPicked{PlacementFor: tag, X: ent.Sprite.X, Y: ent.Sprite.Y, Z: currentZ})
		RemoveEntity(ent.Id)
	}

}

func updateCorners(baseCorners [4]image.Point, zBounds [13]image.Rectangle) ([4]image.Point, int) {
	x, y, currentZ := positionPointBasedOnCursorOnZslice(zBounds)

	// Calculate current center of the corners
	currentCenterX := (baseCorners[0].X + baseCorners[1].X + baseCorners[2].X + baseCorners[3].X) / 4
	currentCenterY := (baseCorners[0].Y + baseCorners[1].Y + baseCorners[2].Y + baseCorners[3].Y) / 4

	// Calculate how much to shift to center on cursor
	shiftX := x - currentCenterX
	shiftY := y - currentCenterY

	// Apply the shift to all corners
	for i := 0; i < 4; i++ {
		baseCorners[i].X += shiftX
		baseCorners[i].Y += shiftY
	}

	// Apply X boundary constraints (your existing logic)
	leftMostCornerIndex := 0
	rightMostCornerIndex := 0
	for i := 1; i < 4; i++ {
		if baseCorners[i].X < baseCorners[leftMostCornerIndex].X {
			leftMostCornerIndex = i
		}
		if baseCorners[i].X > baseCorners[rightMostCornerIndex].X {
			rightMostCornerIndex = i
		}
	}

	if baseCorners[leftMostCornerIndex].X < zBounds[currentZ].Min.X {
		shift := zBounds[currentZ].Min.X - baseCorners[leftMostCornerIndex].X
		for i := 0; i < 4; i++ {
			baseCorners[i].X += shift
		}
	}

	if baseCorners[rightMostCornerIndex].X > zBounds[currentZ].Max.X {
		shift := baseCorners[rightMostCornerIndex].X - zBounds[currentZ].Max.X
		for i := 0; i < 4; i++ {
			baseCorners[i].X -= shift
		}
	}

	// Add Y boundary constraints
	if baseCorners[0].Y < zBounds[0].Max.Y {
		shift := zBounds[0].Max.Y - baseCorners[0].Y
		for i := 0; i < 4; i++ {
			baseCorners[i].Y += shift
		}

	}

	if baseCorners[1].Y < zBounds[0].Max.Y {
		shift := zBounds[0].Max.Y - baseCorners[1].Y
		for i := 0; i < 4; i++ {
			baseCorners[i].Y += shift
		}
	}

	if baseCorners[2].Y > zBounds[11].Max.Y {
		currentZ = 2
		shift := baseCorners[2].Y - zBounds[11].Max.Y
		for i := 0; i < 4; i++ {
			baseCorners[i].Y -= shift
		}
	}

	if baseCorners[3].Y > zBounds[11].Max.Y {
		shift := baseCorners[3].Y - zBounds[11].Max.Y
		currentZ = 11
		for i := 0; i < 4; i++ {
			baseCorners[i].Y -= shift
		}
	}

	return baseCorners, max(currentZ+shiftY, currentZ)
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
