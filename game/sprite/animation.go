package sprite

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Animation struct {
	SpriteSheet
	Img          *ebiten.Image
	NormalImg    *ebiten.Image
	FirstF       int
	LastF        int
	Step         int
	SpeedInTPS   float32
	FrameCounter float32
	frame        int
	OffSetX      float32
	OffSetY      float32
}

func (a *Animation) Update() {
	//this code iteration assumes each animation loops
	a.FrameCounter -= 1.0 // no need to worry about time as ebiten has a locked frame rate
	if a.FrameCounter < 0 {
		a.FrameCounter = a.SpeedInTPS
		a.frame += a.Step
		if a.frame > a.LastF {
			a.frame = a.FirstF
		}
	}
}

func (a *Animation) Frame() int {
	return a.frame
}

func (a *Animation) Reset() {
	a.frame = a.FirstF
}

func (a *Animation) Stop() int {
	return a.LastF
}

type SpriteSheet struct {
	WidthInTiles  int
	HeightInTiles int
	SpriteWidth   int
	SpriteHeight  int
}

func (s SpriteSheet) Rect(index int) image.Rectangle {
	x := index % s.WidthInTiles * s.SpriteWidth
	y := index / s.WidthInTiles * s.SpriteHeight
	lowX := x + s.SpriteWidth
	lowY := y + s.SpriteHeight
	return image.Rect(x, y, lowX, lowY)
}

func NewSpriteSheet(widthTiles, heightTiles, spriteWidth, spriteHeight int) SpriteSheet {
	return SpriteSheet{
		widthTiles, heightTiles, spriteWidth, spriteHeight,
	}
}

func NewAnimation(ss SpriteSheet, firstF int, lastF int, step int, speedinTPS float32) *Animation {
	return &Animation{
		SpriteSheet: ss,
		Step:        step,
		SpeedInTPS:  speedinTPS,
		FirstF:      firstF,
		LastF:       lastF,
	}
}
