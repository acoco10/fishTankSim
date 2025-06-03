package sprite

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Sprite struct {
	Img                *ebiten.Image
	Scale              float64
	X, Y               float32
	Dy, Dx             float32
	Shader             *ebiten.Shader
	ShaderParams       map[string]any
	CPUShaderParams    map[string]any
	UpdateShaderParams func(map[string]any) map[string]any
	UpdateBothParams   func(map[string]any, map[string]any) (map[string]any, map[string]any)
}

func (s *Sprite) UpdateShader() {
	if s.CPUShaderParams != nil {
		s.CPUShaderParams["origin"] = [2]float64{float64(s.X), float64(s.Y)}
	}
	if s.UpdateBothParams != nil {
		shaderParams, cpuParams := s.UpdateBothParams(s.ShaderParams, s.CPUShaderParams)
		s.ShaderParams = shaderParams
		s.CPUShaderParams = cpuParams
		return
	}

	if s.UpdateShaderParams != nil {
		s.ShaderParams = s.UpdateShaderParams(s.ShaderParams)
	}
}

func (s *Sprite) Update() {
	s.UpdateShader()
}

func (s *Sprite) Draw(screen *ebiten.Image) {

	if s.Shader != nil {
		shaderOpts := &ebiten.DrawRectShaderOptions{}
		shaderOpts.GeoM.Translate(float64(s.X), float64(s.Y))
		shaderOpts.Images[0] = s.Img
		shaderOpts.Uniforms = s.ShaderParams
		b := s.Img.Bounds()
		screen.DrawRectShader(b.Dx(), b.Dy(), s.Shader, shaderOpts)
		return
	}

	dOpts := &ebiten.DrawImageOptions{}
	dOpts.GeoM.Translate(float64(s.X), float64(s.Y))
	screen.DrawImage(s.Img, dOpts)

}

func (s *Sprite) SpriteHovered() bool {
	x, y := ebiten.CursorPosition()
	point := image.Point{x, y}
	rect := s.Img.Bounds()

	if rect.Max.X < 50 {
		rect.Max.X += 25
		rect.Min.X -= 25
	}

	if rect.Max.Y < 50 {
		rect.Max.Y += 25
		rect.Min.Y -= 25
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
	s.Shader = shader
}

func (s *Sprite) UnLoadShader() {
	s.Shader = nil
}

func (s *Sprite) CheckOverlap(sprite Sprite) bool {
	return s.Img.Bounds().Overlaps(sprite.Img.Bounds())
}

func (as *AnimatedSprite) Update() {
	as.UpdateShader()
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
	if as.Shader != nil {
		shaderOpts.Images[0] = img
		shaderOpts.Uniforms = as.ShaderParams
		b := img.Bounds()
		screen.DrawRectShader(b.Dx(), b.Dy(), as.Shader, shaderOpts)
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
