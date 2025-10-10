package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"math"
	"math/rand"
	"slices"
)

type HealthState uint8

const backFishLayer = 2
const frontFishLayer = 11
const dontTurnAroundDist = 50
const minSpeed = 0.6

const (
	Healthy HealthState = iota
	Stressed
	ReallyStressed
	Sick
	Dead
)

func CreatureEventSubscriptions(c *Entity) {

	c.EventHub.Subscribe(PointGenerated{}, func(e tasks.Event) {
		//only add it to queue, state controller decides what to do
		ev := e.(PointGenerated)
		point, exists := GetEntity(ev.PointId)
		if !exists {
			log.Fatal("A recently created point should exist....")
		}
		if c.CreatureData.Hunger == c.CreatureData.MaxHunger && point.ParticleData.PType == util.Food {
			return
		}

		c.EnqueueParticlePoint(ev.PointId)
	})

	c.EventHub.Subscribe(CreatureReachedPoint{}, func(e tasks.Event) {
		ev := e.(CreatureReachedPoint)
		pt, exists := GetEntity(ev.PointID)
		if !exists {
			log.Println("WARNING: fish point subscription  received creature received point with entity id that is not initiated")
			return
		}

		if ev.CreatureID == c.Id || pt.ParticleData.PType == util.Structure {
			//if its your own point(we will process in the queue properly) or a structure, dont worry about it
			return
		}

		if pt.ParticleData == nil {
			log.Println("fish point subscription received point has no particle data")
			return
		}

		for _, p := range c.CreatureData.ParticlePointQueue {
			if p == ev.PointID && pt.ParticleData.PType == util.Food {
				//remove food that has been eaten by another fish
				c.RemoveParticlePoint(ev.PointID)
				c.SortParticleQueue()
			}
		}

		c.StateMachine.Transition(c) //send to to state machine to find next point
	})

	c.CreatureData.EventHub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		c.CreatureData.Hunger = 0
		c.CreatureData.energy = c.CreatureData.maxEnergy
	})

	c.EventHub.Subscribe(events.NewProp{}, func(e tasks.Event) {
		/*ev := e.(events.NewProp)

		c.effectDeInitHandler = LoadFollowEffectAsEnt("exclamation", c.Id, c.EventHub, EntityParameters{})*/

		/*prop, exists := GetEntity(ev.PropId)
		if !exists {
			log.Fatal("some weird shit happened when fish was attracted to a new structure")
		}*/

		/*TargetPoint := &util.Point{X: float32(prop.Sprite.TranslatedMidX()), Y: float32(prop.Sprite.TranslatedMidY()), PType: util.Structure}
		println("making creature target point new prop")
		c.CreatureData.TargetParticleId = ev.PropId
		c.CreatureData.ParticlePointQueue[ev.PropId] = TargetPoint
		c.SetTargetPoint(TargetPoint)
		c.CreatureData.TargetZ = min(prop.Z+1, 12)*/
	})

	c.CreatureData.EventHub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		// add normal-map bac.CreatureData. when it's a new day since it was turned off for night-scene
		//c.CreatureData.Shader = registry.ShaderMap["NormalMap"]
		println("new day  received  for fish:", c)
		c.CreatureData.CalcDailyFishHealthState()
		c.CheckAndLevelUp()
		c.CreatureData.age += 1
	})
}

func (ent *Entity) RemoveParticlePoint(particleID uint32) {
	for i, id := range ent.CreatureData.ParticlePointQueue {
		if id == particleID {
			// Remove element by slicing around it
			ent.CreatureData.ParticlePointQueue = append(
				ent.CreatureData.ParticlePointQueue[:i],
				ent.CreatureData.ParticlePointQueue[i+1:]...,
			)
		}
	}
}

func (ent *Entity) EnqueueParticlePoint(entID uint32) {
	// Check if point already exists in queue
	for _, existingPoint := range ent.CreatureData.ParticlePointQueue {
		if existingPoint == entID {
			return // Don't add duplicate
		}
	}

	// Add point to end of queue
	ent.CreatureData.ParticlePointQueue = append(ent.CreatureData.ParticlePointQueue, entID)
}

func (ent *Entity) DequeueParticlePoint() *util.Point {
	if len(ent.CreatureData.ParticlePointQueue) == 0 {
		return nil
	}

	// Get first point
	id := ent.CreatureData.ParticlePointQueue[0]
	point, exists := GetEntity(id)
	if !exists {
		return nil
	}
	ent.CreatureData.TargetParticleId = id

	// Remove first point from queue
	ent.CreatureData.ParticlePointQueue = ent.CreatureData.ParticlePointQueue[1:]

	return point.ParticleData.Point
}

