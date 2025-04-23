package gameEntities

import (
	"fishTankWebGame/shaders"
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

func loadOutlineShader() *ebiten.Shader {
	ols := []byte(shaders.OutlineShader)
	s, err := ebiten.NewShader(ols)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

type Sprite struct {
	Img    *ebiten.Image
	X, Y   float32
	Dy, Dx float32
}

func (s Sprite) SpriteHovered() bool {
	x, y := ebiten.CursorPosition()
	point := image.Point{x, y}
	rect := s.Img.Bounds()
	rect.Min.X += int(s.X)
	rect.Min.Y += int(s.Y)
	rect.Max.X += int(s.X)
	rect.Max.Y += int(s.Y)
	return point.In(rect)
}

type UiSprite struct {
	*Sprite
}

func AddYellowOutlineShader(spriteImg *ebiten.Image, sprite Sprite, screen *ebiten.Image) {
	var options ebiten.DrawRectShaderOptions

	width, height := ebiten.WindowSize()

	options.Images[0] = spriteImg // the sprite to outline
	options.Uniforms = map[string]interface{}{
		"Resolution": []float32{float32(width), float32(height)},
	}

	options.GeoM.Translate(float64(sprite.X), float64(sprite.Y))
	s := loadOutlineShader()
	DrawShader(sprite, spriteImg, s, screen)
}

func ApplyOutlineShaderToAnimation(sprite AnimatedSprite, screen *ebiten.Image) {
	frame := sprite.Frame()
	frameRect := sprite.SpriteSheet.Rect(frame)
	eImg := ebiten.NewImageFromImage(frameRect)
	AddYellowOutlineShader(eImg, *sprite.Sprite, screen)
}
func (us *UiSprite) Draw(screen *ebiten.Image) {

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(us.X), float64(us.Y))
	screen.DrawImage(us.Img, &opts)
	
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

func (as *AnimatedSprite) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if as.SpriteHovered() {
		ApplyOutlineShaderToAnimation(*as, screen)
	} else {
		frame := as.Frame()
		frameRect := as.SpriteSheet.Rect(frame)
		screen.DrawImage(as.Img.SubImage(frameRect).(*ebiten.Image), opts)
	}
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
