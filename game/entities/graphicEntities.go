package entities

import (
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/hajimehoshi/ebiten/v2"
)

type ImageGraphic struct {
	*sprite.Sprite
	dst    *ebiten.Image
	X      float64
	Y      float64
	shader *ebiten.Shader
}

func (ig ImageGraphic) Draw() {

	shaderOpts := &ebiten.DrawRectShaderOptions{}
	shaderOpts.Uniforms = ig.ShaderParams
	shaderOpts.Images[0] = ig.Img

	ig.dst.DrawRectShader(ig.dst.Bounds().Dx(), ig.dst.Bounds().Dy(), ig.shader, shaderOpts)
}
