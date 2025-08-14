package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

type TextWithShader struct {
	insets            [2]float64
	renderShaderImage *ebiten.Image
	shader            *ebiten.Shader
	shaderParams      map[string]any
	text              string
	textYinset        float64
	textXinset        float64
	FullyDrawn        bool
	face              text.Face
	clr               ebiten.ColorScale
	drawnToRender     bool
	buffer            *ebiten.Image
	Id                int
	updateFunc        func(map[string]any) map[string]any
}

func NewTextWithMarkerShader(text string, buff *ebiten.Image, insets [2]float64, clr ebiten.ColorScale, yInset float64, xInset float64) *TextWithShader {
	ts := &TextWithShader{}
	ts.text = text
	ts.face = registry.FontMap["RockSalt"]

	ts.insets = insets
	ts.textYinset = yInset
	ts.textXinset = xInset
	ts.updateFunc = shaders.UpdateCounterOneShot
	ts.shader = registry.ShaderMap["HandWriting"]
	ts.clr = clr
	length, _ := util.MeasureText(text, 16, "RockSalt")

	ts.shaderParams = make(map[string]any)
	ts.shaderParams["Counter"] = 0 ///dumb way to not have to fix the logic of this shader
	ts.shaderParams["Speed"] = 1.6
	ts.shaderParams["MaxCounter"] = int(insets[0]) + int(length)

	ts.renderShaderImage = ebiten.NewImage(buff.Bounds().Dx(), buff.Bounds().Dy())
	ts.buffer = buff
	return ts
}

func (t *TextWithShader) Scaled() ScaledType {
	return UnScaled
}

func (t *TextWithShader) IsFullyDrawn() bool {
	return t.FullyDrawn
}

func (t *TextWithShader) Update() {

	if t.FullyDrawn {
		return
	}

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

	topts := &text.DrawOptions{}
	dopts := &ebiten.DrawImageOptions{}

	topts.ColorScale.SetR(0.1)
	topts.ColorScale.SetG(0.2)
	topts.ColorScale.SetB(0.6)
	topts.ColorScale.SetA(1)
	topts.GeoM.Translate(30, t.textYinset)
	text.Draw(t.renderShaderImage, t.text, t.face, topts)

	shaderOpts := &ebiten.DrawRectShaderOptions{}
	shaderOpts.Uniforms = t.shaderParams
	shaderOpts.Images[0] = t.renderShaderImage
	t.buffer.DrawRectShader(t.renderShaderImage.Bounds().Dx(), t.renderShaderImage.Bounds().Dy(), t.shader, shaderOpts)

	dopts.GeoM.Translate(t.insets[0], t.insets[1])

	//dst.DrawImage(t.buffer, dopts)
}
