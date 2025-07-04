package entities

import (
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"math/rand"
)

var n int = 0

type Particle struct {
	*geometry.Point
	counter           int
	underWater        bool
	waterLevel        float32
	floorLevel        float32
	underWaterCounter int
	eventHub          *tasks.EventHub
}

func (p *Particle) SpriteHovered() bool {
	//TODO implement me
	panic("implement me")
}

func (p *Particle) float() {
	vy := 10.0
	dx := float32(-5.0)
	if p.underWater {
		dx = -0.01
		p.underWaterCounter++
		vy -= 2 * float64(p.underWaterCounter)
		vy = max(vy, 0.15)
	}

	if p.counter%5 == 0 && p.underWater {
		vx := math.Sin(float64(p.counter)*0.5) * 0.3 * 1
		noise := rand.Float64()*0.1 - 0.05
		p.X = p.X + float32(vx+noise)
	}

	p.Y += float32(vy)
	p.X += dx

}

func (p *Particle) Update() {
	p.counter++

	if !p.underWater && p.Y > p.waterLevel+10 {
		ev := SendData{Data: "particle entered water",
			DataFor: "soundFx"}
		p.eventHub.Publish(ev)

		initialNoise := math.Sin(rand.Float64()*10) * 30
		p.X += float32(initialNoise)
		p.underWater = true
	}
	if p.Y < p.floorLevel {
		p.float()
	}
}

func (p *Particle) Draw(screen *ebiten.Image) {
	clr := color.RGBA{255, 234, 0, 255}
	vector.DrawFilledCircle(screen, p.X, p.Y, 2, clr, false)
}

func NewParticle(point *geometry.Point, rect geometry.Rect, hub *tasks.EventHub) *Particle {
	println("calling new particle function", n)
	n++
	p := Particle{
		point,
		0,
		false,
		rect.Y1,
		rect.Y2,
		0,
		hub,
	}
	pointEvent := PointGenerated{Point: point, Source: "new particle function"}
	p.eventHub.Publish(pointEvent)
	return &p
}
