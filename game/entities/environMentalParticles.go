package entities

import (
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"math"
	"math/rand"
)

var particleImg, _ = util.LoadImageAssetAsEbitenImage("textures/particleTexture")

//0= bright gold 1x1
//1= plus sign 3x3
//2 = minus sign 3x3
//3= bubble 2x2
//4 = bubble 4x4
//5 = bone white 2x2
//6 = 1x1 bubble
//7 = 1x1 bubbleLessOpaque
//8 = 2x2 gold

var Textures = [9]*ebiten.Image{
	particleImg.SubImage(image.Rect(8, 0, 9, 1)).(*ebiten.Image),
	particleImg.SubImage(image.Rect(18, 0, 27, 9)).(*ebiten.Image),
	particleImg.SubImage(image.Rect(27, 2, 33, 6)).(*ebiten.Image),
	particleImg.SubImage(image.Rect(1, 0, 3, 2)).(*ebiten.Image),
	particleImg.SubImage(image.Rect(0, 2, 4, 6)).(*ebiten.Image),
	particleImg.SubImage(image.Rect(16, 0, 18, 2)).(*ebiten.Image),
	particleImg.SubImage(image.Rect(0, 0, 1, 1)).(*ebiten.Image),
	particleImg.SubImage(image.Rect(3, 0, 4, 1)).(*ebiten.Image),
	particleImg.SubImage(image.Rect(8, 0, 10, 2)).(*ebiten.Image),
}

//0=  prop spawn rate
//1= ph boost spawn rate
//2 = ph min
//3= bubble 2x2
//4 = bubble 4x4
//5 = bone white 2x2
//6 = 1x1 bubble

var spawnRate = []float64{20, 5, 5}

type ParticleSystem struct {
	On           bool
	img          *ebiten.Image
	img2         *ebiten.Image
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
	EndAfter     float64
	flags        map[string]bool
	PConfig      *ParticleConfig
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
	typeFlag      uint32
	AlphaDecay    float64
}

func NewBubbleSystem(x float64, y float64, bounds image.Rectangle) *ParticleSystem {

	return &ParticleSystem{
		Particles:    make([]GenParticle, 0, 400),
		Sprite:       &sprite.Sprite{Img: ebiten.NewImage(bounds.Dx(), bounds.Dy()), X: float32(bounds.Min.X), Y: float32(bounds.Min.Y), UnFocusable: true},
		SpawnRate:    50.0, // particles per second
		MaxParticles: 1000,
		Bounds:       bounds,
		SpawnPointX:  x,
		SpawnPointY:  y,
		tag:          "bubble",
		img:          Textures[3],
		img2:         Textures[4],
		flags:        make(map[string]bool),
	}

}

func NewPhReaderParticleSystem(x, y float64, bounds image.Rectangle) *ParticleSystem {
	return &ParticleSystem{
		Particles:    make([]GenParticle, 0, 200),
		SpawnRate:    400,
		Sprite:       &sprite.Sprite{Img: ebiten.NewImage(bounds.Dx(), bounds.Dy()), X: float32(bounds.Min.X), Y: float32(bounds.Min.Y), UnFocusable: true},
		MaxParticles: 200,
		SpawnPointX:  x - float64(bounds.Min.X),
		SpawnPointY:  y - float64(bounds.Min.Y),
		Bounds:       bounds,
		tag:          "downwardBubble",
		On:           true,
		EndAfter:     2.0,
		flags:        make(map[string]bool),
		img:          Textures[6],
	}
}

