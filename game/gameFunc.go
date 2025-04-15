package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"math"
	"math/rand"
)

func (g *Game) AddTargetPointForCreatures(point image.Point) {
	for _, creature := range g.Creatures {
		creature.AddTargetPointToQueue(point)
	}
}

type Particle struct {
	x, y       float32
	counter    int
	underWater bool
}

func (p *Particle) float() {
	vy := 10.0
	if p.underWater {
		vy = 0.1
	}
	if p.counter%5 == 0 && p.underWater {
		vx := math.Sin(float64(p.counter)*0.5) * 0.3 * 5
		noise := rand.Float64()*0.1 - 0.05
		p.x = p.x + float32(vx+noise)
	}
	// uncomment if using randomness

	p.y += float32(vy)

}

func (p *Particle) Update() {
	p.counter++
	p.float()
	if !p.underWater && p.y > 100 {
		p.underWater = true
	}
}

func (p *Particle) Draw(screen *ebiten.Image) {
	x := float32(p.x)
	y := float32(p.y)
	clr := color.RGBA{255, 234, 0, 255}
	vector.DrawFilledCircle(screen, x, y, 2, clr, false)
}

func NewParticle(x, y int) Particle {
	originX := rand.Float32()*100 + float32(x)
	originY := float32(y)

	p := Particle{
		originX,
		originY,
		0,
		false,
	}

	return p
}