func (ent *Entity) ProcessTargetPointQueue() {
	//deque point or set to random point if no targets queued
	ent.SortParticleQueue()
	newTarg := ent.DequeueParticlePoint()
	if newTarg != nil {
		ent.SetTargetPoint(newTarg)
	} else {
		ent.SetTargetPoint(ent.RandomTarget())
	}
}

func (c *CreatureData) CalcDailyFishHealthState() {
	c.stressContributors = []string{}
	if math.Abs(float64(c.Environment.Temperature-c.IdealTemperature)) > 10 {
		c.stressContributors = append(c.stressContributors, "temperature")
		c.Happiness -= 1
		c.Stress += 2
	}
	//compare environment ph to ideal
	if math.Abs(float64(c.Environment.NaturalPHLevel-c.IdealPH)) > 0.5 {
		c.stressContributors = append(c.stressContributors, "ph")
		c.Stress += 2
	}

	if len(c.stressContributors) == 0 {
		c.Stress -= 1
	}

	if c.HealthState != Sick {
		if c.Stress > 5 {
			c.HealthState = ReallyStressed
			c.DaysStressed += 2
		} else if c.Stress > 3 {
			c.HealthState = Stressed
			c.Hunger = 3 //less growth possible if stressed
			c.DaysStressed++
		}

		//check if fish gets sick if chronically stressed
		if c.DaysStressed > 1 && c.HealthState != Sick {
			sickChance := rand.Intn(10) + c.DaysStressed
			if sickChance > 6 {
				c.HealthState = Sick
			}
		}
	}

	if c.HealthState == Sick {
		c.DaysSick++
	}
	//roll for death modified by total stress of environment
	if c.DaysSick > 3 {
		deathChance := rand.Float32()*10 + c.Stress
		if deathChance > 7 {
			c.HealthState = Dead
		}
	}

}

func (e *Entity) calcSpeed() {

	if e.CreatureData == nil {
		return
	}
	c := e.CreatureData
	ratio := minSpeed + (float32(e.Z)-1.0)/11.0*(1.0-minSpeed)
	if len(c.ParticlePointQueue) > 1 {
		// has destination
		c.speed = c.maxSpeed * ratio
	} else {
		c.speed = (float32(rand.NormFloat64())*c.stdDevSpeed + c.avgSpeed) * ratio
	}

	fmt.Printf("random speed generated = %f\n", c.speed)

}

func DistanceFunc(x, x2, y, y2 float32) float64 {
	xDis := float64(x - x2)
	yDis := float64(y - y2)
	dis := math.Sqrt(math.Pow(xDis, 2) + math.Pow(yDis, 2))
	return dis
}

func (c *Entity) Type() util.InterestPoint {
	return util.OtherCreature
}

func ClosestParticle(x, y float32, points map[uint32]*util.Point) uint32 {

	closest := 1000.0
	var pointToReturnID uint32

	for key, val := range points {
		dist := DistanceFunc(x, val.X, y, val.Y)
		if dist < closest {
			pointToReturnID = key
		}
	}

	return pointToReturnID
}

func (e *Entity) Move(collisions []FishCollision) {
	if e.CreatureData == nil {
		return
	}
	e.CalculateVelocity(collisions)

	/*if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		e.Sprite.Dx -= 0.1 // Move left (negative X direction)
	} else if e.Sprite.Dx < 0 {
		e.Sprite.Dx += 0.05 // Slow down leftward movement when key released
		if e.Sprite.Dx > -0.01 {
			e.Sprite.Dx = 0 // Stop when close to zero
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		e.Sprite.Dx += 0.1 // Move right (positive X direction)
	} else if e.Sprite.Dx > 0 {
		e.Sprite.Dx -= 0.05 // Slow down rightward movement when key released
		if e.Sprite.Dx < 0.01 {
			e.Sprite.Dx = 0 // Stop when close to zero
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		e.Sprite.Dy += 0.1 // Move down (positive Y direction)
	} else if e.Sprite.Dy > 0 {
		e.Sprite.Dy -= 0.05 // Slow down downward movement when key released
		if e.Sprite.Dy < 0.01 {
			e.Sprite.Dy = 0 // Stop when close to zero
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		e.Sprite.Dy -= 0.1 // Move up (negative Y direction)
	} else if e.Sprite.Dy < 0 {
		e.Sprite.Dy += 0.05 // Slow down upward movement when key released
		if e.Sprite.Dy > -0.01 {
			e.Sprite.Dy = 0 // Stop when close to zero
		}
	}
	*/
	// Apply movement with max speed limits
	/*maxSpeed := float32(2.0)
	e.Sprite.Dx = max(-maxSpeed, min(maxSpeed, e.Sprite.Dx))
	e.Sprite.Dy = max(-maxSpeed, min(maxSpeed, e.Sprite.Dy))*/

	e.Sprite.X += e.Sprite.Dx
	e.Sprite.Y += e.Sprite.Dy

	e.CreatureData.distanceTraveled += 1

	e.EnforceBoundaries()
	e.CheckPointReached()

}

