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

func LoadRotatingHighlightShader() *ebiten.Shader {
	rhls := []byte(RotatingHighlightOutline)
	s, err := ebiten.NewShader(rhls)
	if err != nil {
		log.Printf("Couldnt Load Rotating Highlight Shader %q", err)
	}
	return s
}

func LoadOnePointLighting() *ebiten.Shader {
	opl := []byte(OnePointLighting)
	s, err := ebiten.NewShader(opl)
	if err != nil {
		log.Printf("Couldnt Load one point lighting Shader %q", err)
	}
	return s
}

func LoadSpriteLighting() *ebiten.Shader {
	osl := []byte(SpriteLightingEffect)
	s, err := ebiten.NewShader(osl)
	if err != nil {
		log.Printf("Couldnt sprite lighting Shader %q", err)
	}
	return s
}
