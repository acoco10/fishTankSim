package gameEntities

import (
	"github.com/acoco10/fishTankWebGame/game/debug"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math/rand"
)

type FishState uint8

const (
	Swimming FishState = iota
	Eating
	Resting
)

type FishList string

const (
	Fish      FishList = "fish"
	MollyFish FishList = "mollyFish"
)

type FishPersonality string

const (
	shy    FishPersonality = "shy"
	social FishPersonality = "social"
)

type Creature struct {
	pointQueue     []*Point
	eventHub       *EventHub
	tankBoundaries debug.Rect
	timers         map[FishState]*Timer
	state          FishState
	selected       bool
	*FishStats
	*AnimatedSprite
}

func NewFish(hub *EventHub, tankSize debug.Rect, saveData SavedFish) *Creature {

	timers := make(map[FishState]*Timer)
	randDuration := rand.Float64() * 50
	timers[Swimming] = NewTimer(randDuration)
	timers[Eating] = NewTimer(0.5)
	timers[Resting] = NewTimer(10)

	fs, err := GenFishStats(FishList(saveData.FishType))
	if err != nil {
		log.Fatal(err)
	}
	if fs == nil {
		println("Fish stats returning empty pointer")
	}

	fs.Size = saveData.Size
	fs.name = saveData.Name

	c := Creature{
		[]*Point{},
		hub,
		tankSize,
		timers,
		Swimming,
		false,
		fs,
		NewAnimatedSprite(),
	}

	c.AnimatedSprite = LoadFishSprite(c.fishType, c.Size)

	shaderParams := make(map[string]any)
	shaderParams["OutlineColor"] = [4]float64{255, 255, 0, 255}
	c.shaderParams = shaderParams
	firstPoint := c.RandomTarget()
	c.AddTargetPointToQueue(firstPoint)

	c.X = rand.Float32()*200 + tankSize.X1
	c.Y = rand.Float32()*100 + tankSize.Y1

	CreatureEventSubscriptions(&c)

	return &c
}

func (c *Creature) Update() {
	switch c.state {

	case Swimming:
		c.swimmingUpdate()
	case Resting:
		c.restingUpdate()
	case Eating:
		c.eatingUpdate()
	}

	c.LevelUp()

	c.AnimatedSprite.Update()

	if c.selected {
		c.publishStats("statsMenu")
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			c.selected = false
			c.shader = nil
		}
	}

	if c.SpriteHovered() {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			c.selected = true
			c.shader = LoadOutlineShader()
		}
	}

}

func (c *Creature) swimmingUpdate() {
	c.Move()
	tState := c.timers[Swimming].Update()

	if tState == Done {
		c.timers[Swimming].Duration = rand.Intn(40)
		if c.energy > 0 {
			c.state = Swimming
		} else {
			c.state = Resting
		}
	}
}

func (c *Creature) restingUpdate() {
	c.speed = 0.1
	c.Move()

	if c.timers[Resting].on == false {
		c.timers[Resting].TurnOn()
	}

	tState := c.timers[Resting].Update()
	if tState == Done {

		c.energy += 2
		c.state = Swimming
	}
}

func (c *Creature) eatingUpdate() {
	c.timers[Eating].TurnOn()
	tState := c.timers[Eating].Update()
	if tState == Done {
		c.state = Swimming
		c.energy++
	}
}

func (c *Creature) Draw(screen *ebiten.Image) {

	opts := c.TranSlateFishOpts()
	shaderOpts := c.TranSlateFishShaderOpts()

	c.AnimatedSprite.Draw(screen, opts, shaderOpts)

}