func (c *Entity) EnforceBoundaries() {
	s := c.Sprite
	rect := s.GetSpriteRect() // This already accounts for flip correctly

	// Use the actual rectangle bounds
	if rect.Min.X < c.CreatureData.TankBoundaries.Min.X {
		// Move sprite so left edge aligns with boundary
		c.CreatureData.SetHitBoundary()
		offset := c.CreatureData.TankBoundaries.Min.X - rect.Min.X
		s.X += float32(offset)
		c.SetTargetPoint(c.RandomTarget())
	}
	if rect.Max.X > c.CreatureData.TankBoundaries.Max.X {
		c.CreatureData.SetHitBoundary()
		// Move sprite so right edge aligns with boundary
		offset := rect.Max.X - c.CreatureData.TankBoundaries.Max.X
		s.X -= float32(offset)
		c.SetTargetPoint(c.RandomTarget())

	}
	if rect.Min.Y < c.CreatureData.TankBoundaries.Min.Y {
		offset := c.CreatureData.TankBoundaries.Min.Y - rect.Min.Y
		s.Y += float32(offset)
		c.SetTargetPoint(c.RandomTarget())

	}
	if rect.Max.Y > c.CreatureData.TankBoundaries.Max.Y {
		offset := rect.Max.Y - c.CreatureData.TankBoundaries.Max.Y
		s.Y -= float32(offset)
		c.SetTargetPoint(c.RandomTarget())

	}
}

func SetTurningPoint(ent *Entity) *util.Point {
	c := ent.CreatureData
	newPoint := util.Point{}
	//check if we cant move any more back in the tank

	zChance := rand.Intn(2)
	//randomly choose to swim forward or backwards during turn within the tank z layers
	if zChance == 1 {
		if ent.Z > backFishLayer {
			//dont want layer 0/1 to be accesible by fish since back of tank is there
			c.TargetZ = ent.Z - 1
			//no need to check if we can move forward if Z < 2
		} else {
			c.TargetZ = ent.Z + 1

		}
	} else {
		if ent.Z < frontFishLayer {
			//dont want layer 0/1 to be accesible by fish since back of tank is there
			c.TargetZ = ent.Z + 1
			//no need to check if we can move Backward if z > 2
		} else {
			c.TargetZ = ent.Z - 1
		}
	}

	if c.IsGoingLeft() {
		newPoint = util.Point{X: ent.Sprite.X - float32(rand.Float64())*10 - 15, Y: ent.Sprite.Y + float32(rand.NormFloat64())*5}
	} else {
		newPoint = util.Point{X: ent.Sprite.X + float32(rand.Float64())*-10 + 15, Y: ent.Sprite.Y + float32(rand.NormFloat64())*5}
	}

	xOffSet := ent.Sprite.SpriteWidth() / 2
	yOffSet := ent.Sprite.SpriteHeight() / 2

	newPoint = restrictTargetPointWithinBounds(xOffSet, yOffSet, newPoint, ent.CreatureData.TankBoundaries)
	return &newPoint
}

func restrictTargetPointWithinBounds(xOffSet, yOffSet int, point util.Point, rectangle image.Rectangle) util.Point {
	// The sprite center must stay within these bounds to keep sprite edges inside rectangle
	minX := float32(rectangle.Min.X + xOffSet) // left boundary + half width
	maxX := float32(rectangle.Max.X - xOffSet) // right boundary - half width
	minY := float32(rectangle.Min.Y + yOffSet) // top boundary + half height
	maxY := float32(rectangle.Max.Y - yOffSet) // bottom boundary - half height

	// Clamp the point to stay within bounds
	x := max(minX, min(maxX, point.X))
	y := max(minY, min(maxY, point.Y))

	return util.Point{X: x, Y: y, PType: point.PType, Tag: point.Tag}
}

func (e *Entity) CheckPointReached() {
	if e.CreatureData == nil {
		return
	}

	c := e.CreatureData
	s := e.Sprite

	if c.inBetweenPoint != nil {
		x := c.inBetweenPoint.X - s.X
		y := c.inBetweenPoint.Y - s.Y
		dist := math.Hypot(float64(x), float64(y))
		if dist < 5 {
			c.inBetweenPoint = nil
			e.Z = c.TargetZ
		}
		return
	}

	if c.TargetParticleId != 0 {
		xdist := c.TargetPoint.X - e.Sprite.X
		ydist := c.TargetPoint.Y - e.Sprite.Y

		dist := math.Hypot(float64(xdist), float64(ydist))

		if dist < 5 {
			e.Sprite.Y = c.TargetPoint.Y + float32(e.Sprite.GetSpriteRect().Dy()/4)
			e.PointReached(e.CreatureData.MovementFlags)
			c.energy = c.energy - 0.5
			if c.energy < 0 {
				c.energy = 0
			}
		}
	} else {
		x := c.TargetPoint.X - s.X
		y := c.TargetPoint.Y - s.Y
		dist := math.Hypot(float64(x), float64(y))
		if dist < 5 {
			e.PointReached(e.CreatureData.MovementFlags)
			c.energy = c.energy - 0.5
			if c.energy < 0 {
				c.energy = 0
			}
		}
	}
}