func NewGenericParticleSystem(x float64, y float64, bounds image.Rectangle, textureTag uint32) *ParticleSystem {
	zb := bounds
	zb.Min.Y -= 5
	var SpawnRate uint32
	SpawnRate = textureTag
	if textureTag > uint32(len(spawnRate)) {
		SpawnRate = 1
	}

	return &ParticleSystem{
		Particles:    make([]GenParticle, 0, 200),
		SpawnRate:    spawnRate[SpawnRate], // particles per second
		Sprite:       &sprite.Sprite{Img: ebiten.NewImage(bounds.Dx(), bounds.Dy()), X: float32(bounds.Min.X), Y: float32(bounds.Min.Y), UnFocusable: true, DOptsUpdaterParams: make(map[string]float64)},
		MaxParticles: 300,
		SpawnPointX:  x - float64(bounds.Min.X),
		SpawnPointY:  y - float64(bounds.Min.Y),
		Bounds:       bounds,
		tag:          "plant",
		On:           true,
		EndAfter:     2.5,
		flags:        make(map[string]bool),
		img:          Textures[textureTag],
	}
}

func NewFertilizerParticleSystem(x float64, y float64, bounds image.Rectangle) *ParticleSystem {
	return &ParticleSystem{
		Particles:    make([]GenParticle, 0, 200),
		Sprite:       &sprite.Sprite{Img: ebiten.NewImage(bounds.Dx(), bounds.Dy()), X: float32(bounds.Min.X), Y: float32(bounds.Min.Y), UnFocusable: true},
		SpawnRate:    20.0, // particles per second
		MaxParticles: 200,
		SpawnPointX:  x - float64(bounds.Min.X),
		SpawnPointY:  y - float64(bounds.Min.Y),
		Bounds:       bounds,
		tag:          "fertilizer",
		On:           true,
		EndAfter:     3.0,
		img:          Textures[5],
		flags:        make(map[string]bool),
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
		p.Alpha -= p.AlphaDecay
		// Physics
		p.X += p.VelX * deltaTime
		p.Y += p.VelY * deltaTime * ((p.Y + 20) / float64(ps.Bounds.Dy())) //very hacky buoyancy
		p.Rotation += p.RotationSpeed * deltaTime

		// Life
		p.Life -= deltaTime
		p.Alpha = p.Life / p.MaxLife // Fade out over time

		ps.CheckBounds(p, ps.Bounds)

		// Remove dead particles
		if p.Life <= 0 {
			ps.Particles = append(ps.Particles[:i], ps.Particles[i+1:]...)
		}
	}
	if ps.EndAfter != 0 && ps.lifeTime >= ps.EndAfter {
		ps.On = false
	}
}

func (ps *ParticleSystem) CheckBounds(p *GenParticle, bounds image.Rectangle) {
	pt := image.Point{int(p.X), int(p.Y)}

	if !pt.In(bounds) {
		if pt.X > ps.Bounds.Dx() {
			pt.X = ps.Bounds.Dx()
			p.VelX = -(p.VelX / 2)
		}
		if pt.X < 0 {
			pt.X = 0
			p.VelX = -(p.VelX / 2)
		}
		if pt.Y > ps.Bounds.Dy() {
			pt.Y = ps.Bounds.Dy()
			p.VelY = -(p.VelY / 2)
		}
		if pt.Y < 0 {
			pt.Y = 0
			if ps.tag == "bubble" {
				p.VelY = 1
				p.VelX = rand.NormFloat64() * 5
			} else {
				p.VelY = -(p.VelY / 2)
			}

		}
	}
}

type ParticleConfig struct {
	XVariance         float64
	YVariance         float64
	XVelocityVariance float64
	YVelocityVariance float64
	BaseYVelocity     float64
	MaxLife           float64
	Scale             float64
	AlphaDecay        float64
	RotationSpeed     float64
}

