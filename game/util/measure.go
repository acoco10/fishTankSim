package util

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func MeasureText(outputText string, fontsize float64, font string) (float64, float64) {
	face := registry.FontMap[font]
	if fontsize != 16 {
		face, _ = LoadFont(fontsize, font)
	}

	return text.Measure(outputText, face, 2)
}

func Lerp64(A float64, B float64, t float64) float64 {
	return A + (B-A)*t
}

func Lerp32(A float32, B float32, t float32) float32 {
	return A + (B-A)*t
}

func GetScaledCursorPosition() (int, int) {
	x, y := ebiten.CursorPosition()
	if registry.Config.Zoom {

		scaledX := float64(x) / registry.Config.ResolutionScalingF / registry.Config.ZoomFactor
		xOffSet := registry.Config.ZoomOffSetX / registry.Config.ZoomFactor
		xf := scaledX - xOffSet

		scaledY := float64(y-registry.Config.YOffset) / registry.Config.ResolutionScalingF / registry.Config.ZoomFactor
		yOffSet := registry.Config.ZoomOffSetY / registry.Config.ZoomFactor
		yf := scaledY - yOffSet

		return int(xf), int(yf)

	}
	return int(float64(x) / registry.Config.ResolutionScalingF), int(float64(y-registry.Config.YOffset) / registry.Config.ResolutionScalingF)

}
