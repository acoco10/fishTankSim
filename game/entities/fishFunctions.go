package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"math"
	"math/rand"
)

type HealthState uint8

const (
	Healthy HealthState = iota
	Stressed
	Sick
	Dead
)

func CreatureEventSubscriptions(c *Entity) {

	c.EventHub.Subscribe(PointGenerated{}, func(e tasks.Event) {
		ev := e.(PointGenerated)
		point, exists := GetEntity(ev.PointId)
		if !exists {
			log.Fatal("A recently created point should exist....")
		}
		if point.ParticleData.PType == util.Food && !c.CreatureData.TickClicked {
			c.CreatureData.ParticlePointQueue[ev.PointId] = point.ParticleData.Point
			c.CreatureData.TickClicked = true
			if c.CreatureData.Hunger < c.CreatureData.MaxHunger {
				c.goToFood()
			}
		}

	})

	c.EventHub.Subscribe(CreatureReachedPoint{}, func(e tasks.Event) {
		ev := e.(CreatureReachedPoint)
		delete(c.CreatureData.ParticlePointQueue, ev.PointID)
		if c.CreatureData.Hunger < c.CreatureData.MaxHunger && ev.CreatureID != c.Id {
			c.goToFood()
		}
	})

	c.CreatureData.EventHub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		c.CreatureData.Hunger = 0
		c.CreatureData.energy = c.CreatureData.maxEnergy
	})

	c.EventHub.Subscribe(events.NewProp{}, func(e tasks.Event) {
		ev := e.(events.NewProp)

		LoadFollowEffectAsEnt("exclamation", c.Id, c.EventHub)

		prop, exists := GetEntity(ev.PropId)
		if !exists {
			log.Fatal("some weird shit happened when fish was attracted to a new structure")
		}

		TargetPoint := &util.Point{X: float32(prop.Sprite.X) + 20, Y: float32(prop.Sprite.Y) + float32(prop.Sprite.GetSpriteRect().Dy()/2), PType: util.Structure}
		println("making creature target point new prop")
		c.CreatureData.TargetParticleId = ev.PropId
		c.CreatureData.ParticlePointQueue[ev.PropId] = TargetPoint
		c.MakeTargetPoint(TargetPoint)

		c.CreatureData.TargetZ = min(prop.Z+1, 12)
	})

	c.CreatureData.EventHub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		// add normal-map bac.CreatureData. when it's a new day since it was turned off for night-scene
		//c.CreatureData.Shader = registry.ShaderMap["NormalMap"]
		println("new day  received  for fish:", c)
		c.CreatureData.CalcDailyFishHealthState()
		c.CheckAndLevelUp()
	})
}

