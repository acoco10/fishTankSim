package game

import (
	"fishTankWebGame/game/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
)

type Creature struct {
	*sprite.Sprite
	ActionQueue []func()
	counter     int
	xDirection  bool
	yDirection  bool
	Size        int
}

func NewFish() *Creature {
	f := Creature{
		&sprite.Sprite{},
		[]func(){},
		0,
		true,
		true,
		1,
	}
	img := LoadImageAssetAsEbitenImage("fishSprite")
	f.Img = img
	f.X = 100
	f.Y = 100
	return &f
}

func (c *Creature) Move() {

	c.X += c.Dx
	c.Y += c.Dy

	if c.X > 480 {
		c.X = 480
	}
	if c.X < 60 {
		c.X = 60
	}

	if c.Y > 280 {
		c.Y = 280
	}

	if c.Y < 0 {
		c.Y = 0
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
}

func randomBool() bool {
	return rand.Intn(2) == 0
}

func (c *Creature) AddRandomMovement() {
	dX := rand.Float64()
	dY := rand.Float64()

	if c.counter == 120 {
		c.xDirection = randomBool()
		c.yDirection = randomBool()
	}

	if !c.xDirection {
		dX = -dX
	}

	if !c.yDirection {
		dY = -dY
	}

	c.Dx = dX
	c.Dy = dY
}

func (c *Creature) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	if c.xDirection {
		opts.GeoM.Scale(-1, 1)
	}
	opts.GeoM.Translate(c.X, c.Y)
	screen.DrawImage(c.Img, &opts)
}
