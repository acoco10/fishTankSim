package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type TextWithShader struct {
	insets            [2]float64
	renderShaderImage *ebiten.Image
	shader            *ebiten.Shader
	ShaderParams      map[string]any
	text              string
	textYinset        float64
	textXinset        float64
	FullyDrawn        bool
	face              text.Face
	clr               ebiten.ColorScale
	buffer            *ebiten.Image
	updateFunc        func(map[string]any) map[string]any
}

func NewTextWithMarkerShader(txt string, buff *ebiten.Image, insets [2]float64, clr ebiten.ColorScale, yInset float64, xInset float64) *TextWithShader {
	ts := &TextWithShader{}
	ts.text = txt
	ts.face = registry.FontMap["RockSalt_12"]

	ts.insets = insets
	ts.textYinset = yInset
	ts.textXinset = xInset
	ts.updateFunc = shaders.UpdateCounterOneShot
	ts.shader = registry.ShaderMap["HandWriting"]
	ts.clr = clr
	length, _ := util.MeasureText(txt, 16, "RockSalt_12")

	ts.ShaderParams = make(map[string]any)
	ts.ShaderParams["Speed"] = 1.5
	ts.ShaderParams["MaxCounter"] = int(length) * 2
	ts.ShaderParams["Counter"] = int(xInset - 2) //counter is space based

	ts.ShaderParams["MaxOpacity"] = 0.9

	ts.renderShaderImage = ebiten.NewImage(buff.Bounds().Dx(), buff.Bounds().Dy())
	ts.buffer = buff

	topts := &text.DrawOptions{}

	topts.ColorScale.SetR(242.0 / 255.0)
	topts.ColorScale.SetG(68.0 / 255.0)
	topts.ColorScale.SetB(51.0 / 255.0)
	topts.ColorScale.SetA(1)
	topts.GeoM.Translate(ts.textXinset, ts.textYinset)
	text.Draw(ts.renderShaderImage, ts.text, ts.face, topts)

	return ts
}

func (t *TextWithShader) IsFullyDrawn() bool {
	return t.FullyDrawn
}

func (t *TextWithShader) Update() {

	if t.ShaderParams == nil {
		println("nil map for params in text with shader")
		return
	}

	if t.updateFunc == nil {
		println("nil update func for draw text w/ shader")
		return
	}
	if t.ShaderParams["Counter"].(int) == t.ShaderParams["MaxCounter"].(int)-50 {
		t.FullyDrawn = true
	}
	t.ShaderParams = t.updateFunc(t.ShaderParams)

}

func (t *TextWithShader) Draw() {
	shaderOpts := &ebiten.DrawRectShaderOptions{}
	shaderOpts.Uniforms = t.ShaderParams
	shaderOpts.Images[0] = t.renderShaderImage
	t.buffer.DrawRectShader(t.renderShaderImage.Bounds().Dx(), t.renderShaderImage.Bounds().Dy(), t.shader, shaderOpts)

}
