package entities

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/system"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"image"
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
	GoldFish  FishList = "goldFish"
	MollyFish FishList = "mollyFish"
	Guppy     FishList = "guppy"
	Kirbensis FishList = "kirbensis"
	AngelFish FishList = "angelFish"
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

type CreatureData struct {
	TargetPoint        *util.Point
	distanceTraveled   float64
	TargetZ            int
	TargetParticleId   uint32
	ParticlePointQueue map[uint32]*util.Point
	EventHub           *tasks.EventHub
	TankBoundaries     image.Rectangle
	Timers             map[FishState]*util.Timer
	State              FishState
	TickClicked        bool
	Environment        *system.Environment
	stressContributors []string
	MovementFlags      [10]uint32
	*FishStats
	Flip bool
}

func (e *Entity) FishUpdate(state *GameState) {
	c := e.CreatureData

	c.TankBoundaries = state.Zbounds[c.TargetZ]

	c.TickClicked = false
	switch c.State {

	case Swimming:
		e.swimmingUpdate(state.ActiveCollisions)
	case Resting:
		e.restingUpdate(state.ActiveCollisions)
	case Eating:

		e.eatingUpdate()
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

	e.updateAnimation()

	e.publishStats("statsMenu")

}

func (e *Entity) updateAnimation() {

	switch e.CreatureData.State {
	case Swimming:

		e.Sprite.CurrentAnimation = "swimming"

		if e.CreatureData.TargetZ < e.Z {
			e.Sprite.CurrentAnimation = "backwards"
		}
		if e.CreatureData.TargetZ > e.Z {
			e.Sprite.CurrentAnimation = "forward"
		}
		/*if e.CreatureData.TargetZ == e.Z && e.Z < 4 {
			if registry.Config.Zoom {
				e.Sprite.CurrentAnimation = "depth"
			}
		}*/
	case Eating:

		e.Sprite.CurrentAnimation = "eating"

	case Resting:
		// Use swimming animation for resting or create a resting animation
		e.Sprite.CurrentAnimation = "swimming"
	}

}

func (e *Entity) swimmingUpdate(collisions []FishCollision) {

	e.Move(collisions)
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

func (e *Entity) restingUpdate(collisions []FishCollision) {
	c := e.CreatureData
	c.speed = 0.4
	e.Move(collisions)

	if c.Timers[Resting].On == false {
		c.Timers[Resting].TurnOn()
	}

	tState := c.Timers[Resting].Update()
	if tState == util.Done {
		c.energy += 10
		if c.energy > 15 {
			c.State = Swimming
		}
	}
}

func (e *Entity) eatingUpdate() {

	if e.CreatureData == nil {
		log.Fatal("called a creature data func on a non creature some how")
	}
	c := e.CreatureData
	if !c.Timers[Eating].On {
		c.Timers[Eating].TurnOn()
	}

	tState := c.Timers[Eating].Update()
	if tState == util.Done {
		c.State = Swimming
		c.energy += 4
		DoneEating(e)
	}
}
