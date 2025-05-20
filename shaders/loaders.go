package shaders

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func LoadOutlineShader() *ebiten.Shader {
	ols := []byte(OutlineShader)
	s, err := ebiten.NewShader(ols)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func LoadPulseOutlineShader() *ebiten.Shader {
	ols := []byte(PulseOutline)
	s, err := ebiten.NewShader(ols)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func LoadSolidColorShader() *ebiten.Shader {
	sls := []byte(SolidColor)
	s, err := ebiten.NewShader(sls)
	if err != nil {
		log.Fatal(err)
	}
	return s
}
