package entImportableLoaders

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
)

func AddNormalMap(sprite *sprite.Sprite) {
	normalMapShader := registry.ShaderMap["NormalMap"]
	sprite.Shader = normalMapShader
	sprite.ShaderParams = make(map[string]any)
	sprite.ShaderParams["Cursor"] = []float64{400, 50}
}
