package gameEntities

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"math/rand"
	"sort"
)

func CreatureEventSubscriptions(c *Creature) {
	c.eventHub.Subscribe(PointGenerated{}, func(e Event) {
		ev := e.(PointGenerated)
		if ev.Point.PType == Food {
			if c.hunger > 0 {
				c.AddTargetPointToQueue(ev.Point)
				c.sortPoints()
			}
		}
	})

	c.eventHub.Subscribe(CreatureReachedPoint{}, func(e Event) {
		ev := e.(CreatureReachedPoint)
		if ev.Creature == c {
			c.ownPointReached(ev)
		} else {
			c.otherFishPoint(ev)
		}
	})
}

func (c *Creature) ownPointReached(ev CreatureReachedPoint) {

	if len(c.pointQueue) > 0 {
		c.pointQueue = c.pointQueue[1:]
	}

	switch ev.Point.PType {
	case Food:
		println("reached own point, eating food and going for next piece")
		c.progress += 1
		c.hunger -= 1
		if c.hunger < 0 {
			c.hunger = 0
		}
		c.state = Eating
		if len(c.pointQueue) > 0 {
			if c.hunger > 0 {
				c.goToFood()
			}
		}

	default:
		println("reached own point, setting random speed and random next target")

		if c.energy == 0 {
			c.state = Resting
		}
		c.CalcSpeed()
		//other creature behaviour
	}

}

func (c *Creature) otherFishPoint(ev CreatureReachedPoint) {
	pointThere := false
	var pointIndex int
	for i, point := range c.pointQueue {
		if ev.Point == point {
			pointThere = true
			pointIndex = i
		}
	}

	endIndex := len(c.pointQueue) - 1

	if pointThere {
		switch ev.Point.PType {
		case Food:
			c.pointQueue[pointIndex] = c.pointQueue[endIndex]
			c.pointQueue = c.pointQueue[0:endIndex]

		default:
			//placeholder, no modification when another Fish reaches their point that's not food as of now
		}
	}
}

func (c *Creature) goToFood() {
	c.sortPoints()
	c.CalcSpeed()
}

func (c *Creature) CalcSpeed() {
	if len(c.pointQueue) < 0 {
		if c.pointQueue[0].PType == Food {
			c.speed = c.maxSpeed
		} else {
			c.speed = rand.Float32()*c.maxSpeed + 0.5*(c.energy/c.maxEnergy)
		}
	}

	fmt.Printf("random speed generated = %f", c.speed)

}

func DistanceFunc(x, x2, y, y2 float32) float64 {
	xDis := float64(x - x2)
	yDis := float64(y - y2)
	dis := math.Sqrt(math.Pow(xDis, 2) + math.Pow(yDis, 2))
	return dis
}

func (c *Creature) Type() InterestPoint {
	return OtherCreature
}

