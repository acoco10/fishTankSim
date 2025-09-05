package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"math"
	"math/rand"
)

type ParticleSystem struct {
	On           bool
	img          *ebiten.Image
	Particles    []GenParticle
	Texture      *ebiten.Image
	SpawnRate    float64
	SpawnTimer   float64
	MaxParticles int
	Bounds       image.Rectangle
	SpawnPointX  float64
	SpawnPointY  float64
	tag          string
	lifeTime     float64
	Sprite       *sprite.Sprite
	endAfter     float64
}

type GenParticle struct {
	X, Y          float64
	img           *ebiten.Image
	VelX, VelY    float64
	Life          float64
	MaxLife       float64
	Size          float64
	Rotation      float64
	RotationSpeed float64
	MaxHeight     float64
	Alpha         float64
}

func NewBubbleSystem(x float64, y float64, bounds image.Rectangle) *ParticleSystem {
	img, err := util.LoadImageAssetAsEbitenImage("textures/particleTexture")
	if err != nil {
		log.Fatal(err)
	}
	return &ParticleSystem{
		Particles:    make([]GenParticle, 0, 200),
		Sprite:       &sprite.Sprite{Img: ebiten.NewImage(700, 1200), X: 0, Y: 0},
		SpawnRate:    10.0, // particles per second
		MaxParticles: 200,
		Bounds:       bounds,
		SpawnPointX:  300,
		SpawnPointY:  300,
		tag:          "bubble",
		img:          img.SubImage(image.Rect(1, 0, 3, 2)).(*ebiten.Image),
	}

}

func NewPhReaderParticleSystem(x, y float64, bounds image.Rectangle) *ParticleSystem {
	img, err := util.LoadImageAssetAsEbitenImage("textures/particleTexture")
	if err != nil {
		log.Fatal(err)
	}
	return &ParticleSystem{
		Particles:    make([]GenParticle, 0, 200),
		SpawnRate:    10.0, // particles per second
		MaxParticles: 200,
		SpawnPointX:  x,
		SpawnPointY:  y,
		Bounds:       bounds,
		tag:          "downwardBubble",
		On:           true,
		endAfter:     2.0,
		img:          img.SubImage(image.Rect(2, 0, 3, 1)).(*ebiten.Image),
	}
}

func NewPlantParticleSystem(x float64, y float64, bounds image.Rectangle) *ParticleSystem {
	img, err := util.LoadImageAssetAsEbitenImage("textures/particleTexture")
	if err != nil {
		log.Fatal(err)
	}
	return &ParticleSystem{
		Particles:    make([]GenParticle, 0, 200),
		SpawnRate:    20.0, // particles per second
		MaxParticles: 200,
		SpawnPointX:  x,
		SpawnPointY:  y,
		Bounds:       bounds,
		tag:          "plant",
		On:           true,
		endAfter:     2.5,
		img:          img.SubImage(image.Rect(2, 0, 3, 1)).(*ebiten.Image),
	}
}

func NewFertilizerParticleSystem(x float64, y float64, bounds image.Rectangle) *ParticleSystem {
	img, err := util.LoadImageAssetAsEbitenImage("textures/particleTexture")
	if err != nil {
		log.Fatal(err)
	}
	return &ParticleSystem{
		Particles:    make([]GenParticle, 0, 200),
		SpawnRate:    20.0, // particles per second
		MaxParticles: 200,
		SpawnPointX:  x,
		SpawnPointY:  y,
		Bounds:       bounds,
		tag:          "fertilizer",
		On:           true,
		endAfter:     2.0,
		img:          img.SubImage(image.Rect(2, 0, 3, 1)).(*ebiten.Image),
	}
}

func NewDebrisSystem(texture *ebiten.Image) *ParticleSystem {
	return &ParticleSystem{
		Particles:    make([]GenParticle, 0, 100),
		Texture:      texture,
		SpawnRate:    1.0,
		MaxParticles: 100,
	}
}

