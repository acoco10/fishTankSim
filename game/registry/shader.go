package registry

import (
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
)

var ShaderMap map[string]*ebiten.Shader

const (
	Outline          = "Outline"
	Erase            = "Erase"
	HandWriting      = "HandWriting"
	NormalMap        = "NormalNap"
	PH               = "PH"
	Highlight        = "Highlight"
	Stomach          = "Stomach"
	NormalMapOutline = "NormalMapOutline"
	PulseHighlight   = "PulseHighlight"
	Lowlight         = "Lowlight"
)

func LoadShaderRegistry() {

	ShaderMap = make(map[string]*ebiten.Shader)

	ShaderMap["Outline"] = shaders.LoadOutlineShader()
	ShaderMap["Erase"] = shaders.LoadEraseShader()
	ShaderMap["HandWriting"] = shaders.LoadHandWritingShader()
	ShaderMap["NormalMap"] = shaders.LoadNormalMapShader()
	ShaderMap["PH"] = shaders.LoadPHShader()
	ShaderMap["Highlight"] = shaders.LoadHighlightShader()
	ShaderMap["Stomach"] = shaders.LoadStomachFillShader()
	ShaderMap["NormalMapOutline"] = shaders.LoadNormalMapOutlined()
	ShaderMap["PulseHighlight"] = shaders.LoadPulseHighlight()
	ShaderMap[Lowlight] = shaders.LoadLowLightShader()
	//ShaderMap["PulseOutlineNormal"] = shaders.LoadPulseOutlineNormal()
}
