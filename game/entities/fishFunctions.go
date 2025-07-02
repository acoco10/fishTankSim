package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"math/rand"
	"sort"
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
		c.CheckAndLevelUp()
	})

}

func (c *Creature) ownPointReached(ev CreatureReachedPoint) {

	if len(c.PointQueue) > 0 {
		c.PointQueue = c.PointQueue[1:]

	}

	switch ev.Point.PType {
	case geometry.Food:
		println("reached own point, eating food and going for next piece")
		c.progress += 1
		c.Hunger -= 1
		if c.Hunger < 0 {
			c.Hunger = 0
		}
		c.State = Eating
		if len(c.PointQueue) > 0 {
			if c.Hunger > 0 {
				c.goToFood()
			}
		}

	default:
		println("reached own point, setting random speed and random next target")

		if c.energy == 0 {
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

		xI, yI := c.PointQueue[i].Coord()
		xJ, yJ := c.PointQueue[j].Coord()

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
		x, y := c.X, c.Y
		c.X = x
		c.Y = y
		c.nextLevel *= 1.2
		c.progress = 0

		ev := FishLevelUp{Fish: c}
		c.EventHub.Publish(ev)
	}
}

func randomBool() bool {
	return rand.Intn(2) == 0
}

func (c *Creature) RandomTarget() *geometry.Point {
	randomTargetX := max(c.TankBoundaries.X1+float32(c.SpriteWidth), rand.Float32()*c.TankBoundaries.X2-float32(c.SpriteWidth))

	//normally distributed y based on avg depth stat
	//standard dev = entire tank?
	sample := float32(rand.NormFloat64())*50 + c.TankBoundaries.Y2 - c.avgDepth - 100
	randomTargetY := max(c.TankBoundaries.Y1+float32(c.SpriteHeight), sample)

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
	if c.State == Resting {
		state = "resting"
	}
	if c.State == Swimming {
		state = "swimming"
	}
	if c.State == Eating {
		state = "eating"
	}

	if len(c.PointQueue) > 0 {
		targetPoint = fmt.Sprintf("Target Point: %f, %f", c.PointQueue[0].X, c.PointQueue[0].Y)
	}

	nameString := fmt.Sprintf("Name: %s\n", c.name)
	hungerString := fmt.Sprintf("Hunger : %f/%f\n", c.Hunger, c.maxHunger)
	energyString := fmt.Sprintf("Energy : %f/%f\n", c.energy, c.maxEnergy)
	SizeString := fmt.Sprintf("Size : %d\n", c.Size)
	experienceString := fmt.Sprintf("Growth : %f/%f\n", c.progress, c.nextLevel)
	stateString := fmt.Sprintf("State: %s\n", state)
	speedString := fmt.Sprintf("Speed: %f/%f\n", c.speed, c.maxSpeed)

	ev.Data = nameString + stateString + SizeString + hungerString + energyString + experienceString + speedString

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
	name        string
	Hunger      float32
	maxHunger   float32
	maxEnergy   float32
	energy      float32
	maxSpeed    float32
	avgSpeed    float32
	avgDepth    float32
	speed       float32
	Size        int
	progress    float32
	nextLevel   float32
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
	}
	return nil, nil
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
	fs.avgDepth = 100
	fs.progress = 0
	fs.nextLevel = 10

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
	fs.avgDepth = 0.0
	fs.avgSpeed = 2.0
	fs.maxSpeed = rand.Float32() + 0.5
	fs.speed = rand.Float32()*fs.maxSpeed + 0.3
	fs.FishType = Fish
	fs.maxEnergy = 25
	fs.energy = fs.maxEnergy / 2
	fs.Hunger = 4
	fs.progress = 0
	fs.nextLevel = 10
	fs.maxHunger = 10*rand.Float32() + 4
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
	fs.maxHunger = 10*rand.Float32() + 4
	persRoll := rand.Intn(10)

	if persRoll < 8 {
		fs.Personality = social
	} else {
		fs.Personality = shy
	}

	return fs, nil
}
