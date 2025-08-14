package sprite

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"math"
)

type Sprite struct {
	LayerIndex                     int
	Img                            *ebiten.Image
	NormalMap                      *ebiten.Image
	Flip                           bool
	LinkedSprite                   *Sprite
	Scale                          float64
	X, Y                           float32 //these could really be 64
	Z                              int     //not used yet, layer based not math based haha
	Dy, Dx                         float32 //these could really be 64
	Shader                         *ebiten.Shader
	ShaderParams                   map[string]any
	CPUShaderParams                map[string]any
	UpdateShaderParams             func(map[string]any) map[string]any
	UpdateBothParams               func(map[string]any, map[string]any) (map[string]any, map[string]any)
	Remove                         bool //stop drawing this sprite
	UpdateFunc                     func(s *Sprite)
	PublishedGraphicId             []int
	Focused                        bool
	Unfocusable                    bool
	AbleToBeUnfocusedAutomatically bool
	highlight                      bool
	DOptsUpdaterTag                string //for animating with external function
	DOptsUpdaterParams             map[string]float64
	*XYUpdater
	*animations.Animation
	*spritesheet.SpriteSheet
	frameImg   *ebiten.Image
	drawOpts   *ebiten.DrawImageOptions
	shaderOpts *ebiten.DrawRectShaderOptions //flexible interface for animation checks needed at draw or update time
}

func (s *Sprite) SavePublishedGraphicID(id int) {
	s.PublishedGraphicId = append(s.PublishedGraphicId, id)
}
func (s *Sprite) Update() {
	if s.XYUpdater != nil {
		s.XYUpdater.Update(s)
	}
	if s.LinkedSprite != nil {
		s.LinkedSprite.Update()
	}
	if s.Animation != nil {
		UpdateSpriteAnimation(s)
	}
	if s.DOptsUpdaterParams == nil {
		s.DOptsUpdaterParams = make(map[string]float64)
	}

	if s.UpdateFunc != nil {
		s.UpdateFunc(s)
	}

	s.UpdateShader()
}

func (s *Sprite) Focus() {
	if s.NormalMap == nil {
		s.Shader = registry.ShaderMap["Outline"]
		if s.ShaderParams == nil {
			log.Fatal("nil shader param issue when focusing sprite")
		}
		s.ShaderParams["OutlineColor"] = [4]float32{0.9, 0.9, 0.2, 1.0}
		s.Focused = true
	} else {
		s.Shader = registry.ShaderMap["NormalMapOutline"]
		if s.ShaderParams == nil {
			log.Fatal("nil shader param issue when focusing sprite")
		}
		s.ShaderParams["OutlineColor"] = [4]float32{0.9, 0.9, 0.2, 1.0}
		s.Focused = true
	}
}

func (s *Sprite) UnFocus() {
	if s.NormalMap != nil {
		s.Shader = registry.ShaderMap["NormalMap"]
	} else {
		s.Shader = nil
	}
	s.Focused = false

}

func (s *Sprite) Draw(screen *ebiten.Image) {
	if s.LinkedSprite != nil {
		s.LinkedSprite.Draw(screen)
	}
	if s.Animation != nil {
		DrawAnimation(s, screen)
		return
	}

	if s.Img == nil {
		log.Println("sprite with no img trying to draw")
		return
	}

	if s.Shader != nil {
		shaderOpts := &ebiten.DrawRectShaderOptions{}

		degrees, exists := s.DOptsUpdaterParams["degree"]
		if exists {
			shaderOpts.GeoM.Rotate(degrees)
		}
		if s.Scale != 0.0 {
			shaderOpts.GeoM.Scale(s.Scale, s.Scale)
		}
		shaderOpts.GeoM.Translate(float64(s.X), float64(s.Y))
		shaderOpts.Images[0] = s.Img
		shaderOpts.Uniforms = s.ShaderParams
		b := s.Img.Bounds()
		screen.DrawRectShader(b.Dx(), b.Dy(), s.Shader, shaderOpts)

		return
	}

	dOpts := &ebiten.DrawImageOptions{}
	if s.DOptsUpdaterTag == "spin" {
		spinAnimation(s, dOpts)
	}

	degrees, exists := s.DOptsUpdaterParams["degree"]
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
	return s.Remove
}

