package entities

import (
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
	"math/rand"
	"time"
)

var n int = 0

type FoodParticle struct {
	*util.Point
	image             *ebiten.Image
	counter           int
	underWater        bool
	waterLevel        int
	floorLevel        int
	underWaterCounter int
	remove            bool
	eventHub          *tasks.EventHub
	stop              bool
	targX             float32
	targY             float32
	*sprite.Sprite
}

func (p *FoodParticle) ShouldRemove() bool {
	return p.remove
}

func (p *FoodParticle) float() {
	vy := 10.0
	dx := float32(-5.0)
	if p.Point.Tag == "left" {
		dx = float32(5)
	}
	if p.underWater {
		dx = -0.01
		p.underWaterCounter++
		vy -= 2 * float64(p.underWaterCounter)
		vy = max(vy, 0.07)
	}

	if p.counter%5 == 0 && p.underWater {
		vx := math.Sin(float64(p.counter)*0.5) * 0.3 * 1
		noise := rand.Float64()*0.1 - 0.05
		p.Point.X = p.Point.X + float32(vx+noise)
	}

	p.Point.Y += float32(vy)
	p.Point.X += dx

}

func (p *FoodParticle) Update() {

	if p.stop {
		return
	}
	p.counter++

	if !p.underWater && int(p.Point.Y) > p.waterLevel+10 {
		ev := SendData{Data: "particle entered water",
			DataFor: "soundFx"}
		p.eventHub.Publish(ev)

		initialNoise := math.Sin(rand.Float64()*10) * 30
		p.Point.X += float32(initialNoise)
		p.underWater = true
	}
	if int(p.Point.Y) < p.floorLevel {
		p.float()
	}
}

func (p *FoodParticle) Draw(screen *ebiten.Image) {
	drawOpts := &ebiten.DrawImageOptions{}
	drawOpts.GeoM.Translate(float64(p.Point.X), float64(p.Point.Y))

	if p.Point.Z >= 8 {
		screen.DrawImage(p.image.SubImage(image.Rect(4, 0, 8, 4)).(*ebiten.Image), drawOpts)
	} else if p.Point.Z >= 6 {
		screen.DrawImage(p.image.SubImage(image.Rect(0, 0, 4, 3)).(*ebiten.Image), drawOpts)
	} else {
		screen.DrawImage(p.image.SubImage(image.Rect(8, 0, 12, 3)).(*ebiten.Image), drawOpts)
	}
}

func NewParticle(point *util.Point, rect image.Rectangle, hub *tasks.EventHub) *Entity {
	println("calling new particle function", n)
	n++
	img, _ := util.LoadImageAssetAsEbitenImage("textures/foodParticle")

	p := &FoodParticle{
		Point:             point,
		counter:           0,
		underWater:        false,
		waterLevel:        rect.Min.Y,
		floorLevel:        rect.Max.Y,
		underWaterCounter: 0,
		eventHub:          hub,
		remove:            false,
		Sprite:            &sprite.Sprite{},
		image:             img,
	}

	//empty image so we can z sort, we should give entities position data for this
	z := rand.Intn(4)
	sp := &sprite.Sprite{Img: ebiten.NewImage(10, 10)}
	foodEntity := &Entity{ParticleData: p, Sprite: sp}
	foodEntity.Z = 6 + z
	foodEntity.ParticleData.floorLevel -= 12 - z
	foodEntity.ParticleData.Point.Z = foodEntity.Z
	foodEntity.EventHub = hub
	RegisterEntity(foodEntity)

	pointEvent := PointGenerated{PointId: foodEntity.Id, Source: "new particle function"}
	p.eventHub.Publish(pointEvent)

	foodEntity.subscribe()
	return foodEntity
}

func (p *Entity) subscribe() {
	p.EventHub.Subscribe(CreatureReachedPoint{}, func(e tasks.Event) {
		ev := e.(CreatureReachedPoint)
		if ev.PointID == p.Id {
			time.AfterFunc(300*time.Millisecond, func() { p.ParticleData.stop = true })
			time.AfterFunc(600*time.Millisecond, func() { RemoveEntity(p.Id) })
		}
	})
}
