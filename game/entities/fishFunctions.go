package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math"
	"math/rand"
	"sort"
)

type HealthState uint8

const (
	Healthy HealthState = iota
	Stressed
	Sick
	Dead
)

func CreatureEventSubscriptions(c *Creature) {

	c.EventHub.Subscribe(PointGenerated{}, func(e tasks.Event) {
		ev := e.(PointGenerated)
		if ev.Point.PType == geometry.Food && !c.TickClicked {
			c.TickClicked = true
			if c.Hunger > 0 {
				c.AddTargetPointToQueue(ev.Point)
			}
		}
	})

	c.EventHub.Subscribe(CreatureReachedPoint{}, func(e tasks.Event) {
		ev := e.(CreatureReachedPoint)
		if ev.Creature == c {
			c.ownPointReached(ev)
		} else {
			c.otherFishPoint(ev)
		}
	})

	c.EventHub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		c.Hunger = c.maxHunger
		c.energy = c.maxEnergy

		c.Shader = nil
		c.CheckAndLevelUp()
	})

	c.EventHub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		// add normal-map back when it's a new day since it was turned off for night-scene
		c.Shader = registry.ShaderMap["NormalMap"]
		println("new day received for fish:", c)
		c.CalcDailyFishHealthState()
	})
}

func (c *Creature) CalcDailyFishHealthState() {
	if math.Abs(float64(c.Environment.Temperature-c.idealTemperature)) > 10 {
		c.Happiness -= 1
		c.Stress += 2
	}
	//compare environment ph to ideal
	if math.Abs(float64(c.Environment.NaturalPHLevel-c.idealPH)) > 1.5 {
		c.Stress += 2
	}
	// reduce hunger if chronically stressed
	if c.Stress > 3 {
		c.maxHunger = 3
		c.DaysStressed++
	}
	//check if fish gets sick if chronically stressed
	if c.DaysStressed > 3 && c.HealthState != Sick {
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

func (c *Creature) Highlighted() bool {
	return c.Selected
}

func (c *Creature) ownPointReached(ev CreatureReachedPoint) {
	//more points?deque this one safely
	if len(c.PointQueue) > 0 {
		c.PointQueue = c.PointQueue[1:]

	}

	switch ev.Point.PType {
	//reached food: eat
	case geometry.Food:
		println("reached own point, eating food and going for next piece")
		c.progress += 1
		c.Hunger -= 1
		if c.Hunger < 0 {
			c.Hunger = 0
		}
		c.State = Eating
		if len(c.PointQueue) > 0 {
			//more points? still hungry? go eat
			if c.Hunger > 0 {
				c.goToFood()
			}
		}

	default:
		println("reached own point, setting random speed and random next target")
		restCheck := rand.Intn(2)
		if c.energy == 0 || restCheck > 1 {
			c.State = Resting
		}
		//newTarg := c.RandomTarget()
		//c.PointQueue = append(c.PointQueue, newTarg)
		c.calcSpeed()
		//other creature behaviour
	}
	c.NextPoint()
}

func (c *Creature) otherFishPoint(ev CreatureReachedPoint) {
	pointThere := false
	var pointIndex int

	for i, point := range c.PointQueue {
		if ev.Point == point {
			pointThere = true
			pointIndex = i
		}
	}

	endIndex := len(c.PointQueue) - 1

	if pointThere {
		switch ev.Point.PType {
		case geometry.Food:
			c.PointQueue[pointIndex] = c.PointQueue[endIndex]
			c.PointQueue = c.PointQueue[0:endIndex]
		default:
			//placeholder, no modification when another Fish reaches their point that's not food as of now
		}
	}
	c.NextPoint()
}

func (c *Creature) goToFood() {
	c.sortPoints()
	c.calcSpeed()
}

func (c *Creature) NextPoint() {
	if len(c.PointQueue) == 0 {
		c.AddTargetPointToQueue(c.RandomTarget())
	}

	//scaling needs smoothing
	/*	x := rand.Intn(3)

		switch x {

		case 0:
			c.Scale = 1.0
		case 1:
			c.Scale = 1.0
		case 2:
			c.Scale = 0.9
		}*/

	if c.Flip && c.PointQueue[0].X-c.X < -50 {
		c.Flip = false
	}

	if !c.Flip && c.PointQueue[0].X-c.X > 50 {
		c.Flip = true
	}

}

func (c *Creature) calcSpeed() {
	if len(c.PointQueue) < 0 {
		if c.PointQueue[0].PType == geometry.Food {
			c.speed = c.maxSpeed
		}
	} else {
		c.speed = float32(math.Min(rand.Float64()*float64(c.maxSpeed)+float64(c.avgSpeed), float64(c.maxSpeed)))
	}

	fmt.Printf("random speed generated = %f\n", c.speed)

}

func DistanceFunc(x, x2, y, y2 float32) float64 {
	xDis := float64(x - x2)
	yDis := float64(y - y2)
	dis := math.Sqrt(math.Pow(xDis, 2) + math.Pow(yDis, 2))
	return dis
}

func (c *Creature) Type() geometry.InterestPoint {
	return geometry.OtherCreature
}

func (c *Creature) sortPoints() {

	sort.Slice(c.PointQueue, func(i, j int) bool {

		xI, yI := c.PointQueue[i].PointCoord()
		xJ, yJ := c.PointQueue[j].PointCoord()

		distI := math.Hypot(float64(c.X-xI), float64(c.Y-yI))
		distJ := math.Hypot(float64(c.X-xJ), float64(c.Y-yJ))

		pTypeI := c.PointQueue[i].PType
		pTypeJ := c.PointQueue[j].PType

		if pTypeI == geometry.Food && pTypeJ != geometry.Food {
			return true
		}
		return distI < distJ
	})

}

func (c *Creature) Move() {

	c.AddRandomMovement()

	c.X += c.Dx
	c.Y += c.Dy

	c.EnforceBoundaries()

	c.UpdateToNextPoint()

}

func (c *Creature) EnforceBoundaries() {
	c.X = max(c.TankBoundaries.X1, c.X)
	c.Y = max(c.TankBoundaries.Y1, c.Y)

	c.X = min(c.TankBoundaries.X2, c.X)
	c.Y = min(c.TankBoundaries.Y2, c.Y)
}

func (c *Creature) UpdateToNextPoint() {
	if len(c.PointQueue) > 0 {

		tgtPoint := c.PointQueue[0]

		x := tgtPoint.X - c.X
		y := tgtPoint.Y - c.Y

		dist := math.Hypot(float64(x), float64(y))

		if dist < 10 {

			c.PointReached()
			c.energy = c.energy - 0.5
			if c.energy < 0 {
				c.energy = 0
			}
		}
	}
}

func (c *Creature) CheckAndLevelUp() {
	if c.progress >= c.nextLevel && c.Size < 3 {
		c.Size += 1
		c.nextLevel *= 1.2
		c.progress = 0
		c.defaultMaxHunger *= 1.1

		ev := FishLevelUp{Fish: c}
		c.EventHub.Publish(ev)
	}
}

func randomBool() bool {
	return rand.Intn(2) == 0
}

func (c *Creature) RandomTarget() *geometry.Point {
	nextX := c.X + float32(rand.NormFloat64())*c.TankBoundaries.X2/4
	randomTargetX := max(c.TankBoundaries.X1+float32(c.SpriteWidth), nextX)

	if randomTargetX > c.TankBoundaries.X2-float32(c.SpriteWidth) {
		randomTargetX = c.TankBoundaries.X2 - float32(c.SpriteWidth)
	}

	//normally distributed y based on avg depth stat
	//standard dev = entire tank?
	// lowest point (highest y) - (a randomly, normally distributed number * std dev(50) +
	//then we subtract the mean depth of our species and the height since were dealing with a left corner of sprite)
	offsetFromBottom := float32(rand.NormFloat64())*20 + c.avgDepth + float32(c.SpriteHeight)
	println("fish offset =", offsetFromBottom)
	sample := c.TankBoundaries.Y2 - offsetFromBottom
	//clamp below highest possible point
	randomTargetY := max(c.TankBoundaries.Y1+float32(c.SpriteHeight), sample)
	randomTargetY = min(c.TankBoundaries.Y2-float32(c.SpriteHeight), randomTargetY)

	targetX := randomTargetX
	targetY := randomTargetY

	fmt.Printf("random point generated of x: %f y: %f\n", targetX, targetY)

	newPoint := geometry.Point{X: targetX, Y: targetY, PType: geometry.Structure}
	return &newPoint
}

func (c *Creature) AddRandomMovement() {

	// Calculate the desired direction
	desiredDx := c.PointQueue[0].X - c.X
	desiredDy := c.PointQueue[0].Y - c.Y

	// Normalize it
	length := float32(math.Hypot(float64(desiredDx), float64(desiredDy)))

	if length > 0 {
		desiredDx /= length
		desiredDy /= length
	}

	// Scale by desired speed
	desiredDx *= c.speed
	desiredDy *= c.speed

	// Smooth steering: blend current velocity toward desired
	steeringFactor := float32(0.05) // tweak for responsiveness

	c.Dx += (desiredDx - c.Dx) * steeringFactor
	c.Dy += (desiredDy - c.Dy) * steeringFactor

	c.ChangeAnimationSpeed(10 - (desiredDx * 3))
	if c.speed < 0.2 {
		c.ChangeAnimationSpeed(30)
	}

}

func (c *Creature) PointReached() {
	ev := CreatureReachedPoint{
		Point:    c.PointQueue[0],
		Creature: c,
	}

	c.EventHub.Publish(ev)
}

func (c *Creature) TranSlateFishShaderOpts() *ebiten.DrawRectShaderOptions {

	opts := &ebiten.DrawRectShaderOptions{}

	if c.Flip {
		opts.GeoM.Scale(-1, 1) // flip horizontally
		opts.GeoM.Translate(float64(c.SpriteWidth), 0)
	}

	if c.Dy < -0.5 {
		if c.Flip {
			opts.GeoM.Rotate(-0.3)
		} else {
			opts.GeoM.Rotate(0.3)
		}
	}
	if c.Dy > 0.5 {
		if c.Flip {
			opts.GeoM.Rotate(0.3)
		} else {
			opts.GeoM.Rotate(-0.3)
		}
	} /*
		b := c.Img.Bounds()
		midpoint := float32(b.Dy() / 2)*/

	y := float64(c.Y - float32(c.SpriteHeight/2))
	x := float64(c.X)

	if c.Flip {
		x = x - float64(c.SpriteWidth)
	}
	opts.GeoM.Translate(x, y)

	return opts
}

func (c *Creature) TranSlateFishOpts() *ebiten.DrawImageOptions {

	flip := c.Dx > 0

	opts := &ebiten.DrawImageOptions{}

	if flip {
		FlipSprite(float64(c.SpriteSheet.SpriteWidth/2), opts)
	}

	if c.Dy < -0.5 {
		if flip {
			opts.GeoM.Rotate(-0.3)
		} else {
			opts.GeoM.Rotate(0.3)
		}
	}

	if c.Dy > 0.5 {
		if flip {
			opts.GeoM.Rotate(0.3)
		} else {
			opts.GeoM.Rotate(-0.3)
		}
	}
	if flip {
		opts.GeoM.Translate(float64(c.X-float32(c.SpriteWidth)), float64(c.Y))
	} else {
		opts.GeoM.Translate(float64(c.X), float64(c.Y))
	}

	return opts
}

func FlipSprite(spriteWidth float64, opts *ebiten.DrawImageOptions) {
	opts.GeoM.Scale(-1, 1) // flip horizontally
	opts.GeoM.Translate(spriteWidth, 0)
}

func (c *Creature) AddTargetPointToQueue(point *geometry.Point) {
	if point.X > c.X {

	}

	c.PointQueue = append(c.PointQueue, point)
}

func (c *Creature) publishStats(sendTo string) {
	ev := SendData{}
	ev.DataFor = sendTo

	var targetPoint string
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

	switch c.HealthState {
	case Healthy:
		healthState = "healthy"
	case Stressed:
		healthState = "stressed"
	case Sick:
		healthState = "sick"
	}

	if len(c.PointQueue) > 0 {
		targetPoint = fmt.Sprintf("Target Point: %d, %d", int(c.PointQueue[0].X), int(c.PointQueue[0].Y))
	}

	nameString := fmt.Sprintf("Name: %s\n", c.name)
	hungerString := fmt.Sprintf("Hunger : %d/%d\n", int(c.Hunger), int(c.maxHunger))
	energyString := fmt.Sprintf("Energy : %d/%d\n", int(c.energy), int(c.maxEnergy))
	SizeString := fmt.Sprintf("Size : %d\n", c.Size)
	experienceString := fmt.Sprintf("Growth : %d/%d\n", int(c.progress), int(c.nextLevel))
	stateString := fmt.Sprintf("State: %s\n", state)
	healthStateString := fmt.Sprintf("Health State: %s\n", healthState)
	speedString := fmt.Sprintf("Speed: %d/%d\n", int(c.speed), int(c.maxSpeed))
	stressString := fmt.Sprintf("Stress: %d\n", int(c.Stress))

	ev.Data = nameString + stateString + SizeString + hungerString +
		energyString + experienceString + speedString + stressString + healthStateString

	if targetPoint != "" {
		ev.Data += targetPoint
	}

	c.EventHub.Publish(ev)
}

func GameFishToSaveFish(creature *Creature) SavedFish {

	var s SavedFish

	s.Name = creature.name
	s.Size = creature.Size
	s.Progress = creature.progress
	s.NextLevel = creature.nextLevel
	s.FishType = string(creature.FishType)

	return s
}

type FishStats struct {
	name             string
	Stress           float32
	Happiness        float32
	Hunger           float32
	maxHunger        float32
	defaultMaxHunger float32
	maxEnergy        float32
	energy           float32
	maxSpeed         float32
	avgSpeed         float32
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
		println("loading molly Fish")
		fs, err := GenMollyFishStats()
		fs.name = name
		if err != nil {
			return fs, err
		}
		return fs, nil
	case Fish:
		println("loading gold Fish")
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
	fs.avgSpeed = 1.0
	fs.speed = rand.Float32()*fs.maxSpeed + 0.3
	fs.FishType = MollyFish
	fs.maxEnergy = 25
	fs.energy = fs.maxEnergy / 2
	fs.Hunger = 5
	fs.maxHunger = 8*rand.Float32() + 4
	fs.defaultMaxHunger = fs.maxHunger
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
	fs.avgSpeed = 1.2
	fs.maxSpeed = rand.Float32() + 0.5
	fs.speed = rand.Float32()*fs.maxSpeed + 0.3
	fs.FishType = Fish
	fs.maxEnergy = 25
	fs.energy = fs.maxEnergy / 2
	fs.Hunger = 4
	fs.progress = 0
	fs.nextLevel = 10
	fs.idealPH = 6.5
	fs.maxHunger = 10*rand.Float32() + 4
	fs.defaultMaxHunger = fs.maxHunger
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
	fs.avgSpeed = 2.0
	fs.maxSpeed = rand.Float32() + 0.5
	fs.speed = rand.Float32()*fs.maxSpeed + 0.3
	fs.FishType = Guppy
	fs.maxEnergy = 25
	fs.energy = fs.maxEnergy / 2
	fs.Hunger = 4
	fs.progress = 0
	fs.nextLevel = 10
	fs.idealTemperature = 80
	fs.idealPH = 7.5
	fs.maxHunger = 10*rand.Float32() + 4
	fs.defaultMaxHunger = fs.maxHunger
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
	fs.avgSpeed = 1.5
	fs.maxSpeed = rand.Float32() + 0.5
	fs.speed = rand.Float32()*fs.maxSpeed + 0.3
	fs.FishType = Guppy
	fs.maxEnergy = 25
	fs.energy = fs.maxEnergy / 2
	fs.Hunger = 4
	fs.progress = 0
	fs.nextLevel = 10
	fs.idealTemperature = 77
	fs.idealPH = 7.0
	fs.maxHunger = 10*rand.Float32() + 4
	fs.defaultMaxHunger = fs.maxHunger
	persRoll := rand.Intn(10)

	if persRoll < 8 {
		fs.Personality = social
	} else {
		fs.Personality = shy
	}

	return fs, nil
}