func spinAnimation(s *Sprite, dOpts *ebiten.DrawImageOptions) {
	theta := s.DOptsUpdaterParams["degree"]
	dOpts.GeoM.Rotate(theta)
	dOpts.GeoM.Translate(float64(s.Img.Bounds().Dx()/2), float64(s.Img.Bounds().Dy()/2))
	s.DOptsUpdaterParams["degree"] += 0.1
	if s.DOptsUpdaterParams["degree"] >= math.Pi {
		s.DOptsUpdaterTag = ""
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
	x, y := util.GetScaledCursorPosition()

	cpt := image.Point{x, y}
	rect := s.getSpriteRect()
	return cpt.In(rect)

	/*	pixelX := point.X - int(s.X)
		pixelY := point.Y - int(s.Y)
		_, _, _, a := s.Img.At(pixelX, pixelY).RGBA()
		s.Img.Bounds()

		if a != 0.0 {
			return true
		}*/

}

func (s *Sprite) getSpriteRect() image.Rectangle {
	b := s.Img.Bounds()
	width := b.Dx()
	height := b.Dy()
	rect := image.Rect(int(s.X), int(s.Y), int(s.X)+width, int(s.Y)+height)

	if s.Animation != nil && s.frameImg != nil {
		b = s.frameImg.Bounds()
		width = b.Dx()
		height = b.Dy()
		rect = image.Rect(int(s.X), int(s.Y), int(s.X)+width, int(s.Y)+height)
		if s.Flip {
			rect = image.Rect(int(s.X)-width, int(s.Y), int(s.X), int(s.Y)+height)
		}
	}

	return rect
}

func (s *Sprite) SpriteHoveredWithBuffer(buffer int) bool {
	x, y := util.GetScaledCursorPosition()
	point := image.Point{X: x, Y: y}

	bounds := s.Img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	rect := image.Rect(
		int(s.X)-buffer,
		int(s.Y)-buffer,
		int(s.X)+width+buffer,
		int(s.Y)+height+buffer+20, // If you need extra bottom padding
	)

	return point.In(rect)
}

type AnimatedSprite struct {
	*Sprite
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

func UpdateSpriteAnimation(as *Sprite) {

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
	UpdateSpriteFrameImg(as)
}

func UpdateSpriteFrameImg(as *Sprite) {
	frame := as.Frame()
	frameRect := as.SpriteSheet.Rect(frame)
	img := as.Img.SubImage(frameRect).(*ebiten.Image)
	as.frameImg = img
}

func (as *Sprite) GetFirstFrameAsStaticImage() *ebiten.Image {
	frameRect := as.SpriteSheet.Rect(1)
	img := as.Img.SubImage(frameRect).(*ebiten.Image)
	return img
}

func DrawAnimation(as *Sprite, screen *ebiten.Image) {
	frame := as.Frame()
	frameRect := as.SpriteSheet.Rect(frame)
	img := as.Img.SubImage(frameRect).(*ebiten.Image)
	if as.frameImg != nil {
		// for debugging cursor hovered
		/*	rect := as.getSpriteRect()
			if as.SpriteHovered() {
				vector.StrokeRect(screen, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), 2.0, colornames.Teal, false)
			} else {
				vector.StrokeRect(screen, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), 2.0, colornames.Crimson, false)
			}*/
	}
	if as.NormalMap != nil {
		DrawNormal(as, screen)
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

func DrawNormal(as *Sprite, screen *ebiten.Image) {

	if as.shaderOpts == nil {
		log.Printf("nil shader opts")
		as.shaderOpts = &ebiten.DrawRectShaderOptions{}
		as.shaderOpts.GeoM.Translate(float64(as.X), float64(as.Y))
	}

	if as.Shader == nil {
		shader := registry.ShaderMap["NormalMap"]
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

func (as *Sprite) UpdateOpts(options any) {

	opts, ok := options.(*ebiten.DrawImageOptions)
	if ok {
		as.drawOpts = opts
	}

	shaderOpts, ok := options.(*ebiten.DrawRectShaderOptions)
	if ok {
		as.shaderOpts = shaderOpts
	}

}

func (s *Sprite) ChangeAnimationSpeed(newSpeed float32) {
	if s.Animation != nil {
		s.Animation.SpeedInTPS = newSpeed
	}
}

func LoadPulseOutlineShader(us *Sprite) {
	ols := registry.ShaderMap["Outline"]
	us.Shader = ols
	us.ShaderParams["Opacity"] = float32(0.0)
	us.ShaderParams["OutlineColor"] = [4]float32{0.2, 0.7, 0.2, 1.0}
	us.UpdateShaderParams = shaders.UpdatePulseWithText
}
