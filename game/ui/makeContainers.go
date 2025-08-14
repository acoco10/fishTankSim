package ui

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
	"image/color"
)

func MakeToolTipContainer(tipText []string) *widget.Container {
	tooltipContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionVertical))),
		widget.ContainerOpts.AutoDisableChildren(),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 230, A: 255})),
	)
	if len(tipText) > 0 {
		for _, text := range tipText {
			option1 := widget.NewText(
				widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
				widget.TextOpts.Text(fmt.Sprint(text), registry.FontMap["nk57"], colornames.Greenyellow),
				widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 0)),
			)
			tooltipContainer.AddChild(option1)
		}
	}
	return tooltipContainer
}

func LoadHeader(headerText string, face text.Face) *widget.Container {
	headerContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	headerLbl := widget.NewText(
		widget.TextOpts.Text(headerText, face, color.RGBA{R: 0, G: 160, B: 0, A: 255}),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.TextOpts.Insets(widget.Insets{
			Left:   30,
			Right:  10,
			Top:    50,
			Bottom: 600,
		}),
	)

	headerContainer.AddChild(headerLbl)

	return headerContainer
}
