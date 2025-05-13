package gameEntities

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

func LoadOutlineShader() *ebiten.Shader {
	ols := []byte(shaders.OutlineShader)
	s, err := ebiten.NewShader(ols)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func LoadSolidColorShader() *ebiten.Shader {
	sls := []byte(shaders.SolidColor)
	s, err := ebiten.NewShader(sls)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

type Sprite struct {
	Img          *ebiten.Image
	X, Y         float32
	Dy, Dx       float32
	shader       *ebiten.Shader
	shaderParams map[string]any
}

func (s Sprite) SpriteHovered() bool {
	x, y := ebiten.CursorPosition()
	point := image.Point{x, y}
	rect := s.Img.Bounds()

	if rect.Max.X < 50 {
		rect.Max.X += 50
	}

	if rect.Max.Y < 50 {
		rect.Max.Y += 50
	}

	rect.Min.X += int(s.X)
	rect.Min.Y += int(s.Y)
	rect.Max.X += int(s.X)
	rect.Max.Y += int(s.Y)
	return point.In(rect)
}

type AnimatedSprite struct {
	*Sprite
	*animations.Animation
	*spritesheet.SpriteSheet
	frameImg *ebiten.Image
	effect   *ebiten.Image
}

func (s *Sprite) Coord() (x, y float32) {
	return s.X, s.Y
}

func (s *Sprite) LoadShader(shader *ebiten.Shader) {
	println("loading shader")
	s.shader = shader
}

func (s *Sprite) UnLoadShader() {
	s.shader = nil
}

func (as *AnimatedSprite) Update() {
	as.Animation.Update()
	as.UpdateSpriteFrameImg()
}

func (as *AnimatedSprite) UpdateSpriteFrameImg() {
	frame := as.Frame()
	frameRect := as.SpriteSheet.Rect(frame)
	img := as.Img.SubImage(frameRect).(*ebiten.Image)
	as.frameImg = img
}

func (as *AnimatedSprite) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions, shaderOpts *ebiten.DrawRectShaderOptions) {
	frame := as.Frame()
	frameRect := as.SpriteSheet.Rect(frame)
	img := as.Img.SubImage(frameRect).(*ebiten.Image)
	if as.shader != nil {
		shaderOpts.Images[0] = img
		shaderOpts.Uniforms = as.shaderParams
		b := img.Bounds()
		screen.DrawRectShader(b.Dx(), b.Dy(), as.shader, shaderOpts)
		return
	}
	screen.DrawImage(img, opts)
}

func NewAnimatedSprite() *AnimatedSprite {
	as := AnimatedSprite{
		&Sprite{},
		&animations.Animation{},
		&spritesheet.SpriteSheet{},
		&ebiten.Image{},
		nil,
	}
	return &as
}

func (as *AnimatedSprite) ChangeAnimationSpeed(newSpeed float32) {
	as.Animation.SpeedInTPS = newSpeed
}

func (as *AnimatedSprite) TriggerEffect(image *ebiten.Image) {
	as.effect = image
}
