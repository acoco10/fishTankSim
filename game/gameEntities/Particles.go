package gameEntities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"math/rand"
)

type Particle struct {
	*Point
	counter    int
	underWater bool
}

func (p *Particle) float() {
	vy := 10.0
	if p.underWater {
		vy = 0.1
	}
	if p.counter%5 == 0 && p.underWater {
		vx := math.Sin(float64(p.counter)*0.5) * 0.3 * 1
		noise := rand.Float64()*0.1 - 0.05
		p.X = p.X + float32(vx+noise)
	}
	// uncomment if using randomness

	p.Y += float32(vy)

}

func (p *Particle) Update() {
	p.counter++
	p.float()
	if !p.underWater && p.Y > 100 {
		p.underWater = true
	}
}

func (p *Particle) Draw(screen *ebiten.Image) {
	clr := color.RGBA{255, 234, 0, 255}
	vector.DrawFilledCircle(screen, p.X, p.Y, 2, clr, false)
}

func NewParticle(point *Point) Particle {

	p := Particle{
		point,
		0,
		false,
	}

	return p
}
