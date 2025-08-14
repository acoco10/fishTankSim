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

	turnOff := shaders.LoadTurnOffTheLights()

	dayLight := shaders.LoadDayLight()

	water := shaders.LoadWaterShader()

	wall := shaders.LoadWallShader()

	ph := shaders.LoadPHShader()

	highlight := shaders.LoadHighlightShader()

	stomach := shaders.LoadStomachFillShader()

	normalMapOutline := shaders.LoadNormalMapOutlined()

	registry.ShaderMap = make(map[string]*ebiten.Shader)

	// I gave these capital letters to mirror go's way of making structs importable/ kage but im not sure i like it
	registry.ShaderMap["Outline"] = ols
	registry.ShaderMap["Erase"] = erase
	registry.ShaderMap["HandWriting"] = hwr
	registry.ShaderMap["NormalMap"] = normalMap
	registry.ShaderMap["OnePointLighting"] = opl
	registry.ShaderMap["TurnOff"] = turnOff
	registry.ShaderMap["DayLight"] = dayLight
	registry.ShaderMap["Water"] = water
	registry.ShaderMap["Wall"] = wall
	registry.ShaderMap["PH"] = ph
	registry.ShaderMap["Highlight"] = highlight
	registry.ShaderMap["Stomach"] = stomach
	registry.ShaderMap["NormalMapOutline"] = normalMapOutline
}