func (ps *ParticleSystem) Update() {
	ps.lifeTime += 0.016
	deltaTime := 0.016
	ps.SpawnTimer += deltaTime

	if ps.tag == "bubble" {
		randomInterval := 0.3 + rand.Float64()*0.4
		if ps.SpawnTimer >= randomInterval {
			ps.SpawnRate = 50
		} else {
			ps.SpawnRate = 20
		}
	}
	// Spawn new particles
	if ps.SpawnTimer >= 1.0/ps.SpawnRate && len(ps.Particles) < ps.MaxParticles && ps.On {
		ps.SpawnParticle()
		ps.SpawnTimer = 0
	}

	// Update existing particles
	for i := len(ps.Particles) - 1; i >= 0; i-- {
		p := &ps.Particles[i]

		// Physics
		p.X += p.VelX * deltaTime
		p.Y += p.VelY * deltaTime
		p.Rotation += p.RotationSpeed * deltaTime

		// Life
		p.Life -= deltaTime
		p.Alpha = p.Life / p.MaxLife // Fade out over time

		pt := image.Point{int(p.X), int(p.Y)}
		if !pt.In(ps.Bounds) {
			if pt.X > ps.Bounds.Max.X {
				pt.X = ps.Bounds.Max.X
				p.VelX = -(p.VelX / 2)
			}
			if pt.X < ps.Bounds.Min.X {
				pt.X = ps.Bounds.Min.X
				p.VelX = -(p.VelX / 2)
			}
			if pt.Y > ps.Bounds.Max.Y+5 {
				pt.Y = ps.Bounds.Max.Y
				p.VelY = -(p.VelY / 2)
			}
			if pt.Y < ps.Bounds.Min.Y {
				pt.Y = ps.Bounds.Min.Y - 5
				if ps.tag == "bubble" {
					p.VelY = 1
				} else {
					p.VelY = -(p.VelY / 2)
				}

			}
		}
		// Remove dead particles
		if p.Life <= 0 {
			ps.Particles = append(ps.Particles[:i], ps.Particles[i+1:]...)
		}
	}
	if ps.endAfter != 0 && ps.lifeTime >= ps.endAfter {
		ps.On = false
	}
}

func (ps *ParticleSystem) SpawnInteractionBubble() {
	if len(ps.Particles) >= ps.MaxParticles {
		return
	}
	particle := GenParticle{
		X:             ps.SpawnPointX + (rand.NormFloat64())*20, // Random spread
		Y:             ps.SpawnPointY + (rand.Float64()-0.5)*20,
		VelX:          (rand.NormFloat64()) * 20, // Random horizontal drift
		VelY:          10 + rand.Float64()*10,    // Float downward
		Life:          2.0 + rand.Float64()*3.0,  // 2-5 seconds
		MaxLife:       0,
		MaxHeight:     0,
		Size:          1.0,
		Rotation:      0,
		RotationSpeed: (rand.Float64() - 0.5) * 2, // Gentle spin
		Alpha:         1.0,
	}
	particle.MaxLife = particle.Life

	ps.Particles = append(ps.Particles, particle)
}

func (ps *ParticleSystem) SpawnBubble() {
	if len(ps.Particles) >= ps.MaxParticles {
		return
	}
	particle := GenParticle{
		X:             ps.SpawnPointX, // Random spread
		Y:             ps.SpawnPointY,
		VelX:          (rand.Float64() - 0.5) * 2, // Random horizontal drift
		VelY:          -80 - rand.Float64()*30,    // Float upward
		Life:          2.0 + rand.Float64()*3.0,   // 2-5 seconds
		MaxLife:       0,
		Size:          0.5 + rand.Float64()*1.0, // Random size
		Rotation:      0,
		RotationSpeed: (rand.Float64() - 0.5) * 2, // Gentle spin

	}
	particle.MaxLife = particle.Life

	ps.Particles = append(ps.Particles, particle)
}

