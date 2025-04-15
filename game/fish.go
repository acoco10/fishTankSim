package game

import (
	"fishTankWebGame/game/events"
	"fishTankWebGame/game/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
	"math/rand"
	"sort"
)

type Creature struct {
	*sprite.Sprite
	counter        int
	xDirection     bool
	yDirection     bool
	Size           int
	progress       int
	nextLevel      int
	appetite       int
	pointQueue     []image.Point
	eventHub       *events.EventHub
	tankBoundaries image.Rectangle
}

func NewFish(hub *events.EventHub, tankSize image.Rectangle) *Creature {
	c := Creature{
		&sprite.Sprite{},
		0,
		true,
		true,
		1,
		0,
		10,
		3,
		[]image.Point{},
		hub,
		tankSize,
	}
	img := LoadImageAssetAsEbitenImage("fishSprite")
	c.Img = img
	c.X = rand.Float64() * 200
	c.Y = rand.Float64() * 100
	c.eventHub.Subscribe(events.MouseButtonPressed{}, func(e events.Event) {
		ev := e.(events.MouseButtonPressed)
		c.AddTargetPointToQueue(ev.Point)
		c.sortPoints()
	})
	c.eventHub.Subscribe(events.CreatureReachedPoint{}, func(e events.Event) {
		ev := e.(events.CreatureReachedPoint)
		for i, point := range c.pointQueue {
			if ev.Point == point {
				c.pointQueue[i] = c.pointQueue[len(c.pointQueue)-1]
				c.pointQueue = c.pointQueue[:len(c.pointQueue)-1]
			}
		}
		c.sortPoints()
	})

	return &c
}

func DistanceFunc(x, x2, y, y2 int) float64 {
	xDis := float64(x - x2)
	yDis := float64(y - y2)

	dis := math.Sqrt(math.Pow(xDis, 2) + math.Pow(yDis, 2))
	return dis
}

func (c *Creature) sortPoints() {
	sort.Slice(c.pointQueue, func(i, j int) bool {
		distI := DistanceFunc(c.pointQueue[i].X, int(c.X), c.pointQueue[i].Y, int(c.Y))
		distJ := DistanceFunc(c.pointQueue[j].X, int(c.X), c.pointQueue[j].Y, int(c.Y))

		return distI < distJ
	})
}

func (c *Creature) Move() {

	c.X += c.Dx
	c.Y += c.Dy

	if c.X > float64(c.tankBoundaries.Max.X) {
		c.X = float64(c.tankBoundaries.Max.X)
	}
	if c.X < float64(c.tankBoundaries.Min.X) {
		c.X = float64(c.tankBoundaries.Min.X)
	}

	if c.Y > float64(c.tankBoundaries.Max.Y) {
		c.Y = float64(c.tankBoundaries.Max.Y)
	}

	if c.Y < float64(c.tankBoundaries.Min.Y) {
		c.Y = float64(c.tankBoundaries.Min.Y)
	}

	c.Dx = 0
	c.Dy = 0

}

func (c *Creature) Update() {
	c.AddRandomMovement()
	c.Move()
	c.counter++
	if c.counter == 121 {
		c.counter = 0
	}
	if c.progress >= c.nextLevel {
		newImg := LoadImageAssetAsEbitenImage("goldFishGrowth1")
		c.Img = newImg
	}
}

func randomBool() bool {
	return rand.Intn(2) == 0
}

func (c *Creature) AddRandomMovement() {
	dX := rand.Float64()
	dY := rand.Float64()

	if c.counter == 120 && len(c.pointQueue) == 0 {
		c.xDirection = randomBool()
		c.yDirection = randomBool()
	}

	if !c.xDirection {
		dX = -dX
	}

	if !c.yDirection {
		dY = -dY
	}

	if len(c.pointQueue) > 0 {
		c.TravelToPoint(c.pointQueue[0])
	}

	c.Dx = dX
	c.Dy = dY
}

func (c *Creature) PointReached() {
	ev := events.CreatureReachedPoint{
		Point: c.pointQueue[0],
	}
	c.progress += 1
	c.eventHub.Publish(ev)
}

func (c *Creature) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	if c.xDirection {
		opts.GeoM.Scale(-1, 1)
	}
	opts.GeoM.Translate(c.X, c.Y)
	screen.DrawImage(c.Img, &opts)
}

func (c *Creature) TravelToPoint(point image.Point) {
	difX := float64(point.X) - c.X
	difY := float64(point.Y) - c.Y

	if math.Abs(difX) > 5 {
		if difX > 0 {
			c.xDirection = true
		}
		if difX < 0 {
			c.xDirection = false
		}
	}
	if math.Abs(difY) > 5 {
		if difY > 0 {
			c.yDirection = true
		}
		if difY < 0 {
			c.yDirection = false
		}
	}
	if math.Abs(difY) < 5 && math.Abs(difX) < 5 {
		c.PointReached()
	}
}

func (c *Creature) AddTargetPointToQueue(point image.Point) {
	c.pointQueue = append(c.pointQueue, point)
}
