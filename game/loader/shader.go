package loader

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
)

func LoadShaderRegistry() {

	ols := shaders.LoadOutlineShader()

	normalMap := shaders.LoadNormalMapShader()

	erase := shaders.LoadEraseShader()

	hwr := shaders.LoadHandWritingShader()

	opl := shaders.LoadOnePointLightingBlue()

	registry.ShaderMap = make(map[string]*ebiten.Shader)

	registry.ShaderMap["Outline"] = ols
	registry.ShaderMap["Erase"] = erase
	registry.ShaderMap["HandWriting"] = hwr
	registry.ShaderMap["NormalMap"] = normalMap
	registry.ShaderMap["OnePointLighting"] = opl

}
