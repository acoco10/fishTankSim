package entities

import (
	"github.com/ByteArena/box2d"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/stringConstants"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	physics "github.com/acoco10/fishTankWebGame/game/testPhysics"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
	"math/rand"
)

type IndividualParticleEntity uint8

const (
	Debris IndividualParticleEntity = iota
	Food
)

var n int = 0

type EntityParticle struct {
	body *box2d.B2Body
	*util.Point
	image *ebiten.Image
	IndividualParticleEntity
	tankBounds        image.Rectangle
	justSpawned       bool
	counter           int
	UnderWater        bool
	waterLevel        int
	floorLevel        int
	underWaterCounter int
	bounds            *image.Rectangle
	remove            bool
	eventHub          *tasks.EventHub
	stop              bool
	rotationSpeed     float64
	baseVy            float64
	targX             float32
	targY             float32
	*sprite.Sprite
}

func (p *EntityParticle) ShouldRemove() bool {
	return p.remove
}

func fallingParticleUpdater(ent *Entity, gs GameState) {
	p := ent.ParticleData
	p.counter++
	vy := p.baseVy
	dx := (rand.Float32() - 0.5) * 0.1
	if p.IndividualParticleEntity == Food {
		dx = float32(-5.0)
		if p.Point.Tag == "left" {
			dx = float32(5)
		}
	}

	p.Point.Y += float32(vy)
	p.Point.X += dx

	if p.GetSpriteRect().In(p.tankBounds) {
		ent.StateMachine.Transition(ent)
	}

	p.Sprite.X = p.Point.X
	p.Sprite.Y = p.Point.Y
}

func captureUpdater(ent *Entity, gs GameState) {
	p := ent.ParticleData
	p.counter++

	if p.bounds == nil || p.Point == nil {
		ent.StateMachine.Transition(ent)
		return
	}

	newPoint := restrictTargetPointWithinBounds(p.Sprite.GetSpriteRect().Dx(), p.Sprite.GetSpriteRect().Dy(), *p.Point, *p.bounds)
	p.Point = &newPoint
	p.rotationSpeed = 0

	p.Sprite.X = p.Point.X
	p.Sprite.Y = p.Point.Y

}

func box2dUpdater(ent *Entity, gs GameState) {
	posVec := ent.ParticleData.body.GetPosition()
	ent.Sprite.X = float32(posVec.X) * 50
	ent.Sprite.Y = float32(posVec.Y) * 50

	if ent.ParticleData.bounds != nil {
		ent.StateMachine.Transition(ent)
	}

}

func underWaterUpdater(ent *Entity, gs GameState) {

	p := ent.ParticleData
	p.counter++

	if p.Point.Y >= float32(p.floorLevel) {
		if p.PType == util.Food {
			return
		}
		p.underWaterCounter = 0
		p.baseVy = -2
	}

	if p.stop {
		return
	}

	vy := p.baseVy
	dx := (rand.Float32() - 0.5) * 0.1

	dx = -0.01
	p.underWaterCounter++
	if p.baseVy > 0 {
		vy -= 2 * float64(p.underWaterCounter)
		if p.PType == util.Food {
			vy = max(vy, 0.08)
		} else {
			vy = max(vy, 0.5)
		}
	} else {
		vy += 0.02 * float64(p.underWaterCounter)
	}

	if p.counter%5 == 0 {
		vx := math.Sin(float64(p.counter)*0.5) * 0.3 * 1
		noise := rand.Float64()*0.1 - 0.05
		p.Point.X = p.Point.X + float32(vx+noise)
	}

	p.Point.Y += float32(vy)
	p.Point.X += dx

	p.Sprite.DOptsUpdaterParams["degree"] += p.rotationSpeed

	p.Sprite.X = p.Point.X
	p.Sprite.Y = p.Point.Y

	if p.bounds != nil {
		ent.StateMachine.Transition(ent)
	}
}

func transitionToUnderWater(ent *Entity) {
	p := ent.ParticleData
	if p.IndividualParticleEntity == Food {
		ev := SendData{Data: "particle entered water",
			DataFor: "soundFx"}
		p.eventHub.Publish(ev)
	}
}

func (ent *Entity) StopParticle() {
	ent.ParticleData.stop = true
}

