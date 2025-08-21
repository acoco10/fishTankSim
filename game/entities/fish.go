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
		if e.AnimationMap["eating"] != nil {
			e.Sprite = e.AnimationMap["eating"]
			e.Sprite.X = e.AnimationMap["swimming"].X
			e.Sprite.Y = e.AnimationMap["swimming"].Y
		}
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

	//e.publishStats("statsMenu")

}

func (e *Entity) swimmingUpdate() {
	if e.Sprite != e.AnimationMap["swimming"] {
		e.Sprite = e.AnimationMap["swimming"]
		e.Sprite.X = e.AnimationMap["eating"].X
		e.Sprite.Y = e.AnimationMap["eating"].Y
	}

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
		DoneEating(e)
		c.State = Swimming
		c.energy += 4
	}
}