func (ps *ParticleSystem) SpawnDebris(x, y float64) {
	if len(ps.Particles) >= ps.MaxParticles {
		return
	}

	particle := GenParticle{
		X:             x,
		Y:             y,
		VelX:          (rand.Float64() - 0.5) * 100, // Random scatter
		VelY:          -50 - rand.Float64()*50,      // Float upward
		Life:          5.0 + rand.Float64()*5.0,     // 5-10 seconds
		MaxLife:       0,
		Size:          0.3 + rand.Float64()*0.7, // Smaller debris
		Rotation:      rand.Float64() * math.Pi * 2,
		RotationSpeed: (rand.Float64() - 0.5) * 5, // Tumble
		Alpha:         1.0,
	}
	particle.MaxLife = particle.Life

	ps.Particles = append(ps.Particles, particle)
}

func (ps *ParticleSystem) SpawnParticle() {
	switch ps.tag {
	case "bubble":
		ps.SpawnBubble()
	case "plant":
		ps.PropSpawnParticle()
	case "downwardBubble":
		ps.SpawnInteractionBubble()
	case "fertilizer":
		ps.FertilizeSpawnParticle()
	}
}

func (ps *ParticleSystem) PropSpawnParticle() {
	if len(ps.Particles) >= ps.MaxParticles {
		return
	}

	particle := GenParticle{
		X:             ps.SpawnPointX + (rand.NormFloat64())*10, // Random spread
		Y:             ps.SpawnPointY - (rand.Float64()-0.5)*20,
		VelX:          (rand.Float64() - 0.5) * 2, // Random horizontal drift
		VelY:          -20 - rand.Float64()*30,    // Float upward
		Life:          1.0 + rand.Float64()*0.5,   // 2-5 seconds
		MaxLife:       2.0,
		MaxHeight:     50,
		Size:          1.0, // Random size
		Rotation:      0,
		RotationSpeed: (rand.Float64() - 0.5) * 2, // Gentle spin
		Alpha:         1.0,
		img:           ps.Texture.SubImage(image.Rect(2, 0, 5, 2)).(*ebiten.Image),
	}
	particle.MaxLife = particle.Life

	ps.Particles = append(ps.Particles, particle)
}

func (ps *ParticleSystem) FertilizeSpawnParticle() {
	if len(ps.Particles) >= ps.MaxParticles {
		return
	}

	particle := GenParticle{
		X:             ps.SpawnPointX + (rand.NormFloat64())*2, // Random spread
		Y:             ps.SpawnPointY,
		VelX:          (rand.Float64() - 0.5) * 2, // Random horizontal drift
		VelY:          30 + rand.Float64()*30,     // Float upward
		Life:          1.0 + rand.Float64()*0.5,   // 2-5 seconds
		MaxLife:       2.0,
		MaxHeight:     50,
		Size:          1.0, // Random size
		Rotation:      0,
		RotationSpeed: (rand.Float64() - 0.5) * 2, // Gentle spin
		Alpha:         1.0,
		img:           ps.Texture.SubImage(image.Rect(2, 0, 3, 1)).(*ebiten.Image),
	}
	particle.MaxLife = particle.Life

	ps.Particles = append(ps.Particles, particle)
}

func (ps *ParticleSystem) Draw() {
	ps.Sprite.Img.Clear()
	if ps.On {
		opts := &ebiten.DrawImageOptions{}
		for _, p := range ps.Particles {
			opts.GeoM.Reset()

			//opts.GeoM.Rotate(p.Rotation)

			// Position
			opts.GeoM.Translate(p.X, p.Y)

			fmt.Printf("Drawing at: (%f, %f)",
				p.X, p.Y,
			)
			ps.Sprite.Img.DrawImage(ps.img, opts)
		}
	}
}

func (ps *ParticleSystem) BubbleSubscriptions(hub *tasks.EventHub) {
	hub.Subscribe(TurnOnBubbles{}, func(e tasks.Event) {
		ps.On = true
	})
	hub.Subscribe(TurnOffBubbles{}, func(e tasks.Event) {
		ps.On = false
	})
}

type Physics struct {
	weight      float64
	velocity    Vec2
	floatFactor float64 //in water how much does it float 1.0 not affect by gravity at all
}

type Vec2 struct {
	x, y float64
}

func UnderWaterPhysicsUpDater(particle GenParticle) {

}
