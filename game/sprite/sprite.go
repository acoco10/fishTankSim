package sprite

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"math"
)

type Sprite struct {
	AnimationMap                   map[string]*Animation
	LayerIndex                     int
	Img                            *ebiten.Image
	NormalMap                      *ebiten.Image
	Flip                           bool
	LinkedSprite                   *Sprite
	Scale                          float64
	X, Y                           float32 //these could really be 64
	Dy, Dx                         float32 //these could really be 64
	Shader                         *ebiten.Shader
	ShaderParams                   map[string]any
	CPUShaderParams                map[string]any
	UpdateShaderParams             func(map[string]any) map[string]any
	UpdateBothParams               func(map[string]any, map[string]any) (map[string]any, map[string]any)
	Remove                         bool            //remove the entity that has this sprite
	UpdateFunc                     func(s *Sprite) //quicky script
	PublishedGraphicId             []int
	Focused                        bool
	Unfocusable                    bool
	AbleToBeUnfocusedAutomatically bool
	highlight                      bool
	DOptsUpdaterTag                string //for animating with external function at draw call
	DOptsUpdaterParams             map[string]float64
	CurrentAnimation               string
	IsBuffer                       bool
	BufferDst                      *ebiten.Image
	*XYUpdater
	frameImg   *ebiten.Image
	drawOpts   *ebiten.DrawImageOptions
	shaderOpts *ebiten.DrawRectShaderOptions
}

func (s *Sprite) GetAnimation() *Animation {
	if s.CurrentAnimation == "" {
		return nil
	}
	if s.AnimationMap[s.CurrentAnimation] == nil {
		log.Fatal("nil animation being checked error name:", s.CurrentAnimation)
	}
	if s.AnimationMap[s.CurrentAnimation].Img == nil {
		log.Fatal("select animation exists but image is nil:", s.CurrentAnimation)
	}
	return s.AnimationMap[s.CurrentAnimation]
}

func (s *Sprite) SpriteWidth() int {
	if s.GetAnimation() != nil {
		return s.GetAnimation().SpriteWidth
	} else {
		return s.Img.Bounds().Dx()
	}
}