func (ps *ParticleSystem) configurableParticleSpawn(config ParticleConfig) {
	if len(ps.Particles) >= ps.MaxParticles {
		return
	}

	if ps.PConfig.Scale == 0 {
		ps.PConfig.Scale = 1
	}

	particle := GenParticle{
		X:             ps.SpawnPointX + (rand.NormFloat64())*config.XVariance, // Random spread
		Y:             ps.SpawnPointY + (rand.NormFloat64())*config.YVariance,
		VelX:          (rand.NormFloat64()) * config.XVelocityVariance,                      // Random horizontal drift
		VelY:          (rand.NormFloat64())*config.YVelocityVariance + config.BaseYVelocity, // Float downward
		Life:          2.0 + rand.Float64()*config.MaxLife,                                  // 2-5 seconds
		MaxLife:       0,
		MaxHeight:     0,
		Size:          ps.PConfig.Scale,
		Rotation:      0,
		RotationSpeed: ps.PConfig.RotationSpeed, // Gentle spin
		Alpha:         1.0,
		AlphaDecay:    config.AlphaDecay,
	}
	particle.MaxLife = particle.Life

	ps.Particles = append(ps.Particles, particle)
}

func (ps *ParticleSystem) SpawnInteractionBubble() {
	if len(ps.Particles) >= ps.MaxParticles {
		return
	}
	particle := GenParticle{
		X:             ps.SpawnPointX + (rand.NormFloat64())*20,
		Y:             ps.SpawnPointY + (rand.Float64()-0.5)*5,
		VelX:          (rand.NormFloat64()) * 20, // Random horizontal drift
		VelY:          10 + rand.Float64()*20,    // Float downward
		Life:          2.0 + rand.Float64()*15,   // 2-5 seconds
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
	var typeFlag uint32 = 0
	if rand.Intn(10) < 1 {
		typeFlag = 1
	}
	if rand.Intn(4) < 3 {
		typeFlag = 2
	}

	particle := GenParticle{
		X:             ps.SpawnPointX + (rand.Float64()-0.5)*20, // Random spread
		Y:             ps.SpawnPointY - (rand.Float64()-0.5)*20,
		VelX:          (rand.Float64() - 0.5) * 2, // Random horizontal drift
		VelY:          -120 - rand.Float64()*30,   // Float upward
		Life:          2.0 + rand.Float64()*3.0,   // 2-5 seconds
		MaxLife:       0,
		Size:          0.5 + rand.Float64()*1.0, // Random size
		Rotation:      0,
		RotationSpeed: (rand.Float64() - 0.5) * 0.01,
		typeFlag:      typeFlag,
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
	if ps.PConfig != nil {
		ps.configurableParticleSpawn(*ps.PConfig)
		return
	}

	switch ps.tag {
	case "bubble":
		ps.SpawnBubble()
	case "plant":
		ps.PropSpawnParticle()
	case "downwardBubble":
		ps.SpawnInteractionBubble()
	case "fertilizer":
		ps.SpawnFertilizerParticle()
	}
}

func (ps *ParticleSystem) PropSpawnParticle() {
	if len(ps.Particles) >= ps.MaxParticles {
		return
	}

	particle := GenParticle{
		X:             ps.SpawnPointX + (rand.Float64()-0.5)*50, // Random spread
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
	}
	particle.MaxLife = particle.Life

	ps.Particles = append(ps.Particles, particle)
}

func (ps *ParticleSystem) SpawnFertilizerParticle() {
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
		Size:          1.0, // Random size
		Rotation:      0,
		RotationSpeed: (rand.Float64() - 0.5) * 2, // Gentle spin
		Alpha:         1.0,
	}
	particle.MaxLife = particle.Life

	ps.Particles = append(ps.Particles, particle)
}

func (ps *ParticleSystem) Draw() {
	ps.Sprite.Img.Clear()

	opts := &ebiten.DrawImageOptions{}
	for _, p := range ps.Particles {
		opts.GeoM.Reset()
		if p.Size != 0 {
			opts.GeoM.Scale(p.Size, p.Size)
		}
		opts.GeoM.Translate(p.X, p.Y)

		opts.GeoM.Rotate(p.Rotation)
		opts.ColorScale.SetA(float32(p.Alpha))
		if ps.img == nil {
			log.Fatal("error particle system image in draw call")
		}

		switch p.typeFlag {
		case 1:
			ps.Sprite.Img.DrawImage(Textures[4], opts)
		case 2:
			ps.Sprite.Img.DrawImage(Textures[6], opts)
		default:
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
