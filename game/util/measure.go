package util

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

func MeasureNK57(outputText string, fontsize float64) (float64, float64) {

	face, err := LoadFont(fontsize, "nk57")
	if err != nil {
		log.Fatal("invalid font selected", err)
	}

	return text.Measure(outputText, face, 2)
}