func (e *Entity) CheckAndLevelUp() {

	c := e.CreatureData

	if c.progress >= c.nextLevel && c.Size < 3 {
		c.Size += 1
		c.nextLevel *= 1.2
		c.progress = 0
		c.defaultMaxHunger += c.defaultMaxHunger / 3
		LoadFishSprite(e)
	}
}

func randomBool() bool {
	return rand.Intn(2) == 0
}

func (e *Entity) RandomTarget() *util.Point {
	//not actually ranomd at all, highly controlled behaviour

	c := e.CreatureData
	s := e.Sprite

	//didn't like having fish swim back and forth across the whole screen so i divide by 4 for smaller destination points
	// a carry direction thing could be set up to make each smaller point be in the same direction or something

	randomTargetX := randomX(s.X, s.GetSpriteRect().Dx(), c.TankBoundaries, e.CreatureData.Flags)

	randZ := rand.Intn(12)
	if !c.IsDecreasingZ() || c.IsIncreasingZ() {
		if randZ > 8 && randZ < 10 {
			if e.Z-1 > backFishLayer {
				c.SetDecreasingZ()
				//dont want layer 0/1 to be accesible by fish since back of tank is there
				c.TargetZ = e.Z - 1
				newPoint := util.Point{X: e.Sprite.X + float32(rand.NormFloat64())*10 + 5, Y: e.Sprite.Y + float32(rand.NormFloat64())*5}
				newPoint = restrictTargetPointWithinBounds(e.Sprite.GetSpriteRect().Dx()/2, e.Sprite.GetSpriteRect().Dy()/2, newPoint, e.CreatureData.TankBoundaries)
				return &newPoint
			}
		} else if randZ > 10 {
			if e.Z+1 < frontFishLayer {
				c.SetIncreasingZ()
				//last layer of tank
				c.TargetZ = e.Z + 1
				newPoint := util.Point{X: e.Sprite.X + float32(rand.NormFloat64())*10 + 5, Y: e.Sprite.Y + float32(rand.NormFloat64())*5}
				newPoint = restrictTargetPointWithinBounds(e.Sprite.GetSpriteRect().Dx()/2, e.Sprite.GetSpriteRect().Dy()/2, newPoint, e.CreatureData.TankBoundaries)
				return &newPoint
			}
		}
	} else {
		c.ClearZChange()
	}

	//normally distributed y based on avg depth stat
	//standard dev = entire tank?
	// lowest point (highest y) - (a randomly, normally distributed number * std dev(50) +
	//then we subtract the mean depth of our species and the height since were dealing with a left corner of sprite)
	offsetFromBottom := float32(rand.NormFloat64())*80 + c.avgDepth + float32(s.SpriteHeight()/2)
	println("fish offset =", offsetFromBottom)
	sample := float32(c.TankBoundaries.Max.Y) - offsetFromBottom
	//clamp below highest possible point
	randomTargetY := max(float32(c.TankBoundaries.Min.Y+s.SpriteHeight()), sample)
	randomTargetY = min(float32(c.TankBoundaries.Max.Y-s.SpriteHeight()), randomTargetY)

	targetX := randomTargetX
	targetY := randomTargetY

	fmt.Printf("random point generated of x: %f y: %f\n", targetX, targetY)

	newPoint := util.Point{X: targetX, Y: targetY}
	newPoint = restrictTargetPointWithinBounds(e.Sprite.GetSpriteRect().Dx()/2, e.Sprite.GetSpriteRect().Dy()/2, newPoint, e.CreatureData.TankBoundaries)
	return &newPoint

}

func CheckCollisionFlag() {
	/*	if flags[0] == 1 {
		// When hitting collision, pick a target that's definitely in the opposite direction
		if e.Sprite.Dx > 5 {
			e.Sprite.X -= 1 + float32(e.Sprite.GetSpriteRect().Dx())
			// Was moving right, now go left
			randomTargetX = s.X - float32(rand.Intn(100)+50) // 50-150 pixels left
		} else {
			e.Sprite.X += 1 + float32(e.Sprite.GetSpriteRect().Dx())
			randomTargetX = s.X + float32(rand.Intn(100)+50) // 50-150 pixels right
		}
		c.MovementFlags[0] = 0
		// Clamp to tank boundaries
		randomTargetX = max(float32(c.TankBoundaries.Min.X+s.SpriteWidth()/2), randomTargetX)
		randomTargetX = min(float32(c.TankBoundaries.Max.X-s.SpriteWidth()/2), randomTargetX)
	}*/
}

