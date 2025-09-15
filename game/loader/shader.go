package loader

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
)

func LoadShaderRegistry() {

	registry.ShaderMap = make(map[string]*ebiten.Shader)

	registry.ShaderMap["Outline"] = shaders.LoadOutlineShader()
	registry.ShaderMap["Erase"] = shaders.LoadEraseShader()
	registry.ShaderMap["HandWriting"] = shaders.LoadHandWritingShader()
	registry.ShaderMap["NormalMap"] = shaders.LoadNormalMapShader()
	registry.ShaderMap["PH"] = shaders.LoadPHShader()
	registry.ShaderMap["Highlight"] = shaders.LoadHighlightShader()
	registry.ShaderMap["Stomach"] = shaders.LoadStomachFillShader()
	registry.ShaderMap["NormalMapOutline"] = shaders.LoadNormalMapOutlined()
	registry.ShaderMap["PulseHighlight"] = shaders.LoadPulseHighlight()
	//registry.ShaderMap["PulseOutlineNormal"] = shaders.LoadPulseOutlineNormal()
}
