package util

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func MeasureText(outputText string, fontsize float64, font string) (float64, float64) {
	face := registry.FontMap[font]
	if fontsize != 16 {
		face, _ = LoadFont(fontsize, font)
	}

	return text.Measure(outputText, face, 2)
}
