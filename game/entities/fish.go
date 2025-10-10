package entities

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/system"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

type State uint8

type FishState uint8

const (
	Swimming FishState = iota
	Eating
	Resting
	Captured
	Falling
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
	ArrivedAtPoint     uint32
	inBetweenPoint     *util.Point
	distanceTraveled   float64
	TargetZ            int
	TargetParticleId   uint32
	ParticlePointQueue []uint32
	EventHub           *tasks.EventHub
	TankBoundaries     image.Rectangle
	Timers             map[FishState]*util.Timer
	State              FishState
	StateWas           FishState
	Environment        *system.Environment
	stressContributors []string
	MovementFlags      [10]uint32
	Flags              FishFlags
	*FishStats
	Flip bool
}

func InitFishStateMachine() *StateMachine {
	sm := &StateMachine{}

	swimmingUpdater := StateHandler{Updater: swimmingUpdate, TransitionTo: 2, TransitionToFunc: TransitionToSwimming}
	restingUpdater := StateHandler{Updater: restingUpdate, TransitionTo: 3, TransitionOutFunc: nil}
	eatingUpdater := StateHandler{Updater: eatingUpdate, TransitionTo: 1, TransitionToFunc: TransitionToEating}
	fallingUpdater := StateHandler{Updater: FallingUpdate, TransitionTo: 4, TransitionOutFunc: nil}
	fishCaptureUpdater := StateHandler{Updater: fishCapturedUpdater, TransitionTo: 3, TransitionToFunc: transitionToCaptured, TransitionOutFunc: transitionFromCaptured}

	sm.States = map[int]*StateHandler{
		1: &swimmingUpdater,
		2: &restingUpdater,
		3: &eatingUpdater,
		4: &fallingUpdater,
		5: &fishCaptureUpdater,
	}

	sm.stateController = FishStateController
	sm.CurrentState = 1
	sm.EveryUpdate = append(sm.EveryUpdate, FishUpdate)

	return sm
}

func TransitionToEating(ent *Entity) {
	point, exists := GetEntity(ent.CreatureData.ArrivedAtPoint)
	if !exists {
		log.Fatal("WARNING fish is is trying to eat unregistered point")
	}

	ent.CreatureData.State = Eating
	ent.CreatureData.Timers[Eating].TurnOn()

	ent.CreatureData.Hunger++
	ent.Add1expGraphic()
	ent.CreatureData.progress += 1

	ev := CreatureReachedPoint{
		PointTypeReached: point.ParticleData.PType,
		PointID:          ent.CreatureData.TargetParticleId,
		CreatureID:       ent.Id,
	}

	ent.EventHub.Publish(ev)

	ent.CreatureData.TargetParticleId = 0

	point.StopParticle()
}

func TransitionToSwimming(ent *Entity) {

	ent.CreatureData.State = Swimming
}

func FishStateController(ent *Entity) int {
	//core assumption: transition is actively called within a state, not arbitrarily enacted at set intervals
	//has boundaries == captured
	if ent.CreatureData.TankBoundaries.Dy() < 150 {
		return 5
	}

	if ent.CreatureData.State == 5 {
		return 4
	}

	if ent.CreatureData.ArrivedAtPoint != 0 {
		point, exists := GetEntity(ent.CreatureData.ArrivedAtPoint)
		if !exists {
			log.Println("WARNING fish arrived at unregistered point")
			return 1
		}
		if point.ParticleData.PType == util.Food {
			return 3
		}
	}

	if ent.CreatureData.TargetPoint == nil {
		ent.ProcessTargetPointQueue()
	}

	//resting timer triggerd = resting
	if ent.CreatureData.Timers[Resting].On {
		return 2
	}

	//swimming
	return 1
}

func FishUpdate(ent *Entity, gs GameState) {
	c := ent.CreatureData

	if c.State != Captured {
		//update always for z changes if not captured
		c.TankBoundaries = gs.Zbounds[ent.Z]
	}

	dopts := ent.TranSlateFishOpts(ebiten.DrawImageOptions{})
	ent.Sprite.UpdateOpts(&dopts)

	if ent.Sprite.Shader != nil {
		sopts := ent.TranSlateFishShaderOpts()
		ent.Sprite.UpdateOpts(sopts)

	}

	if registry.Config.Zoom {
		ent.Sprite.UnFocusable = false
	} else {
		if ent.Sprite.Focused {
			UnFocus(ent.Id)
		}
		ent.Sprite.UnFocusable = true
	}

	ent.updateAnimation()
	//e.publishStats("statsMenu")

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
		e.Sprite.GetAnimation().SpeedInTPS = 12

	case Resting:
		// Use swimming animation for resting or create a resting animation
		e.Sprite.CurrentAnimation = "swimming"
	}

}

func swimmingUpdate(ent *Entity, gs GameState) {
	if len(ent.CreatureData.ParticlePointQueue) > 0 && ent.CreatureData.ArrivedAtPoint != 0 {
		ent.Transition()
	}
	ent.Move(gs.ActiveCollisions)

}

func restingUpdate(ent *Entity, gs GameState) {
	c := ent.CreatureData
	c.speed = 0.4
	ent.Move(gs.ActiveCollisions)

	tState := c.Timers[Resting].Update()
	if tState == util.Done {
		c.Timers[Resting].TurnOff()
		c.energy += 10
		ent.Transition()
	}
}

func eatingUpdate(ent *Entity, gs GameState) {

	c := ent.CreatureData

	tState := c.Timers[Eating].Update()
	if tState == util.Done {
		c.energy += 4
		c.Timers[Eating].TurnOff()
		RemoveEntity(ent.CreatureData.ArrivedAtPoint)
		ent.CreatureData.ArrivedAtPoint = 0
		ent.StateMachine.Transition(ent)
	}
}

func fishCapturedUpdater(ent *Entity, gs GameState) {
	ent.Z = 5
	ent.CreatureData.TargetZ = 5
	ent.Move(gs.ActiveCollisions)
}

func FallingUpdate(ent *Entity, gs GameState) {
	if ent.Sprite.Y < float32(ent.CreatureData.TankBoundaries.Min.Y) {
		ent.Sprite.Y += 9.8
		//fall until we are back in the tank boundaries which were reset when captured condition was removed
	} else {
		ent.Transition()
	}
}

func transitionToCaptured(ent *Entity) {
	ent.CreatureData.State = Captured
	ent.AddDeInitHandler(LoadFollowEffectAsEnt("angry", ent.Id, ent.EventHub, EntityParameters{}))
}

func transitionFromCaptured(ent *Entity) {
	ent.DeInitEffects()
}