func GoRight(currentX float32, spriteWidth int, TankBoundaries image.Rectangle) float32 {
	nextX := currentX + float32(rand.Float64()*float64(TankBoundaries.Dx()/4))
	return min(float32(TankBoundaries.Max.X-spriteWidth/2), nextX)
}

func GoLeft(currentX float32, spriteWidth int, TankBoundaries image.Rectangle) float32 {
	nextX := currentX - float32(rand.Float64()*float64(TankBoundaries.Dx()/4))
	return max(float32(TankBoundaries.Min.X+spriteWidth/2), nextX)
}

func randomX(currentX float32, spriteWidth int, TankBoundaries image.Rectangle, Flags FishFlags) float32 {
	var randomTargetX float32

	if Flags&HitEdgeBoundary != 0 {
		switch Flags & FlagSwimmingDirection {
		case 0:
			randomTargetX = GoRight(currentX, spriteWidth, TankBoundaries)
			return randomTargetX
		case FlagSwimmingDirection:
			randomTargetX = GoLeft(currentX, spriteWidth, TankBoundaries)
			return randomTargetX
		}
	}

	//if recently changed direction dont change direction again
	if Flags&FlagChangedDirection == 0 {
		nextX := currentX + float32(rand.NormFloat64()*float64(TankBoundaries.Max.X/4)) //norm float 64 is between -1, 1
		randomTargetX = max(float32(TankBoundaries.Min.X+spriteWidth/2), nextX)
		randomTargetX = min(float32(TankBoundaries.Max.X-spriteWidth/2), randomTargetX)

	} else {
		switch Flags & FlagSwimmingDirection {
		case 0:
			randomTargetX = GoLeft(currentX, spriteWidth, TankBoundaries)
		case FlagSwimmingDirection:
			randomTargetX = GoRight(currentX, spriteWidth, TankBoundaries)
		}
	}

	return randomTargetX

}

func (e *Entity) CalculateVelocity(collisions []FishCollision) {
	c := e.CreatureData
	s := e.Sprite
	// Calculate the desired direction
	desiredDx := c.TargetPoint.X - s.X
	desiredDy := c.TargetPoint.Y - s.Y

	if c.inBetweenPoint != nil {
		desiredDx = c.inBetweenPoint.X - s.X
		desiredDy = c.inBetweenPoint.Y - s.Y
	}

	// Calculate distance BEFORE normalizing
	dist := math.Hypot(float64(desiredDx), float64(desiredDy))
	//arrivalRadius := float64(15) // Start slowing down at this distance

	// Normalize the direction
	length := float32(dist) // Use the distance we just calculated
	if length > 0 {
		desiredDx /= length
		desiredDy /= length
	}

	// Scale by desired speed with arrival behavior
	/*if dist < arrivalRadius {
		// Scale speed based on distance - closer = slower
		targetSpeed := c.speed * float32(dist/arrivalRadius)
		targetSpeed = max(targetSpeed, c.speed*0.8) // Minimum speed
		desiredDx *= targetSpeed
		desiredDy *= targetSpeed
	} else {*/
	desiredDx *= c.speed
	desiredDy *= c.speed

	steeringFactor := float32(0.05) // tweak for responsiveness

	s.Dx += (desiredDx - s.Dx) * steeringFactor
	s.Dy += (desiredDy - s.Dy) * steeringFactor

	/*if !TestNewPos(collisions, *s) {
		e.CreatureData.MovementFlags[0] = 1
		e.PointReached(e.CreatureData.MovementFlags)
	}*/

	s.ChangeAnimationSpeed(10 - (desiredDx * 3))
	if c.speed < 0.2 {
		s.ChangeAnimationSpeed(20)
	}

}

func TestNewPos(collisions []FishCollision, sp sprite.Sprite) bool {
	x := sp.X
	y := sp.Y

	x += sp.Dx
	y += sp.Dy

	rect := sp.GetSpriteRect()
	for _, col := range collisions {
		if rect.Overlaps(col.Rectangle) {
			return false
		}
	}
	return true
}

func (c *Entity) PointReached(flags [10]uint32) {
	if c.Z != c.CreatureData.TargetZ {
		c.Z = c.CreatureData.TargetZ
		ZSortEntities()
	}

	if c.CreatureData.TargetPoint.PType != util.Nada {
		c.CreatureData.ArrivedAtPoint = c.CreatureData.TargetParticleId
	}
	c.CreatureData.TargetPoint = nil
	c.Transition()

}

