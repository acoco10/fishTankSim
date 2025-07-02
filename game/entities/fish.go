package entities

import (
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	Guppy     FishList = "guppy"
)

type Direction uint8

const (
	Left Direction = iota
	Right
)

type FishPersonality string

const (
	shy    FishPersonality = "shy"
	social FishPersonality = "social"
)

type Creature struct {
	PointQueue     []*geometry.Point
	EventHub       *tasks.EventHub
	TankBoundaries geometry.Rect
	Timers         map[FishState]*Timer
	State          FishState
	Selected       bool
	TickClicked    bool
	*FishStats
	*sprite.AnimatedSprite
	Flip bool
}

func (c *Creature) Update() {
	c.TickClicked = false
	switch c.State {

	case Swimming:
		c.swimmingUpdate()
	case Resting:
		c.restingUpdate()
	case Eating:
		c.eatingUpdate()
	}

	c.AnimatedSprite.Update()

	dopts := c.TranSlateFishOpts()
	sopts := c.TranSlateFishShaderOpts()
	c.UpdateOpts(sopts)
	c.UpdateOpts(dopts)

	if c.Selected {
		c.publishStats("statsMenu")
		if ebiten.IsKeyPressed(ebiten.KeyEscape) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !c.SpriteHovered() {
			c.Selected = false
			c.Shader = nil
			ev := SendData{Data: "fish deselect", DataFor: "statsMenu"}
			c.EventHub.Publish(ev)
		}
	}

}

func (c *Creature) swimmingUpdate() {
	c.Move()
	tState := c.Timers[Swimming].Update()

	if tState == Done {
		c.Timers[Swimming].Duration = rand.Intn(40)
		if c.energy > 0 {
			c.State = Swimming
		} else {
			c.State = Resting
		}
	}
}

func (c *Creature) restingUpdate() {
	c.speed = 0.4
	c.Move()

	if c.Timers[Resting].on == false {
		c.Timers[Resting].TurnOn()
	}

	tState := c.Timers[Resting].Update()
	if tState == Done {
		c.energy += 10
		c.State = Swimming
	}
}

func (c *Creature) eatingUpdate() {
	c.Timers[Eating].TurnOn()
	tState := c.Timers[Eating].Update()
	if tState == Done {
		c.State = Swimming
		c.energy += 4
	}
}

func (c *Creature) Draw(screen *ebiten.Image) {
	c.AnimatedSprite.Draw(screen)
}