func (s *Sprite) SpriteHeight() int {
	if s.GetAnimation() != nil {
		return s.GetAnimation().SpriteHeight
	} else {
		return s.Img.Bounds().Dy()
	}
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
	if s.CurrentAnimation != "" {
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
			s.ShaderParams = make(map[string]any)
			log.Printf("nil shader param issue when focusing sprite")
		}
		s.ShaderParams["OutlineColor"] = [4]float32{0.9, 0.9, 0.2, 1.0}
		s.Focused = true
	} else {
		s.Shader = registry.ShaderMap["NormalMapOutline"]
		if s.ShaderParams == nil {
			s.ShaderParams = make(map[string]any)
			log.Printf("nil shader param issue when focusing sprite")
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

	if s.CurrentAnimation != "" {
		DrawAnimation(s, screen)
		return
	}

	if s.Img == nil {
		log.Println("sprite with no img trying to draw")
		return
	}

	if s.DOptsUpdaterTag == "sway" {
		DrawSwayAnimation(s, screen)
		return
	}

	if s.DOptsUpdaterTag == "swirl" {
		DrawSwirlSprite(s, screen)
		return
	}

	if s.Shader != nil {
		shaderOpts := &ebiten.DrawRectShaderOptions{}
		if s.ShaderParams == nil {
			s.ShaderParams = make(map[string]any)
		}

		if s.Scale != 0.0 {
			shaderOpts.GeoM.Scale(s.Scale, s.Scale)
		}

		degrees, exists := s.DOptsUpdaterParams["degree"]
		if exists {
			// Get the center of the sprite (accounting for scale)
			b := s.Img.Bounds()
			centerX := float64(b.Dx()) / 2
			centerY := float64(b.Dy()) / 2
			if s.Scale != 0.0 {
				centerX *= s.Scale
				centerY *= s.Scale
			}

			// Translate to center, rotate, translate back
			shaderOpts.GeoM.Translate(-centerX, -centerY)
			shaderOpts.GeoM.Rotate(degrees)
			shaderOpts.GeoM.Translate(centerX, centerY)
		}

		if s.DOptsUpdaterTag == "flip" {
			shaderOpts.GeoM.Scale(-1, 1) // flip horizontally
			shaderOpts.GeoM.Translate(float64(s.Img.Bounds().Dx()), 0)
		}

		shaderOpts.GeoM.Translate(float64(s.X), float64(s.Y))
		shaderOpts.Images[0] = s.Img
		if s.NormalMap != nil {
			shaderOpts.Images[1] = s.NormalMap
		}

		shaderOpts.Uniforms = s.ShaderParams
		b := s.Img.Bounds()
		screen.DrawRectShader(b.Dx(), b.Dy(), s.Shader, shaderOpts)

		return
	}

	dOpts := &ebiten.DrawImageOptions{}
	if s.DOptsUpdaterTag == "spin" {
		spinAnimation(s, dOpts)
	}

	if s.DOptsUpdaterTag == "flip" {
		FlipSprite(s, dOpts)
	}

	if s.DOptsUpdaterParams["opacity"] > 0 {
		op := float32(s.DOptsUpdaterParams["opacity"])
		dOpts.ColorScale.ScaleAlpha(op) // Only affects alpha

	}

	degrees, exists := s.DOptsUpdaterParams["degree"]
	if exists {
		dOpts.GeoM.Rotate(degrees)
	}

	if s.Scale != 0.0 {
		dOpts.GeoM.Scale(s.Scale, s.Scale)
	}

	dOpts.GeoM.Translate(float64(s.X), float64(s.Y))

	if s.DOptsUpdaterParams != nil {
		dOpts.GeoM.Translate(s.DOptsUpdaterParams["offSetX"], s.DOptsUpdaterParams["offSetY"])
	}
	screen.DrawImage(s.Img, dOpts)

}

func UpdateNormalCursorForDebuggin(s *Sprite) {
	cx, cy := util.GetScaledCursorPosition()
	s.ShaderParams["Cursor"] = []float64{float64(cx), float64(cy)}
}

func FlipSprite(sprite *Sprite, dopts any) {
	opts, ok := dopts.(ebiten.DrawImageOptions)
	if !ok {
		sopts, _ := dopts.(ebiten.DrawRectShaderOptions)
		if sprite.AnimationMap[sprite.CurrentAnimation] != nil {
			sopts.GeoM.Scale(-1, 1) // flip horizontally
			sopts.GeoM.Translate(float64(sprite.AnimationMap[sprite.CurrentAnimation].SpriteWidth), 0)
		} else {
			sopts.GeoM.Translate(float64(sprite.Img.Bounds().Dx()), 0)

		}
		return
	}
	opts.GeoM.Scale(-1, 1) // flip horizontally
	if sprite.AnimationMap[sprite.CurrentAnimation] != nil {
		opts.GeoM.Translate(float64(sprite.AnimationMap[sprite.CurrentAnimation].SpriteWidth), 0)
	} else {
		opts.GeoM.Translate(float64(sprite.Img.Bounds().Dx()), 0)
	}
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
	rect := s.GetSpriteRect()
	return cpt.In(rect)

	/*	pixelX := point.X - int(s.X)
		pixelY := point.Y - int(s.Y)
		_, _, _, a := s.Img.At(pixelX, pixelY).RGBA()
		s.Img.Bounds()

		if a != 0.0 {
			return true
		}*/

}

func (s *Sprite) GetSpriteRect() image.Rectangle {

	width := s.SpriteWidth()
	height := s.SpriteHeight()
	rect := image.Rect(int(s.X), int(s.Y), int(s.X)+width, int(s.Y)+height)

	return rect
}

func (s *Sprite) SpriteHoveredWithBuffer(buffer int) bool {
	x, y := util.GetScaledCursorPosition()
	point := image.Point{X: x, Y: y}

	var bounds image.Rectangle
	if s.Img == nil {
		//prevent crash if sprite has no image, or get rect from animation img
		if s.CurrentAnimation != "" {
			bounds = s.GetSpriteRect().Bounds()
		} else {
			return false
		}
	} else {
		bounds = s.Img.Bounds()
	}
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

	ani := as.GetAnimation()
	shaderOpts := &ebiten.DrawRectShaderOptions{}

	if as.Scale > 0 {
		shaderOpts.GeoM.Scale(as.Scale, as.Scale)
	}

	shaderOpts.GeoM.Translate(float64(as.X), float64(as.Y-ani.OffSetY))

	as.shaderOpts = shaderOpts
	/*if as.ShaderParams != nil {
		cx, cy := util.GetScaledCursorPosition()
		as.ShaderParams["Cursor"] = []float64{float64(cx), float64(cy)}
	}*/
	drawOpts := &ebiten.DrawImageOptions{}

	if as.Scale > 0 {
		drawOpts.GeoM.Scale(as.Scale, as.Scale)
	}

	drawOpts.GeoM.Translate(float64(as.X), float64(as.Y-ani.OffSetY))
	as.drawOpts = drawOpts
	as.UpdateShader()

	as.GetAnimation().Update()

	UpdateSpriteFrameImg(as)
}

func (s *Sprite) SetAnimation(Ani string) {
	s.CurrentAnimation = Ani
	_, exists := s.AnimationMap[Ani]
	if !exists {
		log.Fatal("set a sprite animation that doesnt exist")
	}
	UpdateSpriteAnimation(s)
}

func UpdateSpriteFrameImg(as *Sprite) {

	ani := as.GetAnimation()
	if ani == nil {
		log.Fatal("Why is animation being checked in sprite with current animation empty")
	}
	frame := ani.Frame()
	frameRect := ani.Rect(frame)
	if ani.Img == nil {
		log.Fatal(as.CurrentAnimation, "this animation returns a nil image")
	}
	img := ani.Img.SubImage(frameRect).(*ebiten.Image)
	as.frameImg = img
}

func (a *Animation) GetFirstFrameAsStaticImage() *ebiten.Image {
	frameRect := a.SpriteSheet.Rect(1)
	img := a.Img.SubImage(frameRect).(*ebiten.Image)
	return img
}

func (a *Animation) GetLastFrameAsStaticImage() *ebiten.Image {
	frameRect := a.SpriteSheet.Rect(a.LastF)
	img := a.Img.SubImage(frameRect).(*ebiten.Image)
	return img
}

func (a *Animation) GetLastFrameNormalAsStaticImage() *ebiten.Image {
	frameRect := a.SpriteSheet.Rect(a.LastF)
	img := a.NormalImg.SubImage(frameRect).(*ebiten.Image)
	return img
}

func DrawAnimation(as *Sprite, screen *ebiten.Image) {
	ani := as.GetAnimation()
	frame := as.AnimationMap[as.CurrentAnimation].Frame()
	frameRect := as.AnimationMap[as.CurrentAnimation].Rect(frame)
	img := ani.Img.SubImage(frameRect).(*ebiten.Image)

	if ani.NormalImg != nil {
		DrawNormal(as, screen)
		return
	}

	if as.Shader != nil {
		if as.shaderOpts == nil {
			as.shaderOpts = &ebiten.DrawRectShaderOptions{}
		}
		as.shaderOpts.Images[0] = img
		as.shaderOpts.Uniforms = as.ShaderParams
		b := img.Bounds()
		screen.DrawRectShader(b.Dx(), b.Dy(), as.Shader, as.shaderOpts)
		return
	}
	if as.drawOpts == nil {
		as.drawOpts = &ebiten.DrawImageOptions{}
	}

	screen.DrawImage(img, as.drawOpts)
}

func DrawNormal(as *Sprite, screen *ebiten.Image) {
	ani := as.GetAnimation()
	if as.shaderOpts == nil {
		log.Printf("nil shader opts")
		as.shaderOpts = &ebiten.DrawRectShaderOptions{}
		as.shaderOpts.GeoM.Translate(float64(as.X), float64(as.Y))
	}

	if as.Shader == nil {
		shader := registry.ShaderMap["NormalMap"]
		as.Shader = shader
	}

	frame := ani.Frame()

	frameRect := ani.Rect(frame)

	diffuseImg := ani.Img.SubImage(frameRect).(*ebiten.Image)
	if diffuseImg == nil {
		log.Fatal("normal map sub rect is disposed")
	}

	normalImg := ani.NormalImg.SubImage(frameRect).(*ebiten.Image)
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
	if s.AnimationMap[s.CurrentAnimation] != nil {
		s.AnimationMap[s.CurrentAnimation].SpeedInTPS = newSpeed
	}
}

func LoadPulseOutlineNormalShader(us *Sprite) {
	ols := registry.ShaderMap["NormalMapOutline"]
	us.Shader = ols
	us.ShaderParams["Opacity"] = float32(0.0)
	us.ShaderParams["Cursor"] = []float64{0, 0, 100}
	us.ShaderParams["OutlineColor"] = [4]float32{0.2, 0.7, 0.2, 1.0}
	us.UpdateShaderParams = shaders.UpdatePulseWithText
}

func InitSwayAnimation(sp *Sprite, baseAmp float64) {
	sp.DOptsUpdaterTag = "sway"
	sp.DOptsUpdaterParams["time"] = 0
	sp.DOptsUpdaterParams["amp"] = baseAmp
}

func DrawSwayAnimation(sp *Sprite, screen *ebiten.Image) {
	sp.DOptsUpdaterParams["time"] += 0.016
	time := sp.DOptsUpdaterParams["time"]
	amp := sp.DOptsUpdaterParams["amp"]

	imgBounds := sp.Img.Bounds()
	imgHeight := imgBounds.Dy()
	sliceHeight := 2 // smaller slices = smoother sway

	ani := sp.AnimationMap["StartUp"]
	frameRect := ani.SpriteSheet.Rect(ani.LastF)
	amp = util.Lerp64(amp, 5, time/100)

	// sway amplitude
	baseSway := math.Sin(sp.DOptsUpdaterParams["time"]) * amp
	sp.DOptsUpdaterParams["amp"] = amp
	for y := 0; y < imgHeight; y += sliceHeight {
		// Factor goes from 0.0 at bottom to 1.0 at top
		// flip it so bottom = 0 sway, top = full sway
		progress := float64(imgHeight-y) / float64(imgHeight)

		// slice sway decreases smoothly toward the base
		sliceSway := baseSway * progress

		// slice rectangle (careful with slice bottom clamp)
		sliceBottom := y + sliceHeight
		if sliceBottom > imgHeight {
			sliceBottom = imgHeight
		}

		sliceRect := image.Rect(frameRect.Min.X, y, frameRect.Max.X, sliceBottom)
		subImg := ani.Img.SubImage(sliceRect).(*ebiten.Image)

		if ani.NormalImg != nil {
			shaderOpts := &ebiten.DrawRectShaderOptions{}
			normalSub := ani.NormalImg.SubImage(sliceRect).(*ebiten.Image)
			shaderOpts.GeoM.Translate(float64(sp.X)+sliceSway, float64(sp.Y)+float64(y))
			shaderOpts.Uniforms = sp.ShaderParams
			shaderOpts.Images[0] = subImg
			shaderOpts.Images[1] = normalSub
			b := subImg.Bounds()
			screen.DrawRectShader(b.Dx(), b.Dy(), sp.Shader, shaderOpts)
			continue
		}
		// draw slice
		dOpts := &ebiten.DrawImageOptions{}
		dOpts.GeoM.Translate(float64(sp.X)+sliceSway, float64(sp.Y)+float64(y))
		screen.DrawImage(subImg, dOpts)

	}
}

func DrawSwirlSprite(sp *Sprite, screen *ebiten.Image) {
	sp.DOptsUpdaterParams["time"] += 0.2 // Slower time increment
	//maxAngle := 5.0 * (math.Pi / 180)    // 3 degrees in radians

	if sp.Shader != nil {
		sopts := &ebiten.DrawRectShaderOptions{}
		sopts.GeoM.Translate(math.Sin(sp.DOptsUpdaterParams["time"])*2, math.Sin(sp.DOptsUpdaterParams["time"]))
		sopts.GeoM.Translate(float64(sp.X), float64(sp.Y))

		sopts.Images[0] = sp.Img
		sopts.Uniforms = sp.ShaderParams
		b := sp.Img.Bounds()
		screen.DrawRectShader(b.Dx(), b.Dy(), sp.Shader, sopts)
		return
	}

	dopts := &ebiten.DrawImageOptions{}
	dopts.GeoM.Translate(math.Sin(sp.DOptsUpdaterParams["time"])*2, math.Sin(sp.DOptsUpdaterParams["time"]))
	dopts.GeoM.Translate(float64(sp.X), float64(sp.Y))
	screen.DrawImage(sp.Img, dopts)

}