func (c *Creature) sortPoints() {

	sort.Slice(c.pointQueue, func(i, j int) bool {
		xI, yI := c.pointQueue[i].Coord()
		xJ, yJ := c.pointQueue[j].Coord()

		distI := math.Hypot(float64(c.X-xI), float64(c.Y-yI))
		distJ := math.Hypot(float64(c.X-xJ), float64(c.Y-yJ))

		pTypeI := c.pointQueue[i].PType
		pTypeJ := c.pointQueue[j].PType

		if pTypeI == Food && pTypeJ != Food {
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

	c.NewPoint()

}

func (c *Creature) EnforceBoundaries() {
	c.X = max(c.tankBoundaries.X1, c.X)
	c.Y = max(c.tankBoundaries.Y1, c.Y)

	c.X = min(c.tankBoundaries.X2, c.X)
	c.Y = min(c.tankBoundaries.Y2, c.Y)
}

func (c *Creature) NewPoint() {
	if len(c.pointQueue) > 0 {
		tgtPoint := c.pointQueue[0]
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

func (c *Creature) LevelUp() {
	if c.progress >= c.nextLevel {
		c.Size += 1
		x, y := c.X, c.Y
		c.AnimatedSprite = LoadFishSprite(c.fishType, c.Size)
		c.X = x
		c.Y = y
		c.nextLevel *= 1.2
		c.progress = 0
	}
}

func randomBool() bool {
	return rand.Intn(2) == 0
}

func (c *Creature) RandomTarget() *Point {
	randomTargetX := max(c.tankBoundaries.X1, rand.Float32()*c.tankBoundaries.X2)
	randomTargetY := max(c.tankBoundaries.Y1, rand.Float32()*c.tankBoundaries.Y2)

	targetX := randomTargetX
	targetY := randomTargetY

	fmt.Printf("random point generated of x: %f y: %f\n", targetX, targetY)

	newPoint := Point{targetX, targetY, Structure}
	return &newPoint
}

func (c *Creature) AddRandomMovement() {
	if len(c.pointQueue) == 0 {
		c.AddTargetPointToQueue(c.RandomTarget())
	}

	// Calculate the desired direction
	desiredDx := c.pointQueue[0].X - c.X
	desiredDy := c.pointQueue[0].Y - c.Y

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
	println("publishing creature reached point")
	ev := CreatureReachedPoint{
		Point:    c.pointQueue[0],
		Creature: c,
	}
	c.eventHub.Publish(ev)
}

func (c *Creature) TranSlateFishShaderOpts() *ebiten.DrawRectShaderOptions {
	flip := c.Dx > 0

	opts := &ebiten.DrawRectShaderOptions{}

	if flip {
		opts.GeoM.Scale(-1, 1) // flip horizontally
		opts.GeoM.Translate(float64(c.SpriteWidth/2), 0)
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
	opts.GeoM.Translate(float64(c.X), float64(c.Y))

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
	opts.GeoM.Translate(float64(c.X), float64(c.Y))

	return opts
}

func FlipSprite(spriteWidth float64, opts *ebiten.DrawImageOptions) {
	opts.GeoM.Scale(-1, 1) // flip horizontally
	opts.GeoM.Translate(spriteWidth/2, 0)
}

func (c *Creature) AddTargetPointToQueue(point *Point) {

	c.pointQueue = append(c.pointQueue, point)
}

func (c *Creature) publishStats(sendTo string) {
	ev := SendData{}
	ev.DataFor = sendTo

	var state string
	if c.state == Resting {
		state = "resting"
	}
	if c.state == Swimming {
		state = "swimming"
	}
	if c.state == Eating {
		state = "eating"
	}

	nameString := fmt.Sprintf("Name: %s\n", c.name)
	hungerString := fmt.Sprintf("Hunger : %f/%f\n", c.hunger, c.maxHunger)
	energyString := fmt.Sprintf("Energy : %f/%f\n", c.energy, c.maxEnergy)
	SizeString := fmt.Sprintf("Size : %d\n", c.Size)
	experienceString := fmt.Sprintf("Growth : %f/%f\n", c.progress, c.nextLevel)
	stateString := fmt.Sprintf("State: %s\n", state)
	speedString := fmt.Sprintf("Speed: %f/%f\n", c.speed, c.maxSpeed)

	ev.Data = nameString + stateString + SizeString + hungerString + energyString + experienceString + speedString

	c.eventHub.Publish(ev)
}

func GameFishToSaveFish(creature *Creature) SavedFish {

	var s SavedFish

	s.Size = creature.Size
	s.Progress = creature.progress
	s.NextLevel = creature.nextLevel
	s.FishType = string(creature.fishType)

	return s
}

type FishStats struct {
	name        string
	hunger      float32
	maxHunger   float32
	maxEnergy   float32
	energy      float32
	maxSpeed    float32
	speed       float32
	Size        int
	progress    float32
	nextLevel   float32
	personality FishPersonality
	fishType    FishList
}

func GenFishStats(fType FishList) (*FishStats, error) {
	switch fType {
	case MollyFish:
		println("loading molly Fish")
		fs, err := GenMollyFishStats()
		if err != nil {
			return fs, err
		}
		return fs, nil
	case Fish:
		println("loading gold Fish")
		fs, err := GenGoldFishStats()
		if err != nil {
			return fs, err
		}
		return fs, nil
	default:
		fs, err := GenGoldFishStats()
		if err != nil {
			return fs, err
		}
		return fs, nil
	}
}

func GenMollyFishStats() (*FishStats, error) {
	fs := &FishStats{}

	fs.Size = 1
	fs.maxSpeed = rand.Float32() + 0.6
	fs.speed = rand.Float32()*fs.maxSpeed + 0.1
	fs.fishType = MollyFish
	fs.maxEnergy = 10*rand.Float32() + 5
	fs.energy = fs.maxEnergy / 2
	fs.hunger = 2
	fs.maxHunger = 8*rand.Float32() + 4

	fs.progress = 0
	fs.nextLevel = 10

	persRoll := rand.Intn(10)

	if persRoll < 4 {
		fs.personality = social
	} else {
		fs.personality = shy
	}

	return fs, nil
}

func GenGoldFishStats() (*FishStats, error) {
	fs := &FishStats{}
	fs.Size = 1
	fs.maxSpeed = rand.Float32() + 0.2
	fs.speed = rand.Float32()*fs.maxSpeed + 0.1
	fs.fishType = Fish
	fs.maxEnergy = 10 * rand.Float32()
	fs.energy = fs.maxEnergy / 2
	fs.hunger = 2
	fs.progress = 0
	fs.nextLevel = 10
	fs.maxHunger = 10*rand.Float32() + 4
	persRoll := rand.Intn(10)

	if persRoll < 8 {
		fs.personality = social
	} else {
		fs.personality = shy
	}

	return fs, nil
}
