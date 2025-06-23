package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

type FadeInText struct {
	text       string
	x, y       float64
	face       text.Face
	opacity    float32
	GraphicId  int
	pulse      bool
	colorScale ebiten.ColorScale
}

func NewGraphicText(outputText string, fontsize float64, x float64, y float64, pulse bool, color ebiten.ColorScale, spriteSize float64) int {

	face, err := util.LoadFont(fontsize, "nk57")
	if err != nil {
		log.Fatal("invalid font selected", err)
	}
	ft := FadeInText{
		text: outputText,
		face: face,
		x:    x,
		y:    y}

	width, height := text.Measure(outputText, face, 2)

	ft.x = (x + (spriteSize / 2)) - width/2
	ft.y = y - height

	ft.opacity = 0.0
	ft.pulse = pulse

	ft.colorScale = color

	ft.GraphicId = AssignAndIncrement(&ft)

	return ft.GraphicId
}

func (ft *FadeInText) Draw(screen *ebiten.Image) {

	// Draw the text once onto the img, if not already
	tOpts := &text.DrawOptions{}
	cs := ft.colorScale
	cs.ScaleAlpha(ft.opacity)

	tOpts.ColorScale = cs

	tOpts.GeoM.Translate(ft.x, ft.y)
	text.Draw(screen, ft.text, ft.face, tOpts)

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