func (e *Entity) TranSlateFishShaderOpts() *ebiten.DrawRectShaderOptions {
	if e.CreatureData == nil {
		return nil
	}
	s := e.Sprite
	c := e.CreatureData
	opts := &ebiten.DrawRectShaderOptions{}

	if c.IsGoingRight() {
		e.Sprite.Flip = true
		opts.GeoM.Scale(-1, 1) // flip horizontally
		opts.GeoM.Translate(float64(e.Sprite.GetSpriteRect().Dx()), 0)
	}
	if c.IsGoingLeft() {
		e.Sprite.Flip = false
	}

	if e.Sprite.Dy < -0.5 {
		if c.IsGoingRight() {
			opts.GeoM.Rotate(-0.3)
		} else {
			opts.GeoM.Rotate(0.3)
		}
	}

	if e.Sprite.Dy > 0.5 {
		if c.IsGoingRight() {
			opts.GeoM.Rotate(0.3)
		} else {
			opts.GeoM.Rotate(-0.3)
		}
	} /*
		b := c.Img.Bounds()
		midpoint := float32(b.Dy() / 2)*/

	opts.GeoM.Translate(float64(s.X-float32(s.SpriteWidth()/2)), float64(s.Y-float32(s.SpriteHeight()/2)))

	return opts
}

func (ent *Entity) SortParticleQueue() {
	slices.SortFunc(ent.CreatureData.ParticlePointQueue, func(a, b uint32) int {
		// Calculate Z distance from entity

		aPoint, exists := GetEntity(a)
		if !exists {
			return 0
		}
		bPoint, exists := GetEntity(b)
		if !exists {
			return 0
		}

		aZDist := math.Abs(float64(aPoint.Z - ent.Z))
		bZDist := math.Abs(float64(bPoint.Z - ent.Z))

		// First priority: sort by Z distance
		if aZDist != bZDist {
			if aZDist < bZDist {
				return -1 // a is closer in Z
			}
			return 1 // b is closer in Z
		}

		// If Z distance is equal, sort by XY distance
		aDist := entityEntityDistanceFunc(*ent, *aPoint)
		bDist := entityEntityDistanceFunc(*ent, *bPoint)

		if aDist < bDist {
			return -1 // a is closer in XY
		} else if aDist > bDist {
			return 1 // b is closer in XY
		}
		return 0 // equal distance
	})
}

func entityEntityDistanceFunc(ent1 Entity, ent2 Entity) float64 {
	return DistanceFunc(ent1.Sprite.X, ent2.Sprite.X, ent1.Sprite.Y, ent2.Sprite.Y)
}

func entityPointDistanceFunc(ent1 Entity, point util.Point) float64 {
	return DistanceFunc(ent1.Sprite.X, point.X, ent1.Sprite.Y, point.Y)
}

func (ent *Entity) TranSlateFishOpts(options ebiten.DrawImageOptions) ebiten.DrawImageOptions {
	if ent.CreatureData == nil {
		return options
	}

	s := ent.Sprite
	c := ent.CreatureData

	opts := ebiten.DrawImageOptions{}

	if c.IsGoingRight() {
		ent.Sprite.Flip = true
		opts.GeoM.Scale(-1, 1)
		opts.GeoM.Translate(float64(ent.Sprite.GetSpriteRect().Dx()), 0) // flip horizontally
	}
	if c.IsGoingLeft() {
		ent.Sprite.Flip = false
	}

	if ent.Sprite.Dy < -0.5 {
		if c.IsGoingRight() {
			opts.GeoM.Rotate(-0.3)
		} else {
			opts.GeoM.Rotate(0.3)
		}
	}

	if ent.Sprite.Dy > 0.5 {
		if c.IsGoingRight() {
			opts.GeoM.Rotate(0.3)
		} else {
			opts.GeoM.Rotate(-0.3)
		}
	}

	opts.GeoM.Translate(float64(s.X-float32(s.SpriteWidth()/2)), float64(s.Y-float32(s.SpriteHeight()/2)))

	if registry.Config.Zoom {
		sprite.ColorScaleBasedOnZ(ent.Z, &opts)
	}

	return opts
}

func (e *Entity) SetTargetPoint(point *util.Point) {
	if e.CreatureData == nil {
		return
	}

	if point == nil {
		log.Fatal("Nil point sent to fish to make target point")
	}
	c := e.CreatureData

	c.TargetPoint = point

	if c.TargetPoint.PType == util.Food {
		if c.TargetPoint.Z != 0 {
			c.TargetZ = c.TargetPoint.Z
		}
	}

	if c.Flip {
		if point.X-dontTurnAroundDist < e.Sprite.X {
			c.SetDirection(Left)
			c.Flip = false
			c.inBetweenPoint = SetTurningPoint(e)
			c.SetRecentlyChangedDirection()
		} else {
			c.ClearRecentlyChangedDirection()
		}
	} else {
		if point.X+dontTurnAroundDist > e.Sprite.X {
			c.SetDirection(Right)
			c.Flip = true
			c.inBetweenPoint = SetTurningPoint(e)
			c.SetRecentlyChangedDirection()
		} else {
			c.ClearRecentlyChangedDirection()
		}
	}
	c.ClearHitBoundary()
	e.calcSpeed()
}

