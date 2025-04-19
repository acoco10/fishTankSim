package gameEntities

import (
	"fishTankWebGame/game/helperFunc"
	"fishTankWebGame/game/sprite"
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
	"math/rand"
	"sort"
)

type FishState uint8

const (
	Swimming FishState = iota
	Eating
	Resting
)

type Creature struct {
	Size           int
	progress       int
	nextLevel      int
	appetite       int
	pointQueue     []*Point
	eventHub       *EventHub
	tankBoundaries image.Rectangle
	timers         map[FishState]*Timer
	state          FishState
	maxSpeed       float32
	*sprite.AnimatedSprite
}

func NewFish(hub *EventHub, tankSize image.Rectangle) *Creature {

	timers := make(map[FishState]*Timer)
	randDuration := rand.Float64() * 50
	timers[Swimming] = NewTimer(randDuration)
	timers[Eating] = NewTimer(0.5)
	timers[Resting] = NewTimer(30)
	maxSpeed := rand.Float32()
	c := Creature{
		1,
		0,
		10,
		3,
		[]*Point{},
		hub,
		tankSize,
		timers,
		Swimming,
		maxSpeed,
		sprite.NewAnimatedSprite(),
	}

	firstPoint := c.RandomTarget()
	c.AddTargetPointToQueue(firstPoint)

	img := helperFunc.LoadImageAssetAsEbitenImage("fishSpriteSheet")
	c.Img = img
	c.Animation = animations.NewAnimation(0, 3, 1, 4)
	c.SpriteSheet = spritesheet.NewSpritesheet(4, 2, 19, 7)

	//water Level does not reach the top of  the tank
	// should automate this somehow
	maxX := c.tankBoundaries.Max.X - c.SpriteWidth - 2
	minX := c.tankBoundaries.Min.X + 2
	maxY := c.tankBoundaries.Max.Y - c.SpriteHeight - 10
	minY := c.tankBoundaries.Min.Y + c.SpriteHeight + 50

	tankSize.Max.X = maxX
	tankSize.Min.X += minX
	tankSize.Max.Y -= maxY
	tankSize.Min.Y += minY

	c.X = rand.Float32()*200 + float32(tankSize.Min.X)
	c.Y = rand.Float32()*100 + float32(tankSize.Min.Y)
	c.eventHub.Subscribe(MouseButtonPressed{}, func(e Event) {
		ev := e.(MouseButtonPressed)
		if c.pointQueue[0].PType != Food {
			c.pointQueue = []*Point{}
		}
		c.AddTargetPointToQueue(ev.Point)
		c.sortPoints()
	})
	c.eventHub.Subscribe(CreatureReachedPoint{}, func(e Event) {
		ev := e.(CreatureReachedPoint)
		for i, point := range c.pointQueue {
			if ev.Point == point {
				c.pointQueue[i] = c.pointQueue[len(c.pointQueue)-1]
				c.pointQueue = c.pointQueue[:len(c.pointQueue)-1]
				if ev.Creature == &c && ev.Point.PType == Food {
					c.state = Eating
				}
				if ev.Creature == &c && ev.Point.PType != Food {
					c.RandomTarget()
				}
			}
		}
		c.sortPoints()
	})

	return &c
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

func (c *Creature) Direction() float64 {
	dX := float64(c.Dx)
	dY := float64(c.Dy)
	angle := math.Atan2(dX, dY)

	tiltThreshold := 0.3

	switch {
	case angle < -tiltThreshold:
		angle = -0.3
	case angle > tiltThreshold:
		angle = 0.3
	default:
		angle = 0
	}

	return angle
}

func (c *Creature) sortPoints() {

	sort.Slice(c.pointQueue, func(i, j int) bool {
		xI, yI := c.pointQueue[i].Coord()
		xJ, yJ := c.pointQueue[j].Coord()

		distI := DistanceFunc(xI, c.X, yI, c.Y)
		distJ := DistanceFunc(xJ, c.X, yJ, c.Y)

		return distI < distJ
	})

}

func (c *Creature) Move() {
	c.AddRandomMovement()
	c.X += c.Dx
	c.Y += c.Dy

	maxX := float32(c.tankBoundaries.Max.X)
	minX := float32(c.tankBoundaries.Min.X)
	maxY := float32(c.tankBoundaries.Max.Y)
	minY := float32(c.tankBoundaries.Min.Y)

	if c.X > maxX {
		c.X = maxX
	}
	if c.X < minX {
		c.X = minX
	}
	if c.Y > maxY {
		c.Y = maxY
	}
	if c.Y < minY {
		c.Y = minY
	}

}

func (c *Creature) Update() {
	switch c.state {
	case Swimming:
		c.Move()
		tState := c.timers[Swimming].Update()
		if tState == Done {
			c.timers[Swimming].Duration = rand.Intn(40)
		}
	case Resting:
		tState := c.timers[Resting].Update()
		if tState == Done {
			c.state = Swimming
		}
	case Eating:
		tState := c.timers[Eating].Update()
		if tState == Done {
			c.state = Swimming
		}
	}
	if c.progress >= c.nextLevel {
		newImg := helperFunc.LoadImageAssetAsEbitenImage("goldFishGrowth1")
		c.Img = newImg
	}
	c.AnimatedSprite.Update()
}

func randomBool() bool {
	return rand.Intn(2) == 0
}

func (c *Creature) RandomTarget() *Point {
	randomTargetX := rand.Intn(c.tankBoundaries.Max.X-c.tankBoundaries.Min.X) + c.tankBoundaries.Min.X
	randomTargetY := rand.Intn(c.tankBoundaries.Max.Y-c.tankBoundaries.Min.Y) + c.tankBoundaries.Min.Y

	targetX := float32(randomTargetX)
	targetY := float32(randomTargetY)

	newPoint := Point{targetX, targetY, Structure}
	fmt.Printf("Random target point: (%f, %f)\n", targetX, targetY)
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
	desiredDx *= c.maxSpeed
	desiredDy *= c.maxSpeed

	// Smooth steering: blend current velocity toward desired
	steeringFactor := float32(0.05) // tweak for responsiveness

	c.Dx += (desiredDx - c.Dx) * steeringFactor
	c.Dy += (desiredDy - c.Dy) * steeringFactor

	c.ChangeAnimationSpeed(10 - (desiredDx * 3))

	if len(c.pointQueue) > 0 {
		c.CheckDistFromPoint(c.pointQueue[0].Coord())
	}

}

func (c *Creature) PointReached() {

	ev := CreatureReachedPoint{
		Point:    c.pointQueue[0],
		Creature: c,
	}
	c.eventHub.Publish(ev)
}

func (c *Creature) Draw(screen *ebiten.Image) {
	flip := c.Dx > 0

	opts := ebiten.DrawImageOptions{}

	if flip {
		opts.GeoM.Scale(-1, 1) // flip horizontally
		opts.GeoM.Translate(float64(c.SpriteSheet.SpriteWidth), 0)
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

	c.DrawSprite(screen, &opts)

}

func (c *Creature) CheckDistFromPoint(x, y float32) {
	difX := x - c.X
	difY := y - c.Y

	if math.Abs(float64(difY)) < 5 && math.Abs(float64(difX)) < 5 {
		c.PointReached()
	}
}

func (c *Creature) AddTargetPointToQueue(point *Point) {

	c.pointQueue = append(c.pointQueue, point)
}
