package textEffects

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"log"
)

type TextWithShader struct {
	insets            [2]float64
	renderShaderImage *ebiten.Image
	image             *ebiten.Image
	shader            *ebiten.Shader
	shaderParams      map[string]any
	text              string
	FullyDrawn        bool
	face              text.Face
	clr               ebiten.ColorScale
	drawnToRender     bool
	updateFunc        func(map[string]any) map[string]any
}

func NewTextWithMarkerShader(text string, rect image.Rectangle, insets [2]float64, clr ebiten.ColorScale) *TextWithShader {
	ts := &TextWithShader{}
	ts.text = text
	ts.face = registry.FontMap["RockSalt"]

	ts.insets = insets

	ts.updateFunc = shaders.UpdateCounterOneShot
	ts.shader = registry.ShaderMap["HandWriting"]
	ts.clr = clr
	length, _ := util.MeasureText(text, 16, "RockSalt")

	ts.shaderParams = make(map[string]any)
	ts.shaderParams["Counter"] = 0
	ts.shaderParams["MaxCounter"] = int(length)

	ts.renderShaderImage = ebiten.NewImage(rect.Dx(), rect.Dy())

	return ts
}

func (t *TextWithShader) IsFullyDrawn() bool {
	return t.FullyDrawn
}

func (t *TextWithShader) Update() {

	if t.shaderParams == nil {
		println("nil map for params in text with shader")
		return
	}

	if t.updateFunc == nil {
		println("nil update func for draw text w/ shader")
		return
	}

	t.shaderParams = t.updateFunc(t.shaderParams)

	counter, ok := t.shaderParams["Counter"].(int)
	if !ok {
		log.Printf("Text Shader shader Parameters were reset but function is still updating skipping to avoid nil pointer errors\n")
		return
	}

	if counter == 1 {
		log.Printf("updating shader: |%s|", t.text)
	}

	maxCounter, ok := t.shaderParams["MaxCounter"].(int)
	if !ok {
		log.Printf("Nil max counter value in text shader updater paramaters \n")
		return
	}

	if counter >= maxCounter && !t.FullyDrawn {
		log.Printf("Text shader: |%s|is fully Drawn", t.text)
		t.FullyDrawn = true
	}

}

func (t *TextWithShader) Draw(dst *ebiten.Image) {
	if !t.FullyDrawn {
		log.Printf("Drawing text shader: |%s|", t.text)
		if t.shader == nil {
			t.FullyDrawn = true
			return
		}
		if t.shaderParams == nil {
			t.FullyDrawn = true
			log.Printf("Text Shader shader Parameters were reset but draw is being called skipping to avoid nil pointer errors\n")
			return
		}
		dopts := &text.DrawOptions{}

		dopts.ColorScale.SetR(0)
		dopts.ColorScale.SetG(1)
		dopts.ColorScale.SetB(1)
		dopts.ColorScale.SetA(1)

		dopts.GeoM.Translate(t.insets[0], t.insets[1])

		text.Draw(t.renderShaderImage, t.text, t.face, dopts)

		shaderOpts := &ebiten.DrawRectShaderOptions{}
		shaderOpts.Uniforms = t.shaderParams
		shaderOpts.Images[0] = t.renderShaderImage
		dst.DrawRectShader(dst.Bounds().Dx(), dst.Bounds().Dy(), t.shader, shaderOpts)
		return
	}
	//text.Draw(dst, t.text, t.face, dopts)
}
