package sprite

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

type Sprite struct {
	Img                *ebiten.Image
	NormalMap          *ebiten.Image
	Scale              float64
	X, Y               float32
	Dy, Dx             float32
	Shader             *ebiten.Shader
	ShaderParams       map[string]any
	CPUShaderParams    map[string]any
	UpdateShaderParams func(map[string]any) map[string]any
	UpdateBothParams   func(map[string]any, map[string]any) (map[string]any, map[string]any)
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

func (s *Sprite) SpriteHovered() bool {
	x, y := ebiten.CursorPosition()
	point := image.Point{X: x, Y: y}
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
	frameImg   *ebiten.Image
	effect     *ebiten.Image
	drawOpts   *ebiten.DrawImageOptions
	shaderOpts *ebiten.DrawRectShaderOptions
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

	shaderOpts := &ebiten.DrawRectShaderOptions{}

	if as.Scale > 0 {
		shaderOpts.GeoM.Scale(as.Scale, as.Scale)
	}

	shaderOpts.GeoM.Translate(float64(as.X), float64(as.Y))

	as.shaderOpts = shaderOpts

	drawOpts := &ebiten.DrawImageOptions{}

	if as.Scale > 0 {
		drawOpts.GeoM.Scale(as.Scale, as.Scale)
	}

	drawOpts.GeoM.Translate(float64(as.X), float64(as.Y))

	as.drawOpts = drawOpts
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

func (as *AnimatedSprite) Draw(screen *ebiten.Image) {
	frame := as.Frame()
	frameRect := as.SpriteSheet.Rect(frame)
	img := as.Img.SubImage(frameRect).(*ebiten.Image)

	if as.NormalMap != nil {
		as.DrawNormal(screen, as.shaderOpts)
		return
	}

	if as.Shader != nil {
		as.shaderOpts.Images[0] = img
		as.shaderOpts.Uniforms = as.ShaderParams
		b := img.Bounds()
		screen.DrawRectShader(b.Dx(), b.Dy(), as.Shader, as.shaderOpts)
		return
	}
	if as.drawOpts == nil {
		as.drawOpts = &ebiten.DrawImageOptions{}
		as.drawOpts.GeoM.Translate(float64(as.X), float64(as.Y))
	}
	screen.DrawImage(img, as.drawOpts)
}

func (as *AnimatedSprite) DrawNormal(screen *ebiten.Image, shaderOpts *ebiten.DrawRectShaderOptions) {

	if as.shaderOpts == nil {
		as.shaderOpts = &ebiten.DrawRectShaderOptions{}
		as.shaderOpts.GeoM.Translate(float64(as.X), float64(as.Y))
	}

	if as.Shader == nil {
		shader := shaders.LoadNormalMapShader()
		as.Shader = shader
	}

	frame := as.Frame()

	frameRect := as.SpriteSheet.Rect(frame)

	diffuseImg := as.Img.SubImage(frameRect).(*ebiten.Image)
	if diffuseImg == nil {
		log.Fatal("normal map sub rect is disposed")
	}

	normalImg := as.NormalMap.SubImage(frameRect).(*ebiten.Image)
	if normalImg == nil {
		log.Fatal("normal map sub rect is disposed")
	}

	as.shaderOpts.Images[0] = diffuseImg
	as.shaderOpts.Images[1] = normalImg

	as.shaderOpts.Uniforms = as.ShaderParams

	b := diffuseImg.Bounds()
	screen.DrawRectShader(b.Dx(), b.Dy(), as.Shader, as.shaderOpts)
}

func (as *AnimatedSprite) UpdateOpts(options any) {

	opts, ok := options.(*ebiten.DrawImageOptions)
	if ok {
		as.drawOpts = opts
	}

	shaderOpts, ok := options.(*ebiten.DrawRectShaderOptions)
	if ok {
		as.shaderOpts = shaderOpts
	}

}

func NewAnimatedSprite() *AnimatedSprite {

	as := AnimatedSprite{
		Sprite:      &Sprite{},
		Animation:   &animations.Animation{},
		SpriteSheet: &spritesheet.SpriteSheet{},
		frameImg:    &ebiten.Image{},
		effect:      &ebiten.Image{},
	}

	return &as
}

func (as *AnimatedSprite) ChangeAnimationSpeed(newSpeed float32) {
	as.Animation.SpeedInTPS = newSpeed
}

func (as *AnimatedSprite) TriggerEffect(image *ebiten.Image) {
	as.effect = image
}

func LoadPulseOutlineShader(us *Sprite) {
	ols := shaders.LoadOutlineShader()
	us.Shader = ols
	us.ShaderParams["Opacity"] = float32(0.0)
	us.ShaderParams["OutlineColor"] = [4]float32{0.2, 0.7, 0.2, 1.0}
	us.UpdateShaderParams = shaders.UpdatePulseWithText
}
