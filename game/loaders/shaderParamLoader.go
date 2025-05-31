package loaders

import (
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/shaders"
)

func LoadRotatingHighlightOutline(sprite *sprite.Sprite) {

	b := sprite.Img.Bounds()
	shaderParams := make(map[string]any)
	cpuShaderParams := make(map[string]any)
	X0 := float64(sprite.X)
	Y0 := float64(sprite.Y)

	shaderParams["HighLightRect"] = [2]float64{X0, Y0}

	shaderParams["OutlineColor"] = [4]float64{0.7, 0.7, 0.4, 255}
	shaderParams["HighLightColor"] = [4]float64{0.9, 0, 0, 255}
	shaderParams["HLRectSize"] = 5.0

	cpuShaderParams["imageRect"] = [2]float64{X0 + float64(b.Dx()), Y0 + float64(b.Dy())}
	cpuShaderParams["origin"] = [2]float64{float64(sprite.X), float64(sprite.Y)}
	cpuShaderParams["direction"] = geometry.Right

	sprite.ShaderParams = shaderParams
	sprite.CPUShaderParams = cpuShaderParams

	sprite.UpdateBothParams = shaders.CpuUpdateRotatingOutlineHighlight
}

func LoadRotatingHighlightOutlineAnimated(sprite *sprite.AnimatedSprite) {
	subImgRect := sprite.Rect(0)
	b := subImgRect.Bounds()
	scale := sprite.Scale
	if scale == 0 {
		scale = 1
	}
	shaderParams := make(map[string]any)
	cpuShaderParams := make(map[string]any)

	shaderParams["HighLightRect"] = [2]float64{0, 0}
	shaderParams["OutlineColor"] = [4]float64{0.8, 0.8, 0.0, 200}
	shaderParams["HighLightColor"] = [4]float64{0.99, 0.99, 0.99, 255}
	shaderParams["HLRectSize"] = 5.0

	cpuShaderParams["imageRect"] = [2]float64{float64(b.Dx()) * scale, float64(b.Dy()) * scale}
	cpuShaderParams["origin"] = [2]float64{float64(sprite.X), float64(sprite.Y)}
	cpuShaderParams["direction"] = geometry.Right
	cpuShaderParams["hlRectPoint"] = [2]float64{0, 0}

	sprite.ShaderParams = shaderParams
	sprite.CPUShaderParams = cpuShaderParams

	sprite.UpdateBothParams = shaders.CpuUpdateRotatingOutlineHighlight

}

func LoadLightingParamaters(sprite *sprite.AnimatedSprite) {

	shaderParams := make(map[string]any)
	cpuShaderParams := make(map[string]any)

	b := sprite.Img.Bounds()
	shaderParams["ImgRect"] = [4]float64{float64(b.Dx()), float64(b.Dy())}
	shaderParams["LightPoint"] = [2]float64{500, 500}
	sprite.ShaderParams = shaderParams
	sprite.CPUShaderParams = cpuShaderParams
}

func LoadSpriteLightingParams(sprite *sprite.AnimatedSprite) {
	shaderParams := make(map[string]any)
	b := sprite.Img.Bounds()
	shaderParams["ImgRect"] = [4]float64{float64(b.Dx()), float64(b.Dy())}
	shaderParams["LightPoint"] = [2]float64{150, 0}
	sprite.ShaderParams = shaderParams
}
