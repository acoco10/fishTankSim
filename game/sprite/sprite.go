package sprite

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img    *ebiten.Image
	X, Y   float32
	Dy, Dx float32
}

type AnimatedSprite struct {
	*Sprite
	*animations.Animation
	*spritesheet.SpriteSheet
}

func (s Sprite) Coord() (x, y float32) {
	return s.X, s.Y
}

func (as *AnimatedSprite) Update() {
	as.Animation.Update()
}

func (as *AnimatedSprite) DrawSprite(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	frame := as.Frame()
	frameRect := as.SpriteSheet.Rect(frame)
	screen.DrawImage(as.Img.SubImage(frameRect).(*ebiten.Image), opts)
}

func NewAnimatedSprite() *AnimatedSprite {
	as := AnimatedSprite{
		&Sprite{},
		&animations.Animation{},
		&spritesheet.SpriteSheet{},
	}
	return &as
}

func (as *AnimatedSprite) ChangeAnimationSpeed(newSpeed float32) {
	as.Animation.SpeedInTPS = newSpeed
}
