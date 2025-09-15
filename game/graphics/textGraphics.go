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
	text               *string
	x, y               float64
	face               text.Face
	face2              text.Face
	opacity            float32
	GraphicId          int
	pulse              bool
	colorScale         ebiten.ColorScale
	altColoreScale     ebiten.ColorScale
	offScreen          *ebiten.Image
	shaderParams       map[string]any
	deposeAfterNFrames int
	lifeTime           int
}

func NewGraphicText(face text.Face, face2 text.Face, outputText *string, x float64, y float64, pulse bool, color ebiten.ColorScale, spriteSize float64, fade bool, frameLife int) int {
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

	cs := &ebiten.ColorScale{}
	cs.SetA(0.75)
	cs.SetR(0.0)
	cs.SetB(0.0)
	cs.SetG(0.0)

	ft.altColoreScale = color
	ft.colorScale = *cs

	ft.GraphicId = AssignAndIncrement(&ft)
	ft.deposeAfterNFrames = frameLife
	if !fade {
		ft.opacity = 1.0
	}

	ft.face2 = face2
	return ft.GraphicId

}

func NewOutlineGraphicText(outputText *string, fontsize float64, x float64, y float64, pulse bool, color ebiten.ColorScale, spriteSize float64, fade bool, frameLife int) int {

	face := registry.FontMap["RockSalt"]

	face2 := registry.FontMap["RockSalt"]
	return NewGraphicText(face, face2, outputText, x, y, pulse, color, spriteSize, fade, frameLife)
}

func NewNkTextGraphic(outputText *string, fontsize float64, x float64, y float64, pulse bool, color ebiten.ColorScale, spriteSize float64, fade bool, frameLife int) int {
	face, err := util.LoadFont(fontsize, "nk57")
	if err != nil {
		log.Fatal("font dont exist or something")
	}
	return NewGraphicText(face, nil, outputText, x, y, pulse, color, spriteSize, fade, frameLife)
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
	tOpts.GeoM.Translate(+2, -2)
	text.Draw(screen, *ft.text, ft.face, tOpts)
}

func (ft *FadeInText) Scaled() ScaledType {
	return UnScaled
}

func (ft *FadeInText) Update() {
	if ft.opacity < 1.0 {
		ft.opacity += 0.05 // adjust fade-in speed here
	}
	if ft.opacity > 1.0 && ft.pulse {
		ft.opacity = 0.0
	}
	if ft.opacity > 1.0 && !ft.pulse {
		ft.opacity = 1.0
	}
	if ft.lifeTime >= ft.deposeAfterNFrames && ft.deposeAfterNFrames != 0 {
		DeInitGraphicId(ft.GraphicId)
	}

	ft.lifeTime++

}

func (ft *FadeInText) Coords() (float64, float64) {
	return ft.x, ft.y
}
