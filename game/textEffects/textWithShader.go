package textEffects

import (
	"github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

type TextWithShader struct {
	image        *ebiten.Image
	shader       *ebiten.Shader
	shaderParams map[string]any
	text         string
	face         text.Face
	updateFunc   func(map[string]any) map[string]any
}

func NewTextWithShader(text string, dst *ebiten.Image) *TextWithShader {
	face, err := ui.LoadFont(18, "rockSalt")
	if err != nil {
		log.Fatal(err)
	}

	ts := &TextWithShader{}
	ts.text = text
	ts.face = face

	ts.updateFunc = shaders.UpdateCounter
	shader := shaders.LoadHandWritingShader()
	ts.shader = shader

	ts.shaderParams = make(map[string]any)
	ts.shaderParams["Counter"] = 0

	ts.image = dst

	return ts
}

func (t *TextWithShader) Update() {
	t.shaderParams = t.updateFunc(t.shaderParams)
}

func (t *TextWithShader) Draw(dst *ebiten.Image) {

	dopts := &text.DrawOptions{}
	shaderOpts := ebiten.DrawRectShaderOptions{}
	dopts.ColorScale.Scale(0, 0, 0, 1)
	dopts.GeoM.Translate(float64(10), float64(10))

	text.Draw(t.image, t.text, t.face, dopts)
	shaderOpts.Uniforms = t.shaderParams
	shaderOpts.Images[0] = t.image

	dst.DrawRectShader(t.image.Bounds().Dx(), t.image.Bounds().Dy(), t.shader, &shaderOpts)
}