func (e *Entity) publishStats(sendTo string) {
	if e.CreatureData == nil {
		return
	}

	c := e.CreatureData

	ev := SendData{}
	ev.DataFor = sendTo

	/*var TargetPoint string*/
	var state string
	var healthState string

	switch c.State {
	case Resting:
		state = "resting"
	case Swimming:
		state = "swimming"
	case Eating:
		state = "eating "
	}
	/*	TargetPoint = fmt.Sprintf("Target Point: %d, %d", int(c.TargetPoint.X), int(c.TargetPoint.Y))
	 */
	switch c.HealthState {
	case Healthy:
		healthState = "healthy"
	case Stressed:
		healthState = "stressed"
	case Sick:
		healthState = "sick"
	}

	nameString := fmt.Sprintf("Name: %s\n", c.name)
	hungerString := fmt.Sprintf("Hunger : %d/%d\n", int(c.Hunger), int(c.MaxHunger))
	/*energyString := fmt.Sprintf("Energy : %d/%d\n", int(c.energy), int(c.maxEnergy))*/
	SizeString := fmt.Sprintf("Size : %d\n", c.Size)
	experienceString := fmt.Sprintf("Growth : %d/%d\n", int(c.progress), int(c.nextLevel))
	stateString := fmt.Sprintf("State: %s\n", state)
	healthStateString := fmt.Sprintf("Health State: %s\n", healthState)
	stressString := ""
	targetZ := fmt.Sprintf("target Z: %d\n", e.CreatureData.TargetZ)
	currentZ := fmt.Sprintf("current Z: %d\n", e.Z)

	for i, fact := range e.CreatureData.stressContributors {
		stressString += fmt.Sprintf("Stress Factor%d:. %s\n", i+1, fact)
	}

	ev.Data = nameString + stateString + SizeString + hungerString +
		experienceString + healthStateString + stressString + targetZ + currentZ

	/*	if TargetPoint != "" {
		ev.Data += TargetPoint
	}*/

	c.EventHub.Publish(ev)
}

func GameFishToSaveFish(creature *Entity) SavedFish {

	var s SavedFish

	s.Name = creature.CreatureData.name
	s.Size = creature.CreatureData.Size
	s.Progress = creature.CreatureData.progress
	s.NextLevel = creature.CreatureData.nextLevel
	s.FishType = string(creature.CreatureData.FishType)

	return s
}

const BaseFishSpeed = 0.5

type FishStats struct {
	name             string
	Stress           float32
	Happiness        float32
	Hunger           int
	MaxHunger        int
	age              int
	defaultMaxHunger int
	maxEnergy        float32
	energy           float32
	maxSpeed         float32
	avgSpeed         float32
	stdDevSpeed      float32
	avgDepth         float32
	speed            float32
	Size             int
	progress         float32
	nextLevel        float32
	IdealTemperature int
	IdealPH          float64
	DaysStressed     int
	DaysSick         int
	HealthState
	Personality FishPersonality
	FishType    FishList
}

func GenFishStats(fType FishList, name string) (*FishStats, error) {
	switch fType {
	case MollyFish:
		println("loading molly GoldFish")
		fs, err := GenMollyFishStats()
		fs.name = name
		fs.age = 1
		if err != nil {
			return fs, err
		}
		return fs, nil
	case GoldFish:
		println("loading gold GoldFish")
		fs, err := GenGoldFishStats()
		fs.name = name
		fs.age = 1
		if err != nil {
			return fs, err
		}
		return fs, nil
	case Guppy:
		fs, err := GenGuppyFishStats()
		fs.name = name
		fs.age = 1
		if err != nil {
			return fs, err
		}
		return fs, nil
	case Kirbensis:
		fs, err := GenKirbensisFishStats()
		fs.FishType = Kirbensis
		fs.age = 1
		fs.name = name
		if err != nil {
			return fs, err
		}
		return fs, nil
	default:
		//no fish stat generator for this species yet
		log.Println("Warning: No fish stat generator for:", string(fType))
		fs, err := GenKirbensisFishStats()
		fs.age = 1
		fs.FishType = fType
		fs.name = name
		if err != nil {
			return fs, err
		}
		return fs, nil
	}
}

func GenMollyFishStats() (*FishStats, error) {
	fs := &FishStats{}

	fs.Size = 1
	fs.maxSpeed = rand.Float32() + 0.7
	fs.avgSpeed = BaseFishSpeed + 0.2
	fs.IdealPH = 7.2
	fs.stdDevSpeed = 0.2
	fs.speed = rand.Float32()*fs.maxSpeed + 0.3
	fs.FishType = MollyFish
	fs.maxEnergy = 25
	fs.energy = fs.maxEnergy / 2
	fs.Hunger = 0
	fs.MaxHunger = 5
	fs.defaultMaxHunger = fs.MaxHunger
	fs.avgDepth = 100
	fs.progress = 0
	fs.nextLevel = 10
	fs.IdealTemperature = 75

	persRoll := rand.Intn(10)

	if persRoll < 4 {
		fs.Personality = social
	} else {
		fs.Personality = shy
	}

	return fs, nil
}

