package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

type HorizontalPosition uint8

const (
	Center HorizontalPosition = iota
	Custom
)

type Effect uint8

const (
	Outline Effect = iota
	OffSet  Effect = iota
	None
)

type FadeInText struct {
	text           *string
	x, y           float64
	face           text.Face
	face2          text.Face
	opacity        float32
	GraphicId      int
	pulse          bool
	colorScale     ebiten.ColorScale
	altColoreScale ebiten.ColorScale
	offScreen      *ebiten.Image
	shaderParams   map[string]any
	deInit         bool
}

func NewGraphicText(face text.Face, face2 text.Face, outputText *string, x float64, y float64, pulse bool, color ebiten.ColorScale, spriteSize float64, fade bool) int {
	ft := FadeInText{
		text: outputText,
		face: face,
		x:    x,
		y:    y}

	width, height := text.Measure(*outputText, face, 2)

	ft.x = x*registry.Config.ResolutionScalingF + (spriteSize * registry.Config.ResolutionScalingF / 2) - width/2
	ft.y = y*registry.Config.ResolutionScalingF - height + registry.Config.YOffsetF

	ft.opacity = 0.0
	ft.pulse = pulse

	ft.colorScale = color

	cs := &ebiten.ColorScale{}
	cs.SetA(1.0)
	cs.SetR(0.0)
	cs.SetB(0.0)
	cs.SetG(0.0)

	ft.altColoreScale = *cs

	ft.GraphicId = AssignAndIncrement(&ft)

	if !fade {
		ft.opacity = 1.0
	}

	ft.face2 = face2
	return ft.GraphicId

}

func NewOutlineGraphicText(outputText *string, fontsize float64, x float64, y float64, pulse bool, color ebiten.ColorScale, spriteSize float64, fade bool) int {

	face, err := util.LoadFont(fontsize, "reglisseOutlined")
	if err != nil {
		log.Fatal("invalid font selected", err)
	}
	face2, err := util.LoadFont(fontsize, "reglisseOutline")
	if err != nil {
		log.Fatal("invalid font selected", err)
	}
	return NewGraphicText(face, face2, outputText, x, y, pulse, color, spriteSize, fade)
}

func NewNkTextGraphic(outputText *string, fontsize float64, x float64, y float64, pulse bool, color ebiten.ColorScale, spriteSize float64, fade bool) int {
	face := registry.FontMap["nk57"]
	return NewGraphicText(face, nil, outputText, x, y, pulse, color, spriteSize, fade)
}

func (ft *FadeInText) AutoDeInit() bool {
	return ft.deInit
}

func (ft *FadeInText) Draw(screen *ebiten.Image) {

	tOpts := &text.DrawOptions{}

	cs := ft.colorScale
	cs.ScaleAlpha(ft.opacity)
	tOpts.ColorScale = cs
	tOpts.GeoM.Translate(ft.x, ft.y)

	if ft.face2 != nil {
		text.Draw(screen, *ft.text, ft.face2, tOpts)
	}

	cs2 := ft.altColoreScale
	cs2.ScaleAlpha(ft.opacity)
	tOpts.ColorScale = cs2
	//tOpts.GeoM.Translate(+2, -2)
	text.Draw(screen, *ft.text, ft.face, tOpts)
}

func (ft *FadeInText) Scaled() ScaledType {
	return UnScaled
}

func (ft *FadeInText) Update() {
	if ft.opacity < 1.0 {
		ft.opacity += 0.01 // adjust fade-in speed here
	}
	if ft.opacity > 1.0 && ft.pulse {
		ft.opacity = 0.0
	}
	if ft.opacity > 1.0 && !ft.pulse {
		ft.opacity = 1.0
	}
}

func (ft *FadeInText) Coords() (float64, float64) {
	return ft.x, ft.y
}
