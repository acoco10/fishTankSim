package interactableUIObjects

import (
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
)

type FishFoodSprite struct {
	*UiSprite
}

func (ff *FishFoodSprite) Draw(screen *ebiten.Image) {

	var paramaMapa = make(map[string]any)
	sopts := ebiten.DrawRectShaderOptions{}
	baseColor := [4]float64{0.2, 0.1, 0.05, 255}
	paramaMapa["BaseColor"] = baseColor

	shader := shaders.LoadSolidColorShader()

	sopts.GeoM.Translate(float64(ff.baseX), float64(ff.baseY))
	b := ff.Img.Bounds()

	sopts.Images[0] = ff.Img
	sopts.Uniforms = paramaMapa

	screen.DrawRectShader(b.Dx(), b.Dy(), shader, &sopts)
	opts := ebiten.DrawImageOptions{}

	if ff.state == Idle {
		opts.GeoM.Translate(float64(ff.X), float64(ff.Y))
		screen.DrawImage(ff.Img, &opts)
		opts.GeoM.Reset()
	} else if ff.state == Selected || ff.state == HoveredOver {
		if ff.HoverImg != nil {
			opts.GeoM.Translate(float64(ff.X), float64(ff.Y))
			opts.GeoM.Translate(float64(ff.AltOffsetX), float64(ff.AltOffsetY))
			screen.DrawImage(ff.HoverImg, &opts)
			opts.GeoM.Reset()
		}
	} else if ff.state == ClickedWhileBeingSelected {
		opts.GeoM.Translate(float64(ff.X), float64(ff.Y))
		screen.DrawImage(ff.AltImg, &opts)
		opts.GeoM.Reset()
	}

}

func (ff *FishFoodSprite) Subscribe() {

}