func NewEntityParticle(point *util.Point, rect image.Rectangle, hub *tasks.EventHub, zbounds [13]image.Rectangle,
	ePType IndividualParticleEntity, world *box2d.B2World) *Entity {
	println("calling new particle function", n)
	n++
	pImg, _ := util.LoadImageAssetAsEbitenImage("textures/foodParticles")

	p := &EntityParticle{
		Point:                    point,
		IndividualParticleEntity: ePType,
		counter:                  0,
		UnderWater:               false,
		waterLevel:               rect.Min.Y,
		floorLevel:               rect.Max.Y,
		tankBounds:               rect,
		underWaterCounter:        0,
		justSpawned:              true,
		eventHub:                 hub,
		remove:                   false,
		rotationSpeed:            rand.NormFloat64() * .001,
		Sprite:                   &sprite.Sprite{Img: particleImg},
	}

	sm := StateMachine{}
	fallingUpdater := &StateHandler{Updater: fallingParticleUpdater, TransitionTo: 2, TransitionOutFunc: transitionToUnderWater}
	underWaterStateUpdater := &StateHandler{Updater: underWaterUpdater, TransitionTo: 3, TransitionOutFunc: nil}
	capturedStateUpdater := &StateHandler{Updater: captureUpdater, TransitionTo: 1, TransitionOutFunc: nil}

	sm.CurrentState = 1
	sm.States = map[int]*StateHandler{
		1: fallingUpdater,
		2: underWaterStateUpdater,
		3: capturedStateUpdater,
	}

	z := rand.Intn(7)
	entity := &Entity{ParticleData: p, Sprite: p.Sprite, StateMachine: &sm}
	entity.Z = 5 + z

	switch ePType {
	case Food:
		GetFoodTextureImgAndZ(pImg, entity)
		p.rotationSpeed = 0.0

	case Debris:
		GetDebrisTextureImgAndZSpeed(pImg, entity)
		p.baseVy += rand.NormFloat64() * 5
		p.Dx = rand.Float32() - 0.5
	}

	entity.ParticleData.floorLevel = zbounds[z].Max.Y
	entity.ParticleData.Point.Z = entity.Z

	entity.Sprite.DOptsUpdaterTag = stringConstants.DepthFilter
	entity.Sprite.DOptsUpdaterParams = make(map[string]float64)
	entity.Sprite.DOptsUpdaterParams["Z"] = float64(entity.Z)

	entity.EventHub = hub
	RegisterEntity(entity)

	pointEvent := PointGenerated{PointId: entity.Id, Source: "new particle function"}
	p.eventHub.Publish(pointEvent)

	return entity
}

func MakeDebrisBody(p *EntityParticle, world *box2d.B2World) {
	body := physics.MakeDebris(box2d.MakeB2Vec2(float64(p.Point.X), float64(p.Point.Y)), world)
	randomForceX := (rand.Float64() - 0.5) * .05  // -1.0 to 1.0
	randomForceY := (rand.Float64() - 0.5) * 0.01 // Small vertical variance

	impulse := box2d.MakeB2Vec2(randomForceX, randomForceY)

	// Apply impulse at center of mass
	body.ApplyLinearImpulseToCenter(impulse, true)
	p.body = body

	underWaterStateUpdater1 := &StateHandler{Updater: box2dUpdater, TransitionTo: 3, TransitionOutFunc: nil}
	sm := StateMachine{}
	sm.States[1] = underWaterStateUpdater1
	capturedStateUpdater := &StateHandler{Updater: captureUpdater, TransitionTo: 1, TransitionOutFunc: nil}

	sm.States[2] = capturedStateUpdater

	sm.States[2].TransitionTo = 1
}

func GetDebrisTextureImgAndZSpeed(pImg *ebiten.Image, debrisEntity *Entity) {
	if debrisEntity.Z >= 8 {
		debrisEntity.ParticleData.baseVy = 10
		debrisEntity.Sprite.Img = pImg.SubImage(image.Rect(0, 6, 18, 13)).(*ebiten.Image)
	} else if debrisEntity.Z >= 6 {
		debrisEntity.ParticleData.baseVy = 9.8
		debrisEntity.Sprite.Img = pImg.SubImage(image.Rect(0, 6, 11, 13)).(*ebiten.Image)
	} else {
		debrisEntity.ParticleData.baseVy = 9.6
		debrisEntity.Sprite.Img = pImg.SubImage(image.Rect(0, 6, 11, 13)).(*ebiten.Image)
	}
}

func GetFoodTextureImgAndZ(pImg *ebiten.Image, foodEntity *Entity) {
	//maps 3 different textures to different depths and applies an offset to base velocity based on depth and mutates particle
	if foodEntity.Z >= 8 {
		foodEntity.ParticleData.baseVy = 10
		foodEntity.Sprite.Img = pImg.SubImage(image.Rect(0, 0, 6, 6)).(*ebiten.Image)
	} else if foodEntity.Z >= 6 {
		foodEntity.ParticleData.baseVy = 9.8
		foodEntity.Sprite.Img = pImg.SubImage(image.Rect(6, 0, 12, 5)).(*ebiten.Image)
	} else {
		foodEntity.ParticleData.baseVy = 9.6
		foodEntity.Sprite.Img = pImg.SubImage(image.Rect(12, 0, 16, 4)).(*ebiten.Image)

	}
}