func (c *CreatureData) CalcDailyFishHealthState() {
	if math.Abs(float64(c.Environment.Temperature-c.idealTemperature)) > 10 {
		c.stressContributors = append(c.stressContributors, "temperature")
		c.Happiness -= 1
		c.Stress += 2
	}
	//compare environment ph to ideal
	if math.Abs(float64(c.Environment.NaturalPHLevel-c.idealPH)) > 1.5 {
		c.stressContributors = append(c.stressContributors, "ph")
		c.Stress += 2
	}
	// reduce hunger if chronically stressed
	if c.Stress > 3 {
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

func DoneEating(c *Entity) {
	if c.CreatureData.Hunger < c.CreatureData.MaxHunger {
		c.goToFood()
	} else {
		c.MakeTargetPoint(c.RandomTarget(c.CreatureData.MovementFlags))
	}
}

func (e *Entity) otherFishPoint(point *util.Point) {
	c := e.CreatureData

	delete(c.ParticlePointQueue, c.TargetParticleId)

	if len(c.ParticlePointQueue) > 0 {
		e.goToFood()
	}

}

func (c *Entity) goToFood() {
	if c.CreatureData.Hunger < c.CreatureData.MaxHunger {
		newTargid := ClosestParticle(c.Sprite.X, c.Sprite.Y, c.CreatureData.ParticlePointQueue)
		c.CreatureData.TargetParticleId = newTargid
		if newTargid == 0 {
			c.MakeTargetPoint(c.RandomTarget(c.CreatureData.MovementFlags))
		} else {
			targPoint := c.CreatureData.ParticlePointQueue[newTargid]
			c.MakeTargetPoint(targPoint)
		}
	} else {
		c.MakeTargetPoint(c.RandomTarget(c.CreatureData.MovementFlags))
	}
}

func (e *Entity) calcSpeed() {

	if e.CreatureData == nil {
		return
	}
	c := e.CreatureData
	const minSpeed = 0.5 // farthest fish still move at 20% speed
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
	e.AddRandomMovement(collisions)

	e.Sprite.X += e.Sprite.Dx
	e.Sprite.Y += e.Sprite.Dy

	e.CreatureData.distanceTraveled += 1

	e.EnforceBoundaries()

	e.CheckPointReached()

}

func (c *Entity) EnforceBoundaries() {
	s := c.Sprite
	bounds := c.CreatureData.TankBoundaries

	s.X = max(float32(bounds.Min.X), s.X)
	s.Y = max(float32(bounds.Min.Y), s.Y)
	s.X = min(float32(bounds.Max.X), s.X)
	s.Y = min(float32(bounds.Max.Y), s.Y)
}

func (e *Entity) CheckPointReached() {
	if e.CreatureData == nil {
		return
	}

	c := e.CreatureData
	s := e.Sprite

	if c.TargetParticleId != 0 {
		xdist := c.TargetPoint.X - e.Sprite.X
		ydist := c.TargetPoint.Y - e.Sprite.Y

		dist := math.Hypot(float64(xdist), float64(ydist))

		if dist < 5 {
			e.Sprite.X -= 4
			if e.Sprite.Flip {
				e.Sprite.X += 4
			}
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

func (e *Entity) RandomTarget(flags [10]uint32) *util.Point {

	c := e.CreatureData
	s := e.Sprite

	//didn't like having fish swim back and forth across the whole screen so i divide by 4 for smaller destination points
	// a carry direction thing could be set up to make each smaller point be in the same direction or something

	randomTargetX := randomX(s.X, s.SpriteWidth(), c.TankBoundaries)

	if flags[0] == 1 {
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
		randomTargetX = max(float32(c.TankBoundaries.Min.X+s.SpriteWidth()+5), randomTargetX)
		randomTargetX = min(float32(c.TankBoundaries.Max.X-s.SpriteWidth()-5), randomTargetX)
	}

	randZ := rand.Intn(12)
	if randZ > 8 && randZ < 10 {
		if c.TargetZ-1 >= 1 {
			c.TargetZ -= 1
		}
	} else if randZ > 10 {
		if c.TargetZ+1 <= 12 {
			c.TargetZ += 1
		}
	}

	//normally distributed y based on avg depth stat
	//standard dev = entire tank?
	// lowest point (highest y) - (a randomly, normally distributed number * std dev(50) +
	//then we subtract the mean depth of our species and the height since were dealing with a left corner of sprite)
	offsetFromBottom := float32(rand.NormFloat64())*20 + c.avgDepth + float32(s.SpriteHeight())
	println("fish offset =", offsetFromBottom)
	sample := float32(c.TankBoundaries.Max.Y) - offsetFromBottom
	//clamp below highest possible point
	randomTargetY := max(float32(c.TankBoundaries.Min.Y+s.SpriteHeight()), sample)
	randomTargetY = min(float32(c.TankBoundaries.Max.Y-s.SpriteHeight()), randomTargetY)

	targetX := randomTargetX
	targetY := randomTargetY

	fmt.Printf("random point generated of x: %f y: %f\n", targetX, targetY)

	newPoint := util.Point{X: targetX, Y: targetY}
	return &newPoint
}

func randomX(currentX float32, spriteWidth int, TankBoundaries image.Rectangle) float32 {
	nextX := currentX + float32(rand.NormFloat64()*float64(TankBoundaries.Max.X/4)) //norm float 64 is between -1, 1
	randomTargetX := max(float32(TankBoundaries.Min.X+spriteWidth+5), nextX)

	if randomTargetX > float32(TankBoundaries.Max.X-spriteWidth) {
		randomTargetX = float32(TankBoundaries.Max.X - spriteWidth - 5)
	}

	return randomTargetX
}

func (e *Entity) AddRandomMovement(collisions []FishCollision) {
	c := e.CreatureData
	s := e.Sprite
	// Calculate the desired direction
	desiredDx := c.TargetPoint.X - s.X
	desiredDy := c.TargetPoint.Y - s.Y

	// Calculate distance BEFORE normalizing
	dist := math.Hypot(float64(desiredDx), float64(desiredDy))
	arrivalRadius := float64(15) // Start slowing down at this distance

	// Normalize the direction
	length := float32(dist) // Use the distance we just calculated
	if length > 0 {
		desiredDx /= length
		desiredDy /= length
	}

	// Scale by desired speed with arrival behavior
	if dist < arrivalRadius {
		// Scale speed based on distance - closer = slower
		targetSpeed := c.speed * float32(dist/arrivalRadius)
		targetSpeed = max(targetSpeed, c.speed*0.4) // Minimum speed
		desiredDx *= targetSpeed
		desiredDy *= targetSpeed
	} else {
		desiredDx *= c.speed
		desiredDy *= c.speed
	}

	steeringFactor := float32(0.05) // tweak for responsiveness

	s.Dx += (desiredDx - s.Dx) * steeringFactor
	s.Dy += (desiredDy - s.Dy) * steeringFactor

	if !TestNewPos(collisions, *s) {
		e.CreatureData.MovementFlags[0] = 1
		e.PointReached(e.CreatureData.MovementFlags)
	}

	s.ChangeAnimationSpeed(10 - (desiredDx * 3))
	if c.speed < 0.2 {
		s.ChangeAnimationSpeed(30)
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
	c.Z = c.CreatureData.TargetZ
	c.CreatureData.distanceTraveled = 0
	if flags[0] == 1 {
		c.MakeTargetPoint(c.RandomTarget(flags))
	}

	if c.CreatureData.TargetParticleId != 0 {
		if c.CreatureData.TargetPoint.PType == util.Food {
			c.CreatureData.Hunger++
			c.Add1expGraphic()
			c.CreatureData.progress += 1
			c.CreatureData.State = Eating
		}
		ev := CreatureReachedPoint{
			PointID:    c.CreatureData.TargetParticleId,
			CreatureID: c.Id,
		}
		c.CreatureData.EventHub.Publish(ev)
		c.CreatureData.TargetParticleId = 0
		return
	} else if c.CreatureData.TargetPoint.PType == util.Structure {
		c.CreatureData.TargetParticleId = 0
		c.MakeTargetPoint(c.RandomTarget(flags))
	} else {
		c.MakeTargetPoint(c.RandomTarget(flags))
	}
	//reorg target point and q

}

func (e *Entity) TranSlateFishShaderOpts() *ebiten.DrawRectShaderOptions {
	if e.CreatureData == nil {
		return nil
	}

	c := e.CreatureData
	opts := &ebiten.DrawRectShaderOptions{}

	if c.Flip {
		e.Sprite.Flip = true
		opts.GeoM.Scale(-1, 1) // flip horizontally
		opts.GeoM.Translate(float64(e.Sprite.SpriteWidth()), 0)
	} else {
		e.Sprite.Flip = false
	}

	if e.Sprite.Dy < -0.5 {
		if c.Flip {
			opts.GeoM.Rotate(-0.3)
		} else {
			opts.GeoM.Rotate(0.3)
		}
	}
	if e.Sprite.Dy > 0.5 {
		if c.Flip {
			opts.GeoM.Rotate(0.3)
		} else {
			opts.GeoM.Rotate(-0.3)
		}
	} /*
		b := c.Img.Bounds()
		midpoint := float32(b.Dy() / 2)*/

	y := float64(e.Sprite.Y - float32(e.Sprite.SpriteHeight()/2))
	x := float64(e.Sprite.X)

	if c.Flip {
		x = x - float64(e.Sprite.SpriteWidth())
	}
	opts.GeoM.Translate(x, y)

	return opts
}

func (e *Entity) TranSlateFishOpts() *ebiten.DrawImageOptions {
	if e.CreatureData == nil {
		return nil
	}

	s := e.Sprite
	c := e.CreatureData

	opts := &ebiten.DrawImageOptions{}

	if c.Flip {
		sprite.FlipSprite(e.Sprite, opts)
	}

	if e.Sprite.Dy < -0.5 {
		if c.Flip {
			opts.GeoM.Rotate(-0.3)
		} else {
			opts.GeoM.Rotate(0.3)
		}
	}

	if e.Sprite.Dy > 0.5 {
		if c.Flip {
			opts.GeoM.Rotate(0.3)
		} else {
			opts.GeoM.Rotate(-0.3)
		}
	}
	if c.Flip {
		opts.GeoM.Translate(float64(s.X-float32(s.SpriteWidth())), float64(s.Y))
	} else {
		opts.GeoM.Translate(float64(s.X), float64(s.Y))
	}

	return opts
}

func (e *Entity) MakeTargetPoint(point *util.Point) {
	if e.CreatureData == nil {
		return
	}

	if point == nil {
		log.Fatal("Nil point sent to fish to make target point")
	}

	c := e.CreatureData
	if c.Flip {
		if point.X-20 < e.Sprite.X {
			c.Flip = false
		}
	} else {
		if point.X+20 > e.Sprite.X {
			c.Flip = true
		}
	}

	c.TargetPoint = point
	if c.TargetPoint.PType == util.Food {
		if c.TargetPoint.Z != 0 {
			c.TargetZ = c.TargetPoint.Z
		}
	}
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
		stressString += fmt.Sprintf("%d. %s\n", i+1, fact)
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

type FishStats struct {
	name             string
	Stress           float32
	Happiness        float32
	Hunger           int
	MaxHunger        int
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
	idealTemperature int
	idealPH          float64
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
		if err != nil {
			return fs, err
		}
		return fs, nil
	case GoldFish:
		println("loading gold GoldFish")
		fs, err := GenGoldFishStats()
		fs.name = name
		if err != nil {
			return fs, err
		}
		return fs, nil
	case Guppy:
		fs, err := GenGuppyFishStats()
		fs.name = name
		if err != nil {
			return fs, err
		}
		return fs, nil
	case Kirbensis:
		fs, err := GenKirbensisFishStats()
		fs.FishType = Kirbensis
		fs.name = name
		if err != nil {
			return fs, err
		}
		return fs, nil
	default:
		//no fish stat generator for this species yet
		log.Println("Warning: No fish stat generator for:", string(fType))
		fs, err := GenKirbensisFishStats()
		fs.FishType = Kirbensis
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
	fs.avgSpeed = 1.2
	fs.stdDevSpeed = 0.1
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
	fs.idealTemperature = 75

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
	fs.idealTemperature = 70
	fs.idealPH = 6.5
	fs.avgDepth = 40.0
	fs.avgSpeed = 1.0
	fs.stdDevSpeed = 0.2
	fs.maxSpeed = rand.Float32()*0.5 + 0.2
	fs.speed = rand.Float32()*fs.maxSpeed + 0.3
	fs.FishType = GoldFish
	fs.maxEnergy = 25
	fs.energy = fs.maxEnergy / 2
	fs.Hunger = 0
	fs.progress = 0
	fs.nextLevel = 10
	fs.idealPH = 6.5
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
	fs.avgSpeed = 1.1
	fs.stdDevSpeed = 0.1
	fs.maxSpeed = rand.Float32()*0.5 + 0.2
	fs.speed = rand.Float32()*fs.maxSpeed + 0.3
	fs.FishType = Guppy
	fs.maxEnergy = 25
	fs.energy = fs.maxEnergy / 2
	fs.Hunger = 0
	fs.progress = 0
	fs.nextLevel = 10
	fs.idealTemperature = 80
	fs.idealPH = 7.5
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
	fs.avgSpeed = 0.9
	fs.stdDevSpeed = 0.25
	fs.maxSpeed = rand.Float32()*0.5 + 0.2
	fs.speed = rand.Float32()*fs.maxSpeed + 0.3
	fs.FishType = Guppy
	fs.maxEnergy = 25
	fs.energy = fs.maxEnergy / 2
	fs.Hunger = 0
	fs.progress = 0
	fs.nextLevel = 10
	fs.idealTemperature = 77
	fs.idealPH = 7.0
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
