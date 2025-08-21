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

func LoadTextOutline() *ebiten.Shader {
	ols := []byte(textOutline)
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
		log.Printf("Couldnt Load Rotating Highlight shader %q", err)
	}
	return s
}

func LoadOnePointLightingBlue() *ebiten.Shader {
	opl := []byte(OnePointLightingBlue)
	s, err := ebiten.NewShader(opl)
	if err != nil {
		log.Fatal("Couldnt load one point shader", err)
	}
	return s
}

func LoadOnePointLightingNeutral() *ebiten.Shader {
	opl := []byte(OnePointLightingNeutral)
	s, err := ebiten.NewShader(opl)
	if err != nil {
		log.Printf("Couldnt Load one point lighting shader %q", err)
	}
	return s
}

func LoadSpriteLighting() *ebiten.Shader {
	osl := []byte(SpriteLightingEffect)
	s, err := ebiten.NewShader(osl)
	if err != nil {
		log.Printf("Couldnt sprite lighting shader %q", err)
	}
	return s
}

func LoadHandWritingShader() *ebiten.Shader {
	hws := []byte(HandWritingEffect)
	s, err := ebiten.NewShader(hws)
	if err != nil {
		log.Printf("Couldnt load handwriting shader %q", err)
	}
	return s
}

func LoadEraseShader() *ebiten.Shader {
	es := []byte(EraseEffect)
	s, err := ebiten.NewShader(es)
	if err != nil {
		log.Printf("Couldnt load erase shader %q", err)
	}
	return s
}

func LoadNormalMapShader() *ebiten.Shader {
	es := []byte(NormalMap)
	s, err := ebiten.NewShader(es)
	if err != nil {
		log.Fatal("Couldnt load normal map shader %q", err)
	}

	return s
}

func LoadTurnOffTheLights() *ebiten.Shader {
	es := []byte(TurnOffLight)
	s, err := ebiten.NewShader(es)
	if err != nil {
		log.Printf("Couldnt load turn off the lights shader %q", err)
	}

	return s
}

func LoadDayLight() *ebiten.Shader {
	es := []byte(DayLight)
	s, err := ebiten.NewShader(es)
	if err != nil {
		log.Printf("Couldnt load daylight%q", err)
	}

	return s
}

func LoadWaterShader() *ebiten.Shader {
	es := []byte(Water)
	s, err := ebiten.NewShader(es)
	if err != nil {
		log.Fatal("couldnt load new shader ", err)
	}
	return s
}

func LoadWallShader() *ebiten.Shader {
	es := []byte(Wall)
	s, err := ebiten.NewShader(es)
	if err != nil {
		log.Fatal("couldnt load wall shader ", err)
	}
	return s
}

func LoadPHShader() *ebiten.Shader {
	es := []byte(PHEffect)
	s, err := ebiten.NewShader(es)
	if err != nil {
		log.Fatal("couldnt load ph shader ", err)
	}
	return s
}

func LoadExperiment() *ebiten.Shader {

	es := []byte(FigureOutCoords)

	s, err := ebiten.NewShader(es)

	if err != nil {
		log.Fatal("couldnt load experiment shader ", err)
	}

	return s
}

func LoadHighlightShader() *ebiten.Shader {
	es := []byte(Highlight)
	s, err := ebiten.NewShader(es)
	if err != nil {
		log.Fatal("Couldn't load highlight shader ", err)
	}
	return s
}

func LoadStomachFillShader() *ebiten.Shader {
	es := []byte(StomachFill)
	s, err := ebiten.NewShader(es)
	if err != nil {
		log.Fatal("Couldn't load stomach shader ", err)
	}
	return s
}

func LoadNormalMapOutlined() *ebiten.Shader {
	es := []byte(NormalMapOutlined)
	s, err := ebiten.NewShader(es)
	if err != nil {
		log.Fatal("Couldn't load stomach shader ", err)
	}
	return s
}
