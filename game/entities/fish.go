package entities

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/system"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"image"
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
	Kirbensis FishList = "kirbensis"
)

func IsValidFishType(fishType string) bool {
	switch FishList(fishType) {
	case MollyFish, Fish, Guppy, Kirbensis:
		return true
	default:
		return false
	}
}

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

type CreatureData struct {
	TargetPoint        *util.Point
	TargetParticleId   uint32
	ParticlePointQueue map[uint32]*util.Point
	EventHub           *tasks.EventHub
	TankBoundaries     image.Rectangle
	Timers             map[FishState]*util.Timer
	State              FishState
	TickClicked        bool
	Environment        *system.Environment
	stressContributors []string
	*FishStats
	Flip bool
}

func (e *Entity) FishUpdate() {
	c := e.CreatureData

	c.TickClicked = false
	switch c.State {

	case Swimming:
		e.swimmingUpdate()
	case Resting:
		e.restingUpdate()
	case Eating:
		c.eatingUpdate()
	}

	dopts := e.TranSlateFishOpts()
	sopts := e.TranSlateFishShaderOpts()
	e.Sprite.UpdateOpts(sopts)
	e.Sprite.UpdateOpts(dopts)

	if registry.Config.Zoom {
		e.Sprite.Unfocusable = false
	} else {
		if e.Sprite.Focused {
			UnFocus(e.Id)
		}
		e.Sprite.Unfocusable = true
	}

	e.publishStats("statsMenu")

}

func (e *Entity) swimmingUpdate() {
	e.Move()
	c := e.CreatureData
	tState := c.Timers[Swimming].Update()

	if tState == util.Done {
		c.Timers[Swimming].Duration = rand.Intn(40)
		if c.energy > 0 {
			c.State = Swimming
		} else {
			c.State = Resting
		}
	}
}

func (e *Entity) restingUpdate() {
	c := e.CreatureData
	c.speed = 0.4
	e.Move()

	if c.Timers[Resting].On == false {
		c.Timers[Resting].TurnOn()
		//c.ChangeAnimationSpeed(3)
	}

	tState := c.Timers[Resting].Update()
	if tState == util.Done {
		c.energy += 10
	}
}

func (c *CreatureData) eatingUpdate() {
	c.Timers[Eating].TurnOn()
	tState := c.Timers[Eating].Update()
	if tState == util.Done {
		c.State = Swimming
		c.energy += 4
	}
}
