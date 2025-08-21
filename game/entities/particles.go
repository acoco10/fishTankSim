package entities

import (
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"
)

var n int = 0

type Particle struct {
	*util.Point
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

func (p *Particle) ShouldRemove() bool {
	return p.remove
}

func (p *Particle) float() {
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

func (p *Particle) Update() {

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

func (p *Particle) Draw(screen *ebiten.Image) {
	clr := color.RGBA{255, 234, 0, 255}
	vector.DrawFilledCircle(screen, p.Point.X, p.Point.Y, 1, clr, false)
}

func NewParticle(point *util.Point, rect image.Rectangle, hub *tasks.EventHub) *Entity {
	println("calling new particle function", n)
	n++
	p := &Particle{
		Point:             point,
		counter:           0,
		underWater:        false,
		waterLevel:        rect.Min.Y,
		floorLevel:        rect.Max.Y - rand.Intn(10) - 10,
		underWaterCounter: 0,
		eventHub:          hub,
		remove:            false,
		Sprite:            &sprite.Sprite{},
	}

	//empty image so we can z sort, we should give entities position data for this
	sp := &sprite.Sprite{Z: 3, Img: ebiten.NewImage(10, 10)}
	foodEntity := &Entity{ParticleData: p, Sprite: sp}
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