func GenGoldFishStats() (*FishStats, error) {
	fs := &FishStats{}
	fs.Size = 1
	fs.IdealTemperature = 70
	fs.IdealPH = 6.5
	fs.avgDepth = 40.0
	fs.avgSpeed = BaseFishSpeed
	fs.stdDevSpeed = 0.2
	fs.maxSpeed = rand.Float32()*0.5 + 0.2
	fs.speed = rand.Float32()*fs.maxSpeed + 0.3
	fs.FishType = GoldFish
	fs.maxEnergy = 25
	fs.energy = fs.maxEnergy / 2
	fs.Hunger = 0
	fs.progress = 0
	fs.nextLevel = 10
	fs.IdealPH = 6.5
	fs.MaxHunger = 5
	fs.defaultMaxHunger = fs.MaxHunger
	persRoll := rand.Intn(10)

	if persRoll < 8 {
		fs.Personality = social
	} else {
		fs.Personality = shy
	}

	return fs, nil
}

func GenGuppyFishStats() (*FishStats, error) {
	fs := &FishStats{}
	fs.Size = 1
	fs.avgDepth = 150
	fs.avgSpeed = BaseFishSpeed
	fs.stdDevSpeed = 0.1
	fs.speed = rand.Float32()*fs.stdDevSpeed + fs.avgSpeed
	fs.FishType = Guppy
	fs.maxEnergy = 25
	fs.energy = fs.maxEnergy / 2
	fs.Hunger = 0
	fs.progress = 0
	fs.nextLevel = 10
	fs.IdealTemperature = 80
	fs.IdealPH = 7.5
	fs.MaxHunger = 6
	fs.defaultMaxHunger = fs.MaxHunger
	persRoll := rand.Intn(10)

	if persRoll < 8 {
		fs.Personality = social
	} else {
		fs.Personality = shy
	}

	return fs, nil
}

func GenKirbensisFishStats() (*FishStats, error) {
	fs := &FishStats{}
	fs.Size = 1
	fs.avgDepth = 100
	fs.avgSpeed = BaseFishSpeed
	fs.stdDevSpeed = 0.25
	fs.maxSpeed = 0.8
	fs.speed = rand.Float32()*fs.maxSpeed + 0.3
	fs.FishType = Kirbensis
	fs.maxEnergy = 25
	fs.energy = fs.maxEnergy / 2
	fs.Hunger = 0
	fs.progress = 0
	fs.nextLevel = 10
	fs.IdealTemperature = 77
	fs.IdealPH = 7.0
	fs.MaxHunger = 4
	fs.defaultMaxHunger = fs.MaxHunger
	persRoll := rand.Intn(10)

	if persRoll < 8 {
		fs.Personality = social
	} else {
		fs.Personality = shy
	}

	return fs, nil
}

func (e *Entity) Add1expGraphic() {
	spEff := entImportableLoaders.LoadStaticEffect("+1", e.Sprite.X, e.Sprite.Y, "")
	params := make(map[string]any)
	params["opacity"] = float32(1.0)
	graphics.NewSpriteGraphic(*spEff, MoveSpriteToDestinationUp, params)
}

func MoveSpriteToDestinationUp(sp *graphics.SpriteGraphic) {
	sp.SetDrawFunc(graphics.FadeIn)

	destinationX := float64(sp.Sprite.X)
	destinationY := 100.0
	speed := 2.0

	// Calculate rotation needed to reach π (flipped)
	alpha := sp.Parameters["opacity"].(float32)
	alpha -= 0.02
	sp.Parameters["opacity"] = alpha
	// Calculate the distance to destination
	dx := destinationX - float64(sp.Sprite.X)
	dy := destinationY - float64(sp.Sprite.Y)

	// Calculate the total distance
	distance := math.Sqrt(dx*dx + dy*dy)

	// If we're close enough, stop moving
	if distance < speed {
		sp.Sprite.X = float32(destinationX)
		sp.Sprite.Y = float32(destinationY)
		graphics.DeInitGraphicId(sp.Id)
		return
	}
	sp.Sprite.DOptsUpdaterTag = "opacity"

	sp.Sprite.X += float32(dx / distance * speed)
	sp.Sprite.Y += float32(dy / distance * speed)

}

func CheckIfAllFishFed() bool {
	fed := true

	for _, ent := range LiveList {
		if ent.CreatureData != nil {
			if ent.CreatureData.Hunger < ent.CreatureData.MaxHunger {
				fed = false
			}
		}
	}
	return fed
}
