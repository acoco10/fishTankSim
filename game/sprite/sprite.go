package sprite

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"math"
)

type Sprite struct {
	Img                *ebiten.Image
	NormalMap          *ebiten.Image
	Scale              float64
	X, Y               float32
	Z                  int
	Dy, Dx             float32
	Shader             *ebiten.Shader
	ShaderParams       map[string]any
	CPUShaderParams    map[string]any
	UpdateShaderParams func(map[string]any) map[string]any
	UpdateBothParams   func(map[string]any, map[string]any) (map[string]any, map[string]any)
	remove             bool
	UpdateFunc         func(s *Sprite)
	effect             *AnimatedSprite
	highlight          bool
	DoptsUpdaterTag    string
	DoptsUpdaterParams map[string]float64
}

func (s *Sprite) Update() {
	if s.DoptsUpdaterParams == nil {
		s.DoptsUpdaterParams = make(map[string]float64)
	}

	if s.UpdateFunc != nil {
		s.UpdateFunc(s)
	}

	s.UpdateShader()
}

func (s *Sprite) Draw(screen *ebiten.Image) {

	if s.Shader != nil {
		shaderOpts := &ebiten.DrawRectShaderOptions{}

		degrees, exists := s.DoptsUpdaterParams["degree"]
		if exists {
			shaderOpts.GeoM.Rotate(degrees)
		}
		shaderOpts.GeoM.Translate(float64(s.X), float64(s.Y))
		shaderOpts.Images[0] = s.Img
		shaderOpts.Uniforms = s.ShaderParams
		b := s.Img.Bounds()
		screen.DrawRectShader(b.Dx(), b.Dy(), s.Shader, shaderOpts)

		return
	}

	dOpts := &ebiten.DrawImageOptions{}
	if s.DoptsUpdaterTag == "spin" {
		spinAnimation(s, dOpts)
	}

	degrees, exists := s.DoptsUpdaterParams["degree"]
	if exists {
		dOpts.GeoM.Rotate(degrees)
	}

	if s.Scale != 0.0 {
		dOpts.GeoM.Scale(s.Scale, s.Scale)
	}

	dOpts.GeoM.Translate(float64(s.X), float64(s.Y))

	screen.DrawImage(s.Img, dOpts)

}

func (s *Sprite) ShouldRemove() bool {
	return s.remove
}

func spinAnimation(s *Sprite, dOpts *ebiten.DrawImageOptions) {
	theta := s.DoptsUpdaterParams["degree"]
	dOpts.GeoM.Rotate(theta)
	dOpts.GeoM.Translate(float64(s.Img.Bounds().Dx()/2), float64(s.Img.Bounds().Dy()/2))
	s.DoptsUpdaterParams["degree"] += 0.1
	if s.DoptsUpdaterParams["degree"] >= math.Pi {
		s.DoptsUpdaterTag = ""
	}
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
	drawOpts   *ebiten.DrawImageOptions
	shaderOpts *ebiten.DrawRectShaderOptions
}

func (s *Sprite) AddEffect() {

}

func (s *Sprite) Coord() (float32, float32) {
	return s.X, s.Y
}

func (s *Sprite) LoadShader(shader *ebiten.Shader) {
	println("loading shader")
	s.Shader = shader
}

func (s *Sprite) Highlighted() bool {
	return s.highlight
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
func (as *AnimatedSprite) GetFirstFrameAsStaticImage() *ebiten.Image {
	frameRect := as.SpriteSheet.Rect(1)
	img := as.Img.SubImage(frameRect).(*ebiten.Image)
	return img
}
func (as *AnimatedSprite) Draw(screen *ebiten.Image) {
	frame := as.Frame()
	frameRect := as.SpriteSheet.Rect(frame)
	img := as.Img.SubImage(frameRect).(*ebiten.Image)

	if as.NormalMap != nil {
		as.DrawNormal(screen)
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

func (as *AnimatedSprite) DrawNormal(screen *ebiten.Image) {

	if as.shaderOpts == nil {
		log.Printf("nil shader opts")
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
	}

	return &as
}

func (as *AnimatedSprite) ChangeAnimationSpeed(newSpeed float32) {
	as.Animation.SpeedInTPS = newSpeed
}

func LoadPulseOutlineShader(us *Sprite) {
	ols := shaders.LoadOutlineShader()
	us.Shader = ols
	us.ShaderParams["Opacity"] = float32(0.0)
	us.ShaderParams["OutlineColor"] = [4]float32{0.2, 0.7, 0.2, 1.0}
	us.UpdateShaderParams = shaders.UpdatePulseWithText
}
