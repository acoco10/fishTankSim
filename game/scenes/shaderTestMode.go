//go:build old

package scenes

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func ShaderSwapper(f *FishScene2) {
	if inpututil.IsKeyJustPressed(ebiten.KeyE) && f.lightingShader == registry.ShaderMap["OnePointLighting"] {
		f.lightingShader = registry.ShaderMap["TurnOff"]
		f.shaderUpdater = shaders.FadeLightIntensityForTurnOff
		f.lightingShaderParams["LightIntensity"] = 1.0
		f.offScreenShader = nil
		f.lightingShaderParams["Counter"] = 0.0
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) && f.lightingShader == registry.ShaderMap["DayLight"] {
		f.SetNightLight()
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) && f.lightingShader == registry.ShaderMap["TurnOff"] {
		f.lightingShader = registry.ShaderMap["DayLight"]
		f.backGroundParams["Cursor"] = [2]float64{0, 0}
		f.shaderUpdater = nil
		f.lightingState = Day
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		f.lightingShader = nil
	}

}
